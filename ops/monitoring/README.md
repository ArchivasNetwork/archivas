# Archivas Monitoring Stack

Prometheus + Grafana for Archivas blockchain monitoring.

## Quick Start

```bash
cd ops/monitoring
docker-compose up -d
```

## Access

- **Prometheus:** http://localhost:9090
- **Grafana:** http://localhost:3000
  - Username: `admin`
  - Password: `archivas`

## Dashboards

Pre-loaded dashboard: "Archivas Network Overview"

Shows:
- Tip height per node
- Connected peers
- Difficulty
- Blocks mined
- RPC request rates

## Configuration

### Add/Remove Nodes

Edit `prometheus.yml`:
```yaml
scrape_configs:
  - job_name: 'archivas-nodes'
    static_configs:
      - targets:
          - 'your-node:8080'  # Add here
```

Reload: `docker-compose restart prometheus`

### Custom Dashboards

1. Go to Grafana (http://localhost:3000)
2. Create → Dashboard
3. Add panels with Archivas metrics
4. Export JSON → save to `dashboards/`

## Metrics Reference

### Node Metrics (port 8080)
- `archivas_tip_height` - Current blockchain height
- `archivas_peer_count` - Connected peers
- `archivas_peers_known` - Total known peers
- `archivas_blocks_total` - Total blocks processed
- `archivas_difficulty` - Current mining difficulty
- `archivas_rpc_requests_total{endpoint}` - RPC requests by endpoint

### Timelord Metrics (port 9101)
- `archivas_vdf_iterations_total` - Total VDF computations
- `archivas_vdf_tick_duration_seconds` - VDF computation time
- `archivas_vdf_seed_changed_total` - Seed changes (new blocks)

### Farmer Metrics (port 9102)
- `archivas_plots_loaded` - Number of plot files
- `archivas_qualities_checked_total` - Total quality checks
- `archivas_blocks_won_total` - Blocks won
- `archivas_farming_loop_seconds` - Farming iteration time

## Troubleshooting

**Prometheus not scraping:**
- Check targets at http://localhost:9090/targets
- Verify node/timelord/farmer have /metrics exposed
- Check firewall rules

**Grafana not showing data:**
- Verify datasource at http://localhost:3000/datasources
- Check Prometheus is accessible from Grafana container
- Test query in Explore tab

## Production Deployment

For production, use:
- Persistent volumes for data
- Reverse proxy (nginx) for authentication
- Alertmanager for notifications

See `../../OPERATIONS.md` for details.

