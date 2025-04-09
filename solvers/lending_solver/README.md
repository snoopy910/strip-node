# Lending Solver for StripChain

## Overview

This solver enables StripChain users to interact with the LendingPool contract through the intent-based StripVM architecture. It handles lending operations such as supplying collateral, borrowing assets, repaying debt, and withdrawing collateral.

## Features

- Supply tokens as collateral
- Borrow assets against supplied collateral
- Repay borrowed debt
- Withdraw collateral
- Real-time transaction status monitoring
- Secure transaction parameter management with auto-expiry
- Dynamic gas estimation and fee calculation
- Transaction signature verification with sender recovery

## Implementation Details

### Transaction Parameter Management
- Parameters are stored with a 5-minute expiration window
- Automatic cleanup of expired parameters via background goroutine
- Thread-safe parameter storage using mutex protection

### Gas Management
- Dynamic gas estimation based on operation type
- Uses network-suggested gas prices with dynamic fee caps
- Proper error handling for failed transactions

### Security Features
- Signature verification with sender address recovery
- Expiring transaction parameters to prevent replay attacks
- Proper error handling and validation at each step

## API Endpoints

### GET /health
Health check endpoint that returns 200 OK.

### POST /construct
Constructs transaction data for a lending operation.

Request body:
```json
{
  "operations": [{
    "solverMetadata": {
      "action": "supply",
      "token": "0x...",
      "amount": {
        "int": "1000000000000000000"
      },
      "isCollateral": true
    }
  }],
  "identity": "0x..." // caller's address
}
```

### POST /solve
Executes a signed lending operation.

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
  "amount": {
    "int": "1000000000000000000"
  },
  "token": {
    "address": "0x..."
  }
}
```

### Error Output
```json
{
  "error": "detailed error message"
}
```

## Operation Types

### Supply
- Supply tokens as collateral
- Optionally enable/disable collateral usage

### Borrow
- Borrow assets against supplied collateral
- Requires sufficient collateral value

### Repay
- Repay borrowed debt
- Full or partial repayment supported

### Withdraw
- Withdraw supplied collateral
- Subject to health factor requirements
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
