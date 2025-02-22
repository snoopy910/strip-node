package identity

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func TestVerifySignatureBitcoinSuccess(t *testing.T) {
	// Generate a new Bitcoin private key
	privateKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Get the public key in uncompressed format (65 bytes, starting with 0x04)
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	publicKey = append([]byte{0x04}, publicKey...) // Prepend 0x04 for uncompressed format
	publicKeyHex := hex.EncodeToString(publicKey)

	// Create a message to sign
	message := "Hello, Bitcoin!"

	// Hash the message using double SHA-256 (as required by Bitcoin)
	firstHash := sha256.Sum256([]byte(message))
	hash := sha256.Sum256(firstHash[:])

	// Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		t.Fatalf("Failed to sign message: %v", err)
	}

	// Encode the signature (r and s) as a hex string
	signature := append(r.Bytes(), s.Bytes()...)
	signatureHex := hex.EncodeToString(signature)

	// Verify the signature
	valid, err := VerifySignature(publicKeyHex, BITCOIN_CURVE, message, signatureHex)
	if err != nil {
		t.Fatalf("Failed to verify signature: %v", err)
	}

	if !valid {
		t.Fatal("Expected signature to be valid, but it was invalid")
	}
}

func TestVerifySignatureBitcoinFail(t *testing.T) {
	// Generate a new Bitcoin private key
	privateKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Get the public key in uncompressed format (65 bytes, starting with 0x04)
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	publicKey = append([]byte{0x04}, publicKey...) // Prepend 0x04 for uncompressed format
	publicKeyHex := hex.EncodeToString(publicKey)

	// Create a message to sign
	message := "Hello, Bitcoin!"

	// Hash the message using double SHA-256 (as required by Bitcoin)
	firstHash := sha256.Sum256([]byte(message))
	hash := sha256.Sum256(firstHash[:])

	// Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		t.Fatalf("Failed to sign message: %v", err)
	}

	// Modify the signature to make it invalid (e.g., change the last byte of s)
	sBytes := s.Bytes()
	sBytes[len(sBytes)-1] ^= 0xFF // Flip the last byte
	invalidSignature := append(r.Bytes(), sBytes...)
	invalidSignatureHex := hex.EncodeToString(invalidSignature)

	// Verify the invalid signature
	valid, err := VerifySignature(publicKeyHex, BITCOIN_CURVE, message, invalidSignatureHex)
	if err != nil {
		t.Fatalf("Failed to verify signature: %v", err)
	}

	if valid {
		t.Fatal("Expected signature to be invalid, but it was valid")
	}
}
