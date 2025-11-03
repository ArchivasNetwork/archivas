# Changelog

All notable changes to the Archivas project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [v1.2.0] - 2025-11-03

### Added
- **Explorer listing endpoints** (read-only, non-consensus)
  - `GET /blocks/recent?limit=N` - List recent blocks (default 20, max 100)
  - `GET /block/<height>` - Get block details by height
  - `GET /tx/recent?limit=N` - List recent transactions (default 50, max 200)
- **API improvements**
  - All numeric fields returned as strings for consistency
  - Proper error handling for invalid heights/limits
  - CORS-enabled for public access

### Changed
- **Terminology:** Renamed `miner` → `farmer` in block JSON responses
  - Reflects Proof-of-Space-and-Time terminology
  - `miner` kept as deprecated alias until v1.3 (backward compatibility)
  - Primary field is now `farmer`

### Fixed
- **CORS duplication:** Fixed "Access-Control-Allow-Origin: *, *" browser error
  - Nginx no longer adds CORS headers
  - Backend (RPC middleware) solely responsible for CORS
  - Single, clean CORS header set

### Notes
- **No protocol changes:** No modifications to consensus, storage, or validation
- **Backward compatible:** All existing endpoints unchanged
- This is an API-additive release for explorer support

---

## [v1.1.1-infra] - 2025-11-02

### Added
- **Seed node infrastructure** at `https://seed.archivas.ai`
  - Nginx reverse proxy with TLS termination
  - Let's Encrypt certificate (auto-renewing)
  - HTTP/2 enabled for performance
  - CORS configured for public API access
  - Rate limiting on `/submit` endpoint (10 req/min per IP)
  - Security headers (HSTS, X-Frame-Options, X-Content-Type-Options, Referrer-Policy)
- **Deployment scripts**
  - `scripts/setup-seed-nginx.sh` - Idempotent installer
  - `scripts/renew-cert-hook.sh` - Certificate renewal hook
  - `scripts/check-seed.sh` - Health validation
- **Documentation**
  - `docs/SEED_HOST.md` - Complete seed node setup guide
  - README updated with public RPC endpoint examples
  - Firewall and security configuration documented

### Security
- Node RPC bound to `127.0.0.1:8080` only (not externally accessible)
- `/metrics` endpoint blocked from public access (internal monitoring only)
- Firewall configured to deny external access to port 8080
- Rate limiting on transaction submission
- TLS 1.2+ with modern cipher suites only

### Infrastructure
- Public RPC endpoint: `https://seed.archivas.ai`
- Certificate expires: 2026-01-31 (auto-renewal enabled)
- HTTP redirects to HTTPS
- Nginx serves as reverse proxy to localhost:8080

### Notes
- **No protocol changes:** No modifications to consensus, database schema, or RPC response formats
- **No API changes:** Fully compatible with v1.1.0 wallet API
- This is an infrastructure-only release

---

## [v1.1.0] - 2025-11-01

### Added
- **Wallet Primitives + Public API Freeze**
- Ed25519 keypairs with BIP39 + SLIP-0010 derivation
- Bech32 addresses (`arcv` prefix)
- Transaction v1 schema with RFC 8785 canonical JSON
- Blake2b-256 hashing with domain separation (`Archivas-TxV1`)
- Complete RPC API for wallet operations
- TypeScript SDK ready for publication

### API Endpoints
- `GET /account/<address>` - Account state
- `GET /chainTip` - Chain tip info
- `GET /mempool` - Pending transactions
- `GET /tx/<hash>` - Transaction details
- `GET /estimateFee?bytes=<n>` - Fee estimation
- `POST /submit` - Submit signed transaction

### Packages
- `pkg/crypto/` - BIP39, SLIP-0010, address codec
- `pkg/tx/v1/` - Transaction model and signing
- `cmd/archivas-cli/` - Wallet CLI tools

---

## Earlier Versions

See Git history for v1.0.x releases (observability, IBD, PoSpace fixes, etc.)

