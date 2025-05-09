// Package ripple provides functions for interacting with the XRP Ledger.
//
// The functions in this package use the XRP Ledger API to interact with the XRP
// blockchain. The functions are typically used by the sequencer to retrieve,
// submit and check transactions on the XRP Ledger.
package ripple

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util"
	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/websockets"
)

const decimalsPlaces = 6

var (
	clientMutex sync.Mutex
	clientMap   = make(map[string]*websockets.Remote) // map due to tests
)

// getClient returns a singleton websockets.Remote client for the given chain URL
func getClient(chainUrl string) (*websockets.Remote, error) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	if client, exists := clientMap[chainUrl]; exists {
		return client, nil
	}

	newClient, err := websockets.NewRemote(chainUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create ripple client: %v", err)
	}
	clientMap[chainUrl] = newClient

	return newClient, nil
}

// GetRippleTransfers takes the chain ID and the transaction hash as input and returns
// a list of Transfer objects representing the transfers associated with the transaction.
func GetRippleTransfers(chainId string, txHash string) ([]common.Transfer, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain config: %v", err)
	}

	client, err := getClient(chain.ChainUrl)
	if err != nil {
		return nil, err
	}

	// Get transaction details
	hash, err := data.NewHash256(strings.TrimPrefix(txHash, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction hash: %v", err)
	}
	tx, err := client.Tx(*hash)
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
				ScaledAmount: fmt.Sprintf("%d", int64(txn.Amount.Value.Float()*math.Pow(10, float64(decimalsPlaces)))),
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported transaction type: %T", txn)
	}
}

// SendRippleTransaction submits a signed XRP Ledger transaction to the network.
func SendRippleTransaction(serializedTxn string, chainId string, keyCurve string, publicKey string, signatureHex string) (string, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", fmt.Errorf("error getting chain: %v", err)
	}

	client, err := getClient(chain.ChainUrl)
	if err != nil {
		return "", err
	}

	// Decode the serialized transaction
	txBytes, err := hex.DecodeString(strings.TrimPrefix(serializedTxn, "0x"))
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
	publicKeyBytes, err := hex.DecodeString(strings.TrimPrefix(publicKey, "0x"))
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
	response, err := client.Submit(&tx)
	if err != nil {
		return "", fmt.Errorf("error submitting transaction: %v", err)
	}

	txMap := response.Tx.(map[string]interface{})
	txHash := txMap["hash"].(string)
	return txHash, nil
}

// CheckRippleTransactionConfirmed checks whether a XRP Ledger transaction has been confirmed
func CheckRippleTransactionConfirmed(chainId string, txHash string) (bool, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return false, fmt.Errorf("error getting chain: %v", err)
	}

	client, err := getClient(chain.ChainUrl)
	if err != nil {
		return false, err
	}

	// Get transaction details
	hash, err := data.NewHash256(strings.TrimPrefix(txHash, "0x"))
	if err != nil {
		return false, fmt.Errorf("failed to parse transaction hash: %v", err)
	}

	tx, err := client.Tx(*hash)
	if err != nil {
		if strings.Contains(err.Error(), "txnNotFound") {
			return false, nil
		}
		return false, fmt.Errorf("error getting transaction: %v", err)
	}

	// Check if transaction is validated
	return tx.MetaData.TransactionResult.Success(), nil
}
