package ripple

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/rubblelabs/ripple/data"
)

// WithdrawRippleNativeGetSignature creates a transaction for withdrawing native XRP and returns the transaction and data to sign
func WithdrawRippleNativeGetSignature(
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
	client, err := getClient(rpcURL)
	if err != nil {
		return "", "", err
	}

	amount2, err := data.NewAmount(amount)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse amount: %v", err)
	}

	bridgeAccount, err := data.NewAccountFromAddress(bridgeAddress)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse bridge address: %v", err)
	}

	userAccount, err := data.NewAccountFromAddress(userAddress)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse user address: %v", err)
	}

	// Get current network fee
	fee, err := client.Fee()
	if err != nil {
		return "", "", fmt.Errorf("failed to get network fee: %v", err)
	}

	// Create payment transaction
	payment := &data.Payment{
		TxBase: data.TxBase{
			TransactionType: data.PAYMENT,
			Account:         *bridgeAccount,
			Fee:             fee.Drops.BaseFee,
		},
		Destination: *userAccount,
		Amount:      *amount2,
	}

	// Get account sequence
	account, err := client.AccountInfo(*bridgeAccount)
	if err != nil {
		return "", "", fmt.Errorf("failed to get account info: %v", err)
	}

	payment.Sequence = *account.AccountData.Sequence

	// Serialize transaction
	tx := data.NewTransactionWithMetadata(data.PAYMENT)
	tx.Transaction = payment
	txBytes, err := tx.MarshalJSON()
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal transaction: %v", err)
	}

	hash, _, err := data.SigningHash(payment)
	if err != nil {
		return "", "", fmt.Errorf("failed to get signing hash: %v", err)
	}

	serializedTxn := hex.EncodeToString(txBytes)
	dataToSign := hex.EncodeToString(hash.Bytes())

	return serializedTxn, dataToSign, nil
}

// WithdrawRippleTokenGetSignature creates a transaction for withdrawing a non-native token and returns the transaction and data to sign
func WithdrawRippleTokenGetSignature(
	rpcURL string,
	bridgeAddress string,
	solverOutput string,
	userAddress string,
	// tokenCode string,
	tokenIssuer string,
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
	client, err := getClient(rpcURL)
	if err != nil {
		return "", "", err
	}

	bridgeAccount, err := data.NewAccountFromAddress(bridgeAddress)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse bridge address: %v", err)
	}

	userAccount, err := data.NewAccountFromAddress(userAddress)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse user address: %v", err)
	}

	// tokenCurrency, err := data.NewCurrency(tokenCode)
	// if err != nil {
	// 	return "", "", fmt.Errorf("failed to parse token code: %v", err)
	// }

	tokenIssuerAccount, err := data.NewAccountFromAddress(tokenIssuer)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse token issuer: %v", err)
	}

	value, err := data.NewValue(amount, true)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse amount: %v", err)
	}
	// Create token amount
	tokenAmount := data.Amount{
		Value: value,
		// Currency: tokenCurrency,
		Issuer: *tokenIssuerAccount,
	}

	// Get current network fee
	fee, err := client.Fee()
	if err != nil {
		return "", "", fmt.Errorf("failed to get network fee: %v", err)
	}

	// Create payment transaction
	payment := &data.Payment{
		TxBase: data.TxBase{
			TransactionType: data.PAYMENT,
			Account:         *bridgeAccount,
			Fee:             fee.Drops.BaseFee,
		},
		Destination: *userAccount,
		Amount:      tokenAmount,
	}

	// Get account sequence
	account, err := client.AccountInfo(*bridgeAccount)
	if err != nil {
		return "", "", fmt.Errorf("failed to get account info: %v", err)
	}

	payment.Sequence = *account.AccountData.Sequence

	// Serialize transaction
	tx := data.NewTransactionWithMetadata(data.PAYMENT)
	tx.Transaction = payment
	txBytes, err := tx.MarshalJSON()
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal transaction: %v", err)
	}

	hash, _, err := data.SigningHash(payment)
	if err != nil {
		return "", "", fmt.Errorf("failed to get signing hash: %v", err)
	}

	serializedTxn := hex.EncodeToString(txBytes)
	dataToSign := hex.EncodeToString(hash.Bytes())

	return serializedTxn, dataToSign, nil
}
