# Fork Bitcoin Core and Rebrand into "Archivas Core" (GUI for Archivas Node + Farmer)

## Goal

We are forking the Bitcoin Core repo, but we are not building a Bitcoin client. We want to reuse the mature Qt desktop shell (menus, status bar, debug window, RPC console layout, platform integration) and replace Bitcoin-specific logic with Archivas-specific logic. The final app should be a desktop GUI that can start/stop an Archivas node and an Archivas farmer, display chain status from the Archivas RPC, and show logs.

**Repo to fork from:**
- https://github.com/bitcoin/bitcoin

**Create a new repo:**
- `https://github.com/ArchivasNetwork/archivas-core-gui`

(This is the GUI wrapper, not the Go node itself.)

---

## 1. Strip Bitcoin-specific Identity

In the new repo, do the following renames and removals:

### 1.1 Global Rename

- "Bitcoin Core" → "Archivas Core"
- "bitcoin-qt" → "archivas-qt"
- "Bitcoin" → "Archivas"
- "BTC" tickers and coin units → "RCHV"

### 1.2 Remove/Disable Features We Do Not Need from Bitcoin

- UTXO-specific views
- Coin control
- Mining controls for PoW
- Peers list that shows Bitcoin p2p message types
- All RPC commands that assume Bitcoin's JSON-RPC schema

### 1.3 Keep

- Qt application skeleton
- Main window layout
- Status bar (we will replace the content)
- Debug window (we will repurpose it to show logs from Archivas)
- Settings dialog

The idea is to keep the battle-tested Qt frame and swap the data sources.

---

## 2. Integration Model

Archivas node and farmer are separate Go binaries. We are not rewriting those into C++. The GUI should act as a **supervisor** that:

- checks if `archivas-node` is present on disk
- checks if `archivas-farmer` is present on disk
- can start/stop them as background processes
- can read their stdout/stderr and show it in the GUI
- can poll the local Archivas RPC (HTTP) for status

So in Qt we need a small process manager layer.

### 2.1 Add New C++ Class: `ArchivasProcessManager`

Create `src/qt/archivasprocessmanager.{h,cpp}`:

**Key Methods:**
- `startNode()` - runs something like:
  ```bash
  ./archivas-node --network=archivas-devnet-v4 --rpc=0.0.0.0:8080
  ```
- `startFarmer()` - runs something like:
  ```bash
  ./archivas-farmer farm --node http://127.0.0.1:8080 --plots ./plots --farmer-privkey <key>
  ```
- `stopNode()` - gracefully stops the node process
- `stopFarmer()` - gracefully stops the farmer process
- `isNodeRunning()` - returns bool
- `isFarmerRunning()` - returns bool

**Implementation Details:**
- Uses `QProcess` to manage external processes
- Captures stdout/stderr and emits Qt signals (`nodeOutput(QString)`, `farmerOutput(QString)`)
- Handles process crashes and auto-restart (optional)
- Emits status signals when processes start/stop

**Example Signal:**
```cpp
signals:
    void nodeStarted();
    void nodeStopped();
    void nodeOutput(const QString &line);
    void farmerStarted();
    void farmerStopped();
    void farmerOutput(const QString &line);
```

This replaces Bitcoin's internal node loop with "control an external node."

---

## 3. Replace RPC Layer

Bitcoin Core has its own JSON-RPC server and client expectations. Archivas already serves HTTP JSON at:

- `GET /chainTip`
- `GET /blocks/recent?limit=N`
- `GET /account/<addr>`
- `GET /tx/recent?limit=N`
- `POST /submit`

### 3.1 Create New Module: `ArchivasRpcClient`

Create `src/qt/archivasrpcclient.{h,cpp}`:

**Key Features:**
- Configurable base URL (default `http://127.0.0.1:8080`)
- Uses `QNetworkAccessManager` to call Archivas endpoints
- Parses JSON into simple structs
- Emits signals when data is refreshed

**Data Structures:**
```cpp
struct ChainTip {
    QString height;
    QString hash;
    QString difficulty;
    QString timestamp;
};

struct BlockInfo {
    QString height;
    QString hash;
    QString farmer;
    int txCount;
    QString timestamp;
};

struct TransactionInfo {
    QString hash;
    QString from;
    QString to;
    QString amount;
    QString fee;
    QString height;
};
```

**Key Methods:**
- `getChainTip()` - async call to `/chainTip`
- `getRecentBlocks(int limit)` - async call to `/blocks/recent?limit=N`
- `getRecentTransactions(int limit)` - async call to `/tx/recent?limit=N`
- `getAccount(const QString &address)` - async call to `/account/<addr>`
- `submitTransaction(const QByteArray &txData)` - POST to `/submit`

**Signals:**
```cpp
signals:
    void chainTipUpdated(const ChainTip &tip);
    void blocksUpdated(const QList<BlockInfo> &blocks);
    void transactionsUpdated(const QList<TransactionInfo> &txs);
    void accountUpdated(const AccountInfo &account);
    void error(const QString &message);
```

**Fallback Behavior:**
- If localhost RPC is down, fall back to `https://seed.archivas.ai` in read-only mode
- Show connection status in the UI

---

## 4. GUI Changes to Make It a Real Farmer/Node GUI

Update the main window to have these pages in the left sidebar:

### 4.1 Overview Page

**Display:**
- Chain height (from `/chainTip`)
- Difficulty
- Last block hash
- Node status: Running / Stopped (with indicator)
- Farmer status: Running / Stopped (with indicator)
- Network: archivas-devnet-v4 / archivas-mainnet
- RPC connection status

**Auto-refresh:** Every 5 seconds

### 4.2 Node Page

**Controls:**
- Buttons: "Start Node", "Stop Node", "Restart Node"
- Status indicator (green/red)
- Show node command being used
- Show node executable path (editable)
- Show last 50 lines of node log in a scrollable text area
- Show peer count (if `/peers` endpoint exists later)
- Show sync status (synced / syncing / behind)

**Log Display:**
- Real-time log output from `archivas-node` stdout/stderr
- Color-coded (info/warn/error)
- Auto-scroll to bottom option

### 4.3 Farmer Page

**Controls:**
- Buttons: "Start Farmer", "Stop Farmer", "Restart Farmer"
- Status indicator (green/red)
- Show farmer executable path (editable)
- Show plots path (editable)
- Show farmer private key path (masked, editable)
- Show last qualities submitted
- Show last successful proof
- Show total plots loaded
- Show total plot size (TB)

**Log Display:**
- Real-time log output from `archivas-farmer` stdout/stderr
- Color-coded (info/warn/error)
- Highlight winning proofs

### 4.4 Blocks Page

**Table View:**
- Bound to `/blocks/recent?limit=20`
- Columns:
  - Height
  - Hash (truncated, clickable to copy)
  - Farmer Address
  - Transaction Count
  - Timestamp
  - Difficulty
- Auto-refresh every 10 seconds
- Double-click row to show block details (future)

### 4.5 Transactions Page

**Table View:**
- Bound to `/tx/recent?limit=50`
- Columns:
  - Hash (truncated, clickable to copy)
  - From Address
  - To Address
  - Amount (RCHV)
  - Fee (RCHV)
  - Height
  - Timestamp
- Filter by address (optional)
- Auto-refresh every 10 seconds
- Double-click row to show transaction details (future)

### 4.6 Logs Page

**Two Tabs:**
1. **Node Logs** - QPlainTextEdit bound to node process output
2. **Farmer Logs** - QPlainTextEdit bound to farmer process output

**Features:**
- Clear button for each tab
- Save logs to file
- Search/filter
- Timestamp toggle
- Auto-scroll toggle

This makes it clearly an Archivas app and not a Bitcoin wallet.

---

## 5. Remove Bitcoin Consensus/Wallet Internals from the Build

In `src/` of Bitcoin Core there are a lot of consensus and wallet files. For this fork:

- **Keep:**
  - Qt app (`src/qt/`)
  - Utility libraries (logging, args, fs)
  - Platform integration code

- **Remove/Stub:**
  - Bitcoin validation (`src/consensus/`, `src/validation/`)
  - Bitcoin wallet (`src/wallet/`)
  - Bitcoin P2P (`src/net_processing.cpp`, etc.)
  - Bitcoin RPC server (we use external node)

- **CMake/qmake Changes:**
  - Adjust build system to only compile GUI app
  - Include new `ArchivasProcessManager` and `ArchivasRpcClient`
  - Remove Bitcoin-specific dependencies where possible
  - Keep Qt dependencies

The goal is not to ship Bitcoin's consensus. The goal is to ship a desktop controller for Archivas.

---

## 6. Config and Paths

Add a simple config file under the user's home directory:

**Paths:**
- Linux: `~/.archivas-core/config.json`
- Windows: `%APPDATA%\ArchivasCore\config.json`
- macOS: `~/Library/Application Support/ArchivasCore/config.json`

**Config File Structure:**
```json
{
  "node": {
    "executable_path": "/home/user/archivas/archivas-node",
    "network": "archivas-devnet-v4",
    "rpc_bind": "0.0.0.0:8080",
    "data_dir": "/home/user/.archivas/data",
    "auto_start": false
  },
  "farmer": {
    "executable_path": "/home/user/archivas/archivas-farmer",
    "plots_path": "/home/user/archivas/plots",
    "farmer_privkey_path": "/home/user/.archivas/farmer.key",
    "node_url": "http://127.0.0.1:8080",
    "auto_start": false
  },
  "rpc": {
    "url": "http://127.0.0.1:8080",
    "fallback_url": "https://seed.archivas.ai",
    "poll_interval_ms": 5000
  },
  "ui": {
    "theme": "dark",
    "language": "en",
    "minimize_to_tray": true
  }
}
```

**Settings Dialog:**
- Allow user to edit all paths
- Validate paths exist before saving
- Show file picker for executables and directories
- Mask private key path in UI (show as `****`)

---

## 7. Build Targets

Rename build targets from:

- `bitcoin-qt` → `archivas-qt`
- `bitcoin-cli` → **remove** (not needed)
- `bitcoind` → **remove** (node is external, in Go)

We only need to build the Qt desktop app here.

**Build Output:**
- Linux: `archivas-qt` binary
- Windows: `archivas-qt.exe`
- macOS: `Archivas Core.app`

---

## 8. Platform-Specific Considerations

### 8.1 Linux

- Systemd integration (optional): Create `.desktop` file for launcher
- AppIndicator support for system tray
- Use `QStandardPaths` for config location

### 8.2 Windows

- Windows installer (NSIS or WiX)
- Add to Start Menu
- System tray icon
- Windows service integration (optional, for auto-start)

### 8.3 macOS

- Bundle structure: `Archivas Core.app/Contents/`
- Code signing (for distribution)
- macOS menu bar integration
- Use `QStandardPaths` for config location

---

## 9. Implementation Phases

### Phase 1: Foundation
1. Fork Bitcoin Core repo
2. Global rename (Bitcoin → Archivas)
3. Remove Bitcoin-specific UI elements
4. Create `ArchivasProcessManager` skeleton
5. Create `ArchivasRpcClient` skeleton

### Phase 2: Core Functionality
1. Implement process management (start/stop node/farmer)
2. Implement RPC client (chainTip, blocks, transactions)
3. Update Overview page with real data
4. Add Node and Farmer pages with controls

### Phase 3: Polish
1. Add Blocks and Transactions tables
2. Add Logs page with real-time output
3. Implement config file loading/saving
4. Add Settings dialog
5. Error handling and user feedback

### Phase 4: Distribution
1. Build system cleanup (remove Bitcoin dependencies)
2. Platform-specific packaging
3. Documentation
4. Release builds

---

## 10. Deliverables

1. **New repo initialized from Bitcoin Core**
   - All "Bitcoin" strings renamed to "Archivas"
   - Build system updated

2. **New Classes:**
   - `ArchivasProcessManager` (start/stop Go binaries)
   - `ArchivasRpcClient` (poll Archivas HTTP API)

3. **Updated Main Window:**
   - Overview page with chain status
   - Node page with controls and logs
   - Farmer page with controls and logs
   - Blocks table
   - Transactions table
   - Logs viewer

4. **Config File Support:**
   - JSON config with all paths
   - Settings dialog for editing

5. **Build Instructions:**
   - Linux (CMake/qmake)
   - macOS (Xcode/CMake)
   - Windows (Visual Studio/CMake)

6. **Documentation:**
   - README with setup instructions
   - User guide for GUI features
   - Developer guide for extending

---

## 11. Getting Started with Cursor

Tell Cursor to start by:

1. **Creating the new repo structure**
   - Fork/clone Bitcoin Core
   - Create new repo `archivas-core-gui`
   - Initial commit with Bitcoin Core codebase

2. **Renaming the Qt app and wiring a QProcess-based process manager**
   - Global find/replace: Bitcoin → Archivas
   - Create `ArchivasProcessManager` class
   - Wire up basic start/stop functionality

3. **Adding the Archivas RPC client in Qt**
   - Create `ArchivasRpcClient` class
   - Implement `getChainTip()` first
   - Test with local node

4. **Replacing the overview page with Archivas data**
   - Remove Bitcoin balance/UTXO displays
   - Add chain height, difficulty, node/farmer status
   - Wire up auto-refresh

That will get you a true farmer/node GUI on top of a forked Bitcoin Core UI.

---

## 12. Technical Notes

### 12.1 Process Management

- Use `QProcess::startDetached()` or `QProcess` with proper signal handling
- Capture both stdout and stderr
- Handle process crashes gracefully
- Show error messages if binaries not found

### 12.2 RPC Client

- Use `QNetworkAccessManager` for async HTTP requests
- Implement timeout handling (default 5 seconds)
- Parse JSON with `QJsonDocument` and `QJsonObject`
- Handle network errors and show user-friendly messages

### 12.3 Threading

- RPC calls should be async (non-blocking UI)
- Process output should be captured in separate thread or via signals
- Use Qt's signal/slot mechanism for thread-safe updates

### 12.4 Error Handling

- Validate all user inputs (paths, URLs)
- Show clear error messages
- Log errors to file for debugging
- Graceful degradation (fallback to seed node if local fails)

---

## 13. Future Enhancements

- Plot management UI (list plots, verify integrity)
- Wallet integration (send/receive RCHV)
- Transaction builder (create and sign transactions)
- Network statistics dashboard
- Farmer pool integration (if pools are implemented)
- Dark/light theme toggle
- Multi-language support
- Auto-update mechanism

---

## 14. License and Attribution

- Bitcoin Core is licensed under MIT
- Keep original license and attribution
- Add `THIRD_PARTY.md` documenting the fork
- Clearly state this is a fork/derivative work

---

**End of Specification**


