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

	"github.com/StripChain/strip-node/common"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func WithdrawBitcoinGetSignature(
	chainId string,
	account string,
	amount string,
	recipient string,
) (string, error) {
	// Create a new Bitcoin transaction
	var msgTx wire.MsgTx
	msgTx.Version = wire.TxVersion

	// Parse amount
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse amount: %w", err)
	}

	// Convert amount to satoshis
	amountSatoshis := int64(amountFloat * 100000000)

	// Get chain parameters based on chainId
	chainParams, err := GetChainParams(chainId)
	if err != nil {
		return "", err
	}

	// Create transaction output
	addr, err := btcutil.DecodeAddress(recipient, chainParams)
	if err != nil {
		return "", fmt.Errorf("failed to decode recipient address: %w", err)
	}

	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return "", fmt.Errorf("failed to create output script: %w", err)
	}

	// Add the main transaction output
	txOut := wire.NewTxOut(amountSatoshis, pkScript)
	msgTx.AddTxOut(txOut)

	// Add a dummy input (will be updated with actual UTXO later)
	dummyHash := chainhash.Hash{}
	dummyOutpoint := wire.NewOutPoint(&dummyHash, 0)
	txIn := wire.NewTxIn(dummyOutpoint, nil, nil)
	msgTx.AddTxIn(txIn)

	// Create P2WPKH script for the input
	_, err = btcutil.DecodeAddress(account, chainParams)
	if err != nil {
		return "", fmt.Errorf("failed to decode from address: %w", err)
	}

	// For P2WPKH, we use empty SignatureScript and put the actual script in witness
	txIn.SignatureScript = []byte{}

	// Serialize the transaction
	var buf bytes.Buffer
	if err := msgTx.Serialize(&buf); err != nil {
		return "", fmt.Errorf("failed to serialize transaction: %w", err)
	}

	return hex.EncodeToString(buf.Bytes()), nil
}

func WithdrawBitcoinTxn(
	chainId string,
	transaction string,
	signature string,
) (string, error) {
	ctx := context.Background()

	// Get chain URL from chainId
	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", fmt.Errorf("failed to get chain: %w", err)
	}

	// Decode transaction and signature
	txBytes, err := hex.DecodeString(transaction)
	if err != nil {
		return "", fmt.Errorf("failed to decode transaction: %w", err)
	}

	// Create DER signature
	derSignatureHex, err := derEncode(signature)
	if err != nil {
		return "", fmt.Errorf("failed to encode signature: %w", err)
	}
	derSignature, err := hex.DecodeString(derSignatureHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode DER signature: %w", err)
	}

	// Deserialize transaction
	msgTx := wire.NewMsgTx(wire.TxVersion)
	if err := msgTx.Deserialize(bytes.NewReader(txBytes)); err != nil {
		return "", fmt.Errorf("failed to deserialize transaction: %w", err)
	}

	// For SegWit transactions, we use witness data instead of signature script
	if len(msgTx.TxIn) == 0 {
		return "", fmt.Errorf("transaction has no inputs")
	}

	// Add SIGHASH_ALL byte to signature
	derSignature = append(derSignature, byte(txscript.SigHashAll))

	// Set witness data for the input
	// For P2WPKH, the witness stack must contain exactly two items:
	// 1. The signature (including sighash type byte)
	// 2. The public key
	witness := wire.TxWitness{derSignature}
	msgTx.TxIn[0].Witness = witness

	// Serialize signed transaction
	var signedTxBuffer bytes.Buffer
	if err := msgTx.Serialize(&signedTxBuffer); err != nil {
		return "", fmt.Errorf("failed to serialize signed transaction: %w", err)
	}
	signedTxHex := hex.EncodeToString(signedTxBuffer.Bytes())

	// Prepare and send RPC request
	rpcRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "sendrawtransaction",
		"params":  []interface{}{signedTxHex},
	}

	jsonData, err := json.Marshal(rpcRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal RPC request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", chain.ChainUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if chain.RpcUsername != "" {
		req.SetBasicAuth(chain.RpcUsername, chain.RpcPassword)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var rpcResponse struct {
		Result string `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &rpcResponse); err != nil {
		return "", fmt.Errorf("failed to parse RPC response: %w", err)
	}

	if rpcResponse.Error != nil {
		return "", fmt.Errorf("RPC error: %v", rpcResponse.Error.Message)
	}

	return rpcResponse.Result, nil
}
