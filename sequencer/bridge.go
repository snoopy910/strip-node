package sequencer

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/StripChain/strip-node/bridge"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func initiaiseBridge() {

	// Generate bridge accounts
	// Configure SC
	// Topup bridge EVM account on L2

	// intents won't be signed using this identity for bridge operations
	// this identity is just used to identity the bridge accounts
	identity := BridgeContractAddress
	identityCurve := "ecdsa"

	_createWallet := false

	fmt.Println("Creating bridge wallet", identity, identityCurve)

	_, err := GetWallet(identity, identityCurve)
	if err != nil {
		if err.Error() == "pg: no rows in result set" {
			_createWallet = true
		} else {
			fmt.Println("Panic")
			panic(err)
		}
	}

	if !_createWallet {
		fmt.Println("wallet already exists")
		return
	}

	err = createWallet(identity, identityCurve)
	if err != nil {
		panic(err)
	}

	fmt.Println("Bridge wallet created")

	wallet, err := GetWallet(identity, identityCurve)
	if err != nil {
		panic(err)
	}

	fmt.Println("Bridge authority is: ", wallet.ECDSAPublicKey)

	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	instance, err := bridge.NewBridge(common.HexToAddress(BridgeContractAddress), client)
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

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasPrice = gasPrice
	auth.GasLimit = 972978

	nonce, err = client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))

	tx, err := instance.SetAuthority(auth, common.HexToAddress(wallet.ECDSAPublicKey))
	if err != nil {
		log.Fatal(err)
	}

	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bridge authority set")
}

func mintBridge(amount string, account string, token string, signature string) (string, error) {
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		return "", err
	}

	instance, err := bridge.NewBridge(common.HexToAddress(BridgeContractAddress), client)
	if err != nil {
		return "", err
	}

	amountBigInt, _ := big.NewInt(0).SetString(amount, 10)
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", err
	}

	nonce, err := instance.MintNonces(&bind.CallOpts{}, common.HexToAddress(account))
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

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Value = big.NewInt(0) // in wei
	auth.GasPrice = gasPrice
	auth.GasLimit = 972978

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	txnNonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	auth.Nonce = big.NewInt(int64(txnNonce))

	tx, err := instance.Mint(
		auth,
		amountBigInt,
		common.HexToAddress(account),
		common.HexToAddress(token),
		nonce,
		signatureBytes,
	)

	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
