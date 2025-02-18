package algorand

import (
	"context"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

type SuggestedParamsRequester interface {
	Do(ctx context.Context) (types.SuggestedParams, error)
}

type PendingTransactionInformationRequester interface {
	Do(ctx context.Context) (*algod.PendingTransactionInformation, error)
}

type LookupTransactionRequester interface {
	Do(ctx context.Context) (*indexer.LookupTransaction, error)
}

type LookupAssetByIDRequester interface {
	Do(ctx context.Context) (*indexer.LookupAssetByID, error)
}

type AlgodClient interface {
	SendRawTransaction(txn []byte) (txid string, err error)
	PendingTransactionInformation(txid string) PendingTransactionInformationRequester
	SuggestedParams() SuggestedParamsRequester
}

type IndexerClient interface {
	LookupTransaction(txid string) LookupTransactionRequester
	LookupAssetByID(assetId uint64) LookupAssetByIDRequester
}

type Client interface {
	SendAlgorandTransaction(serializedTxn string, genesisHash string, signatureBase64 string) (string, error)
	GetAlgorandTransfers(genesisHash string, txnHash string) ([]common.Transfer, error)
	CheckAlgorandTransactionConfirmed(genesisHash string, txnHash string) (bool, error)
	WithdrawAlgorandNativeGetSignature(account string, amount string, recipient string) (string, *types.Transaction, error)
	WithdrawAlgorandASAGetSignature(account string, amount string, recipient string, assetId string) (string, *types.Transaction, error)
	WithdrawAlgorandTxn(signature string, tx *types.Transaction) (string, error)
}

type Clients struct {
	algodClient   *algod.Client
	indexerClient *indexer.Client
}

var _ Client = (*Clients)(nil)

func NewClients(genesisHash string, algodURL string, indexerURL string, createAlgod bool, createIndexer bool) (c *Clients, err error) {
	aURL := algodURL
	iURL := indexerURL

	nodeClient := &algod.Client{}
	indexerClient := &indexer.Client{}

	if createAlgod {
		if algodURL == "" {
			chain, err := common.GetChain(genesisHash)
			if err != nil {
				return nil, err
			}
			aURL = chain.ChainUrl
		}
		nodeClient, err = algod.MakeClient(aURL, "")
		if err != nil {
			return nil, fmt.Errorf("failed to create algod client: %v", err)
		}
	}

	if createIndexer {
		if indexerURL == "" {
			chain, err := common.GetChain(genesisHash)
			if err != nil {
				return nil, err
			}
			iURL = chain.IndexerUrl
		}
		// Create an indexer client (no API key needed for AlgoNode)
		indexerClient, err = indexer.MakeClient(iURL, "")
		if err != nil {
			return nil, fmt.Errorf("failed to create indexer client: %v", err)
		}
	}
	return &Clients{algodClient: nodeClient, indexerClient: indexerClient}, nil
}

func GetAlgorandTransfers(genesisHash string, txnHash string) ([]common.Transfer, error) {
	client, err := NewClients(genesisHash, "", "", false, true)
	if err != nil {
		return nil, err
	}
	return client.GetAlgorandTransfers(genesisHash, txnHash)
}

// GetAlgorandTransfers retrieves transfer information from an Algorand transaction
// It handles both native ALGO transfers and ASA (Algorand Standard Asset) transfers
func (client *Clients) GetAlgorandTransfers(genesisHash string, txnHash string) ([]common.Transfer, error) {

	// Create context with timeout
	// why 10 seconds? Needs to be computed ?
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Look up the transaction
	txnResponse, err := client.indexerClient.LookupTransaction(txnHash).Do(ctx)
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

		_, asset, err := client.indexerClient.LookupAssetByID(txn.AssetTransferTransaction.AssetId).Do(ctx)
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

func CheckAlgorandTransactionConfirmed(genesisHash string, txnHash string) (bool, error) {
	client, err := NewClients(genesisHash, "", "", true, true)
	if err != nil {
		return false, err
	}
	return client.CheckAlgorandTransactionConfirmed(genesisHash, txnHash)
}

// CheckAlgorandTransactionConfirmed checks if an Algorand transaction is confirmed
// It first tries the Algod API, then falls back to the Indexer if needed
func (client *Clients) CheckAlgorandTransactionConfirmed(genesisHash string, txnHash string) (bool, error) {
	// First try using native Algod API (Priority 1)
	// Get pending transaction information
	pendingTxn, _, err := client.algodClient.PendingTransactionInformation(txnHash).Do(context.Background())
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
	txnResponse, err := client.indexerClient.LookupTransaction(txnHash).Do(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed to lookup transaction: %v", err)
	}

	// If we can find the transaction in the indexer, it means it's confirmed
	// The indexer only indexes confirmed transactions
	return txnResponse.Transaction.ConfirmedRound > 0, nil
}

func SendAlgorandTransaction(serializedTxn string, genesisHash string, signatureBase64 string) (string, error) {
	client, err := NewClients(genesisHash, "", "", true, false)
	if err != nil {
		return "", err
	}
	return client.SendAlgorandTransaction(serializedTxn, genesisHash, signatureBase64)
}

// SendAlgorandTransaction sends a signed Algorand transaction to the network
func (client *Clients) SendAlgorandTransaction(serializedTxn string, genesisHash string, signatureBase64 string) (string, error) {
	// Decode the serialized transaction (base32 encoded)
	txnBytes, err := base32.StdEncoding.DecodeString(serializedTxn)
	if err != nil {
		return "", fmt.Errorf("failed to decode serialized transaction: %v", err)
	}

	// Decode the signature (base32 encoded)
	sigBytes, err := base32.StdEncoding.DecodeString(signatureBase64)
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
	txid, err := client.algodClient.SendRawTransaction(signedTxnBytes).Do(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return txid, nil
}

// bridge withdraw

func WithdrawAlgorandNativeGetSignature(
	algodURL string,
	account string,
	amount string,
	recipient string,
) (string, *types.Transaction, error) {
	client, err := NewClients("", algodURL, "", true, false)
	if err != nil {
		return "", nil, err
	}
	return client.WithdrawAlgorandNativeGetSignature(account, amount, recipient)

}

func (client *Clients) WithdrawAlgorandNativeGetSignature(
	account string,
	amount string,
	recipient string,
) (string, *types.Transaction, error) {

	sp, err := client.algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		return "", nil, fmt.Errorf("failed to get suggested params: %w", err)
	}

	amt, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return "", nil, fmt.Errorf("invalid amount: %w", err)
	}

	tx, err := future.MakePaymentTxn(account, recipient, amt, nil, "", sp)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create payment transaction: %w", err)
	}

	// algorand sdk v1 doesn't support tx.ID()
	txHash := crypto.TransactionID(tx)

	return hex.EncodeToString(txHash), &tx, nil
}

func WithdrawAlgorandASAGetSignature(
	algodURL string,
	account string,
	amount string,
	recipient string,
	assetId string,
) (string, *types.Transaction, error) {
	client, err := NewClients("", algodURL, "", true, false)
	if err != nil {
		return "", nil, err
	}
	return client.WithdrawAlgorandASAGetSignature(account, amount, recipient, assetId)
}

func (client *Clients) WithdrawAlgorandASAGetSignature(
	account string,
	amount string,
	recipient string,
	assetId string,
) (string, *types.Transaction, error) {

	sp, err := client.algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		return "", nil, fmt.Errorf("failed to get suggested params: %w", err)
	}

	amt, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return "", nil, fmt.Errorf("invalid amount: %w", err)
	}

	assetID, err := strconv.ParseUint(assetId, 10, 64)
	if err != nil {
		return "", nil, fmt.Errorf("invalid asset id: %w", err)
	}

	tx, err := future.MakeAssetTransferTxn(account, recipient, amt, []byte(""), sp, "", assetID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create asset transfer transaction: %w", err)
	}

	// algorand sdk v1 doesn't support tx.ID()
	txHash := crypto.TransactionID(tx)

	return hex.EncodeToString(txHash), &tx, nil
}

func WithdrawAlgorandTxn(
	algodURL string,
	signature string,
	tx *types.Transaction,
) (string, error) {
	client, err := NewClients("", algodURL, "", true, false)
	if err != nil {
		return "", err
	}
	return client.WithdrawAlgorandTxn(signature, tx)
}

func (client *Clients) WithdrawAlgorandTxn(
	signature string,
	tx *types.Transaction,
) (string, error) {

	// Decode the signature (base32 encoded)
	sigBytes, err := base32.StdEncoding.DecodeString(signature)
	if err != nil {
		return "", fmt.Errorf("failed to decode signature: %v", err)
	}

	// Create a signed transaction with the provided signature
	// In go1.19, we can't convert sigBytes directly to types.Signature
	var sig types.Signature
	copy(sig[:], sigBytes)
	signedTxn := types.SignedTxn{
		Txn: *tx,
		Sig: sig,
	}

	// Encode the signed transaction using msgpack
	signedTxnBytes := msgpack.Encode(signedTxn)

	// Send the transaction
	txid, err := client.algodClient.SendRawTransaction(signedTxnBytes).Do(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return txid, nil
}
