package evm

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ArchivasNetwork/archivas/address"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

// VMExecutor provides full EVM transaction execution using go-ethereum
type VMExecutor struct {
	chainID     *big.Int
	chainConfig *params.ChainConfig
}

// NewVMExecutor creates a new EVM executor with Betanet configuration
func NewVMExecutor(chainID uint64) *VMExecutor {
	return &VMExecutor{
		chainID: big.NewInt(int64(chainID)),
		chainConfig: &params.ChainConfig{
			ChainID:             big.NewInt(int64(chainID)),
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			BerlinBlock:         big.NewInt(0),
			LondonBlock:         big.NewInt(0), // EIP-1559 enabled
		},
	}
}

// TxExecutionResult contains the result of a single EVM transaction execution
type TxExecutionResult struct {
	GasUsed         uint64
	Status          uint8 // 1 = success, 0 = failure
	ReturnData      []byte
	Logs            []*types.Log
	ContractAddress *common.Address
	Err             error
}

// ExecuteTransaction executes an EVM transaction in the given StateDB
func (e *VMExecutor) ExecuteTransaction(
	statedb *state.StateDB,
	evmTx *EvmTx,
	blockNumber uint64,
	blockTime uint64,
	coinbase common.Address,
	blockGasLimit uint64,
	baseFee *big.Int,
) (*TxExecutionResult, error) {
	
	// Create block context
	blockContext := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash: func(n uint64) common.Hash {
			// For now, return zero hash
			return common.Hash{}
		},
		Coinbase:    coinbase,
		BlockNumber: big.NewInt(int64(blockNumber)),
		Time:        blockTime,
		Difficulty:  big.NewInt(0),
		GasLimit:    blockGasLimit,
		BaseFee:     baseFee,
	}

	// Create transaction context
	sender := common.Address(evmTx.From)
	
	// Create EVM instance
	vmConfig := vm.Config{}
	evm := vm.NewEVM(blockContext, statedb, e.chainConfig, vmConfig)

	// Prepare gas
	gasPool := new(core.GasPool).AddGas(blockGasLimit)
	
	// Check sender balance (in Wei)
	senderBalance := statedb.GetBalance(sender)
	
	// Calculate total cost (value + max gas cost)
	maxGasCost := new(uint256.Int).Mul(
		uint256.MustFromBig(evmTx.GasPrice),
		uint256.NewInt(evmTx.GasLimit),
	)
	totalCost := new(uint256.Int).Add(uint256.MustFromBig(evmTx.Value), maxGasCost)
	
	if senderBalance.Cmp(totalCost) < 0 {
		return nil, fmt.Errorf("insufficient funds: have %s Wei, need %s Wei", senderBalance, totalCost)
	}

	// Prepare recipient
	var recipient *common.Address
	if evmTx.To != nil {
		addr := common.Address(*evmTx.To)
		recipient = &addr
	}

	// Increment nonce BEFORE execution
	currentNonce := statedb.GetNonce(sender)
	if evmTx.Nonce != currentNonce {
		return nil, fmt.Errorf("invalid nonce: have %d, want %d", evmTx.Nonce, currentNonce)
	}
	statedb.SetNonce(sender, evmTx.Nonce+1, tracing.NonceChangeUnspecified)

	// Deduct max gas cost upfront (will be refunded)
	statedb.SubBalance(sender, maxGasCost, 0)

	var (
		ret         []byte
		gasUsed     uint64
		failed      bool
		vmerr       error
		createdAddr *common.Address
	)

	// Convert value to uint256
	value := uint256.MustFromBig(evmTx.Value)

	if recipient == nil {
		// Contract creation
		ret, contractAddr, leftoverGas, vmerr := evm.Create(
			sender, // Use sender address directly
			evmTx.Data,
			gasPool.Gas(),
			value,
		)
		gasUsed = evmTx.GasLimit - leftoverGas
		
		if vmerr != nil {
			failed = true
			log.Printf("[EVM] Contract creation failed: %v", vmerr)
		} else {
			createdAddr = &contractAddr
			log.Printf("[EVM] Contract created at %s, gas used: %d", contractAddr.Hex(), gasUsed)
		}
		_ = ret // Avoid unused warning
	} else {
		// Contract call or value transfer
		ret, leftoverGas, vmerr := evm.Call(
			sender, // Use sender address directly
			*recipient,
			evmTx.Data,
			gasPool.Gas(),
			value,
		)
		gasUsed = evmTx.GasLimit - leftoverGas
		
		if vmerr != nil {
			failed = true
			log.Printf("[EVM] Call to %s failed: %v", recipient.Hex(), vmerr)
		} else {
			log.Printf("[EVM] Call to %s succeeded, gas used: %d", recipient.Hex(), gasUsed)
		}
		_ = ret // Avoid unused warning
	}

	// Refund unused gas
	gasRefund := new(uint256.Int).Mul(
		uint256.MustFromBig(evmTx.GasPrice),
		uint256.NewInt(evmTx.GasLimit-gasUsed),
	)
	statedb.AddBalance(sender, gasRefund, 0)

	// Pay gas fees to coinbase (miner/farmer)
	gasFee := new(uint256.Int).Mul(
		uint256.MustFromBig(evmTx.GasPrice),
		uint256.NewInt(gasUsed),
	)
	statedb.AddBalance(coinbase, gasFee, 0)

	// Collect logs from StateDB
	// Note: In go-ethereum v1.16+, logs are tracked internally
	logs := statedb.Logs()

	// Build result
	status := uint8(1) // Success
	if failed {
		status = 0
	}

	result := &TxExecutionResult{
		GasUsed:         gasUsed,
		Status:          status,
		ReturnData:      ret,
		Logs:            logs,
		ContractAddress: createdAddr,
		Err:             vmerr,
	}

	return result, nil
}

// ConvertLogsToArchivas converts go-ethereum logs to Archivas log format
// Note: Log type is defined in statedb.go
func ConvertLogsToArchivas(
	ethLogs []*types.Log,
	txHash [32]byte,
	txIndex uint32,
	blockHeight uint64,
) []Log {
	logs := make([]Log, len(ethLogs))
	for i, ethLog := range ethLogs {
		topics := make([][32]byte, len(ethLog.Topics))
		for j, topic := range ethLog.Topics {
			topics[j] = topic
		}
		
		logs[i] = Log{
			Address:     address.EVMAddress(ethLog.Address),
			Topics:      topics,
			Data:        ethLog.Data,
			TxHash:      txHash,
			TxIndex:     txIndex,
			BlockHeight: blockHeight,
			Index:       uint32(i),
		}
	}
	return logs
}
