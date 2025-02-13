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

var ctx = context.Background()

// WithdrawAptosNativeGetSignature returns transaction and dataToSign for
// native APT withdrawl operation
func WithdrawAptosNativeGetSignature(
	rpcURL string,
	account string,
	amount string,
	recipient string,
) (string, string, error) {
	addr0x1, _ := aptosModels.HexToAccountAddress("0x1")

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
		Address: addr0x1,
		Module:  "aptos_coin",
		Name:    "AptosCoin",
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
		SetGasUnitPrice(uint64(100)).
		SetMaxGasAmount(uint64(5000)).
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
		SetGasUnitPrice(uint64(100)).
		SetMaxGasAmount(uint64(5000)).
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
func WithdrawAptosTxn(
	rpcURL string,
	transaction string,
	publicKey string,
	signatureHex string,
) (string, error) {
	client := aptosClient.NewAptosClient(rpcURL)

	// Construct the transaction from seralizedTxn
	tx := &aptosModels.Transaction{}

	rawTxn := &aptosModels.RawTransaction{}

	serializedTxn := strings.TrimPrefix(transaction, "0x")
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
