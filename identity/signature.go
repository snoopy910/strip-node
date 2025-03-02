package identity

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/StripChain/strip-node/sequencer"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/blake2b"
)

var (
	ECDSA_CURVE       = "ecdsa"
	EDDSA_CURVE       = "eddsa"
	APTOS_EDDSA_CURVE = "aptos_eddsa"
	SECP256K1_CURVE   = "secp256k1"
	SUI_EDDSA_CURVE   = "sui_eddsa" // Sui uses Ed25519 for native transactions
	STELLAR_CURVE     = "stellar_eddsa"
	ALGORAND_CURVE    = "algorand_eddsa"
	RIPPLE_CURVE      = "ripple_eddsa"
	CARDANO_CURVE     = "cardano_eddsa"
)

type OperationForSigning struct {
	SerializedTxn  string `json:"serializedTxn"`
	DataToSign     string `json:"dataToSign"`
	ChainId        string `json:"chainId"`
	GenesisHash    string `json:"genesisHash"`
	KeyCurve       string `json:"keyCurve"`
	Type           string `json:"type"`
	Solver         string `json:"solver"`
	SolverMetadata string `json:"solverMetadata"`
}

type IntentForSigning struct {
	Operations    []OperationForSigning `json:"operations"`
	Identity      string                `json:"identity"`
	IdentityCurve string                `json:"identityCurve"`
	Expiry        uint64                `json:"expiry"`
}

func VerifySignature(
	identity string,
	identityCurve string,
	message string,
	signature string,
) (bool, error) {

	fmt.Printf("[VERIFY] Starting signature verification for identity: %s with curve: %s\n", identity, identityCurve)
	fmt.Println(message, signature)

	if identityCurve == ECDSA_CURVE {
		fmt.Println("[VERIFY ECDSA] Verifying ECDSA signature")
		// Hash the unsigned message using EIP-191
		hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(message)) + message)
		hash := crypto.Keccak256Hash(hashedMessage)
		fmt.Printf("[VERIFY ECDSA] Message hash: %s\n", hash.Hex())

		// Get the bytes of the signed message
		decodedMessage := hexutil.MustDecode(signature)
		fmt.Printf("[VERIFY ECDSA] Decoded signature length: %d bytes\n", len(decodedMessage))

		// Handles cases where EIP-115 is not implemented (most wallets don't implement it)
		if decodedMessage[64] == 27 || decodedMessage[64] == 28 {
			fmt.Printf("[VERIFY ECDSA] Adjusting V value from: %d\n", decodedMessage[64])
			decodedMessage[64] -= 27
			fmt.Printf("[VERIFY ECDSA] New V value: %d\n", decodedMessage[64])
		}

		// Recover a public key from the signed message
		sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), decodedMessage)
		if sigPublicKeyECDSA == nil {
			fmt.Println("[VERIFY ECDSA] Failed to recover public key from signature")
			err = errors.New("Could not get a public get from the message signature")
		}
		if err != nil {
			fmt.Printf("[VERIFY ECDSA] Error recovering public key: %v\n", err)
			return false, err
		}

		addr := crypto.PubkeyToAddress(*sigPublicKeyECDSA).String()
		fmt.Printf("[VERIFY ECDSA] Recovered address: %s\n", addr)
		fmt.Printf("[VERIFY ECDSA] Expected address: %s\n", identity)

		if addr == identity {
			fmt.Println("[VERIFY ECDSA] Signature is valid")
			return true, nil
		}

		fmt.Println("[VERIFY ECDSA] Signature is invalid")

		return false, nil

	} else if identityCurve == EDDSA_CURVE {
		fmt.Println("[VERIFY EDDSA] Verifying EdDSA signature")
		publicKeyBytes, _ := base58.Decode(identity)
		signatureBytes, _ := base58.Decode(signature)

		messageBytes := []byte(message)

		if ed25519.Verify(publicKeyBytes, messageBytes, signatureBytes) {
			fmt.Println("[VERIFY EDDSA] Signature is valid")
			return true, nil
		}

		fmt.Println("[VERIFY EDDSA] Signature is invalid")
		return false, nil
	} else if identityCurve == SECP256K1_CURVE {
		fmt.Println("[VERIFY SECP256K1] Verifying secp256k1 signature")
		// Parse the public key
		pubKeyBytes, err := hex.DecodeString(identity)
		if err != nil {
			fmt.Printf("[VERIFY SECP256K1] Error decoding public key: %v\n", err)
			return false, fmt.Errorf("failed to decode public key: %v", err)
		}

		// Uncompressed public keys start with 0x04 and are 65 bytes long
		if len(pubKeyBytes) != 65 || pubKeyBytes[0] != 0x04 {
			return false, errors.New("public key must be uncompressed and 65 bytes long")
		}

		// Extract X and Y coordinates
		x := new(big.Int).SetBytes(pubKeyBytes[1:33])
		y := new(big.Int).SetBytes(pubKeyBytes[33:65])

		// Create the ECDSA public key
		pubKey := &ecdsa.PublicKey{
			Curve: secp256k1.S256(), // Use the SECP256K1 curve from go-ethereum
			X:     x,
			Y:     y,
		}

		// Parse the signature
		sigBytes, err := hex.DecodeString(signature)
		if err != nil {
			fmt.Printf("[VERIFY SECP256K1] Error decoding signature: %v\n", err)
			return false, fmt.Errorf("failed to decode signature: %v", err)
		}

		// The signature should be exactly 64 bytes (32 bytes for r, 32 bytes for s)
		if len(sigBytes) != 64 {
			return false, errors.New("signature must be 64 bytes long")
		}

		// Extract r and s values
		r := new(big.Int).SetBytes(sigBytes[:32])
		s := new(big.Int).SetBytes(sigBytes[32:64])

		// Check if this is a Dogecoin message
		var hash [32]byte
		if strings.HasPrefix(message, "DOGE:") {
			// For Dogecoin, we use a special message prefix
			messageBytes := []byte("Dogecoin Signed Message:\n" + strings.TrimPrefix(message, "DOGE:"))
			firstHash := sha256.Sum256(messageBytes)
			hash = sha256.Sum256(firstHash[:]) // Double SHA256 like Bitcoin
			fmt.Println("[VERIFY SECP256K1] Using Dogecoin message format")
		} else {
			// Default Bitcoin double SHA-256
			firstHash := sha256.Sum256([]byte(message))
			hash = sha256.Sum256(firstHash[:]) // Second round of SHA-256
		}

		// Verify the signature using ECDSA
		valid := ecdsa.Verify(pubKey, hash[:], r, s)
		fmt.Printf("[VERIFY SECP256K1] Signature is %svalid\n", func() string {
			if valid {
				return ""
			}
			return "in"
		}())
		return valid, nil
	} else if identityCurve == SUI_EDDSA_CURVE {
		fmt.Println("[VERIFY SUI_EDDSA] Verifying Sui EdDSA signature")

		// Remove 0x prefix from public key
		identity = strings.TrimPrefix(identity, "0x")
		if len(identity) != 64 {
			return false, fmt.Errorf("invalid public key length: expected 64 hex chars, got %d", len(identity))
		}

		// For Sui, we need to verify that the message has the correct prefix
		if !strings.HasPrefix(message, "Sui Message:") {
			return false, fmt.Errorf("invalid message format: must start with 'Sui Message:'")
		}

		// Convert message to bytes
		messageBytes := []byte(message)
		fmt.Printf("[VERIFY SUI_EDDSA] Message bytes: %x\n", messageBytes)

		// Decode the public key from hex
		publicKeyBytes, err := hex.DecodeString(identity)
		if err != nil {
			return false, fmt.Errorf("failed to decode public key: %v", err)
		}
		fmt.Printf("[VERIFY SUI_EDDSA] Public key bytes: %x\n", publicKeyBytes)

		// Decode base64 signature
		signatureBytes, err := base64.StdEncoding.DecodeString(signature)
		if err != nil {
			return false, fmt.Errorf("failed to decode base64 signature: %v", err)
		}
		fmt.Printf("[VERIFY SUI_EDDSA] Signature bytes: %x\n", signatureBytes)

		// Hash the message with Blake2b as per Sui's requirement
		hasher, _ := blake2b.New256(nil) // Using nil key for keyless hashing
		hasher.Write(messageBytes)
		msgHash := hasher.Sum(nil)

		// Verify using ed25519 which internally uses SHA-512 as per RFC 8032
		verified := ed25519.Verify(publicKeyBytes, msgHash, signatureBytes)

		fmt.Printf("[VERIFY SUI_EDDSA] Verification result: %v\n", verified)
		return verified, nil

	} else if identityCurve == ALGORAND_CURVE {
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
		// msgBytes := []byte(message) ?

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
		// Decode signature from base64 (Algorand standard)
		fmt.Println("verify signature: ", signature)
		sigBytes, err := base64.StdEncoding.DecodeString(signature)
		if err != nil {
			return false, fmt.Errorf("invalid Algorand signature encoding: %v", err)
		}
		verified := ed25519.Verify(pubKey, msgBytes, sigBytes)
		fmt.Println("verified signature algorand: ", verified)
		return ed25519.Verify(pubKey, msgBytes, sigBytes), nil
	} else if identityCurve == APTOS_EDDSA_CURVE || identityCurve == RIPPLE_CURVE || identityCurve == CARDANO_CURVE {
		fmt.Println("[VERIFY APTOS_EDDSA] Verifying Aptos EdDSA signature")

		// Remove 0x prefix from public key
		identity = strings.TrimPrefix(identity, "0x")
		if len(identity) != 64 {
			return false, fmt.Errorf("invalid public key length: expected 64 hex chars, got %d", len(identity))
		}

		// Create the prefix message format that matches Aptos wallet
		// Note: we need to use \n not \r\n for line endings to match the client
		// The message should be the raw JSON string without additional encoding
		prefixedMsg := fmt.Sprintf("APTOS\nmessage: %s\nnonce: random_string", message)

		// Convert message to bytes using TextEncoder equivalent
		messageBytes := []byte(prefixedMsg)
		fmt.Printf("[VERIFY APTOS_EDDSA] Message bytes: %x\n", messageBytes)
		fmt.Printf("[VERIFY APTOS_EDDSA] Message string:\n%s\n", prefixedMsg)

		// Decode the public key from hex
		publicKeyBytes, err := hex.DecodeString(identity)
		if err != nil {
			return false, fmt.Errorf("failed to decode public key: %v", err)
		}
		fmt.Printf("[VERIFY APTOS_EDDSA] Public key bytes: %x\n", publicKeyBytes)

		// Remove 0x prefix from signature and decode
		signature = strings.TrimPrefix(signature, "0x")
		signatureBytes, err := hex.DecodeString(signature)
		if err != nil {
			return false, fmt.Errorf("failed to decode signature: %v", err)
		}
		fmt.Printf("[VERIFY APTOS_EDDSA] Signature bytes: %x\n", signatureBytes)

		// Verify using ed25519 which is equivalent to nacl.sign.detached.verify
		// Both use the same Ed25519 verification algorithm
		verified := ed25519.Verify(publicKeyBytes, messageBytes, signatureBytes)

		fmt.Printf("[VERIFY APTOS_EDDSA] Verification result: %v\n", verified)
		return verified, nil
	} else {
		fmt.Printf("unsupported curve: %s", identityCurve)
		return false, fmt.Errorf("unsupported curve: %s", identityCurve)
	}
}

func SanitiseIntent(intent sequencer.Intent) (string, error) {
	intentForSigning := IntentForSigning{
		Identity:      intent.Identity,
		IdentityCurve: intent.IdentityCurve,
		Expiry:        intent.Expiry,
	}

	for _, operation := range intent.Operations {
		operationForSigning := OperationForSigning{
			SerializedTxn:  operation.SerializedTxn,
			DataToSign:     operation.DataToSign,
			ChainId:        operation.ChainId,
			GenesisHash:    operation.GenesisHash,
			KeyCurve:       operation.KeyCurve,
			Type:           operation.Type,
			Solver:         operation.Solver,
			SolverMetadata: operation.SolverMetadata,
		}

		intentForSigning.Operations = append(intentForSigning.Operations, operationForSigning)
	}

	jsonBytes, err := json.Marshal(intentForSigning)
	if err != nil {
		return "", err
	}

	dst := &bytes.Buffer{}
	if err := json.Compact(dst, jsonBytes); err != nil {
		panic(err)
	}

	return dst.String(), nil
}
