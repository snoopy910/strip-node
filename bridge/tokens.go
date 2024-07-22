package bridge

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/StripChain/strip-node/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TokenExists(rpcURL string, bridgeContractAddress string, chainId string, srcToken string) (bool, string, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return false, "", err
	}

	instance, err := NewBridge(common.HexToAddress(bridgeContractAddress), client)
	if err != nil {
		return false, "", err
	}

	peggedToken, err := instance.PeggedTokens(&bind.CallOpts{}, chainId, srcToken)

	if err != nil {
		return false, "", err
	}

	if peggedToken != common.HexToAddress(util.ZERO_ADDRESS) {
		return true, peggedToken.Hex(), nil
	}

	return false, "", nil
}

func AddToken(rpcURL string, bridgeContractAddress string, privKey string, chainId string, srcToken string, peggedToken string) error {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return err
	}

	instance, err := NewBridge(common.HexToAddress(bridgeContractAddress), client)
	if err != nil {
		return err
	}

	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
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

	fmt.Println("Adding token...", bridgeContractAddress, chainId, srcToken, peggedToken)

	tx, err := instance.AddToken(auth, chainId, srcToken, common.HexToAddress(peggedToken))
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Token added successfully")

	return nil
}
