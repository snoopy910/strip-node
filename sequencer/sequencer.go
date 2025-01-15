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
var CLIENT_ID, CLIENT_SECRET, JWT_SECRET, SESSION_SECRET, REDIRECT_URL string
var oauthInfo *OAuthParameters

func StartSequencer(
	httpPort string,
	rpcURL string,
	intentOperatorsRegistryContractAddress string,
	solversRegistryContractAddress string,
	heliusApiKey string,
	bridgeContractAddress string,
	privateKey string,
	clientId string,
	clientSecret string,
	jwtSecret string,
	redirectUrl string,
	sessionSecret string,
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

	instance := intentoperatorsregistry.GetIntentOperatorsRegistryContract(rpcURL, intentOperatorsRegistryContractAddress)
	_maxSigners, err := instance.MAXIMUMSIGNERS(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}

	MaximumSigners = int(_maxSigners.Int64())

	initiaiseBridge()

	CLIENT_ID = clientId
	CLIENT_SECRET = clientSecret
	JWT_SECRET = jwtSecret
	REDIRECT_URL = redirectUrl

	oauthInfo = InitializeGoogleOauth(REDIRECT_URL, CLIENT_ID, CLIENT_SECRET, SESSION_SECRET)

	go startHTTPServer(httpPort)

	<-keepAlive
}
