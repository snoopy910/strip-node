package bridge

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/StripChain/strip-node/util/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func BridgeSwapDataToSign(
	rpcURL string,
	bridgeContractAddress string,
	account string,
	tokenIn string,
	tokenOut string,
	amountIn string,
	deadline int64,
) (string, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", err
	}

	instance, err := NewBridge(common.HexToAddress(bridgeContractAddress), client)
	if err != nil {
		return "", err
	}

	_amountIn, _ := new(big.Int).SetString(amountIn, 10)

	params := ISwapRouterExactInputSingleParams{
		AmountIn:          _amountIn,
		AmountOutMinimum:  big.NewInt(0),
		TokenIn:           common.HexToAddress(tokenIn),
		TokenOut:          common.HexToAddress(tokenOut),
		Fee:               big.NewInt(500),
		Recipient:         common.HexToAddress(account),
		Deadline:          big.NewInt(0).SetInt64(deadline),
		SqrtPriceLimitX96: big.NewInt(0),
	}

	nonce, err := instance.Nonces(&bind.CallOpts{}, common.HexToAddress(account))
	if err != nil {
		return "", err
	}

	messageHash, err := instance.GetSwapMessageHash(
		&bind.CallOpts{},
		params,
		common.HexToAddress(account),
		nonce,
	)

	if err != nil {
		return "", err
	}

	hash, err := instance.GetEthSignedMessageHash(
		&bind.CallOpts{},
		messageHash,
	)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash[:]), nil
}

// decodeHexString safely decodes a hex string with or without "0x" prefix
func decodeHexString(hexStr string) ([]byte, error) {
	// Strip "0x" prefix if present
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}
	return hex.DecodeString(hexStr)
}

// BridgeSwapMultiplePoolsDataToSign prepares the data to sign for a multi-pool swap intent.
func BridgeSwapMultiplePoolsDataToSign(
	rpcURL string,
	bridgeContractAddress string,
	account string,
	tokenIn string,
	path string, // hex-encoded path
	amountIn string,
	deadline int64,
) (string, error) {
	logger.Sugar().Infow("BridgeSwapMultiplePoolsDataToSign called",
		"bridgeContract", bridgeContractAddress,
		"account", account,
		"tokenIn", tokenIn,
		"pathLength", len(path),
		"path", path[:min(len(path), 20)]+"...", // Log only the beginning to avoid overwhelming logs
		"amountIn", amountIn,
		"deadline", deadline)

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		logger.Sugar().Errorw("Failed to connect to Ethereum client", "error", err, "rpcURL", rpcURL)
		return "", err
	}

	instance, err := NewBridge(common.HexToAddress(bridgeContractAddress), client)
	if err != nil {
		logger.Sugar().Errorw("Failed to instantiate Bridge contract", "error", err, "bridgeContract", bridgeContractAddress)
		return "", err
	}

	_amountIn, _ := new(big.Int).SetString(amountIn, 10)

	// Log before hex decoding
	logger.Sugar().Infow("Preparing to decode path hex string",
		"pathRaw", path,
		"hasPrefix", strings.HasPrefix(path, "0x"))

	// Strip "0x" prefix if present before hex decoding
	pathForDecoding := path
	if strings.HasPrefix(pathForDecoding, "0x") {
		pathForDecoding = pathForDecoding[2:]
		logger.Sugar().Infow("Removed 0x prefix from path", "newPath", pathForDecoding[:min(len(pathForDecoding), 20)]+"...")
	}

	pathBytes, err := hex.DecodeString(pathForDecoding)
	if err != nil {
		logger.Sugar().Errorw("Failed to decode path from hex string",
			"error", err,
			"pathOriginal", path,
			"pathForDecoding", pathForDecoding)
		return "", err
	}

	logger.Sugar().Infow("Successfully decoded path hex string",
		"decodedLengthBytes", len(pathBytes))

	params := ISwapRouterExactInputParams{
		Path:             pathBytes,
		Recipient:        common.HexToAddress(account),
		Deadline:         big.NewInt(0).SetInt64(deadline),
		AmountIn:         _amountIn,
		AmountOutMinimum: big.NewInt(0),
	}

	logger.Sugar().Infow("Created swap parameters",
		"recipient", params.Recipient.Hex(),
		"amountIn", params.AmountIn.String(),
		"pathLength", len(params.Path))

	nonce, err := instance.Nonces(&bind.CallOpts{}, common.HexToAddress(account))
	if err != nil {
		logger.Sugar().Errorw("Failed to get nonce", "error", err, "account", account)
		return "", err
	}

	logger.Sugar().Infow("Retrieved account nonce", "nonce", nonce.String(), "account", account)

	// Extract tokenOut from path - for Uniswap V3 path format
	// The path format is [tokenIn, fee, tokenOut] for single hop
	// or [tokenIn, fee, token1, fee, token2, ...] for multi-hop
	var tokenOutAddress string
	if len(pathBytes) >= 40 {
		// Extract the last 20 bytes as tokenOut address
		tokenOutAddress = "0x" + hex.EncodeToString(pathBytes[len(pathBytes)-20:])
		logger.Sugar().Infow("Extracted tokenOut from path", "tokenOutAddress", tokenOutAddress)
	} else {
		logger.Sugar().Errorw("Path is too short to extract tokenOut", "pathLength", len(pathBytes))
		return "", fmt.Errorf("invalid path format: too short to extract tokenOut")
	}

	// IMPORTANT: Use the SwapMultiplePoolsMessageHash to match contract's expectations
	// This is critical to ensure the hash generated here matches what the contract expects
	messageHash, err := instance.GetSwapMultiplePoolsMessageHash(
		&bind.CallOpts{},
		params,
		common.HexToAddress(tokenIn),
		common.HexToAddress(account),
		nonce,
	)
	if err != nil {
		logger.Sugar().Errorw("Failed to get swap message hash", "error", err)
		return "", err
	}

	logger.Sugar().Infow("Generated message hash", "messageHash", "0x"+hex.EncodeToString(messageHash[:]))

	// Dump complete parameters for debugging
	logger.Sugar().Infow("BridgeSwapMultiplePoolsDataToSign parameters",
		"path", hex.EncodeToString(params.Path[:min(30, len(params.Path))])+"...",
		"recipient", params.Recipient.Hex(),
		"deadline", params.Deadline.String(),
		"amountIn", params.AmountIn.String(),
		"amountOutMinimum", params.AmountOutMinimum.String(),
		"tokenIn", tokenIn,
		"tokenOut", tokenOutAddress,
		"account", account,
		"nonce", nonce.String())

	// Get chain ID for comparison
	chainID, err := client.ChainID(context.Background())
	if err == nil {
		logger.Sugar().Infow("Chain ID for swap hash", "chainID", chainID)
	}

	hash, err := instance.GetEthSignedMessageHash(
		&bind.CallOpts{},
		messageHash,
	)
	if err != nil {
		logger.Sugar().Errorw("Failed to get eth signed message hash", "error", err)
		return "", err
	}

	result := hex.EncodeToString(hash[:])
	logger.Sugar().Infow("Generated final hash for signing",
		"signatureData", result,
		"length", len(result))

	return result, nil
}

// Helper function for min - used for logging truncation
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
