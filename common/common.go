package common

import (
	"encoding/hex"
	"errors"
	"log"
	"os"
	"unsafe"
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
