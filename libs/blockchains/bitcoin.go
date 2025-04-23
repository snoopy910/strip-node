package blockchains

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/libs"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const (
	BTC_ZERO_ADDRESS = "0x0000000000000000000000000000000000000000"
)

// NewBitcoinBlockchain creates a new Bitcoin blockchain instance
func NewBitcoinBlockchain(networkType NetworkType) (IBlockchain, error) {
	user := "your_rpc_user"
	pass := "your_rpc_password"
	network := Network{
		networkType: networkType,
		nodeURL:     "172.17.0.1:8332",
		networkID:   "mainnet",
	}
	chainParams := &chaincfg.MainNetParams

	if networkType == Testnet {
		network.nodeURL = "172.17.0.1:18332"
		network.networkID = "testnet"
		chainParams = &chaincfg.TestNet3Params
	}

	if networkType == Regnet {
		network.nodeURL = "bitcoind:8332"
		network.networkID = "regtest"
		chainParams = &chaincfg.RegressionNetParams
		user = "bitcoin"
		pass = "bitcoin"
	}

	// Configure RPC client connection
	connCfg := &rpcclient.ConnConfig{
		Host:         network.nodeURL,
		User:         user,
		Pass:         pass,
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	// Create a new RPC client
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create bitcoin RPC client: %w", err)
	}

	return &BitcoinBlockchain{
		BaseBlockchain: BaseBlockchain{
			chainName:       Bitcoin,
			network:         network,
			keyCurve:        common.CurveEcdsa,
			signingEncoding: "hex",
			tokenSymbol:     "BTC",
			decimals:        8,
			opTimeout:       time.Minute * 1,
		},
		client:        client,
		chainParams:   chainParams,
		confirmations: 3,
	}, nil
}

// This is a type assertion to ensure that the BitcoinBlockchain implements the IBlockchain interface
var _ IBlockchain = &BitcoinBlockchain{}

// BitcoinBlockchain implements the IBlockchain interface for Bitcoin
type BitcoinBlockchain struct {
	BaseBlockchain
	client        *rpcclient.Client
	chainParams   *chaincfg.Params
	confirmations int
}

func (b *BitcoinBlockchain) BroadcastTransaction(txn string, signatureHex string, publicKey *string) (string, error) {
	if publicKey == nil {
		return "", errors.New("public key is required for bitcoin broadcast")
	}
	var scriptPubKey []byte
	var isSegWit bool

	// Check if the input is a hex public key
	if len(*publicKey) == 66 || len(*publicKey) == 130 { // Compressed (33 bytes * 2) or uncompressed (65 bytes * 2)
		pubKeyBytes, err := hex.DecodeString(*publicKey)
		if err != nil {
			return "", fmt.Errorf("failed to decode public key: %v", err)
		}
		// For public key input, we'll use P2PKH (legacy)
		scriptPubKey = pubKeyBytes
		isSegWit = false
	} else {
		// Try to decode as address
		decodedAddress, err := btcutil.DecodeAddress(*publicKey, b.chainParams)
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
	msgTx, err := parseSerializedTransaction(txn)
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
		pubKeyBytes, err := hex.DecodeString(*publicKey)
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

	// Send the raw transaction
	txHash, err := b.client.SendRawTransaction(msgTx, true) // Allow high fees if necessary
	if err != nil {
		return "", fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	return txHash.String(), nil
}

// GetTransfers currently only extracts outputs (potential receivers and amounts).
// Determining the "from" address accurately requires fetching previous transactions.
func (b *BitcoinBlockchain) GetTransfers(txHash string, address *string) ([]common.Transfer, error) {
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction hash format: %w", err)
	}

	// Get verbose transaction details
	txVerbose, err := b.client.GetRawTransactionVerbose(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get raw transaction details: %w", err)
	}

	var transfers []common.Transfer

	for _, input := range txVerbose.Vin {
		prevHash, err := chainhash.NewHashFromStr(input.Txid)
		if err != nil {
			return nil, fmt.Errorf("invalid transaction hash format: %w", err)
		}
		// Fetch previous transaction to get input address
		prevTx, err := b.client.GetRawTransactionVerbose(prevHash)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch previous tx %s: %v", input.Txid, err)
		}
		if int(input.Vout) >= len(prevTx.Vout) {
			continue // Skip invalid vout
		}
		prevOutput := prevTx.Vout[input.Vout]
		fromAddress := prevOutput.ScriptPubKey.Address
		// inputValue := int64(prevOutput.Value * 1e8) // BTC to satoshis
		// totalInputValue += inputValue

		// Process outputs
		for _, output := range txVerbose.Vout {
			outputValue := int64(output.Value * 1e8) // BTC to satoshis
			scaledAmount := fmt.Sprintf("%d", outputValue)
			formattedAmount, err := getFormattedAmount(scaledAmount, int(b.Decimals()))
			if err != nil {
				return nil, fmt.Errorf("error formatting amount: %w", err)
			}

			transfers = append(transfers, common.Transfer{
				From:         fromAddress,
				To:           output.ScriptPubKey.Address,
				Amount:       formattedAmount,
				Token:        b.TokenSymbol(),
				IsNative:     true,
				TokenAddress: BTC_ZERO_ADDRESS,
				ScaledAmount: scaledAmount,
			})
		}
	}

	if len(transfers) == 0 {
		return nil, fmt.Errorf("no transfers found")
	}

	return transfers, nil
}

func (b *BitcoinBlockchain) IsTransactionBroadcastedAndConfirmed(txHash string) (bool, error) {
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return false, fmt.Errorf("invalid transaction hash format: %w", err)
	}

	txDetails, err := b.client.GetTransaction(hash)
	if err != nil {
		return false, fmt.Errorf("failed to get transaction details: %w", err)
	}

	// Consider a transaction confirmed if it has at least 1 confirmation.
	// This threshold might need to be configurable.
	return txDetails.Confirmations >= int64(b.confirmations), nil
}

func (b *BitcoinBlockchain) BuildWithdrawTx(account string,
	solverOutput string,
	recipient string,
	tokenAddress *string,
) (string, string, error) {
	if tokenAddress != nil {
		return "", "", errors.New("token transfers are not supported for Bitcoin")
	}
	var solverData map[string]interface{}
	if err := json.Unmarshal([]byte(solverOutput), &solverData); err != nil {
		return "", "", fmt.Errorf("failed to parse solver output: %v", err)
	}

	amount, ok := solverData["amount"].(string)
	if !ok {
		return "", "", fmt.Errorf("amount not found in solver output")
	}

	// Create a new Bitcoin transaction
	var msgTx wire.MsgTx
	msgTx.Version = wire.TxVersion

	// Parse amount
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse amount: %w", err)
	}

	// Convert amount to satoshis
	amountSatoshis := int64(amountFloat * 100000000)

	// Create transaction output
	addr, err := btcutil.DecodeAddress(recipient, b.chainParams)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode recipient address: %w", err)
	}

	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return "", "", fmt.Errorf("failed to create output script: %w", err)
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
	_, err = btcutil.DecodeAddress(account, b.chainParams)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode from address: %w", err)
	}

	// For P2WPKH, we use empty SignatureScript and put the actual script in witness
	txIn.SignatureScript = []byte{}

	// Serialize the transaction
	var buf bytes.Buffer
	if err := msgTx.Serialize(&buf); err != nil {
		return "", "", fmt.Errorf("failed to serialize transaction: %w", err)
	}

	// firstInputScriptPubKey := inputScriptPubKeys[0]
	// firstInputValue := inputValues[0]

	// // Assuming the first input is P2WPKH for sighash calculation.
	// // This is a major assumption and needs to be verified or made dynamic.
	// sigHashes := txscript.NewTxSigHashes(&msgTx)
	// sighashBytes, err := txscript.CalcWitnessSigHash(firstInputScriptPubKey, sigHashes, txscript.SigHashAll, &msgTx, 0, int64(firstInputValue))
	// if err != nil {
	// 	return "", "", fmt.Errorf("failed to calculate sighash for input 0: %w", err)
	// }
	// dataToSign := hex.EncodeToString(sighashBytes)
	// TODO: Missing data to sign implementation
	dataToSign := ""
	return hex.EncodeToString(buf.Bytes()), dataToSign, nil
}

func (b *BitcoinBlockchain) RawPublicKeyBytesToAddress(pkBytes []byte, networkType NetworkType) (string, error) {
	return "", errors.New("RawPublicKeyBytesToAddress not implemented yet")
}

func (b *BitcoinBlockchain) RawPublicKeyToPublicKeyStr(pkBytes []byte) (string, error) {
	return "", errors.New("RawPublicKeyToPublicKeyStr not implemented yet")
}

func (b *BitcoinBlockchain) ExtractDestinationAddress(operation *libs.Operation) (string, error) {
	// For Bitcoin, decode the serialized transaction to get output address
	var tx wire.MsgTx
	txBytes, err := hex.DecodeString(*operation.SerializedTxn)
	if err != nil {
		return "", err
	}
	if err := tx.Deserialize(bytes.NewReader(txBytes)); err != nil {
		return "", err
	}
	// Get the first output's address (assuming it's the bridge address)
	if len(tx.TxOut) > 0 {
		_, addrs, _, err := txscript.ExtractPkScriptAddrs(tx.TxOut[0].PkScript, nil)
		if err != nil || len(addrs) == 0 {
			return "", err
		}
		return addrs[0].String(), nil
	}
	return "", fmt.Errorf("no output destination bitcoin address")
}
