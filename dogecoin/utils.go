package dogecoin

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

// PublicKeyToAddress converts a hex-encoded public key to a Dogecoin address
// The public key should be in compressed or uncompressed format
func PublicKeyToAddress(publicKeyHex string) (string, error) {
	// Decode the hex public key
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key hex: %v", err)
	}

	// Parse the public key
	publicKey, err := btcec.ParsePubKey(publicKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %v", err)
	}

	// Convert the public key to a Dogecoin address
	// Using mainnet parameters with Dogecoin-specific settings
	params := chaincfg.MainNetParams
	params.PubKeyHashAddrID = 0x1E // Dogecoin mainnet P2PKH prefix (30)

	pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &params)
	if err != nil {
		return "", fmt.Errorf("failed to create address: %v", err)
	}

	// Validate that we got a proper Dogecoin mainnet address (starts with 'D')
	addrStr := addr.String()
	if !strings.HasPrefix(addrStr, "D") {
		return "", fmt.Errorf("generated address does not have correct Dogecoin mainnet prefix: %s", addrStr)
	}

	return addrStr, nil
}

// PublicKeyToTestnetAddress converts a hex-encoded public key to a Dogecoin testnet address
func PublicKeyToTestnetAddress(publicKeyHex string) (string, error) {
	// Decode the hex public key
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key hex: %v", err)
	}

	// Parse the public key
	publicKey, err := btcec.ParsePubKey(publicKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %v", err)
	}

	// Convert the public key to a Dogecoin testnet address
	params := chaincfg.TestNet3Params
	params.PubKeyHashAddrID = 0x71 // Dogecoin testnet P2PKH prefix (113)

	pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &params)
	if err != nil {
		return "", fmt.Errorf("failed to create address: %v", err)
	}

	// Validate that we got a proper Dogecoin testnet address (starts with 'n')
	addrStr := addr.String()
	if !strings.HasPrefix(addrStr, "n") {
		return "", fmt.Errorf("generated address does not have correct Dogecoin testnet prefix: %s", addrStr)
	}

	return addrStr, nil
}
