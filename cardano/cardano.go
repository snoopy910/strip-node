// Package cardano provides functions for interacting with the Cardano blockchain.
//
// The functions in this package use the Blockfrost API to interact with the Cardano
// blockchain. The functions are typically used by the sequencer to retrieve,
// submit and check transactions on the Cardano blockchain.
package cardano

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util"
	"github.com/blockfrost/blockfrost-go"
	"github.com/fxamacker/cbor/v2"
)

const (
	decimalsPlaces = 6 // Lovelace (ADA) has 6 decimal places
)

var (
	clientMutex sync.Mutex
	clientMap   = make(map[string]blockfrost.APIClient) // map for different network clients
)

// getClient returns a singleton Blockfrost client for the given chain URL and project ID
func getClient(chainUrl string) (blockfrost.APIClient, error) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	if client, exists := clientMap[chainUrl]; exists {
		return client, nil
	}

	projectID := "mainnet2NiLEqB498izHpWxY90otIxo8UxQ9YSC"

	// Determine network from chain URL
	network := blockfrost.CardanoMainNet
	if strings.Contains(strings.ToLower(chainUrl), "preprod") {
		network = blockfrost.CardanoPreProd
		projectID = "preprodQqR2rJIwZJFQmMCmjcTol3HqCWKi9ZKQ"
	}

	// Create new client
	newClient := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{
			ProjectID: projectID,
			Server:    network,
		},
	)
	clientMap[chainUrl] = newClient

	return newClient, nil
}

// getAssetInfo fetches both decimals and token name for a given asset
func getAssetInfo(client blockfrost.APIClient, unit string) (uint, string, error) {
	if unit == "lovelace" {
		return decimalsPlaces, "ADA", nil
	}

	// Get asset details from Blockfrost
	asset, err := client.Asset(context.Background(), unit)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get asset info: %v", err)
	}

	tokenName, err := hex.DecodeString(asset.AssetName)
	if err != nil {
		return 0, "", fmt.Errorf("failed to decode asset name: %v", err)
	}

	if asset.Metadata == nil {
		// Assuming decimals
		return 0, string(tokenName), nil
	}

	return uint(asset.Metadata.Decimals), string(tokenName), nil
}

// GetCardanoTransfers takes the chain ID and transaction hash as input and returns
// a list of Transfer objects representing the transfers associated with the transaction.
func GetCardanoTransfers(chainId string, txHash string) ([]common.Transfer, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain config: %v", err)
	}

	client, err := getClient(chain.ChainUrl)
	if err != nil {
		return nil, err
	}

	// Get transaction UTXOs
	utxos, err := client.TransactionUTXOs(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction UTXOs: %v", err)
	}

	var transfers []common.Transfer

	// Process outputs (transfers)
	for _, output := range utxos.Outputs {
		for _, amount := range output.Amount {
			decimals, tokenName, err := getAssetInfo(client, amount.Unit)
			if err != nil {
				return nil, fmt.Errorf("failed to get asset info: %v", err)
			}

			amountInt, err := strconv.ParseInt(amount.Quantity, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse amount: %v", err)
			}

			formattedAmount := fmt.Sprintf("%f", float64(amountInt)/math.Pow(10, float64(decimals)))

			if amount.Unit == "lovelace" {
				transfers = append(transfers, common.Transfer{
					From:         utxos.Inputs[0].Address,
					To:           output.Address,
					Amount:       formattedAmount,
					Token:        tokenName,
					IsNative:     true,
					TokenAddress: util.ZERO_ADDRESS,
					ScaledAmount: amount.Quantity,
				})
				continue
			}

			transfers = append(transfers, common.Transfer{
				From:         utxos.Inputs[0].Address,
				To:           output.Address,
				Amount:       formattedAmount,
				Token:        tokenName,
				IsNative:     false,
				TokenAddress: amount.Unit,
				ScaledAmount: amount.Quantity,
			})
		}
	}

	return transfers, nil
}

// SendCardanoTransaction submits a signed Cardano transaction to the network.
func SendCardanoTransaction(serializedTxn string, chainId string, keyCurve string, publicKey string, signatureHex string) (string, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", fmt.Errorf("error getting chain: %v", err)
	}

	client, err := getClient(chain.ChainUrl)
	if err != nil {
		return "", err
	}

	// Decode the transaction bytes
	txBytes, err := hex.DecodeString(strings.TrimPrefix(serializedTxn, "0x"))
	if err != nil {
		return "", fmt.Errorf("failed to decode transaction: %v", err)
	}

	// Decode signature and public key
	sigBytes, err := hex.DecodeString(strings.TrimPrefix(signatureHex, "0x"))
	if err != nil {
		return "", fmt.Errorf("failed to decode signature: %v", err)
	}

	pubKeyBytes, err := hex.DecodeString(strings.TrimPrefix(publicKey, "0x"))
	if err != nil {
		return "", fmt.Errorf("failed to decode public key: %v", err)
	}

	// Use a custom CBOR decoder that preserves the exact structure
	decMode, err := cbor.DecOptions{
		TagsMd:            cbor.TagsAllowed,
		ExtraReturnErrors: cbor.ExtraDecErrorUnknownField,
	}.DecMode()
	if err != nil {
		return "", fmt.Errorf("failed to create decoder: %v", err)
	}

	// Decode the transaction to get its structure
	var txData interface{}
	err = decMode.Unmarshal(txBytes, &txData)
	if err != nil {
		return "", fmt.Errorf("failed to decode transaction: %v", err)
	}

	// For Cardano Shelley transactions, the structure is typically:
	// [transaction_body, transaction_witness_set, transaction_metadata, auxiliary_data]
	// We need to ensure we maintain this exact structure

	txArray, ok := txData.([]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected transaction structure: %T", txData)
	}

	// Create a witness with the public key and signature
	vkeyWitness := []interface{}{
		pubKeyBytes,
		sigBytes,
	}

	// Check if we have a witness set (index 1)
	if len(txArray) > 1 {
		witnessSet, ok := txArray[1].(map[interface{}]interface{})
		if !ok {
			// If it's not a map, create a new one
			witnessSet = make(map[interface{}]interface{})
			txArray[1] = witnessSet
		}

		// Check if we have vkey witnesses
		vkeyWitnesses, ok := witnessSet[uint64(0)].([]interface{})
		if !ok {
			// If not, create a new array with our witness
			witnessSet[uint64(0)] = []interface{}{vkeyWitness}
		} else {
			// If we do, append our witness
			witnessSet[uint64(0)] = append(vkeyWitnesses, vkeyWitness)
		}
	} else {
		// If we don't have a witness set, create one
		witnessSet := map[interface{}]interface{}{
			uint64(0): []interface{}{vkeyWitness},
		}
		txArray = append(txArray, witnessSet)
	}

	// Ensure we have all required elements (at least 3 for a valid Shelley transaction)
	for len(txArray) < 3 {
		txArray = append(txArray, nil)
	}

	// Re-encode the transaction
	encMode, err := cbor.EncOptions{
		Sort:   cbor.SortCanonical,
		TagsMd: cbor.TagsAllowed,
	}.EncMode()
	if err != nil {
		return "", fmt.Errorf("failed to create encoder: %v", err)
	}

	signedTxBytes, err := encMode.Marshal(txArray)
	if err != nil {
		return "", fmt.Errorf("failed to marshal signed transaction: %v", err)
	}

	// Submit the transaction
	txHash, err := client.TransactionSubmit(context.Background(), signedTxBytes)
	if err != nil {
		return "", fmt.Errorf("error submitting transaction: %v", err)
	}

	return txHash, nil
}

// CheckCardanoTransactionConfirmed checks whether a Cardano transaction has been confirmed
func CheckCardanoTransactionConfirmed(chainId string, txHash string) (bool, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return false, fmt.Errorf("error getting chain: %v", err)
	}

	client, err := getClient(chain.ChainUrl)
	if err != nil {
		return false, err
	}

	// Get transaction details
	tx, err := client.Transaction(context.Background(), txHash)
	if err != nil {
		if strings.Contains(err.Error(), "StatusCode:404") {
			return false, nil
		}
		return false, fmt.Errorf("error getting transaction: %v", err)
	}

	// Transaction is confirmed if it has a block number
	return tx.Block != "", nil
}
