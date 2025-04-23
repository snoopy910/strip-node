package blockchains

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/libs"
	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
)

const (
	SUI_TYPE = "0x2::sui::SUI"
)

// NewSuiBlockchain creates a new Stellar blockchain instance
func NewSuiBlockchain(networkType NetworkType) (IBlockchain, error) {
	network := Network{
		networkType: networkType,
		nodeURL:     "https://fullnode.mainnet.sui.io:443",
		networkID:   "mainnet",
	}

	if networkType == Testnet {
		network.nodeURL = "https://fullnode.devnet.sui.io:443"
		network.networkID = "testnet"
	}

	client, err := client.Dial(network.nodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Sui node: %v", err)
	}
	return &SuiBlockchain{
		BaseBlockchain: BaseBlockchain{
			chainName:       Sui,
			network:         network,
			keyCurve:        common.CurveEddsa,
			signingEncoding: "base64",
			decimals:        7,
			opTimeout:       time.Second * 10,
		},
		client: client,
	}, nil
}

// This is a type assertion to ensure that the SuiBlockchain implements the IBlockchain interface
var _ IBlockchain = &SuiBlockchain{}

// SuiBlockchain implements the IBlockchain interface for Stellar
type SuiBlockchain struct {
	BaseBlockchain
	client *client.Client
}

func (b *SuiBlockchain) BroadcastTransaction(txn string, signatureBase64 string, _ *string) (string, error) {
	txBytes, _ := lib.NewBase64Data(txn)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Submit the transaction
	resp, err := b.client.ExecuteTransactionBlock(
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

func (b *SuiBlockchain) GetTransfers(txHash string, address *string) ([]common.Transfer, error) {
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
	resp, err := b.client.GetTransactionBlock(
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

	coinMetadata, err := b.client.GetCoinMetadata(ctx, resp.BalanceChanges[1].CoinType)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	from = resp.BalanceChanges[0].Owner.AddressOwner.String()
	to = resp.BalanceChanges[1].Owner.AddressOwner.String()
	amount, err = getFormattedAmount(resp.BalanceChanges[1].Amount, int(coinMetadata.Decimals))
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

func (b *SuiBlockchain) IsTransactionBroadcastedAndConfirmed(txHash string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Convert transaction hash to Digest
	digest, err := sui_types.NewDigest(txHash)
	if err != nil {
		return false, fmt.Errorf("invalid transaction hash format: %v", err)
	}

	// Get transaction
	resp, err := b.client.GetTransactionBlock(
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

func (b *SuiBlockchain) BuildWithdrawTx(account string,
	solverOutput string,
	recipient string,
	tokenAddress *string,
) (string, string, error) {
	ctx := context.Background()
	var solverData map[string]interface{}
	if err := json.Unmarshal([]byte(solverOutput), &solverData); err != nil {
		return "", "", fmt.Errorf("failed to parse solver output: %v", err)
	}

	amount, ok := solverData["amount"].(string)
	if !ok {
		return "", "", fmt.Errorf("amount not found in solver output")
	}
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
	gasPrice, err := b.client.GetReferenceGasPrice(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to get gas price: %v", err)
	}

	// Calculate gas budget (adjust multiplier as needed)
	gasBudget := uint64(1000) * gasPrice.Uint64()

	if tokenAddress == nil {

		// Get SUI coins
		coins, err := b.client.GetSuiCoinsOwnedByAddress(ctx, *sender)
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
		txBytes, err := b.client.PaySui(
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

	// Get gas coins (SUI)
	gasCoins, err := b.client.GetSuiCoinsOwnedByAddress(ctx, *sender)
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
		coinPage, err := b.client.GetCoins(ctx, *sender, tokenAddress, cursor, 50) // Get 50 coins at a time
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
	txBytes, err := b.client.Pay(
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

func (b *SuiBlockchain) RawPublicKeyBytesToAddress(pkBytes []byte, networkType NetworkType) (string, error) {
	return "", errors.New("RawPublicKeyBytesToAddress not implemented")
}

func (b *SuiBlockchain) RawPublicKeyToPublicKeyStr(pkBytes []byte) (string, error) {
	return "", errors.New("RawPublicKeyToPublicKeyStr not implemented")
}

func (b *SuiBlockchain) ExtractDestinationAddress(operation *libs.Operation) (string, error) {
	var tx sui_types.TransactionData
	txBytes, err := base64.StdEncoding.DecodeString(*operation.SerializedTxn)
	if err != nil {
		return "", fmt.Errorf("error decoding Sui transaction", err)
	}
	if err := json.Unmarshal(txBytes, &tx); err != nil {
		return "", fmt.Errorf("error parsing Sui transaction", err)
	}
	if len(tx.V1.Kind.ProgrammableTransaction.Inputs) < 1 {
		return "", fmt.Errorf("wrong format sui transaction")
	}
	destAddress := string(*tx.V1.Kind.ProgrammableTransaction.Inputs[0].Pure)
	return destAddress, nil
}
