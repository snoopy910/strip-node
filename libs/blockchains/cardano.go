package blockchains

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util"
	"github.com/blockfrost/blockfrost-go"
	"github.com/echovl/cardano-go"
	"github.com/fxamacker/cbor/v2"
)

// NewCardanoBlockchain creates a new Cardano blockchain instance
func NewCardanoBlockchain(networkType NetworkType) (IBlockchain, error) {
	mainnetAPIKey := "mainnet2NiLEqB498izHpWxY90otIxo8UxQ9YSC"
	testnetAPIKey := "preprodQqR2rJIwZJFQmMCmjcTol3HqCWKi9ZKQ"
	network := Network{
		networkType: networkType,
		nodeURL:     blockfrost.CardanoMainNet,
		apiKey:      &mainnetAPIKey,
		networkID:   "mainnet",
	}

	if networkType == Testnet {
		network.nodeURL = blockfrost.CardanoPreProd
		network.apiKey = &testnetAPIKey
		network.networkID = "preprod"
	}

	newClient := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{
			ProjectID: *network.apiKey,
			Server:    network.nodeURL,
		},
	)

	return &CardanoBlockchain{
		BaseBlockchain: BaseBlockchain{
			chainName:       Cardano,
			network:         network,
			keyCurve:        common.CurveEddsa,
			signingEncoding: "hex",
			decimals:        6,
			opTimeout:       time.Second * 10,
		},
		client: newClient,
	}, nil
}

// This is a type assertion to ensure that the CardanoBlockchain implements the IBlockchain interface
var _ IBlockchain = &CardanoBlockchain{}

// CardanoBlockchain implements the IBlockchain interface for Cardano
type CardanoBlockchain struct {
	BaseBlockchain
	client blockfrost.APIClient
}

func (b *CardanoBlockchain) BroadcastTransaction(txn string, signatureHex string, publicKey *string) (string, error) {

	// Decode the transaction bytes
	txBytes, err := hex.DecodeString(strings.TrimPrefix(txn, "0x"))
	if err != nil {
		return "", fmt.Errorf("failed to decode transaction: %v", err)
	}

	// Decode signature and public key
	sigBytes, err := hex.DecodeString(strings.TrimPrefix(signatureHex, "0x"))
	if err != nil {
		return "", fmt.Errorf("failed to decode signature: %v", err)
	}

	pubKeyBytes, err := hex.DecodeString(strings.TrimPrefix(*publicKey, "0x"))
	if err != nil {
		return "", fmt.Errorf("failed to decode public key: %v", err)
	}

	// Use a custom CBOR decoder that preserves the exact structure
	decMode, err := cbor.DecOptions{
		TagsMd:            cbor.TagsAllowed,
		ExtraReturnErrors: cbor.ExtraDecErrorUnknownField,
	}.DecMode()
	if err != nil {
		return "", fmt.Errorf("failed to create decoder: %v", err)
	}

	// Decode the transaction to get its structure
	var txData interface{}
	err = decMode.Unmarshal(txBytes, &txData)
	if err != nil {
		return "", fmt.Errorf("failed to decode transaction: %v", err)
	}

	// For Cardano Shelley transactions, the structure is typically:
	// [transaction_body, transaction_witness_set, transaction_metadata, auxiliary_data]
	// We need to ensure we maintain this exact structure

	txArray, ok := txData.([]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected transaction structure: %T", txData)
	}

	// Create a witness with the public key and signature
	vkeyWitness := []interface{}{
		pubKeyBytes,
		sigBytes,
	}

	// Check if we have a witness set (index 1)
	if len(txArray) > 1 {
		witnessSet, ok := txArray[1].(map[interface{}]interface{})
		if !ok {
			// If it's not a map, create a new one
			witnessSet = make(map[interface{}]interface{})
			txArray[1] = witnessSet
		}

		// Check if we have vkey witnesses
		vkeyWitnesses, ok := witnessSet[uint64(0)].([]interface{})
		if !ok {
			// If not, create a new array with our witness
			witnessSet[uint64(0)] = []interface{}{vkeyWitness}
		} else {
			// If we do, append our witness
			witnessSet[uint64(0)] = append(vkeyWitnesses, vkeyWitness)
		}
	} else {
		// If we don't have a witness set, create one
		witnessSet := map[interface{}]interface{}{
			uint64(0): []interface{}{vkeyWitness},
		}
		txArray = append(txArray, witnessSet)
	}

	// Ensure we have all required elements (at least 3 for a valid Shelley transaction)
	for len(txArray) < 3 {
		txArray = append(txArray, nil)
	}

	// Re-encode the transaction
	encMode, err := cbor.EncOptions{
		Sort:   cbor.SortCanonical,
		TagsMd: cbor.TagsAllowed,
	}.EncMode()
	if err != nil {
		return "", fmt.Errorf("failed to create encoder: %v", err)
	}

	signedTxBytes, err := encMode.Marshal(txArray)
	if err != nil {
		return "", fmt.Errorf("failed to marshal signed transaction: %v", err)
	}

	// Submit the transaction
	txHash, err := b.client.TransactionSubmit(context.Background(), signedTxBytes)
	if err != nil {
		return "", fmt.Errorf("error submitting transaction: %v", err)
	}

	return txHash, nil
}

func (b *CardanoBlockchain) GetTransfers(txHash string, address *string) ([]common.Transfer, error) {
	utxos, err := b.client.TransactionUTXOs(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction UTXOs: %v", err)
	}

	var transfers []common.Transfer

	// Process outputs (transfers)
	for _, output := range utxos.Outputs {
		for _, amount := range output.Amount {
			decimals, tokenName, err := b.getAssetInfo(amount.Unit)
			if err != nil {
				return nil, fmt.Errorf("failed to get asset info: %v", err)
			}

			amountInt, err := strconv.ParseInt(amount.Quantity, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse amount: %v", err)
			}

			formattedAmount := fmt.Sprintf("%f", float64(amountInt)/math.Pow(10, float64(decimals)))

			if amount.Unit == "lovelace" {
				transfers = append(transfers, common.Transfer{
					From:         utxos.Inputs[0].Address,
					To:           output.Address,
					Amount:       formattedAmount,
					Token:        tokenName,
					IsNative:     true,
					TokenAddress: util.ZERO_ADDRESS,
					ScaledAmount: amount.Quantity,
				})
				continue
			}

			transfers = append(transfers, common.Transfer{
				From:         utxos.Inputs[0].Address,
				To:           output.Address,
				Amount:       formattedAmount,
				Token:        tokenName,
				IsNative:     false,
				TokenAddress: amount.Unit,
				ScaledAmount: amount.Quantity,
			})
		}
	}

	return transfers, nil
}

// getAssetInfo fetches both decimals and token name for a given asset
func (b *CardanoBlockchain) getAssetInfo(unit string) (uint, string, error) {
	if unit == "lovelace" {
		return b.Decimals(), "ADA", nil
	}

	// Get asset details from Blockfrost
	asset, err := b.client.Asset(context.Background(), unit)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get asset info: %v", err)
	}

	tokenName, err := hex.DecodeString(asset.AssetName)
	if err != nil {
		return 0, "", fmt.Errorf("failed to decode asset name: %v", err)
	}

	if asset.Metadata == nil {
		// Assuming decimals
		return 0, string(tokenName), nil
	}

	return uint(asset.Metadata.Decimals), string(tokenName), nil
}

func (b *CardanoBlockchain) IsTransactionBroadcastedAndConfirmed(txHash string) (bool, error) {
	// Get transaction details
	tx, err := b.client.Transaction(context.Background(), txHash)
	if err != nil {
		if strings.Contains(err.Error(), "StatusCode:404") {
			return false, nil
		}
		return false, fmt.Errorf("error getting transaction: %v", err)
	}

	// Transaction is confirmed if it has a block number
	return tx.Block != "", nil
}

func (b *CardanoBlockchain) BuildWithdrawTx(bridgeAddress string,
	solverOutput string,
	userAddress string,
	tokenAddress *string,
) (string, string, error) {
	if tokenAddress == nil {
		// Parse solver output to get amount
		var solverData map[string]interface{}
		if err := json.Unmarshal([]byte(solverOutput), &solverData); err != nil {
			return "", "", fmt.Errorf("failed to parse solver output: %v", err)
		}

		txBuilder := cardano.NewTxBuilder(&cardano.ProtocolParams{})
		amountStr, ok := solverData["amount"].(string)
		if !ok {
			return "", "", fmt.Errorf("amount not found in solver output")
		}

		amount, err := strconv.ParseUint(amountStr, 10, 64)
		if err != nil {
			return "", "", fmt.Errorf("failed to parse amount: %v", err)
		}

		address, err := cardano.NewAddress(userAddress)
		if err != nil {
			return "", "", fmt.Errorf("failed to parse address: %v", err)
		}

		inputTxHash := strings.ToLower(solverData["txHash"].(string))
		txBuilder.AddInputs(&cardano.TxInput{
			TxHash: cardano.Hash32(inputTxHash),
			Amount: &cardano.Value{
				Coin: cardano.Coin(amount),
			},
		})
		txBuilder.AddOutputs(&cardano.TxOutput{
			Address: address,
			Amount: &cardano.Value{
				Coin: cardano.Coin(amount),
			},
		})

		tx, err := txBuilder.Build()
		if err != nil {
			return "", "", fmt.Errorf("failed to build transaction: %v", err)
		}
		hash, err := tx.Hash()
		if err != nil {
			return "", "", fmt.Errorf("failed to hash transaction: %v", err)
		}
		serializedTxn := tx.Hex()
		dataToSign := hash.String()
		return serializedTxn, dataToSign, nil
	}
	// Parse solver output
	var solverData map[string]interface{}
	if err := json.Unmarshal([]byte(solverOutput), &solverData); err != nil {
		return "", "", fmt.Errorf("failed to parse solver output: %v", err)
	}

	txBuilder := cardano.NewTxBuilder(&cardano.ProtocolParams{})
	amountStr, ok := solverData["amount"].(string)
	if !ok {
		return "", "", fmt.Errorf("amount not found in solver output")
	}

	amount, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse amount: %v", err)
	}

	address, err := cardano.NewAddress(userAddress)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse address: %v", err)
	}

	inputTxHash := strings.ToLower(solverData["txHash"].(string))

	// FIX: This is a temporary fillin and should be addressed
	policyID := cardano.NewPolicyIDFromHash(cardano.Hash28(*tokenAddress))
	assetName := cardano.NewAssetName(*tokenAddress)
	txBuilder.AddInputs(&cardano.TxInput{
		TxHash: cardano.Hash32(inputTxHash),
		Amount: &cardano.Value{
			MultiAsset: cardano.NewMultiAsset().Set(policyID, cardano.NewAssets().Set(assetName, cardano.BigNum(amount))),
		},
	})
	txBuilder.AddOutputs(&cardano.TxOutput{
		Address: address,
		Amount: &cardano.Value{
			MultiAsset: cardano.NewMultiAsset().Set(policyID, cardano.NewAssets().Set(assetName, cardano.BigNum(amount))),
		},
	})

	tx, err := txBuilder.Build()
	if err != nil {
		return "", "", fmt.Errorf("failed to build transaction: %v", err)
	}
	hash, err := tx.Hash()
	if err != nil {
		return "", "", fmt.Errorf("failed to hash transaction: %v", err)
	}
	serializedTxn := tx.Hex()
	dataToSign := hash.String()

	return serializedTxn, dataToSign, nil
}

func (b *CardanoBlockchain) RawPublicKeyBytesToAddress(pkBytes []byte, networkType NetworkType) (string, error) {
	keyCredential, err := cardano.NewKeyCredential(pkBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create key credential: %v", err)
	}

	if networkType == Mainnet {
		mainnetAddress, err := cardano.NewBaseAddress(cardano.Mainnet, keyCredential, keyCredential)
		if err != nil {
			return "", fmt.Errorf("failed to create address: %v", err)
		}

		return mainnetAddress.Bech32(), nil
	}

	testnetAddress, err := cardano.NewBaseAddress(cardano.Testnet, keyCredential, keyCredential)
	if err != nil {
		return "", fmt.Errorf("failed to create address: %v", err)
	}

	return testnetAddress.Bech32(), nil
}

func (b *CardanoBlockchain) RawPublicKeyToPublicKeyStr(pkBytes []byte) (string, error) {
	return "", nil
}

func (b *CardanoBlockchain) ExtractDestinationAddress(serializedTxn string) (string, string, error) {
	var tx cardano.Tx
	txBytes, err := hex.DecodeString(serializedTxn)
	destAddress := ""
	if err != nil {
		return "", "", fmt.Errorf("error decoding Cardano transaction", err)
	}
	if err := json.Unmarshal(txBytes, &tx); err != nil {
		return "", "", fmt.Errorf("error parsing Cardano transaction", err)
	}
	destAddress = tx.Body.Outputs[0].Address.String()
	return destAddress, "", nil
}
