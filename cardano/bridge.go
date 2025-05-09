package cardano

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/echovl/cardano-go"
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

	txBuilder := cardano.NewTxBuilder(&cardano.ProtocolParams{})
	amountStr, ok := solverData["amount"].(string)
	if !ok {
		return "", "", fmt.Errorf("amount not found in solver output")
	}

	amount, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse amount: %v", err)
	}

	address, err := cardano.NewAddress(userAddress)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse address: %v", err)
	}

	inputTxHash := strings.ToLower(solverData["txHash"].(string))
	txBuilder.AddInputs(&cardano.TxInput{
		TxHash: cardano.Hash32(inputTxHash),
		Amount: &cardano.Value{
			Coin: cardano.Coin(amount),
		},
	})
	txBuilder.AddOutputs(&cardano.TxOutput{
		Address: address,
		Amount: &cardano.Value{
			Coin: cardano.Coin(amount),
		},
	})

	tx, err := txBuilder.Build()
	if err != nil {
		return "", "", fmt.Errorf("failed to build transaction: %v", err)
	}
	hash, err := tx.Hash()
	if err != nil {
		return "", "", fmt.Errorf("failed to hash transaction: %v", err)
	}
	serializedTxn := tx.Hex()
	dataToSign := hash.String()

	return serializedTxn, dataToSign, nil
}

// WithdrawCardanoTokenGetSignature creates a transaction for withdrawing a native token and returns the transaction and data to sign
func WithdrawCardanoTokenGetSignature(
	rpcURL string,
	bridgeAddress string,
	solverOutput string,
	userAddress string,
	token string,
) (string, string, error) {
	// Parse solver output
	var solverData map[string]interface{}
	if err := json.Unmarshal([]byte(solverOutput), &solverData); err != nil {
		return "", "", fmt.Errorf("failed to parse solver output: %v", err)
	}

	txBuilder := cardano.NewTxBuilder(&cardano.ProtocolParams{})
	amountStr, ok := solverData["amount"].(string)
	if !ok {
		return "", "", fmt.Errorf("amount not found in solver output")
	}

	amount, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse amount: %v", err)
	}

	address, err := cardano.NewAddress(userAddress)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse address: %v", err)
	}

	inputTxHash := strings.ToLower(solverData["txHash"].(string))

	// FIX: This is a temporary fillin and should be addressed
	policyID := cardano.NewPolicyIDFromHash(cardano.Hash28(token))
	assetName := cardano.NewAssetName(token)
	txBuilder.AddInputs(&cardano.TxInput{
		TxHash: cardano.Hash32(inputTxHash),
		Amount: &cardano.Value{
			MultiAsset: cardano.NewMultiAsset().Set(policyID, cardano.NewAssets().Set(assetName, cardano.BigNum(amount))),
		},
	})
	txBuilder.AddOutputs(&cardano.TxOutput{
		Address: address,
		Amount: &cardano.Value{
			MultiAsset: cardano.NewMultiAsset().Set(policyID, cardano.NewAssets().Set(assetName, cardano.BigNum(amount))),
		},
	})

	tx, err := txBuilder.Build()
	if err != nil {
		return "", "", fmt.Errorf("failed to build transaction: %v", err)
	}
	hash, err := tx.Hash()
	if err != nil {
		return "", "", fmt.Errorf("failed to hash transaction: %v", err)
	}
	serializedTxn := tx.Hex()
	dataToSign := hash.String()

	return serializedTxn, dataToSign, nil
}
