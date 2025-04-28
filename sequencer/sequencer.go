package sequencer

import (
	"log"
	"os"

	intentoperatorsregistry "github.com/StripChain/strip-node/intentOperatorsRegistry"
	"github.com/StripChain/strip-node/libs"
	db "github.com/StripChain/strip-node/libs/database"
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
	os.Setenv("HELIUS_API_KEY", heliusApiKey)

	intents, err := db.GetIntentsWithStatus(libs.IntentStatusProcessing)
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
	go startCheckingSigner()

	<-keepAlive
}
