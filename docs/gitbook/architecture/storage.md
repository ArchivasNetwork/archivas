# Storage Engine

Archivas uses BadgerDB for persistence.

## Data Stored

- Blocks (by height)
- Account state (balance, nonce)
- Metadata (tip, difficulty)
- Peer list

## Size

Current: ~2 GB for 64,000 blocks
Growth: ~30 MB per 1,000 blocks

## Backup

```bash
tar -czf backup.tar.gz data/
```
