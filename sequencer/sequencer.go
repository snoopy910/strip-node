package sequencer

import (
	"fmt"
	"log"

	intentoperatorsregistry "github.com/StripChain/strip-node/intentOperatorsRegistry"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
)

var MaximumSigners int
var RPC_URL, IntentOperatorsRegistryContractAddress, SolversRegistryContractAddress, BridgeContractAddress string
var HeliusApiKey string
var PrivateKey string
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
	redirectUrl string,
	jwtSecret string,
	sessionSecret string,
	message string,
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

	if enableOAuth {
		fmt.Println("Initializing Google OAuth")
		if redirectUrl != "" && clientId != "" && clientSecret != "" && sessionSecret != "" && jwtSecret != "" && message != "" {
			oauthInfo = initializeGoogleOauth(redirectUrl, clientId, clientSecret, sessionSecret, jwtSecret, message)
			fmt.Println("Initializing Google OAuth done")
		} else {
			panic("Missing OAuth parameters")
		}
	}

	router := mux.NewRouter()
	// oauthRouter := router.PathPrefix("/asoauth").Subrouter()
	router.Use(ValidateAccessMiddleware)

	go startHTTPServer(httpPort, router)

	<-keepAlive
}
