# Lending Solver for StripChain

## Overview

This solver enables StripChain users to interact with the LendingPool contract through the intent-based StripVM architecture. It handles lending operations such as supplying collateral, borrowing StripUSD, repaying debt, and withdrawing collateral.

## Features

- Supply tokens as collateral
- Borrow StripUSD against supplied collateral
- Repay StripUSD debt
- Withdraw collateral
- Real-time health factor monitoring
- Integration with StripChain's intent system

## Structure

The solver is organized into the following files:

- `types.go`: Defines custom types and interfaces for the solver
- `solver.go`: Implements the core logic for lending operations
- `handler.go`: Handles HTTP requests and responses for the solver

## API Endpoints

The solver implements the standard StripVM solver API:

### POST /construct
Constructs transaction data for a lending operation.

### POST /solve
Executes a lending operation with the given signature.

### POST /status
Checks the status of a lending operation.

### POST /output
Retrieves the result of a lending operation.

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
- Operation data

Response:
```json
{
    "dataToSign": "0x7a69f0..." // Hash to be signed by the user
}
```

### 2. Solve Phase
```http
POST /solve?operationIndex=<operation_index>&signature=<signature>
```

The solve phase takes the user's signature, recreates the transaction, and executes it on the blockchain. It:
1. Recreates the exact transaction from the construct phase
2. Adds the user's signature
3. Broadcasts the signed transaction

Response:
```json
{
    "result": "0x8fd92c..." // Transaction hash
}
```

### 3. Status Phase
```http
POST /status?operationIndex=<operation_index>
```

Monitors the transaction progress and returns one of three statuses:
- pending: Transaction is still being processed
- success: Transaction was successfully executed
- failure: Transaction failed or was reverted

Response:
```json
{
    "status": "pending|success|failure"
}
```

### 4. Output Phase
```http
POST /output?operationIndex=<operation_index>
```

Once the operation succeeds, returns the operation results:
```json
{
    "txHash": "0x8fd92c...",
    "amount": "1000000000000000000",
    "collateralUSD": "2000000000000000000",
    "borrowedUSD": "1000000000000000000",
    "healthFactor": "2000000000000000000"
}
```

## Operation Types

### Supply
Supplies collateral tokens to the lending pool.
```json
{
    "action": "supply",
    "token": "0x...",
    "amount": {
        "int": "1000000000000000000"
    }
}
```

### Borrow
Borrows StripUSD against supplied collateral.
```json
{
    "action": "borrow",
    "amount": {
        "int": "1000000000000000000"
    }
}
```

### Repay
Repays borrowed StripUSD.
```json
{
    "action": "repay",
    "amount": {
        "int": "1000000000000000000"
    }
}
```

### Withdraw
Withdraws collateral tokens.
```json
{
    "action": "withdraw",
    "token": "0x...",
    "amount": {
        "int": "1000000000000000000"
    }
}
```

## Error Handling

The solver implements comprehensive error handling:
- Invalid operation parameters
- Transaction construction failures
- Network errors
- Transaction reversion
- Invalid signatures

All errors are returned with appropriate error messages and HTTP status codes.

## Implementation Details

The solver uses:
- go-ethereum for blockchain interactions
- EIP-1559 transaction type for better gas handling
- LondonSigner for transaction signing
- In-memory transaction status tracking
- Proper nonce management
- Dynamic gas price estimation

## Example Usage

```json
{
  "intent": {
    "operations": [
      {
        "id": 1,
        "type": "lending",
        "chainId": "1",
        "solverMetadata": {
          "action": "supply",
          "token": "0x...",
          "amount": {
            "int": "1000000000000000000"
          },
          "isCollateral": true
        }
      }
    ]
  },
  "opIndex": 0
}
```

## Development

To run the solver locally:

1. Set up your Go environment
2. Install dependencies
3. Configure the RPC URL and contract addresses
4. Run the server: `go run main.go`
