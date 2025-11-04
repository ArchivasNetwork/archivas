# Network Protocol

Archivas P2P networking.

## Messages

- NEW_BLOCK: Block announcement
- GET_STATUS: Peer status request
- BLOCK_DATA: Block sync
- GOSSIP_PEERS: Peer discovery

## Handshake

Validates genesis hash and network ID before accepting peers.

Port: 9090 (default)
