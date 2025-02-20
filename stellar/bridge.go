package stellar

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

// WithdrawNativeGetSignature creates a transaction for withdrawing native XLM and returns the transaction and data to sign
func WithdrawStellarNativeGetSignature(
	client *horizonclient.Client,
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

	// Get bridge account to use as source account
	account, err := client.AccountDetail(horizonclient.AccountRequest{AccountID: bridgeAddress})
	if err != nil {
		hzErr, ok := err.(*horizonclient.Error)
		if ok {
			return "", "", fmt.Errorf("error getting bridge account: %v", hzErr.Problem)
		}

		return "", "", fmt.Errorf("error getting bridge account: %v", err)
	}

	// Set up time bounds
	tb := txnbuild.NewTimeout(300)

	// Build the transaction
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &account,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: userAddress,
					Amount:      amount,
					Asset:       txnbuild.NativeAsset{},
				},
			},
			BaseFee:       txnbuild.MinBaseFee,
			Memo:          nil,
			Preconditions: txnbuild.Preconditions{TimeBounds: tb},
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("error creating transaction: %v", err)
	}

	// Get the transaction in XDR format
	txeB64, err := tx.Base64()
	if err != nil {
		return "", "", fmt.Errorf("error encoding transaction: %v", err)
	}

	// Get the transaction hash to sign
	hash, err := tx.Hash(network.PublicNetworkPassphrase)
	if err != nil {
		return "", "", fmt.Errorf("error getting transaction hash: %v", err)
	}

	return txeB64, base64.StdEncoding.EncodeToString(hash[:]), nil
}

// WithdrawAssetGetSignature creates a transaction for withdrawing a non-native Stellar asset and returns the transaction and data to sign
func WithdrawStellarAssetGetSignature(
	client *horizonclient.Client,
	bridgeAddress string,
	solverOutput string,
	userAddress string,
	assetCode string,
	assetIssuer string,
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

	// Get bridge account to use as source account
	account, err := client.AccountDetail(horizonclient.AccountRequest{AccountID: bridgeAddress})
	if err != nil {
		hzErr, ok := err.(*horizonclient.Error)
		if ok {
			return "", "", fmt.Errorf("error getting bridge account: %v", hzErr.Problem)
		}

		return "", "", fmt.Errorf("error getting bridge account: %v", err)
	}

	// Set up time bounds
	tb := txnbuild.NewTimeout(300)

	// Build the transaction
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &account,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: userAddress,
					Amount:      amount,
					Asset:       txnbuild.CreditAsset{Code: assetCode, Issuer: assetIssuer},
				},
			},
			BaseFee:       txnbuild.MinBaseFee,
			Memo:          nil,
			Preconditions: txnbuild.Preconditions{TimeBounds: tb},
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("error creating transaction: %v", err)
	}

	// Get the transaction in XDR format
	txeB64, err := tx.Base64()
	if err != nil {
		return "", "", fmt.Errorf("error encoding transaction: %v", err)
	}

	// Get the transaction hash to sign
	hash, err := tx.Hash(network.PublicNetworkPassphrase)
	if err != nil {
		return "", "", fmt.Errorf("error getting transaction hash: %v", err)
	}

	return txeB64, base64.StdEncoding.EncodeToString(hash[:]), nil
}

// WithdrawTxn submits a signed Stellar transaction for withdrawal
func WithdrawStellarTxn(
	client *horizonclient.Client,
	serializedTxn string,
	signature string,
) (string, error) {
	// Decode the transaction
	var tx xdr.TransactionEnvelope
	err := xdr.SafeUnmarshalBase64(serializedTxn, &tx)
	if err != nil {
		return "", fmt.Errorf("error decoding transaction: %v", err)
	}

	// Decode the signature
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %w", err)
	}

	// Create decorated signature
	decorated := xdr.DecoratedSignature{
		Hint:      [4]byte{0, 0, 0, 0}, // Hint is not used for Ed25519
		Signature: sigBytes,
	}

	// Add the signature to the appropriate envelope based on type
	switch tx.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		tx.V1.Signatures = append(tx.V1.Signatures, decorated)
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		tx.V0.Signatures = append(tx.V0.Signatures, decorated)
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
		tx.FeeBump.Signatures = append(tx.FeeBump.Signatures, decorated)
	default:
		return "", fmt.Errorf("unsupported transaction envelope type: %v", tx.Type)
	}

	// Encode the transaction for submission
	txeB64, err := xdr.MarshalBase64(tx)
	if err != nil {
		return "", fmt.Errorf("error encoding transaction envelope: %v", err)
	}

	// Submit the transaction
	resp, err := client.SubmitTransactionXDR(txeB64)
	if err != nil {
		hzErr, ok := err.(*horizonclient.Error)
		if ok {
			resultCodes, _ := hzErr.ResultCodes()
			return "", fmt.Errorf("transaction submission failed: %w, Result Codes: %v", hzErr.Problem, resultCodes)
		}
		return "", fmt.Errorf("error submitting transaction: %v", err)
	}

	return resp.Hash, nil
}
