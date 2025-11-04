# Troubleshooting

Common farmer issues and solutions.

## Farmer Issues

### "No plots found"
**Solution:** Check `--plots` path matches plot directory

### "Connection refused"
**Solution:** Verify node is running: `curl http://localhost:8080/chainTip`

### Not winning blocks
**Normal if:** Network has much more space than you
**Check:** Are plots loading? Is farmer scanning?

## Node Issues

### "Database lock"
**Solution:** `pkill -9 archivas-node` then restart

### Sync stuck
**Solution:** Check logs, verify genesis hash, restart node

### Out of disk
**Solution:** Free space or increase disk size

---

For more help: [GitHub Discussions](https://github.com/ArchivasNetwork/archivas/discussions)
