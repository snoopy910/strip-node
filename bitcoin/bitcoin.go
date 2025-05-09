package bitcoin

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"github.com/StripChain/strip-node/common"
)

// GetBitcoinTransfers fetches Bitcoin transaction details and parses them into transfers
// This function processes a transaction and extracts transfers from inputs and outputs
func GetBitcoinTransfers(chainId string, txHash string) ([]common.Transfer, *FeeDetails, error) {
	// Get chain information
	chain, err := defaultGetChain(chainId) // Assume this exists from your codebase
	if err != nil {
		return nil, nil, fmt.Errorf("chain not found: %v", err)
	}

	// Fetch transaction details
	tx, err := fetchTransaction(chain, txHash)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch transaction: %v", err)
	}

	// Validate transaction has inputs
	if len(tx.Vin) == 0 {
		return nil, nil, fmt.Errorf("transaction has no inputs")
	}

	var transfers []common.Transfer
	var totalInputValue int64
	var totalOutputValue int64

	// Process inputs
	for _, input := range tx.Vin {
		// Fetch previous transaction to get input address
		prevTx, err := fetchTransaction(chain, input.TxID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch previous tx %s: %v", input.TxID, err)
		}
		if int(input.Vout) >= len(prevTx.Vout) {
			continue // Skip invalid vout
		}
		prevOutput := prevTx.Vout[input.Vout]
		fromAddress := prevOutput.ScriptPubKey.Address
		inputValue := int64(prevOutput.Value * 1e8) // BTC to satoshis
		totalInputValue += inputValue

		// Process outputs
		for _, output := range tx.Vout {
			outputValue := int64(output.Value * 1e8) // BTC to satoshis
			scaledAmount := fmt.Sprintf("%d", outputValue)
			formattedAmount, err := getFormattedAmount(scaledAmount, SATOSHI_DECIMALS)
			if err != nil {
				return nil, nil, fmt.Errorf("error formatting amount: %w", err)
			}

			transfers = append(transfers, common.Transfer{
				From:         fromAddress,
				To:           output.ScriptPubKey.Address,
				Amount:       formattedAmount,
				Token:        BTC_TOKEN_SYMBOL,
				IsNative:     true,
				TokenAddress: BTC_ZERO_ADDRESS,
				ScaledAmount: scaledAmount,
			})
		}
	}

	if len(transfers) == 0 {
		return nil, nil, fmt.Errorf("no transfers found")
	}

	// Calculate total output value
	for _, output := range tx.Vout {
		outputValue := int64(output.Value * 1e8)
		totalOutputValue += outputValue
	}

	// Calculate fee
	transactionFee := totalInputValue - totalOutputValue
	formattedFee, err := getFormattedAmount(fmt.Sprintf("%d", transactionFee), SATOSHI_DECIMALS)
	if err != nil {
		return nil, nil, fmt.Errorf("error formatting fee: %w", err)
	}

	feeDetails := &FeeDetails{
		FeeAmount:    transactionFee,
		FormattedFee: formattedFee,
		TotalInputs:  totalInputValue,
		TotalOutputs: totalOutputValue,
	}

	return transfers, feeDetails, nil
}

func SendBitcoinTransaction(serializedTxn string, chainId string, keyCurve string, dataToSign string, address string, signatureHex string) (string, error) {
	log.Println("SendBitcoinTransaction", serializedTxn, chainId, keyCurve, dataToSign, address, signatureHex)
	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", err
	}
	log.Println("rpcURL", chain.ChainUrl)

	// Step 0: Handle input that could be either an address or public key
	var scriptPubKey []byte
	var isSegWit bool

	// First try to decode as address
	netParam, err := GetChainParams(chain.ChainId)
	if err != nil {
		return "", err
	}

	// Check if the input is a hex public key
	if len(address) == 66 || len(address) == 130 { // Compressed (33 bytes * 2) or uncompressed (65 bytes * 2)
		pubKeyBytes, err := hex.DecodeString(address)
		if err != nil {
			return "", fmt.Errorf("failed to decode public key: %v", err)
		}
		// For public key input, we'll use P2PKH (legacy)
		scriptPubKey = pubKeyBytes
		isSegWit = false
	} else {
		// Try to decode as address
		decodedAddress, err := btcutil.DecodeAddress(address, netParam)
		if err != nil {
			return "", fmt.Errorf("failed to decode address: %v", err)
		}

		// Determine if the address is SegWit
		switch addr := decodedAddress.(type) {
		case *btcutil.AddressWitnessPubKeyHash:
			isSegWit = true
			scriptPubKey = addr.WitnessProgram()
		case *btcutil.AddressWitnessScriptHash:
			isSegWit = true
			scriptPubKey = addr.WitnessProgram()
		default:
			isSegWit = false
			scriptPubKey = decodedAddress.ScriptAddress()
		}
	}

	// Step 1: Parse the transaction
	msgTx, err := parseSerializedTransaction(serializedTxn)
	if err != nil {
		return "", fmt.Errorf("error parsing transaction: %v", err)
	}

	// Step 2: Create DER signature
	derSignatureHex, err := derEncode(signatureHex)
	if err != nil {
		return "", fmt.Errorf("error encoding signature: %v", err)
	}
	log.Println("DER signature:", derSignatureHex)
	derSignature, err := hex.DecodeString(derSignatureHex)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %v", err)
	}

	// Step 3: Handle transaction signing based on address type
	if isSegWit {
		// For SegWit, we use witness data
		witness := wire.TxWitness{derSignature, scriptPubKey}
		msgTx.TxIn[0].Witness = witness
		// Empty the signature script for witness transactions
		msgTx.TxIn[0].SignatureScript = []byte{}
	} else {
		// For legacy P2PKH, signature script should be: <sig> <pubkey>
		builder := txscript.NewScriptBuilder()

		// Add DER signature with SIGHASH_ALL if not already present
		if len(derSignature) == 0 || derSignature[len(derSignature)-1] != byte(txscript.SigHashAll) {
			derSignature = append(derSignature, byte(txscript.SigHashAll))
		}
		builder.AddData(derSignature)

		// For P2PKH, we need the full public key, not its hash
		pubKeyBytes, err := hex.DecodeString(address)
		if err != nil {
			return "", fmt.Errorf("error decoding public key: %v", err)
		}
		builder.AddData(pubKeyBytes)

		sigScript, err := builder.Script()
		if err != nil {
			return "", fmt.Errorf("error building signature script: %v", err)
		}
		msgTx.TxIn[0].SignatureScript = sigScript
		// Empty the witness for legacy transactions
		msgTx.TxIn[0].Witness = wire.TxWitness{}
	}

	// Step 4: Serialize the signed transaction
	var signedTxBuffer bytes.Buffer
	if err := msgTx.Serialize(&signedTxBuffer); err != nil {
		return "", fmt.Errorf("error serializing signed transaction: %v", err)
	}
	signedTxHex := hex.EncodeToString(signedTxBuffer.Bytes())
	log.Println("Signed transaction hex:", signedTxHex)

	// Step 5: Prepare and send RPC request
	rpcRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "sendrawtransaction",
		"params":  []interface{}{signedTxHex},
	}

	jsonData, err := json.Marshal(rpcRequest)
	if err != nil {
		return "", fmt.Errorf("error marshaling RPC request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", chain.ChainUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	if chain.RpcUsername != "" {
		req.SetBasicAuth(chain.RpcUsername, chain.RpcPassword)
	}

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending transaction: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	var rpcResponse struct {
		Result string `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &rpcResponse); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	// Check for RPC error
	if rpcResponse.Error != nil {
		return "", fmt.Errorf("RPC error %d: %s", rpcResponse.Error.Code, rpcResponse.Error.Message)
	}

	return rpcResponse.Result, nil
}

func CheckBitcoinTransactionConfirmed(chainId string, txnHash string) (bool, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return false, err
	}

	txn, err := fetchTransaction(chain, txnHash)
	if err != nil {
		return false, err
	}

	// Assuming a transaction is confirmed if it has at least 3 confirmations
	if txn != nil && txn.Confirmations >= 3 {
		return true, nil
	}

	return false, nil
}
