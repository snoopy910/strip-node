package common

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"unsafe"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
)

var (
	ECDSA_CURVE = "ecdsa"
	EDDSA_CURVE = "eddsa"
)

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func PublicKeyStrToBytes32(publicKey string) [32]byte {
	pubkey := string([]rune(publicKey)[2:])
	signerPublicKeyBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		log.Fatal(err)
	}

	return *byte32(signerPublicKeyBytes)
}

func byte32(s []byte) (a *[32]byte) {
	if len(a) <= len(s) {
		a = (*[len(a)]byte)(unsafe.Pointer(&s[0]))
	}
	return a
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
