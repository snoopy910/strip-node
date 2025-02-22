package algorand

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

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

	if tx == nil {
		return "", fmt.Errorf("transaction is nil")
	}

	// Decode the signature (base64 encoded)
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
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
