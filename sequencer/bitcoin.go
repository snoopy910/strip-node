package sequencer

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"regexp"
)

// Bitcoin integration constants
const (
	BTC_TOKEN_SYMBOL = "BTC"                                        // Symbol representing Bitcoin
	SATOSHI_DECIMALS = 8                                            // Number of decimals in Bitcoin (1 BTC = 10^8 satoshis)
	BTC_ZERO_ADDRESS = "0x0000000000000000000000000000000000000000" // Representing a zero address for Bitcoin (not used)
)

// BlockCypherTransaction represents the transaction data structure returned by BlockCypher
type BlockCypherTransaction struct {
	Hash          string              `json:"hash"`          // Transaction hash
	Inputs        []BlockCypherInput  `json:"inputs"`        // List of inputs in the transaction
	Outputs       []BlockCypherOutput `json:"outputs"`       // List of outputs in the transaction
	Fees          int64               `json:"fees"`          // Transaction fees in satoshis
	Confirmations int                 `json:"confirmations"` // Number of confirmations
}

// BlockCypherInput represents the input of a transaction
type BlockCypherInput struct {
	Addresses   []string `json:"addresses"`    // List of addresses involved in the input
	Value       int64    `json:"value"`        // Value of the input in satoshis
	OutputValue int64    `json:"output_value"` // Alternative field for output value
}

// BlockCypherOutput represents the output of a transaction
type BlockCypherOutput struct {
	Addresses   []string `json:"addresses"`    // List of addresses involved in the output
	Value       int64    `json:"value"`        // Value of the output in satoshis
	OutputValue int64    `json:"output_value"` // Alternative field for output value
	Script      string   `json:"script"`       // Output script
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

// GetChainFunc is a function type for getting chain information
type GetChainFunc func(chainId string) (Chain, error)

// defaultGetChain is the default implementation that uses the GetChain function
var defaultGetChain = func(chainId string) (Chain, error) {
	return GetChain(chainId)
}

// FeeDetails represents the transaction fee information
type FeeDetails struct {
	FeeAmount    int64  `json:"feeAmount"`    // Fee amount in satoshis
	FormattedFee string `json:"formattedFee"` // Fee amount formatted in BTC
	TotalInputs  int64  `json:"totalInputs"`  // Total input value in satoshis
	TotalOutputs int64  `json:"totalOutputs"` // Total output value in satoshis
}

// GetBitcoinTransfers fetches Bitcoin transaction details and parses them into transfers
// This function processes a transaction and extracts transfers from inputs and outputs
func GetBitcoinTransfers(chainId string, txHash string) ([]Transfer, *FeeDetails, error) {
	// Get chain information
	chain, err := defaultGetChain(chainId)
	if err != nil {
		return nil, nil, fmt.Errorf("chain not found")
	}

	// Fetch transaction details from BlockCypher using the chain URL and txHash
	tx, err := FetchTransaction(chain.ChainUrl, txHash)
	if err != nil {
		return nil, nil, err
	}

	// Validate transaction has inputs
	if len(tx.Inputs) == 0 {
		return nil, nil, fmt.Errorf("transaction has no inputs")
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

		// Use OutputValue if Value is not available
		inputValue := input.Value
		if inputValue == 0 {
			inputValue = input.OutputValue
		}
		totalInputValue += inputValue // Sum input values

		// Process outputs of the transaction
		for _, output := range tx.Outputs {
			if len(output.Addresses) == 0 {
				continue // Skip output if there are no addresses
			}

			// Extract value (amount) from the output and convert to string
			outputValue := output.Value
			if outputValue == 0 {
				outputValue = output.OutputValue
			}
			scaledAmount := fmt.Sprintf("%d", outputValue)

			// Format the amount using the helper function
			formattedAmount, err := getFormattedAmount(scaledAmount, SATOSHI_DECIMALS)
			if err != nil {
				return nil, nil, fmt.Errorf("error formatting amount: %w", err)
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

	// Validate we found some transfers
	if len(transfers) == 0 {
		return nil, nil, fmt.Errorf("no transfers found")
	}

	// Process outputs of the transaction to calculate the total output value
	for _, output := range tx.Outputs {
		if len(output.Addresses) == 0 {
			continue // Skip output if there are no addresses
		}
		outputValue := output.Value
		if outputValue == 0 {
			outputValue = output.OutputValue
		}
		totalOutputValue += outputValue // Sum output values
	}

	// Get transaction fee from BlockCypher API response
	transactionFee := tx.Fees
	formattedFee, err := getFormattedAmount(fmt.Sprintf("%d", transactionFee), SATOSHI_DECIMALS)
	if err != nil {
		return nil, nil, fmt.Errorf("error formatting fee: %w", err)
	}

	// Create fee details structure
	feeDetails := &FeeDetails{
		FeeAmount:    transactionFee,
		FormattedFee: formattedFee,
		TotalInputs:  totalInputValue,
		TotalOutputs: totalOutputValue,
	}

	// Return the list of transfers and fee details
	return transfers, feeDetails, nil
}

// FetchUTXOValue fetches the value of a UTXO using BlockCypher API
// This function retrieves the value of an unspent transaction output (UTXO)
func FetchUTXOValue(chainUrl string, txHash string) (int64, error) {
	// Fetch the full transaction details
	tx, err := FetchTransaction(chainUrl, txHash)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch transaction: %v", err)
	}

	// For UTXOs, we're interested in the first output value
	// In most Bitcoin transactions, the first output is the actual transfer amount
	if len(tx.Outputs) == 0 {
		return 0, fmt.Errorf("transaction has no outputs")
	}

	// Use OutputValue if Value is not available
	outputValue := tx.Outputs[0].Value
	if outputValue == 0 {
		outputValue = tx.Outputs[0].OutputValue
	}

	return outputValue, nil
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

// isValidBitcoinAddress checks if a Bitcoin address is valid
func isValidBitcoinAddress(address string) bool {
	// Basic regex to validate Bitcoin address format
	re := regexp.MustCompile(`^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$`)
	return re.MatchString(address)
}
