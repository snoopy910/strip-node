// Package aptos provides functions for interacting with the Aptos blockchain.
//
// The functions in this package use the Aptos Go SDK to interact with the Aptos
// blockchain. The functions are typically used by the sequencer to retreive,
// submit and check transactions on the Aptos blockchain.
package aptos

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util"
	aptosClient "github.com/portto/aptos-go-sdk/client"
	aptosModels "github.com/portto/aptos-go-sdk/models"
	"github.com/the729/lcs"
)

// GetAptosTransfers takes the chain ID and the transaction hash as input and returns a list of Transfer objects
// that represent the transfers associated with the transaction.
// GetAptosTransfers retrieves and parses transfer information from an Aptos transaction
// It handles three types of transfers:
// 1. Fungible Asset Transfers (new Aptos token standard)
// 2. Native APT transfers
// 3. Coin transfers (both APT and other tokens)
func GetAptosTransfers(chainId string, txHash string) ([]common.Transfer, error) {
	// Get chain configuration for RPC endpoint and native token info
	chain, err := common.GetChain(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain config: %v", err)
	}

	// Initialize Aptos client to fetch transaction details
	client := aptosClient.NewAptosClient(chain.ChainUrl)

	// Fetch the transaction by its hash
	// This includes the payload with transfer details and arguments
	tx, err := client.GetTransactionByHash(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	var transfers []common.Transfer

	// Switch on the transaction function to handle different transfer types
	switch tx.Payload.Function {

	// Case 1: New Fungible Asset Standard Transfer
	// This is Aptos's newer token standard with enhanced features
	case FUNGIBLE_ASSET_TRANSFER:
		// Extract token address from the nested map structure
		// The 'inner' field contains the actual token address
		address := tx.Payload.Arguments[0].(map[string]interface{})["inner"].(string)

		// Fetch token metadata (symbol, decimals) from the contract
		metadata, err := getMetadata(chain, address)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch token metadata: %v", err)
		}

		// Get transfer amount (3rd argument in the payload)
		amount := tx.Payload.Arguments[2].(string)

		// Format amount according to token decimals
		formattedAmount, err := getFormattedAmount(amount, metadata.Decimal)
		if err != nil {
			return nil, fmt.Errorf("failed to format amount: %v", err)
		}

		// Create transfer record with token details
		transfers = append(transfers, common.Transfer{
			From:         tx.Sender,
			To:           tx.Payload.Arguments[1].(string),
			Amount:       formattedAmount,
			Token:        metadata.Symbol,
			IsNative:     false,
			TokenAddress: address,
			ScaledAmount: amount,
		})

	// Case 2: Native APT Transfer using account transfer
	// This is a direct transfer of native APT tokens
	case ACCOUNT_APT_TRANSFER:
		// Get transfer amount (2nd argument)
		amount := tx.Payload.Arguments[1].(string)
		// APT has 8 decimal places
		formattedAmount, err := getFormattedAmount(amount, 8)
		if err != nil {
			return nil, fmt.Errorf("failed to format amount: %v", err)
		}

		// Create transfer record for native APT
		transfers = append(transfers, common.Transfer{
			From:         tx.Sender,
			To:           tx.Payload.Arguments[0].(string),
			Amount:       formattedAmount,
			Token:        chain.TokenSymbol,
			IsNative:     true,
			TokenAddress: util.ZERO_ADDRESS,
			ScaledAmount: amount,
		})

	// Case 3: Generic Coin Transfers
	// This handles both APT and other token transfers using the coin module
	case ACCOUNT_COIN_TRANSFER, COIN_TRANSFER:
		// Get the type arguments which specify which coins are being transferred
		assetType := tx.Payload.TypeArguments

		// Process each asset type in the transaction
		for _, asset := range assetType {
			// Check if it's a native APT transfer
			if asset == APT_COIN_TYPE {
				// Handle native APT transfer
				amount := tx.Payload.Arguments[1].(string)
				formattedAmount, err := getFormattedAmount(amount, 8)
				if err != nil {
					return nil, fmt.Errorf("failed to format amount: %v", err)
				}

				transfers = append(transfers, common.Transfer{
					From:         tx.Sender,
					To:           tx.Payload.Arguments[0].(string),
					Amount:       formattedAmount,
					Token:        chain.TokenSymbol,
					IsNative:     true,
					TokenAddress: util.ZERO_ADDRESS,
					ScaledAmount: amount,
				})

			} else {
				// Handle other token transfers
				// Extract token address from the type argument (format: address::module::type)
				tokenAddress := strings.Split(asset, "::")[0]

				// Fetch token metadata
				metadata, err := getMetadata(chain, tokenAddress)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch token metadata: %v", err)
				}

				amount := tx.Payload.Arguments[1].(string)

				// Format amount according to token decimals
				formattedAmount, err := getFormattedAmount(amount, metadata.Decimal)
				if err != nil {
					return nil, fmt.Errorf("failed to format amount: %v", err)
				}

				// Create transfer record for the token
				transfers = append(transfers, common.Transfer{
					From:         tx.Sender,
					To:           tx.Payload.Arguments[0].(string),
					Amount:       formattedAmount,
					Token:        metadata.Symbol,
					IsNative:     false,
					TokenAddress: tokenAddress,
					ScaledAmount: amount,
				})
			}
		}
	}

	return transfers, nil
}

// SendAptosTransaction submits an Aptos transaction to the network with the serialized transaction, the chain ID,
// the curve type, the public key and signature as input and returns the transaction hash as a string.
// SendAptosTransaction submits a signed transaction to the Aptos blockchain
// It takes a hex-encoded serialized transaction, chain ID, key curve type, hex-encoded public key,
// and hex-encoded signature as input, and returns the transaction hash
func SendAptosTransaction(serializedTxn string, chainId string, keyCurve string, publicKey string, signatureHex string) (string, error) {
	// Get chain configuration to determine the RPC endpoint
	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", fmt.Errorf("error getting chain: %v", err)
	}

	// Initialize Aptos client with the chain's RPC URL
	client := aptosClient.NewAptosClient(chain.ChainUrl)

	// Initialize transaction containers
	// Transaction is the outer wrapper that includes both raw transaction and auth data
	tx := &aptosModels.Transaction{}
	// RawTransaction contains the unsigned transaction payload
	rawTxn := &aptosModels.RawTransaction{}

	// Remove '0x' prefix if present and decode hex-encoded transaction
	serializedTxn = strings.TrimPrefix(serializedTxn, "0x")
	decodedTransactionData, err := hex.DecodeString(serializedTxn)
	if err != nil {
		return "", fmt.Errorf("error decoding transaction data: %v", err)
	}

	// Deserialize the transaction using Aptos's LCS (Libra Canonical Serialization)
	// This reconstructs the RawTransaction with sender, payload, gas parameters, etc.
	err = lcs.Unmarshal(decodedTransactionData, rawTxn)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshall raw transaction: %v", err)
	}

	// Assign the deserialized raw transaction to our transaction wrapper
	tx.RawTransaction = *rawTxn

	// Decode the hex-encoded Ed25519 signature
	// Aptos uses 64-byte Ed25519 signatures
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %v", err)
	}

	// Decode the hex-encoded Ed25519 public key (32 bytes)
	// and create the transaction authenticator
	publicKeyBytes, err := hex.DecodeString(strings.TrimPrefix(publicKey, "0x"))
	if err != nil {
		return "", fmt.Errorf("error decoding public key: %v", err)
	}

	// Create a signed transaction by adding the Ed25519 authenticator
	// This combines the public key and signature into Aptos's authentication structure
	signedTx := tx.SetAuthenticator(aptosModels.TransactionAuthenticatorEd25519{
		PublicKey: publicKeyBytes,
		Signature: signature,
	})

	// Submit the signed transaction to the Aptos network
	// The response includes the transaction hash for tracking
	response, err := client.SubmitTransaction(context.Background(), signedTx.UserTransaction)
	if err != nil {
		return "", fmt.Errorf("error submitting transaction: %v", err)
	}

	// Return the transaction hash which can be used to track the transaction status
	return response.Hash, nil
}

// CheckAptosTransactionConfirmed checks whether an Aptos transaction has been confirmed
func CheckAptosTransactionConfirmed(chainId string, txnHash string) (bool, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return false, fmt.Errorf("error getting chain: %v", err)
	}

	client := aptosClient.NewAptosClient(chain.ChainUrl)

	tx, err := client.GetTransactionByHash(context.Background(), txnHash)
	if err != nil {
		return false, fmt.Errorf("error getting transaction by hash: %v", err)
	}

	if tx.Success {
		return true, nil
	} else {
		return false, nil
	}
}
