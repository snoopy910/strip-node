Relationship Between the Components in the TEE Integration
Let me clarify the relationships between the components in our new TEE-integrated architecture.
Original Architecture vs New Architecture
Original System
In the original signer implementation:

API endpoints in api.go triggered operations
Key shares were generated through generateKeygen() and stored directly in PostgreSQL via db.go
Multiple nodes communicated via P2P network to run distributed TSS protocols
Each node operator had access to the key shares in their database

New TEE-Protected System
We've now created a layered architecture with three main components:

Parent Instance (regular Go application)
TEE Package (interface/communication layer)
Enclave Application (protected environment)

Component Relationships
1. Parent Instance
This is our familiar signer code, but modified to delegate sensitive operations:

Modified api.go: Same HTTP endpoints, but now delegates sensitive operations to the TEE client
Modified db.go: Now uses a TEE client for key operations and only stores non-sensitive metadata locally
Modified keygen.go: Simplified to delegate to TEE client for key generation
Modified signature.go: Simplified to delegate to TEE client for signing
P2P network: Still handles inter-node communication, but now it transfers messages between enclaves

The parent instance is the public-facing component that receives requests, validates them, and coordinates with the TEE.

2. TEE Package
This is a new interface layer that connects the parent instance to the enclave:

tee/client.go: Defines the standard interface for TEE operations
tee/mock/client.go: Mock implementation for development and testing
tee/nitro/client.go: Real implementation using AWS Nitro Enclaves

This layer abstracts away the details of how to communicate with the enclave, letting the parent instance use a consistent API regardless of the TEE technology.

3. Enclave Application
This is the protected environment where sensitive operations occur:

enclave/main.go: Entry point that sets up vsock server to listen for parent requests
enclave/keygen.go: Implements actual key generation inside the TEE
enclave/signing.go: Implements actual signing operations inside the TEE
enclave/kms.go: Handles key sealing/unsealing with AWS KMS
enclave/keystore.go: Manages sealed key storage

The enclave contains all the sensitive cryptographic operations and key material.
Data and Control Flow

Key Generation Flow:

Client → HTTP API (parent) → TEE Client → Enclave → Sealed Storage → Parent → Client
The parent validates the request and passes it to the TEE
The enclave generates key material, seals it, and returns a reference
The parent stores non-sensitive metadata (like signers list)


Signing Flow:

Client → HTTP API (parent) → TEE Client → Enclave → Signing Operation → Parent → P2P Network
The parent validates the request and passes it to the TEE
The enclave unseals the key, signs the message, and returns the signature
The parent broadcasts the signature via P2P and returns it to the client



TEE-to-TEE Communication
This is the most interesting part. In the original design, nodes communicated directly via P2P to run the distributed TSS protocol. In the TEE design:

When Enclave A needs to send a message to Enclave B:

Enclave A sends message to its Parent A via vsock
Parent A broadcasts message via existing P2P network to Parent B
Parent B forwards message to its Enclave B via vsock
Enclave B processes the message and responds similarly

The parents act as message proxies, but they can't see the actual key material because:

Messages between enclaves are encrypted
Key material is sealed using KMS with attestation (can only be unsealed by same enclave)


This architecture preserves the distributed security of TSS while adding hardware-level protection for key material. The node operator can still run the service but cannot extract or misuse the protected key shares.
