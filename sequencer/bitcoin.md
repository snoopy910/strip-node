# Bitcoin Module Documentation

## Overview

This module adds Bitcoin (BTC) support to the `strip-node` sequencer. It enables handling native BTC transfers and UTXO-based transactions, including parsing, transfer event tracking, and transaction fee calculation.

## Features
1. Native BTC transfers with support for both mainnet and testnet
2. UTXO-based transaction parsing and management
3. Bitcoin amount formatting (satoshis to BTC)
4. Transaction fee calculation
5. Integration with existing sequencer logic
6. Comprehensive error handling for all operations

## File Structure

### `bitcoin.go`
- Contains the implementation of Bitcoin module
- Handles Bitcoin transfers using the `Transfer` struct
- Provides utilities for formatting BTC amounts and calculating transaction fees
- Implements BlockCypher API integration for transaction data

### `bitcoin_test.go`
- Unit tests for Bitcoin functionality, including:
  - Mainnet and testnet transaction validation
  - Error handling scenarios
  - UTXO value fetching
  - Transaction fee calculation

## Implementation Details

### Chain Configuration

Bitcoin chain configuration is managed through `chain.go`, which defines the following chains:
- Mainnet (ChainId: "1000"): Uses "https://api.blockcypher.com/v1/btc/main"
- Testnet (ChainId: "1001"): Uses "https://api.blockcypher.com/v1/btc/test3"
- Regtest (ChainId: "1002"): Uses "http://localhost:18443/v1/btc/regtest"

### Transaction Processing
1. **Chain Selection**
   - Mainnet (chainId: "1000"): Uses "https://api.blockcypher.com/v1/btc/main"
   - Testnet (chainId: "1001"): Uses "https://api.blockcypher.com/v1/btc/test3"
   - Regtest (chainId: "1002"): Uses "http://localhost:18443/v1/btc/regtest" (local development network)

2. **UTXO Management**
   ```go
   type BlockCypherTransaction struct {
       Inputs  []BlockCypherInput  `json:"inputs"`
       Outputs []BlockCypherOutput `json:"outputs"`
       Fees    int64              `json:"fees"`
   }
   ```

3. **Fee Calculation**
   ```go
   feeDetails := &FeeDetails{
       FeeAmount:    tx.Fees,
       TotalInputs:  totalInputValue,
       TotalOutputs: totalOutputValue,
   }
   ```

4. **Amount Formatting**
   - Converts satoshis to BTC (8 decimal places)
   - Example: 100000000 satoshis → "1.00000000" BTC

### Error Handling

The module handles various error scenarios:
1. Invalid chain ID
2. Invalid transaction hash
3. API server errors
4. Network timeouts
5. Empty responses
6. Missing input/output addresses
7. Malformed JSON responses
8. Rate limit exceeded errors
9. Invalid address format
10. Insufficient confirmations

#### Rate Limiting
- BlockCypher API has a rate limit of 200 requests per hour for free tier
- Implements exponential backoff for rate limit errors
- Provides clear error messages when rate limits are exceeded

### API Integration

#### BlockCypher API
- **Base URLs**:
  - Mainnet: https://api.blockcypher.com/v1/btc/main
  - Testnet: https://api.blockcypher.com/v1/btc/test3
- **Endpoints**:
  - Transaction details: `/txs/{hash}`
  - UTXO details: `/txs/{hash}?includeHex=true`

## Testing

### Test Coverage
1. **Basic Transfer Tests**
   - Mainnet transaction validation
   - Fee calculation verification

2. **Error Handling Tests**
   - Invalid chain ID
   - Invalid transaction hash
   - API server errors
   - Malformed JSON responses
   - Network timeouts
   - Empty responses
   - Missing output addresses

3. **UTXO Value Tests**
   - Valid Bitcoin transaction
   - Invalid transaction hash
   - Invalid chain URL

### Test Results
All test cases pass successfully, including:
- Transaction fetching and parsing
- Transfer formatting
- Fee calculation
- Error handling scenarios

## Configuration

### Constants
```go
const (
    BTC_TOKEN_SYMBOL = "BTC"
    SATOSHI_DECIMALS = 8
    BTC_ZERO_ADDRESS = "0x0000000000000000000000000000000000000000"
)
```

## Future Enhancements
1. Add support for additional Bitcoin APIs (e.g., Blockstream API) for redundancy
2. Implement mempool monitoring for unconfirmed transactions
3. Add support for SegWit addresses
4. Enhance UTXO selection logic for better fee optimization
5. Add support for multi-signature transactions

## References
- [BlockCypher API Documentation](https://www.blockcypher.com/dev/bitcoin/)
- [Bitcoin Transaction Format](https://en.bitcoin.it/wiki/Transaction)
- [UTXO Model](https://en.wikipedia.org/wiki/Unspent_transaction_output)

Last Updated: 2025-01-12
