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
func GetAptosTransfers(chainId string, txHash string) ([]common.Transfer, error) {
	// Get chain configuration
	chain, err := common.GetChain(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain config: %v", err)
	}

	// Intialize Aptos client
	client := aptosClient.NewAptosClient(chain.ChainUrl)

	tx, err := client.GetTransactionByHash(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	var transfers []common.Transfer

	switch tx.Payload.Function {

	case FUNGIBLE_ASSET_TRANSFER:

		address := tx.Payload.Arguments[0].(map[string]interface{})["inner"].(string)

		metadata, err := getMetadata(chain, address)
		if err != nil {
			fmt.Println("Error fetching token metadata, ", err)
		}

		amount := tx.Payload.Arguments[2].(string)

		formattedAmount, err := getFormattedAmount(amount, metadata.Decimal)
		if err != nil {
			fmt.Println("Error formatting amount, %w", err)
		}

		transfers = append(transfers, common.Transfer{
			From:         tx.Sender,
			To:           tx.Payload.Arguments[1].(string),
			Amount:       formattedAmount,
			Token:        metadata.Symbol,
			IsNative:     false,
			TokenAddress: address,
			ScaledAmount: amount,
		})

	case ACCOUNT_APT_TRANSFER:

		amount := tx.Payload.Arguments[1].(string)
		formattedAmount, err := getFormattedAmount(amount, 8)
		if err != nil {
			fmt.Println("Error formatting amount, ", err)
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

	case ACCOUNT_COIN_TRANSFER, COIN_TRANSFER:

		assetType := tx.Payload.TypeArguments

		for _, asset := range assetType {

			if asset == APT_COIN_TYPE {

				amount := tx.Payload.Arguments[1].(string)
				formattedAmount, err := getFormattedAmount(amount, 8)
				if err != nil {
					fmt.Println("Error formatting amount, ", err)
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
				tokenAddress := strings.Split(asset, "::")[0]

				metadata, err := getMetadata(chain, tokenAddress)
				if err != nil {
					fmt.Println("Error fetching token metadata, ", err)
				}

				amount := tx.Payload.Arguments[1].(string)

				formattedAmount, err := getFormattedAmount(amount, metadata.Decimal)
				if err != nil {
					fmt.Println("Error formatting amount, %w", err)
				}

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
func SendAptosTransaction(serializedTxn string, chainId string, keyCurve string, publicKey string, signatureHex string) (string, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", fmt.Errorf("error getting chain: %v", err)
	}

	client := aptosClient.NewAptosClient(chain.ChainUrl)

	// Construct the transaction from seralizedTxn
	tx := &aptosModels.Transaction{}

	rawTxn := &aptosModels.RawTransaction{}

	serializedTxn = strings.TrimPrefix(serializedTxn, "0x")
	decodedTransactionData, err := hex.DecodeString(serializedTxn)
	if err != nil {
		return "", fmt.Errorf("error decoding transaction data: %v", err)
	}

	err = lcs.Unmarshal(decodedTransactionData, rawTxn)
	if err != nil {
		fmt.Println("error unmarshalling raw transaction: ", err)
	}

	tx.RawTransaction = *rawTxn

	// Retreive signatureHex
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %v", err)
	}

	// Sign transaction with pubKey and signature
	publicKeyBytes, err := hex.DecodeString(strings.TrimPrefix(publicKey, "0x"))
	if err != nil {
		return "", fmt.Errorf("error decoding public key: %v", err)
	}
	signedTx := tx.SetAuthenticator(aptosModels.TransactionAuthenticatorEd25519{
		PublicKey: publicKeyBytes,
		Signature: signature,
	})

	// Submit transaction
	response, err := client.SubmitTransaction(context.Background(), signedTx.UserTransaction)
	if err != nil {
		return "", fmt.Errorf("error submitting transaction: %v", err)
	}

	fmt.Println("Submitted aptos transaction with hash:", response.Hash)

	return response.Hash, nil
}

// CheckAptosTransactionConfirmed checks whether an Aptos transaction has been confirmed
func CheckAptosTransactionConfirmed(chainId string, txnHash string) (bool, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return false, err
	}

	client := aptosClient.NewAptosClient(chain.ChainUrl)

	tx, err := client.GetTransactionByHash(context.Background(), txnHash)
	if err != nil {
		return false, err
	}

	if tx.Success {
		return true, nil
	} else {
		return false, nil
	}
}
