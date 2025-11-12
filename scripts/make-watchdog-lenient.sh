#!/bin/bash

# Make watchdog much more lenient for seed node
# Instead of disabling, make it only restart after many failures

echo "Making watchdog more lenient for seed node..."
echo ""

# Update watchdog service to be more lenient
sudo sed -i 's/Environment="TIMEOUT=30"/Environment="TIMEOUT=60"/' /etc/systemd/system/archivas-watchdog.service
sudo sed -i 's/Environment="MAX_FAILURES=5"/Environment="MAX_FAILURES=20"/' /etc/systemd/system/archivas-watchdog.service
sudo sed -i 's/Environment="CHECK_INTERVAL=60"/Environment="CHECK_INTERVAL=120"/' /etc/systemd/system/archivas-watchdog.service

# Reload and restart
sudo systemctl daemon-reload
sudo systemctl restart archivas-watchdog.service

echo "✅ Watchdog updated:"
echo "  - Timeout: 60s (was 30s)"
echo "  - Max failures: 20 (was 5) = 40 minutes before restart"
echo "  - Check interval: 120s (was 60s)"
echo ""

