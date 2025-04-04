# Uniswap V3 Solver

This solver handles Uniswap V3 liquidity operations through the NonfungiblePositionManager (NPM) contract. It supports minting new positions and exiting existing positions.

## Features

- Mint new liquidity positions with specified parameters
- Exit existing liquidity positions
- Transaction status tracking
- Asynchronous operation handling

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
