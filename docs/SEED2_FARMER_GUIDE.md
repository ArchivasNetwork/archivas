# Seed2 Farmer Guide

## 🎉 Seed2 is Now Live!

Seed2 (`seed2.archivas.ai`) is now fully synced and ready for farmers!

## ✅ Seed2 Status
- **Height**: Fully synced with Seed1
- **RPC**: ✅ Available at `https://seed2.archivas.ai`
- **P2P**: ✅ Available at `seed2.archivas.ai:30303`
- **Uptime**: 24/7 monitoring
- **Location**: High-availability server

## For Farmers

### Using Seed2 for Farming (Recommended)

Farmers can now use Seed2's RPC endpoint to reduce load on Seed1:

```bash
# Windows (PowerShell)
.\archivas-farmer.exe farm `
  --plots .\plots `
  --node https://seed2.archivas.ai `
  --farmer-privkey YOUR_PRIVATE_KEY

# Linux/Mac
./archivas-farmer farm \
  --plots ./plots \
  --node https://seed2.archivas.ai \
  --farmer-privkey YOUR_PRIVATE_KEY
```

### Benefits of Using Seed2
- 🚀 Reduces load on Seed1
- ⚡ Fast response times
- 🛡️ Redundancy (if Seed1 has issues)
- 🌐 Globally distributed

## For Node Operators

If you're running your own Archivas node, you can peer with Seed2:

```bash
./archivas-node \
  --data-dir ./data \
  --p2p 0.0.0.0:9090 \
  --peer seed.archivas.ai:9090 \
  --peer seed2.archivas.ai:30303 \
  --rpc 127.0.0.1:8080
```

**Note**: Seed2 uses port `30303` for P2P, while Seed1 uses port `9090`.

## Dual-Seed Setup (Maximum Reliability)

For maximum reliability, configure your farmer to fail over between seeds:

```bash
# Primary: Seed2 (fast, cached)
# Fallback: Seed1 (main seed)

# Start with Seed2
./archivas-farmer farm \
  --plots ./plots \
  --node https://seed2.archivas.ai \
  --farmer-privkey YOUR_KEY

# If Seed2 is down, restart with Seed1
./archivas-farmer farm \
  --plots ./plots \
  --node https://seed.archivas.ai \
  --farmer-privkey YOUR_KEY
```

## RPC Endpoints

Both seeds support the same endpoints:

- `GET /challenge` - Get current farming challenge
- `POST /submit` - Submit a block proof
- `GET /chainTip` - Get current chain height
- `GET /sync/status` - Check sync status
- `GET /account/{address}` - Get account balance

## Troubleshooting

### Cannot connect to Seed2
1. Check your internet connection
2. Ensure you're using HTTPS: `https://seed2.archivas.ai`
3. Try Seed1 as fallback: `https://seed.archivas.ai`

### Slow response times
1. Check Seed2 status at `/sync/status`
2. Switch to Seed1 temporarily
3. Report persistent issues to the team

### P2P connection refused
1. Ensure you're using port 30303: `seed2.archivas.ai:30303`
2. Check your firewall allows outbound connections
3. Verify your node is configured correctly

## Support

- **Documentation**: https://docs.archivas.ai
- **Issues**: https://github.com/ArchivasNetwork/archivas/issues
- **Discord**: [Community Server]

---

**Last Updated**: 2025-11-14  
**Seed2 Version**: v1.1.1-ibd  
**Chain**: Archivas Devnet V4

