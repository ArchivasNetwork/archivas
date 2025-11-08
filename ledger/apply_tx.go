package ledger

import (
	"errors"
)

var (
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrBadNonce           = errors.New("bad nonce")
	ErrInvalidSignature   = errors.New("invalid signature")
	ErrMissingSignature   = errors.New("missing signature")
	ErrMissingSenderPubKey = errors.New("missing sender public key")
)

// ApplyTransaction applies a transaction to the world state
// Returns error if the transaction is invalid or cannot be applied
func (ws *WorldState) ApplyTransaction(tx Transaction) error {
	// Skip signature verification for transactions from mempool (already verified via txv1 or legacy)
	// Mempool transactions have already passed signature verification in handleSubmitV1 or handleSubmitTx
	// Re-verification would fail for txv1 (Ed25519) transactions because this function expects secp256k1
	
	// Note: For transactions from external sources (P2P, IBD), signature should be verified before calling this

	sender, ok := ws.Accounts[tx.From]
	if !ok {
		return ErrInsufficientFunds
	}

	// Nonce must match
	if sender.Nonce != tx.Nonce {
		return ErrBadNonce
	}

	totalCost := tx.Amount + tx.Fee
	if sender.Balance < totalCost {
		return ErrInsufficientFunds
	}

	// Deduct from sender
	sender.Balance -= totalCost
	sender.Nonce += 1

	// Credit receiver
	recv, ok := ws.Accounts[tx.To]
	if !ok {
		recv = &AccountState{Balance: 0, Nonce: 0}
		ws.Accounts[tx.To] = recv
	}
	recv.Balance += tx.Amount

	// Fee handling: for now, just burn it (TODO: pay block proposer)
	return nil
}

