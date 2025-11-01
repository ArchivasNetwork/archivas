# v1.1.1 Hotfix — POST /submit Method Handling

## Issue
Tester reported HTTP 405 "Method not allowed" when submitting transactions via `./archivas-cli`. The node was not correctly handling POST requests to `/submit`.

## Fixes

### 1. Method Validation
- **Fixed**: `/submit` now returns HTTP 405 with `Allow: POST` header for non-POST methods
- **Before**: Generic "Method not allowed" without `Allow` header
- **After**: Proper 405 response with `Allow: POST` header

### 2. Content-Type Validation
- **Fixed**: `/submit` now requires `Content-Type: application/json`
- **Before**: No Content-Type validation
- **After**: Returns HTTP 415 "Unsupported Media Type" if Content-Type is missing or incorrect

### 3. Compatibility Alias
- **Fixed**: POST `/broadcast` now routes to same handler as `/submit` when Content-Type is `application/json`
- **Before**: `/broadcast` only handled legacy `ledger.Transaction` format
- **After**: `/broadcast` accepts both v1.1.0 format (routes to `handleSubmitV1`) and legacy format

### 4. CLI Updates
- **Fixed**: CLI now uses `http.Client` with `CheckRedirect` to prevent POST→GET conversions
- **Before**: Used `http.Post` which could follow redirects
- **After**: Explicitly prevents redirects that would convert POST to GET

## Changes Made

### Files Modified
- `rpc/farming.go`:
  - `handleSubmitV1`: Added 405 with `Allow: POST` header, Content-Type validation (415)
  - `handleBroadcast`: Routes to `handleSubmitV1` for v1.1.0 format
- `cmd/archivas-cli/main.go`:
  - `broadcast`: Updated to use `http.Client` with redirect prevention

### Files Added
- `rpc/submit_hotfix_test.go`: Comprehensive tests for all fixes

## Tests

All tests pass:
```bash
go test ./rpc/... -run "TestSubmit|TestBroadcast" -v
```

Tests cover:
- ✅ GET `/submit` → 405 with `Allow: POST`
- ✅ POST `/submit` without Content-Type → 415
- ✅ POST `/submit` with valid JSON → proper handling (no 405/415)
- ✅ POST `/broadcast` with v1.1.0 format → routes correctly

## Backward Compatibility

✅ **Fully backward compatible**:
- Existing `/submit` endpoint behavior unchanged (only adds validation)
- Legacy `/broadcast` endpoint still works for old `ledger.Transaction` format
- No breaking changes to API, metrics, or consensus

## Deployment

```bash
# Build
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-cli ./cmd/archivas-cli

# Test locally
go test ./rpc/... -v

# Deploy to server (follow DEPLOY-v1.1.0.md steps)
```

## Verification

After deployment, verify:

```bash
# Test GET /submit (should return 405 with Allow: POST)
curl -v http://localhost:8080/submit

# Test POST /submit without Content-Type (should return 415)
curl -v -X POST http://localhost:8080/submit -d '{}'

# Test POST /submit with valid JSON (should work)
curl -v -X POST http://localhost:8080/submit \
  -H "Content-Type: application/json" \
  -d @tx.json

# Test POST /broadcast with v1.1.0 format (should work)
curl -v -X POST http://localhost:8080/broadcast \
  -H "Content-Type: application/json" \
  -d @tx.json
```

## Release Tag

**v1.1.1-hotfix**

