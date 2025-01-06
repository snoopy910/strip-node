package sequencer

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
)

// Bitcoin integration constants
const (
	BTC_TOKEN_SYMBOL = "BTC"
	SATOSHI_DECIMALS = 8
	BTC_ZERO_ADDRESS = "0x0000000000000000000000000000000000000000"
)

// BlockCypherTransaction represents the transaction data structure returned by BlockCypher
type BlockCypherTransaction struct {
	Inputs  []BlockCypherInput  `json:"inputs"`
	Outputs []BlockCypherOutput `json:"outputs"`
}

type BlockCypherInput struct {
	Addresses []string `json:"addresses"`
}

type BlockCypherOutput struct {
	Addresses []string `json:"addresses"`
	Value     int64    `json:"value"` // Value in satoshis
}

// FetchTransaction fetches transaction details from BlockCypher
func FetchTransaction(chainUrl string, txHash string) (*BlockCypherTransaction, error) {
	url := fmt.Sprintf("%s/txs/%s", chainUrl, txHash)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var tx BlockCypherTransaction
	if err := json.NewDecoder(resp.Body).Decode(&tx); err != nil {
		return nil, fmt.Errorf("failed to decode transaction response: %v", err)
	}

	return &tx, nil
}

// GetBitcoinTransfers fetches Bitcoin transaction details and parses them into transfers
func GetBitcoinTransfers(chainId string, txHash string) ([]Transfer, error) {
	// Get chain configuration
	chain, err := GetChain(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain config: %v", err)
	}

	// Fetch transaction details from BlockCypher
	tx, err := FetchTransaction(chain.ChainUrl, txHash)
	if err != nil {
		return nil, err
	}

	var transfers []Transfer
	var totalInputValue int64
	var totalOutputValue int64

	// Process inputs and outputs
	for _, input := range tx.Inputs {
		if len(input.Addresses) == 0 {
			continue
		}
		fromAddress := input.Addresses[0]

		for _, output := range tx.Outputs {
			if len(output.Addresses) == 0 {
				continue
			}

			amount := output.Value // Value is already in satoshis

			scaledAmount := fmt.Sprintf("%d", amount)

			formattedAmount, err := getFormattedAmount(scaledAmount, SATOSHI_DECIMALS)
			if err != nil {
				return nil, fmt.Errorf("error formatting amount: %w", err)
			}

			transfers = append(transfers, Transfer{
				From:         fromAddress,
				To:           output.Addresses[0],
				Amount:       formattedAmount,
				Token:        BTC_TOKEN_SYMBOL,
				IsNative:     true,
				TokenAddress: BTC_ZERO_ADDRESS,
				ScaledAmount: scaledAmount,
			})
		}
	}

	// Process outputs and calculate total output value
	for _, output := range tx.Outputs {
		if len(output.Addresses) == 0 {
			continue
		}
		totalOutputValue += output.Value
	}

	// Calculate transaction fee
	transactionFee := totalInputValue - totalOutputValue
	_, err = getFormattedAmount(fmt.Sprintf("%d", transactionFee), SATOSHI_DECIMALS)
	if err != nil {
		return nil, fmt.Errorf("error formatting fee: %w", err)
	}
	// TODO: output fee

	return transfers, nil
}

// FetchUTXOValue fetches the value of a UTXO (mock function)
func FetchUTXOValue(chainUrl string, txHash string) (int64, error) {
	// Use BlockCypher or other API to fetch UTXO details by txHash
	// Example: Use a dummy value for now
	return 100000, nil // Replace with actual UTXO fetch logic
}

// Helper function to format amounts
func getFormattedAmount(amount string, decimal int) (string, error) {
	bigIntAmount := new(big.Int)

	_, success := bigIntAmount.SetString(amount, 10)
	if !success {
		return "", fmt.Errorf("error: Invalid number string")
	}

	formattedAmount, err := FormatUnits(bigIntAmount, decimal)
	if err != nil {
		return "", err
	}

	return formattedAmount, nil
}
