// Package cardano provides functions for interacting with the Cardano blockchain.
//
// The functions in this package use the Blockfrost API to interact with the Cardano
// blockchain. The functions are typically used by the sequencer to retrieve,
// submit and check transactions on the Cardano blockchain.
package cardano

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util"
	"github.com/blockfrost/blockfrost-go"
	"github.com/echovl/cardano-go"
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
		fmt.Printf("error decoding transaction: %v\n", err)
		return "", fmt.Errorf("failed to decode transaction: %v", err)
	}

	var tx cardano.Tx
	err = tx.UnmarshalCBOR(txBytes)
	if err != nil {
		fmt.Printf("error unmarshalling transaction: %+v\n", err)
		return "", fmt.Errorf("failed to unmarshal transaction: %v", err)
	}

	// Calculate the transaction hash that needs to be signed
	txHash, err := tx.Hash()
	if err != nil {
		return "", fmt.Errorf("failed to calculate transaction hash: %v", err)
	}

	fmt.Printf("Transaction hash to be signed: %+v\n", hex.EncodeToString(txHash))

	// Decode public key
	pubKeyBytes, err := hex.DecodeString(strings.TrimPrefix(publicKey, "0x"))
	if err != nil {
		return "", fmt.Errorf("failed to decode public key: %v", err)
	}

	txHash, err = tx.Hash()
	if err != nil {
		return "", fmt.Errorf("failed to calculate transaction hash: %v", err)
	}

	fmt.Printf("Transaction hash to be signed: %+v\n", hex.EncodeToString(txHash))

	// Decode signature
	signatureHex = strings.TrimPrefix(signatureHex, "0x")
	sigBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		fmt.Printf("error decoding signature: %v\n", err)
		return "", fmt.Errorf("failed to decode signature: %v", err)
	}

	// Verify that the signature is actually for this transaction's hash
	// using the provided public key
	if !ed25519.Verify(pubKeyBytes, txHash, sigBytes) {
		fmt.Printf("signature verification failed - signature doesn't match transaction hash")
		return "", fmt.Errorf("signature verification failed - signature doesn't match transaction hash")
	}

	// Create the witness set
	tx.WitnessSet = cardano.WitnessSet{
		VKeyWitnessSet: []cardano.VKeyWitness{
			{
				VKey:      pubKeyBytes,
				Signature: sigBytes,
			},
		},
	}

	// Log the transaction for debugging
	fmt.Printf("Submitting transaction: %+v\n", tx)

	// Submit the transaction
	txHash2, err := client.TransactionSubmit(context.Background(), tx.Bytes())
	if err != nil {
		fmt.Printf("error submitting transaction: %v\n", err)
		return "", fmt.Errorf("error submitting transaction: %v", err)
	}

	return txHash2, nil
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

// func ProcessCardanoTransaction(serializedTxn string, publicKeyHex string, signatureHex string) (string, error) {
// 	// Remove 0x prefix if present
// 	serializedTxn = strings.TrimPrefix(serializedTxn, "0x")
// 	publicKeyHex = strings.TrimPrefix(publicKeyHex, "0x")
// 	signatureHex = strings.TrimPrefix(signatureHex, "0x")

// 	// Decode the transaction bytes
// 	txBytes, err := hex.DecodeString(serializedTxn)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to decode transaction: %v", err)
// 	}

// 	// Parse the transaction
// 	// tx, err := cardano.ParseTx(txBytes)
// 	// if err != nil {
// 	// 	return "", fmt.Errorf("failed to parse transaction: %v", err)
// 	// }

// 	sigBytes, err := hex.DecodeString(signatureHex)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to decode signature: %v", err)
// 	}

// 	// Decode public key
// 	pubKeyBytes, err := hex.DecodeString(publicKeyHex)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to decode public key: %v", err)
// 	}

// 	// txBuilder := cardano.NewTxBuilder(&cardano.ProtocolParams{})

// 	// tx2, err := txBuilder.Build()
// 	// if err != nil {
// 	// 	return "", fmt.Errorf("failed to build transaction: %v", err)
// 	// }

// 	var tx cardano.Tx
// 	err = json.Unmarshal(txBytes, &tx)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to unmarshal transaction: %v", err)
// 	}

// 	tx.WitnessSet.VKeyWitnessSet = []cardano.VKeyWitness{
// 		{
// 			VKey:      pubKeyBytes,
// 			Signature: sigBytes,
// 		},
// 	}

// 	// Create a verification key from the public key bytes
// 	// vkey, err := crypto.NewVerificationKey(pubKeyBytes)
// 	// if err != nil {
// 	// 	return "", fmt.Errorf("failed to create verification key: %v", err)
// 	// }

// 	// Decode signature

// 	// Create a signature from bytes
// 	// signature, err := crypto.NewSignature(sigBytes)
// 	// if err != nil {
// 	// 	return "", fmt.Errorf("failed to create signature: %v", err)
// 	// }

// 	// // Create a witness from the verification key and signature
// 	// witness := cardano.NewWitness(vkey, signature)

// 	// // Add witness to transaction
// 	// tx.AddWitness(witness)

// 	// Serialize the transaction
// 	signedTxBytes, err := tx.Bytes()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to serialize signed transaction: %v", err)
// 	}

// 	return hex.EncodeToString(signedTxBytes), nil
// }

func GetAddressPublicKey(address string) ([]byte, error) {
	// Parse the address using cardano-go
	addr, err := cardano.NewAddress(address)
	if err != nil {
		return nil, fmt.Errorf("failed to parse address: %v", err)
	}

	// Get the payment credential
	paymentCred := addr.Payment.Hash()
	if paymentCred == nil {
		return nil, fmt.Errorf("address has no payment credential")
	}

	// The key hash is what we need - it's derived from the public key
	fmt.Printf("Key hash from address: %x\n", paymentCred)
	return paymentCred, nil
}

// Alternative using Blockfrost API
func GetAddressDetailsFromBlockfrost(chainId string, address string) ([]byte, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return nil, fmt.Errorf("error getting chain: %v", err)
	}

	client, err := getClient(chain.ChainUrl)
	if err != nil {
		return nil, err
	}

	// Get address details from Blockfrost
	addrDetails, err := client.Address(context.Background(), address)
	if err != nil {
		return nil, fmt.Errorf("failed to get address details: %v", err)
	}

	fmt.Printf("Address details: %+v\n", addrDetails)
	return nil, nil // You'll need to extract the key hash from the response
}
