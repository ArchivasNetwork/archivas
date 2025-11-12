package ledger

// AccountState represents the state of a single account
type AccountState struct {
	Balance int64  // in base units of RCHV
	Nonce   uint64 // increments each tx
}

// WorldState represents the global state of all accounts
type WorldState struct {
	Accounts map[string]*AccountState
}

// NewWorldState creates a new world state with a funded genesis account.
func NewWorldState(genesisAlloc map[string]int64) *WorldState {
	ws := &WorldState{
		Accounts: make(map[string]*AccountState),
	}
	for addr, amt := range genesisAlloc {
		ws.Accounts[addr] = &AccountState{
			Balance: amt,
			Nonce:   0,
		}
	}
	return ws
}

// GetBalance returns the balance of an account (0 if account doesn't exist)
func (ws *WorldState) GetBalance(addr string) int64 {
	acct, ok := ws.Accounts[addr]
	if !ok {
		return 0
	}
	return acct.Balance
}

// GetNonce returns the nonce of an account (0 if account doesn't exist)
func (ws *WorldState) GetNonce(addr string) uint64 {
	acct, ok := ws.Accounts[addr]
	if !ok {
		return 0
	}
	return acct.Nonce
}

// GetAllAccountsWithBalance returns all addresses that have a non-zero balance
func (ws *WorldState) GetAllAccountsWithBalance() map[string]*AccountState {
	result := make(map[string]*AccountState)
	for addr, acct := range ws.Accounts {
		if acct.Balance > 0 {
			result[addr] = acct
		}
	}
	return result
}

