package sequencer

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/StripChain/strip-node/ERC20"
	"github.com/StripChain/strip-node/bridge"
	tssCommon "github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/libs/blockchains"
	db "github.com/StripChain/strip-node/libs/database"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
)

func initialiseBridge() {

	// Generate bridge accounts
	// Configure SC
	// Topup bridge EVM account on L2

	// intents won't be signed using this identity for bridge operations
	// this identity is just used to identity the bridge accounts
	identity := BridgeContractAddress
	blockchainID := blockchains.Ethereum

	_createWallet := false

	logger.Sugar().Infow("Creating bridge wallet", "identity", identity, "blockchainID", blockchainID)

	_, err := db.GetWallet(identity, blockchainID)
	if err != nil {
		if err.Error() == "pg: no rows in result set" {
			_createWallet = true
		} else {
			logger.Sugar().Errorw("failed to get wallet", "error", err)
			panic(err)
		}
	}

	if !_createWallet {
		logger.Sugar().Info("wallet already exists")
		return
	}

	err = createWallet(identity, blockchainID)
	if err != nil {
		logger.Sugar().Errorw("failed to create wallet", "error", err)
		panic(err)
	}

	logger.Sugar().Info("Bridge wallet created")

	wallet, err := db.GetWallet(identity, blockchainID)
	if err != nil {
		logger.Sugar().Errorw("failed to get wallet", "error", err)
		panic(err)
	}

	logger.Sugar().Infow("Bridge authority is: ", "authority", wallet.EthereumPublicKey)

	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		logger.Sugar().Errorw("failed to dial ethclient", "error", err)
		panic(err)
	}

	privateKey, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		logger.Sugar().Errorw("failed to convert private key to ECDSA", "error", err)
		panic(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		logger.Sugar().Errorw("error casting public key to ECDSA")
		panic("error casting public key to ECDSA")
	}

	toAddress := common.HexToAddress(BridgeContractAddress)

	instance, err := bridge.NewBridge(toAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	abi, err := bridge.BridgeMetaData.GetAbi()
	if err != nil {
		log.Fatal(err)
	}

	data, err := abi.Pack("setAuthority", common.HexToAddress(wallet.EthereumPublicKey))
	if err != nil {
		log.Fatal(err)
	}

	gas, err := tssCommon.EstimateTransactionGas(fromAddress, &toAddress, 0, gasPrice, nil, nil, data, client, 1.2)
	if err != nil {
		log.Fatalf("failed to estimate gas: %v", err)
	}
	logger.Sugar().Infof("gas estimate %d", gas)

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasPrice = gasPrice
	auth.GasLimit = gas
	// auth.GasLimit = 972978

	nonce, err = client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))

	tx, err := instance.SetAuthority(auth, common.HexToAddress(wallet.EthereumPublicKey))
	if err != nil {
		log.Fatal(err)
	}

	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}

	logger.Sugar().Info("Bridge authority set")
}

func mintBridge(amount string, account string, token string, signature string) (string, error) {
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		return "", err
	}

	toAddress := common.HexToAddress(BridgeContractAddress)

	instance, err := bridge.NewBridge(toAddress, client)
	if err != nil {
		return "", err
	}

	amountBigInt := new(big.Int)
	_, ok := amountBigInt.SetString(amount, 10)
	if !ok {
		return "", fmt.Errorf("invalid amount: %s", amount)
	}
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", err
	}

	nonce, err := instance.Nonces(&bind.CallOpts{}, common.HexToAddress(account))
	if err != nil {
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	abi, err := bridge.BridgeMetaData.GetAbi()
	if err != nil {
		return "", err
	}

	ethSigHex := hexutil.Encode(signatureBytes[:])
	recoveryParam := ethSigHex[len(ethSigHex)-2:]
	ethSigHex = ethSigHex[:len(ethSigHex)-2]

	if recoveryParam == "00" {
		ethSigHex = ethSigHex + "1b"
	} else {
		ethSigHex = ethSigHex + "1c"
	}

	ethSigHex = strings.Replace(ethSigHex, "0x", "", -1)

	ethSigHexBytes, err := hex.DecodeString(ethSigHex)
	if err != nil {
		return "", err
	}

	data, err := abi.Pack("mint", amountBigInt, common.HexToAddress(token), common.HexToAddress(account), nonce, ethSigHexBytes)
	if err != nil {
		return "", err
	}

	gas, err := tssCommon.EstimateTransactionGas(fromAddress, &toAddress, 0, gasPrice, nil, nil, data, client, 1.2)
	if err != nil {
		return "", fmt.Errorf("failed to estimate gas: %v", err)
	}
	logger.Sugar().Infof("gas estimate %d", gas)

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Value = big.NewInt(0) // in wei
	auth.GasPrice = gasPrice
	auth.GasLimit = gas
	// auth.GasLimit = 972978

	txnNonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	auth.Nonce = big.NewInt(int64(txnNonce))

	tx, err := instance.Mint(
		auth,
		amountBigInt,
		common.HexToAddress(token),
		common.HexToAddress(account),
		nonce,
		ethSigHexBytes,
	)

	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func swapBridge(
	account string,
	tokenIn string,
	tokenOut string,
	amountIn string,
	deadline int64,
	signature string,
) (string, error) {
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		return "", err
	}

	// Validate input parameters
	if tokenIn == "" {
		return "", fmt.Errorf("tokenIn cannot be empty")
	}

	// If tokenOut is empty, we need to handle it
	if tokenOut == "" {
		return "", fmt.Errorf("tokenOut cannot be empty for swap operation")
	}

	// Check if tokenIn and tokenOut are the same
	if strings.EqualFold(tokenIn, tokenOut) {
		return "", fmt.Errorf("tokenIn and tokenOut cannot be the same: %s", tokenIn)
	}

	toAddress := common.HexToAddress(BridgeContractAddress)

	// Log the bridge address for debugging
	logger.Sugar().Infow("Swap details",
		"bridgeAddress", BridgeContractAddress,
		"account", account,
		"tokenIn", tokenIn,
		"tokenOut", tokenOut,
		"amountIn", amountIn)

	instance, err := bridge.NewBridge(toAddress, client)
	if err != nil {
		return "", err
	}

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", fmt.Errorf("failed to decode signature: %v", err)
	}

	// Enhanced logging for signature debugging
	logger.Sugar().Infow("Signature details",
		"rawSignature", signature,
		"length", len(signatureBytes),
		"hexEncoded", hex.EncodeToString(signatureBytes))

	nonce, err := instance.Nonces(&bind.CallOpts{}, common.HexToAddress(account))
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	logger.Sugar().Infow("Account nonce", "account", account, "nonce", nonce)

	privateKey, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Log chain ID for debugging
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		logger.Sugar().Warnw("Failed to get chain ID", "error", err)
	} else {
		logger.Sugar().Infow("Chain ID", "chainID", chainID)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	abi, err := bridge.BridgeMetaData.GetAbi()
	if err != nil {
		return "", fmt.Errorf("failed to get ABI: %v", err)
	}

	_amountIn := new(big.Int)
	_, ok = _amountIn.SetString(amountIn, 10)
	if !ok {
		return "", fmt.Errorf("invalid amount: %s", amountIn)
	}

	params := bridge.ISwapRouterExactInputSingleParams{
		AmountIn:          _amountIn,
		AmountOutMinimum:  big.NewInt(0),
		TokenIn:           common.HexToAddress(tokenIn),
		TokenOut:          common.HexToAddress(tokenOut),
		Fee:               big.NewInt(500),
		Recipient:         common.HexToAddress(account),
		Deadline:          big.NewInt(0).SetInt64(deadline),
		SqrtPriceLimitX96: big.NewInt(0),
	}

	logger.Sugar().Infow("Swap parameters",
		"tokenIn", params.TokenIn.Hex(),
		"tokenOut", params.TokenOut.Hex(),
		"amountIn", params.AmountIn.String(),
		"fee", params.Fee.String(),
		"recipient", params.Recipient.Hex(),
		"deadline", params.Deadline.String())

	// Get message hash directly from contract for comparison
	messageHash, err := instance.GetSwapMessageHash(
		&bind.CallOpts{},
		params,
		common.HexToAddress(account),
		nonce,
	)
	if err == nil {
		logger.Sugar().Infow("Contract message hash",
			"messageHash", "0x"+hex.EncodeToString(messageHash[:]))

		// Dump raw parameters for debugging
		logger.Sugar().Infow("Message hash parameters",
			"tokenIn", params.TokenIn.Hex(),
			"tokenOut", params.TokenOut.Hex(),
			"fee", params.Fee.String(),
			"recipient", params.Recipient.Hex(),
			"deadline", params.Deadline.String(),
			"amountIn", params.AmountIn.String(),
			"amountOutMinimum", params.AmountOutMinimum.String(),
			"sqrtPriceLimitX96", params.SqrtPriceLimitX96.String(),
			"nonce", nonce.String(),
			"account", account,
			"chainID", chainID)

		ethSignedMessageHash, err := instance.GetEthSignedMessageHash(
			&bind.CallOpts{},
			messageHash,
		)
		if err == nil {
			logger.Sugar().Infow("Contract eth signed message hash",
				"ethSignedMessageHash", "0x"+hex.EncodeToString(ethSignedMessageHash[:]))
		}
	}

	ethSigHex := hexutil.Encode(signatureBytes[:])
	logger.Sugar().Infow("Signature transformation - Initial", "ethSigHex", ethSigHex)

	recoveryParam := ethSigHex[len(ethSigHex)-2:]
	ethSigHex = ethSigHex[:len(ethSigHex)-2]

	logger.Sugar().Infow("Signature transformation",
		"recoveryParam", recoveryParam,
		"ethSigHexBeforeV", ethSigHex)

	if recoveryParam == "00" {
		ethSigHex = ethSigHex + "1b"
	} else {
		ethSigHex = ethSigHex + "1c"
	}
	logger.Sugar().Infow("Signature transformation - After V", "ethSigHex", ethSigHex)

	ethSigHex = strings.Replace(ethSigHex, "0x", "", -1)
	logger.Sugar().Infow("Signature transformation - Final", "ethSigHex", ethSigHex)

	ethSigHexBytes, err := hex.DecodeString(ethSigHex)
	if err != nil {
		return "", err
	}
	logger.Sugar().Infow("Signature bytes for contract",
		"length", len(ethSigHexBytes),
		"hexEncoded", hex.EncodeToString(ethSigHexBytes))

	data, err := abi.Pack("swap", params, common.HexToAddress(account), nonce, ethSigHexBytes)
	if err != nil {
		return "", err
	}

	gas, err := tssCommon.EstimateTransactionGas(fromAddress, &toAddress, 0, gasPrice, nil, nil, data, client, 1.2)
	if err != nil {
		return "", fmt.Errorf("failed to estimate gas: %v", err)
	}
	logger.Sugar().Infof("gas estimate %d", gas)

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Value = big.NewInt(0) // in wei
	auth.GasPrice = gasPrice
	auth.GasLimit = gas
	// auth.GasLimit = 972978

	txnNonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	auth.Nonce = big.NewInt(int64(txnNonce))

	logger.Sugar().Infow("Calling swap function on contract",
		"fromAddress", fromAddress.Hex(),
		"toAddress", toAddress.Hex(),
		"paramLength", len(ethSigHexBytes),
		"txnNonce", txnNonce)

	tx, err := instance.Swap(
		auth,
		params,
		common.HexToAddress(account),
		nonce,
		ethSigHexBytes,
	)

	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func swapMultiplePoolsBridge(
	account string,
	tokenIn string,
	path string, // hex-encoded path instead of tokenOut
	amountIn string,
	deadline int64,
	signature string,
) (string, error) {
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		return "", err
	}

	// Validate input parameters
	if tokenIn == "" {
		return "", fmt.Errorf("tokenIn cannot be empty")
	}

	if path == "" {
		return "", fmt.Errorf("path cannot be empty for multi-pool swap operation")
	}

	// Strip "0x" prefix if present before hex decoding
	pathForDecoding := path
	if strings.HasPrefix(pathForDecoding, "0x") {
		pathForDecoding = pathForDecoding[2:]
	}

	// Decode path from hex
	pathBytes, err := hex.DecodeString(pathForDecoding)
	if err != nil {
		return "", fmt.Errorf("invalid path format: %v", err)
	}

	// Extract tokenOut for logging (last 20 bytes of path)
	var tokenOutForLogging string
	if len(pathBytes) >= 40 {
		tokenOutForLogging = "0x" + hex.EncodeToString(pathBytes[len(pathBytes)-20:])
	} else {
		return "", fmt.Errorf("path too short, cannot extract tokenOut")
	}

	// Extract first token from path and ensure it matches tokenIn
	if len(pathBytes) >= 20 {
		firstTokenInPath := "0x" + hex.EncodeToString(pathBytes[:20])
		logger.Sugar().Infow("First token in path", "firstTokenInPath", firstTokenInPath, "providedTokenIn", tokenIn)

		// Convert both to lowercase for case-insensitive comparison
		if !strings.EqualFold(firstTokenInPath, tokenIn) {
			logger.Sugar().Errorw("TokenIn parameter doesn't match first token in path",
				"tokenIn", tokenIn,
				"firstTokenInPath", firstTokenInPath)
			return "", fmt.Errorf("tokenIn parameter (%s) doesn't match first token in path (%s)",
				tokenIn, firstTokenInPath)
		}
	}

	toAddress := common.HexToAddress(BridgeContractAddress)

	// Log the bridge address for debugging
	logger.Sugar().Infow("Swap multiple pools details",
		"bridgeAddress", BridgeContractAddress,
		"account", account,
		"tokenIn", tokenIn,
		"pathLength", len(pathBytes),
		"path", path[:min(len(path), 40)]+"...", // Log only beginning for brevity
		"extractedTokenOut", tokenOutForLogging,
		"amountIn", amountIn)

	instance, err := bridge.NewBridge(toAddress, client)
	if err != nil {
		return "", err
	}

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", fmt.Errorf("failed to decode signature: %v", err)
	}

	// Enhanced logging for signature debugging
	logger.Sugar().Infow("Signature details",
		"rawSignature", signature,
		"length", len(signatureBytes),
		"hexEncoded", hex.EncodeToString(signatureBytes))

	nonce, err := instance.Nonces(&bind.CallOpts{}, common.HexToAddress(account))
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	logger.Sugar().Infow("Account nonce", "account", account, "nonce", nonce)

	privateKey, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Log chain ID for debugging
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		logger.Sugar().Warnw("Failed to get chain ID", "error", err)
	} else {
		logger.Sugar().Infow("Chain ID", "chainID", chainID)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	abi, err := bridge.BridgeMetaData.GetAbi()
	if err != nil {
		return "", fmt.Errorf("failed to get ABI: %v", err)
	}

	_amountIn := new(big.Int)
	_, ok = _amountIn.SetString(amountIn, 10)
	if !ok {
		return "", fmt.Errorf("invalid amount: %s", amountIn)
	}

	// Create ExactInputParams for multiple pools
	params := bridge.ISwapRouterExactInputParams{
		Path:             pathBytes,
		Recipient:        common.HexToAddress(account),
		Deadline:         big.NewInt(0).SetInt64(deadline),
		AmountIn:         _amountIn,
		AmountOutMinimum: big.NewInt(0),
	}

	// Check token balance of bridge contract for the input token
	tokenContract, err := ERC20.NewERC20(common.HexToAddress(tokenIn), client)
	if err != nil {
		logger.Sugar().Warnw("Failed to instantiate token contract", "error", err, "tokenIn", tokenIn)
	} else {
		// Check bridge contract's balance
		bridgeBalance, err := tokenContract.BalanceOf(&bind.CallOpts{}, toAddress)
		if err != nil {
			logger.Sugar().Warnw("Failed to get bridge balance", "error", err)
		} else {
			logger.Sugar().Infow("Bridge contract token balance",
				"tokenIn", tokenIn,
				"bridgeBalance", bridgeBalance.String(),
				"requiredAmount", _amountIn.String())

			if bridgeBalance.Cmp(_amountIn) < 0 {
				logger.Sugar().Errorw("Bridge contract has insufficient token balance",
					"tokenIn", tokenIn,
					"bridgeBalance", bridgeBalance.String(),
					"requiredAmount", _amountIn.String())
			}
		}

		// Check if bridge contract has allowance to spend user's tokens
		// In the contract, it uses sweep which directly transfers tokens
		// But let's check allowance anyway for debugging
		allowance, err := tokenContract.Allowance(&bind.CallOpts{},
			common.HexToAddress(account), toAddress)
		if err != nil {
			logger.Sugar().Warnw("Failed to get token allowance", "error", err)
		} else {
			logger.Sugar().Infow("Token allowance for bridge contract",
				"tokenIn", tokenIn,
				"account", account,
				"allowance", allowance.String(),
				"requiredAmount", _amountIn.String())

			if allowance.Cmp(_amountIn) < 0 {
				logger.Sugar().Warnw("Bridge contract may have insufficient allowance",
					"tokenIn", tokenIn,
					"allowance", allowance.String(),
					"requiredAmount", _amountIn.String())
			}
		}

		// Check if the account has sufficient token balance
		accountBalance, err := tokenContract.BalanceOf(&bind.CallOpts{}, common.HexToAddress(account))
		if err != nil {
			logger.Sugar().Warnw("Failed to get account balance", "error", err)
		} else {
			logger.Sugar().Infow("Account token balance",
				"tokenIn", tokenIn,
				"account", account,
				"balance", accountBalance.String(),
				"requiredAmount", _amountIn.String())

			if accountBalance.Cmp(_amountIn) < 0 {
				logger.Sugar().Errorw("Account has insufficient token balance",
					"tokenIn", tokenIn,
					"accountBalance", accountBalance.String(),
					"requiredAmount", _amountIn.String())
			}
		}
	}

	// Check if bridge contract has router approval
	swapRouterAddress, err := instance.SwapRouter(&bind.CallOpts{})
	if err != nil {
		logger.Sugar().Warnw("Failed to get swap router address", "error", err)
	} else {
		logger.Sugar().Infow("Swap router address", "swapRouter", swapRouterAddress.Hex())

		// Check if bridge has approval for the swap router
		if tokenContract != nil {
			routerAllowance, err := tokenContract.Allowance(&bind.CallOpts{},
				toAddress, swapRouterAddress)
			if err != nil {
				logger.Sugar().Warnw("Failed to get router allowance", "error", err)
			} else {
				logger.Sugar().Infow("Router allowance from bridge",
					"tokenIn", tokenIn,
					"swapRouter", swapRouterAddress.Hex(),
					"allowance", routerAllowance.String(),
					"requiredAmount", _amountIn.String())

				if routerAllowance.Cmp(_amountIn) < 0 {
					logger.Sugar().Warnw("Swap router may have insufficient allowance",
						"tokenIn", tokenIn,
						"allowance", routerAllowance.String(),
						"requiredAmount", _amountIn.String())
				}
			}
		}
	}

	logger.Sugar().Infow("Swap multiple pools parameters",
		"pathLength", len(params.Path),
		"recipient", params.Recipient.Hex(),
		"deadline", params.Deadline.String(),
		"amountIn", params.AmountIn.String(),
		"amountOutMinimum", params.AmountOutMinimum.String())

	// Get message hash directly from contract for comparison
	messageHash, err := instance.GetSwapMultiplePoolsMessageHash(
		&bind.CallOpts{},
		params,
		common.HexToAddress(tokenIn),
		common.HexToAddress(account),
		nonce,
	)
	if err == nil {
		logger.Sugar().Infow("Contract message hash for multiple pools",
			"messageHash", "0x"+hex.EncodeToString(messageHash[:]))

		// Dump raw parameters for debugging
		logger.Sugar().Infow("Multiple pools message hash parameters",
			"tokenIn", tokenIn,
			"pathLength", len(params.Path),
			"recipient", params.Recipient.Hex(),
			"deadline", params.Deadline.String(),
			"amountIn", params.AmountIn.String(),
			"amountOutMinimum", params.AmountOutMinimum.String(),
			"nonce", nonce.String(),
			"account", account,
			"chainID", chainID)

		ethSignedMessageHash, err := instance.GetEthSignedMessageHash(
			&bind.CallOpts{},
			messageHash,
		)
		if err == nil {
			logger.Sugar().Infow("Contract eth signed message hash for multiple pools",
				"ethSignedMessageHash", "0x"+hex.EncodeToString(ethSignedMessageHash[:]))
		}
	}

	ethSigHex := hexutil.Encode(signatureBytes[:])
	logger.Sugar().Infow("Signature transformation - Initial", "ethSigHex", ethSigHex)

	recoveryParam := ethSigHex[len(ethSigHex)-2:]
	ethSigHex = ethSigHex[:len(ethSigHex)-2]

	logger.Sugar().Infow("Signature transformation",
		"recoveryParam", recoveryParam,
		"ethSigHexBeforeV", ethSigHex)

	if recoveryParam == "00" {
		ethSigHex = ethSigHex + "1b"
	} else {
		ethSigHex = ethSigHex + "1c"
	}
	logger.Sugar().Infow("Signature transformation - After V", "ethSigHex", ethSigHex)

	ethSigHex = strings.Replace(ethSigHex, "0x", "", -1)
	logger.Sugar().Infow("Signature transformation - Final", "ethSigHex", ethSigHex)

	ethSigHexBytes, err := hex.DecodeString(ethSigHex)
	if err != nil {
		return "", err
	}
	logger.Sugar().Infow("Signature bytes for contract",
		"length", len(ethSigHexBytes),
		"hexEncoded", hex.EncodeToString(ethSigHexBytes))

	// Pack the parameters for swapMultiplePools function
	data, err := abi.Pack("swapMultiplePools", params, common.HexToAddress(tokenIn), common.HexToAddress(account), nonce, ethSigHexBytes)
	if err != nil {
		return "", err
	}

	gas, err := tssCommon.EstimateTransactionGas(fromAddress, &toAddress, 0, gasPrice, nil, nil, data, client, 1.2)
	if err != nil {
		return "", fmt.Errorf("failed to estimate gas: %v", err)
	}
	logger.Sugar().Infof("gas estimate %d", gas)

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Value = big.NewInt(0) // in wei
	auth.GasPrice = gasPrice
	auth.GasLimit = gas

	txnNonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	auth.Nonce = big.NewInt(int64(txnNonce))

	logger.Sugar().Infow("Calling swapMultiplePools function on contract",
		"fromAddress", fromAddress.Hex(),
		"toAddress", toAddress.Hex(),
		"paramLength", len(ethSigHexBytes),
		"nonce", nonce)

	// Call swapMultiplePools on the contract
	tx, err := instance.SwapMultiplePools(
		auth,
		params,
		common.HexToAddress(tokenIn),
		common.HexToAddress(account),
		nonce,
		ethSigHexBytes,
	)

	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func burnTokens(
	account string,
	amount string,
	token string,
	signature string,
) (string, error) {
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		return "", err
	}

	toAddress := common.HexToAddress(BridgeContractAddress)

	instance, err := bridge.NewBridge(toAddress, client)
	if err != nil {
		return "", err
	}

	amountBigInt := new(big.Int)
	_, ok := amountBigInt.SetString(amount, 10)
	if !ok {
		return "", fmt.Errorf("invalid amount: %s", amount)
	}
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", err
	}

	nonce, err := instance.Nonces(&bind.CallOpts{}, common.HexToAddress(account))
	if err != nil {
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	abi, err := bridge.BridgeMetaData.GetAbi()
	if err != nil {
		return "", err
	}

	ethSigHex := hexutil.Encode(signatureBytes[:])
	recoveryParam := ethSigHex[len(ethSigHex)-2:]
	ethSigHex = ethSigHex[:len(ethSigHex)-2]

	if recoveryParam == "00" {
		ethSigHex = ethSigHex + "1b"
	} else {
		ethSigHex = ethSigHex + "1c"
	}

	ethSigHex = strings.Replace(ethSigHex, "0x", "", -1)

	ethSigHexBytes, err := hex.DecodeString(ethSigHex)
	if err != nil {
		return "", err
	}

	data, err := abi.Pack("burn", common.HexToAddress(account), amountBigInt, common.HexToAddress(token), nonce, ethSigHexBytes)
	if err != nil {
		return "", err
	}

	gas, err := tssCommon.EstimateTransactionGas(fromAddress, &toAddress, 0, gasPrice, nil, nil, data, client, 1.2)
	if err != nil {
		return "", fmt.Errorf("failed to estimate gas: %v", err)
	}
	logger.Sugar().Infof("gas estimate %d", gas)

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Value = big.NewInt(0) // in wei
	auth.GasPrice = gasPrice
	auth.GasLimit = gas
	// auth.GasLimit = 972978

	txnNonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	auth.Nonce = big.NewInt(int64(txnNonce))

	tx, err := instance.Burn(
		auth,
		common.HexToAddress(account),
		amountBigInt,
		common.HexToAddress(token),
		nonce,
		ethSigHexBytes,
	)

	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func withdrawEVMNativeGetSignature(
	rpcURL string,
	account string,
	amount string,
	recipient string,
	chainId string,
) (string, *types.Transaction, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", nil, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(account))
	if err != nil {
		return "", nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", nil, err
	}

	gasLimit := uint64(60000)

	amountBigInt := new(big.Int)
	_, ok := amountBigInt.SetString(amount, 10)
	if !ok {
		return "", nil, fmt.Errorf("invalid amount: %s", amount)
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(recipient), amountBigInt, gasLimit, gasPrice, nil)
	chainIdBigInt, _ := new(big.Int).SetString(chainId, 10)
	signer := types.NewEIP155Signer(chainIdBigInt)
	txHash := signer.Hash(tx)

	return hex.EncodeToString(txHash.Bytes()), tx, nil
}

func withdrawEVMTxn(
	rpcURL string,
	signature string,
	tx *types.Transaction,
	chainId string,
) (string, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", err
	}

	chainIdBigInt, _ := new(big.Int).SetString(chainId, 10)
	signer := types.NewEIP155Signer(chainIdBigInt)

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", err
	}

	signedTx, err := tx.WithSignature(signer, signatureBytes)
	if err != nil {
		return "", err
	}

	signedTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		return "", err
	}

	logger.Sugar().Infof("Signed transaction: 0x%x", signedTxBytes)

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

func withdrawERC20GetSignature(
	rpcURL string,
	account string,
	amount string,
	recipient string,
	chainId string,
	token string,
) (string, *types.Transaction, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", nil, err
	}

	const erc20ABI = `[{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(account))
	if err != nil {
		return "", nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", nil, err
	}

	gasLimit := uint64(60000)

	amountBigInt := new(big.Int)
	_, ok := amountBigInt.SetString(amount, 10)
	if !ok {
		return "", nil, fmt.Errorf("invalid amount: %s", amount)
	}

	parsedABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return "", nil, err
	}

	data, err := parsedABI.Pack("transfer", common.HexToAddress(recipient), amountBigInt)
	if err != nil {
		return "", nil, err
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(token), big.NewInt(0), gasLimit, gasPrice, data)
	chainIdBigInt, _ := new(big.Int).SetString(chainId, 10)
	signer := types.NewEIP155Signer(chainIdBigInt)
	txHash := signer.Hash(tx)

	return hex.EncodeToString(txHash.Bytes()), tx, nil
}

func withdrawSolanaNativeGetSignature(
	rpcURL string,
	account string,
	amount string,
	recipient string,
) (string, string, error) {
	accountFrom := solana.MustPublicKeyFromBase58(account)
	accountTo := solana.MustPublicKeyFromBase58(recipient)

	// convert amount to uint64
	_amount, _ := big.NewInt(0).SetString(amount, 10)
	amountUint64 := _amount.Uint64()

	c := rpc.New(rpcURL)
	recentHash, err := c.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return "", "", err
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				amountUint64,
				accountFrom,
				accountTo,
			).Build(),
		},
		recentHash.Value.Blockhash,
		solana.TransactionPayer(accountFrom),
	)

	if err != nil {
		return "", "", err
	}

	_msg, err := tx.ToBase64()
	if err != nil {
		return "", "", err
	}

	_msgBytes, _ := base64.StdEncoding.DecodeString(_msg)
	_msgBase58 := base58.Encode(_msgBytes)

	msg, err := tx.Message.MarshalBinary()
	if err != nil {
		return "", "", err
	}

	return _msgBase58, base58.Encode(msg), nil
}

func withdrawSolanaSPLGetSignature(
	rpcURL string,
	account string,
	amount string,
	recipient string,
	tokenAddr string,
) (string, string, error) {
	accountFrom := solana.MustPublicKeyFromBase58(account)
	accountTo := solana.MustPublicKeyFromBase58(recipient)
	tokenMint := solana.MustPublicKeyFromBase58(tokenAddr)

	senderTokenAccount, _, err := solana.FindAssociatedTokenAddress(accountFrom, tokenMint)
	if err != nil {
		return "", "", fmt.Errorf("failed to get sender token account: %v", err)
	}

	recipientTokenAccount, _, err := solana.FindAssociatedTokenAddress(accountTo, tokenMint)
	if err != nil {
		return "", "", fmt.Errorf("failed to get recipient token account: %v", err)
	}

	// convert amount to uint64
	_amount, _ := big.NewInt(0).SetString(amount, 10)
	amountUint64 := _amount.Uint64()

	c := rpc.New(rpcURL)
	recentHash, err := c.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return "", "", err
	}

	transferInstruction := token.NewTransferInstruction(
		amountUint64,
		senderTokenAccount,
		recipientTokenAccount,
		accountFrom,
		nil, // No multisig signers
	).Build()

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			transferInstruction,
		},
		recentHash.Value.Blockhash,
		solana.TransactionPayer(accountFrom),
	)

	if err != nil {
		return "", "", err
	}

	_msg, err := tx.ToBase64()
	if err != nil {
		return "", "", err
	}

	_msgBytes, _ := base64.StdEncoding.DecodeString(_msg)
	_msgBase58 := base58.Encode(_msgBytes)

	msg, err := tx.Message.MarshalBinary()
	if err != nil {
		return "", "", err
	}

	return _msgBase58, base58.Encode(msg), nil
}

func withdrawSolanaTxn(
	rpcURL string,
	transaction string,
	signature string,
) (string, error) {
	c := rpc.New(rpcURL)

	decodedTransactionData, err := base58.Decode(transaction)
	if err != nil {
		return "", err
	}

	_tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(decodedTransactionData))
	if err != nil {
		return "", err
	}

	sig, _ := base58.Decode(signature)
	_signature := solana.SignatureFromBytes(sig)

	_tx.Signatures = append(_tx.Signatures, _signature)

	err = _tx.VerifySignatures()

	if err != nil {
		return "", err
	}

	hash, err := c.SendTransaction(context.Background(), _tx)

	if err != nil {
		return "", err
	}

	return hash.String(), nil
}
