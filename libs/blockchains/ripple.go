package blockchains

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util"
	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/websockets"
)

// NewRippleBlockchain creates a new Stellar blockchain instance
func NewRippleBlockchain(networkType NetworkType) (IBlockchain, error) {
	network := Network{
		networkType: networkType,
		nodeURL:     "wss://s1.ripple.com:51233",
		networkID:   "mainnet",
	}

	if networkType == Testnet {
		network.nodeURL = "wss://s.altnet.rippletest.net:51233"
		network.networkID = "testnet"
	}

	client, err := websockets.NewRemote(network.nodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create ripple client: %v", err)
	}

	return &RippleBlockchain{
		BaseBlockchain: BaseBlockchain{
			chainName:       Ripple,
			network:         network,
			keyCurve:        common.CurveEddsa,
			signingEncoding: "base64",
			decimals:        6,
			opTimeout:       time.Second * 10,
		},
		client: client,
	}, nil
}

// This is a type assertion to ensure that the RippleBlockchain implements the IBlockchain interface
var _ IBlockchain = &RippleBlockchain{}

// RippleBlockchain implements the IBlockchain interface for Stellar
type RippleBlockchain struct {
	BaseBlockchain
	client *websockets.Remote
}

func (b *RippleBlockchain) BroadcastTransaction(txn string, signatureHex string, publicKey *string) (string, error) {
	txBytes, err := hex.DecodeString(strings.TrimPrefix(txn, "0x"))
	if err != nil {
		return "", fmt.Errorf("error decoding transaction: %v", err)
	}

	// Parse the transaction
	var tx data.Payment
	err = json.Unmarshal(txBytes, &tx)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling transaction: %v", err)
	}

	// Add the signature
	sig, err := hex.DecodeString(strings.TrimPrefix(signatureHex, "0x"))
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %v", err)
	}

	// Convert public key to bytes and create data.PublicKey
	publicKeyBytes, err := hex.DecodeString(strings.TrimPrefix(*publicKey, "0x"))
	if err != nil {
		return "", fmt.Errorf("error decoding public key: %v", err)
	}

	var ripplePublicKey data.PublicKey
	copy(ripplePublicKey[:], publicKeyBytes)

	// Set the signature
	sigVar := data.VariableLength(sig)
	tx.SigningPubKey = &ripplePublicKey
	tx.TxnSignature = &sigVar

	// Submit the transaction
	response, err := b.client.Submit(&tx)
	if err != nil {
		return "", fmt.Errorf("error submitting transaction: %v", err)
	}

	txMap := response.Tx.(map[string]interface{})
	txHash := txMap["hash"].(string)
	return txHash, nil
}

func (b *RippleBlockchain) GetTransfers(txHash string, address *string) ([]common.Transfer, error) {
	hash, err := data.NewHash256(strings.TrimPrefix(txHash, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction hash: %v", err)
	}
	tx, err := b.client.Tx(*hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	// Handle different transaction types
	switch txn := tx.Transaction.(type) {
	case *data.Payment:
		tokenAddress := txn.Amount.Issuer.String()
		if txn.Amount.IsNative() {
			tokenAddress = util.ZERO_ADDRESS
		}
		return []common.Transfer{
			{
				From:         txn.Account.String(),
				To:           txn.Destination.String(),
				Amount:       txn.Amount.Value.String(),
				Token:        txn.Amount.Currency.String(),
				IsNative:     txn.Amount.IsNative(),
				TokenAddress: tokenAddress,
				ScaledAmount: fmt.Sprintf("%d", int64(txn.Amount.Value.Float()*math.Pow(10, float64(b.Decimals())))),
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported transaction type: %T", txn)
	}
}

func (b *RippleBlockchain) IsTransactionBroadcastedAndConfirmed(txHash string) (bool, error) {
	hash, err := data.NewHash256(strings.TrimPrefix(txHash, "0x"))
	if err != nil {
		return false, fmt.Errorf("failed to parse transaction hash: %v", err)
	}

	tx, err := b.client.Tx(*hash)
	if err != nil {
		if strings.Contains(err.Error(), "txnNotFound") {
			return false, nil
		}
		return false, fmt.Errorf("error getting transaction: %v", err)
	}

	// Check if transaction is validated
	return tx.MetaData.TransactionResult.Success(), nil
}

func (b *RippleBlockchain) BuildWithdrawTx(bridgeAddress string,
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

		amount2, err := data.NewAmount(amount)
		if err != nil {
			return "", "", fmt.Errorf("failed to parse amount: %v", err)
		}

		bridgeAccount, err := data.NewAccountFromAddress(bridgeAddress)
		if err != nil {
			return "", "", fmt.Errorf("failed to parse bridge address: %v", err)
		}

		userAccount, err := data.NewAccountFromAddress(userAddress)
		if err != nil {
			return "", "", fmt.Errorf("failed to parse user address: %v", err)
		}

		// Get current network fee
		fee, err := b.client.Fee()
		if err != nil {
			return "", "", fmt.Errorf("failed to get network fee: %v", err)
		}

		// Create payment transaction
		payment := &data.Payment{
			TxBase: data.TxBase{
				TransactionType: data.PAYMENT,
				Account:         *bridgeAccount,
				Fee:             fee.Drops.BaseFee,
			},
			Destination: *userAccount,
			Amount:      *amount2,
		}

		// Get account sequence
		account, err := b.client.AccountInfo(*bridgeAccount)
		if err != nil {
			return "", "", fmt.Errorf("failed to get account info: %v", err)
		}

		payment.Sequence = *account.AccountData.Sequence

		// Serialize transaction
		tx := data.NewTransactionWithMetadata(data.PAYMENT)
		tx.Transaction = payment
		txBytes, err := tx.MarshalJSON()
		if err != nil {
			return "", "", fmt.Errorf("failed to marshal transaction: %v", err)
		}

		hash, _, err := data.SigningHash(payment)
		if err != nil {
			return "", "", fmt.Errorf("failed to get signing hash: %v", err)
		}

		serializedTxn := hex.EncodeToString(txBytes)
		dataToSign := hex.EncodeToString(hash.Bytes())

		return serializedTxn, dataToSign, nil
	}

	var solverData map[string]interface{}
	if err := json.Unmarshal([]byte(solverOutput), &solverData); err != nil {
		return "", "", fmt.Errorf("failed to parse solver output: %v", err)
	}

	amount, ok := solverData["amount"].(string)
	if !ok {
		return "", "", fmt.Errorf("amount not found in solver output")
	}

	bridgeAccount, err := data.NewAccountFromAddress(bridgeAddress)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse bridge address: %v", err)
	}

	userAccount, err := data.NewAccountFromAddress(userAddress)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse user address: %v", err)
	}

	// tokenCurrency, err := data.NewCurrency(tokenCode)
	// if err != nil {
	// 	return "", "", fmt.Errorf("failed to parse token code: %v", err)
	// }

	tokenIssuerAccount, err := data.NewAccountFromAddress(*tokenAddress)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse token issuer: %v", err)
	}

	value, err := data.NewValue(amount, true)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse amount: %v", err)
	}
	// Create token amount
	tokenAmount := data.Amount{
		Value: value,
		// Currency: tokenCurrency,
		Issuer: *tokenIssuerAccount,
	}

	// Get current network fee
	fee, err := b.client.Fee()
	if err != nil {
		return "", "", fmt.Errorf("failed to get network fee: %v", err)
	}

	// Create payment transaction
	payment := &data.Payment{
		TxBase: data.TxBase{
			TransactionType: data.PAYMENT,
			Account:         *bridgeAccount,
			Fee:             fee.Drops.BaseFee,
		},
		Destination: *userAccount,
		Amount:      tokenAmount,
	}

	// Get account sequence
	account, err := b.client.AccountInfo(*bridgeAccount)
	if err != nil {
		return "", "", fmt.Errorf("failed to get account info: %v", err)
	}

	payment.Sequence = *account.AccountData.Sequence

	// Serialize transaction
	tx := data.NewTransactionWithMetadata(data.PAYMENT)
	tx.Transaction = payment
	txBytes, err := tx.MarshalJSON()
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal transaction: %v", err)
	}

	hash, _, err := data.SigningHash(payment)
	if err != nil {
		return "", "", fmt.Errorf("failed to get signing hash: %v", err)
	}

	serializedTxn := hex.EncodeToString(txBytes)
	dataToSign := hex.EncodeToString(hash.Bytes())

	return serializedTxn, dataToSign, nil
}

func (b *RippleBlockchain) RawPublicKeyBytesToAddress(pkBytes []byte, networkType NetworkType) (string, error) {
	return "", errors.New("RawPublicKeyBytesToAddress not implemented")
}

func (b *RippleBlockchain) RawPublicKeyToPublicKeyStr(pkBytes []byte) (string, error) {
	return "", errors.New("RawPublicKeyToPublicKeyStr not implemented")
}
