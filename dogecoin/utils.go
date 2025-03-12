package dogecoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
)

// PublicKeyToAddress converts a hex-encoded public key to a Dogecoin address
// The public key should be in compressed or uncompressed format
func PublicKeyToAddress(publicKeyHex string) (string, error) {
	// Decode the hex public key
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key hex: %v", err)
	}

	// Parse the public key
	publicKey, err := btcec.ParsePubKey(publicKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %v", err)
	}

	// Convert the public key to a Dogecoin address
	// Using mainnet parameters with Dogecoin-specific settings
	params := chaincfg.MainNetParams
	params.PubKeyHashAddrID = 0x1E // Dogecoin mainnet P2PKH prefix (30)

	pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &params)
	if err != nil {
		return "", fmt.Errorf("failed to create address: %v", err)
	}

	// Validate that we got a proper Dogecoin mainnet address (starts with 'D')
	addrStr := addr.String()
	if !strings.HasPrefix(addrStr, "D") {
		return "", fmt.Errorf("generated address does not have correct Dogecoin mainnet prefix: %s", addrStr)
	}

	return addrStr, nil
}

// PublicKeyToTestnetAddress converts a hex-encoded public key to a Dogecoin testnet address
func PublicKeyToTestnetAddress(publicKeyHex string) (string, error) {
	// Decode the hex public key
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key hex: %v", err)
	}

	// Parse the public key
	publicKey, err := btcec.ParsePubKey(publicKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %v", err)
	}

	// Convert the public key to a Dogecoin testnet address
	params := chaincfg.TestNet3Params
	params.PubKeyHashAddrID = 0x71 // Dogecoin testnet P2PKH prefix (113)

	pubKeyHash := btcutil.Hash160(publicKey.SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &params)
	if err != nil {
		return "", fmt.Errorf("failed to create address: %v", err)
	}

	// Validate that we got a proper Dogecoin testnet address (starts with 'n')
	addrStr := addr.String()
	if !strings.HasPrefix(addrStr, "n") {
		return "", fmt.Errorf("generated address does not have correct Dogecoin testnet prefix: %s", addrStr)
	}

	return addrStr, nil
}

func getTatumApiKey(chainID string) (string, error) {
	var key string
	if chainID == "2000" {
		key = "TATUM_API_KEY_MAINNET"
		if val, ok := os.LookupEnv(key); ok {
			return val, nil
		}
		// TODO: secure tatum mainnet api key
		return "t-67cb0e957f1a5a5a2483e093-03079b3a500a4f39bf4d651b", nil
	} else if chainID == "2001" {
		key = "TATUM_API_KEY_TESTNET"
		if val, ok := os.LookupEnv(key); ok {
			return val, nil
		}
		// TODO: secure tatum testnet api key
		return "t-67cb0e957f1a5a5a2483e093-eeb92712de9c4144a0edbcca", nil
	}
	return "", fmt.Errorf("[Error] Wrong chainID for tatum api key")
}

// GetChainParams returns the appropriate Dogecoin chain parameters based on chainId
// chainId mapping:
// - "2000": MainNet
// - "2001": TestNet
func GetChainParams(chainId string) (*chaincfg.Params, error) {
	switch chainId {
	case "2000":
		chainParams := &chaincfg.MainNetParams
		chainParams.PubKeyHashAddrID = 0x1e
		chainParams.ScriptHashAddrID = 0x16
		chainParams.PrivateKeyID = 0x9e
		return chainParams, nil
	case "2001":
		chainParams := &chaincfg.TestNet3Params
		chainParams.PubKeyHashAddrID = 0x71
		chainParams.ScriptHashAddrID = 0xc4
		chainParams.PrivateKeyID = 0xf1
		return chainParams, nil
	default:
		return nil, fmt.Errorf("unsupported chain ID: %s", chainId)
	}
}

// parseSerializedTransaction parses a base64 encoded serialized raw transaction
// and returns the unsigned transaction as a wire.MsgTx.
func parseSerializedTransaction(serializedTxn string) (*wire.MsgTx, error) {
	txBytes, err := hex.DecodeString(serializedTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex transaction: %v", err)
	}

	return parseRawTransaction(txBytes)
}

// parseRawTransaction parses a raw transaction bytes into a wire.MsgTx.
func parseRawTransaction(txBytes []byte) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(wire.TxVersion)
	err := tx.Deserialize(bytes.NewReader(txBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %v", err)
	}
	return tx, nil
}

func derEncode(signature string) (string, error) {
	// Decode hex signature
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %v", err)
	}

	if len(sigBytes) != 64 {
		return "", fmt.Errorf("invalid signature length: expected 64 bytes, got %d", len(sigBytes))
	}

	// Split into r and s components (32 bytes each)
	r := new(big.Int).SetBytes(sigBytes[:32])
	s := new(big.Int).SetBytes(sigBytes[32:])

	// Get curve parameters
	curve := btcec.S256()
	halfOrder := new(big.Int).Rsh(curve.N, 1)

	// Normalize S value to be in the lower half of the curve
	if s.Cmp(halfOrder) > 0 {
		s = new(big.Int).Sub(curve.N, s)
	}

	// Convert r and s to bytes, removing leading zeros
	rBytes := r.Bytes()
	sBytes := s.Bytes()

	// Add 0x00 prefix if the highest bit is set (to ensure positive number)
	if rBytes[0]&0x80 == 0x80 {
		rBytes = append([]byte{0x00}, rBytes...)
	}
	if sBytes[0]&0x80 == 0x80 {
		sBytes = append([]byte{0x00}, sBytes...)
	}

	// Calculate lengths
	rLen := len(rBytes)
	sLen := len(sBytes)
	totalLen := rLen + sLen + 4 // 4 additional bytes for DER sequence

	// Create DER signature
	derSig := make([]byte, 0, totalLen+1)   // +1 for sighash type
	derSig = append(derSig, 0x30)           // sequence tag
	derSig = append(derSig, byte(totalLen)) // length of sequence

	// Encode R value
	derSig = append(derSig, 0x02) // integer tag
	derSig = append(derSig, byte(rLen))
	derSig = append(derSig, rBytes...)

	// Encode S value
	derSig = append(derSig, 0x02) // integer tag
	derSig = append(derSig, byte(sLen))
	derSig = append(derSig, sBytes...)

	// Add SIGHASH_ALL
	derSig = append(derSig, 0x01)

	return hex.EncodeToString(derSig), nil
}
