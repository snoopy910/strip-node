package dogecoin

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

)

const (
	DEFAULT_FEE_RATE     = 0.001 // DOGE per kB
	DEFAULT_CONFIRMATIONS = 6
	DOGE_DECIMALS        = 8
	DOGE_TOKEN_SYMBOL    = "DOGE"
	DOGE_ZERO_ADDRESS    = "0x0000000000000000000000000000000000000000"
)

// WithdrawDogeNativeGetSignature returns transaction and dataToSign for
// native DOGE withdrawal operation
func WithdrawDogeNativeGetSignature(
	rpcURL string,
	account string,
	amount string,
	recipient string,
) (string, string, error) {
	// Validate recipient address
	if !ValidateDogeAddress(recipient) {
		return "", "", fmt.Errorf("invalid recipient address: %s", recipient)
	}

	// Convert amount to satoshis
	amountBig, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return "", "", fmt.Errorf("invalid amount format: %s", amount)
	}

	// Create a new transaction
	msgTx := wire.NewMsgTx(wire.TxVersion)

	// Add inputs (will be populated by the node)
	// For now, we add a dummy input that will be replaced
	dummyHash := make([]byte, 32)
	dummyOutpoint := wire.NewOutPoint((*chainhash.Hash)(dummyHash), 0)
	txIn := wire.NewTxIn(dummyOutpoint, nil, nil)
	msgTx.AddTxIn(txIn)

	// Add the output
	script, err := txscript.NewScriptBuilder().
		AddOp(txscript.OP_DUP).
		AddOp(txscript.OP_HASH160).
		AddData([]byte(recipient)).
		AddOp(txscript.OP_EQUALVERIFY).
		AddOp(txscript.OP_CHECKSIG).
		Script()
	if err != nil {
		return "", "", fmt.Errorf("failed to create output script: %v", err)
	}

	txOut := wire.NewTxOut(amountBig.Int64(), script)
	msgTx.AddTxOut(txOut)

	// Serialize the transaction
	var txBuf strings.Builder
	err = msgTx.Serialize(hex.NewEncoder(&txBuf))
	if err != nil {
		return "", "", fmt.Errorf("failed to serialize transaction: %v", err)
	}

	// The data to sign will be the transaction hash
	txHash := msgTx.TxHash()
	dataToSign := hex.EncodeToString(txHash[:])

	return txBuf.String(), dataToSign, nil
}

// WithdrawDogeTxn submits transaction to withdraw assets and returns
// the txHash as the result
func WithdrawDogeTxn(
	rpcURL string,
	transaction string,
	publicKey string,
	signatureHex string,
) (string, error) {
	// Decode the transaction
	txBytes, err := hex.DecodeString(transaction)
	if err != nil {
		return "", fmt.Errorf("failed to decode transaction: %v", err)
	}

	var msgTx wire.MsgTx
	if err := msgTx.Deserialize(strings.NewReader(string(txBytes))); err != nil {
		return "", fmt.Errorf("failed to deserialize transaction: %v", err)
	}

	// Decode the signature
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode signature: %v", err)
	}

	// Create signature script
	builder := txscript.NewScriptBuilder()
	builder.AddData(signature)
	builder.AddData([]byte(publicKey))
	signatureScript, err := builder.Script()
	if err != nil {
		return "", fmt.Errorf("failed to create signature script: %v", err)
	}

	// Apply signature script to all inputs
	for i := range msgTx.TxIn {
		msgTx.TxIn[i].SignatureScript = signatureScript
	}

	// Serialize the signed transaction
	var signedTxBuf strings.Builder
	err = msgTx.Serialize(hex.NewEncoder(&signedTxBuf))
	if err != nil {
		return "", fmt.Errorf("failed to serialize signed transaction: %v", err)
	}

	// Submit the transaction using RPC
	client := NewDogeRPCClient(rpcURL)
	txHash, err := client.SendRawTransaction(signedTxBuf.String())
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return txHash, nil
}

// Helper function to validate Dogecoin address using regex
