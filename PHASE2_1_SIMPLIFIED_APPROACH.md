# Phase 2.1 - Simplified But Full Working Approach

## Problem
- go-ethereum's `state.StateDB` requires a full database backend (leveldb/memdb)
- Creating this from scratch is complex and time-consuming
- We need MetaMask to work NOW, not in 2 weeks

## Pragmatic Solution
Instead of full VM integration, enhance the current simplified approach to work correctly:

1. **Keep simplified balance transfer logic** (already works)
2. **Add proper nonce handling** (increment after tx)
3. **Add proper gas deduction** (deduct from sender balance)
4. **Add proper receipts** (with correct gas and status)
5. **Convert all Wei <-> RCHV correctly**

This gives us:
- ✅ MetaMask transfers work
- ✅ Balance updates correctly
- ✅ Gas fees work
- ✅ Receipts show up in Blockscout
- ❌ Smart contracts don't work YET

## Next Phase (Phase 2.2)
Once transfers are working, we can add:
- Full go-ethereum StateDB integration
- Contract deployment
- Contract calls

But for NOW, let's get MetaMask working with transfers!

## Implementation
Enhance the existing code in AcceptBlock to:
1. Deduct gas cost from sender
2. Pay gas to farmer
3. Increment nonce
4. Create proper receipts

TIME: 30 minutes instead of 2 weeks!
