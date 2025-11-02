#!/usr/bin/env bash
# renew-cert-hook.sh - Post-renewal hook for certbot
# Installed in: /etc/letsencrypt/renewal-hooks/post/renew-cert-hook.sh
# Usage: certbot runs this automatically after certificate renewal

set -euo pipefail

echo "🔄 Certificate renewed, reloading Nginx..."

# Reload Nginx to pick up new certificates
if systemctl is-active --quiet nginx; then
    systemctl reload nginx
    echo "✅ Nginx reloaded successfully"
else
    echo "⚠️  Nginx not running, skipping reload"
fi

# Optional: Log renewal timestamp
echo "$(date -u +%Y-%m-%dT%H:%M:%SZ) - Certificate renewed" >> /var/log/archivas-cert-renewal.log

