# Bitcoin Development Environment

This setup provides a complete Bitcoin development environment for the StripChain project, featuring Bitcoin Core (bitcoind), Electrs (Electrum Rust Server), and a suite of utility scripts for common operations.

## Components

### 1. Bitcoin Core (bitcoind)
- **Image:** blockstream/bitcoind:latest
- **Network Mode:** Regtest (local development network)
- **Ports:**
  - RPC: 8332 (Credentials: bitcoin:bitcoin)
  - P2P: 8333
  - ZMQ: 28332 (rawblock), 28333 (rawtx)
- **Data Volume:** bitcoin-data
- **Configuration:** Mounted from ./bitcoin.conf

### 2. Electrs (Electrum Rust Server)
- **Image:** mempool/electrs:latest
- **Ports:**
  - Electrum Protocol & REST API: 50001
- **Features:**
  - JSON-RPC Import: Uses Bitcoin Core's RPC instead of direct block files
  - Address Search: Enables prefix-based address search
  - Transaction History: Up to 1000 transactions per address
  - CORS Support: Enabled for all origins (*)
- **Data Volume:** electrs-data

## Utility Scripts

The following utility scripts are provided in the `bitcoin-scripts` directory:

### 1. check-address-balance.sh
- **Usage:** `./check-address-balance.sh <bitcoin_address>`
- **Description:** Checks the balance of any Bitcoin address using the Electrs REST API.
- **Features:**
  - Supports all Bitcoin address formats (legacy, P2SH, Bech32)
  - Shows confirmed and unconfirmed balances
  - Displays detailed UTXO information (transaction IDs, amounts, block heights)
  - Falls back to JSON-RPC if REST API fails

### 2. send-bitcoins.sh
- **Usage:** `./send-bitcoins.sh <address> <amount>`
- **Description:** Sends bitcoins from the regtest-wallet to a specified address.
- **Features:**
  - Checks wallet balance before sending
  - Shows transaction ID after successful sending
  - Displays updated wallet balance

### 3. generate-blocks.sh
- **Usage:** `./generate-blocks.sh <number_of_blocks>`
- **Description:** Generates new blocks in the regtest network.
- **Features:**
  - Default: Generates 1 block if no count is specified
  - Shows block hash for each generated block

### 4. fund-test-wallet.sh
- **Usage:** `./fund-test-wallet.sh`
- **Description:** Funds the regtest-wallet with newly mined bitcoins.
- **Features:**
  - Mines blocks to generate coins
  - Automatically sends coins to the wallet

### 5. bitcoin-transaction-test.sh
- **Usage:** `./bitcoin-transaction-test.sh`
- **Description:** Runs a complete test of Bitcoin transaction functionality.
- **Features:**
  - Creates test transactions
  - Verifies transaction processing

## Getting Started

1. **Start the Bitcoin Services:**
   ```bash
   # Start the Bitcoin Core and Electrs services
   docker-compose up -d bitcoind electrs
   ```

2. **Wait for Initialization:**
   - Allow a few moments for the services to initialize
   - Bitcoin Core needs to create the initial regtest blockchain
   - Electrs needs to connect to Bitcoin Core and index the blockchain

3. **Fund the Test Wallet:**
   ```bash
   cd bitcoin-scripts
   ./fund-test-wallet.sh
   ```

4. **Generate Some Blocks:**
   ```bash
   ./generate-blocks.sh 10  # Generate 10 new blocks
   ```

5. **Check Wallet Balance:**
   ```bash
   # Using bitcoin-cli
   docker exec bitcoind bitcoin-cli -regtest -rpcuser=bitcoin -rpcpassword=bitcoin getbalance
   
   # Or create and check a new address
   ADDRESS=$(docker exec bitcoind bitcoin-cli -regtest -rpcuser=bitcoin -rpcpassword=bitcoin -rpcwallet=regtest-wallet getnewaddress)
   ./send-bitcoins.sh $ADDRESS 1.0
   ./check-address-balance.sh $ADDRESS
   ```

## API Access

### Bitcoin RPC
- Connect to `http://localhost:8332`
- Use credentials: bitcoin:bitcoin
- Example:
  ```bash
  curl --user bitcoin:bitcoin --data-binary '{"jsonrpc":"1.0","id":"curltest","method":"getblockcount","params":[]}' -H 'content-type: text/plain;' http://localhost:8332/
  ```

### Electrs REST API
- Base URL: `http://localhost:50001`
- Endpoints:
  - `/address/{addr}` - Get address information
  - `/address/{addr}/utxo` - Get UTXOs for an address
  - Check the script `check-address-balance.sh` for usage examples

## Troubleshooting

- **Reset the Environment:**
  ```bash
  # Stop the services
  docker-compose stop bitcoind electrs
  
  # Remove the containers
  docker-compose rm -f bitcoind electrs
  
  # Remove the volumes
  docker volume rm strip-node_bitcoin-data strip-node_electrs-data
  
  # Start fresh
  docker-compose up -d bitcoind electrs
  ```

- **Check Logs:**
  ```bash
  # Bitcoin Core logs
  docker logs bitcoind
  
  # Electrs logs
  docker logs electrs
  ```
