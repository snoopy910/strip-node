# StripCNode

## Quick Start

To create a complete network with bootnode, 2 SIOs (Signer Interfacing Oracles), and sequencer:

```sh
# First, build the contracts image from the contracts repository
cd /path/to/contracts
docker build -t strip-contracts .

# Then return to strip-node repository and run the network
cd /path/to/strip-node
./runNetwork.sh
```

## Running with Docker

StripNode is designed to run as a multi-container application using Docker and Docker Compose. The system includes several components that work together:

### Prerequisites

- Docker and Docker Compose installed on your system
- Git (to clone the repositories)
- ./contracts repository (for building the contracts image)

### Available Services

The Docker Compose configuration includes the following services:

1. **ganache**: An Ethereum development blockchain for testing
2. **foundry**: Handles smart contract deployment
3. **bootnode**: P2P network bootstrap node
4. **validatorX** (validator1, validator2): Signer nodes that participate in the distributed signing process
5. **signerXpostgres**: PostgreSQL databases for each signer
6. **sequencer**: Processes transactions and communicates with signers
7. **sequencerpostgres**: PostgreSQL database for the sequencer
8. **solver**: Test solver service

### Setup and Usage

1. **Build the contracts image** (must be done in the contracts repository):

   ```sh
   git clone git@github.com:snoopy910/contracts.git
   # Navigate to the contracts repository
   cd contracts
   # Build the contracts image
   docker build -t strip-contracts .
   ```

2. **Start the complete network** (in the strip-node repository):

   ```sh
   git clone git@github.com:snoopy910/strip-node.git
   # Navigate to the strip-node repository
   cd strip-node
   # Start the network
   docker-compose up -d
   ```

3. **Access the services**:

   - Sequencer API: http://localhost:80
   - Validator1 API: http://localhost:8080
   - Validator2 API: http://localhost:8081
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

## Bitcoin Support

StripChain includes a full Bitcoin development environment for testing and integration. The setup features:

- Bitcoin Core (bitcoind) running in regtest mode
- Electrs (Electrum Rust Server) for REST API access
- Utility scripts for common Bitcoin operations

### Quick Bitcoin Guide

1. **Start Bitcoin Services**:
   ```sh
   docker-compose up -d bitcoind electrs
   ```

2. **Fund the Test Wallet**:
   ```sh
   cd bitcoin-scripts
   ./fund-test-wallet.sh
   ```

3. **Check Address Balance**:
   ```sh
   ./check-address-balance.sh <bitcoin_address>
   ```

4. **Send Test Bitcoins**:
   ```sh
   ./send-bitcoins.sh <address> <amount>
   ```

5. **Generate Blocks**:
   ```sh
   ./generate-blocks.sh <number_of_blocks>
   ```

For complete documentation of the Bitcoin environment and available utilities, see the [Bitcoin README](BITCOIN-README.md).

## Running Uniswap V3 Solver

To run the Uniswap V3 solver, use the following command:

```sh
go run ./main.go --isUniswapSolver=true \
    --rpcURL=<rpc-url> \
    --httpPort=<port> \
    --uniswapV3Factory=<uniswap-v3-factory-address> \
    --npmAddress=<npm-address> \
    --chainId=<chain-id>
```

Required parameters:
- `rpcURL`: Ethereum RPC URL for the network
- `httpPort`: Port number for the solver's HTTP server
- `uniswapV3Factory`: Address of the Uniswap V3 Factory contract
- `npmAddress`: Address of the Uniswap V3 NonfungiblePositionManager contract
- `chainId`: Chain ID of the network

## Docker Volumes

StripChain uses Docker volumes to persist data across container restarts. The following volumes are created:

1. **postgres-data-X**: Volumes for each PostgreSQL database (sequencer and signers)
2. **bitcoin-data**: Contains the Bitcoin blockchain data, wallet information, and configurations
3. **electrs-data**: Stores the Electrs database and index information

These volumes ensure that your data is preserved even when containers are stopped or restarted. You can view all volumes with:

```sh
docker volume ls | grep strip-node
```

## Resetting the System

There are several approaches to reset the system depending on your needs:

### Soft Reset (Preserve Data)

To restart all services while preserving data in volumes:

```sh
docker-compose down
docker-compose up -d
```

### Reset Database Only

To reset just the PostgreSQL databases while preserving other data:

```sh
docker-compose down
docker volume rm strip-node_postgres-data-sequencer strip-node_postgres-data-validator1 strip-node_postgres-data-validator2
docker-compose up -d
```

### Reset Bitcoin Data Only

To reset just the Bitcoin blockchain data (useful during testing):

```sh
docker-compose down bitcoind electrs
docker volume rm strip-node_bitcoin-data strip-node_electrs-data
docker-compose up -d bitcoind electrs
# Refund the wallet after reset
cd bitcoin-scripts
./fund-test-wallet.sh
```

### Full Reset

To completely reset the system and all data:

```sh
docker-compose down
docker volume rm $(docker volume ls -q | grep strip-node)
docker-compose up -d
```

This will remove all data and start with a clean system. You'll need to re-initialize wallets and other state after a full reset.
