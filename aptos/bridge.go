package aptos

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	aptosClient "github.com/portto/aptos-go-sdk/client"
	aptosModels "github.com/portto/aptos-go-sdk/models"
	"github.com/the729/lcs"
)

const (
	DEFAULT_GAS_UNIT_PRICE = 100
	DEFAULT_MAX_GAS_AMOUNT = 5000
)

var ctx = context.Background()

// WithdrawAptosNativeGetSignature returns transaction and dataToSign for
// native APT withdrawl operation
// WithdrawAptosNativeGetSignature creates a transaction for withdrawing native APT tokens
// Returns:
// - serializedTxn: Hex-encoded transaction for submission
// - dataToSign: Hex-encoded message that needs to be signed
// - error: Any error encountered during transaction creation
func WithdrawAptosNativeGetSignature(
	rpcURL string,
	account string,
	amount string,
	recipient string,
) (string, string, error) {
	// Get the core framework account address (0x1)
	// This is where the coin module is deployed
	addr0x1, _ := aptosModels.HexToAccountAddress("0x1")

	// Convert recipient address to Aptos account address format
	// Aptos addresses are 32 bytes, hex-encoded with 0x prefix
	accountTo, err := aptosModels.HexToAccountAddress(recipient)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert account address: %v", err)
	}

	// Initialize Aptos client to fetch account information
	client := aptosClient.NewAptosClient(rpcURL)

	// Get sender's account info to retrieve the current sequence number
	// Sequence number prevents transaction replay and must be incremented for each tx
	accountInfo, err := client.GetAccount(ctx, account)
	if err != nil {
		return "", "", fmt.Errorf("failed to get accountInfo: %v", err)
	}

	// Initialize transaction wrapper
	tx := &aptosModels.Transaction{}

	// Create type tag for AptosCoin
	// This identifies that we're transferring the native APT token
	// Located at 0x1::aptos_coin::AptosCoin
	aptosCoinTypeTag := aptosModels.TypeTagStruct{
		Address: addr0x1,
		Module:  "aptos_coin",
		Name:    "AptosCoin",
	}

	// Build the transaction with a fluent interface
	// 1. Set the sender account
	// 2. Set the entry function call payload (0x1::coin::transfer)
	// 3. Set expiration (10 minutes from now)
	// 4. Set gas parameters
	// 5. Set sequence number from account info
	err = tx.SetSender(account).
		SetPayload(aptosModels.EntryFunctionPayload{
			// Use the coin module from core framework
			Module: aptosModels.Module{
				Address: addr0x1,
				Name:    "coin",
			},
			// Call the transfer function
			Function: "transfer",
			// Specify we're transferring AptosCoin
			TypeArguments: []aptosModels.TypeTag{aptosCoinTypeTag},
			// Arguments: recipient address and amount
			Arguments: []interface{}{accountTo, amount},
		}).
		SetExpirationTimestampSecs(uint64(time.Now().Add(10 * time.Minute).Unix())).
		SetGasUnitPrice(DEFAULT_GAS_UNIT_PRICE).
		SetMaxGasAmount(DEFAULT_MAX_GAS_AMOUNT).
		SetSequenceNumber(accountInfo.SequenceNumber).Error()
	if err != nil {
		return "", "", fmt.Errorf("failed to set transaction data: %v", err)
	}

	// Serialize the transaction using BCS (Binary Canonical Serialization)
	// This is required for transaction submission
	bcsBytes, err := lcs.Marshal(tx.RawTransaction)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert bcs format: %v", err)
	}
	// Convert to hex string and remove 0x prefix
	serializedTxn := strings.TrimPrefix(hex.EncodeToString(bcsBytes), "0x")

	// Get the signing message that needs to be signed by the sender
	// This includes the transaction hash and chain-specific data
	msgBytes, err := tx.GetSigningMessage()
	if err != nil {
		return "", "", fmt.Errorf("failed to get message: %v", err)
	}
	// Convert to hex string and remove 0x prefix
	dataToSign := strings.TrimPrefix(hex.EncodeToString(msgBytes), "0x")

	return serializedTxn, dataToSign, nil
}

// WithdrawAptosTokenGetSignature returns transaction and dataToSign for
// custom Aptos token
func WithdrawAptosTokenGetSignature(
	rpcURL string,
	account string,
	amount string,
	recipient string,
	tokenAddr string,
) (string, string, error) {
	addr0x1, _ := aptosModels.HexToAccountAddress("0x1")

	// get token info
	if len(strings.Split(tokenAddr, "::")) != 3 {
		return "", "", fmt.Errorf("invalid token address format")
	}
	ownerAddr, err := aptosModels.HexToAccountAddress(strings.Split(tokenAddr, "::")[0])
	collenctionName := strings.Split(tokenAddr, "::")[1]
	name := strings.Split(tokenAddr, "::")[2]
	if err != nil {
		return "", "", fmt.Errorf("failed to convert token address: %v", err)
	}

	accountTo, err := aptosModels.HexToAccountAddress(recipient)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert account address: %v", err)
	}

	client := aptosClient.NewAptosClient(rpcURL)

	accountInfo, err := client.GetAccount(ctx, account)
	if err != nil {
		return "", "", fmt.Errorf("failed to get accountInfo: %v", err)
	}

	tx := &aptosModels.Transaction{}

	aptosCoinTypeTag := aptosModels.TypeTagStruct{
		Address: ownerAddr,
		Module:  collenctionName,
		Name:    name,
	}

	err = tx.SetSender(account).
		SetPayload(aptosModels.EntryFunctionPayload{
			Module: aptosModels.Module{
				Address: addr0x1,
				Name:    "coin",
			},
			Function:      "transfer",
			TypeArguments: []aptosModels.TypeTag{aptosCoinTypeTag},
			Arguments:     []interface{}{accountTo, amount},
		}).
		SetExpirationTimestampSecs(uint64(time.Now().Add(10 * time.Minute).Unix())).
		SetGasUnitPrice(DEFAULT_GAS_UNIT_PRICE).
		SetMaxGasAmount(DEFAULT_MAX_GAS_AMOUNT).
		SetSequenceNumber(accountInfo.SequenceNumber).Error()
	if err != nil {
		return "", "", fmt.Errorf("failed to set transaction data: %v", err)
	}

	bcsBytes, err := lcs.Marshal(tx.RawTransaction)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert bcs format: %v", err)
	}
	serializedTxn := strings.TrimPrefix(hex.EncodeToString(bcsBytes), "0x")

	msgBytes, err := tx.GetSigningMessage()
	if err != nil {
		return "", "", fmt.Errorf("failed to get message: %v", err)
	}
	dataToSign := strings.TrimPrefix(hex.EncodeToString(msgBytes), "0x")

	return serializedTxn, dataToSign, nil
}

// WithdrawAptosTxn submits transaction to witdraw assets and return
// the txHash as the result
// WithdrawAptosTxn submits a withdrawal transaction to the Aptos network
// It takes:
// - rpcURL: The URL of the Aptos node to submit the transaction to
// - transaction: Hex-encoded serialized transaction data
// - publicKey: Hex-encoded Ed25519 public key (32 bytes)
// - signatureHex: Hex-encoded Ed25519 signature (64 bytes)
// Returns the transaction hash if successful
func WithdrawAptosTxn(
	rpcURL string,
	transaction string,
	publicKey string,
	signatureHex string,
) (string, error) {
	// Initialize Aptos client with the provided RPC URL
	// This client will be used to submit the transaction
	client := aptosClient.NewAptosClient(rpcURL)

	// Initialize transaction containers
	// Transaction is the outer wrapper that includes both raw transaction and auth data
	tx := &aptosModels.Transaction{}
	// RawTransaction contains the unsigned transaction payload
	rawTxn := &aptosModels.RawTransaction{}

	// Remove '0x' prefix if present and decode the hex-encoded transaction
	// The transaction is serialized using Aptos's LCS format
	serializedTxn := strings.TrimPrefix(transaction, "0x")
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

	// Decode the hex-encoded Ed25519 signature (64 bytes)
	// This signature is created by signing the transaction hash
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %v", err)
	}

	// Decode the hex-encoded Ed25519 public key (32 bytes)
	// Remove '0x' prefix if present
	publicKeyBytes, err := hex.DecodeString(strings.TrimPrefix(publicKey, "0x"))
	if err != nil {
		return "", fmt.Errorf("error decoding public key: %v", err)
	}

	// Create a signed transaction by adding the Ed25519 authenticator
	// The authenticator combines:
	// 1. The public key that can verify the signature
	// 2. The signature of the transaction
	signedTx := tx.SetAuthenticator(aptosModels.TransactionAuthenticatorEd25519{
		PublicKey: publicKeyBytes,
		Signature: signature,
	})

	// Submit the signed transaction to the Aptos network
	// The transaction will be validated and included in a block if valid
	response, err := client.SubmitTransaction(context.Background(), signedTx.UserTransaction)
	if err != nil {
		return "", fmt.Errorf("error submitting transaction: %v", err)
	}

	// Return the transaction hash which can be used to track the transaction status
	return response.Hash, nil
}
