package dogecoin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/StripChain/strip-node/common"
)

// DogeRPCClient represents a Dogecoin RPC client
type DogeRPCClient struct {
	endpoint string
	client   *http.Client
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
func NewDogeRPCClient(endpoint string) *DogeRPCClient {
	return &DogeRPCClient{
		endpoint: endpoint,
		client:   &http.Client{},
	}
}

// CheckDogeTransactionConfirmed checks if a Dogecoin transaction is confirmed
func CheckDogeTransactionConfirmed(chainId string, txHash string) (bool, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return false, err
	}

	client := NewDogeRPCClient(chain.ChainUrl)
	tx, err := client.GetTransaction(txHash)
	if err != nil {
		return false, fmt.Errorf("failed to get transaction: %v", err)
	}

	return tx.Confirmations >= DEFAULT_CONFIRMATIONS, nil
}

// SendDogeTransaction sends a signed Dogecoin transaction
func SendDogeTransaction(serializedTx string, chainId string, keyCurve string, dataToSign string, signatureHex string) (string, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", err
	}

	client := NewDogeRPCClient(chain.ChainUrl)
	return client.SendRawTransaction(serializedTx)
}

// GetDogeTransfers gets transfers from a Dogecoin transaction
func GetDogeTransfers(chainId string, txHash string) ([]common.Transfer, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return nil, err
	}

	client := NewDogeRPCClient(chain.ChainUrl)
	tx, err := client.GetTransaction(txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	var transfers []common.Transfer

	// Process inputs (from addresses)
	inputMap := make(map[string]float64)
	for _, input := range tx.Vin {
		inputTx, err := client.GetTransaction(input.Txid)
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
			From:   addr,
			To:     "",
			Amount: fmt.Sprintf("%d", int64(value*1e8)), // Convert to satoshis
			Token:  DOGE_TOKEN_SYMBOL,
		})
	}

	// Process outputs (to addresses)
	for _, output := range tx.Vout {
		if len(output.ScriptPubKey.Addresses) > 0 {
			transfers = append(transfers, common.Transfer{
				From:   "",
				To:     output.ScriptPubKey.Addresses[0],
				Amount: fmt.Sprintf("%d", int64(output.Value*1e8)), // Convert to satoshis
				Token:  DOGE_TOKEN_SYMBOL,
			})
		}
	}

	return transfers, nil
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

	req.Header.Set("Content-Type", "application/json")

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

// ValidateDogeAddress validates a Dogecoin address format using regex
func ValidateDogeAddress(address string) bool {
	// Dogecoin address patterns
	// D: Standard address (mainnet)
	// A: Multi-signature address (mainnet)
	// 9: Testnet address
	patterns := []string{
		"^[D][a-km-zA-HJ-NP-Z1-9]{33}$", // Standard mainnet
		"^[A][a-km-zA-HJ-NP-Z1-9]{33}$", // Multisig mainnet
		"^[9][a-km-zA-HJ-NP-Z1-9]{33}$", // Testnet
	}

	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, address)
		if err == nil && matched {
			return true
		}
	}

	return false
}
