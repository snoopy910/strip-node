package sui

import (
	"fmt"
	"math/big"
)

const (
	SUI_TYPE = "0x2::sui::SUI"
)

func getFormattedAmount(amount string, decimal uint8) (string, error) {
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
