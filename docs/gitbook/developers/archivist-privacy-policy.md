---
description: Privacy policy for the Archivist browser wallet extension
---

# Archivist Wallet Privacy Policy

_Last updated: November 9, 2025_

Archivist – Archivas Wallet is a self-custodial browser extension. It is designed so your wallet data never leaves your device unless you deliberately broadcast a transaction to the Archivas network. This page explains how the extension handles data, storage, permissions, and user controls.

## Data Collection

* The extension does not collect or transmit personally identifiable information.
* No registration is required. Wallet addresses, balances, mnemonics, or private keys are never sent to Archivas-operated servers.
* All key material remains on the user’s device. The extension only accesses it locally to sign transactions client-side.

## Local Storage

* Sensitive data — private keys, mnemonics, imported keys, account metadata, recent activity, and UI preferences — is stored in the browser’s extension storage.
* Stored data is encrypted at rest using the user’s password-derived key.
* Archivist does not sync data to any cloud service. Removing the extension erases this local data.

## Network Requests

* Archivist communicates exclusively with Archivas RPC endpoints such as `https://seed.archivas.ai` or any custom node the user configures.
* Requests are limited to fetching account balances, chain tip height, nonce/history, and submitting signed transactions.
* The extension does not call analytics, advertising, or third-party tracking services.

## Permission Usage

* `contextMenus` — provides a shortcut to open the wallet in a new tab.
* `storage` and `unlimitedStorage` — store encrypted wallet data and settings.
* `notifications` — optionally alert users about transaction status or reminders.
* `content_scripts` and `host_permissions` — inject the in-page provider so decentralized apps can request signatures. No page content is scraped or analyzed.

## Data Sharing

* Archivist does not share user data with third parties.
* When a user submits a transaction, the signed payload becomes part of the public Archivas blockchain, just like any other blockchain transaction.

## Security Practices

* Mnemonics and private keys never leave the device; they are decrypted only in memory for signing.
* Transaction signing happens locally in the extension before any data is broadcast.
* All communications with RPC endpoints use HTTPS.

## User Controls

* Users can delete the extension at any time to remove all locally stored data.
* Within the extension, the “Reset wallet” option (if enabled) clears stored key material and forces a fresh setup.
* Users may export mnemonics or private keys for backup; safeguarding exports is the user’s responsibility.

For questions about this policy, contact the Archivas team via the official support channels.

