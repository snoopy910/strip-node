package blockchains

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/stellar/go/clients/horizonclient"
)

const (
	DOGECOIN_DEFAULT_CONFIRMATIONS = 6
)

// DogeRPCClient represents a Dogecoin RPC client
type DogeRPCClient struct {
	endpoint string
	client   *http.Client
	apiKey   string
}

// RPCRequest represents a JSON-RPC request
type RPCRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// RPCResponse represents a JSON-RPC response
type RPCResponse struct {
	JsonRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error"`
	ID      int             `json:"id"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Transaction represents a Dogecoin transaction
type Transaction struct {
	TxID          string  `json:"txid"`
	Version       int     `json:"version"`
	LockTime      uint32  `json:"locktime"`
	Vin           []TxIn  `json:"vin"`
	Vout          []TxOut `json:"vout"`
	Confirmations int     `json:"confirmations"`
}

// TxIn represents a transaction input
type TxIn struct {
	Txid      string `json:"txid"`
	Vout      uint32 `json:"vout"`
	ScriptSig struct {
		Asm string `json:"asm"`
		Hex string `json:"hex"`
	} `json:"scriptSig"`
	Sequence uint32 `json:"sequence"`
}

// TxOut represents a transaction output
type TxOut struct {
	Value        float64 `json:"value"`
	N            uint32  `json:"n"`
	ScriptPubKey struct {
		Asm       string   `json:"asm"`
		Hex       string   `json:"hex"`
		ReqSigs   int      `json:"reqSigs"`
		Type      string   `json:"type"`
		Addresses []string `json:"addresses"`
	} `json:"scriptPubKey"`
}

// NewDogeRPCClient creates a new Dogecoin RPC client
const (
	defaultTimeout = 15 * time.Second
	maxRetries     = 3
)

func NewDogeRPCClient(endpoint string, apiKey string) *DogeRPCClient {
	client := &http.Client{
		Timeout: defaultTimeout,
		Transport: &http.Transport{
			MaxIdleConns:          10,
			IdleConnTimeout:       defaultTimeout,
			TLSHandshakeTimeout:   defaultTimeout,
			ExpectContinueTimeout: defaultTimeout,
			MaxConnsPerHost:       maxRetries,
		},
	}

	return &DogeRPCClient{
		endpoint: endpoint,
		client:   client,
		apiKey:   apiKey,
	}
}

// NewDogecoinBlockchain creates a new Stellar blockchain instance
func NewDogecoinBlockchain(networkType NetworkType) (IBlockchain, error) {
	apiKey := os.Getenv("TATUM_API_KEY")
	network := Network{
		networkType: networkType,
		nodeURL:     horizonclient.DefaultPublicNetClient.HorizonURL,
		networkID:   "mainnet",
	}

	if networkType == Testnet {
		network.nodeURL = horizonclient.DefaultTestNetClient.HorizonURL
		network.networkID = "testnet"
		if apiKey == "" {
			apiKey = "t-67cb0e957f1a5a5a2483e093-eeb92712de9c4144a0edbcca"
		}
	}

	if apiKey == "" {
		apiKey = "t-67cb0e957f1a5a5a2483e093-03079b3a500a4f39bf4d651b"
	}

	client := NewDogeRPCClient(network.nodeURL, apiKey)

	return &DogecoinBlockchain{
		BaseBlockchain: BaseBlockchain{
			chainName:       Dogecoin,
			network:         network,
			keyCurve:        common.CurveEcdsa,
			signingEncoding: "hex",
			tokenSymbol:     "DOGE",
			decimals:        8,
			opTimeout:       time.Second * 10,
		},
		client: client,
	}, nil
}

// This is a type assertion to ensure that the DogecoinBlockchain implements the IBlockchain interface
var _ IBlockchain = &DogecoinBlockchain{}

// DogecoinBlockchain implements the IBlockchain interface for Stellar
type DogecoinBlockchain struct {
	BaseBlockchain
	client *DogeRPCClient
}

func (b *DogecoinBlockchain) BroadcastTransaction(txn string, signatureHex string, pubkey *string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(*pubkey)
	if err != nil {
		return "", fmt.Errorf("error decoding public key: %v", err)
	}

	// Step 1: Parse the transaction first
	msgTx, err := parseSerializedTransaction(txn)
	if err != nil {
		return "", fmt.Errorf("error parsing transaction: %v", err)
	}

	// Step 2: Create DER signature
	derSignatureHex, err := derEncode(signatureHex)
	if err != nil {
		return "", fmt.Errorf("error encoding signature: %v", err)
	}
	derSignature, err := hex.DecodeString(derSignatureHex)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %v", err)
	}

	// Step 3: Add signature to the transaction
	sigScript, err := txscript.NewScriptBuilder().
		AddData(derSignature).
		AddData(pubKeyBytes).
		Script()
	if err != nil {
		return "", fmt.Errorf("error creating signature script: %v", err)
	}
	// Set the signature script for the input
	for i := range msgTx.TxIn {
		msgTx.TxIn[i].SignatureScript = sigScript
	}

	// Step 4: Serialize the transaction
	var signedTxBuffer bytes.Buffer
	if err := msgTx.Serialize(&signedTxBuffer); err != nil {
		return "", fmt.Errorf("error serializing signed transaction: %v", err)
	}
	signedTxHex := hex.EncodeToString(signedTxBuffer.Bytes())
	fmt.Println("signed serialized txn:", signedTxHex)

	return b.client.SendRawTransaction(signedTxHex)
}

func (b *DogecoinBlockchain) GetTransfers(txHash string, address *string) ([]common.Transfer, error) {
	tx, err := b.client.GetTransaction(txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	var transfers []common.Transfer

	// Process inputs (from addresses)
	inputMap := make(map[string]float64)
	for _, input := range tx.Vin {
		inputTx, err := b.client.GetTransaction(input.Txid)
		if err != nil {
			continue
		}
		if int(input.Vout) < len(inputTx.Vout) {
			vout := inputTx.Vout[input.Vout]
			if len(vout.ScriptPubKey.Addresses) > 0 {
				addr := vout.ScriptPubKey.Addresses[0]
				inputMap[addr] += vout.Value
			}
		}
	}

	// Convert inputs to transfers
	for addr, value := range inputMap {
		transfers = append(transfers, common.Transfer{
			From:         addr,
			To:           "",
			Amount:       fmt.Sprintf("%d", int64(value*1e8)), // Convert to satoshis
			Token:        b.TokenSymbol(),
			IsNative:     true,
			ScaledAmount: fmt.Sprintf("%d", int64(value)),
		})
	}

	// Process outputs (to addresses)
	for _, output := range tx.Vout {
		if len(output.ScriptPubKey.Addresses) > 0 {
			transfers = append(transfers, common.Transfer{
				From:         "",
				To:           output.ScriptPubKey.Addresses[0],
				Amount:       fmt.Sprintf("%d", int64(output.Value*1e8)), // Convert to satoshis
				Token:        b.TokenSymbol(),
				IsNative:     true,
				ScaledAmount: fmt.Sprintf("%d", int64(output.Value)),
			})
		}
	}

	return transfers, nil
}

func (b *DogecoinBlockchain) IsTransactionBroadcastedAndConfirmed(txHash string) (bool, error) {
	tx, err := b.client.GetTransaction(txHash)
	if err != nil {
		return false, fmt.Errorf("failed to get transaction: %v", err)
	}

	return tx.Confirmations >= DOGECOIN_DEFAULT_CONFIRMATIONS, nil
}

func (b *DogecoinBlockchain) BuildWithdrawTx(bridgeAddress string,
	solverOutput string,
	recipient string,
	tokenAddress *string,
) (string, string, error) {
	var solverData map[string]interface{}
	if err := json.Unmarshal([]byte(solverOutput), &solverData); err != nil {
		return "", "", fmt.Errorf("failed to parse solver output: %v", err)
	}

	amount, ok := solverData["amount"].(string)
	if !ok {
		return "", "", fmt.Errorf("amount not found in solver output")
	}

	if tokenAddress == nil {
		// Validate recipient address
		// if !ValidateDogeAddress(recipient) {
		// 	return "", "", fmt.Errorf("invalid recipient address: %s", recipient)
		// }

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
		var txBuf bytes.Buffer
		err = msgTx.Serialize(&txBuf)
		if err != nil {
			return "", "", fmt.Errorf("failed to serialize transaction: %v", err)
		}

		// Create a new transaction for dataToSign
		msgTxToSign := wire.NewMsgTx(wire.TxVersion)
		txIn = wire.NewTxIn(dummyOutpoint, script, nil)
		msgTxToSign.AddTxIn(txIn)
		msgTxToSign.AddTxOut(txOut)

		// Serialize the transaction to sign
		var txBufToSign bytes.Buffer
		err = msgTxToSign.Serialize(&txBufToSign)
		if err != nil {
			return "", "", fmt.Errorf("failed to serialize transaction: %v", err)
		}

		// Hash-256 the transaction to get dataToSign
		dataToSign := sha256.Sum256(txBufToSign.Bytes())

		return hex.EncodeToString(txBuf.Bytes()), hex.EncodeToString(dataToSign[:]), nil
	}
	return "", "", errors.New("BuildWithdrawTx not implemented")
}

func (b *DogecoinBlockchain) RawPublicKeyBytesToAddress(pkBytes []byte, networkType NetworkType) (string, error) {
	return "", errors.New("RawPublicKeyBytesToAddress not implemented")
}

func (b *DogecoinBlockchain) RawPublicKeyToPublicKeyStr(pkBytes []byte) (string, error) {
	return "", errors.New("RawPublicKeyToPublicKeyStr not implemented")
}

// GetTransaction gets a transaction by its hash
func (c *DogeRPCClient) GetTransaction(txHash string) (*Transaction, error) {
	response, err := c.call("getrawtransaction", []interface{}{txHash, true})
	if err != nil {
		return nil, err
	}

	var tx Transaction
	if err := json.Unmarshal(response, &tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %v", err)
	}

	return &tx, nil
}

// SendRawTransaction sends a raw transaction
func (c *DogeRPCClient) SendRawTransaction(txHex string) (string, error) {
	response, err := c.call("sendrawtransaction", []interface{}{txHex})
	if err != nil {
		return "", err
	}

	var txHash string
	if err := json.Unmarshal(response, &txHash); err != nil {
		return "", fmt.Errorf("failed to unmarshal transaction hash: %v", err)
	}
	fmt.Println("Dogecoin tx sent", txHash)

	return txHash, nil
}

// call makes an RPC call to the Dogecoin node
func (c *DogeRPCClient) call(method string, params []interface{}) (json.RawMessage, error) {
	request := RPCRequest{
		JsonRPC: "1.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("x-api-key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var rpcResponse RPCResponse
	if err := json.Unmarshal(body, &rpcResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if rpcResponse.Error != nil {
		return nil, fmt.Errorf("RPC error: %v", rpcResponse.Error.Message)
	}

	return rpcResponse.Result, nil
}

// parseSerializedTransaction parses a base64 encoded serialized raw transaction
// and returns the unsigned transaction as a wire.MsgTx.
func parseSerializedTransaction(serializedTxn string) (*wire.MsgTx, error) {
	txBytes, err := hex.DecodeString(serializedTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex transaction: %v", err)
	}

	return parseRawTransaction(txBytes)
}

// parseRawTransaction parses a raw transaction bytes into a wire.MsgTx.
func parseRawTransaction(txBytes []byte) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(wire.TxVersion)
	err := tx.Deserialize(bytes.NewReader(txBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %v", err)
	}
	return tx, nil
}

func derEncode(signature string) (string, error) {
	// Decode hex signature
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %v", err)
	}

	if len(sigBytes) != 64 {
		return "", fmt.Errorf("invalid signature length: expected 64 bytes, got %d", len(sigBytes))
	}

	// Split into r and s components (32 bytes each)
	r := new(big.Int).SetBytes(sigBytes[:32])
	s := new(big.Int).SetBytes(sigBytes[32:])

	// Get curve parameters
	curve := btcec.S256()
	halfOrder := new(big.Int).Rsh(curve.N, 1)

	// Normalize S value to be in the lower half of the curve
	if s.Cmp(halfOrder) > 0 {
		s = new(big.Int).Sub(curve.N, s)
	}

	// Convert r and s to bytes, removing leading zeros
	rBytes := r.Bytes()
	sBytes := s.Bytes()

	// Add 0x00 prefix if the highest bit is set (to ensure positive number)
	if rBytes[0]&0x80 == 0x80 {
		rBytes = append([]byte{0x00}, rBytes...)
	}
	if sBytes[0]&0x80 == 0x80 {
		sBytes = append([]byte{0x00}, sBytes...)
	}

	// Calculate lengths
	rLen := len(rBytes)
	sLen := len(sBytes)
	totalLen := rLen + sLen + 4 // 4 additional bytes for DER sequence

	// Create DER signature
	derSig := make([]byte, 0, totalLen+1)   // +1 for sighash type
	derSig = append(derSig, 0x30)           // sequence tag
	derSig = append(derSig, byte(totalLen)) // length of sequence

	// Encode R value
	derSig = append(derSig, 0x02) // integer tag
	derSig = append(derSig, byte(rLen))
	derSig = append(derSig, rBytes...)

	// Encode S value
	derSig = append(derSig, 0x02) // integer tag
	derSig = append(derSig, byte(sLen))
	derSig = append(derSig, sBytes...)

	// Add SIGHASH_ALL
	derSig = append(derSig, 0x01)

	return hex.EncodeToString(derSig), nil
}
