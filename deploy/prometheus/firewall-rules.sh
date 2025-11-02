#!/bin/bash
# v1.1.1: Firewall rules for Prometheus scraping
# Run this on each target host (57.129.148.132 and 72.251.11.191)

PROM_IP="57.129.148.132"  # Prometheus server IP

echo "🔒 Setting up firewall rules for Prometheus scraping..."

# Allow Prometheus to scrape metrics from nodes (port 8080)
sudo ufw allow from $PROM_IP to any port 8080 proto tcp comment "Prometheus: archivas-nodes"

# Allow Prometheus to scrape metrics from timelords (port 9101)
sudo ufw allow from $PROM_IP to any port 9101 proto tcp comment "Prometheus: archivas-timelords"

# Allow Prometheus to scrape metrics from farmers (port 9102)
sudo ufw allow from $PROM_IP to any port 9102 proto tcp comment "Prometheus: archivas-farmers"

sudo ufw reload

echo "✅ Firewall rules added"
echo ""
echo "Verify rules:"
sudo ufw status numbered | grep -E "8080|9101|9102"

