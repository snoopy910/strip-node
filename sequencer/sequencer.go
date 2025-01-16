package sequencer

import (
	"fmt"
	"log"

	intentoperatorsregistry "github.com/StripChain/strip-node/intentOperatorsRegistry"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

var MaximumSigners int
var RPC_URL, IntentOperatorsRegistryContractAddress, SolversRegistryContractAddress, BridgeContractAddress string
var HeliusApiKey string
var PrivateKey string
var JWT_SECRET string
var oauthInfo *OAuthParameters

func StartSequencer(
	httpPort string,
	rpcURL string,
	intentOperatorsRegistryContractAddress string,
	solversRegistryContractAddress string,
	heliusApiKey string,
	bridgeContractAddress string,
	privateKey string,
	enableOAuth bool,
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

	fmt.Println("Bridge initialized")

	fmt.Println("enableOAuth")
	fmt.Println(enableOAuth)
	fmt.Println("clientId")
	fmt.Println(clientId)
	fmt.Println(rpcURL)
	if enableOAuth {
		fmt.Println("Initializing Google OAuth")
		// check if != ""
		fmt.Println(clientId)
		fmt.Println(clientSecret)
		fmt.Println(sessionSecret)
		fmt.Println(redirectUrl)
		// oauthinfo secret
		JWT_SECRET = jwtSecret
		oauthInfo = initializeGoogleOauth(redirectUrl, clientId, clientSecret, sessionSecret)
		fmt.Println("Initializing Google OAuth done")
	}

	go startHTTPServer(httpPort)

	<-keepAlive
}
