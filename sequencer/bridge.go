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

	amountBigInt, _ := new(big.Int).SetString(amount, 10)
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

	// Log signature details
	logger.Sugar().Debugw("Signature bytes length", "length", len(signatureBytes))

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

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	abi, err := bridge.BridgeMetaData.GetAbi()
	if err != nil {
		return "", fmt.Errorf("failed to get ABI: %v", err)
	}

	_amountIn, _ := new(big.Int).SetString(amountIn, 10)

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

	amountBigInt, _ := new(big.Int).SetString(amount, 10)
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

	amountBigInt, _ := new(big.Int).SetString(amount, 10)

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

	amountBigInt, _ := new(big.Int).SetString(amount, 10)

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
	recentHash, err := c.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
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
