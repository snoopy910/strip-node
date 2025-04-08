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

## API Flow

### 1. Construct Phase
- Creates transaction parameters with 5-minute validity
- Returns hash for signing
- Stores parameters securely for later retrieval

### 2. Solve Phase
- Retrieves stored parameters using operation hash
- Verifies signature and recovers sender
- Executes transaction with proper gas settings

### 3. Status & Output Phase
- Monitors transaction status
- Parses relevant events (IncreaseLiquidity/DecreaseLiquidity)
- Returns structured output with position details

## Example Intent

```json
{
  "operations": [{
    "solverMetadata": {
      "action": "mint",
      "token0": "0x...",
      "token1": "0x...",
      "fee": 3000,
      "tickLower": -180,
      "tickUpper": 180,
      "amount0Desired": "1000000000000000000",
      "amount1Desired": "1000000000000000000",
      "amount0Min": "990000000000000000",
      "amount1Min": "990000000000000000"
    }
  }],
  "identity": "0x..." // caller's address
}
```

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
