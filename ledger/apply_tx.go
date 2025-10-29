package ledger

import (
	"errors"
	"fmt"
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
	// Verify signature is present
	if len(tx.Signature) == 0 {
		return ErrMissingSignature
	}

	if len(tx.SenderPubKey) == 0 {
		return ErrMissingSenderPubKey
	}

	// Verify transaction signature (checks that SenderPubKey matches From and signature is valid)
	if err := VerifyTransactionSignature(tx); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidSignature, err)
	}

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

