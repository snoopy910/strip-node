package sequencer

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/StripChain/strip-node/util"
	aptosClient "github.com/portto/aptos-go-sdk/client"
)

const (
	FUNGIBLE_ASSET_TRANSFER = "0x1::primary_fungible_store::transfer"
	ACCOUNT_APT_TRANSFER    = "0x1::aptos_account::transfer"
	ACCOUNT_COIN_TRANSFER   = "0x1::aptos_account::transfer_coins"
	COIN_TRANSFER           = "0x1::coin::transfer"
	APT_COIN_TYPE           = "0x1::aptos_coin::AptosCoin"
	FUNGIBLE_ASSET_TYPE     = "0x1::fungible_asset::Metadata"
)

func GetAptosTransfers(chainId string, txHash string) ([]Transfer, error) {
	// Get chain configuration
	chain, err := GetChain(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain config: %v", err)
	}

	// Intialize Aptos client
	client := aptosClient.NewAptosClient(chain.ChainUrl)

	tx, err := client.GetTransactionByHash(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	var transfers []Transfer

	switch tx.Payload.Function {

	case FUNGIBLE_ASSET_TRANSFER:

		address := tx.Payload.Arguments[0].(map[string]interface{})["inner"].(string)

		metadata, err := getMetadata(chain, address)
		if err != nil {
			fmt.Println("Error fetching token metadata, ", err)
		}

		amount := tx.Payload.Arguments[2].(string)

		formattedAmount, err := getFormattedAmount(amount, metadata.Decimal)
		if err != nil {
			fmt.Println("Error formatting amount, %w", err)
		}

		transfers = append(transfers, Transfer{
			From:         tx.Sender,
			To:           tx.Payload.Arguments[1].(string),
			Amount:       formattedAmount,
			Token:        metadata.Symbol,
			IsNative:     false,
			TokenAddress: address,
			ScaledAmount: amount,
		})

	case ACCOUNT_APT_TRANSFER:

		amount := tx.Payload.Arguments[1].(string)
		formattedAmount, err := getFormattedAmount(amount, 8)
		if err != nil {
			fmt.Println("Error formatting amount, ", err)
		}

		transfers = append(transfers, Transfer{
			From:         tx.Sender,
			To:           tx.Payload.Arguments[0].(string),
			Amount:       formattedAmount,
			Token:        chain.TokenSymbol,
			IsNative:     true,
			TokenAddress: util.ZERO_ADDRESS,
			ScaledAmount: amount,
		})

	case ACCOUNT_COIN_TRANSFER, COIN_TRANSFER:

		assetType := tx.Payload.TypeArguments

		for _, asset := range assetType {

			if asset == APT_COIN_TYPE {

				amount := tx.Payload.Arguments[1].(string)
				formattedAmount, err := getFormattedAmount(amount, 8)
				if err != nil {
					fmt.Println("Error formatting amount, ", err)
				}

				transfers = append(transfers, Transfer{
					From:         tx.Sender,
					To:           tx.Payload.Arguments[0].(string),
					Amount:       formattedAmount,
					Token:        chain.TokenSymbol,
					IsNative:     true,
					TokenAddress: util.ZERO_ADDRESS,
					ScaledAmount: amount,
				})

			} else {
				tokenAddress := strings.Split(asset, "::")[0]

				metadata, err := getMetadata(chain, tokenAddress)
				if err != nil {
					fmt.Println("Error fetching token metadata, ", err)
				}

				amount := tx.Payload.Arguments[1].(string)

				formattedAmount, err := getFormattedAmount(amount, metadata.Decimal)
				if err != nil {
					fmt.Println("Error formatting amount, %w", err)
				}

				transfers = append(transfers, Transfer{
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

type AssetData struct {
	Decimal int    `json:"decimals"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
}

type AssetInfo struct {
	Type string `json:"type"`
	Data AssetData
}

func getMetadata(chain Chain, address string) (*AssetData, error) {
	url := fmt.Sprintf("%s/v1/accounts/%s/resources", chain.ChainUrl, address)

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

	formattedAmount, err := FormatUnits(bigIntAmount, int(decimal))
	if err != nil {
		return "", err
	}

	return formattedAmount, nil
}
