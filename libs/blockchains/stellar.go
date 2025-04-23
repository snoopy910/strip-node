package blockchains

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/libs"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon/operations"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

const (
	MaxStellarOperationsPerTx = 100
)

// NewStellarBlockchain creates a new Stellar blockchain instance
func NewStellarBlockchain(networkType NetworkType) (IBlockchain, error) {
	network := Network{
		networkType: networkType,
		nodeURL:     horizonclient.DefaultPublicNetClient.HorizonURL,
		networkID:   "mainnet",
	}

	if networkType == Testnet {
		network.nodeURL = horizonclient.DefaultTestNetClient.HorizonURL
		network.networkID = "testnet"
	}

	client := &horizonclient.Client{
		HorizonURL: network.nodeURL,
		HTTP:       http.DefaultClient,
	}
	// Set timeout using the SDK's constant
	client.SetHorizonTimeout(horizonclient.HorizonTimeout)

	return &StellarBlockchain{
		BaseBlockchain: BaseBlockchain{
			chainName:       Stellar,
			network:         network,
			keyCurve:        common.CurveEddsa,
			signingEncoding: "base64",
			decimals:        7,
			opTimeout:       time.Second * 10,
		},
		client: client,
	}, nil
}

// This is a type assertion to ensure that the StellarBlockchain implements the IBlockchain interface
var _ IBlockchain = &StellarBlockchain{}

// StellarBlockchain implements the IBlockchain interface for Stellar
type StellarBlockchain struct {
	BaseBlockchain
	client *horizonclient.Client
}

func (b *StellarBlockchain) BroadcastTransaction(txn string, signatureHex string, _ *string) (string, error) {
	var envelope xdr.TransactionEnvelope
	err := xdr.SafeUnmarshalBase64(txn, &envelope)
	if err != nil {
		return "", fmt.Errorf("failed to decode transaction XDR: %w", err)
	}

	// Convert hex signature to XDR format
	sigBytes, err := hex.DecodeString(signatureHex)
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

	// Convert back to base64 for submission
	txeB64, err := xdr.MarshalBase64(envelope)
	if err != nil {
		return "", fmt.Errorf("failed to add signature to transaction: %w", err)
	}

	// Submit the transaction using the base64-encoded XDR
	resp, err := b.client.SubmitTransactionXDR(txeB64)
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

func (b *StellarBlockchain) GetTransfers(txHash string, address *string) ([]common.Transfer, error) {
	tx, err := b.client.TransactionDetail(txHash)
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
		ForTransaction: txHash,
		Limit:          MaxStellarOperationsPerTx,
	}
	ops, err := b.client.Operations(opReq)
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
			logger.Sugar().Errorw("unknown operation type", "type", op)
		}
	}
	return transfers, nil
}

func (b *StellarBlockchain) IsTransactionBroadcastedAndConfirmed(txHash string) (bool, error) {
	tx, err := b.client.TransactionDetail(txHash)
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

func (b *StellarBlockchain) BuildWithdrawTx(bridgeAddress string,
	solverOutput string,
	userAddress string,
	tokenAddress *string,
) (string, string, error) {
	if tokenAddress == nil {
		var solverData map[string]interface{}
		if err := json.Unmarshal([]byte(solverOutput), &solverData); err != nil {
			return "", "", fmt.Errorf("failed to parse solver output: %v", err)
		}

		amount, ok := solverData["amount"].(string)
		if !ok {
			return "", "", fmt.Errorf("amount not found in solver output")
		}

		// Get bridge account to use as source account
		account, err := b.client.AccountDetail(horizonclient.AccountRequest{AccountID: bridgeAddress})
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
	account, err := b.client.AccountDetail(horizonclient.AccountRequest{AccountID: bridgeAddress})
	if err != nil {
		hzErr, ok := err.(*horizonclient.Error)
		if ok {
			return "", "", fmt.Errorf("error getting bridge account: %v", hzErr.Problem)
		}

		return "", "", fmt.Errorf("error getting bridge account: %v", err)
	}

	// Set up time bounds
	tb := txnbuild.NewTimeout(300)

	assetInfo := strings.Split(*tokenAddress, ":")
	assetCode := assetInfo[0]
	assetIssuer := assetInfo[1]
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

func (b *StellarBlockchain) RawPublicKeyBytesToAddress(pkBytes []byte, networkType NetworkType) (string, error) {
	return "", errors.New("RawPublicKeyBytesToAddress not implemented")
}

func (b *StellarBlockchain) RawPublicKeyToPublicKeyStr(pkBytes []byte) (string, error) {
	return "", errors.New("RawPublicKeyToPublicKeyStr not implemented")
}

func (b *StellarBlockchain) ExtractDestinationAddress(operation *libs.Operation) (string, error) {
	// For Stellar, parse the XDR transaction envelope
	var txEnv xdr.TransactionEnvelope
	destAddress := ""
	err := xdr.SafeUnmarshalBase64(*operation.SerializedTxn, &txEnv)
	if err != nil {
		return "", fmt.Errorf("error parsing Stellar transaction", err)
	}

	// Get the first operation's destination
	if len(txEnv.Operations()) > 0 {
		if paymentOp, ok := txEnv.Operations()[0].Body.GetPaymentOp(); ok {
			destAddress = paymentOp.Destination.Address()
		}
	}
	return destAddress, nil
}
