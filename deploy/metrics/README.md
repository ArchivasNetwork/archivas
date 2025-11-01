# Archivas Metrics Gateway

Reference configuration for publishing public Prometheus metrics, scrape targets, and exporter health.

## Contents

- `gateway.nginx.conf` – Nginx reverse proxy exposing Prometheus `/prometheus/`, node autodiscovery `/targets.json`, and exporter health `/healthz`.

## Usage

1. Copy `gateway.nginx.conf` to `/etc/nginx/sites-available/archivas-metrics`.
2. Adjust upstream addresses if your Prometheus or node run on different hosts/ports.
3. Optionally enable basic authentication for `/prometheus/` or `/node-metrics` by uncommenting the sample location block and supplying an `.htpasswd` file.
4. Symlink into `sites-enabled` and reload Nginx:

   ```bash
   sudo ln -s /etc/nginx/sites-available/archivas-metrics /etc/nginx/sites-enabled/archivas-metrics
   sudo nginx -t
   sudo systemctl reload nginx
   ```

## Prometheus Auto-Discovery

The node now exposes `GET /metrics/targets.json`. Populate a Prometheus `file_sd_configs` entry with periodic fetch (e.g. `curl ... > /etc/prometheus/archivas-targets.json`). Targets include the local node and, when `?includePeers=true`, peer-derived endpoints using the provided `peerPort` (default `8080`).

Example `crontab` entry:

```cron
* * * * * curl -sf "http://seed.archivas.ai/targets.json?includePeers=true" -o /etc/prometheus/archivas-targets.json
```

Then reference the file in `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: archivas-autodiscover
    file_sd_configs:
      - files:
          - /etc/prometheus/archivas-targets.json
    relabel_configs:
      - source_labels: [__address__]
        regex: "(.*)"
        target_label: instance
        replacement: "$1"
```

## Exporter Health

`GET /metrics/health` returns an aggregated view of metric freshness backed by watchdog timers. The endpoint is suitable for uptime monitoring platforms (Statuscake, Pingdom, etc.). A non-`ok` status indicates at least one metric has not updated within its SLA window.


