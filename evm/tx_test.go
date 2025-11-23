package evm

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Test vector: Known Ethereum transaction
func TestFromRawTransaction_Legacy(t *testing.T) {
	// This is a real Ethereum legacy transaction
	// Sender: 0x39a028dfdcae40bf277ec1ec268d62665d36c073
	// To: 0x676979f6aca633dadabdad81f9ef0310f0eae39e
	// Value: 1000000000 (10 RCHV in base units)
	// Nonce: 0
	// GasPrice: 1 gwei
	// GasLimit: 21000
	// ChainID: 1644
	
	privKeyHex := "fe4e0b573e892b9abc7692782aa14bc7560e6e7637948403827a6242dcf0f2b9"
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		t.Fatalf("Failed to load private key: %v", err)
	}
	
	// Create a legacy transaction
	chainID := big.NewInt(1644)
	signer := types.NewEIP155Signer(chainID)
	
	tx := types.NewTransaction(
		0,                                                         // nonce
		common.HexToAddress("0x676979f6aca633dadabdad81f9ef0310f0eae39e"), // to
		big.NewInt(1000000000),                                    // value
		21000,                                                     // gas limit
		big.NewInt(1000000000),                                    // gas price (1 gwei)
		nil,                                                       // data
	)
	
	// Sign the transaction
	signedTx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}
	
	// Encode to raw bytes
	rawTx, err := signedTx.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal transaction: %v", err)
	}
	
	// Test our decoder
	evmTx, err := FromRawTransaction(rawTx, chainID)
	if err != nil {
		t.Fatalf("FromRawTransaction failed: %v", err)
	}
	
	// Verify fields
	expectedFrom := common.HexToAddress("0x39a028dfdcae40bf277ec1ec268d62665d36c073")
	if evmTx.From != expectedFrom {
		t.Errorf("From mismatch: got %s, want %s", evmTx.From.Hex(), expectedFrom.Hex())
	}
	
	if evmTx.Nonce != 0 {
		t.Errorf("Nonce mismatch: got %d, want 0", evmTx.Nonce)
	}
	
	if evmTx.Value.Cmp(big.NewInt(1000000000)) != 0 {
		t.Errorf("Value mismatch: got %s, want 1000000000", evmTx.Value)
	}
	
	if evmTx.GasLimit != 21000 {
		t.Errorf("GasLimit mismatch: got %d, want 21000", evmTx.GasLimit)
	}
	
	if evmTx.TxType != types.LegacyTxType {
		t.Errorf("TxType mismatch: got %d, want %d (legacy)", evmTx.TxType, types.LegacyTxType)
	}
	
	t.Logf("✅ Successfully decoded legacy transaction")
	t.Logf("   From: %s", evmTx.From.Hex())
	t.Logf("   To: %s", evmTx.To.Hex())
	t.Logf("   Hash: %s", evmTx.Hash.Hex())
}

// Test EIP-1559 transaction
func TestFromRawTransaction_EIP1559(t *testing.T) {
	privKeyHex := "fe4e0b573e892b9abc7692782aa14bc7560e6e7637948403827a6242dcf0f2b9"
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		t.Fatalf("Failed to load private key: %v", err)
	}
	
	chainID := big.NewInt(1644)
	signer := types.NewLondonSigner(chainID)
	
	// Create EIP-1559 transaction
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     1,
		GasTipCap: big.NewInt(2000000000), // 2 gwei priority fee
		GasFeeCap: big.NewInt(3000000000), // 3 gwei max fee
		Gas:       21000,
		To:        &common.Address{0x67, 0x69, 0x79, 0xf6},
		Value:     big.NewInt(5000000000),
		Data:      nil,
	})
	
	signedTx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		t.Fatalf("Failed to sign EIP-1559 transaction: %v", err)
	}
	
	rawTx, err := signedTx.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal transaction: %v", err)
	}
	
	// Decode
	evmTx, err := FromRawTransaction(rawTx, chainID)
	if err != nil {
		t.Fatalf("FromRawTransaction failed for EIP-1559: %v", err)
	}
	
	// Verify EIP-1559 specific fields
	if evmTx.TxType != types.DynamicFeeTxType {
		t.Errorf("TxType mismatch: got %d, want %d (EIP-1559)", evmTx.TxType, types.DynamicFeeTxType)
	}
	
	if evmTx.GasFeeCap.Cmp(big.NewInt(3000000000)) != 0 {
		t.Errorf("GasFeeCap mismatch: got %s, want 3000000000", evmTx.GasFeeCap)
	}
	
	if evmTx.GasTipCap.Cmp(big.NewInt(2000000000)) != 0 {
		t.Errorf("GasTipCap mismatch: got %s, want 2000000000", evmTx.GasTipCap)
	}
	
	t.Logf("✅ Successfully decoded EIP-1559 transaction")
	t.Logf("   MaxFeePerGas: %s", evmTx.GasFeeCap)
	t.Logf("   MaxPriorityFeePerGas: %s", evmTx.GasTipCap)
}

// Test contract creation transaction
func TestFromRawTransaction_ContractCreation(t *testing.T) {
	privKeyHex := "fe4e0b573e892b9abc7692782aa14bc7560e6e7637948403827a6242dcf0f2b9"
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		t.Fatalf("Failed to load private key: %v", err)
	}
	
	chainID := big.NewInt(1644)
	signer := types.NewLondonSigner(chainID)
	
	// Simple contract bytecode (just returns 42)
	contractCode, _ := hex.DecodeString("6042600052600160206000f3")
	
	// Contract creation (To = nil)
	tx := types.NewContractCreation(
		0,                      // nonce
		big.NewInt(0),          // value
		100000,                 // gas limit
		big.NewInt(1000000000), // gas price
		contractCode,           // contract bytecode
	)
	
	signedTx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		t.Fatalf("Failed to sign contract creation: %v", err)
	}
	
	rawTx, err := signedTx.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal transaction: %v", err)
	}
	
	// Decode
	evmTx, err := FromRawTransaction(rawTx, chainID)
	if err != nil {
		t.Fatalf("FromRawTransaction failed for contract creation: %v", err)
	}
	
	// Verify it's a contract creation
	if !evmTx.IsContractCreation() {
		t.Errorf("Expected contract creation (To=nil), got To=%v", evmTx.To)
	}
	
	if len(evmTx.Data) != len(contractCode) {
		t.Errorf("Contract code length mismatch: got %d, want %d", len(evmTx.Data), len(contractCode))
	}
	
	t.Logf("✅ Successfully decoded contract creation transaction")
	t.Logf("   Contract code: %x", evmTx.Data)
}

// Test signature recovery with wrong chain ID
func TestFromRawTransaction_WrongChainID(t *testing.T) {
	privKeyHex := "fe4e0b573e892b9abc7692782aa14bc7560e6e7637948403827a6242dcf0f2b9"
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		t.Fatalf("Failed to load private key: %v", err)
	}
	
	// Sign with chain ID 1644
	chainID := big.NewInt(1644)
	signer := types.NewEIP155Signer(chainID)
	
	tx := types.NewTransaction(
		0,
		common.HexToAddress("0x676979f6aca633dadabdad81f9ef0310f0eae39e"),
		big.NewInt(1000000000),
		21000,
		big.NewInt(1000000000),
		nil,
	)
	
	signedTx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}
	
	rawTx, err := signedTx.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal transaction: %v", err)
	}
	
	// Try to decode with WRONG chain ID (should fail)
	wrongChainID := big.NewInt(1) // Ethereum mainnet
	_, err = FromRawTransaction(rawTx, wrongChainID)
	if err == nil {
		t.Errorf("Expected error when using wrong chain ID, got nil")
	} else {
		t.Logf("✅ Correctly rejected transaction with wrong chain ID: %v", err)
	}
}

// Test ARCV address conversion
func TestToARCVAddress(t *testing.T) {
	privKeyHex := "fe4e0b573e892b9abc7692782aa14bc7560e6e7637948403827a6242dcf0f2b9"
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		t.Fatalf("Failed to load private key: %v", err)
	}
	
	chainID := big.NewInt(1644)
	signer := types.NewLondonSigner(chainID)
	
	tx := types.NewTransaction(
		0,
		common.HexToAddress("0x676979f6aca633dadabdad81f9ef0310f0eae39e"),
		big.NewInt(1000000000),
		21000,
		big.NewInt(1000000000),
		nil,
	)
	
	signedTx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}
	
	rawTx, err := signedTx.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal transaction: %v", err)
	}
	
	evmTx, err := FromRawTransaction(rawTx, chainID)
	if err != nil {
		t.Fatalf("FromRawTransaction failed: %v", err)
	}
	
	// Convert to ARCV address
	arcvAddr, err := evmTx.ToARCVAddress()
	if err != nil {
		t.Fatalf("ToARCVAddress failed: %v", err)
	}
	
	// Should match the known ARCV address for this private key
	expectedARCV := "arcv18xsz3h7u4eqt7fm7c8kzdrtzvewndsrn2jd8su"
	if arcvAddr != expectedARCV {
		t.Errorf("ARCV address mismatch: got %s, want %s", arcvAddr, expectedARCV)
	}
	
	t.Logf("✅ Successfully converted EVM address to ARCV")
	t.Logf("   EVM:  %s", evmTx.From.Hex())
	t.Logf("   ARCV: %s", arcvAddr)
}

// Test effective gas price calculation
func TestEffectiveGasPrice(t *testing.T) {
	// Test EIP-1559 effective gas price
	evmTx := &EvmTx{
		TxType:    types.DynamicFeeTxType,
		GasFeeCap: big.NewInt(100), // max 100 gwei
		GasTipCap: big.NewInt(10),  // priority 10 gwei
	}
	
	// With baseFee = 50, effective = min(100, 50+10) = 60
	baseFee := big.NewInt(50)
	effectivePrice := evmTx.EffectiveGasPrice(baseFee)
	if effectivePrice.Cmp(big.NewInt(60)) != 0 {
		t.Errorf("Effective gas price mismatch: got %s, want 60", effectivePrice)
	}
	
	// With baseFee = 95, effective = min(100, 95+10) = 100 (capped)
	baseFee = big.NewInt(95)
	effectivePrice = evmTx.EffectiveGasPrice(baseFee)
	if effectivePrice.Cmp(big.NewInt(100)) != 0 {
		t.Errorf("Effective gas price (capped) mismatch: got %s, want 100", effectivePrice)
	}
	
	t.Logf("✅ EIP-1559 gas price calculation correct")
}

// Test validation
func TestValidate(t *testing.T) {
	chainID := big.NewInt(1644)
	
	tests := []struct {
		name      string
		tx        *EvmTx
		wantError bool
	}{
		{
			name: "valid transaction (with sender set)",
			tx: &EvmTx{
				TxType:   types.LegacyTxType,
				Nonce:    0,
				GasLimit: 21000,
				GasPrice: big.NewInt(1000000000),
				Value:    big.NewInt(1000),
				V:        big.NewInt(3323),
				R:        big.NewInt(12345),
				S:        big.NewInt(67890),
				From:     common.HexToAddress("0x39a028dfdcae40bf277ec1ec268d62665d36c073"), // Set sender to skip recovery
			},
			wantError: false,
		},
		{
			name: "zero gas limit",
			tx: &EvmTx{
				TxType:   types.LegacyTxType,
				GasLimit: 0,
				Value:    big.NewInt(1000),
			},
			wantError: true,
		},
		{
			name: "negative value",
			tx: &EvmTx{
				TxType:   types.LegacyTxType,
				GasLimit: 21000,
				Value:    big.NewInt(-1),
			},
			wantError: true,
		},
		{
			name: "EIP-1559 tip exceeds cap",
			tx: &EvmTx{
				TxType:    types.DynamicFeeTxType,
				GasLimit:  21000,
				GasFeeCap: big.NewInt(100),
				GasTipCap: big.NewInt(200), // tip > cap!
				Value:     big.NewInt(1000),
			},
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tx.Validate(chainID)
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

