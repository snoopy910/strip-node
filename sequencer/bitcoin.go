package sequencer

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
)

// Transfer defines the structure of a transfer event
type Transfer struct {
	From         string
	To           string
	Amount       string
	Token        string
	IsNative     bool
	TokenAddress string
	ScaledAmount string
}

// Bitcoin integration constants
const (
	BTC_TOKEN_SYMBOL = "BTC"
	SATOSHI_DECIMALS = 8
	BTC_ZERO_ADDRESS = "0x0000000000000000000000000000000000000000"
)

// GetBitcoinTransfers fetches Bitcoin transaction details and parses them into transfers
func GetBitcoinTransfers(chainId string, txHash string) ([]Transfer, error) {
	// Get chain configuration
	chain, err := GetChain(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain config: %v", err)
	}

	// Fetch transaction details using Bitcoin Core RPC
	tx, err := getBitcoinTransaction(chain.ChainUrl, chain.RpcUsername, chain.RpcPassword, txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Bitcoin transaction: %v", err)
	}

	var transfers []Transfer

	for _, output := range tx.Vout {
		// Ignore outputs with no addresses
		if len(output.ScriptPubKey.Addresses) == 0 {
			continue
		}

		amount := big.NewFloat(output.Value)
		amount.Mul(amount, big.NewFloat(1e8)) // Convert BTC to Satoshis
		scaledAmount := fmt.Sprintf("%d", int64(amount.Float64()))

		formattedAmount, err := getFormattedAmount(scaledAmount, SATOSHI_DECIMALS)
		if err != nil {
			return nil, fmt.Errorf("error formatting amount: %w", err)
		}

		transfers = append(transfers, Transfer{
			From:         tx.Vin[0].PrevOut.Addresses[0],
			To:           output.ScriptPubKey.Addresses[0],
			Amount:       formattedAmount,
			Token:        BTC_TOKEN_SYMBOL,
			IsNative:     true,
			TokenAddress: BTC_ZERO_ADDRESS,
			ScaledAmount: scaledAmount,
		})
	}

	return transfers, nil
}

type BitcoinTransaction struct {
	Txid string `json:"txid"`
	Vin  []struct {
		PrevOut struct {
			Addresses []string `json:"addresses"`
		} `json:"prev_out"`
	} `json:"vin"`
	Vout []struct {
		Value        float64 `json:"value"`
		ScriptPubKey struct {
			Addresses []string `json:"addresses"`
		} `json:"scriptPubKey"`
	} `json:"vout"`
}

func getBitcoinTransaction(url, username, password, txHash string) (*BitcoinTransaction, error) {
	// Create RPC request payload
	type RpcRequest struct {
		Jsonrpc string        `json:"jsonrpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
		Id      string        `json:"id"`
	}

	rpcPayload := RpcRequest{
		Jsonrpc: "1.0",
		Method:  "getrawtransaction",
		Params:  []interface{}{txHash, true},
		Id:      "1",
	}

	payload, err := json.Marshal(rpcPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	var rpcResponse struct {
		Result BitcoinTransaction `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rpcResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &rpcResponse.Result, nil
}

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
