package stellar

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon/operations"
	"github.com/stellar/go/xdr"

	"github.com/StripChain/strip-node/common"
)

// CheckTransactionConfirmed checks if a Stellar transaction has been confirmed
func CheckStellarTransactionConfirmed(chainId string, txnHash string) (bool, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return false, err
	}

	// Initialize Horizon client
	client := GetClient(chain.ChainType, chain.ChainUrl)

	// Get transaction details
	tx, err := client.TransactionDetail(txnHash)
	if err != nil {
		hzErr, ok := err.(*horizonclient.Error)
		if ok && hzErr.Response.StatusCode == 404 {
			// Transaction not found yet
			return false, nil
		}
		return false, err
	}

	// Check if transaction is successful
	return tx.Successful, nil
}

// GetStellarTransfers retrieves transfer information from a Stellar transaction.
// It directly uses the Stellar Horizon API to get transaction and operation details.
func GetStellarTransfers(chainId string, txnHash string) ([]common.Transfer, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain: %w", err)
	}

	// Initialize Horizon client
	client := GetClient(chain.ChainType, chain.ChainUrl)

	// Get transaction details
	tx, err := client.TransactionDetail(txnHash)
	if err != nil {
		hzErr, ok := err.(*horizonclient.Error)
		if ok && hzErr.Response.StatusCode == 404 {
			// Transaction not found yet
			return nil, nil
		}
		return nil, fmt.Errorf("error getting transaction: %w", err)
	}

	var transfers []common.Transfer

	// Helper function to create a transfer from a Stellar asset
	createTransfer := func(from, to, amount string, assetType, assetCode, assetIssuer string) common.Transfer {
		if assetType == "native" {
			return common.Transfer{
				From:         from,
				To:           to,
				Amount:       amount,
				Token:        "XLM",
				IsNative:     true,
				TokenAddress: "XLM",
				ScaledAmount: amount, // Stellar amounts are already in decimal format
			}
		}

		tokenAddress := fmt.Sprintf("%s:%s", assetCode, assetIssuer)
		return common.Transfer{
			From:         from,
			To:           to,
			Amount:       amount,
			Token:        assetCode,
			IsNative:     false,
			TokenAddress: tokenAddress,
			ScaledAmount: amount,
		}
	}

	// Get the transaction operations with a reasonable limit
	opReq := horizonclient.OperationRequest{
		ForTransaction: txnHash,
		Limit:          200, // Maximum operations we'll process
	}
	ops, err := client.Operations(opReq)
	if err != nil {
		return nil, fmt.Errorf("error getting operations: %w", err)
	}

	for _, rawOp := range ops.Embedded.Records {
		// Get the source account for this operation
		sourceAccount := ""
		switch op := rawOp.(type) {
		case *operations.Payment:
			sourceAccount = op.From
		case *operations.PathPayment:
			sourceAccount = op.From
		}

		if sourceAccount == "" {
			sourceAccount = tx.Account
		}

		// Handle different types of payment operations
		// We only care about operations that represent actual value transfers
		switch op := rawOp.(type) {
		case *operations.Payment:
			transfers = append(transfers, createTransfer(
				sourceAccount,
				op.To,
				op.Amount,
				op.Asset.Type,
				op.Asset.Code,
				op.Asset.Issuer,
			))



		case *operations.PathPayment:
			transfers = append(transfers, createTransfer(
				sourceAccount,
				op.To,
				op.Amount,
				op.Asset.Type,
				op.Asset.Code,
				op.Asset.Issuer,
			))
		}
	}
	return transfers, nil
}

// SendTransaction submits a signed Stellar transaction to the network
func SendStellarTxn(serializedTxn string, chainId string, keyCurve string, dataToSign string, signature string) (string, error) {
	// Convert hex signature to XDR format
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %w", err)
	}

	// Create XDR signature
	xdrSig := xdr.DecoratedSignature{
		Hint:      [4]byte{}, // Hint will be first 4 bytes of the public key
		Signature: sigBytes,
	}

	// Convert to base64 XDR
	signatureXDR, err := xdr.MarshalBase64(xdrSig)
	if err != nil {
		return "", fmt.Errorf("error encoding XDR signature: %w", err)
	}

	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get chain: %w", err)
	}

	// Initialize Horizon client
	client := GetClient(chain.ChainType, chain.ChainUrl)

	// Decode the serialized transaction
	var envelope xdr.TransactionEnvelope
	err = xdr.SafeUnmarshalBase64(serializedTxn, &envelope)
	if err != nil {
		return "", fmt.Errorf("failed to decode transaction XDR: %w", err)
	}

	// Convert base32 signature to raw bytes
	sigBytes, err = base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(signatureXDR)
	if err != nil {
		return "", fmt.Errorf("failed to decode base32 signature: %w", err)
	}

	// Create a decorated signature
	decoratedSig := xdr.DecoratedSignature{
		Signature: sigBytes,
		Hint:      [4]byte{}, // Hint is not needed for Ed25519 signatures
	}

	// Add the signature to the appropriate envelope based on type
	switch envelope.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		envelope.V1.Signatures = append(envelope.V1.Signatures, decoratedSig)
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		envelope.V0.Signatures = append(envelope.V0.Signatures, decoratedSig)
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
		envelope.FeeBump.Signatures = append(envelope.FeeBump.Signatures, decoratedSig)
	default:
		return "", fmt.Errorf("unsupported transaction envelope type: %v", envelope.Type)
	}

	// Convert back to base64 for submission
	txeB64, err := xdr.MarshalBase64(envelope)
	if err != nil {
		return "", fmt.Errorf("failed to add signature to transaction: %w", err)
	}

	// Submit the transaction using the base64-encoded XDR
	resp, err := client.SubmitTransactionXDR(txeB64)
	if err != nil {
		hzErr, ok := err.(*horizonclient.Error)
		if ok {
			return "", fmt.Errorf("transaction submission failed: %v (result codes: %v)",
				hzErr.Problem.Title,
				hzErr.Problem.Extras["result_codes"])
		}
		return "", fmt.Errorf("failed to submit transaction: %w", err)
	}

	return resp.Hash, nil
}
