package stellar

import (
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
		sourceAccount := tx.Account
		switch op := rawOp.(type) {
		case operations.Payment:
			if op.From != "" {
				sourceAccount = op.From
			}
			transfers = append(transfers, createTransfer(
				sourceAccount,
				op.To,
				op.Amount,
				op.Asset.Type,
				op.Asset.Code,
				op.Asset.Issuer,
			))
		case operations.PathPayment:
			if op.From != "" {
				sourceAccount = op.From
			}
			transfers = append(transfers, createTransfer(
				sourceAccount,
				op.To,
				op.Amount,
				op.Asset.Type,
				op.Asset.Code,
				op.Asset.Issuer,
			))
		default:
			fmt.Printf("unknown operation type: %T\n", op)
		}
	}
	return transfers, nil
}

// SendTransaction submits a signed Stellar transaction to the network
func SendStellarTxn(serializedTxn string, chainId string, keyCurve string, dataToSign string, signature string) (string, error) {
	// Decode the transaction
	var envelope xdr.TransactionEnvelope
	err := xdr.SafeUnmarshalBase64(serializedTxn, &envelope)
	if err != nil {
		return "", fmt.Errorf("failed to decode transaction XDR: %w", err)
	}

	// Convert hex signature to XDR format
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %w", err)
	}

	// Get the last 4 bytes of the public key for the hint
	publicKey := envelope.SourceAccount().ToAccountId().Ed25519[:]
	hint := [4]byte{publicKey[28], publicKey[29], publicKey[30], publicKey[31]}

	// Add signature
	xdrSig := xdr.DecoratedSignature{
		Hint:      hint,
		Signature: sigBytes,
	}

	// Add the signature to the appropriate envelope based on type
	switch envelope.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		envelope.V1.Signatures = []xdr.DecoratedSignature{xdrSig}
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		envelope.V0.Signatures = []xdr.DecoratedSignature{xdrSig}
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
		envelope.FeeBump.Signatures = []xdr.DecoratedSignature{xdrSig}
	default:
		return "", fmt.Errorf("unsupported transaction envelope type: %v", envelope.Type)
	}

	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get chain: %w", err)
	}

	// Initialize Horizon client
	client := GetClient(chain.ChainId, chain.ChainUrl)

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
