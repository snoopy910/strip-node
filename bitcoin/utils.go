package bitcoin

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"regexp"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

type BitcoinNetworkConfig struct {
	Name          string
	NetworkType   string
	AddressPrefix []string
	RPCPort       int
	Params        chaincfg.Params
}

var (
	MainnetConfig = BitcoinNetworkConfig{
		Name:          "mainnet",
		NetworkType:   "mainnet",
		AddressPrefix: []string{"1", "3", "bc1"},
		RPCPort:       8332,
		Params:        chaincfg.MainNetParams,
	}

	TestnetConfig = BitcoinNetworkConfig{
		Name:          "testnet",
		NetworkType:   "testnet3",
		AddressPrefix: []string{"m", "n", "tb1"},
		RPCPort:       18332,
		Params:        chaincfg.TestNet3Params,
	}

	RegtestConfig = BitcoinNetworkConfig{
		Name:          "regtest",
		NetworkType:   "regtest",
		AddressPrefix: []string{"m", "n", "bcrt1"},
		RPCPort:       18443,
		Params:        chaincfg.RegressionNetParams,
	}
)

// Bitcoin integration constants
const (
	BTC_TOKEN_SYMBOL = "BTC"                                        // Symbol representing Bitcoin
	SATOSHI_DECIMALS = 8                                            // Number of decimals in Bitcoin (1 BTC = 10^8 satoshis)
	BTC_ZERO_ADDRESS = "0x0000000000000000000000000000000000000000" // Representing a zero address for Bitcoin (not used)
)

// RPCResponse mimics Bitcoin Core's JSON-RPC response
type RPCResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	ID int `json:"id"`
}

// TxDetails mimics getrawtransaction verbose output
type TxDetails struct {
	TxID          string `json:"txid"`
	Vin           []Vin  `json:"vin"`
	Vout          []Vout `json:"vout"`
	BlockHash     string `json:"blockhash,omitempty"`
	Confirmations int64  `json:"confirmations"`
}
type Vin struct {
	TxID      string `json:"txid"`
	Vout      uint32 `json:"vout"`
	ScriptSig struct {
		Asm string `json:"asm"`
		Hex string `json:"hex"`
	} `json:"scriptSig"`
	Witness []string `json:"txinwitness,omitempty"`
}

type Vout struct {
	Value        float64 `json:"value"` // BTC
	N            uint32  `json:"n"`
	ScriptPubKey struct {
		Asm     string `json:"asm"`
		Hex     string `json:"hex"`
		Address string `json:"address"`
	} `json:"scriptPubKey"`
}

// GetChainFunc is a function type for getting chain information
type GetChainFunc func(chainId string) (common.Chain, error)

// defaultGetChain is the default implementation that uses the GetChain function
var defaultGetChain = func(chainId string) (common.Chain, error) {
	return common.GetChain(chainId)
}

// FeeDetails represents the transaction fee information
type FeeDetails struct {
	FeeAmount    int64  `json:"feeAmount"`    // Fee amount in satoshis
	FormattedFee string `json:"formattedFee"` // Fee amount formatted in BTC
	TotalInputs  int64  `json:"totalInputs"`  // Total input value in satoshis
	TotalOutputs int64  `json:"totalOutputs"` // Total output value in satoshis
}

func PublicKeyToBitcoinAddresses(pubkey []byte) (string, string, string) {
	logger.Sugar().Infof("Bitcoin pubkey: %x", pubkey) // NOTE: don't remove this log
	pubKeyHash := btcutil.Hash160(pubkey)
	mainnetPubkey, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &MainnetConfig.Params)
	if err != nil {
		return "", "", ""
	}
	logger.Sugar().Infof("Bitcoin mainnet address: %s", mainnetPubkey)

	testnetPubkey, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &TestnetConfig.Params)
	if err != nil {
		return mainnetPubkey.EncodeAddress(), "", ""
	}
	logger.Sugar().Infof("Bitcoin testnet address: %s", testnetPubkey)

	regtestPubkey, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &RegtestConfig.Params)
	if err != nil {
		return mainnetPubkey.EncodeAddress(), testnetPubkey.EncodeAddress(), ""
	}
	logger.Sugar().Infof("Bitcoin regtest address: %s", regtestPubkey)

	return mainnetPubkey.EncodeAddress(), testnetPubkey.EncodeAddress(), regtestPubkey.EncodeAddress()
}

// Convert uncompressed public key to compressed public key
func ConvertToCompressedPublicKey(uncompressedPubKey string) (string, error) {
	uncompressedBytes, _ := hex.DecodeString(uncompressedPubKey)

	// Parse the uncompressed key
	pubKey, err := secp256k1.ParsePubKey(uncompressedBytes)
	if err != nil {
		return "", err
	}
	compressedBytes := pubKey.SerializeCompressed()
	return hex.EncodeToString(compressedBytes), nil
}

func VerifyECDSASignature(messageHashHex, signatureHex, pubKeyHex string) bool {
	// Decode message hash
	messageHash, err := hex.DecodeString(messageHashHex)
	if err != nil {
		log.Fatal("Invalid message hash:", err)
	}

	// Decode public key
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		log.Fatal("Invalid public key:", err)
	}
	pubKey, err := btcec.ParsePubKey(pubKeyBytes)
	if err != nil {
		log.Fatal("Failed to parse public key:", err)
	}

	// Decode signature
	sigBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		log.Fatal("Invalid signature:", err)
	}
	signature, err := ecdsa.ParseSignature(sigBytes)
	if err != nil {
		log.Fatal("Failed to parse signature:", err)
	}

	// Verify the signature
	return signature.Verify(messageHash, pubKey)
}

func fetchTransaction(chain common.Chain, txHash string) (*TxDetails, error) {
	// Prepare JSON-RPC request for getrawtransaction
	rpcRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getrawtransaction",
		"params":  []interface{}{txHash, 1}, // 1 for verbose output
	}

	jsonData, err := json.Marshal(rpcRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshaling RPC request: %v", err)
	}

	// Send POST request to chain.ChainUrl
	req, err := http.NewRequest("POST", chain.ChainUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if chain.RpcUsername != "" && chain.RpcPassword != "" {
		req.SetBasicAuth(chain.RpcUsername, chain.RpcPassword)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Parse RPC response
	var rpcResp RPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, fmt.Errorf("error parsing RPC response: %v", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	// Decode transaction details
	var tx TxDetails
	if err := json.Unmarshal(rpcResp.Result, &tx); err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %v", err)
	}

	return &tx, nil
}

// fetchUTXOValue fetches the value of a UTXO using BlockCypher API
// This function retrieves the value of an unspent transaction output (UTXO)
func fetchUTXOValue(chain common.Chain, txHash string) (int64, error) {
	// Fetch the full transaction details
	tx, err := fetchTransaction(chain, txHash)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch transaction: %v", err)
	}

	// For UTXOs, we're interested in the first output value
	// In most Bitcoin transactions, the first output is the actual transfer amount
	if len(tx.Vout) == 0 {
		return 0, fmt.Errorf("transaction has no outputs")
	}

	// Use OutputValue if Value is not available
	outputValue := int64(tx.Vout[0].Value * 1e8)

	return outputValue, nil
}

// Helper function to format amounts
// This function converts a string representation of an amount into a properly formatted BTC value
func getFormattedAmount(amount string, decimal int) (string, error) {
	// Convert the amount string into a big integer
	bigIntAmount := new(big.Int)
	_, success := bigIntAmount.SetString(amount, 10)
	if !success {
		return "", fmt.Errorf("error: Invalid number string")
	}

	// Format the amount into a human-readable BTC value with the specified decimal precision
	formattedAmount, err := formatUnits(bigIntAmount, decimal)
	if err != nil {
		return "", err
	}

	return formattedAmount, nil
}

// isValidBitcoinAddress checks if a Bitcoin address is valid
func isValidBitcoinAddress(address string) bool {
	// Basic regex to validate Bitcoin address format
	re := regexp.MustCompile(`^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$`)
	return re.MatchString(address)
}

// parseSerializedTransaction parses a base64 encoded serialized transaction (PSBT or raw transaction)
// and returns the unsigned transaction as a wire.MsgTx.
// If the input is a PSBT, it extracts the unsigned transaction from the global map.
// If the input is a raw transaction, it deserializes it directly.
func parseSerializedTransaction(serializedTxn string) (*wire.MsgTx, error) {
	txBytes, err := hex.DecodeString(serializedTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex transaction: %v", err)
	}

	// Check if this is a PSBT by looking for magic bytes
	if bytes.HasPrefix(txBytes, []byte("psbt\xff")) {
		return parsePSBT(txBytes)
	}

	// If not a PSBT, try to parse as a raw transaction
	return parseRawTransaction(txBytes)
}

// parsePSBT parses a Partially Signed Bitcoin Transaction and returns the unsigned transaction.
func parsePSBT(txBytes []byte) (*wire.MsgTx, error) {
	// Skip PSBT magic bytes and separator
	txBytes = txBytes[len("psbt\xff"):]

	// Extract the unsigned transaction from the global map
	var unsignedTxBytes []byte
	var offset int
	for offset < len(txBytes) {
		// Read key length
		if offset >= len(txBytes) {
			break
		}
		keyLen := int(txBytes[offset])
		offset++

		if keyLen == 0 {
			// Separator to inputs map
			offset++
			break
		}

		// Read key
		if offset+keyLen >= len(txBytes) {
			return nil, fmt.Errorf("invalid key length")
		}
		key := txBytes[offset : offset+keyLen]
		offset += keyLen

		// Read value length
		if offset >= len(txBytes) {
			return nil, fmt.Errorf("no value length")
		}
		valueLen := int(txBytes[offset])
		offset++

		if offset+valueLen > len(txBytes) {
			return nil, fmt.Errorf("invalid value length")
		}
		value := txBytes[offset : offset+valueLen]
		offset += valueLen

		// If this is the unsigned transaction (key = 0x00)
		if len(key) == 1 && key[0] == 0x00 {
			unsignedTxBytes = value
		}
	}

	if unsignedTxBytes == nil {
		return nil, fmt.Errorf("no unsigned transaction found in PSBT")
	}

	return parseRawTransaction(unsignedTxBytes)
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

func formatUnits(value *big.Int, decimals int) (string, error) {
	// Create the scaling factor as 10^decimals
	scalingFactor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))

	// Convert the value to a big.Float
	valueFloat := new(big.Float).SetInt(value)

	// Divide the value by the scaling factor
	result := new(big.Float).Quo(valueFloat, scalingFactor)

	// Convert the result to a string with the appropriate precision
	return result.Text('f', decimals), nil
}

func derEncode(signature string) (string, error) {
	// Decode hex signature
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %v", err)
	}

	if len(sigBytes) >= 70 && len(sigBytes) <= 72 {
		return signature, nil
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

// GetChainParams returns the appropriate Bitcoin chain parameters based on chainId
// chainId mapping:
// - "1000": MainNet
// - "1001": TestNet3
// - "1002": RegTest
// - "1003": SimNet
func GetChainParams(chainId string) (*chaincfg.Params, error) {
	switch chainId {
	case "1000":
		return &chaincfg.MainNetParams, nil
	case "1001":
		return &chaincfg.TestNet3Params, nil
	case "1002":
		return &chaincfg.RegressionNetParams, nil
	case "1003":
		return &chaincfg.SimNetParams, nil
	default:
		return nil, fmt.Errorf("unsupported chain ID: %s", chainId)
	}
}
