package identity

import (
	"bytes"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/Silent-Protocol/go-sio/sequencer"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
)

var (
	ECDSA_CURVE = "ecdsa"
	EDDSA_CURVE = "eddsa"
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
	} else {
		return false, nil
	}
}

func SanitiseIntent(intent sequencer.Intent) (string, error) {
	intentForSigning := IntentForSigning{
		Identity:      intent.Identity,
		IdentityCurve: intent.IdentityCurve,
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
