package node

import (
	"math/big"
	"testing"

	"github.com/ArchivasNetwork/archivas/evm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Mock state functions
func mockGetBalance(addr common.Address) *big.Int {
	// Default: 1000 RCHV = 10^18 Wei
	return new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e18))
}

func mockGetNonce(addr common.Address) uint64 {
	return 0
}

// Helper to create a test transaction
func createTestTx(t *testing.T, nonce uint64, gasPrice *big.Int, value *big.Int) *evm.EvmTx {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	
	chainID := big.NewInt(1644)
	signer := types.NewLondonSigner(chainID)
	
	to := common.HexToAddress("0x676979f6aca633dadabdad81f9ef0310f0eae39e")
	
	tx := types.NewTransaction(
		nonce,
		to,
		value,
		21000,
		gasPrice,
		nil,
	)
	
	signedTx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		t.Fatalf("Failed to sign tx: %v", err)
	}
	
	rawTx, err := signedTx.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal tx: %v", err)
	}
	
	evmTx, err := evm.FromRawTransaction(rawTx, chainID)
	if err != nil {
		t.Fatalf("Failed to decode tx: %v", err)
	}
	
	return evmTx
}

func TestEVMMempool_Add(t *testing.T) {
	mp := NewEVMMempool(mockGetBalance, mockGetNonce)
	
	tx := createTestTx(t, 0, big.NewInt(1000000000), big.NewInt(1000))
	
	err := mp.Add(tx)
	if err != nil {
		t.Fatalf("Failed to add valid transaction: %v", err)
	}
	
	if mp.Size() != 1 {
		t.Errorf("Expected size 1, got %d", mp.Size())
	}
	
	// Try to add same tx again
	err = mp.Add(tx)
	if err == nil {
		t.Error("Expected error when adding duplicate transaction")
	}
	
	t.Logf("✅ Successfully added transaction to mempool")
}

func TestEVMMempool_NonceTooLow(t *testing.T) {
	getNonce := func(addr common.Address) uint64 {
		return 5 // Current nonce is 5
	}
	
	mp := NewEVMMempool(mockGetBalance, getNonce)
	
	// Try to add tx with nonce 3 (too low)
	tx := createTestTx(t, 3, big.NewInt(1000000000), big.NewInt(1000))
	
	err := mp.Add(tx)
	if err == nil {
		t.Error("Expected error for nonce too low")
	} else if err.Error() != "nonce too low: have 3, want >= 5" {
		t.Errorf("Wrong error message: %v", err)
	}
	
	t.Logf("✅ Correctly rejected transaction with low nonce")
}

func TestEVMMempool_InsufficientBalance(t *testing.T) {
	getBalance := func(addr common.Address) *big.Int {
		return big.NewInt(1000) // Only 1000 wei
	}
	
	mp := NewEVMMempool(getBalance, mockGetNonce)
	
	// Try to send 10000 wei with 1 gwei gas price and 21000 gas limit
	// Total cost = 10000 + (1000000000 * 21000) = 21000010000 wei
	tx := createTestTx(t, 0, big.NewInt(1000000000), big.NewInt(10000))
	
	err := mp.Add(tx)
	if err == nil {
		t.Error("Expected error for insufficient balance")
	}
	
	t.Logf("✅ Correctly rejected transaction with insufficient balance: %v", err)
}

func TestEVMMempool_GetExecutable(t *testing.T) {
	mp := NewEVMMempool(mockGetBalance, mockGetNonce)
	
	// Add transactions from different senders (all with nonce 0, which matches mock)
	tx0 := createTestTx(t, 0, big.NewInt(2000000000), big.NewInt(1000)) // 2 gwei
	tx1 := createTestTx(t, 0, big.NewInt(1000000000), big.NewInt(1000)) // 1 gwei
	tx2 := createTestTx(t, 0, big.NewInt(3000000000), big.NewInt(1000)) // 3 gwei (highest)
	
	mp.Add(tx0)
	mp.Add(tx1)
	mp.Add(tx2)
	
	// Get executable (should return all 3 since they all have nonce=0 matching current)
	// They should be sorted by gas price (highest first)
	executable := mp.GetExecutable(10)
	
	if len(executable) != 3 {
		t.Errorf("Expected 3 executable transactions, got %d", len(executable))
	}
	
	// First should have highest gas price (3 gwei)
	if executable[0].GasPrice.Cmp(big.NewInt(3000000000)) != 0 {
		t.Errorf("Expected highest gas price first, got %s", executable[0].GasPrice)
	}
	
	t.Logf("✅ GetExecutable returned %d transactions sorted by price", len(executable))
	for i, tx := range executable {
		t.Logf("   %d: nonce=%d, price=%s gwei", i+1, tx.Nonce, tx.GasPrice)
	}
}

func TestEVMMempool_SortByGasPrice(t *testing.T) {
	mp := NewEVMMempool(mockGetBalance, mockGetNonce)
	
	// Create transactions from different senders with different gas prices
	tx1 := createTestTx(t, 0, big.NewInt(1000000000), big.NewInt(1000)) // 1 gwei
	tx2 := createTestTx(t, 0, big.NewInt(3000000000), big.NewInt(1000)) // 3 gwei (highest)
	tx3 := createTestTx(t, 0, big.NewInt(2000000000), big.NewInt(1000)) // 2 gwei
	
	mp.Add(tx1)
	mp.Add(tx2)
	mp.Add(tx3)
	
	// Get executable (should be sorted by gas price, highest first)
	executable := mp.GetExecutable(10)
	
	if len(executable) != 3 {
		t.Errorf("Expected 3 executable transactions, got %d", len(executable))
	}
	
	// Verify sorted by gas price (descending)
	for i := 0; i < len(executable)-1; i++ {
		price1 := executable[i].EffectiveGasPrice(mp.baseFee)
		price2 := executable[i+1].EffectiveGasPrice(mp.baseFee)
		if price1.Cmp(price2) < 0 {
			t.Errorf("Transactions not sorted by gas price: %s < %s", price1, price2)
		}
	}
	
	t.Logf("✅ Transactions correctly sorted by gas price")
	t.Logf("   1st: %s gwei", executable[0].GasPrice)
	t.Logf("   2nd: %s gwei", executable[1].GasPrice)
	t.Logf("   3rd: %s gwei", executable[2].GasPrice)
}

func TestEVMMempool_Remove(t *testing.T) {
	mp := NewEVMMempool(mockGetBalance, mockGetNonce)
	
	tx := createTestTx(t, 0, big.NewInt(1000000000), big.NewInt(1000))
	
	mp.Add(tx)
	if mp.Size() != 1 {
		t.Fatalf("Expected size 1 after add")
	}
	
	mp.Remove(tx.Hash)
	if mp.Size() != 0 {
		t.Errorf("Expected size 0 after remove, got %d", mp.Size())
	}
	
	t.Logf("✅ Successfully removed transaction from mempool")
}

func TestEVMMempool_Revalidate(t *testing.T) {
	// Start with nonce 0
	currentNonce := uint64(0)
	getNonce := func(addr common.Address) uint64 {
		return currentNonce
	}
	
	mp := NewEVMMempool(mockGetBalance, getNonce)
	
	// Add transactions with nonces 0, 1, 2
	tx0 := createTestTx(t, 0, big.NewInt(1000000000), big.NewInt(1000))
	tx1 := createTestTx(t, 1, big.NewInt(1000000000), big.NewInt(1000))
	tx2 := createTestTx(t, 2, big.NewInt(1000000000), big.NewInt(1000))
	
	// Use same sender
	sender := common.HexToAddress("0x2222222222222222222222222222222222222222")
	tx0.From = sender
	tx1.From = sender
	tx2.From = sender
	
	mp.Add(tx0)
	mp.Add(tx1)
	mp.Add(tx2)
	
	if mp.Size() != 3 {
		t.Fatalf("Expected 3 transactions in mempool, got %d", mp.Size())
	}
	
	// Simulate tx0 being mined (nonce advances to 1)
	currentNonce = 1
	
	// Revalidate should remove tx0 (nonce too low)
	removed := mp.Revalidate()
	
	if removed != 1 {
		t.Errorf("Expected 1 transaction removed, got %d", removed)
	}
	
	if mp.Size() != 2 {
		t.Errorf("Expected 2 transactions remaining, got %d", mp.Size())
	}
	
	t.Logf("✅ Revalidate correctly removed stale transaction")
}

func TestEVMMempool_PerAccountLimit(t *testing.T) {
	mp := NewEVMMempool(mockGetBalance, mockGetNonce)
	mp.maxTxsPerAccount = 3 // Limit to 3 txs per account
	
	sender := common.HexToAddress("0x3333333333333333333333333333333333333333")
	
	// Add 3 transactions (should succeed)
	for i := uint64(0); i < 3; i++ {
		tx := createTestTx(t, i, big.NewInt(1000000000), big.NewInt(1000))
		tx.From = sender
		if err := mp.Add(tx); err != nil {
			t.Fatalf("Failed to add tx %d: %v", i, err)
		}
	}
	
	// Try to add 4th transaction (should fail)
	tx4 := createTestTx(t, 3, big.NewInt(1000000000), big.NewInt(1000))
	tx4.From = sender
	
	err := mp.Add(tx4)
	if err == nil {
		t.Error("Expected error when exceeding per-account limit")
	}
	
	t.Logf("✅ Per-account limit correctly enforced: %v", err)
}

func TestEVMMempool_Stats(t *testing.T) {
	mp := NewEVMMempool(mockGetBalance, mockGetNonce)
	
	// Add some transactions
	tx1 := createTestTx(t, 0, big.NewInt(1000000000), big.NewInt(1000))
	tx2 := createTestTx(t, 0, big.NewInt(1000000000), big.NewInt(1000))
	
	mp.Add(tx1)
	mp.Add(tx2)
	
	stats := mp.Stats()
	
	totalTxs, ok := stats["total_transactions"].(int)
	if !ok || totalTxs != 2 {
		t.Errorf("Expected 2 total transactions, got %v", stats["total_transactions"])
	}
	
	uniqueAccounts, ok := stats["unique_accounts"].(int)
	if !ok || uniqueAccounts != 2 {
		t.Errorf("Expected 2 unique accounts, got %v", stats["unique_accounts"])
	}
	
	t.Logf("✅ Stats: %v", stats)
}

