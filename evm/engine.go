package evm

import (
	"fmt"
	"math/big"

	"github.com/ArchivasNetwork/archivas/address"
	"github.com/ArchivasNetwork/archivas/types"
)

// Engine is the Archivas EVM execution engine
// This wraps go-ethereum's EVM for Betanet
type Engine struct {
	chainConfig *ChainConfig
	stateDB     StateDB
}

// ChainConfig contains EVM chain configuration
type ChainConfig struct {
	ChainID             *big.Int // Network chain ID
	HomesteadBlock      *big.Int
	EIP150Block         *big.Int
	EIP155Block         *big.Int
	EIP158Block         *big.Int
	ByzantiumBlock      *big.Int
	ConstantinopleBlock *big.Int
	PetersburgBlock     *big.Int
	IstanbulBlock       *big.Int
	BerlinBlock         *big.Int
	LondonBlock         *big.Int
}

// NewEngine creates a new EVM engine
func NewEngine(chainConfig *ChainConfig, stateDB StateDB) *Engine {
	return &Engine{
		chainConfig: chainConfig,
		stateDB:     stateDB,
	}
}

// ExecuteBlock executes all transactions in a block
// Returns new state root, receipts root, receipts, and total gas used
func (e *Engine) ExecuteBlock(
	block *types.Block,
	parentStateRoot [32]byte,
) (*ExecutionResult, error) {
	// Load state from parent
	if err := e.stateDB.SetRoot(parentStateRoot); err != nil {
		return nil, fmt.Errorf("failed to set state root: %w", err)
	}

	receipts := make([]*types.Receipt, 0, len(block.Txs))
	totalGasUsed := uint64(0)
	cumulativeGasUsed := uint64(0)

	// Execute each transaction
	for txIndex, tx := range block.Txs {
		evmTx, ok := tx.(*types.EVMTransaction)
		if !ok {
			// Skip non-EVM transactions (legacy transactions)
			continue
		}

		// Execute transaction
		receipt, err := e.ExecuteTransaction(evmTx, block, uint32(txIndex))
		if err != nil {
			// Transaction execution failed, but we still create a receipt
			receipt = &types.Receipt{
				TxHash:            tx.Hash(),
				BlockHeight:       block.Height,
				TxIndex:           uint32(txIndex),
				From:              evmTx.From(),
				To:                evmTx.To(),
				GasUsed:           evmTx.Gas(), // Consume all gas on failure
				Status:            0,            // Failed
				Logs:              []types.Log{},
				CumulativeGasUsed: cumulativeGasUsed + evmTx.Gas(),
			}
		}

		cumulativeGasUsed = receipt.CumulativeGasUsed
		totalGasUsed += receipt.GasUsed
		receipts = append(receipts, receipt)

		// Check gas limit
		if totalGasUsed > block.GasLimit {
			return nil, fmt.Errorf("block gas limit exceeded: used %d, limit %d", totalGasUsed, block.GasLimit)
		}
	}

	// Commit state changes
	newStateRoot, err := e.stateDB.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit state: %w", err)
	}

	// Compute receipts root
	receiptsRoot := computeReceiptsRoot(receipts)

	return &ExecutionResult{
		StateRoot:    newStateRoot,
		ReceiptsRoot: receiptsRoot,
		Receipts:     receipts,
		GasUsed:      totalGasUsed,
	}, nil
}

// ExecuteTransaction executes a single EVM transaction
func (e *Engine) ExecuteTransaction(
	tx *types.EVMTransaction,
	block *types.Block,
	txIndex uint32,
) (*types.Receipt, error) {
	// Create execution context
	ctx := &ExecutionContext{
		Origin:      tx.From(),
		GasPrice:    tx.GasPrice(),
		BlockHeight: block.Height,
		BlockTime:   block.TimestampUnix,
		Coinbase:    block.FarmerAddr, // Farmer receives gas fees
		GasLimit:    tx.Gas(),
	}

	// Verify sender has enough balance
	senderBalance := e.stateDB.GetBalance(tx.From())
	txCost := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(tx.Gas())))
	txCost.Add(txCost, tx.Value())

	if senderBalance.Cmp(txCost) < 0 {
		return nil, fmt.Errorf("insufficient balance: have %s, need %s", senderBalance.String(), txCost.String())
	}

	// Verify nonce
	expectedNonce := e.stateDB.GetNonce(tx.From())
	if tx.Nonce() != expectedNonce {
		return nil, fmt.Errorf("invalid nonce: have %d, want %d", tx.Nonce(), expectedNonce)
	}

	// Deduct gas cost upfront (but not value - that's handled in execution)
	gasCost := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(tx.Gas())))
	e.stateDB.SubBalance(tx.From(), gasCost)

	var (
		contractAddr *address.EVMAddress
		gasUsed      uint64
		logs         []types.Log
		failed       bool
	)

	if tx.IsContractCreation() {
		// Contract deployment
		addr, gas, err := e.executeContractCreation(tx, ctx)
		if err != nil {
			failed = true
			gasUsed = tx.Gas() // Consume all gas on failure
		} else {
			contractAddr = &addr
			gasUsed = gas
		}
	} else {
		// Contract call or transfer
		gas, txLogs, err := e.executeCall(tx, ctx)
		if err != nil {
			failed = true
			gasUsed = tx.Gas() // Consume all gas on failure
		} else {
			gasUsed = gas
			logs = txLogs
		}
	}

	// Refund unused gas
	refund := new(big.Int).Mul(
		tx.GasPrice(),
		big.NewInt(int64(tx.Gas()-gasUsed)),
	)
	e.stateDB.AddBalance(tx.From(), refund)

	// Pay gas fees to farmer (coinbase)
	gasFee := new(big.Int).Mul(
		tx.GasPrice(),
		big.NewInt(int64(gasUsed)),
	)
	e.stateDB.AddBalance(block.FarmerAddr, gasFee)

	// Increment sender nonce
	e.stateDB.SetNonce(tx.From(), expectedNonce+1)

	// Create receipt
	status := uint8(1)
	if failed {
		status = 0
	}

	receipt := &types.Receipt{
		TxHash:            tx.Hash(),
		BlockHeight:       block.Height,
		TxIndex:           txIndex,
		From:              tx.From(),
		To:                tx.To(),
		ContractAddress:   contractAddr,
		GasUsed:           gasUsed,
		Status:            status,
		Logs:              logs,
		CumulativeGasUsed: 0, // Will be set by caller
	}

	return receipt, nil
}

// executeContractCreation deploys a new contract
func (e *Engine) executeContractCreation(
	tx *types.EVMTransaction,
	ctx *ExecutionContext,
) (address.EVMAddress, uint64, error) {
	// Generate contract address (simplified - use CREATE address generation)
	contractAddr := e.generateContractAddress(tx.From(), tx.Nonce())

	// Check if address already exists
	if e.stateDB.GetCodeSize(contractAddr) > 0 {
		return address.ZeroAddress(), 0, fmt.Errorf("contract address already exists")
	}

	// For Phase 2, we'll implement a simplified execution
	// In Phase 2, we'll integrate go-ethereum's actual EVM

	// Deploy contract code
	code := tx.Data()
	e.stateDB.SetCode(contractAddr, code)

	// Transfer value if any
	if tx.Value().Sign() > 0 {
		e.stateDB.SubBalance(tx.From(), tx.Value())
		e.stateDB.AddBalance(contractAddr, tx.Value())
	}

	// Simplified gas calculation
	// Real EVM would compute based on opcodes executed
	gasUsed := uint64(53000) + uint64(len(code))*200

	if gasUsed > ctx.GasLimit {
		return address.ZeroAddress(), 0, fmt.Errorf("out of gas")
	}

	return contractAddr, gasUsed, nil
}

// executeCall executes a contract call or simple transfer
func (e *Engine) executeCall(
	tx *types.EVMTransaction,
	ctx *ExecutionContext,
) (uint64, []types.Log, error) {
	to := *tx.To()

	// Simple transfer (no data)
	if len(tx.Data()) == 0 {
		if tx.Value().Sign() > 0 {
			e.stateDB.SubBalance(tx.From(), tx.Value())
			e.stateDB.AddBalance(to, tx.Value())
		}
		return 21000, []types.Log{}, nil // Base gas cost
	}

	// Contract call
	// For Phase 2, simplified execution
	// In full implementation, would execute EVM bytecode

	// Transfer value if any
	if tx.Value().Sign() > 0 {
		e.stateDB.SubBalance(tx.From(), tx.Value())
		e.stateDB.AddBalance(to, tx.Value())
	}

	// Simplified gas calculation
	gasUsed := uint64(21000) + uint64(len(tx.Data()))*68

	if gasUsed > ctx.GasLimit {
		return 0, nil, fmt.Errorf("out of gas")
	}

	// No logs in simplified version
	return gasUsed, []types.Log{}, nil
}

// generateContractAddress generates a contract address from sender and nonce
// Uses CREATE address generation: keccak256(rlp([sender, nonce]))[:20]
func (e *Engine) generateContractAddress(sender address.EVMAddress, nonce uint64) address.EVMAddress {
	// Simplified version - just hash sender + nonce
	// Real implementation would use RLP encoding
	data := append(sender.Bytes(), byte(nonce))
	hash := types.Hash256(data)
	
	var addr address.EVMAddress
	copy(addr[:], hash[:20])
	return addr
}

// ExecutionContext contains block context for transaction execution
type ExecutionContext struct {
	Origin      address.EVMAddress // Transaction origin
	GasPrice    *big.Int           // Gas price
	BlockHeight uint64             // Current block height
	BlockTime   int64              // Current block timestamp
	Coinbase    address.EVMAddress // Block farmer (receives gas fees)
	GasLimit    uint64             // Gas limit for this transaction
}

// ExecutionResult contains the results of block execution
type ExecutionResult struct {
	StateRoot    [32]byte         // New state root after execution
	ReceiptsRoot [32]byte         // Receipts trie root
	Receipts     []*types.Receipt // Transaction receipts
	GasUsed      uint64           // Total gas used
}

// computeReceiptsRoot computes the merkle root of receipts
// Simplified version - proper implementation would use merkle patricia trie
func computeReceiptsRoot(receipts []*types.Receipt) [32]byte {
	if len(receipts) == 0 {
		return types.EmptyReceiptsRoot()
	}

	data := make([]byte, 0)
	for _, receipt := range receipts {
		receiptHash := receipt.Hash()
		data = append(data, receiptHash[:]...)
	}

	return types.Hash256(data)
}

