package blockchains

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util"
	"github.com/the729/lcs"

	aptosClient "github.com/portto/aptos-go-sdk/client"
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
		nodeURL:     "https://fullnode.mainnet.aptoslabs.com",
		networkID:   "mainnet",
	}

	if networkType == Testnet {
		network.nodeURL = "https://fullnode.devnet.aptoslabs.com"
		network.networkID = "testnet"
	}

	client := aptosClient.NewAptosClient(network.nodeURL)

	return &AptosBlockchain{
		BaseBlockchain: BaseBlockchain{
			chainName:       Aptos,
			network:         network,
			keyCurve:        common.CurveEddsa,
			signingEncoding: "base64",
			decimals:        7,
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
	client aptosClient.AptosClient
}

func (b *AptosBlockchain) BroadcastTransaction(txn string, signatureHex string, publicKey *string) (string, error) {
	// Initialize transaction containers
	// Transaction is the outer wrapper that includes both raw transaction and auth data
	tx := &aptosModels.Transaction{}
	// RawTransaction contains the unsigned transaction payload
	rawTxn := &aptosModels.RawTransaction{}

	// Remove '0x' prefix if present and decode hex-encoded transaction
	serializedTxn := strings.TrimPrefix(txn, "0x")
	decodedTransactionData, err := hex.DecodeString(serializedTxn)
	if err != nil {
		return "", fmt.Errorf("error decoding transaction data: %v", err)
	}

	// Deserialize the transaction using Aptos's LCS (Libra Canonical Serialization)
	// This reconstructs the RawTransaction with sender, payload, gas parameters, etc.
	err = lcs.Unmarshal(decodedTransactionData, rawTxn)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshall raw transaction: %v", err)
	}

	// Assign the deserialized raw transaction to our transaction wrapper
	tx.RawTransaction = *rawTxn

	// Decode the hex-encoded Ed25519 signature
	// Aptos uses 64-byte Ed25519 signatures
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %v", err)
	}

	// Decode the hex-encoded Ed25519 public key (32 bytes)
	// and create the transaction authenticator
	publicKeyBytes, err := hex.DecodeString(strings.TrimPrefix(*publicKey, "0x"))
	if err != nil {
		return "", fmt.Errorf("error decoding public key: %v", err)
	}

	// Create a signed transaction by adding the Ed25519 authenticator
	// This combines the public key and signature into Aptos's authentication structure
	signedTx := tx.SetAuthenticator(aptosModels.TransactionAuthenticatorEd25519{
		PublicKey: publicKeyBytes,
		Signature: signature,
	})

	// Submit the signed transaction to the Aptos network
	// The response includes the transaction hash for tracking
	response, err := b.client.SubmitTransaction(context.Background(), signedTx.UserTransaction)
	if err != nil {
		return "", fmt.Errorf("error submitting transaction: %v", err)
	}

	// Return the transaction hash which can be used to track the transaction status
	return response.Hash, nil
}

func (b *AptosBlockchain) GetTransfers(txHash string, address *string) ([]common.Transfer, error) {
	// Fetch the transaction by its hash
	// This includes the payload with transfer details and arguments
	tx, err := b.client.GetTransactionByHash(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	var transfers []common.Transfer

	// Switch on the transaction function to handle different transfer types
	switch tx.Payload.Function {

	// Case 1: New Fungible Asset Standard Transfer
	// This is Aptos's newer token standard with enhanced features
	case FUNGIBLE_ASSET_TRANSFER:
		// Extract token address from the nested map structure
		// The 'inner' field contains the actual token address
		address := tx.Payload.Arguments[0].(map[string]interface{})["inner"].(string)

		// Fetch token metadata (symbol, decimals) from the contract
		metadata, err := b.getMetadata(address)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch token metadata: %v", err)
		}

		// Get transfer amount (3rd argument in the payload)
		amount := tx.Payload.Arguments[2].(string)

		// Format amount according to token decimals
		formattedAmount, err := getFormattedAmount(amount, metadata.Decimal)
		if err != nil {
			return nil, fmt.Errorf("failed to format amount: %v", err)
		}

		// Create transfer record with token details
		transfers = append(transfers, common.Transfer{
			From:         tx.Sender,
			To:           tx.Payload.Arguments[1].(string),
			Amount:       formattedAmount,
			Token:        metadata.Symbol,
			IsNative:     false,
			TokenAddress: address,
			ScaledAmount: amount,
		})

	// Case 2: Native APT Transfer using account transfer
	// This is a direct transfer of native APT tokens
	case ACCOUNT_APT_TRANSFER:
		// Get transfer amount (2nd argument)
		amount := tx.Payload.Arguments[1].(string)
		// APT has 8 decimal places
		formattedAmount, err := getFormattedAmount(amount, 8)
		if err != nil {
			return nil, fmt.Errorf("failed to format amount: %v", err)
		}

		// Create transfer record for native APT
		transfers = append(transfers, common.Transfer{
			From:         tx.Sender,
			To:           tx.Payload.Arguments[0].(string),
			Amount:       formattedAmount,
			Token:        b.TokenSymbol(),
			IsNative:     true,
			TokenAddress: util.ZERO_ADDRESS,
			ScaledAmount: amount,
		})

	// Case 3: Generic Coin Transfers
	// This handles both APT and other token transfers using the coin module
	case ACCOUNT_COIN_TRANSFER, COIN_TRANSFER:
		// Get the type arguments which specify which coins are being transferred
		assetType := tx.Payload.TypeArguments

		// Process each asset type in the transaction
		for _, asset := range assetType {
			// Check if it's a native APT transfer
			if asset == APT_COIN_TYPE {
				// Handle native APT transfer
				amount := tx.Payload.Arguments[1].(string)
				formattedAmount, err := getFormattedAmount(amount, 8)
				if err != nil {
					return nil, fmt.Errorf("failed to format amount: %v", err)
				}

				transfers = append(transfers, common.Transfer{
					From:         tx.Sender,
					To:           tx.Payload.Arguments[0].(string),
					Amount:       formattedAmount,
					Token:        b.TokenSymbol(),
					IsNative:     true,
					TokenAddress: util.ZERO_ADDRESS,
					ScaledAmount: amount,
				})

			} else {
				// Handle other token transfers
				// Extract token address from the type argument (format: address::module::type)
				tokenAddress := strings.Split(asset, "::")[0]

				// Fetch token metadata
				metadata, err := b.getMetadata(tokenAddress)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch token metadata: %v", err)
				}

				amount := tx.Payload.Arguments[1].(string)

				// Format amount according to token decimals
				formattedAmount, err := getFormattedAmount(amount, metadata.Decimal)
				if err != nil {
					return nil, fmt.Errorf("failed to format amount: %v", err)
				}

				// Create transfer record for the token
				transfers = append(transfers, common.Transfer{
					From:         tx.Sender,
					To:           tx.Payload.Arguments[0].(string),
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
	tx, err := b.client.GetTransactionByHash(context.Background(), txnHash)
	if err != nil {
		return false, fmt.Errorf("error getting transaction by hash: %v", err)
	}

	return tx.Success, nil
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

		// Get sender's account info to retrieve the current sequence number
		// Sequence number prevents transaction replay and must be incremented for each tx
		accountInfo, err := b.client.GetAccount(context.Background(), bridgeAddress)
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

	accountInfo, err := b.client.GetAccount(context.Background(), bridgeAddress)
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
	url := fmt.Sprintf("%s/v1/accounts/%s/resources", b.network.nodeURL, address)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get coin info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	var resources []AssetInfo
	if err := json.NewDecoder(resp.Body).Decode(&resources); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	for _, info := range resources {
		if strings.Contains(info.Type, "CoinInfo") || info.Type == FUNGIBLE_ASSET_TYPE {
			return &info.Data, nil
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
