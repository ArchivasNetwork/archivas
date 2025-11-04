package mempool

import "github.com/ArchivasNetwork/archivas/ledger"

// Mempool stores pending transactions
type Mempool struct {
	txs []ledger.Transaction
}

// NewMempool creates a new empty mempool
func NewMempool() *Mempool {
	return &Mempool{
		txs: make([]ledger.Transaction, 0),
	}
}

// Add adds a transaction to the mempool
func (m *Mempool) Add(tx ledger.Transaction) {
	m.txs = append(m.txs, tx)
}

// Pending returns all pending transactions
func (m *Mempool) Pending() []ledger.Transaction {
	return m.txs
}

// Clear removes all transactions from the mempool
func (m *Mempool) Clear() {
	m.txs = make([]ledger.Transaction, 0)
}

