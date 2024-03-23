package signer

import (
	"log"

	intentoperatorsregistry "github.com/Silent-Protocol/go-sio/intentOperatorsRegistry"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getIntentOperatorsRegistryContract(
	rpcURL string,
	intentOperatorsRegistryContractAddress string,
) *intentoperatorsregistry.IntentOperatorsRegistry {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	instance, err := intentoperatorsregistry.NewIntentOperatorsRegistry(ethCommon.HexToAddress(intentOperatorsRegistryContractAddress), client)

	if err != nil {
		log.Fatal(err)
	}

	return instance
}
