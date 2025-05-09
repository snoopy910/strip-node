# Uniswap V3 Solver for StripChain

## Overview

This solver enables StripChain users to interact with Uniswap V3 through the intent-based StripVM architecture. It handles liquidity operations such as minting new positions and exiting existing positions using the NonfungiblePositionManager (NPM) contract.

## Features

- Mint new liquidity positions
- Exit existing positions
- Real-time transaction status monitoring with event parsing
- Secure transaction parameter management with auto-expiry
- Dynamic gas estimation and fee calculation
- Transaction signature verification with sender recovery

## Implementation Details

### Transaction Parameter Management
- Parameters are stored with a 5-minute expiration window
- Automatic cleanup of expired parameters via background goroutine
- Thread-safe parameter storage using mutex protection

### Event Handling
- Support for `IncreaseLiquidity` and `DecreaseLiquidity` events
- Proper parsing of indexed and non-indexed event parameters
- Detailed error reporting for event parsing failures

### Gas Management
- Dynamic gas estimation based on operation type
- 20% buffer added to estimated gas for safety
- Uses network-suggested gas prices with dynamic fee caps

### Security Features
- Signature verification with sender address recovery
- Expiring transaction parameters to prevent replay attacks
- Proper error handling and validation at each step

## API Endpoints

### POST /construct
Constructs transaction data for a Uniswap V3 operation.

Request body:
```json
{
  "operations": [{
    "solverMetadata": {
      "action": "mint",
      "token0": "0x...",
      "token1": "0x...",
      "fee": 3000,
      "tickLower": -100,
      "tickUpper": 100,
      "amount0Desired": "1000000000000000000",
      "amount1Desired": "1000000000000000000",
      "amount0Min": "990000000000000000",
      "amount1Min": "990000000000000000",
      "recipient": "0x...",
      "deadline": "1234567890"
    }
  }],
  "identity": "0x..." // caller's address
}
```

### POST /solve
Executes a signed Uniswap V3 operation.

Request body:
```json
{
  "operations": [...],
  "identity": "0x...",
  "opIndex": 0,
  "signature": "0x..."
}
```

### GET /status
Checks the status of a transaction.

Query parameters:
- `opIndex`: Operation index
- `txHash`: Transaction hash

### GET /output
Retrieves the result of a completed transaction.

Query parameters:
- `opIndex`: Operation index
- `txHash`: Transaction hash

## Response Format

### Success Output
```json
{
  "txHash": "0x...",
  "tokenId": 123,
  "liquidity": "1000000000000000000",
  "amountA": "1000000000000000000",
  "amountB": "1000000000000000000"
}
```

## Usage

Start the solver with:

```bash
./strip-node --isUniswapSolver \
  --rpcURL=<ethereum_rpc_url> \
  --httpPort=<port> \
  --uniswapV3FactoryAddress=<factory_address> \
  --npmAddress=<npm_address> \
  --chainId=<chain_id>
```

Or using environment variables:

```bash
IS_UNISWAP_SOLVER=true \
RPC_URL=<ethereum_rpc_url> \
HTTP_PORT=<port> \
UNISWAP_V3_FACTORY_ADDRESS=<factory_address> \
NPM_ADDRESS=<npm_address> \
CHAIN_ID=<chain_id> \
./strip-node
```
