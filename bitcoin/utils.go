package bitcoin

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"regexp"

	"github.com/StripChain/strip-node/common"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
)

type BitcoinNetworkConfig struct {
	Name          string
	NetworkType   string
	AddressPrefix []string
	RPCPort       int
	Params        *chaincfg.Params
}

var (
	MainnetConfig = BitcoinNetworkConfig{
		Name:          "mainnet",
		NetworkType:   "mainnet",
		AddressPrefix: []string{"1", "3", "bc1"},
		RPCPort:       8332,
		Params:        &chaincfg.MainNetParams,
	}

	TestnetConfig = BitcoinNetworkConfig{
		Name:          "testnet",
		NetworkType:   "testnet3",
		AddressPrefix: []string{"m", "n", "tb1"},
		RPCPort:       18332,
		Params:        &chaincfg.TestNet3Params,
	}

	RegtestConfig = BitcoinNetworkConfig{
		Name:          "regtest",
		NetworkType:   "regtest",
		AddressPrefix: []string{"m", "n", "bcrt1"},
		RPCPort:       18443,
		Params:        &chaincfg.RegressionNetParams,
	}
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

// GetChainFunc is a function type for getting chain information
type GetChainFunc func(chainId string) (common.Chain, error)

// defaultGetChain is the default implementation that uses the GetChain function
var defaultGetChain = func(chainId string) (common.Chain, error) {
	return common.GetChain(chainId)
}

// FeeDetails represents the transaction fee information
type FeeDetails struct {
	FeeAmount    int64  `json:"feeAmount"`    // Fee amount in satoshis
	FormattedFee string `json:"formattedFee"` // Fee amount formatted in BTC
	TotalInputs  int64  `json:"totalInputs"`  // Total input value in satoshis
	TotalOutputs int64  `json:"totalOutputs"` // Total output value in satoshis
}

// fetchTransaction fetches transaction details from BlockCypher
// This function retrieves a Bitcoin transaction by its hash
func fetchTransaction(chainUrl string, txHash string) (*BlockCypherTransaction, error) {
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

// fetchUTXOValue fetches the value of a UTXO using BlockCypher API
// This function retrieves the value of an unspent transaction output (UTXO)
func fetchUTXOValue(chainUrl string, txHash string) (int64, error) {
	// Fetch the full transaction details
	tx, err := fetchTransaction(chainUrl, txHash)
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
	formattedAmount, err := formatUnits(bigIntAmount, decimal)
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

// parseSerializedTransaction parses a base64 encoded serialized transaction (PSBT or raw transaction)
// and returns the unsigned transaction as a wire.MsgTx.
// If the input is a PSBT, it extracts the unsigned transaction from the global map.
// If the input is a raw transaction, it deserializes it directly.
func parseSerializedTransaction(serializedTxn string) (*wire.MsgTx, error) {
	// Decode base64 transaction
	txBytes, err := base64.StdEncoding.DecodeString(serializedTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 transaction: %v", err)
	}

	// Check if this is a PSBT by looking for magic bytes
	if bytes.HasPrefix(txBytes, []byte("psbt\xff")) {
		return parsePSBT(txBytes)
	}

	// If not a PSBT, try to parse as a raw transaction
	return parseRawTransaction(txBytes)
}

// parsePSBT parses a Partially Signed Bitcoin Transaction and returns the unsigned transaction.
func parsePSBT(txBytes []byte) (*wire.MsgTx, error) {
	// Skip PSBT magic bytes and separator
	txBytes = txBytes[len("psbt\xff"):]

	// Extract the unsigned transaction from the global map
	var unsignedTxBytes []byte
	var offset int
	for offset < len(txBytes) {
		// Read key length
		if offset >= len(txBytes) {
			break
		}
		keyLen := int(txBytes[offset])
		offset++

		if keyLen == 0 {
			// Separator to inputs map
			offset++
			break
		}

		// Read key
		if offset+keyLen >= len(txBytes) {
			return nil, fmt.Errorf("invalid key length")
		}
		key := txBytes[offset : offset+keyLen]
		offset += keyLen

		// Read value length
		if offset >= len(txBytes) {
			return nil, fmt.Errorf("no value length")
		}
		valueLen := int(txBytes[offset])
		offset++

		if offset+valueLen > len(txBytes) {
			return nil, fmt.Errorf("invalid value length")
		}
		value := txBytes[offset : offset+valueLen]
		offset += valueLen

		// If this is the unsigned transaction (key = 0x00)
		if len(key) == 1 && key[0] == 0x00 {
			unsignedTxBytes = value
		}
	}

	if unsignedTxBytes == nil {
		return nil, fmt.Errorf("no unsigned transaction found in PSBT")
	}

	return parseRawTransaction(unsignedTxBytes)
}

// parseRawTransaction parses a raw transaction bytes into a wire.MsgTx.
func parseRawTransaction(txBytes []byte) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(wire.TxVersion)
	err := tx.Deserialize(bytes.NewReader(txBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %v", err)
	}
	return tx, nil
}

func formatUnits(value *big.Int, decimals int) (string, error) {
	// Create the scaling factor as 10^decimals
	scalingFactor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))

	// Convert the value to a big.Float
	valueFloat := new(big.Float).SetInt(value)

	// Divide the value by the scaling factor
	result := new(big.Float).Quo(valueFloat, scalingFactor)

	// Convert the result to a string with the appropriate precision
	return result.Text('f', decimals), nil
}
