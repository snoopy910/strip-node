package identity

import (
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

	"github.com/StripChain/strip-node/algorand"
	"github.com/StripChain/strip-node/libs"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/mr-tron/base58"
	"github.com/stellar/go/keypair"
	"golang.org/x/crypto/blake2b"
)

var (
	ECDSA_CURVE       = "ecdsa"
	EDDSA_CURVE       = "eddsa"
	APTOS_EDDSA_CURVE = "aptos_eddsa"
	BITCOIN_CURVE     = "bitcoin_ecdsa"
	DOGECOIN_CURVE    = "dogecoin_ecdsa"
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
			err = errors.New("could not get a public get from the message signature")
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
	} else if identityCurve == BITCOIN_CURVE {
		fmt.Println("[VERIFY BITCOIN] Verifying bitcoin signature")
		// Parse the public key
		pubKeyBytes, err := hex.DecodeString(identity)
		if err != nil {
			fmt.Printf("[VERIFY BITCOIN] Error decoding public key: %v\n", err)
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
			fmt.Printf("[VERIFY BITCOIN] Error decoding signature: %v\n", err)
			return false, fmt.Errorf("failed to decode signature: %v", err)
		}

		var r, s *big.Int

		// Handle different signature formats
		switch len(sigBytes) {
		case 64: // Compact format: [R (32 bytes) | S (32 bytes)]
			r = new(big.Int).SetBytes(sigBytes[:32])
			s = new(big.Int).SetBytes(sigBytes[32:])
		case 65: // Compact format with recovery ID: [R (32 bytes) | S (32 bytes) | V (1 byte)]
			r = new(big.Int).SetBytes(sigBytes[:32])
			s = new(big.Int).SetBytes(sigBytes[32:64])
			// v := sigBytes[64] // recovery ID (0-3), not needed for verification
		default:
			// Try to parse as DER format
			if len(sigBytes) < 8 || len(sigBytes) > 72 || sigBytes[0] != 0x30 {
				return false, errors.New("invalid signature format: must be 64/65 bytes compact or DER format")
			}

			// Parse DER format
			// DER format: 0x30 [total-length] 0x02 [r-length] [r] 0x02 [s-length] [s]
			var offset int = 2 // Skip 0x30 and total-length

			// Extract R value
			if sigBytes[offset] != 0x02 {
				return false, errors.New("invalid DER signature format: missing R marker")
			}
			offset++
			rLen := int(sigBytes[offset])
			offset++
			r = new(big.Int).SetBytes(sigBytes[offset : offset+rLen])
			offset += rLen

			// Extract S value
			if sigBytes[offset] != 0x02 {
				return false, errors.New("invalid DER signature format: missing S marker")
			}
			offset++
			sLen := int(sigBytes[offset])
			offset++
			s = new(big.Int).SetBytes(sigBytes[offset : offset+sLen])
		}

		// Check if this is a Dogecoin message
		var hash [32]byte
		// if strings.HasPrefix(message, "DOGE:") {
		// 	// For Dogecoin, we use a special message prefix
		// 	messageBytes := []byte("Dogecoin Signed Message:\n" + strings.TrimPrefix(message, "DOGE:"))
		// 	firstHash := sha256.Sum256(messageBytes)
		// 	hash = sha256.Sum256(firstHash[:]) // Double SHA256 like Bitcoin
		// 	fmt.Println("[VERIFY SECP256K1] Using Dogecoin message format")
		// } else {
		// 	// Default Bitcoin double SHA-256
		// 	firstHash := sha256.Sum256([]byte(message))
		// 	hash = sha256.Sum256(firstHash[:]) // Second round of SHA-256
		// }

		firstHash := sha256.Sum256([]byte(message))
		hash = sha256.Sum256(firstHash[:]) // Second round of SHA-256
		// Verify the signature using ECDSA
		valid := ecdsa.Verify(pubKey, hash[:], r, s)
		fmt.Printf("[VERIFY BITCOIN] Signature is %svalid\n", func() string {
			if valid {
				return ""
			}
			return "in"
		}())
		return valid, nil
	} else if identityCurve == DOGECOIN_CURVE {
		fmt.Println("[VERIFY DOGECOIN] Verifying dogecoin signature")
		// Parse the public key
		pubKeyBytes, err := hex.DecodeString(identity)
		if err != nil {
			fmt.Printf("[VERIFY DOGECOIN] Error decoding public key: %v\n", err)
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
			fmt.Printf("[VERIFY DOGECOIN] Error decoding signature: %v\n", err)
			return false, fmt.Errorf("failed to decode signature: %v", err)
		}

		var r, s *big.Int

		// Handle different signature formats
		switch len(sigBytes) {
		case 64: // Compact format: [R (32 bytes) | S (32 bytes)]
			r = new(big.Int).SetBytes(sigBytes[:32])
			s = new(big.Int).SetBytes(sigBytes[32:])
		case 65: // Compact format with recovery ID: [R (32 bytes) | S (32 bytes) | V (1 byte)]
			r = new(big.Int).SetBytes(sigBytes[:32])
			s = new(big.Int).SetBytes(sigBytes[32:64])
			// v := sigBytes[64] // recovery ID (0-3), not needed for verification
		default:
			// Try to parse as DER format
			if len(sigBytes) < 8 || len(sigBytes) > 72 || sigBytes[0] != 0x30 {
				return false, errors.New("invalid signature format: must be 64/65 bytes compact or DER format")
			}

			// Parse DER format
			// DER format: 0x30 [total-length] 0x02 [r-length] [r] 0x02 [s-length] [s]
			var offset int = 2 // Skip 0x30 and total-length

			// Extract R value
			if sigBytes[offset] != 0x02 {
				return false, errors.New("invalid DER signature format: missing R marker")
			}
			offset++
			rLen := int(sigBytes[offset])
			offset++
			r = new(big.Int).SetBytes(sigBytes[offset : offset+rLen])
			offset += rLen

			// Extract S value
			if sigBytes[offset] != 0x02 {
				return false, errors.New("invalid DER signature format: missing S marker")
			}
			offset++
			sLen := int(sigBytes[offset])
			offset++
			s = new(big.Int).SetBytes(sigBytes[offset : offset+sLen])
		}

		// Check if this is a Dogecoin message
		var hash [32]byte
		// For Dogecoin, we use a special message prefix
		messageBytes := []byte("Dogecoin Signed Message:\n" + strings.TrimPrefix(message, "DOGE:"))
		firstHash := sha256.Sum256(messageBytes)
		hash = sha256.Sum256(firstHash[:]) // Double SHA256 like Bitcoin
		fmt.Println("[VERIFY DOGECOIN] Using Dogecoin message format")

		// Verify the signature using ECDSA
		valid := ecdsa.Verify(pubKey, hash[:], r, s)
		fmt.Printf("[VERIFY DOGECOIN] Signature is %svalid\n", func() string {
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
		// First decode the signature
		decoded, err := algorand.DecodeSignature(signature)
		if err != nil {
			return false, err
		}
		fmt.Println("Verify Signature algorand: ", message)
		// Check if message contains AlgorandFlags to determine verification method
		isRealTransaction, msg, _ := algorand.CheckFlags(message)

		// If using direct signature verification path
		if isRealTransaction {
			// This is a direct signature (not a SignedTxn)
			return algorand.VerifyDirectSignature(identity, msg, decoded)
		}

		// Try to verify as a dummy transaction
		return algorand.VerifyDummyTransaction(identity, message, decoded)
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
	} else if identityCurve == STELLAR_CURVE {
		fmt.Println("[VERIFY STELLAR] Verifying Stellar EdDSA signature")

		// Decode the signature from base64
		fmt.Printf("[VERIFY STELLAR] Raw signature: %s\n", signature)

		// First base64 decode
		signatureBytes, err := base64.StdEncoding.DecodeString(signature)
		if err != nil {
			return false, fmt.Errorf("failed to decode first layer of Stellar signature: %v", err)
		}

		// Second base64 decode - according to client code
		decodedSig, err := base64.StdEncoding.DecodeString(string(signatureBytes))
		if err != nil {
			return false, fmt.Errorf("failed to decode second layer of Stellar signature: %v", err)
		}

		fmt.Printf("[VERIFY STELLAR] Decoded signature length: %d bytes\n", len(decodedSig))

		// Decode the public key from Stellar format
		kp, err := keypair.Parse(identity)
		if err != nil {
			return false, fmt.Errorf("failed to parse Stellar keypair: %v", err)
		}

		// Based on the logs, verification succeeds with the original message
		err = kp.Verify([]byte(message), decodedSig)
		if err == nil {
			fmt.Println("[VERIFY STELLAR] Verification succeeded with original message")
			return true, nil
		}

		// If verification failed, return false with error
		fmt.Printf("[VERIFY STELLAR] Verification failed: %v\n", err)
		return false, nil
	} else {
		logger.Sugar().Errorf("unsupported curve: %s", identityCurve)
		return false, fmt.Errorf("unsupported curve: %s", identityCurve)
	}
}

func SanitiseIntent(intent libs.Intent) (string, error) {
	intentForSigning := IntentForSigning{
		Identity:      intent.Identity,
		IdentityCurve: intent.IdentityCurve,
		Operations:    []OperationForSigning{},
		Expiry:        intent.Expiry,
	}

	for _, operation := range intent.Operations {
		intentForSigning.Operations = append(intentForSigning.Operations, OperationForSigning{
			SerializedTxn:  operation.SerializedTxn,
			DataToSign:     operation.DataToSign,
			ChainId:        operation.ChainId,
			GenesisHash:    operation.GenesisHash,
			KeyCurve:       operation.KeyCurve,
			Type:           operation.Type,
			Solver:         operation.Solver,
			SolverMetadata: operation.SolverMetadata,
		})
	}

	data, err := json.Marshal(intentForSigning)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
