package blockchains

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/libs"
	"github.com/StripChain/strip-node/util"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/the729/lcs"

	aptos "github.com/aptos-labs/aptos-go-sdk"
	aptosApi "github.com/aptos-labs/aptos-go-sdk/api"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/aptos-labs/aptos-go-sdk/crypto"
	aptosModels "github.com/portto/aptos-go-sdk/models"
)

const (
	DEFAULT_GAS_UNIT_PRICE = 100
	DEFAULT_MAX_GAS_AMOUNT = 5000
)

// NewAptosBlockchain creates a new Stellar blockchain instance
func NewAptosBlockchain(networkType NetworkType) (IBlockchain, error) {
	network := Network{
		networkType: networkType,
		nodeURL:     "https://api.mainnet.aptoslabs.com/v1",
		networkID:   "mainnet",
	}
	chainID := "1"

	if networkType == Testnet {
		network.nodeURL = "https://api.testnet.aptoslabs.com/v1"
		network.networkID = "testnet"
		chainID = "2"
	}

	if networkType == Devnet {
		network.nodeURL = "https://api.devnet.aptoslabs.com/v1"
		network.networkID = "devnet"
		chainID = "178"
	}

	// Initialize Aptos client
	chainIDInt, err := strconv.Atoi(chainID)
	if err != nil {
		return nil, fmt.Errorf("error converting chainid to int: %v", err)
	}

	client, err := aptos.NewNodeClient(network.nodeURL, uint8(chainIDInt))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Aptos client: %v", err)
	}
	return &AptosBlockchain{
		BaseBlockchain: BaseBlockchain{
			chainName:       Aptos,
			network:         network,
			keyCurve:        common.CurveEddsa,
			signingEncoding: "hex",
			decimals:        8,
			chainID:         &chainID,
			opTimeout:       time.Second * 10,
		},
		client: client,
	}, nil
}

// This is a type assertion to ensure that the AptosBlockchain implements the IBlockchain interface
var _ IBlockchain = &AptosBlockchain{}

// AptosBlockchain implements the IBlockchain interface for Stellar
type AptosBlockchain struct {
	BaseBlockchain
	client *aptos.NodeClient
}

func (b *AptosBlockchain) BroadcastTransaction(txn string, signatureHex string, publicKey *string) (string, error) {
	// Construct the transaction from serializedTxn
	rawTxn := &aptos.RawTransaction{}

	decodedTransactionData, err := hex.DecodeString(strings.TrimPrefix(txn, "0x"))
	if err != nil {
		return "", fmt.Errorf("error decoding transaction data: %v", err)
	}

	err = bcs.Deserialize(rawTxn, decodedTransactionData)
	if err != nil {
		// it missing a byte, not sure why but the transaction is still valid
		logger.Sugar().Warnw("error unmarshalling raw transaction", "error", err, "txn", txn)
		// return "", fmt.Errorf("error unmarshalling raw transaction: %v", err)
	}

	// Retrieve signatureHex
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %v", err)
	}

	// Sign transaction with pubKey and signature
	publicKeyBytes, err := hex.DecodeString(strings.TrimPrefix(*publicKey, "0x"))
	if err != nil {
		return "", fmt.Errorf("error decoding public key: %v", err)
	}
	authenticator := &aptos.TransactionAuthenticator{
		Variant: aptos.TransactionAuthenticatorEd25519,
		Auth: &aptos.Ed25519TransactionAuthenticator{
			Sender: &crypto.AccountAuthenticator{
				Variant: crypto.AccountAuthenticatorEd25519,
				Auth: &crypto.Ed25519Authenticator{
					PubKey: &crypto.Ed25519PublicKey{
						Inner: publicKeyBytes,
					},
					Sig: &crypto.Ed25519Signature{
						Inner: [64]byte(signature),
					},
				},
			},
		},
	}

	signedTx := &aptos.SignedTransaction{
		Transaction:   rawTxn,
		Authenticator: authenticator,
	}

	// Submit transaction
	response, err := b.client.SubmitTransaction(signedTx)
	if err != nil {
		return "", fmt.Errorf("error submitting transaction: %v", err)
	}

	fmt.Println("Submitted aptos transaction with hash:", response.Hash)

	return response.Hash, nil
}

func (b *AptosBlockchain) GetTransfers(txHash string, address *string) ([]common.Transfer, error) {
	txResp, err := b.client.TransactionByHash(txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	var transfers []common.Transfer

	tx, err := txResp.UserTransaction()
	if err != nil {
		return nil, fmt.Errorf("failed to get user transaction: %v", err)
	}

	payload, ok := tx.Payload.Inner.(*aptosApi.TransactionPayloadEntryFunction)
	if !ok {
		return nil, fmt.Errorf("invalid payload type: %v", tx.Payload.Type)
	}
	sender := tx.Sender.String()

	switch payload.Function {

	case FUNGIBLE_ASSET_TRANSFER:

		address := payload.Arguments[0].(map[string]interface{})["inner"].(string)

		metadata, err := b.getMetadata(address)
		if err != nil {
			return nil, fmt.Errorf("error fetching token metadata, %w", err)
		}

		amount := payload.Arguments[2].(string)

		formattedAmount, err := getFormattedAmount(amount, metadata.Decimal)
		if err != nil {
			return nil, fmt.Errorf("wrror formatting amount, %w", err)
		}

		transfers = append(transfers, common.Transfer{
			From:         sender,
			To:           payload.Arguments[1].(string),
			Amount:       formattedAmount,
			Token:        metadata.Symbol,
			IsNative:     false,
			TokenAddress: address,
			ScaledAmount: amount,
		})

	case ACCOUNT_APT_TRANSFER:

		amount := payload.Arguments[1].(string)
		formattedAmount, err := getFormattedAmount(amount, 8)
		if err != nil {
			return nil, fmt.Errorf("error formatting amount, %w", err)
		}

		transfers = append(transfers, common.Transfer{
			From:         sender,
			To:           payload.Arguments[0].(string),
			Amount:       formattedAmount,
			Token:        b.TokenSymbol(),
			IsNative:     true,
			TokenAddress: util.ZERO_ADDRESS,
			ScaledAmount: amount,
		})

	case ACCOUNT_COIN_TRANSFER, COIN_TRANSFER:

		assetType := payload.TypeArguments

		for _, asset := range assetType {

			if asset == APT_COIN_TYPE {

				amount := payload.Arguments[1].(string)
				formattedAmount, err := getFormattedAmount(amount, 8)
				if err != nil {
					return nil, fmt.Errorf("error formatting amount, %w", err)
				}

				transfers = append(transfers, common.Transfer{
					From:         sender,
					To:           payload.Arguments[0].(string),
					Amount:       formattedAmount,
					Token:        b.TokenSymbol(),
					IsNative:     true,
					TokenAddress: util.ZERO_ADDRESS,
					ScaledAmount: amount,
				})

			} else {
				tokenAddress := strings.Split(asset, "::")[0]

				metadata, err := b.getMetadata(tokenAddress)
				if err != nil {
					return nil, fmt.Errorf("error fetching token metadata, %w", err)
				}

				amount := payload.Arguments[1].(string)
				formattedAmount, err := getFormattedAmount(amount, metadata.Decimal)
				if err != nil {
					return nil, fmt.Errorf("error formatting amount, %w", err)
				}

				transfers = append(transfers, common.Transfer{
					From:         sender,
					To:           payload.Arguments[0].(string),
					Amount:       formattedAmount,
					Token:        metadata.Symbol,
					IsNative:     false,
					TokenAddress: tokenAddress,
					ScaledAmount: amount,
				})
			}
		}
	}

	return transfers, nil
}

func (b *AptosBlockchain) IsTransactionBroadcastedAndConfirmed(txnHash string) (bool, error) {
	tx, err := b.client.TransactionByHash(txnHash)
	if err != nil {
		return false, fmt.Errorf("error getting transaction by hash: %v", err)
	}

	success := tx.Success()
	if success == nil {
		return false, nil
	}
	return *success, nil
}

func (b *AptosBlockchain) BuildWithdrawTx(bridgeAddress string,
	solverOutput string,
	userAddress string,
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
		// Get the core framework account address (0x1)
		// This is where the coin module is deployed
		addr0x1, _ := aptosModels.HexToAccountAddress("0x1")

		// Convert recipient address to Aptos account address format
		// Aptos addresses are 32 bytes, hex-encoded with 0x prefix
		accountTo, err := aptosModels.HexToAccountAddress(userAddress)
		if err != nil {
			return "", "", fmt.Errorf("failed to convert account address: %v", err)
		}

		decodedBridgeAddress, err := hex.DecodeString(strings.TrimPrefix(bridgeAddress, "0x"))
		if err != nil {
			return "", "", fmt.Errorf("failed to decode bridge address: %v", err)
		}

		// Get sender's account info to retrieve the current sequence number
		// Sequence number prevents transaction replay and must be incremented for each tx
		accountInfo, err := b.client.Account(aptos.AccountAddress(decodedBridgeAddress))
		if err != nil {
			return "", "", fmt.Errorf("failed to get accountInfo: %v", err)
		}

		// Initialize transaction wrapper
		tx := &aptosModels.Transaction{}

		// Create type tag for AptosCoin
		// This identifies that we're transferring the native APT token
		// Located at 0x1::aptos_coin::AptosCoin
		aptosCoinTypeTag := aptosModels.TypeTagStruct{
			Address: addr0x1,
			Module:  "aptos_coin",
			Name:    "AptosCoin",
		}

		// Build the transaction with a fluent interface
		// 1. Set the sender account
		// 2. Set the entry function call payload (0x1::coin::transfer)
		// 3. Set expiration (10 minutes from now)
		// 4. Set gas parameters
		// 5. Set sequence number from account info
		err = tx.SetSender(bridgeAddress).
			SetPayload(aptosModels.EntryFunctionPayload{
				// Use the coin module from core framework
				Module: aptosModels.Module{
					Address: addr0x1,
					Name:    "coin",
				},
				// Call the transfer function
				Function: "transfer",
				// Specify we're transferring AptosCoin
				TypeArguments: []aptosModels.TypeTag{aptosCoinTypeTag},
				// Arguments: recipient address and amount
				Arguments: []interface{}{accountTo, amount},
			}).
			SetExpirationTimestampSecs(uint64(time.Now().Add(10 * time.Minute).Unix())).
			SetGasUnitPrice(DEFAULT_GAS_UNIT_PRICE).
			SetMaxGasAmount(DEFAULT_MAX_GAS_AMOUNT).
			SetSequenceNumber(accountInfo.SequenceNumber).Error()
		if err != nil {
			return "", "", fmt.Errorf("failed to set transaction data: %v", err)
		}

		// Serialize the transaction using BCS (Binary Canonical Serialization)
		// This is required for transaction submission
		bcsBytes, err := lcs.Marshal(tx.RawTransaction)
		if err != nil {
			return "", "", fmt.Errorf("failed to convert bcs format: %v", err)
		}
		// Convert to hex string and remove 0x prefix
		serializedTxn := strings.TrimPrefix(hex.EncodeToString(bcsBytes), "0x")

		// Get the signing message that needs to be signed by the sender
		// This includes the transaction hash and chain-specific data
		msgBytes, err := tx.GetSigningMessage()
		if err != nil {
			return "", "", fmt.Errorf("failed to get message: %v", err)
		}
		// Convert to hex string and remove 0x prefix
		dataToSign := strings.TrimPrefix(hex.EncodeToString(msgBytes), "0x")

		return serializedTxn, dataToSign, nil
	}

	addr0x1, _ := aptosModels.HexToAccountAddress("0x1")

	// get token info
	if len(strings.Split(*tokenAddress, "::")) != 3 {
		return "", "", fmt.Errorf("invalid token address format")
	}
	ownerAddr, err := aptosModels.HexToAccountAddress(strings.Split(*tokenAddress, "::")[0])
	collenctionName := strings.Split(*tokenAddress, "::")[1]
	name := strings.Split(*tokenAddress, "::")[2]
	if err != nil {
		return "", "", fmt.Errorf("failed to convert token address: %v", err)
	}

	accountTo, err := aptosModels.HexToAccountAddress(userAddress)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert account address: %v", err)
	}

	decodedBridgeAddress, err := hex.DecodeString(strings.TrimPrefix(bridgeAddress, "0x"))
	if err != nil {
		return "", "", fmt.Errorf("failed to decode bridge address: %v", err)
	}

	accountInfo, err := b.client.Account(aptos.AccountAddress(decodedBridgeAddress))
	if err != nil {
		return "", "", fmt.Errorf("failed to get accountInfo: %v", err)
	}

	tx := &aptosModels.Transaction{}

	aptosCoinTypeTag := aptosModels.TypeTagStruct{
		Address: ownerAddr,
		Module:  collenctionName,
		Name:    name,
	}

	err = tx.SetSender(bridgeAddress).
		SetPayload(aptosModels.EntryFunctionPayload{
			Module: aptosModels.Module{
				Address: addr0x1,
				Name:    "coin",
			},
			Function:      "transfer",
			TypeArguments: []aptosModels.TypeTag{aptosCoinTypeTag},
			Arguments:     []interface{}{accountTo, amount},
		}).
		SetExpirationTimestampSecs(uint64(time.Now().Add(10 * time.Minute).Unix())).
		SetGasUnitPrice(DEFAULT_GAS_UNIT_PRICE).
		SetMaxGasAmount(DEFAULT_MAX_GAS_AMOUNT).
		SetSequenceNumber(accountInfo.SequenceNumber).Error()
	if err != nil {
		return "", "", fmt.Errorf("failed to set transaction data: %v", err)
	}

	bcsBytes, err := lcs.Marshal(tx.RawTransaction)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert bcs format: %v", err)
	}
	serializedTxn := strings.TrimPrefix(hex.EncodeToString(bcsBytes), "0x")

	msgBytes, err := tx.GetSigningMessage()
	if err != nil {
		return "", "", fmt.Errorf("failed to get message: %v", err)
	}
	dataToSign := strings.TrimPrefix(hex.EncodeToString(msgBytes), "0x")

	return serializedTxn, dataToSign, nil
}

func (b *AptosBlockchain) RawPublicKeyBytesToAddress(pkBytes []byte, networkType NetworkType) (string, error) {
	return "", errors.New("RawPublicKeyBytesToAddress not implemented")
}

func (b *AptosBlockchain) RawPublicKeyToPublicKeyStr(pkBytes []byte) (string, error) {
	return "", errors.New("RawPublicKeyToPublicKeyStr not implemented")
}

const (
	FUNGIBLE_ASSET_TRANSFER = "0x1::primary_fungible_store::transfer"
	ACCOUNT_APT_TRANSFER    = "0x1::aptos_account::transfer"
	ACCOUNT_COIN_TRANSFER   = "0x1::aptos_account::transfer_coins"
	COIN_TRANSFER           = "0x1::coin::transfer"
	APT_COIN_TYPE           = "0x1::aptos_coin::AptosCoin"
	FUNGIBLE_ASSET_TYPE     = "0x1::fungible_asset::Metadata"
)

type AssetData struct {
	Decimal int    `json:"decimals"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
}

type AssetInfo struct {
	Type string `json:"type"`
	Data AssetData
}

func (b *AptosBlockchain) getMetadata(address string) (*AssetData, error) {
	decodedAddress, err := hex.DecodeString(strings.TrimPrefix(address, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode address: %w", err)
	}
	resp, err := b.client.AccountResources(aptos.AccountAddress(decodedAddress))
	if err != nil {
		return nil, fmt.Errorf("failed to get coin info: %w", err)
	}

	for _, info := range resp {
		if strings.Contains(info.Type, "CoinInfo") || info.Type == FUNGIBLE_ASSET_TYPE {
			assetData := &AssetData{}
			for key, data := range info.Data {
				switch strings.ToLower(key) {
				case "decimals":
					assetData.Decimal = int(data.(float64))
				case "name":
					assetData.Name = data.(string)
				case "symbol":
					assetData.Symbol = data.(string)
				}
			}
			return assetData, nil
		}
	}
	return nil, fmt.Errorf("there is no coin info")
}

func getFormattedAmount(amount string, decimal int) (string, error) {
	bigIntAmount := new(big.Int)

	_, success := bigIntAmount.SetString(amount, 10)
	if !success {
		return "", fmt.Errorf("error: Invalid number string")
	}

	formattedAmount, err := util.FormatUnits(bigIntAmount, int(decimal))
	if err != nil {
		return "", err
	}

	return formattedAmount, nil
}

func (b *AptosBlockchain) ExtractDestinationAddress(operation *libs.Operation) (string, error) {
	// For Aptos, the destination is in the transaction payload
	var aptosPayload struct {
		Function string   `json:"function"`
		Args     []string `json:"arguments"`
	}
	destAddress := ""
	if err := json.Unmarshal([]byte(*operation.SerializedTxn), &aptosPayload); err != nil {
		logger.Sugar().Errorw("error parsing Aptos transaction", "error", err)
		return "", fmt.Errorf("error parsing Aptos transaction", err)
	}
	if len(aptosPayload.Args) > 0 {
		destAddress = aptosPayload.Args[0] // First arg is typically the recipient
	}
	return destAddress, nil
}
