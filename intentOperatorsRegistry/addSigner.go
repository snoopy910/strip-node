package intentoperatorsregistry

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	tssCommon "github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func AddSignerToHub(rpcURL string, contractAddress string, privKey string, validatorPublicKey string, signerNodeURL string) {
	time.Sleep(5 * time.Second)

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
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

	instance, err := NewIntentOperatorsRegistry(common.HexToAddress(contractAddress), client)

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

	toAddress := common.HexToAddress(contractAddress)

	abi, err := IntentOperatorsRegistryMetaData.GetAbi()
	if err != nil {
		log.Fatal(err)
	}

	data, err := abi.Pack("addSigner", tssCommon.PublicKeyStrToBytes32(validatorPublicKey), signerNodeURL)
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

	tx, err := instance.AddSigner(auth, tssCommon.PublicKeyStrToBytes32(validatorPublicKey), signerNodeURL)
	if err != nil {
		log.Fatal(err)
	}

	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	fmt.Println("Signer", validatorPublicKey, " added. Transaction hash: ", tx.Hash().String())
}
