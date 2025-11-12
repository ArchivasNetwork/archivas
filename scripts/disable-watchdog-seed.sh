#!/bin/bash

# Disable watchdog for seed node (Server A)
# The seed node is under heavy load and watchdog causes unnecessary restarts

echo "Disabling watchdog for seed node..."
echo ""

# Stop and disable watchdog
sudo systemctl stop archivas-watchdog.service
sudo systemctl disable archivas-watchdog.service

echo "✅ Watchdog stopped and disabled"
echo ""
echo "To re-enable later:"
echo "  sudo systemctl enable archivas-watchdog.service"
echo "  sudo systemctl start archivas-watchdog.service"
echo ""

