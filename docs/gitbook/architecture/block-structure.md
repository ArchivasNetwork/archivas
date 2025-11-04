# Block Structure

Archivas block format.

## Fields

- height: Block number
- hash: Block hash (SHA256)
- prevHash: Previous block hash
- timestamp: Unix timestamp
- difficulty: Target difficulty
- challenge: PoSpace challenge
- farmer: Block producer address
- txs: Array of transactions
- proof: PoSpace proof

See existing blocks: `curl https://seed.archivas.ai/block/64000`
