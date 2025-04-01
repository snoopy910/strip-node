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
	validatorPublicKeyBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		log.Fatal(err)
	}

	return *byte32(validatorPublicKeyBytes)
}

func byte32(s []byte) (a *[32]byte) {
	if len(a) <= len(s) {
		a = (*[len(a)]byte)(unsafe.Pointer(&s[0]))
	}
	return a
}

type Curve string

const (
	CurveEddsa     Curve = "eddsa"
	CurveEcdsa     Curve = "ecdsa"
	CurveAlgorand  Curve = "algorand_eddsa"
	CurveSecp256k1 Curve = "secp256k1"
	CurveAptos     Curve = "aptos_eddsa"
	CurveStellar   Curve = "stellar_eddsa"
	CurveRipple    Curve = "ripple_eddsa"
	CurveCardano   Curve = "cardano_eddsa"
)

func ParseCurve(curve string) (Curve, error) {
	switch curve {
	case "eddsa":
		return CurveEddsa, nil
	case "ecdsa":
		return CurveEcdsa, nil
	case "algorand_eddsa":
		return CurveAlgorand, nil
	case "secp256k1":
		return CurveSecp256k1, nil
	case "aptos_eddsa":
		return CurveAptos, nil
	case "stellar_eddsa":
		return CurveStellar, nil
	case "ripple_eddsa":
		return CurveRipple, nil
	case "cardano_eddsa":
		return CurveCardano, nil
	default:
		return "", fmt.Errorf("invalid curve")
	}
}
