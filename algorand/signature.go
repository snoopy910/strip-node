package algorand

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

// DecodeSignature decodes an Algorand signature from base64 or hex format
func DecodeSignature(signature string) ([]byte, error) {
	// Try to decode as base64
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		// If not base64, try as a direct hex string
		decoded, err = hex.DecodeString(signature)
		if err != nil {
			return nil, fmt.Errorf("failed to decode algorand signature: %v", err)
		}
	}
	return decoded, nil
}

// CheckFlags checks if a message contains AlgorandFlags with IsRealTransaction set
func CheckFlags(message string) (bool, error) {
	// Check if message is a JSON with AlgorandFlags.IsRealTransaction set to true
	var messageObj struct {
		AlgorandFlags *struct {
			IsRealTransaction bool `json:"isRealTransaction"`
		} `json:"algorandFlags"`
	}

	jsonErr := json.Unmarshal([]byte(message), &messageObj)
	isRealTransaction := false
	if jsonErr == nil && messageObj.AlgorandFlags != nil {
		isRealTransaction = messageObj.AlgorandFlags.IsRealTransaction
	}
	return isRealTransaction, jsonErr
}

// VerifyDirectSignature verifies a direct ed25519 signature for Algorand
func VerifyDirectSignature(identity string, message string, signature []byte) (bool, error) {
	// Decode the public key from the Algorand address (base32 encoded with checksum)
	address, err := types.DecodeAddress(identity)
	if err != nil {
		return false, fmt.Errorf("invalid Algorand address: %v", err)
	}

	// Convert public key bytes to ed25519.PublicKey
	pubKeyBytes := address[:]
	pubKey := make(ed25519.PublicKey, ed25519.PublicKeySize)
	copy(pubKey, pubKeyBytes)

	// Convert message to bytes
	fmt.Println("verify message: ", message)
	var msgBytes []byte
	var js map[string]interface{}
	// Unmarshal the string into the map. If no error, it's valid JSON.
	err = json.Unmarshal([]byte(message), &js)
	if err == nil {
		prefix := []byte("MX")
		messageBytes := []byte(message)
		msgBytes = append(prefix, messageBytes...)
	} else {
		msgBytes, err = base64.StdEncoding.DecodeString(message)
		fmt.Println("verify message algorand bytes: ", message)
		if err != nil {
			return false, fmt.Errorf("invalid Algorand message encoding: %v", err)
		}
	}
	// Verify the signature
	fmt.Println("verify signature: ", signature)
	verified := ed25519.Verify(pubKey, msgBytes, signature)
	fmt.Println("verified signature algorand: ", verified)
	return verified, nil
}

// VerifyDummyTransaction verifies an Algorand signature using the dummy transaction approach
func VerifyDummyTransaction(identity string, signature []byte) (bool, error) {
	// Try to decode as a SignedTxn (dummy transaction)
	var stx types.SignedTxn
	err := msgpack.Decode(signature, &stx)
	if err != nil {
		// If it's not a SignedTxn, try direct signature verification
		return VerifyDirectSignature(identity, "", signature)
	}

	// Extract the sender's address from the transaction
	sender := stx.Txn.Sender.String()

	// Compare with the claimed identity
	if sender != identity {
		return false, fmt.Errorf("sender address %s does not match claimed identity %s", sender, identity)
	}

	// Convert Algorand address to public key
	pubKey, err := AddressToPubKey(identity)
	if err != nil {
		return false, fmt.Errorf("failed to convert address to public key: %v", err)
	}

	// Recreate the canonical transaction bytes that were signed (prefixed with "TX")
	txnBytes := msgpack.Encode(stx.Txn)
	signingBytes := append([]byte("TX"), txnBytes...)

	// Verify the signature
	if !ed25519.Verify(pubKey, signingBytes, stx.Sig[:]) {
		return false, fmt.Errorf("algorand signature verification failed")
	}

	fmt.Println("verified Algorand dummy transaction signature successfully")
	return true, nil
}

// AddressToPubKey converts an Algorand address to its public key
func AddressToPubKey(address string) (ed25519.PublicKey, error) {
	decodedAddress, err := types.DecodeAddress(address)
	if err != nil {
		return nil, err
	}
	return decodedAddress[:], nil
}
