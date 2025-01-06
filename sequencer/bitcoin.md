
# Bitcoin Module Documentation

## Overview

This module adds Bitcoin (BTC) support to the `strip-node` sequencer. It enables handling native BTC transfers and UTXO-based transactions, including parsing, transfer event tracking, and transaction fee calculation.

## Features
1. Native BTC transfers
2. UTXO-based transaction parsing and management
3. Bitcoin amount formatting (satoshis to BTC)
4. Transaction fee calculation
5. Integration with existing sequencer logic

## File Structure

### `bitcoin.go`
- Contains the implementation of Bitcoin module.
- Handles Bitcoin transfers using the `Transfer` struct.
- Provides utilities for formatting BTC amounts and calculating transaction fees.

### `bitcoin_test.go`
- Unit tests for Bitcoin functionality, including:
  - Validating transfers
  - Testing mainnet and testnet configurations

## Key Functions

### `FetchTransaction`
Fetches Bitcoin transaction details from BlockCypher API.

**Parameters:**
- `chainUrl`: The URL of the BlockCypher API.
- `txHash`: The transaction hash.

**Returns:**
- A `BlockCypherTransaction` struct containing transaction details.

---

### `GetBitcoinTransfers`
Fetches and parses Bitcoin transaction details into a list of `Transfer` objects.

**Parameters:**
- `chainId`: The ID of the blockchain.
- `txHash`: The transaction hash.

**Returns:**
- A list of `Transfer` objects representing parsed transactions.

---

### `getFormattedAmount`
Formats BTC amounts by converting satoshis to BTC.

**Parameters:**
- `amount`: The amount in satoshis as a string.
- `decimal`: The number of decimals for BTC (8).

**Returns:**
- A string representation of the formatted BTC amount.

---

### `FetchUTXOValue` (Mock Function)
Fetches the value of a UTXO. Placeholder for future implementation.

**Parameters:**
- `chainUrl`: The URL of the BlockCypher API.
- `txHash`: The transaction hash.

**Returns:**
- The value of the UTXO in satoshis.

## Testing
The `bitcoin_test.go` file provides test cases for validating:
1. Basic BTC transfer parsing.
2. Mainnet and testnet configurations.
3. Amount formatting and fee calculations.

## Configuration
- **Chain Configuration**: The module integrates seamlessly with the sequencer's existing chain configuration system.
- **Libraries Used**: `btcd/btcutil` for handling BTC-specific operations.

## Notes
- The module currently uses BlockCypher API for transaction details.
- Proper error handling and logging are implemented for robustness.

## Future Enhancements
1. Add support for additional Bitcoin APIs to improve redundancy.
2. Optimize UTXO selection logic for better fee management.
3. Enhance documentation and include more examples in test cases.
