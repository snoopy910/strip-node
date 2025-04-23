package blockchains

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/libs"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/stellar/go/clients/horizonclient"
)

// NewAlgorandBlockchain creates a new Stellar blockchain instance
func NewAlgorandBlockchain(networkType NetworkType) (IBlockchain, error) {
	indexerURL := "https://mainnet-idx.4160.nodely.dev"
	network := Network{
		networkType: networkType,
		nodeURL:     horizonclient.DefaultPublicNetClient.HorizonURL,
		networkID:   "mainnet",
	}

	if networkType == Testnet {
		network.nodeURL = "https://mainnet-api.4160.nodely.dev"
		network.networkID = "testnet"
		indexerURL = "https://testnet-idx.4160.nodely.dev"
	}

	algodClient, err := algod.MakeClient(network.nodeURL, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create algod client: %v", err)
	}

	indexerClient, err := indexer.MakeClient(indexerURL, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create indexer client: %v", err)
	}

	return &AlgorandBlockchain{
		BaseBlockchain: BaseBlockchain{
			chainName:       Algorand,
			network:         network,
			keyCurve:        common.CurveEddsa,
			signingEncoding: "base64",
			decimals:        6,
			opTimeout:       time.Second * 10,
		},
		algodClient:   algodClient,
		indexerClient: indexerClient,
	}, nil
}

// This is a type assertion to ensure that the AlgorandBlockchain implements the IBlockchain interface
var _ IBlockchain = &AlgorandBlockchain{}

// AlgorandBlockchain implements the IBlockchain interface for Stellar
type AlgorandBlockchain struct {
	BaseBlockchain
	algodClient   *algod.Client
	indexerClient *indexer.Client
}

func (b *AlgorandBlockchain) BroadcastTransaction(serializedTxn string, signature string, _ *string) (string, error) {
	// Decode the serialized transaction (base64 encoded)
	txnBytes, err := base64.StdEncoding.DecodeString(serializedTxn)
	if err != nil {
		return "", fmt.Errorf("failed to decode serialized transaction: %v", err)
	}

	// Decode the signature (base64 encoded)
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return "", fmt.Errorf("failed to decode signature: %v", err)
	}

	// Deserialize the transaction using msgpack
	var txn types.Transaction
	err = msgpack.Decode(txnBytes, &txn)
	if err != nil {
		return "", fmt.Errorf("failed to deserialize transaction: %v", err)
	}

	// Create a signed transaction with the provided signature
	// In go1.19, we can't convert sigBytes directly to types.Signature
	var sig types.Signature
	copy(sig[:], sigBytes)
	signedTxn := types.SignedTxn{
		Txn: txn,
		Sig: sig,
	}

	// Encode the signed transaction using msgpack
	signedTxnBytes := msgpack.Encode(signedTxn)

	// Send the transaction
	txid, err := b.algodClient.SendRawTransaction(signedTxnBytes).Do(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}
	return txid, nil
}

func (b *AlgorandBlockchain) GetTransfers(txHash string, address *string) ([]common.Transfer, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Look up the transaction
	txnResponse, err := b.indexerClient.LookupTransaction(txHash).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup transaction: %v", err)
	}

	// Validate transaction
	if txnResponse.Transaction.ConfirmedRound == 0 {
		return nil, fmt.Errorf("transaction not found or not confirmed")
	}

	txn := txnResponse.Transaction
	var transfers []common.Transfer

	addressSender, err := types.DecodeAddress(txn.Sender)
	if err != nil {
		return nil, fmt.Errorf("failed to decode sender address: %v", err)
	}
	if addressSender.IsZero() {
		return nil, fmt.Errorf("invalid sender address")
	}

	switch txn.Type {
	// Handle native ALGO transfer
	case string(types.PaymentTx):
		// Validate receiver address
		addressReceiver, err := types.DecodeAddress(txn.PaymentTransaction.Receiver)
		if err != nil {
			return nil, fmt.Errorf("failed to decode receiver address: %v", err)
		}
		if addressReceiver.IsZero() {
			return nil, fmt.Errorf("invalid receiver address")
		}

		// ALGO amounts are in microAlgos (1 ALGO = 1,000,000 microAlgos)
		amount := float64(txn.PaymentTransaction.Amount) / 1_000_000

		transfers = append(transfers, common.Transfer{
			From:   txn.Sender,
			To:     txn.PaymentTransaction.Receiver,
			Amount: fmt.Sprintf("%.6f", amount),
			Token:  "ALGO", // Native token
		})

	// Handle ASA (Algorand Standard Asset) transfer
	case string(types.AssetTransferTx):
		// Validate addresses
		addressReceiver, err := types.DecodeAddress(txn.AssetTransferTransaction.Receiver)
		if err != nil {
			return nil, fmt.Errorf("failed to decode receiver address: %v", err)
		}
		if addressReceiver.IsZero() {
			return nil, fmt.Errorf("invalid receiver address")
		}

		// Get asset info for decimals
		_, asset, err := b.indexerClient.LookupAssetByID(txn.AssetTransferTransaction.AssetId).Do(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to lookup asset info: %v", err)
		}

		// Calculate amount with proper decimals
		decimals := asset.Params.Decimals
		amount := float64(txn.AssetTransferTransaction.Amount) / math.Pow10(int(decimals))

		transfers = append(transfers, common.Transfer{
			From:   txn.Sender,
			To:     txn.AssetTransferTransaction.Receiver,
			Amount: fmt.Sprintf("%.*f", decimals, amount),
			Token:  fmt.Sprintf("%d", asset.Index), // ASA ID as token identifier
		})
	default:
		return nil, fmt.Errorf("unsupported transaction type: %v", txn.Type)
	}

	return transfers, nil
}

func (b *AlgorandBlockchain) IsTransactionBroadcastedAndConfirmed(txHash string) (bool, error) {
	pendingTxn, _, err := b.algodClient.PendingTransactionInformation(txHash).Do(context.Background())
	if err == nil {
		// If confirmed round is non-zero, transaction is confirmed
		if pendingTxn.ConfirmedRound > 0 {
			return true, nil
		}
		// If pool error is empty and confirmed round is zero, transaction is still pending
		if pendingTxn.PoolError == "" {
			return false, nil
		}
	}

	// Fallback to Indexer if Algod fails or transaction not found (Priority 2)
	// Look up the transaction
	txnResponse, err := b.indexerClient.LookupTransaction(txHash).Do(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed to lookup transaction: %v", err)
	}

	// If we can find the transaction in the indexer, it means it's confirmed
	// The indexer only indexes confirmed transactions
	return txnResponse.Transaction.ConfirmedRound > 0, nil
}

func (b *AlgorandBlockchain) BuildWithdrawTx(account string,
	solverOutput string,
	recipient string,
	tokenAddress *string,
) (string, string, error) {
	var solverData map[string]interface{}
	if err := json.Unmarshal([]byte(solverOutput), &solverData); err != nil {
		return "", "", fmt.Errorf("failed to parse solver output: %v", err)
	}
	amountStr, ok := solverData["amount"].(string)
	if !ok {
		return "", "", fmt.Errorf("amount not found in solver output")
	}

	amount, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse amount: %v", err)
	}

	if tokenAddress == nil {
		sp, err := b.algodClient.SuggestedParams().Do(context.Background())
		if err != nil {
			return "", "", fmt.Errorf("failed to get suggested params: %w", err)
		}

		tx, err := future.MakePaymentTxn(account, recipient, amount, nil, "", sp)
		if err != nil {
			return "", "", fmt.Errorf("failed to create payment transaction: %w", err)
		}

		// algorand sdk v1 doesn't support tx.ID()
		txHash := crypto.TransactionID(tx)

		// Encode the transaction as msgpack
		serializedTxn := msgpack.Encode(tx)

		return base64.StdEncoding.EncodeToString(serializedTxn), hex.EncodeToString(txHash), nil
	}

	sp, err := b.algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		return "", "", fmt.Errorf("failed to get suggested params: %w", err)
	}

	assetID, err := strconv.ParseUint(*tokenAddress, 10, 64)
	if err != nil {
		return "", "", fmt.Errorf("invalid asset id: %w", err)
	}

	tx, err := future.MakeAssetTransferTxn(account, recipient, amount, []byte(""), sp, "", assetID)
	if err != nil {
		return "", "", fmt.Errorf("failed to create asset transfer transaction: %w", err)
	}

	// algorand sdk v1 doesn't support tx.ID()
	txHash := crypto.TransactionID(tx)

	// Encode the transaction as msgpack
	serializedTxn := msgpack.Encode(tx)

	return base64.StdEncoding.EncodeToString(serializedTxn), hex.EncodeToString(txHash), nil

}

func (b *AlgorandBlockchain) RawPublicKeyBytesToAddress(pkBytes []byte, networkType NetworkType) (string, error) {
	return "", errors.New("RawPublicKeyBytesToAddress not implemented")
}

func (b *AlgorandBlockchain) RawPublicKeyToPublicKeyStr(pkBytes []byte) (string, error) {
	return "", errors.New("RawPublicKeyToPublicKeyStr not implemented")
}

func (b *AlgorandBlockchain) ExtractDestinationAddress(operation *libs.Operation) (string, error) {
	txnBytes, err := base64.StdEncoding.DecodeString(*operation.SerializedTxn)
	if err != nil {
		return "", fmt.Errorf("failed to decode serialized transaction", err)
	}
	var txn types.Transaction
	err = msgpack.Decode(txnBytes, &txn)
	if err != nil {
		return "", fmt.Errorf("failed to deserialize transaction", err)
	}
	destAddress := ""
	if txn.Type == types.PaymentTx {
		destAddress = txn.PaymentTxnFields.Receiver.String()
	} else if txn.Type == types.AssetTransferTx {
		destAddress = txn.AssetTransferTxnFields.AssetReceiver.String()
	} else {
		return "", fmt.Errorf("unknown transaction type", txn.Type)
	}
	return destAddress, nil
}
