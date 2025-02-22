package bitcoin

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func WithdrawBitcoinGetSignature(
	rpcURL string,
	account string,
	amount string,
	recipient string,
) (string, error) {
	// Create a new Bitcoin transaction
	var msgTx wire.MsgTx
	msgTx.Version = wire.TxVersion

	// Parse solver output for transaction details
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse amount: %w", err)
	}

	// Convert amount to satoshis
	amountSatoshis := int64(amountFloat * 100000000)

	// Create transaction output
	addr, err := btcutil.DecodeAddress(recipient, &chaincfg.MainNetParams)
	if err != nil {
		return "", fmt.Errorf("failed to decode recipient address: %w", err)
	}

	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return "", fmt.Errorf("failed to create output script: %w", err)
	}

	txOut := wire.NewTxOut(amountSatoshis, pkScript)
	msgTx.AddTxOut(txOut)

	// Serialize the transaction
	var buf bytes.Buffer
	if err := msgTx.Serialize(&buf); err != nil {
		return "", fmt.Errorf("failed to serialize transaction: %w", err)
	}

	return hex.EncodeToString(buf.Bytes()), nil
}

func WithdrawBitcoinTxn(
	rpcURL string,
	transaction string,
	signature string,
) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Decode and prepare transaction
	txBytes, _ := hex.DecodeString(transaction)
	sigBytes, _ := hex.DecodeString(signature)

	msgTx := wire.NewMsgTx(wire.TxVersion)
	msgTx.Deserialize(bytes.NewReader(txBytes))

	// Create and apply signature script
	builder := txscript.NewScriptBuilder()
	builder.AddData(sigBytes)
	builder.AddData(txBytes)
	signatureScript, _ := builder.Script()

	for i := range msgTx.TxIn {
		msgTx.TxIn[i].SignatureScript = signatureScript
	}

	// Serialize signed transaction
	var signedTxBuffer bytes.Buffer
	msgTx.Serialize(&signedTxBuffer)
	signedTxHex := hex.EncodeToString(signedTxBuffer.Bytes())

	// Prepare and send RPC request
	rpcRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "sendrawtransaction",
		"params":  []interface{}{signedTxHex},
	}

	jsonData, _ := json.Marshal(rpcRequest)
	req, _ := http.NewRequestWithContext(ctx, "POST", rpcURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var rpcResponse struct {
		Result string `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &rpcResponse)

	if rpcResponse.Error != nil {
		return "", fmt.Errorf("RPC error: %v", rpcResponse.Error.Message)
	}

	return rpcResponse.Result, nil
}
