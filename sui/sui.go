package sui

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
)

const (
	DEFAULT_CONFIRMATIONS = 1 // Sui has instant finality
	SUI_DECIMALS         = 9
	SUI_TOKEN_SYMBOL     = "SUI"
	SUI_TYPE             = "0x2::sui::SUI"
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
			ShowEvents:  true,
			ShowInput:   true,
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

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Submit the transaction
	resp, err := cli.ExecuteTransactionBlock(
		ctx,
		lib.Base64Data(serializedTx),
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
			ShowEvents:  true,
			ShowInput:   true,
			ShowEffects: true,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	// Extract transfers from events
	var transfers []common.Transfer
	if resp.Events == nil {
		return transfers, nil
	}

	for _, event := range resp.Events {
		if event.Type != "0x2::coin::TransferEvent" {
			continue
		}

		// Get parsed JSON data
		parsedData, ok := event.ParsedJson.(map[string]interface{})
		if !ok {
			continue
		}

		// Parse amount
		amountStr, ok := parsedData["amount"].(string)
		if !ok {
			continue
		}
		amount, err := strconv.ParseUint(amountStr, 10, 64)
		if err != nil {
			continue
		}

		// Get coin type
		coinType, ok := parsedData["coin_type"].(string)
		if !ok {
			continue
		}

		// Get sender and recipient
		sender, ok := parsedData["sender"].(string)
		if !ok {
			continue
		}
		recipient, ok := parsedData["recipient"].(string)
		if !ok {
			continue
		}

		// Get token info
		tokenSymbol := getTokenSymbol(coinType)
		decimals := getTokenDecimals(coinType)

		// Create big.Int for amount with proper decimals
		amountBig := new(big.Int).SetUint64(amount)
		// Scale amount by decimals
		scaledBig := new(big.Int).Mul(amountBig, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))

		transfer := common.Transfer{
			From:         sender,
			To:           recipient,
			Amount:       scaledBig.String(),  // Use scaled amount
			Token:        tokenSymbol,
			IsNative:     coinType == SUI_TYPE,
			TokenAddress: coinType,
			ScaledAmount: amountBig.String(),  // Use unscaled amount
		}
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

// Helper functions
func getTokenSymbol(coinType string) string {
	if coinType == SUI_TYPE {
		return SUI_TOKEN_SYMBOL
	}
	// For other tokens, use the last part of the type
	// TODO: Implement token metadata lookup
	return coinType[len(coinType)-3:]
}

// Well-known token decimals
var tokenDecimals = map[string]int{
	SUI_TYPE:                     SUI_DECIMALS,  // Native SUI
	"0x2::devnet_nft::DevNetNFT": 0,           // NFTs have 0 decimals
	// Add more token types here as needed
}

// getTokenDecimals returns the number of decimals for a given coin type.
// For unknown tokens, it returns the default of 9 decimals (most common in Sui).
func getTokenDecimals(coinType string) int {
	// Check well-known tokens first
	if decimals, ok := tokenDecimals[coinType]; ok {
		return decimals
	}

	// For NFTs, return 0 decimals
	if strings.Contains(coinType, "nft") || strings.Contains(coinType, "NFT") {
		return 0
	}

	// For unknown tokens, return default Sui decimals
	return SUI_DECIMALS
}
