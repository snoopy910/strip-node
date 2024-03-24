package sequencer

import (
	"log"

	intentoperatorsregistry "github.com/Silent-Protocol/go-sio/intentOperatorsRegistry"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

var MaximumSigners int

func StartSequencer(
	httpPort string,
	rpcURL string,
	intentOperatorsRegistryContractAddress string,
) {
	keepAlive := make(chan string)

	intents, err := GetIntents(INTENT_STATUS_PROCESSING)
	if err != nil {
		log.Fatal(err)
	}

	for _, intent := range intents {
		go ProcessIntent(intent.ID)
	}

	instance := intentoperatorsregistry.GetIntentOperatorsRegistryContract(rpcURL, intentOperatorsRegistryContractAddress)
	_maxSigners, err := instance.MAXIMUMSIGNERS(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}

	MaximumSigners = int(_maxSigners.Int64())

	go startHTTPServer(httpPort)

	<-keepAlive
}
