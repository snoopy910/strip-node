package identity

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/StripChain/strip-node/sequencer"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/mr-tron/base58"
)

var (
	ECDSA_CURVE       = "ecdsa"
	EDDSA_CURVE       = "eddsa"
	APTOS_EDDSA_CURVE = "aptos_eddsa"
	SECP256K1_CURVE   = "secp256k1"
	ALGORAND_CURVE    = "algorand_eddsa"
)

type OperationForSigning struct {
	SerializedTxn  string `json:"serializedTxn"`
	DataToSign     string `json:"dataToSign"`
	ChainId        string `json:"chainId"`
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

	fmt.Println(message, signature)

	if identityCurve == ECDSA_CURVE {
		// Hash the unsigned message using EIP-191
		hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(message)) + message)
		hash := crypto.Keccak256Hash(hashedMessage)

		// Get the bytes of the signed message
		decodedMessage := hexutil.MustDecode(signature)

		// Handles cases where EIP-115 is not implemented (most wallets don't implement it)
		if decodedMessage[64] == 27 || decodedMessage[64] == 28 {
			decodedMessage[64] -= 27
		}

		// Recover a public key from the signed message
		sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), decodedMessage)
		if sigPublicKeyECDSA == nil {
			err = errors.New("Could not get a public get from the message signature")
		}
		if err != nil {
			return false, err
		}

		addr := crypto.PubkeyToAddress(*sigPublicKeyECDSA).String()

		if addr == identity {
			fmt.Println("Signature is valid")
			return true, nil
		}

		fmt.Println("Signature is invalid")

		return false, nil
	} else if identityCurve == EDDSA_CURVE {
		publicKeyBytes, _ := base58.Decode(identity)
		signatureBytes, _ := base58.Decode(signature)

		messageBytes := []byte(message)

		if ed25519.Verify(publicKeyBytes, messageBytes, signatureBytes) {
			return true, nil
		}

		return false, nil
	} else if identityCurve == APTOS_EDDSA_CURVE {
		signatureBytes, _ := hex.DecodeString(strings.TrimPrefix(signature, "0x"))
		publicKeyBytes, _ := hex.DecodeString(strings.TrimPrefix(identity, "0x"))

		messageBytes := []byte(message)

		// Verify signature
		if ed25519.Verify(publicKeyBytes, messageBytes, signatureBytes) {
			return true, nil
		}

		return false, nil
	} else if identityCurve == SECP256K1_CURVE {
		// Parse the public key
		pubKeyBytes, err := hex.DecodeString(identity)
		if err != nil {
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
			return false, fmt.Errorf("failed to decode signature: %v", err)
		}

		// The signature should be exactly 64 bytes (32 bytes for r, 32 bytes for s)
		if len(sigBytes) != 64 {
			return false, errors.New("signature must be 64 bytes long")
		}

		// Extract r and s values
		r := new(big.Int).SetBytes(sigBytes[:32])
		s := new(big.Int).SetBytes(sigBytes[32:64])

		// Hash the message using double SHA-256 (as required by Bitcoin)
		firstHash := sha256.Sum256([]byte(message))
		hash := sha256.Sum256(firstHash[:]) // Second round of SHA-256

		// Verify the signature using ECDSA
		valid := ecdsa.Verify(pubKey, hash[:], r, s)
		return valid, nil
	} else if identityCurve == ALGORAND_CURVE {
		// Decode the public key from the Algorand address (base32 encoded with checksum)
		pk, err := types.DecodeAddress(identity)
		if err != nil {
			return false, fmt.Errorf("invalid Algorand address: %v", err)
		}

		// Convert message to bytes
		msgBytes := []byte(message)

		// Decode signature from base32 (Algorand standard)
		sigBytes, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(signature)
		if err != nil {
			return false, fmt.Errorf("invalid Algorand signature encoding: %v", err)
		}

		// Convert public key bytes to ed25519.PublicKey
		pubKeyBytes := pk.GetPublicKey()
		pubKey := make(ed25519.PublicKey, ed25519.PublicKeySize)
		copy(pubKey, pubKeyBytes)
		return ed25519.Verify(pubKey, msgBytes, sigBytes), nil
	} else {
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
