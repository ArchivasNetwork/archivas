# Setting Up a Farmer on Windows

Step-by-step guide to start farming Archivas on Windows.

---

## Overview

You'll need to:
1. Install Go and Git
2. Build the binaries
3. Create plots
4. Start the farmer
5. Earn RCHV!

**Time:** ~30 minutes  
**Difficulty:** Intermediate

---

## Step 1: Install Dependencies

### Install Go

1. **Download Go for Windows:**
   - Visit: https://go.dev/dl/
   - Download: `go1.24.0.windows-amd64.msi` (or latest version)

2. **Run the installer:**
   - Double-click the `.msi` file
   - Follow the installation wizard
   - Go will be installed to `C:\Program Files\Go`

3. **Verify installation:**
   - Open **Command Prompt** or **PowerShell**
   - Run: `go version`
   - You should see: `go version go1.24.0 windows/amd64`

### Install Git (if not already installed)

1. **Download Git for Windows:**
   - Visit: https://git-scm.com/download/win
   - Download and run the installer
   - Use default settings

2. **Verify installation:**
   - Open **Command Prompt** or **PowerShell**
   - Run: `git --version`

---

## Step 2: Clone and Build

### Open PowerShell or Command Prompt

Press `Win + X` and select **Windows PowerShell** or **Command Prompt**.

### Clone the Repository

```powershell
# Navigate to your desired directory (e.g., Documents)
cd C:\Users\YourUsername\Documents

# Clone repository
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas
```

### Build Binaries

```powershell
# Build farmer
go build -o archivas-farmer.exe ./cmd/archivas-farmer

# Build CLI (optional, for wallet operations)
go build -o archivas-cli.exe ./cmd/archivas-cli

# Verify
.\archivas-farmer.exe --help
.\archivas-cli.exe --help
```

---

## Step 3: Create Your Farmer Identity

The farmer binary manages its own identity keys using **secp256k1** (32-byte private key, 33-byte compressed public key).  
Let the binary generate the keys for you and save them before plotting continues.

```powershell
# Create a directory for plots
mkdir C:\Users\YourUsername\archivas-plots

# Let the farmer generate its own identity (secp256k1) and start your first plot
.\archivas-farmer.exe plot `
  --size 28 `
  --path C:\Users\YourUsername\archivas-plots

# The command immediately prints something like:
#  Generated new farmer identity:
#    Address:     arcv1...
#    Public Key:  02ab... (use for --farmer-pubkey)
#    Private Key: 1f2c... (use for --farmer-privkey)
#  ⚠️  Save both keys! You'll need the private key to farm.
#
# Copy the Address, Public Key, and Private Key to a safe place (Notepad or password manager).
# After the message, the same command continues generating plot-k28.arcv in your plots directory.
```

**Important Notes:**
- Use backticks (`) for line continuation in PowerShell, or put everything on one line
- In Command Prompt, use `^` for line continuation instead of backticks
- Do not mix these keys with the Ed25519 keys produced by `archivas-cli`
- The farmer currently requires the secp256k1 keys printed by the step above

If you only needed to record the keys (for example, to plot on another machine), you can press `Ctrl+C` once you have copied them, then rerun the command later with `--farmer-pubkey <saved_key>` to start plotting.

---

## Step 4: Create Plots

```powershell
# Create additional k28 plots (use the public key captured in Step 3)
.\archivas-farmer.exe plot `
  --size 28 `
  --path C:\Users\YourUsername\archivas-plots `
  --farmer-pubkey YOUR_PUBKEY_FROM_STEP_3

# Each run creates a file such as C:\Users\YourUsername\archivas-plots\plot-k28.arcv
# Repeat to add more plots (the farmer auto-increments filenames):
.\archivas-farmer.exe plot --size 28 --path C:\Users\YourUsername\archivas-plots --farmer-pubkey YOUR_PUBKEY_FROM_STEP_3
```

**Notes:**
- The `--path` flag points to a directory; the farmer names the plot automatically (`plot-k28.arcv`, `plot-k28-1.arcv`, ...)
- You can run multiple plot commands in parallel if you have CPU and disk headroom
- Plotting can take 10-30 minutes per k28 plot depending on your CPU

---

## Step 5: Start Farming

### Option A: Run in PowerShell (Foreground)

```powershell
# Create logs directory
mkdir C:\Users\YourUsername\archivas-logs

# Start farmer (runs in foreground, press Ctrl+C to stop)
.\archivas-farmer.exe farm `
  --plots C:\Users\YourUsername\archivas-plots `
  --node https://seed.archivas.ai `
  --farmer-privkey YOUR_PRIVKEY_FROM_STEP_3 `
  > C:\Users\YourUsername\archivas-logs\farmer.log 2>&1
```

### Option B: Run in Background (PowerShell)

```powershell
# Start farmer in background
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd C:\Users\YourUsername\Documents\archivas; .\archivas-farmer.exe farm --plots C:\Users\YourUsername\archivas-plots --node https://seed.archivas.ai --farmer-privkey YOUR_PRIVKEY_FROM_STEP_3 > C:\Users\YourUsername\archivas-logs\farmer.log 2>&1"
```

### Option C: Create a Batch File

Create a file `start-farmer.bat` in the `archivas` directory:

```batch
@echo off
cd /d "%~dp0"
.\archivas-farmer.exe farm ^
  --plots C:\Users\YourUsername\archivas-plots ^
  --node https://seed.archivas.ai ^
  --farmer-privkey YOUR_PRIVKEY_FROM_STEP_3
pause
```

Double-click `start-farmer.bat` to start farming.

**Tip:** The `--farmer-privkey` value is the 64-character hex string printed in Step 3 (32-byte secp256k1 key). If you see "must be 32 bytes", double-check that you copied the key from the farmer output and not the Ed25519 key from `archivas-cli`.

**Expected output:**
```
👨‍🌾 Farmer Address: arcv1...
📁 Plots Directory: C:\Users\YourUsername\archivas-plots
🌐 Node: https://seed.archivas.ai

✅ Loaded 3 plot(s)
   - plot-k28-1.arcv (k=28, 268435456 hashes)
   - plot-k28-2.arcv (k=28, 268435456 hashes)
   - plot-k28-3.arcv (k=28, 268435456 hashes)

🚜 Starting farming loop...
🔍 NEW HEIGHT 64500 (difficulty: 1000000)
⚙️  Checking plots...
```

---

## Step 6: Monitor Your Earnings

### Check Your Balance

Open PowerShell and run:

```powershell
# Replace YOUR_ADDRESS with your farmer address from Step 3
Invoke-WebRequest -Uri "https://seed.archivas.ai/account/YOUR_ADDRESS" | Select-Object -ExpandProperty Content
```

Or use a web browser:
```
https://seed.archivas.ai/account/YOUR_ADDRESS
```

### Watch for Wins

Open the log file in Notepad or use PowerShell:

```powershell
# Watch logs in real-time (PowerShell)
Get-Content C:\Users\YourUsername\archivas-logs\farmer.log -Wait -Tail 50
```

**Expected when you win:**
```
🎉 Found winning proof! Quality: 123456 (target: 1000000)
✅ Block submitted successfully for height 64501
```

---

## Running as a Windows Service (Advanced)

### Using NSSM (Non-Sucking Service Manager)

1. **Download NSSM:**
   - Visit: https://nssm.cc/download
   - Download the latest release (e.g., `nssm-2.24.zip`)
   - Extract to `C:\nssm`

2. **Install as Service:**
   ```powershell
   # Open PowerShell as Administrator
   cd C:\nssm\win64
   
   .\nssm.exe install ArchivasFarmer
   ```

3. **Configure Service:**
   - **Path:** `C:\Users\YourUsername\Documents\archivas\archivas-farmer.exe`
   - **Startup directory:** `C:\Users\YourUsername\Documents\archivas`
   - **Arguments:** `farm --plots C:\Users\YourUsername\archivas-plots --node https://seed.archivas.ai --farmer-privkey YOUR_PRIVKEY`
   - **Output:** `C:\Users\YourUsername\archivas-logs\farmer.log`
   - **Error:** `C:\Users\YourUsername\archivas-logs\farmer-error.log`

4. **Start Service:**
   ```powershell
   .\nssm.exe start ArchivasFarmer
   ```

5. **Check Status:**
   ```powershell
   .\nssm.exe status ArchivasFarmer
   ```

---

## Troubleshooting

### "No plots found"

**Problem:** Farmer can't find plot files.

**Solution:**
```powershell
# Check plots exist
dir C:\Users\YourUsername\archivas-plots

# Verify path in farmer command (use forward slashes or escaped backslashes)
.\archivas-farmer.exe farm --plots C:/Users/YourUsername/archivas-plots ...
```

### "Connection refused" or Network Errors

**Problem:** Can't reach node.

**Solution:**
- Verify internet connection
- Test node: `Invoke-WebRequest -Uri "https://seed.archivas.ai/chainTip"`
- Check Windows Firewall isn't blocking HTTPS
- Try disabling antivirus temporarily to test

### "Invalid proof"

**Problem:** Proof rejected by node.

**Solution:**
- Ensure plots were created with correct farmer public key
- Check logs for specific error message
- Verify difficulty target

### Not winning blocks

**Expected behavior** if:
- Network has much more space than you
- You have small plots (k=27 or smaller)
- Bad luck (probability-based)

**Check:**
- How much space do you have vs network total?
- Are plots loading correctly?
- Is farmer scanning on each challenge?

### Antivirus Blocking

Some antivirus software may flag the farmer executable as suspicious.

**Solution:**
- Add `C:\Users\YourUsername\Documents\archivas` to antivirus exclusions
- Or temporarily disable real-time protection during setup

---

## Performance Optimization

### Faster Plot Scanning

```powershell
# Use all CPU cores (PowerShell)
$cores = (Get-WmiObject Win32_ComputerSystem).NumberOfLogicalProcessors
.\archivas-farmer.exe farm `
  --plots C:\Users\YourUsername\archivas-plots `
  --threads $cores `
  --node https://seed.archivas.ai `
  --farmer-privkey YOUR_PRIVKEY
```

### Multiple Plot Directories

```powershell
# Combine multiple directories (comma-separated)
.\archivas-farmer.exe farm `
  --plots D:\plots,E:\plots,F:\plots `
  --node https://seed.archivas.ai `
  --farmer-privkey YOUR_PRIVKEY
```

### Reduce I/O

- Use SSD for plots (if possible)
- Avoid network-mounted storage (Windows network drives)
- Keep plots on local filesystem
- Disable Windows Defender real-time scanning on plot directories

---

## Security

### Protect Your Private Key

**Never commit private keys to git:**
```powershell
# Add to .gitignore
echo "*.key" >> .gitignore
echo "farmer.key" >> .gitignore
```

**Store in environment variable (PowerShell):**
```powershell
# Set environment variable (current session only)
$env:FARMER_PRIVKEY = "your_hex_key_here"

# Use in command
.\archivas-farmer.exe farm `
  --plots C:\Users\YourUsername\archivas-plots `
  --node https://seed.archivas.ai `
  --farmer-privkey $env:FARMER_PRIVKEY
```

**For permanent environment variable:**
1. Press `Win + X` → **System** → **Advanced system settings**
2. Click **Environment Variables**
3. Under **User variables**, click **New**
4. Variable name: `FARMER_PRIVKEY`
5. Variable value: `your_hex_key_here`
6. Click **OK**

---

## Scaling Up

### Add More Plots

```powershell
# Create additional plots (PowerShell loop)
for ($i=4; $i -le 10; $i++) {
    .\archivas-farmer.exe plot `
      --size 28 `
      --path C:\Users\YourUsername\archivas-plots `
      --farmer-pubkey YOUR_PUBKEY
}
```

### Monitor Performance

```powershell
# Watch farmer logs (PowerShell)
Get-Content C:\Users\YourUsername\archivas-logs\farmer.log -Wait -Tail 20 | Select-String -Pattern "Found winning|Quality|NEW HEIGHT"

# Check balance growth (PowerShell, refresh every 30 seconds)
while ($true) {
    Invoke-WebRequest -Uri "https://seed.archivas.ai/account/YOUR_ADDRESS" | Select-Object -ExpandProperty Content
    Start-Sleep -Seconds 30
}
```

---

## Quick Reference

### Common Commands

```powershell
# Build binaries
go build -o archivas-farmer.exe ./cmd/archivas-farmer

# Create plot
.\archivas-farmer.exe plot --size 28 --path C:\Users\YourUsername\archivas-plots

# Start farming
.\archivas-farmer.exe farm --plots C:\Users\YourUsername\archivas-plots --node https://seed.archivas.ai --farmer-privkey YOUR_PRIVKEY

# Check balance
Invoke-WebRequest -Uri "https://seed.archivas.ai/account/YOUR_ADDRESS"
```

### File Locations

- **Binaries:** `C:\Users\YourUsername\Documents\archivas\`
- **Plots:** `C:\Users\YourUsername\archivas-plots\`
- **Logs:** `C:\Users\YourUsername\archivas-logs\`

---

## Next Steps

- [Creating Plots](../farmers/creating-plots.md) - Detailed plotting guide
- [Running a Node](../farmers/running-node.md) - Optional: run your own node
- [Earnings Guide](../farmers/earnings.md) - Understand rewards

---

**Start farming!** You're ready to earn RCHV on Windows! 🌾

