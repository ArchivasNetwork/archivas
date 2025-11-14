#!/bin/bash
set -euo pipefail

# Firewall configuration for Seed2 (Full Node + Relay)
# Run as: sudo bash infra/firewall.seed2.sh

echo "=== Configuring firewall for Seed2 (Full Node + Relay) ==="

# Check if UFW is installed
if ! command -v ufw &> /dev/null; then
    echo "Installing UFW..."
    apt-get update && apt-get install -y ufw
fi

# Reset UFW to default deny
ufw --force reset
ufw default deny incoming
ufw default allow outgoing

# Allow SSH (critical - don't lock yourself out!)
ufw allow 22/tcp comment 'SSH'

# HTTP/HTTPS for RPC relay (public)
ufw allow 80/tcp comment 'HTTP Relay'
ufw allow 443/tcp comment 'HTTPS Relay'

# P2P port for full node (public - farmers need this)
ufw allow 30303/tcp comment 'Archivas P2P TCP'
ufw allow 30303/udp comment 'Archivas P2P UDP'

# Metrics (restricted to trusted IPs only)
# Replace with your monitoring server IP
# ufw allow from <MONITORING_IP> to any port 9102 proto tcp comment 'Prometheus metrics'
# ufw allow from <MONITORING_IP> to any port 9090 proto tcp comment 'Relay service metrics'

# Node RPC (127.0.0.1:8082) - already bound to localhost, no firewall rule needed

# Enable UFW
ufw --force enable

# Show status
echo ""
echo "=== Firewall rules applied ==="
ufw status numbered

echo ""
echo "IMPORTANT: Verify SSH (port 22) is allowed before closing this terminal!"
echo "If you get locked out, use your hosting provider's console to disable UFW."

