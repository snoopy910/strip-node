package sequencer

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
)

// Bitcoin integration constants
const (
	BTC_TOKEN_SYMBOL = "BTC"                                        // Symbol representing Bitcoin
	SATOSHI_DECIMALS = 8                                            // Number of decimals in Bitcoin (1 BTC = 10^8 satoshis)
	BTC_ZERO_ADDRESS = "0x0000000000000000000000000000000000000000" // Representing a zero address for Bitcoin (not used)
)

// BlockCypherTransaction represents the transaction data structure returned by BlockCypher
type BlockCypherTransaction struct {
	Inputs  []BlockCypherInput  `json:"inputs"`  // List of inputs in the transaction
	Outputs []BlockCypherOutput `json:"outputs"` // List of outputs in the transaction
}

// BlockCypherInput represents the input of a transaction
type BlockCypherInput struct {
	Addresses []string `json:"addresses"` // List of addresses involved in the input
}

// BlockCypherOutput represents the output of a transaction
type BlockCypherOutput struct {
	Addresses []string `json:"addresses"` // List of addresses involved in the output
	Value     int64    `json:"value"`     // Value of the output in satoshis
}

// FetchTransaction fetches transaction details from BlockCypher
// This function retrieves a Bitcoin transaction by its hash
func FetchTransaction(chainUrl string, txHash string) (*BlockCypherTransaction, error) {
	// Construct the URL for the API request using the chain URL and transaction hash
	url := fmt.Sprintf("%s/txs/%s", chainUrl, txHash)

	// Send GET request to BlockCypher API
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK (200)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the JSON response into a BlockCypherTransaction struct
	var tx BlockCypherTransaction
	if err := json.NewDecoder(resp.Body).Decode(&tx); err != nil {
		return nil, fmt.Errorf("failed to decode transaction response: %v", err)
	}

	return &tx, nil
}

// GetBitcoinTransfers fetches Bitcoin transaction details and parses them into transfers
// This function processes a transaction and extracts transfers from inputs and outputs
func GetBitcoinTransfers(chainId string, txHash string) ([]Transfer, error) {
	// Get chain configuration using the chain ID
	chain, err := GetChain(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain config: %v", err)
	}

	// Fetch transaction details from BlockCypher using the chain URL and txHash
	tx, err := FetchTransaction(chain.ChainUrl, txHash)
	if err != nil {
		return nil, err
	}

	// Initialize slice to hold transfers and variables to calculate total input/output values
	var transfers []Transfer
	var totalInputValue int64
	var totalOutputValue int64

	// Process inputs of the transaction
	for _, input := range tx.Inputs {
		if len(input.Addresses) == 0 {
			continue // Skip input if there are no addresses
		}
		fromAddress := input.Addresses[0] // Get the first address from the input

		// Process outputs of the transaction
		for _, output := range tx.Outputs {
			if len(output.Addresses) == 0 {
				continue // Skip output if there are no addresses
			}

			// Extract value (amount) from the output and convert to string
			amount := output.Value // Amount is in satoshis (1 satoshi = 0.00000001 BTC)
			scaledAmount := fmt.Sprintf("%d", amount)

			// Format the amount using the helper function
			formattedAmount, err := getFormattedAmount(scaledAmount, SATOSHI_DECIMALS)
			if err != nil {
				return nil, fmt.Errorf("error formatting amount: %w", err)
			}

			// Append the transfer details to the transfers slice
			transfers = append(transfers, Transfer{
				From:         fromAddress,         // From address of the transfer
				To:           output.Addresses[0], // To address of the transfer
				Amount:       formattedAmount,     // Formatted transfer amount in BTC
				Token:        BTC_TOKEN_SYMBOL,    // Token symbol (BTC)
				IsNative:     true,                // Flag indicating it's a native BTC transfer
				TokenAddress: BTC_ZERO_ADDRESS,    // Token address (zero address in this case)
				ScaledAmount: scaledAmount,        // Transfer amount in satoshis
			})
		}
	}

	// Process outputs of the transaction to calculate the total output value
	for _, output := range tx.Outputs {
		if len(output.Addresses) == 0 {
			continue // Skip output if there are no addresses
		}
		totalOutputValue += output.Value // Sum output values
	}

	// Calculate transaction fee by subtracting total output value from total input value
	transactionFee := totalInputValue - totalOutputValue
	_, err = getFormattedAmount(fmt.Sprintf("%d", transactionFee), SATOSHI_DECIMALS)
	if err != nil {
		return nil, fmt.Errorf("error formatting fee: %w", err)
	}

	// TODO: Output the fee details as needed

	// Return the list of transfers
	return transfers, nil
}

// FetchUTXOValue fetches the value of a UTXO (mock function)
// This function simulates fetching the value of an unspent transaction output (UTXO)
func FetchUTXOValue(chainUrl string, txHash string) (int64, error) {
	// Example: Use a dummy value for now, can be replaced with actual logic to fetch UTXO details
	return 100000, nil // Example value (100,000 satoshis)
}

// Helper function to format amounts
// This function converts a string representation of an amount into a properly formatted BTC value
func getFormattedAmount(amount string, decimal int) (string, error) {
	// Convert the amount string into a big integer
	bigIntAmount := new(big.Int)
	_, success := bigIntAmount.SetString(amount, 10)
	if !success {
		return "", fmt.Errorf("error: Invalid number string")
	}

	// Format the amount into a human-readable BTC value with the specified decimal precision
	formattedAmount, err := FormatUnits(bigIntAmount, decimal)
	if err != nil {
		return "", err
	}

	return formattedAmount, nil
}
