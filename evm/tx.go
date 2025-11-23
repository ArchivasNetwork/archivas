package evm

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ArchivasNetwork/archivas/address"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// EvmTx represents an Ethereum-compatible transaction
// This is the native format used by MetaMask and other Ethereum tools
type EvmTx struct {
	// Transaction type (0=legacy, 1=EIP-2930, 2=EIP-1559)
	TxType uint8
	
	// Core fields
	ChainID  *big.Int        // Chain ID (for EIP-155 and EIP-1559)
	Nonce    uint64
	GasLimit uint64
	To       *common.Address // nil for contract creation
	Value    *big.Int
	Data     []byte
	
	// Gas pricing
	GasPrice       *big.Int // For legacy and EIP-2930
	GasFeeCap      *big.Int // For EIP-1559 (maxFeePerGas)
	GasTipCap      *big.Int // For EIP-1559 (maxPriorityFeePerGas)
	
	// Signature (ECDSA)
	V *big.Int // Recovery ID + chain ID (EIP-155)
	R *big.Int // Signature R component
	S *big.Int // Signature S component
	
	// Derived fields (computed, not serialized)
	Hash   common.Hash    // Transaction hash
	From   common.Address // Sender address (recovered from signature)
}

// FromRawTransaction decodes a raw Ethereum transaction and recovers the sender
// This is the canonical way to process transactions from MetaMask
func FromRawTransaction(rawTx []byte, chainID *big.Int) (*EvmTx, error) {
	// Use go-ethereum's transaction decoder
	gethTx := new(types.Transaction)
	if err := gethTx.UnmarshalBinary(rawTx); err != nil {
		return nil, fmt.Errorf("failed to decode RLP transaction: %w", err)
	}
	
	// Create our EvmTx from the go-ethereum transaction
	evmTx := &EvmTx{
		TxType:   gethTx.Type(),
		ChainID:  gethTx.ChainId(),
		Nonce:    gethTx.Nonce(),
		GasLimit: gethTx.Gas(),
		Value:    gethTx.Value(),
		Data:     gethTx.Data(),
		Hash:     gethTx.Hash(),
	}
	
	// Handle 'To' address (nil for contract creation)
	if gethTx.To() != nil {
		to := *gethTx.To()
		evmTx.To = &to
	}
	
	// Extract signature components
	v, r, s := gethTx.RawSignatureValues()
	evmTx.V = v
	evmTx.R = r
	evmTx.S = s
	
	// Handle different transaction types for gas pricing
	switch gethTx.Type() {
	case types.LegacyTxType:
		evmTx.GasPrice = gethTx.GasPrice()
	case types.AccessListTxType:
		evmTx.GasPrice = gethTx.GasPrice()
	case types.DynamicFeeTxType:
		evmTx.GasFeeCap = gethTx.GasFeeCap()
		evmTx.GasTipCap = gethTx.GasTipCap()
	default:
		return nil, fmt.Errorf("unsupported transaction type: %d", gethTx.Type())
	}
	
	// Recover sender address using Ethereum signature verification
	// This is the ONLY correct way to get the sender for Ethereum transactions
	if err := evmTx.RecoverSender(chainID); err != nil {
		return nil, fmt.Errorf("failed to recover sender: %w", err)
	}
	
	return evmTx, nil
}

// RecoverSender recovers the sender address from the transaction signature
// This implements Ethereum's secp256k1 signature recovery
func (tx *EvmTx) RecoverSender(chainID *big.Int) error {
	// Reconstruct go-ethereum transaction for signature recovery
	gethTx := tx.toGethTransaction()
	
	// Get the appropriate signer for the chain ID and transaction type
	// Use LatestSignerForChainID which auto-detects the right signer
	var signer types.Signer
	if chainID != nil {
		signer = types.LatestSignerForChainID(chainID)
	} else {
		signer = types.HomesteadSigner{}
	}
	
	// Recover sender address
	from, err := types.Sender(signer, gethTx)
	if err != nil {
		return fmt.Errorf("signature recovery failed: %w (wrong chainID or invalid signature)", err)
	}
	
	tx.From = from
	return nil
}

// toGethTransaction converts our EvmTx back to go-ethereum's Transaction
// This is needed for signature recovery
func (tx *EvmTx) toGethTransaction() *types.Transaction {
	var gethTx *types.Transaction
	
	switch tx.TxType {
	case types.LegacyTxType:
		gethTx = types.NewTx(&types.LegacyTx{
			Nonce:    tx.Nonce,
			GasPrice: tx.GasPrice,
			Gas:      tx.GasLimit,
			To:       tx.To,
			Value:    tx.Value,
			Data:     tx.Data,
			V:        tx.V,
			R:        tx.R,
			S:        tx.S,
		})
	case types.AccessListTxType:
		gethTx = types.NewTx(&types.AccessListTx{
			Nonce:    tx.Nonce,
			GasPrice: tx.GasPrice,
			Gas:      tx.GasLimit,
			To:       tx.To,
			Value:    tx.Value,
			Data:     tx.Data,
			V:        tx.V,
			R:        tx.R,
			S:        tx.S,
		})
	case types.DynamicFeeTxType:
		gethTx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   tx.ChainID,
			Nonce:     tx.Nonce,
			GasTipCap: tx.GasTipCap,
			GasFeeCap: tx.GasFeeCap,
			Gas:       tx.GasLimit,
			To:        tx.To,
			Value:     tx.Value,
			Data:      tx.Data,
			V:         tx.V,
			R:         tx.R,
			S:         tx.S,
		})
	}
	
	return gethTx
}

// RecoverPublicKey recovers the full public key from the transaction signature
// This is needed for advanced cryptographic operations
func (tx *EvmTx) RecoverPublicKey(chainID *big.Int) (*ecdsa.PublicKey, error) {
	// Build the signature in Ethereum format: R(32) + S(32) + V(1)
	signature := make([]byte, 65)
	
	// Pad R and S to 32 bytes
	rBytes := tx.R.Bytes()
	sBytes := tx.S.Bytes()
	copy(signature[32-len(rBytes):32], rBytes)
	copy(signature[64-len(sBytes):64], sBytes)
	
	// Normalize V to recovery ID (0 or 1)
	// Ethereum uses v = 27 + recoveryID (legacy) or v = chainID*2 + 35 + recoveryID (EIP-155)
	var recoveryID byte
	if tx.V.Cmp(big.NewInt(35)) >= 0 {
		// EIP-155: v = chainID*2 + 35 + recoveryID
		v := new(big.Int).Set(tx.V)
		v.Sub(v, big.NewInt(35))
		if chainID != nil {
			chainIDMul := new(big.Int).Mul(chainID, big.NewInt(2))
			v.Sub(v, chainIDMul)
		}
		recoveryID = byte(v.Uint64())
	} else {
		// Legacy: v = 27 + recoveryID
		recoveryID = byte(tx.V.Uint64() - 27)
	}
	signature[64] = recoveryID
	
	// Get the transaction hash for signature verification
	gethTx := tx.toGethTransaction()
	var signer types.Signer
	if chainID != nil {
		signer = types.LatestSignerForChainID(chainID)
	} else {
		signer = types.HomesteadSigner{}
	}
	txHash := signer.Hash(gethTx)
	
	// Recover the public key
	pubKey, err := crypto.SigToPub(txHash.Bytes(), signature)
	if err != nil {
		return nil, fmt.Errorf("failed to recover public key: %w", err)
	}
	
	return pubKey, nil
}

// ToARCVAddress converts the Ethereum address to ARCV Bech32 format
// This maintains compatibility with Archivas native addressing
func (tx *EvmTx) ToARCVAddress() (string, error) {
	evmAddr := address.EVMAddress(tx.From)
	return address.EncodeARCVAddress(evmAddr, "arcv")
}

// IsContractCreation returns true if this transaction creates a contract
func (tx *EvmTx) IsContractCreation() bool {
	return tx.To == nil
}

// EffectiveGasPrice returns the effective gas price for this transaction
// For EIP-1559, this depends on the base fee, but for now we use maxFeePerGas
func (tx *EvmTx) EffectiveGasPrice(baseFee *big.Int) *big.Int {
	if tx.TxType == types.DynamicFeeTxType {
		// EIP-1559: effectiveGasPrice = min(maxFeePerGas, baseFee + maxPriorityFeePerGas)
		if baseFee == nil {
			baseFee = big.NewInt(0)
		}
		
		// Calculate baseFee + maxPriorityFeePerGas
		effectivePrice := new(big.Int).Add(baseFee, tx.GasTipCap)
		
		// Cap at maxFeePerGas
		if effectivePrice.Cmp(tx.GasFeeCap) > 0 {
			effectivePrice = tx.GasFeeCap
		}
		
		return effectivePrice
	}
	
	// Legacy or EIP-2930: just use gasPrice
	return tx.GasPrice
}

// Cost returns the total upfront cost of the transaction (value + gas)
func (tx *EvmTx) Cost(baseFee *big.Int) *big.Int {
	gasPrice := tx.EffectiveGasPrice(baseFee)
	gasCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(tx.GasLimit)))
	total := new(big.Int).Add(tx.Value, gasCost)
	return total
}

// Validate performs basic validation on the transaction
func (tx *EvmTx) Validate(chainID *big.Int) error {
	// Check signature is present
	if tx.V == nil || tx.R == nil || tx.S == nil {
		return fmt.Errorf("missing signature")
	}
	
	// Check value is non-negative
	if tx.Value == nil || tx.Value.Sign() < 0 {
		return fmt.Errorf("negative value")
	}
	
	// Check gas limit is reasonable
	if tx.GasLimit == 0 {
		return fmt.Errorf("zero gas limit")
	}
	
	// For EIP-1559, check gas caps
	if tx.TxType == types.DynamicFeeTxType {
		if tx.GasFeeCap == nil || tx.GasTipCap == nil {
			return fmt.Errorf("EIP-1559 transaction missing gas caps")
		}
		if tx.GasTipCap.Cmp(tx.GasFeeCap) > 0 {
			return fmt.Errorf("maxPriorityFeePerGas exceeds maxFeePerGas")
		}
	}
	
	// Verify sender can be recovered
	if tx.From == (common.Address{}) {
		if err := tx.RecoverSender(chainID); err != nil {
			return fmt.Errorf("cannot recover sender: %w", err)
		}
	}
	
	return nil
}

// String returns a human-readable representation of the transaction
func (tx *EvmTx) String() string {
	txType := "legacy"
	if tx.TxType == types.AccessListTxType {
		txType = "EIP-2930"
	} else if tx.TxType == types.DynamicFeeTxType {
		txType = "EIP-1559"
	}
	
	to := "contract creation"
	if tx.To != nil {
		to = tx.To.Hex()
	}
	
	return fmt.Sprintf("EvmTx{type=%s, from=%s, to=%s, nonce=%d, value=%s, gas=%d, hash=%s}",
		txType, tx.From.Hex(), to, tx.Nonce, tx.Value, tx.GasLimit, tx.Hash.Hex())
}

