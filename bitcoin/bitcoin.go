package bitcoin

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/btcsuite/btcd/txscript"

	"github.com/StripChain/strip-node/common"
)

// GetBitcoinTransfers fetches Bitcoin transaction details and parses them into transfers
// This function processes a transaction and extracts transfers from inputs and outputs
func GetBitcoinTransfers(chainId string, txHash string) ([]common.Transfer, *FeeDetails, error) {
	// Get chain information
	chain, err := defaultGetChain(chainId)
	if err != nil {
		return nil, nil, fmt.Errorf("chain not found")
	}

	// Fetch transaction details from BlockCypher using the chain URL and txHash
	tx, err := fetchTransaction(chain.ChainUrl, txHash)
	if err != nil {
		return nil, nil, err
	}

	// Validate transaction has inputs
	if len(tx.Inputs) == 0 {
		return nil, nil, fmt.Errorf("transaction has no inputs")
	}

	// Initialize slice to hold transfers and variables to calculate total input/output values
	var transfers []common.Transfer
	var totalInputValue int64
	var totalOutputValue int64

	// Process inputs of the transaction
	for _, input := range tx.Inputs {
		if len(input.Addresses) == 0 {
			continue // Skip input if there are no addresses
		}
		fromAddress := input.Addresses[0] // Get the first address from the input

		// Use OutputValue if Value is not available
		inputValue := input.Value
		if inputValue == 0 {
			inputValue = input.OutputValue
		}
		totalInputValue += inputValue // Sum input values

		// Process outputs of the transaction
		for _, output := range tx.Outputs {
			if len(output.Addresses) == 0 {
				continue // Skip output if there are no addresses
			}

			// Extract value (amount) from the output and convert to string
			outputValue := output.Value
			if outputValue == 0 {
				outputValue = output.OutputValue
			}
			scaledAmount := fmt.Sprintf("%d", outputValue)

			// Format the amount using the helper function
			formattedAmount, err := getFormattedAmount(scaledAmount, SATOSHI_DECIMALS)
			if err != nil {
				return nil, nil, fmt.Errorf("error formatting amount: %w", err)
			}

			// Append the transfer details to the transfers slice
			transfers = append(transfers, common.Transfer{
				From:         fromAddress,         // From address of the transfer
				To:           output.Addresses[0], // To address of the transfer
				Amount:       formattedAmount,     // Formatted transfer amount in BTC
				Token:        BTC_TOKEN_SYMBOL,    // Token symbol (BTC)
				IsNative:     true,                // Flag indicating it's a native BTC transfer
				TokenAddress: BTC_ZERO_ADDRESS,    // Token address (zero address in this case)
				ScaledAmount: scaledAmount,        // Transfer amount in satoshis
			})
		}
	}

	// Validate we found some transfers
	if len(transfers) == 0 {
		return nil, nil, fmt.Errorf("no transfers found")
	}

	// Process outputs of the transaction to calculate the total output value
	for _, output := range tx.Outputs {
		if len(output.Addresses) == 0 {
			continue // Skip output if there are no addresses
		}
		outputValue := output.Value
		if outputValue == 0 {
			outputValue = output.OutputValue
		}
		totalOutputValue += outputValue // Sum output values
	}

	// Get transaction fee from BlockCypher API response
	transactionFee := tx.Fees
	formattedFee, err := getFormattedAmount(fmt.Sprintf("%d", transactionFee), SATOSHI_DECIMALS)
	if err != nil {
		return nil, nil, fmt.Errorf("error formatting fee: %w", err)
	}

	// Create fee details structure
	feeDetails := &FeeDetails{
		FeeAmount:    transactionFee,
		FormattedFee: formattedFee,
		TotalInputs:  totalInputValue,
		TotalOutputs: totalOutputValue,
	}

	// Return the list of transfers and fee details
	return transfers, feeDetails, nil
}

func SendBitcoinTransaction(serializedTxn string, chainId string, keyCurve string, dataToSign string, signatureHex string) (string, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", err
	}
	log.Println("rpcURL", chain.ChainUrl)

	// Step 1: Decode the signature from hex
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %v", err)
	}
	log.Println("Decoded signature length:", len(signature))

	// Step 2: Parse the serialized transaction
	msgTx, err := parseSerializedTransaction(serializedTxn)
	if err != nil {
		return "", fmt.Errorf("error parsing transaction: %v", err)
	}

	// Step 3: Create proper signature script
	if len(msgTx.TxIn) == 0 {
		return "", fmt.Errorf("transaction has no inputs")
	}

	// Create a proper Bitcoin script that only contains push operations
	builder := txscript.NewScriptBuilder()
	builder.AddData(signature)          // Push the signature
	builder.AddData([]byte(dataToSign)) // Push the public key
	signatureScript, err := builder.Script()
	if err != nil {
		return "", fmt.Errorf("error building signature script: %v", err)
	}
	msgTx.TxIn[0].SignatureScript = signatureScript

	// Step 4: Serialize the signed transaction
	var signedTxBuffer bytes.Buffer
	if err := msgTx.Serialize(&signedTxBuffer); err != nil {
		return "", fmt.Errorf("error serializing signed transaction: %v", err)
	}
	signedTxHex := hex.EncodeToString(signedTxBuffer.Bytes())

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
	req.SetBasicAuth("your_rpc_user", "your_rpc_password")

	// Send request
	client := &http.Client{}
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

	txn, err := fetchTransaction(chain.ChainUrl, txnHash)
	if err != nil {
		return false, err
	}

	// Assuming a transaction is confirmed if it has at least 3 confirmations
	if txn != nil && txn.Confirmations >= 3 {
		return true, nil
	}

	return false, nil
}
