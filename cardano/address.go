package cardano

import (
	"fmt"

	"github.com/echovl/cardano-go"
)

// PublicKeyToAddress converts a TSS public key to a Cardano address
func PublicKeyToAddress(pkBytes []byte) (string, string, error) {
	keyCredential, err := cardano.NewKeyCredential(pkBytes)
	if err != nil {
		return "", "", fmt.Errorf("failed to create key credential: %v", err)
	}

	mainnetAddress, err := cardano.NewBaseAddress(cardano.Mainnet, keyCredential, keyCredential)
	if err != nil {
		return "", "", fmt.Errorf("failed to create address: %v", err)
	}

	mainnetAddrStr := mainnetAddress.Bech32()

	testnetAddress, err := cardano.NewBaseAddress(cardano.Testnet, keyCredential, keyCredential)
	if err != nil {
		return "", "", fmt.Errorf("failed to create address: %v", err)
	}

	testnetAddrStr := testnetAddress.Bech32()

	return mainnetAddrStr, testnetAddrStr, nil

}
