# P2P Peer Tracking Fix - Implementation Plan

## Problem

Server A broadcasts to "0 peers" even though Server B is connected.
The peers map is empty despite TCP connections existing.

## Root Cause

- Accepted connections not properly registered in peers map
- Peer deregistration happening prematurely  
- No thread-safe snapshot for broadcasting

## Solution

1. Fix peer registration on accept
2. Add proper handshake validation
3. Use snapshot for broadcasting
4. Add /peers RPC endpoint for debugging
5. Better logging

## Implementation

See commits following this plan.

## Testing

After implementation:
- curl http://localhost:8080/peers should show 1+ peers
- Logs should show "broadcasted block X to 1 peers" (not 0!)
- Server B should receive blocks and sync

## Status

Starting implementation now...
