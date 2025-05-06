package sequencer

import (
	"os"

	intentoperatorsregistry "github.com/StripChain/strip-node/intentOperatorsRegistry"
	"github.com/StripChain/strip-node/libs"
	db "github.com/StripChain/strip-node/libs/database"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

var MaximumSigners int
var RPC_URL, IntentOperatorsRegistryContractAddress, SolversRegistryContractAddress, BridgeContractAddress string
var HeliusApiKey string
var PrivateKey string

var validatorClientManager *ValidatorClientManager

func StartSequencer(
	httpPort string,
	rpcURL string,
	intentOperatorsRegistryContractAddress string,
	solversRegistryContractAddress string,
	heliusApiKey string,
	bridgeContractAddress string,
	privateKey string,
	valClientManager *ValidatorClientManager,
) {
	if valClientManager == nil {
		logger.Sugar().Fatal("ValidatorClientManager instance is nil in StartSequencer")
	}
	validatorClientManager = valClientManager
	logger.Sugar().Info("ValidatorClientManager initialized globally within sequencer package.")

	keepAlive := make(chan string)

	HeliusApiKey = heliusApiKey
	os.Setenv("HELIUS_API_KEY", heliusApiKey)

	intents, err := db.GetIntentsWithStatus(libs.IntentStatusProcessing)
	if err != nil {
		logger.Sugar().Fatalf("Failed to get processing intents: %v", err)
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
		logger.Sugar().Fatalf("Failed to get IntentOperatorsRegistry contract instance: %v", err)
	}

	_maxSigners, err := instance.MAXIMUMSIGNERS(&bind.CallOpts{})
	if err != nil {
		logger.Sugar().Fatalf("Failed to query MAXIMUM_SIGNERS from contract: %v", err)
	}

	MaximumSigners = int(_maxSigners.Int64())

	initialiseBridge()

	go startHTTPServer(httpPort)
	go startCheckingSigner()

	<-keepAlive
}
