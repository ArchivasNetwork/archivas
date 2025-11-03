# Metrics After v1.1.1-infra Hardening

After the v1.1.1-infra security hardening, `/metrics` is **intentionally blocked** from public access. This document explains how to configure Prometheus to scrape metrics from localhost.

---

## Security Posture

**Public Endpoint (https://seed.archivas.ai):**
- ✅ `/metrics` returns **404** (blocked by Nginx)
- ✅ Internal metrics NOT exposed to internet
- ✅ Security hardening in place

**Local Monitoring:**
- ✅ Prometheus scrapes `127.0.0.1` ports
- ✅ Grafana connects to local Prometheus
- ✅ Metrics remain internal to the server

---

## Configuration

### 1. Ensure Exporters Bind to Localhost

The Archivas binaries expose metrics on their RPC ports by default:
- **archivas-node**: Port 8080 (same as RPC)
- **archivas-farmer**: Embedded in process (no separate metrics port)
- **archivas-timelord**: Embedded in process (no separate metrics port)

**All services already bind to `127.0.0.1` by default when using `--rpc 127.0.0.1:8080`.**

Verify listeners:
```bash
# Check what's listening:
sudo lsof -i :8080
# Should show: archivas-node listening on 127.0.0.1:8080

# Or use netstat:
sudo netstat -tlnp | grep -E ':(8080|9101|9102)'
```

---

### 2. Configure Prometheus for Localhost Scraping

**File:** `/etc/prometheus/prometheus.yml`

```yaml
global:
  scrape_interval: 15s
  scrape_timeout: 10s
  evaluation_interval: 15s

scrape_configs:
  # Archivas Node (RPC + Metrics on same port)
  - job_name: 'archivas-nodes'
    static_configs:
      - targets: ['127.0.0.1:8080']
        labels:
          instance: 'server-a'
          role: 'seed'

  # If you have farmers with separate metrics ports:
  # - job_name: 'archivas-farmers'
  #   static_configs:
  #     - targets: ['127.0.0.1:9102']

  # If you have timelords with separate metrics ports:
  # - job_name: 'archivas-timelords'
  #   static_configs:
  #     - targets: ['127.0.0.1:9101']
```

**Important:** Do NOT add `https://seed.archivas.ai/metrics` here. Prometheus must scrape localhost only.

Restart Prometheus:
```bash
sudo systemctl restart prometheus
```

Verify targets:
```bash
# Check Prometheus targets page:
curl http://localhost:9090/targets
# Or visit: http://localhost:9090/targets in browser
```

---

### 3. Verify Nginx Blocks /metrics (Intentional)

**File:** `deploy/seed/nginx-site.conf`

```nginx
# Block internal metrics from public access
location = /metrics {
  return 404;
}

location /metrics/ {
  return 404;
}
```

**This is correct and should NOT be changed.** Public access to `/metrics` is a security risk.

Test:
```bash
curl https://seed.archivas.ai/metrics
# Expected: HTTP/2 404
```

---

### 4. Grafana Verification

After Prometheus restart:

1. **Check Prometheus Targets:**
   - Open: `http://localhost:9090/targets`
   - Verify: `archivas-nodes` job shows **UP**

2. **Check Grafana Dashboard:**
   - Open: `http://localhost:3000` (or your Grafana URL)
   - Navigate to: **Archivas Network Overview**
   - Verify metrics appear within 30 seconds:
     - Tip Height
     - Peer Count
     - Difficulty
     - Blocks Mined

3. **If No Data:**
   - Increase dashboard time range to "Last 1 hour"
   - Check Prometheus data source health in Grafana
   - Verify Prometheus is scraping: `curl http://localhost:9090/api/v1/targets`

---

## Multi-Server Monitoring (Optional)

To collect metrics from Server C (or other nodes) **without exposing /metrics publicly**, use Prometheus **remote_write**.

### On Server C (Farmer Node)

**Install Prometheus Agent:**
```bash
sudo apt-get install -y prometheus
```

**Configure Agent Mode:**

File: `/etc/prometheus/prometheus.yml`

```yaml
global:
  scrape_interval: 15s
  external_labels:
    instance: 'server-c'
    role: 'farmer'

scrape_configs:
  # Local node metrics
  - job_name: 'archivas-nodes'
    static_configs:
      - targets: ['127.0.0.1:8080']

  # Local farmer metrics (if separate port)
  # - job_name: 'archivas-farmers'
  #   static_configs:
  #     - targets: ['127.0.0.1:9102']

# Forward metrics to Server A
remote_write:
  - url: http://57.129.148.132:9090/api/v1/write
    queue_config:
      capacity: 10000
      max_samples_per_send: 5000
      batch_send_deadline: 5s
```

**Restart Prometheus:**
```bash
sudo systemctl restart prometheus
```

### On Server A (Prometheus Server)

**Enable Remote Write Receiver:**

File: `/etc/prometheus/prometheus.yml`

Add to top-level config:
```yaml
# Enable receiving metrics from remote agents
storage:
  tsdb:
    out_of_order_time_window: 30m

# If using Prometheus 2.40+, enable remote write receiver
remote_write_receiver:
  enabled: true
```

**Firewall:**
```bash
# Allow Prometheus remote write from Server C
sudo ufw allow from 57.129.148.134 to any port 9090 proto tcp comment "Prometheus remote write"
```

**Restart Prometheus:**
```bash
sudo systemctl restart prometheus
```

---

## Troubleshooting

### "Connection Refused" Errors

**Symptom:** Prometheus shows targets as DOWN with "connection refused"

**Fix:**
1. **Verify service is running:**
   ```bash
   ps aux | grep archivas-node
   ```

2. **Check binding address:**
   ```bash
   sudo lsof -i :8080
   # Should show: 127.0.0.1:8080
   ```

3. **Check Prometheus targets config:**
   ```bash
   grep -A 5 "job_name:" /etc/prometheus/prometheus.yml
   # Should show: 127.0.0.1:8080 (NOT public IPs)
   ```

4. **Test scrape endpoint:**
   ```bash
   curl http://127.0.0.1:8080/metrics
   # Should show: Prometheus metrics output
   ```

---

### Grafana Shows "No Data"

**Fix:**
1. **Check Prometheus is scraping:**
   ```bash
   curl http://localhost:9090/api/v1/targets | jq
   ```

2. **Check metric names:**
   ```bash
   curl http://localhost:9090/api/v1/label/__name__/values | jq | grep archivas
   ```

3. **Test PromQL query:**
   ```bash
   curl 'http://localhost:9090/api/v1/query?query=archivas_tip_height'
   ```

4. **Increase dashboard time range:**
   - Grafana → Dashboard → Time Range → Last 1 hour

---

### Public /metrics Returns 404

**This is intentional and correct!**

```bash
curl https://seed.archivas.ai/metrics
# HTTP/2 404

curl -I https://seed.archivas.ai/metrics
# HTTP/2 404
# content-type: text/html
```

**Do NOT expose /metrics publicly.** This is a security feature.

**For monitoring:**
- Use localhost scraping (Prometheus on same server)
- Use remote_write (Prometheus agent forwarding)
- Use VPN/bastion for remote access to Prometheus

---

## Architecture

### Before v1.1.1-infra (INSECURE)
```
Internet → https://seed.archivas.ai/metrics → Full metrics exposed ❌
```

### After v1.1.1-infra (SECURE)
```
Internet → https://seed.archivas.ai/metrics → 404 (blocked) ✅

Localhost:
  Prometheus → http://127.0.0.1:8080/metrics → Metrics ✅
  Grafana → Prometheus → Dashboards ✅
```

### Multi-Server (SECURE)
```
Server C:
  Prometheus Agent → http://127.0.0.1:8080/metrics
    ↓ remote_write
Server A:
  Prometheus → receive + scrape local
    ↓
  Grafana → Dashboards
```

---

## Summary

✅ **Public /metrics blocked** - Security hardening in place  
✅ **Localhost scraping** - Prometheus scrapes 127.0.0.1  
✅ **Multi-server** - Use remote_write for distributed monitoring  
✅ **Grafana works** - Dashboards populate from local Prometheus  

**This is the correct and secure configuration for production monitoring.**

