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
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/test-go/testify/mock"
)

type SuggestedParamsRequester interface {
	Do(ctx context.Context) (types.SuggestedParams, error)
}

type MockSuggestedParamsRequester struct {
	mock.Mock
}

func (m *MockSuggestedParamsRequester) Do(ctx context.Context) (types.SuggestedParams, error) {
	args := m.Called(ctx)
	return args.Get(0).(types.SuggestedParams), args.Error(1)
}

type PendingTransactionInformationRequester interface {
	Do(ctx context.Context) (response models.PendingTransactionInfoResponse, stxn types.SignedTxn, err error)
}

type MockPendingTransactionInformationRequester struct {
	mock.Mock
}

func (m *MockPendingTransactionInformationRequester) Do(ctx context.Context) (models.PendingTransactionInfoResponse, types.SignedTxn, error) {
	args := m.Called(ctx)
	return args.Get(0).(models.PendingTransactionInfoResponse), args.Get(1).(types.SignedTxn), args.Error(2)
}

type LookupTransactionRequester interface {
	Do(ctx context.Context) (models.TransactionResponse, error)
}

type MockLookupTransactionRequester struct {
	mock.Mock
}

func (m *MockLookupTransactionRequester) Do(ctx context.Context) (models.Transaction, error) {
	args := m.Called(ctx)
	return args.Get(0).(models.Transaction), args.Error(1)
}

type LookupAssetByIDRequester interface {
	Do(ctx context.Context) (validRound uint64, result models.Asset, err error)
}

type MockLookupAssetByIDRequester struct {
	mock.Mock
}

func (m *MockLookupAssetByIDRequester) Do(ctx context.Context) (validRound uint64, result models.Asset, err error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Get(1).(models.Asset), args.Error(2)
}

type SendRawTransactionRequester interface {
	Do(ctx context.Context) (string, error)
}

type MockSendRawTransactionRequester struct {
	mock.Mock
}

func (m *MockSendRawTransactionRequester) Do(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.Get(0).(string), args.Error(1)
}

type AlgodClient interface {
	SendRawTransaction(txn []byte) SendRawTransactionRequester
	PendingTransactionInformation(txid string) PendingTransactionInformationRequester
	SuggestedParams() SuggestedParamsRequester
}

type IndexerClient interface {
	LookupTransaction(txid string) LookupTransactionRequester
	LookupAssetByID(assetId uint64) LookupAssetByIDRequester
}

type MockAlgodClient struct {
	mock.Mock
}

func (m *MockAlgodClient) SendRawTransaction(txn []byte) SendRawTransactionRequester {
	args := m.Called(txn)
	return args.Get(0).(SendRawTransactionRequester)
}

func (m *MockAlgodClient) PendingTransactionInformation(txid string) PendingTransactionInformationRequester {
	args := m.Called(txid)
	return args.Get(0).(PendingTransactionInformationRequester)
}

func (m *MockAlgodClient) SuggestedParams() SuggestedParamsRequester {
	args := m.Called()
	return args.Get(0).(SuggestedParamsRequester)
}

type MockIndexerClient struct {
	mock.Mock
}

func (m *MockIndexerClient) LookupTransaction(txid string) LookupTransactionRequester {
	args := m.Called(txid)
	return args.Get(0).(LookupTransactionRequester)
}

func (m *MockIndexerClient) LookupAssetByID(assetId uint64) LookupAssetByIDRequester {
	args := m.Called(assetId)
	return args.Get(0).(LookupAssetByIDRequester)
}

type MockClients struct {
	mockAlgod   *MockAlgodClient
	mockIndexer *MockIndexerClient
}

var _ Client = (*MockClients)(nil)

func (mockClient *MockClients) SendAlgorandTransaction(serializedTxn string, genesisHash string, signatureBase64 string) (string, error) {
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
	txid, err := mockClient.mockAlgod.SendRawTransaction(signedTxnBytes).Do(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return txid, nil
}

func (mockClient *MockClients) GetAlgorandTransfers(genesisHash string, txnHash string) ([]common.Transfer, error) {

	// Create context with timeout
	// why 10 seconds? Needs to be computed ?
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Look up the transaction
	txnResponse, err := mockClient.mockIndexer.LookupTransaction(txnHash).Do(ctx)
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

		_, asset, err := mockClient.mockIndexer.LookupAssetByID(txn.AssetTransferTransaction.AssetId).Do(ctx)
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

func (mockClient *MockClients) CheckAlgorandTransactionConfirmed(genesisHash string, txnHash string) (bool, error) {
	// First try using native Algod API (Priority 1)
	// Get pending transaction information
	pendingTxn, _, err := mockClient.mockAlgod.PendingTransactionInformation(txnHash).Do(context.Background())
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
	txnResponse, err := mockClient.mockIndexer.LookupTransaction(txnHash).Do(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed to lookup transaction: %v", err)
	}

	// If we can find the transaction in the indexer, it means it's confirmed
	// The indexer only indexes confirmed transactions
	return txnResponse.Transaction.ConfirmedRound > 0, nil
}

func (mockClient *MockClients) WithdrawAlgorandNativeGetSignature(
	account string,
	amount string,
	recipient string,
) (string, *types.Transaction, error) {

	sp, err := mockClient.mockAlgod.SuggestedParams().Do(context.Background())
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

func (mockClient *MockClients) WithdrawAlgorandASAGetSignature(
	account string,
	amount string,
	recipient string,
	assetId string,
) (string, *types.Transaction, error) {

	sp, err := mockClient.mockAlgod.SuggestedParams().Do(context.Background())
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

func (mockClient *MockClients) WithdrawAlgorandTxn(
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
	txid, err := mockClient.mockAlgod.SendRawTransaction(signedTxnBytes).Do(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return txid, nil
}
