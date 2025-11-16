package evm

import (
	"math/big"
	"testing"

	"github.com/ArchivasNetwork/archivas/address"
	"github.com/ArchivasNetwork/archivas/types"
)

func TestEngineSimpleTransfer(t *testing.T) {
	// Create state DB
	stateDB := NewMemoryStateDB()

	// Setup: Give sender 100000 wei (enough for transfer + gas)
	sender, _ := address.EVMAddressFromHex("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
	receiver, _ := address.EVMAddressFromHex("0x1234567890abcdef1234567890abcdef12345678")
	
	stateDB.SetBalance(sender, big.NewInt(100000))
	stateDB.SetNonce(sender, 0)
	
	initialRoot, _ := stateDB.Commit()

	// Create EVM engine
	config := DefaultBetanetConfig()
	engine := NewEngine(config, stateDB)

	// Create transaction: transfer 100 wei
	tx := &types.EVMTransaction{
		TypeFlag:    types.TxTypeEVMCall,
		NonceVal:    0,
		GasPriceVal: big.NewInt(1),
		GasLimitVal: 21000,
		FromAddr:    sender,
		ToAddr:      &receiver,
		ValueVal:    big.NewInt(100),
		DataVal:     []byte{},
	}

	// Create block
	farmer, _ := address.EVMAddressFromHex("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	block := &types.Block{
		Height:      1,
		TimestampUnix: 1000,
		FarmerAddr:  farmer,
		GasLimit:    100000,
		Txs:         []types.Transaction{tx},
	}

	// Execute block
	result, err := engine.ExecuteBlock(block, initialRoot)
	if err != nil {
		t.Fatalf("ExecuteBlock failed: %v", err)
	}

	// Verify results
	if result.GasUsed != 21000 {
		t.Errorf("Expected gas used 21000, got %d", result.GasUsed)
	}

	if len(result.Receipts) != 1 {
		t.Fatalf("Expected 1 receipt, got %d", len(result.Receipts))
	}

	receipt := result.Receipts[0]
	if receipt.Status != 1 {
		t.Error("Expected successful transaction")
	}

	// Verify balances
	// Sender: 100000 - 100 (value) - 21000 (gas) = 78900
	// Receiver: 100
	// Farmer: 21000 (gas fees)
	
	senderBalance := stateDB.GetBalance(sender)
	receiverBalance := stateDB.GetBalance(receiver)
	farmerBalance := stateDB.GetBalance(farmer)
	
	expectedSender := big.NewInt(78900)
	expectedReceiver := big.NewInt(100)
	expectedFarmer := big.NewInt(21000)
	
	if senderBalance.Cmp(expectedSender) != 0 {
		t.Errorf("Sender balance: expected %s, got %s", expectedSender, senderBalance)
	}
	if receiverBalance.Cmp(expectedReceiver) != 0 {
		t.Errorf("Receiver balance: expected %s, got %s", expectedReceiver, receiverBalance)
	}
	if farmerBalance.Cmp(expectedFarmer) != 0 {
		t.Errorf("Farmer balance: expected %s, got %s", expectedFarmer, farmerBalance)
	}
}

func TestEngineInsufficientBalance(t *testing.T) {
	// Create state DB
	stateDB := NewMemoryStateDB()

	// Setup: Give sender only 100 wei (not enough for gas + value)
	sender, _ := address.EVMAddressFromHex("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
	receiver, _ := address.EVMAddressFromHex("0x1234567890abcdef1234567890abcdef12345678")
	
	stateDB.SetBalance(sender, big.NewInt(100))
	stateDB.SetNonce(sender, 0)
	
	initialRoot, _ := stateDB.Commit()

	// Create EVM engine
	config := DefaultBetanetConfig()
	engine := NewEngine(config, stateDB)

	// Create transaction: transfer 100 wei with 21000 gas (will fail)
	tx := &types.EVMTransaction{
		TypeFlag:    types.TxTypeEVMCall,
		NonceVal:    0,
		GasPriceVal: big.NewInt(1),
		GasLimitVal: 21000,
		FromAddr:    sender,
		ToAddr:      &receiver,
		ValueVal:    big.NewInt(100),
		DataVal:     []byte{},
	}

	// Create block
	farmer, _ := address.EVMAddressFromHex("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	block := &types.Block{
		Height:        1,
		TimestampUnix: 1000,
		FarmerAddr:    farmer,
		GasLimit:      100000,
		Txs:           []types.Transaction{tx},
	}

	// Execute block - should fail but not error
	result, err := engine.ExecuteBlock(block, initialRoot)
	if err != nil {
		t.Fatalf("ExecuteBlock should not error on tx failure: %v", err)
	}

	// Receipt should show failure
	if len(result.Receipts) != 1 {
		t.Fatalf("Expected 1 receipt, got %d", len(result.Receipts))
	}

	receipt := result.Receipts[0]
	if receipt.Status != 0 {
		t.Error("Expected failed transaction (status=0)")
	}
}

func TestEngineNonceValidation(t *testing.T) {
	// Create state DB
	stateDB := NewMemoryStateDB()

	// Setup: Give sender enough balance
	sender, _ := address.EVMAddressFromHex("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
	receiver, _ := address.EVMAddressFromHex("0x1234567890abcdef1234567890abcdef12345678")
	
	stateDB.SetBalance(sender, big.NewInt(100000))
	stateDB.SetNonce(sender, 5) // Current nonce is 5
	
	initialRoot, _ := stateDB.Commit()

	// Create EVM engine
	config := DefaultBetanetConfig()
	engine := NewEngine(config, stateDB)

	// Create transaction with wrong nonce (3 instead of 5)
	tx := &types.EVMTransaction{
		TypeFlag:    types.TxTypeEVMCall,
		NonceVal:    3, // Wrong nonce
		GasPriceVal: big.NewInt(1),
		GasLimitVal: 21000,
		FromAddr:    sender,
		ToAddr:      &receiver,
		ValueVal:    big.NewInt(100),
		DataVal:     []byte{},
	}

	// Create block
	farmer, _ := address.EVMAddressFromHex("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	block := &types.Block{
		Height:        1,
		TimestampUnix: 1000,
		FarmerAddr:    farmer,
		GasLimit:      100000,
		Txs:           []types.Transaction{tx},
	}

	// Execute block
	result, err := engine.ExecuteBlock(block, initialRoot)
	if err != nil {
		t.Fatalf("ExecuteBlock should not error: %v", err)
	}

	// Transaction should have failed
	if len(result.Receipts) != 1 {
		t.Fatalf("Expected 1 receipt, got %d", len(result.Receipts))
	}

	receipt := result.Receipts[0]
	if receipt.Status != 0 {
		t.Error("Expected failed transaction due to nonce mismatch")
	}
}

func TestStateDBSnapshot(t *testing.T) {
	stateDB := NewMemoryStateDB()

	addr, _ := address.EVMAddressFromHex("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
	
	// Set initial balance
	stateDB.SetBalance(addr, big.NewInt(1000))
	
	// Take snapshot
	snap := stateDB.Snapshot()
	
	// Modify balance
	stateDB.SetBalance(addr, big.NewInt(2000))
	balance := stateDB.GetBalance(addr)
	if balance.Cmp(big.NewInt(2000)) != 0 {
		t.Errorf("Expected balance 2000, got %s", balance.String())
	}
	
	// Revert to snapshot
	stateDB.RevertToSnapshot(snap)
	balance = stateDB.GetBalance(addr)
	if balance.Cmp(big.NewInt(1000)) != 0 {
		t.Errorf("Expected balance 1000 after revert, got %s", balance.String())
	}
}

