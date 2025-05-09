package aptos

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/StripChain/strip-node/common"
)

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

func getMetadata(chain common.Chain, address string) (*AssetData, error) {
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

func FormatUnits(value *big.Int, decimals int) (string, error) {
	// Create the scaling factor as 10^decimals
	scalingFactor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))

	// Convert the value to a big.Float
	valueFloat := new(big.Float).SetInt(value)

	// Divide the value by the scaling factor
	result := new(big.Float).Quo(valueFloat, scalingFactor)

	// Convert the result to a string with the appropriate precision
	return result.Text('f', decimals), nil
}
