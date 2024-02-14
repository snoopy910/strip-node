package signer

import (
	"log"

	"github.com/Silent-Protocol/go-sio/signerhub"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getSignerHubContract(
	rpcURL string,
	signerHubContractAddress string,
) *signerhub.Signerhub {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	instance, err := signerhub.NewSignerhub(ethCommon.HexToAddress(signerHubContractAddress), client)

	if err != nil {
		log.Fatal(err)
	}

	return instance
}
