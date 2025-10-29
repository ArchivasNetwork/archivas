package ledger

// Transaction represents a basic "send RCHV" transaction
type Transaction struct {
	From         string // bech32 address
	To           string // bech32 address
	Amount       int64  // base units
	Fee          int64  // base units
	Nonce        uint64 // must match sender's current nonce
	SenderPubKey []byte // sender's public key (used to verify signature and derive From address)
	Signature    []byte // secp256k1 signature over HashTransaction(tx)
}

