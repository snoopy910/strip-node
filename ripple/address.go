package ripple

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/decred/dcrd/dcrec/edwards/v2"

	// RIPEMD160 is required for Ripple address generation despite being deprecated
	rippleCrypto "github.com/rubblelabs/ripple/crypto"
)

const (
	// RIPPLE_ACCOUNT_ID is the version byte for Ripple addresses (0x00)
	RIPPLE_ACCOUNT_ID = 0x00
	// ED25519_PREFIX is the prefix byte for Ed25519 public keys (0xED)
	ED25519_PREFIX = 0xED
)

// PublicKeyToAddress converts a public key to a Ripple address and returns the formatted public key
func PublicKeyToAddress(rawKeyEddsa *eddsaKeygen.LocalPartySaveData) string {
	// Get the 32-byte public key
	pk := edwards.PublicKey{
		Curve: tss.Edwards(),
		X:     rawKeyEddsa.EDDSAPub.X(),
		Y:     rawKeyEddsa.EDDSAPub.Y(),
	}

	pkBytes2 := pk.Serialize()

	// Create 33-byte public key with ED prefix
	pkBytes := make([]byte, 33)
	pkBytes[0] = ED25519_PREFIX // Add ED prefix
	copy(pkBytes[1:], pkBytes2) // Copy the 32-byte public key

	return fmt.Sprintf("%X", pkBytes)
}

// IsValidRippleAddress checks if a string is a valid Ripple address
func IsValidRippleAddress(address string) bool {
	// Check basic format
	if len(address) < 25 || len(address) > 35 {
		return false
	}
	if address[0] != 'r' {
		return false
	}

	// Try to decode the address
	decoded, err := rippleCrypto.Base58Decode(address, rippleCrypto.ALPHABET)
	if err != nil {
		return false
	}

	// Check minimum length (1 version byte + 20 bytes RIPEMD160 + 4 bytes checksum)
	if len(decoded) != 25 {
		return false
	}

	// Verify version byte
	if decoded[0] != RIPPLE_ACCOUNT_ID {
		return false
	}

	// Verify checksum
	versionAndHash := decoded[:21] // Version byte + 20-byte hash
	checksum := decoded[21:]       // 4-byte checksum

	// Calculate checksum (double SHA256)
	sha := sha256.New()
	sha.Write(versionAndHash)
	hash := sha.Sum(nil)
	sha.Reset()
	sha.Write(hash)
	calculatedChecksum := sha.Sum(nil)[:4]

	// Compare checksums
	return bytes.Equal(checksum, calculatedChecksum)
}

// func PublicKeyToAddress(rawKeyEddsa *eddsaKeygen.LocalPartySaveData) string {
// 	pk := edwards.PublicKey{
// 		Curve: tss.Edwards(),
// 		X:     rawKeyEddsa.EDDSAPub.X(),
// 		Y:     rawKeyEddsa.EDDSAPub.Y(),
// 	}

// 	pkBytes := pk.Serialize()

// 	// Convert to compressed SEC format (33 bytes)
// 	// For Ripple, public key should start with 0x02 or 0x03 depending on Y being even or odd
// 	// pkBytes := make([]byte, 33)

// 	// // First byte is 0x02 for even Y, 0x03 for odd Y
// 	// if pk.Y.Bit(0) == 0 {
// 	// 	pkBytes[0] = 0x02
// 	// } else {
// 	// 	pkBytes[0] = 0x03
// 	// }

// 	// // Copy X coordinate to the remaining 32 bytes
// 	// xBytes := pk.X.Bytes()
// 	// copy(pkBytes[1:], xBytes)

// 	// fmt.Printf("pkBytes: %+v\n", pkBytes)
// 	// Return hex-encoded string with 0x prefix
// 	// publicKeyStr := hex.EncodeToString(pkBytes)

// 	rippleCrypto.Sha256RipeMD160()
// 	accountPublicKey, err := rippleCrypto.NewAccountPublicKey(pkBytes)
// 	if err != nil {
// 		fmt.Printf("err: %v\n", err)
// 	}
// 	fmt.Printf("accountPublicKey: %+v\n", accountPublicKey.String())

// 	address := rippleCrypto.Base58Encode(accountPublicKey.Payload(), rippleCrypto.ALPHABET)
// 	fmt.Printf("address: %+v\n", address)
// 	return address

// 	// // Step 1: SHA256 of the public key
// 	// sha := sha256.New()
// 	// sha.Write(pkBytes)
// 	// hash := sha.Sum(nil)

// 	// // Step 2: RIPEMD160 of the SHA256
// 	// ripemd := ripemd160.New()
// 	// ripemd.Write(hash)
// 	// hash = ripemd.Sum(nil)

// 	// // Step 3: Add version byte in front (0x00 for Ripple account addresses)
// 	// versionHash := append([]byte{RIPPLE_ACCOUNT_ID}, hash...)

// 	// // Step 4: Double SHA256 for checksum
// 	// sha.Reset()
// 	// sha.Write(versionHash)
// 	// hash = sha.Sum(nil)
// 	// sha.Reset()
// 	// sha.Write(hash)
// 	// hash = sha.Sum(nil)

// 	// // Step 5: Append first 4 bytes of checksum
// 	// finalBytes := append(versionHash, hash[:4]...)

// 	// // Step 6: Use Ripple's Base58 encoding
// 	// return encodeBase58(finalBytes)
// }
