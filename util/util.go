package util

import (
	"log"
	"math/big"
	"os"
	"strconv"
)

var ZERO_ADDRESS = "0x0000000000000000000000000000000000000000"

func LookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("LookupEnvOrInt[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func LookupEnvOrBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.ParseBool(val)
		if err != nil {
			log.Fatalf("LookupEnvOrInt[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
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
