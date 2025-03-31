# Uniswap V3 Solver for StripChain

## Overview

This solver enables StripChain users to interact with Uniswap V3 through the intent-based StripVM architecture. It handles liquidity operations such as minting new positions and exiting existing positions using the NonfungiblePositionManager (NPM) contract.

## Features

- Mint new liquidity positions
- Exit existing positions
- Real-time transaction status monitoring
- Integration with StripChain's intent system

## Structure

The solver is organized into the following files:

- `types.go`: Defines custom types and interfaces for the solver
- `solver.go`: Implements the core logic for Uniswap V3 operations
- `api.go`: Handles HTTP requests and responses for the solver

## API Endpoints

The solver implements the standard StripVM solver API:

### POST /construct
Constructs transaction data for a Uniswap V3 operation.

### POST /solve
Executes a Uniswap V3 operation with the given signature.

### POST /status
Checks the status of a Uniswap V3 operation.

### POST /output
Retrieves the result of a Uniswap V3 operation.

## Flow

### 1. Construct Phase
```http
POST /construct?operationIndex=<operation_index>
```

The construct phase creates a complete transaction object and returns a hash that needs to be signed. The transaction includes:
- Chain ID
- Nonce
- Gas parameters (tip cap, fee cap, gas limit)
- Contract address
- Function data based on the operation type (mint or exit)

### 2. Solve Phase
```http
POST /solve?operationIndex=<operation_index>
```

The solve phase takes the signed transaction and executes it on the blockchain. The solver:
1. Verifies the signature
2. Submits the transaction
3. Returns the transaction hash

### 3. Status Phase
```http
GET /status?operationIndex=<operation_index>&txHash=<tx_hash>
```

The status phase checks the current state of the transaction:
- `pending`: Transaction is still being processed
- `success`: Transaction was successfully mined
- `failed`: Transaction failed during execution

### 4. Output Phase
```http
GET /output?operationIndex=<operation_index>&txHash=<tx_hash>
```

The output phase returns the result of the transaction after it has been mined.

## Example Intent

Here's an example of a mint operation intent:

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
