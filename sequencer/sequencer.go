package sequencer

import (
	"log"

	intentoperatorsregistry "github.com/StripChain/strip-node/intentOperatorsRegistry"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

var MaximumSigners int
var RPC_URL, IntentOperatorsRegistryContractAddress, SolversRegistryContractAddress, BridgeContractAddress string
var HeliusApiKey string
var PrivateKey string

func StartSequencer(
	httpPort string,
	rpcURL string,
	intentOperatorsRegistryContractAddress string,
	solversRegistryContractAddress string,
	heliusApiKey string,
	bridgeContractAddress string,
	privateKey string,
) {
	keepAlive := make(chan string)

	HeliusApiKey = heliusApiKey

	intents, err := GetIntentsWithStatus(INTENT_STATUS_PROCESSING)
	if err != nil {
		log.Fatal(err)
	}

	for _, intent := range intents {
		go ProcessIntent(intent.ID)
	}

	RPC_URL = rpcURL
	IntentOperatorsRegistryContractAddress = intentOperatorsRegistryContractAddress
	SolversRegistryContractAddress = solversRegistryContractAddress
	BridgeContractAddress = bridgeContractAddress
	PrivateKey = privateKey

	instance, err := intentoperatorsregistry.GetIntentOperatorsRegistryContract(rpcURL, intentOperatorsRegistryContractAddress)
	if err != nil {
		panic(err)
	}

	_maxSigners, err := instance.MAXIMUMSIGNERS(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}

	MaximumSigners = int(_maxSigners.Int64())

	initialiseBridge()

	go startHTTPServer(httpPort)

	<-keepAlive
}
