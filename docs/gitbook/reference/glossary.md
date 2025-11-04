# Glossary

Key terms used in Archivas and Proof-of-Space-and-Time blockchains.

---

## A

**Address**  
A Bech32-encoded identifier for an account (e.g., `arcv1t3huuyd08er3yfnmk9c935rmx3wdh5j6m2uc9d`). Derived from a public key using Blake2b-160 hash.

**Archivas**  
The name of this blockchain. Latin for "archives" - reflects the storage-based consensus.

---

## B

**Base Units**  
The smallest denomination of RCHV. 1 RCHV = 100,000,000 base units (8 decimals).

**Bech32**  
Address encoding format. Archivas uses `arcv` prefix (Archivas prefix).

**BIP39**  
Bitcoin Improvement Proposal 39. Standard for generating mnemonics (24-word phrases).

**Block**  
A group of transactions bundled together with a Proof-of-Space and metadata.

**Block Height**  
The sequential number of a block (genesis = 0, next = 1, etc.).

**Block Reward**  
The RCHV earned by the farmer who produces a block (currently 20 RCHV).

---

## C

**Challenge**  
A random value generated each block that farmers must find proofs for. Computed from VDF output.

**Coinbase Transaction**  
The first transaction in every block - pays the block reward to the farmer.

**Consensus**  
The mechanism by which the network agrees on the blockchain state. Archivas uses Proof-of-Space-and-Time.

---

## D

**Difficulty**  
A value that determines how hard it is to find a winning proof. Adjusts automatically to maintain stable block times.

---

## E

**Ed25519**  
Elliptic curve cryptography used for signatures in Archivas. More efficient than secp256k1.

---

## F

**Farmer**  
A participant who creates plots and searches for winning proofs to earn RCHV block rewards.

**Farmer Address**  
The wallet address that receives block rewards when a farmer wins.

---

## G

**Genesis Block**  
The first block in the blockchain (height 0). Contains initial token allocation.

**Genesis Hash**  
The hash of the genesis block. Used to identify which chain nodes are on.

---

## H

**Hash**  
The output of a cryptographic hash function (Blake2b for Archivas). Used for block IDs, addresses, and challenges.

---

## I

**IBD (Initial Block Download)**  
The process of syncing the full blockchain from peers when first starting a node.

---

## K

**k-size**  
The size parameter for plots. k=28 means 2^28 hashes (~268 million). Larger k = larger plots = more chances to win.

---

## M

**Mempool**  
The pool of pending transactions waiting to be included in the next block.

**Mnemonic**  
A 24-word phrase that can recreate a wallet. Standard BIP39 format.

---

## N

**Network ID**  
Identifier for the blockchain network (e.g., `archivas-devnet-v4`). Prevents cross-network communication.

**Nonce**  
A sequential counter for each account's transactions. Prevents replay attacks.

---

## P

**Plot**  
A file containing precomputed hashes used for Proof-of-Space farming.

**Proof (Proof-of-Space)**  
Evidence that a farmer has allocated disk space. Used to win blocks.

**Proof Quality** (or just **Quality**)  
A numeric value derived from a proof. Lower is better. If quality < difficulty, farmer wins the block.

**PoST (Proof-of-Space-and-Time)**  
Consensus mechanism combining Proof-of-Space with Verifiable Delay Functions.

---

## Q

**QMAX**  
The maximum quality value (1 trillion for Archivas). Quality is normalized to 0-QMAX range.

---

## R

**RCHV**  
The native token of Archivas. Pronounced "archive" or "R-C-H-V".

**Reorg (Reorganization)**  
When the chain switches to a different fork. Not common in PoST due to VDF.

---

## S

**Seed Node**  
A publicly accessible node that helps new nodes join the network. Archivas seed: seed.archivas.ai

**SLIP-0010**  
Standard for hierarchical deterministic key derivation, used with Ed25519.

---

## T

**Timelord**  
A process that computes VDF (Verifiable Delay Functions) to generate challenges and provide temporal ordering.

**Transaction**  
A transfer of RCHV from one address to another, signed with a private key.

**Tip**  
The most recent block in the blockchain (the "tip" of the chain).

---

## V

**VDF (Verifiable Delay Function)**  
A function that takes real time to compute but is fast to verify. Used to prevent grinding attacks.

**VDF Iterations**  
The number of sequential hash operations in a VDF. More iterations = more time.

---

## W

**Wallet**  
Software that manages private keys and creates transactions. Can be a CLI tool, SDK, or web app.

**Winning Proof**  
A Proof-of-Space with quality below the difficulty target. Allows farmer to produce the next block.

---

## Technical Terms

**Blake2b**  
Cryptographic hash function used by Archivas. Faster than SHA-256, secure.

**Canonical JSON (RFC 8785)**  
Deterministic JSON serialization. Ensures same hash for same transaction data.

**Challenge**  
`SHA256(VDF_output || height)` - the value farmers must find proofs for.

**Difficulty Target**  
Maximum quality value that wins a block. Lower = harder to win.

**Quality**  
`SHA256(challenge || plot_hash)` - computed for each plot entry.

---

## Acronyms

- **API** - Application Programming Interface
- **CORS** - Cross-Origin Resource Sharing
- **IBD** - Initial Block Download
- **PoS** - Proof-of-Stake
- **PoST** - Proof-of-Space-and-Time
- **PoW** - Proof-of-Work
- **RPC** - Remote Procedure Call
- **SDK** - Software Development Kit
- **TLS** - Transport Layer Security
- **VDF** - Verifiable Delay Function

---

**Can't find a term?** Ask in [GitHub Discussions](https://github.com/ArchivasNetwork/archivas/discussions)!

