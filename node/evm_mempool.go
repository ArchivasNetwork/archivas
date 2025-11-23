package node

import (
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/ArchivasNetwork/archivas/evm"
	"github.com/ethereum/go-ethereum/common"
)

// EVMMempool tracks pending EVM transactions
// This is separate from the legacy mempool to maintain clean separation
type EVMMempool struct {
	mu          sync.RWMutex
	pending     map[common.Hash]*evm.EvmTx          // All pending transactions by hash
	byAddress   map[common.Address][]*evm.EvmTx     // Transactions grouped by sender
	baseFee     *big.Int                             // Current base fee (for EIP-1559)
	maxTxs      int                                  // Maximum transactions in mempool
	maxTxsPerAccount int                             // Maximum pending txs per account
	
	// Dependencies
	getBalance func(common.Address) *big.Int
	getNonce   func(common.Address) uint64
}

// NewEVMMempool creates a new EVM transaction mempool
func NewEVMMempool(
	getBalance func(common.Address) *big.Int,
	getNonce func(common.Address) uint64,
) *EVMMempool {
	return &EVMMempool{
		pending:          make(map[common.Hash]*evm.EvmTx),
		byAddress:        make(map[common.Address][]*evm.EvmTx),
		baseFee:          big.NewInt(1000000000), // 1 gwei default
		maxTxs:           5000,
		maxTxsPerAccount: 16,
		getBalance:       getBalance,
		getNonce:         getNonce,
	}
}

// Add validates and adds a transaction to the mempool
func (mp *EVMMempool) Add(tx *evm.EvmTx) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	
	// Basic validation
	if err := tx.Validate(tx.ChainID); err != nil {
		return fmt.Errorf("transaction validation failed: %w", err)
	}
	
	// Check if already in mempool
	if _, exists := mp.pending[tx.Hash]; exists {
		return fmt.Errorf("transaction already in mempool")
	}
	
	// Check mempool capacity
	if len(mp.pending) >= mp.maxTxs {
		return fmt.Errorf("mempool full (max %d transactions)", mp.maxTxs)
	}
	
	// Get current state for sender
	currentNonce := mp.getNonce(tx.From)
	balance := mp.getBalance(tx.From)
	
	// Validate nonce
	// Allow nonce >= currentNonce (future transactions are queued)
	if tx.Nonce < currentNonce {
		return fmt.Errorf("nonce too low: have %d, want >= %d", tx.Nonce, currentNonce)
	}
	
	// For immediate execution, nonce must match exactly
	if tx.Nonce > currentNonce+uint64(mp.maxTxsPerAccount) {
		return fmt.Errorf("nonce too high: have %d, current %d (max gap %d)",
			tx.Nonce, currentNonce, mp.maxTxsPerAccount)
	}
	
	// Validate balance
	cost := tx.Cost(mp.baseFee)
	if balance.Cmp(cost) < 0 {
		return fmt.Errorf("insufficient balance: have %s, need %s", balance, cost)
	}
	
	// For EIP-1559, check base fee
	if tx.TxType == 2 { // DynamicFeeTxType
		if tx.GasFeeCap.Cmp(mp.baseFee) < 0 {
			return fmt.Errorf("maxFeePerGas (%s) < baseFee (%s)", tx.GasFeeCap, mp.baseFee)
		}
	}
	
	// Check per-account limit
	senderTxs := mp.byAddress[tx.From]
	if len(senderTxs) >= mp.maxTxsPerAccount {
		// If new tx has higher gas price, consider replacing lowest
		if !mp.shouldReplace(senderTxs, tx) {
			return fmt.Errorf("too many pending transactions for account (max %d)", mp.maxTxsPerAccount)
		}
	}
	
	// Add to mempool
	mp.pending[tx.Hash] = tx
	mp.byAddress[tx.From] = append(mp.byAddress[tx.From], tx)
	
	// Sort by nonce for this sender
	mp.sortByNonce(tx.From)
	
	return nil
}

// shouldReplace determines if a new transaction should replace an existing one
func (mp *EVMMempool) shouldReplace(existing []*evm.EvmTx, newTx *evm.EvmTx) bool {
	// Find the transaction with the lowest gas price
	var lowest *evm.EvmTx
	lowestPrice := new(big.Int).SetUint64(^uint64(0)) // Max uint64
	
	for _, tx := range existing {
		price := tx.EffectiveGasPrice(mp.baseFee)
		if price.Cmp(lowestPrice) < 0 {
			lowestPrice = price
			lowest = tx
		}
	}
	
	// Replace if new tx has at least 10% higher gas price
	newPrice := newTx.EffectiveGasPrice(mp.baseFee)
	threshold := new(big.Int).Mul(lowestPrice, big.NewInt(110))
	threshold.Div(threshold, big.NewInt(100))
	
	if newPrice.Cmp(threshold) >= 0 && lowest != nil {
		// Remove the lowest priced transaction
		mp.remove(lowest.Hash)
		return true
	}
	
	return false
}

// sortByNonce sorts transactions for a given address by nonce
func (mp *EVMMempool) sortByNonce(addr common.Address) {
	txs := mp.byAddress[addr]
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Nonce < txs[j].Nonce
	})
}

// remove removes a transaction from the mempool (caller must hold lock)
func (mp *EVMMempool) remove(hash common.Hash) {
	tx, exists := mp.pending[hash]
	if !exists {
		return
	}
	
	delete(mp.pending, hash)
	
	// Remove from byAddress
	senderTxs := mp.byAddress[tx.From]
	for i, t := range senderTxs {
		if t.Hash == hash {
			mp.byAddress[tx.From] = append(senderTxs[:i], senderTxs[i+1:]...)
			break
		}
	}
	
	// Clean up empty entries
	if len(mp.byAddress[tx.From]) == 0 {
		delete(mp.byAddress, tx.From)
	}
}

// Remove removes a transaction from the mempool (public interface)
func (mp *EVMMempool) Remove(hash common.Hash) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.remove(hash)
}

// Get retrieves a transaction by hash
func (mp *EVMMempool) Get(hash common.Hash) (*evm.EvmTx, bool) {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	tx, exists := mp.pending[hash]
	return tx, exists
}

// GetPending returns all pending transactions for an address, sorted by nonce
func (mp *EVMMempool) GetPending(addr common.Address) []*evm.EvmTx {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	
	txs := mp.byAddress[addr]
	result := make([]*evm.EvmTx, len(txs))
	copy(result, txs)
	return result
}

// GetPendingCount returns the number of pending transactions for an address
func (mp *EVMMempool) GetPendingCount(addr common.Address) int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return len(mp.byAddress[addr])
}

// GetExecutable returns transactions ready for execution
// These are transactions with nonce matching the account's current nonce
func (mp *EVMMempool) GetExecutable(maxCount int) []*evm.EvmTx {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	
	var executable []*evm.EvmTx
	
	// For each address, find the first transaction with matching nonce
	for addr, txs := range mp.byAddress {
		currentNonce := mp.getNonce(addr)
		
		// Find executable transactions in sequence
		for _, tx := range txs {
			if tx.Nonce == currentNonce {
				executable = append(executable, tx)
				currentNonce++
				
				if len(executable) >= maxCount {
					goto done
				}
			} else if tx.Nonce > currentNonce {
				// Gap in nonce sequence, stop for this address
				break
			}
		}
	}
	
done:
	// Sort by gas price (descending)
	sort.Slice(executable, func(i, j int) bool {
		priceI := executable[i].EffectiveGasPrice(mp.baseFee)
		priceJ := executable[j].EffectiveGasPrice(mp.baseFee)
		return priceI.Cmp(priceJ) > 0
	})
	
	// Limit to maxCount
	if len(executable) > maxCount {
		executable = executable[:maxCount]
	}
	
	return executable
}

// GetAll returns all pending transactions
func (mp *EVMMempool) GetAll() []*evm.EvmTx {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	
	all := make([]*evm.EvmTx, 0, len(mp.pending))
	for _, tx := range mp.pending {
		all = append(all, tx)
	}
	return all
}

// Size returns the number of pending transactions
func (mp *EVMMempool) Size() int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return len(mp.pending)
}

// UpdateBaseFee updates the current base fee for EIP-1559
func (mp *EVMMempool) UpdateBaseFee(baseFee *big.Int) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.baseFee = new(big.Int).Set(baseFee)
}

// Clear removes all transactions from the mempool
func (mp *EVMMempool) Clear() {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.pending = make(map[common.Hash]*evm.EvmTx)
	mp.byAddress = make(map[common.Address][]*evm.EvmTx)
}

// Revalidate removes transactions that are no longer valid
// Should be called after blocks are added to the chain
func (mp *EVMMempool) Revalidate() int {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	
	removed := 0
	
	for hash, tx := range mp.pending {
		currentNonce := mp.getNonce(tx.From)
		balance := mp.getBalance(tx.From)
		
		// Remove if nonce is too low
		if tx.Nonce < currentNonce {
			mp.remove(hash)
			removed++
			continue
		}
		
		// Remove if insufficient balance
		cost := tx.Cost(mp.baseFee)
		if balance.Cmp(cost) < 0 {
			mp.remove(hash)
			removed++
			continue
		}
	}
	
	return removed
}

// Stats returns mempool statistics
func (mp *EVMMempool) Stats() map[string]interface{} {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	
	return map[string]interface{}{
		"total_transactions": len(mp.pending),
		"unique_accounts":    len(mp.byAddress),
		"base_fee":          mp.baseFee.String(),
		"max_capacity":      mp.maxTxs,
	}
}

