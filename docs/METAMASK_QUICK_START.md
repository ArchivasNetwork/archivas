# MetaMask Quick Start - Archivas Betanet

Connect your MetaMask wallet to Archivas Betanet in 3 simple steps.

---

## 📋 **Step 1: Add Archivas Network**

1. Open **MetaMask**
2. Click network dropdown → **Add Network** → **Add manually**
3. Enter:
   - **Network Name:** `Archivas Betanet`
   - **RPC URL:** `https://seed3.betanet.archivas.ai`
   - **Chain ID:** `1644`
   - **Currency Symbol:** `RCHV`
4. Click **Save**

---

## 🔑 **Step 2: Generate Ethereum-Compatible Wallet**

⚠️ **Important:** Archivas wallets must be generated using Ethereum standards for MetaMask compatibility.

### **Option A: Use Our Tool**

```bash
cd ~/archivas
go build -o eth-wallet cmd/eth-wallet/main.go
./eth-wallet
```

This generates:
- **Private Key** (64 hex characters)
- **Ethereum Address** (0x...)

### **Option B: Use MetaMask Directly**

Create a new account in MetaMask and export the private key.

---

## 💰 **Step 3: Import Your Wallet**

1. In MetaMask, click **account icon** (top right)
2. Select **Import Account**
3. Choose **Private Key**
4. Paste your **Ethereum-compatible private key**
5. Click **Import**

Your RCHV balance will appear automatically! 🎉

---

## ✅ **Verify Connection**

Check your balance:
```bash
curl -s https://seed3.betanet.archivas.ai -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc":"2.0",
    "method":"eth_getBalance",
    "params":["YOUR_0x_ADDRESS","latest"],
    "id":1
  }'
```

---

## 🚨 **Common Issues**

### **MetaMask shows 0 balance but I have RCHV**

**Cause:** Your RCHV is on an old Archivas-derived address, not Ethereum-compatible.

**Fix:** 
1. Generate a new Ethereum-compatible wallet (see Step 2)
2. Update your farmer to use the new private key
3. Future rewards will be MetaMask-compatible!

---

## 🔗 **Resources**

- Full Guide: [docs/METAMASK_GUIDE.md](METAMASK_GUIDE.md)
- RPC Endpoint: `https://seed3.betanet.archivas.ai`
- Chain Explorer: Coming soon
- Support: GitHub Issues

---

## 🎯 **Key Points**

✅ Use `eth-wallet` tool or MetaMask to generate wallets  
✅ Standard Ethereum private keys work perfectly  
✅ Import private key → MetaMask shows balance automatically  
✅ Use Seed3 RPC for public access  

---

**Need help?** Open an issue on GitHub or check the full documentation.

