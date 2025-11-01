package v1

import (
	"encoding/json"
	"fmt"
)

// PackSignedTx creates a SignedTx from transaction, signature, and public key.
func PackSignedTx(tx *Transfer, privKey []byte) (*SignedTx, error) {
	sig, pubKey, hash, err := Sign(privKey, tx)
	if err != nil {
		return nil, err
	}

	return &SignedTx{
		Tx:     tx,
		PubKey: EncodePubKey(pubKey),
		Sig:    EncodeSig(sig),
		Hash:   EncodeHash(hash),
	}, nil
}

// UnpackSignedTx decodes a SignedTx from JSON bytes.
func UnpackSignedTx(data []byte) (*SignedTx, error) {
	var stx SignedTx
	if err := json.Unmarshal(data, &stx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal signed transaction: %w", err)
	}

	// Validate the transaction
	if err := stx.Tx.Validate(); err != nil {
		return nil, fmt.Errorf("invalid transaction in signed tx: %w", err)
	}

	// Validate encoding formats
	if _, err := DecodePubKey(stx.PubKey); err != nil {
		return nil, fmt.Errorf("invalid public key encoding: %w", err)
	}

	if _, err := DecodeSig(stx.Sig); err != nil {
		return nil, fmt.Errorf("invalid signature encoding: %w", err)
	}

	if _, err := DecodeHash(stx.Hash); err != nil {
		return nil, fmt.Errorf("invalid hash encoding: %w", err)
	}

	return &stx, nil
}

// MarshalJSON marshals a SignedTx to JSON.
func (stx *SignedTx) MarshalJSON() ([]byte, error) {
	type alias SignedTx
	return json.Marshal((*alias)(stx))
}

// UnmarshalJSON unmarshals a SignedTx from JSON.
func (stx *SignedTx) UnmarshalJSON(data []byte) error {
	type alias SignedTx
	if err := json.Unmarshal(data, (*alias)(stx)); err != nil {
		return err
	}

	// Validate on unmarshal
	if err := stx.Tx.Validate(); err != nil {
		return fmt.Errorf("invalid transaction: %w", err)
	}

	return nil
}
