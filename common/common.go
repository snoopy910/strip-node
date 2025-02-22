package common

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"unsafe"
)

type Transfer struct {
	From         string `json:"from"`
	To           string `json:"to"`
	Amount       string `json:"amount"`
	Token        string `json:"token"`
	IsNative     bool   `json:"isNative"`
	TokenAddress string `json:"tokenAddress"`
	ScaledAmount string `json:"scaledAmount"`
}

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

type Curve string

const (
	CurveEddsa Curve = "eddsa"
	CurveEcdsa Curve = "ecdsa"
)

func ParseCurve(curve string) (Curve, error) {
	switch curve {
	case "eddsa":
		return CurveEddsa, nil
	case "ecdsa":
		return CurveEcdsa, nil
	default:
		return "", fmt.Errorf("invalid curve")
	}
}
