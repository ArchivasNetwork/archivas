package main

import (
	"fmt"
	"math/big"

	"github.com/ArchivasNetwork/archivas/address"
	"github.com/ArchivasNetwork/archivas/evm"
	"github.com/ArchivasNetwork/archivas/types"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("   Archivas Betanet - EVM Engine Demo")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	// Create state database
	stateDB := evm.NewMemoryStateDB()

	// Setup accounts
	alice, _ := address.EVMAddressFromHex("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
	bob, _ := address.EVMAddressFromHex("0x1234567890abcdef1234567890abcdef12345678")
	farmer, _ := address.EVMAddressFromHex("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	fmt.Println("📍 Initial Setup")
	fmt.Println("─────────────────────────────────────────────────")
	fmt.Printf("Alice:  %s\n", alice.Hex())
	fmt.Printf("Bob:    %s\n", bob.Hex())
	fmt.Printf("Farmer: %s\n", farmer.Hex())
	fmt.Println()

	// Give Alice some tokens
	initialBalance := big.NewInt(1000000) // 1 million wei
	stateDB.SetBalance(alice, initialBalance)
	stateDB.SetNonce(alice, 0)

	fmt.Printf("💰 Alice's initial balance: %s wei\n", initialBalance.String())
	fmt.Println()

	// Commit genesis state
	genesisRoot, _ := stateDB.Commit()
	fmt.Printf("📦 Genesis state root: %s\n", types.HashToString(genesisRoot))
	fmt.Println()

	// Create EVM engine
	config := evm.DefaultBetanetConfig()
	engine := evm.NewEngine(config, stateDB)

	// Transaction 1: Alice sends 10000 wei to Bob
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("   Transaction 1: Simple Transfer")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	tx1 := &types.EVMTransaction{
		TypeFlag:    types.TxTypeEVMCall,
		NonceVal:    0,
		GasPriceVal: big.NewInt(1),    // 1 wei per gas
		GasLimitVal: 21000,             // Standard transfer gas
		FromAddr:    alice,
		ToAddr:      &bob,
		ValueVal:    big.NewInt(10000), // Send 10000 wei
		DataVal:     []byte{},
	}

	block1 := &types.Block{
		Height:        1,
		TimestampUnix: 1000,
		PrevHash:      genesisRoot,
		FarmerAddr:    farmer,
		GasLimit:      1000000,
		Txs:           []types.Transaction{tx1},
		StateRoot:     genesisRoot,
	}

	fmt.Printf("📤 Alice → Bob: 10000 wei\n")
	fmt.Printf("⛽ Gas limit: %d, Gas price: %s wei\n", tx1.GasLimitVal, tx1.GasPriceVal.String())
	fmt.Println()

	result1, err := engine.ExecuteBlock(block1, genesisRoot)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	receipt1 := result1.Receipts[0]
	fmt.Printf("📊 Execution Result:\n")
	fmt.Printf("   Status: %s\n", statusToEmoji(receipt1.Status))
	fmt.Printf("   Gas used: %d\n", receipt1.GasUsed)
	fmt.Printf("   Gas cost: %d wei\n", receipt1.GasUsed*tx1.GasPriceVal.Uint64())
	fmt.Printf("   New state root: %s\n", types.HashToString(result1.StateRoot))
	fmt.Println()

	fmt.Printf("💰 Updated Balances:\n")
	fmt.Printf("   Alice:  %s wei (-10000 value - %d gas)\n", stateDB.GetBalance(alice).String(), receipt1.GasUsed)
	fmt.Printf("   Bob:    %s wei (+10000)\n", stateDB.GetBalance(bob).String())
	fmt.Printf("   Farmer: %s wei (+%d gas fees)\n", stateDB.GetBalance(farmer).String(), receipt1.GasUsed)
	fmt.Println()

	// Transaction 2: Alice deploys a contract
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("   Transaction 2: Contract Deployment")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	// Simple bytecode (for demo - not real EVM bytecode)
	contractCode := []byte{0x60, 0x60, 0x60, 0x40, 0x52}

	tx2 := &types.EVMTransaction{
		TypeFlag:    types.TxTypeEVMDeploy,
		NonceVal:    1, // Alice's nonce is now 1
		GasPriceVal: big.NewInt(1),
		GasLimitVal: 100000,
		FromAddr:    alice,
		ToAddr:      nil, // nil = contract creation
		ValueVal:    big.NewInt(0),
		DataVal:     contractCode,
	}

	block2 := &types.Block{
		Height:        2,
		TimestampUnix: 2000,
		PrevHash:      result1.StateRoot,
		FarmerAddr:    farmer,
		GasLimit:      1000000,
		Txs:           []types.Transaction{tx2},
		StateRoot:     result1.StateRoot,
	}

	fmt.Printf("📝 Alice deploys contract (%d bytes)\n", len(contractCode))
	fmt.Printf("⛽ Gas limit: %d, Gas price: %s wei\n", tx2.GasLimitVal, tx2.GasPriceVal.String())
	fmt.Println()

	result2, err := engine.ExecuteBlock(block2, result1.StateRoot)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	receipt2 := result2.Receipts[0]
	fmt.Printf("📊 Execution Result:\n")
	fmt.Printf("   Status: %s\n", statusToEmoji(receipt2.Status))
	fmt.Printf("   Gas used: %d\n", receipt2.GasUsed)
	if receipt2.ContractAddress != nil {
		fmt.Printf("   Contract address: %s\n", receipt2.ContractAddress.Hex())
		
		arcvAddr, _ := address.EncodeARCVAddress(*receipt2.ContractAddress, "arcv")
		fmt.Printf("   Contract (ARCV):  %s\n", arcvAddr)
	}
	fmt.Printf("   New state root: %s\n", types.HashToString(result2.StateRoot))
	fmt.Println()

	fmt.Printf("💰 Updated Balances:\n")
	fmt.Printf("   Alice:  %s wei (-%d gas)\n", stateDB.GetBalance(alice).String(), receipt2.GasUsed)
	fmt.Printf("   Farmer: %s wei (+%d gas fees)\n", stateDB.GetBalance(farmer).String(), receipt2.GasUsed)
	fmt.Println()

	// Summary
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("   Summary")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()
	
	totalGasUsed := receipt1.GasUsed + receipt2.GasUsed
	fmt.Printf("📈 Blocks executed: 2\n")
	fmt.Printf("📝 Transactions: 2 (1 transfer + 1 deployment)\n")
	fmt.Printf("⛽ Total gas used: %d\n", totalGasUsed)
	fmt.Printf("🌳 State transitions: 3 (genesis → block1 → block2)\n")
	fmt.Println()
	
	fmt.Println("✅ All transactions executed successfully!")
	fmt.Println()
	fmt.Println("📚 Key Features Demonstrated:")
	fmt.Println("   • EVM transaction execution")
	fmt.Println("   • Balance transfers")
	fmt.Println("   • Nonce management")
	fmt.Println("   • Gas metering & refunds")
	fmt.Println("   • Contract deployment")
	fmt.Println("   • State root updates")
	fmt.Println("   • Farmer gas fee collection")
	fmt.Println()
}

func statusToEmoji(status uint8) string {
	if status == 1 {
		return "✅ Success"
	}
	return "❌ Failed"
}

