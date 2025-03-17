package sui

import (
	"context"
	"fmt"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
)

// CheckSuiTransactionConfirmed checks if a Sui transaction is confirmed
func CheckSuiTransactionConfirmed(chainId string, txHash string) (bool, error) {
	// Get chain URL
	chain, err := common.GetChain(chainId)
	if err != nil {
		return false, fmt.Errorf("failed to get chain info: %v", err)
	}

	// Connect to Sui node
	cli, err := client.Dial(chain.ChainUrl)
	if err != nil {
		return false, fmt.Errorf("failed to connect to Sui node: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Convert transaction hash to Digest
	digest, err := sui_types.NewDigest(txHash)
	if err != nil {
		return false, fmt.Errorf("invalid transaction hash format: %v", err)
	}

	// Get transaction
	resp, err := cli.GetTransactionBlock(
		ctx,
		*digest,
		types.SuiTransactionBlockResponseOptions{
			ShowEffects: true,
		},
	)
	if err != nil {
		return false, fmt.Errorf("failed to get transaction: %v", err)
	}

	// Check if effects exist
	if resp.Effects == nil {
		return false, fmt.Errorf("transaction effects not available")
	}

	// Check status using IsSuccess helper
	return resp.Effects.Data.IsSuccess(), nil
}

// SendSuiTransaction sends a signed Sui transaction
func SendSuiTransaction(
	serializedTx string,
	chainId string,
	keyCurve string,
	dataToSign string,
	signatureBase64 string,
) (string, error) {
	// Get chain URL
	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get chain info: %v", err)
	}

	// Connect to Sui node
	cli, err := client.Dial(chain.ChainUrl)
	if err != nil {
		return "", fmt.Errorf("failed to connect to Sui node: %v", err)
	}

	txBytes, _ := lib.NewBase64Data(serializedTx)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Submit the transaction
	resp, err := cli.ExecuteTransactionBlock(
		ctx,
		*txBytes,
		[]any{signatureBase64},
		&types.SuiTransactionBlockResponseOptions{
			ShowEffects: true,
		},
		types.TxnRequestTypeWaitForEffectsCert,
	)
	if err != nil {
		return "", fmt.Errorf("failed to execute transaction: %v", err)
	}

	// Check if effects exist
	if resp.Effects == nil {
		return "", fmt.Errorf("transaction effects not available")
	}

	// Check transaction status
	if !resp.Effects.Data.IsSuccess() {
		return "", fmt.Errorf("transaction failed: %v", resp.Effects.Data.V1.Status.Error)
	}

	return resp.Digest.String(), nil
}

// GetSuiTransfers gets transfers from a Sui transaction
// NOTE: only for transaction block with 1:1 transfer
func GetSuiTransfers(chainId string, txHash string) ([]common.Transfer, error) {
	// Get chain URL
	chain, err := common.GetChain(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain info: %v", err)
	}

	// Connect to Sui node
	cli, err := client.Dial(chain.ChainUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Sui node: %v", err)
	}

	// Convert transaction hash to Digest
	digest, err := sui_types.NewDigest(txHash)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction hash format: %v", err)
	}

	// Get transaction
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get transaction
	resp, err := cli.GetTransactionBlock(
		ctx,
		*digest,
		types.SuiTransactionBlockResponseOptions{
			ShowBalanceChanges: true,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	var (
		transfers    []common.Transfer
		from         string
		to           string
		amount       string
		token        string
		isNative     bool
		tokenAddress string
		scaledAmount string
	)

	// Extract transfers from balanceChanges
	if resp.BalanceChanges == nil {
		return transfers, nil
	}

	coinMetadata, err := cli.GetCoinMetadata(ctx, resp.BalanceChanges[1].CoinType)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	from = resp.BalanceChanges[0].Owner.AddressOwner.String()
	to = resp.BalanceChanges[1].Owner.AddressOwner.String()
	amount, err = getFormattedAmount(resp.BalanceChanges[1].Amount, coinMetadata.Decimals)
	if err != nil {
		return nil, fmt.Errorf("failed to format amount: %v", err)
	}
	token = coinMetadata.Symbol
	isNative = false
	tokenAddress = resp.BalanceChanges[1].CoinType
	scaledAmount = resp.BalanceChanges[1].Amount
	if resp.BalanceChanges[1].CoinType == SUI_TYPE {
		isNative = true
	}

	transfers = append(transfers, common.Transfer{
		From:         from,
		To:           to,
		Amount:       amount,
		Token:        token,
		IsNative:     isNative,
		TokenAddress: tokenAddress,
		ScaledAmount: scaledAmount,
	})

	return transfers, nil
}
