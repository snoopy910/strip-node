package sui

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
)

// WithdrawSuiNativeGetSignature returns transaction and dataToSign for
// native SUI withdrawal operation
func WithdrawSuiNativeGetSignature(
	rpcURL string,
	account string,
	amount string,
	recipient string,
) (string, string, error) {
	// Connect to Sui node
	cli, err := client.Dial(rpcURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to connect to Sui node: %v", err)
	}

	// Create context
	ctx := context.Background()

	// Convert amount to uint64
	amountBig, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return "", "", fmt.Errorf("invalid amount format: %s", amount)
	}

	// Convert addresses
	sender, err := sui_types.NewAddressFromHex(account)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse sender address: %v", err)
	}
	recipientAddr, err := sui_types.NewAddressFromHex(recipient)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse recipient address: %v", err)
	}

	// Get reference gas price
	gasPrice, err := cli.GetReferenceGasPrice(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to get gas price: %v", err)
	}

	// Calculate gas budget (adjust multiplier as needed)
	gasBudget := uint64(1000) * gasPrice.Uint64()

	// Get SUI coins
	coins, err := cli.GetSuiCoinsOwnedByAddress(ctx, *sender)
	if err != nil {
		return "", "", fmt.Errorf("failed to get coins: %v", err)
	}

	// Find coins with enough balance
	var selectedCoins []sui_types.ObjectID
	totalBalance := uint64(0)
	neededAmount := amountBig.Uint64() + gasBudget

	for _, coin := range coins {
		if coin.Balance.Uint64() == 0 {
			continue
		}
		selectedCoins = append(selectedCoins, coin.CoinObjectId)
		totalBalance += coin.Balance.Uint64()
		if totalBalance >= neededAmount {
			break
		}
	}

	if totalBalance < neededAmount {
		return "", "", fmt.Errorf("insufficient balance. needed %d, got %d", neededAmount, totalBalance)
	}

	// Create PaySui transaction
	txBytes, err := cli.PaySui(
		ctx,
		*sender,
		selectedCoins,
		[]sui_types.SuiAddress{*recipientAddr},
		[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amountBig.Uint64())},
		types.NewSafeSuiBigInt(gasBudget),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create PaySui transaction: %v", err)
	}

	return string(txBytes.TxBytes), string(txBytes.TxBytes), nil
}

// WithdrawSuiTokenGetSignature returns transaction and dataToSign for
// SUI token withdrawal operation
func WithdrawSuiTokenGetSignature(
	rpcURL string,
	account string,
	amount string,
	recipient string,
	coinType string,
) (string, string, error) {
	// Connect to Sui node
	cli, err := client.Dial(rpcURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to connect to Sui node: %v", err)
	}

	// Create context
	ctx := context.Background()

	// Convert amount to uint64
	amountBig, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return "", "", fmt.Errorf("invalid amount format: %s", amount)
	}

	// Convert addresses
	sender, err := sui_types.NewAddressFromHex(account)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse sender address: %v", err)
	}
	recipientAddr, err := sui_types.NewAddressFromHex(recipient)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse recipient address: %v", err)
	}

	// Get reference gas price
	gasPrice, err := cli.GetReferenceGasPrice(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to get gas price: %v", err)
	}

	// Calculate gas budget (adjust multiplier as needed)
	gasBudget := uint64(1000) * gasPrice.Uint64()

	// Get gas coins (SUI)
	gasCoins, err := cli.GetSuiCoinsOwnedByAddress(ctx, *sender)
	if err != nil {
		return "", "", fmt.Errorf("failed to get gas coins: %v", err)
	}
	if len(gasCoins) == 0 {
		return "", "", fmt.Errorf("no gas coins available")
	}

	// Find gas coin with enough balance
	var selectedGasCoin *sui_types.ObjectID
	for _, coin := range gasCoins {
		if coin.Balance.Uint64() >= gasBudget {
			selectedGasCoin = &coin.CoinObjectId
			break
		}
	}
	if selectedGasCoin == nil {
		return "", "", fmt.Errorf("no gas coin with sufficient balance for gas budget %d", gasBudget)
	}

	// Get token coins with pagination
	var tokenIDs []sui_types.ObjectID
	totalBalance := uint64(0)
	var cursor *sui_types.ObjectID

	for {
		coinPage, err := cli.GetCoins(ctx, *sender, &coinType, cursor, 50) // Get 50 coins at a time
		if err != nil {
			return "", "", fmt.Errorf("failed to get token coins: %v", err)
		}
		if len(coinPage.Data) == 0 {
			break
		}

		// Add coins until we have enough balance
		for _, coin := range coinPage.Data {
			if coin.Balance.Uint64() == 0 {
				continue
			}
			tokenIDs = append(tokenIDs, coin.CoinObjectId)
			totalBalance += coin.Balance.Uint64()
			if totalBalance >= amountBig.Uint64() {
				break
			}
		}

		// Stop paginating if either:
		// 1. We have enough balance
		// 2. No more pages according to API
		hasEnoughBalance := totalBalance >= amountBig.Uint64()
		hasNoMorePages := !coinPage.HasNextPage
		if hasEnoughBalance || hasNoMorePages {
			break
		}
		cursor = &coinPage.Data[len(coinPage.Data)-1].CoinObjectId
	}

	if totalBalance < amountBig.Uint64() {
		return "", "", fmt.Errorf("insufficient token balance. needed %d, got %d", amountBig.Uint64(), totalBalance)
	}

	// Create Pay transaction
	txBytes, err := cli.Pay(
		ctx,
		*sender,
		tokenIDs,
		[]sui_types.SuiAddress{*recipientAddr},
		[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amountBig.Uint64())},
		selectedGasCoin,
		types.NewSafeSuiBigInt(gasBudget),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create Pay transaction: %v", err)
	}

	return string(txBytes.TxBytes), string(txBytes.TxBytes), nil
}

// WithdrawSuiTxn submits transaction to withdraw assets and returns
// the txHash as the result
func WithdrawSuiTxn(
	rpcURL string,
	transaction string,
	publicKey string,
	signatureBase64 string,
) (string, error) {
	// Validate inputs
	if transaction == "" {
		return "", fmt.Errorf("transaction bytes cannot be empty")
	}
	if signatureBase64 == "" {
		return "", fmt.Errorf("signature cannot be empty")
	}
	if publicKey == "" {
		return "", fmt.Errorf("public key cannot be empty")
	}

	// Connect to Sui node
	cli, err := client.Dial(rpcURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to Sui node: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute the transaction
	resp, err := cli.ExecuteTransactionBlock(
		ctx,
		lib.Base64Data(transaction),
		[]any{signatureBase64},
		&types.SuiTransactionBlockResponseOptions{
			ShowEffects: true,
			ShowEvents: true,
		},
		types.TxnRequestTypeWaitForEffectsCert,
	)
	if err != nil {
		return "", fmt.Errorf("failed to execute transaction: %v", err)
	}

	return resp.Digest.String(), nil
}
