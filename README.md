# StripChain Node

## Quick Start

To create a complete network with bootnode, 2 SIOs (Signer Interfacing Oracles), and sequencer:

```sh
# First, build the contracts image from the contracts repository
cd /path/to/StripChain/contracts
docker build -t strip-contracts .

# Then return to strip-node repository and run the network
cd /path/to/StripChain/strip-node
./runNetwork.sh
```

## Running with Docker

StripChain is designed to run as a multi-container application using Docker and Docker Compose. The system includes several components that work together:

### Prerequisites

- Docker and Docker Compose installed on your system
- Git (to clone the repositories)
- StripChain/contracts repository (for building the contracts image)

### Available Services

The Docker Compose configuration includes the following services:

1. **ganache**: An Ethereum development blockchain for testing
2. **foundry**: Handles smart contract deployment
3. **bootnode**: P2P network bootstrap node
4. **signerX** (signer1, signer2): Signer nodes that participate in the distributed signing process
5. **signerXpostgres**: PostgreSQL databases for each signer
6. **sequencer**: Processes transactions and communicates with signers
7. **sequencerpostgres**: PostgreSQL database for the sequencer
8. **solver**: Test solver service

### Setup and Usage

1. **Build the contracts image** (must be done in the contracts repository):

   ```sh
   git clone git@github.com:StripChain/contracts.git
   # Navigate to the contracts repository
   cd contracts
   # Build the contracts image
   docker build -t strip-contracts .
   ```

2. **Start the complete network** (in the strip-node repository):

   ```sh
   git clone git@github.com:StripChain/strip-node.git
   # Navigate to the strip-node repository
   cd strip-node
   # Start the network
   docker-compose up -d
   ```

3. **Access the services**:

   - Sequencer API: http://localhost:80
   - Signer1 API: http://localhost:8080
   - Signer2 API: http://localhost:8081
   - Solver API: http://localhost:8083

4. **Check service status**:

   ```sh
   docker-compose ps
   ```

5. **View logs**:

   ```sh
   docker-compose logs -f [service_name]
   ```

6. **Stop the network**:

   ```sh
   docker-compose down
   ```

### Testing

Once the network is running, you can test the functionality by sending requests to the respective service endpoints.

## Testing with Curl

You can test the various services using curl commands. Here are some examples to get you started:

### Testing Sequencer API

**Check Sequencer Status**:

```sh
curl http://localhost:80/status
```

**Create a Wallet**:

```sh
curl "http://localhost:80/createWallet?identity=0x742d35Cc6634C0532925a3b844Bc454e4438f44e&identityCurve=ecdsa"
```

**Get Wallet Information**:

```sh
curl "http://localhost:80/getWallet?identity=0x742d35Cc6634C0532925a3b844Bc454e4438f44e&identityCurve=ecdsa"
```

**Get Bridge Address**:

```sh
curl http://localhost:80/getBridgeAddress
```

**Get an Intent by ID**:

```sh
curl "http://localhost:80/getIntent?id=1"
```

**Get All Intents**:

```sh
curl http://localhost:80/getIntents
```

**Get Intents for a Specific Address**:

```sh
curl "http://localhost:80/getIntentsOfAddress?identity=0x742d35Cc6634C0532925a3b844Bc454e4438f44e&identityCurve=ecdsa"
```

### Testing Signer API

**Check Signer Status**:

```sh
curl http://localhost:8080/status
```

**Get Address Information**:

```sh
curl "http://localhost:8080/address?identity=0x742d35Cc6634C0532925a3b844Bc454e4438f44e&identityCurve=ecdsa&keyCurve=secp256k1"
```

### Testing Solver API

**Check Solver Status**:

```sh
curl http://localhost:8083/status
```

**Construct a Solution**:

```sh
curl http://localhost:8083/construct
```

**Get Solver Output**:

```sh
curl http://localhost:8083/output
```

## Development

For development and testing purposes, you can use the included script to run a complete local network:

```sh
# Make sure you've built the strip-contracts image from the contracts repository first
./runNetwork.sh
```

This script handles starting all the required services using the pre-built Docker images.
