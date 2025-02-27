package cardano

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/blockfrost/blockfrost-go"
)

// WithdrawCardanoNativeGetSignature creates a transaction for withdrawing native ADA and returns the transaction and data to sign
func WithdrawCardanoNativeGetSignature(
	rpcURL string,
	bridgeAddress string,
	solverOutput string,
	userAddress string,
) (string, string, error) {
	// Parse solver output to get amount
	var solverData map[string]interface{}
	if err := json.Unmarshal([]byte(solverOutput), &solverData); err != nil {
		return "", "", fmt.Errorf("failed to parse solver output: %v", err)
	}

	amount, ok := solverData["amount"].(string)
	if !ok {
		return "", "", fmt.Errorf("amount not found in solver output")
	}

	// Initialize client
	// client, err := getClient(rpcURL)
	// if err != nil {
	// 	return "", "", err
	// }

	// Get UTXOs for bridge address
	// utxos, err := client.AddressUTXOs(context.Background(), bridgeAddress, blockfrost.APIQueryParams{})
	// if err != nil {
	// 	return "", "", fmt.Errorf("failed to get UTXOs: %v", err)
	// }

	// Build transaction
	tx := &blockfrost.TransactionUTXOs{
		Inputs: []blockfrost.TransactionInput{
			{
				Address: bridgeAddress,
				Amount:  []blockfrost.TxAmount{{Unit: "lovelace", Quantity: amount}},
			},
		},
		Outputs: []blockfrost.TransactionOutput{
			{
				Address: userAddress,
				Amount:  []blockfrost.TxAmount{{Unit: "lovelace", Quantity: amount}},
			},
		},
	}

	// Serialize transaction
	txBytes, err := json.Marshal(tx)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal transaction: %v", err)
	}

	serializedTxn := string(txBytes)
	dataToSign := strings.TrimPrefix(serializedTxn, "0x")

	return serializedTxn, dataToSign, nil
}

// WithdrawCardanoTokenGetSignature creates a transaction for withdrawing a native token and returns the transaction and data to sign
func WithdrawCardanoTokenGetSignature(
	rpcURL string,
	bridgeAddress string,
	solverOutput string,
	userAddress string,
	policyID string,
) (string, string, error) {
	// Parse solver output
	var solverData map[string]interface{}
	if err := json.Unmarshal([]byte(solverOutput), &solverData); err != nil {
		return "", "", fmt.Errorf("failed to parse solver output: %v", err)
	}

	amount, ok := solverData["amount"].(string)
	if !ok {
		return "", "", fmt.Errorf("amount not found in solver output")
	}

	// Initialize client
	// client, err := getClient(rpcURL)
	// if err != nil {
	// 	return "", "", err
	// }

	// Get UTXOs for bridge address
	// utxos, err := client.AddressUTXOs(context.Background(), bridgeAddress, blockfrost.APIQueryParams{})
	// if err != nil {
	// 	return "", "", fmt.Errorf("failed to get UTXOs: %v", err)
	// }

	// Build token unit identifier

	// Build transaction
	tx := &blockfrost.TransactionUTXOs{
		Inputs: []blockfrost.TransactionInput{
			{
				Address: bridgeAddress,
				Amount: []blockfrost.TxAmount{
					{Unit: policyID, Quantity: amount},
					{Unit: "lovelace", Quantity: "2000000"}, // Min ADA requirement
				},
			},
		},
		Outputs: []blockfrost.TransactionOutput{
			{
				Address: userAddress,
				Amount: []blockfrost.TxAmount{
					{Unit: policyID, Quantity: amount},
					{Unit: "lovelace", Quantity: "2000000"}, // Min ADA requirement
				},
			},
		},
	}

	// Serialize transaction
	txBytes, err := json.Marshal(tx)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal transaction: %v", err)
	}

	serializedTxn := string(txBytes)
	dataToSign := strings.TrimPrefix(serializedTxn, "0x")

	return serializedTxn, dataToSign, nil
}
