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
// Accepts both ARCV (arcv1...) and EVM (0x...) address formats
// For backward compatibility, checks both the input address and its normalized form
func (ws *WorldState) GetBalance(addr string) int64 {
	// Try direct lookup first (for legacy ARCV-keyed accounts)
	if acct, ok := ws.Accounts[addr]; ok {
		return acct.Balance
	}

	// Try normalized form (for new EVM-keyed accounts or when querying with alternate format)
	normAddr, err := NormalizeAddress(addr)
	if err != nil {
		return 0
	}

	if acct, ok := ws.Accounts[normAddr]; ok {
		return acct.Balance
	}

	return 0
}

// GetNonce returns the nonce of an account (0 if account doesn't exist)
// Accepts both ARCV (arcv1...) and EVM (0x...) address formats
// For backward compatibility, checks both the input address and its normalized form
func (ws *WorldState) GetNonce(addr string) uint64 {
	// Try direct lookup first (for legacy ARCV-keyed accounts)
	if acct, ok := ws.Accounts[addr]; ok {
		return acct.Nonce
	}

	// Try normalized form (for new EVM-keyed accounts or when querying with alternate format)
	normAddr, err := NormalizeAddress(addr)
	if err != nil {
		return 0
	}

	if acct, ok := ws.Accounts[normAddr]; ok {
		return acct.Nonce
	}

	return 0
}

// GetAccount returns the account state for an address (nil if doesn't exist)
// Accepts both ARCV (arcv1...) and EVM (0x...) address formats
// For backward compatibility, checks both the input address and its normalized form
func (ws *WorldState) GetAccount(addr string) *AccountState {
	// Try direct lookup first (for legacy ARCV-keyed accounts)
	if acct, ok := ws.Accounts[addr]; ok {
		return acct
	}

	// Try normalized form (for new EVM-keyed accounts or when querying with alternate format)
	normAddr, err := NormalizeAddress(addr)
	if err != nil {
		return nil
	}

	return ws.Accounts[normAddr]
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

