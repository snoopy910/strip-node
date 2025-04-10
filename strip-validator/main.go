package main

import (
	"flag"
	"log"

	intentoperatorsregistry "github.com/StripChain/strip-node/intentOperatorsRegistry"
	"github.com/StripChain/strip-node/util"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
)

var RPC_URL, IntentOperatorsRegistryContractAddress, SolversRegistryContractAddress, BridgeContractAddress, NodePrivateKey, NodePublicKey string
var MaximumSigners int
var HeliusApiKey string
var h host.Host

type PartyProcess struct {
	Party  *tss.Party
	Exists bool
}

var partyProcesses = make(map[string]PartyProcess)

func main() {

	listenHost := flag.String("host", util.LookupEnvOrString("LISTEN_HOST", "0.0.0.0"), "The bootstrap node host listen address\n")
	port := flag.Int("port", util.LookupEnvOrInt("PORT", 4001), "The bootstrap node listen port")
	bootnodeURL := flag.String("bootnode", util.LookupEnvOrString("BOOTNODE_URL", ""), "is the process a signer")
	httpPort := flag.String("httpPort", util.LookupEnvOrString("HTTP_PORT", "8080"), "http API port")
	validatorPublicKey := flag.String("validatorPublicKey", util.LookupEnvOrString("VALIDATOR_PUBLIC_KEY", ""), "public key of the validator nodes")
	validatorPrivateKey := flag.String("validatorPrivateKey", util.LookupEnvOrString("VALIDATOR_PRIVATE_KEY", ""), "private key of the validator nodes")

	intentOperatorsRegistryContractAddress := flag.String("intentOperatorsRegistryAddress", util.LookupEnvOrString("SIGNER_HUB_CONTRACT_ADDRESS", "0x716A4f850809d929F85BF1C589c24FB25F884674"), "address of IntentOperatorsRegistry contract")
	solversRegistryContractAddress := flag.String("solversRegistryAddress", util.LookupEnvOrString("SOLVERS_REGISTRY_CONTRACT_ADDRESS", "0x56A9bCddF533Af1859842074B46B0daD07b7686a"), "address of SolversRegistry contract")
	bridgeContractAddress := flag.String("bridgeContractAddress", util.LookupEnvOrString("BRIDGE_CONTRACT_ADDRESS", "0x79E3A2B39e77dfB5C9C6a370D4a8a4fa42c482c0"), "address of Bridge contract")
	rpcURL := flag.String("rpcURL", util.LookupEnvOrString("RPC_URL", "http://localhost:8545"), "ethereum node RPC URL")
	heliusApiKey := flag.String("heliusApiKey", util.LookupEnvOrString("HELIUS_API_KEY", "6ccb4a2e-a0e6-4af3-afd0-1e06e1439547"), "helius API key")
	// maximumSigners := flag.Int("maximumSigners", util.LookupEnvOrInt("MAXIMUM_SIGNERS", 3), "maximum number of signers for an account")

	postgresHost := flag.String("postgresHost", util.LookupEnvOrString("POSTGRES_HOST", "localhost:5432"), "postgres host")
	postgresDB := flag.String("postgresDB", util.LookupEnvOrString("POSTGRES_DB", "postgres"), "postgres db name")
	postgresUser := flag.String("postgresUser", util.LookupEnvOrString("POSTGRES_USER", "postgres"), "postgres user")
	postgresPassword := flag.String("postgresPassword", util.LookupEnvOrString("POSTGRES_PASSWORD", "password"), "postgres password")

	flag.Parse()

	SolversRegistryContractAddress = *solversRegistryContractAddress
	NodePrivateKey = *validatorPrivateKey
	BridgeContractAddress = *bridgeContractAddress

	NodePublicKey = *validatorPublicKey
	HeliusApiKey = *heliusApiKey

	IntentOperatorsRegistryContractAddress = *intentOperatorsRegistryContractAddress

	RPC_URL = *rpcURL
	instance, err := intentoperatorsregistry.GetIntentOperatorsRegistryContract(RPC_URL, IntentOperatorsRegistryContractAddress)
	if err != nil {
		panic(err)
	}

	_maxSigners, err := instance.MAXIMUMSIGNERS(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}

	if err := logger.Init(); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	MaximumSigners = int(_maxSigners.Int64())

	InitialiseDB(*postgresHost, *postgresDB, *postgresUser, *postgresPassword)

	// Initialize host first
	var addr multiaddr.Multiaddr
	h, addr, err = createHost(*listenHost, *port, *bootnodeURL)
	if err != nil {
		panic(err)
	}

	// Start HTTP server after host is initialized
	go startHTTPServer(*httpPort)

	go discoverPeers(h, []multiaddr.Multiaddr{addr})
	err = subscribe(h)
	if err != nil {
		panic(err)
	}
}
