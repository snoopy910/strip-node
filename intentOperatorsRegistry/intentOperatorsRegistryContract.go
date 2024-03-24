package intentoperatorsregistry

import (
	"log"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetIntentOperatorsRegistryContract(
	rpcURL string,
	intentOperatorsRegistryContractAddress string,
) *IntentOperatorsRegistry {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	instance, err := NewIntentOperatorsRegistry(ethCommon.HexToAddress(intentOperatorsRegistryContractAddress), client)

	if err != nil {
		log.Fatal(err)
	}

	return instance
}
