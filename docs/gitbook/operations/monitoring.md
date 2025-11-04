# '$(basename $file .md | tr '-' ' ' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2)); print}')'

Detailed guide coming soon.

**For now, see:**
- [SEED_HOST.md](../../SEED_HOST.md) - Seed node setup
- [OBSERVABILITY.md](../../OBSERVABILITY.md) - Monitoring
- [CURRENT-STATUS.md](../current-status.md) - Current state
