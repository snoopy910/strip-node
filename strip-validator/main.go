package main

import (
	"flag"
	"log"

	intentoperatorsregistry "github.com/StripChain/strip-node/intentOperatorsRegistry"
	"github.com/StripChain/strip-node/libs/blockchains"
	"github.com/StripChain/strip-node/util"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/multiformats/go-multiaddr"
)

var RPC_URL, IntentOperatorsRegistryContractAddress, SolversRegistryContractAddress, BridgeContractAddress, NodePrivateKey, NodePublicKey string
var MaximumSigners int
var HeliusApiKey string
var SequencerHost string

type PartyProcess struct {
	Party  *tss.Party
	Exists bool
}

var partyProcesses = make(map[string]PartyProcess)

func main() {

	listenHost := flag.String("host", util.LookupEnvOrString("LISTEN_HOST", "0.0.0.0"), "The bootstrap node host listen address\n")
	port := flag.Int("port", util.LookupEnvOrInt("PORT", 4001), "The bootstrap node listen port")
	bootnodeURL := flag.String("bootnode", util.LookupEnvOrString("BOOTNODE_URL", ""), "is the process a signer")
	grpcPort := flag.String("grpcPort", util.LookupEnvOrString("GRPC_PORT", "50051"), "grpc API port")
	validatorPublicKey := flag.String("validatorPublicKey", util.LookupEnvOrString("VALIDATOR_PUBLIC_KEY", ""), "public key of the validator nodes")
	validatorPrivateKey := flag.String("validatorPrivateKey", util.LookupEnvOrString("VALIDATOR_PRIVATE_KEY", ""), "private key of the validator nodes")
	// serverCertARN := flag.String("server-cert-arn", util.LookupEnvOrString("SERVER_CERT_ARN", ""), "ARN of the gRPC server certificate in Secrets Manager")
	// serverKeyARN := flag.String("server-key-arn", util.LookupEnvOrString("SERVER_KEY_ARN", ""), "ARN of the gRPC server private key in Secrets Manager")
	// clientCaARN := flag.String("client-ca-arn", util.LookupEnvOrString("CLIENT_CA_ARN", ""), "ARN of the client CA certificate for gRPC mTLS in Secrets Manager")

	intentOperatorsRegistryContractAddress := flag.String("intentOperatorsRegistryAddress", util.LookupEnvOrString("SIGNER_HUB_CONTRACT_ADDRESS", "0x716A4f850809d929F85BF1C589c24FB25F884674"), "address of IntentOperatorsRegistry contract")
	solversRegistryContractAddress := flag.String("solversRegistryAddress", util.LookupEnvOrString("SOLVERS_REGISTRY_CONTRACT_ADDRESS", "0x56A9bCddF533Af1859842074B46B0daD07b7686a"), "address of SolversRegistry contract")
	bridgeContractAddress := flag.String("bridgeContractAddress", util.LookupEnvOrString("BRIDGE_CONTRACT_ADDRESS", "0x79E3A2B39e77dfB5C9C6a370D4a8a4fa42c482c0"), "address of Bridge contract")
	rpcURL := flag.String("rpcURL", util.LookupEnvOrString("RPC_URL", "http://localhost:8545"), "ethereum node RPC URL")
	heliusApiKey := flag.String("heliusApiKey", util.LookupEnvOrString("HELIUS_API_KEY", "6ccb4a2e-a0e6-4af3-afd0-1e06e1439547"), "helius API key")
	sequencerHost := flag.String("sequencerHost", util.LookupEnvOrString("SEQUENCER_HOST", "http://sequencer:8082"), "sequencer url")
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
	SequencerHost = *sequencerHost

	IntentOperatorsRegistryContractAddress = *intentOperatorsRegistryContractAddress

	RPC_URL = *rpcURL

	if err := logger.Init(); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	// awsCfg, err := config.LoadDefaultConfig(context.TODO())
	// if err != nil {
	// 	logger.Sugar().Fatalf("Failed to load AWS config: %v", err)
	// }
	// smClient := secretsmanager.NewFromConfig(awsCfg)

	// logger.Sugar().Info("Fetching gRPC mTLS credentials from AWS Secrets Manager...")
	// serverCertPEM, err := libs.FetchSecret(context.TODO(), smClient, *serverCertARN)
	// if err != nil {
	// 	logger.Sugar().Fatalf("Failed to fetch server certificate: %v", err)
	// }
	// serverKeyPEM, err := libs.FetchSecret(context.TODO(), smClient, *serverKeyARN)
	// if err != nil {
	// 	logger.Sugar().Fatalf("Failed to fetch server key: %v", err)
	// }
	// clientCaPEM, err := libs.FetchSecret(context.TODO(), smClient, *clientCaARN)
	// if err != nil {
	// 	logger.Sugar().Fatalf("Failed to fetch client CA certificate: %v", err)
	// }
	// logger.Sugar().Info("Successfully fetched gRPC mTLS credentials.")

	instance, err := intentoperatorsregistry.GetIntentOperatorsRegistryContract(RPC_URL, IntentOperatorsRegistryContractAddress)
	if err != nil {
		logger.Sugar().Fatalf("Failed to get IntentOperatorsRegistry contract instance: %v", err)
	}

	_maxSigners, err := instance.MAXIMUMSIGNERS(&bind.CallOpts{})
	if err != nil {
		logger.Sugar().Fatalf("Failed to query MAXIMUM_SIGNERS from contract: %v", err)
	}

	if err := logger.Init(); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	MaximumSigners = int(_maxSigners.Int64())

	InitialiseDB(*postgresHost, *postgresDB, *postgresUser, *postgresPassword)

	blockchains.InitBlockchainRegistry()
	// Initialize host first
	var addr multiaddr.Multiaddr
	h, addr, err := createHost(*listenHost, *port, *bootnodeURL)
	if err != nil {
		logger.Sugar().Fatalf("Failed to create libp2p host: %v", err)
	}

	discoverPeers(h, []multiaddr.Multiaddr{addr})

	// Start HTTP server after host is initialized
	// go startHTTPServer(*httpPort)
	go startGRPCServer(*grpcPort, h, "", "", "")

	err = subscribe(h)
	if err != nil {
		logger.Sugar().Fatalf("Failed to subscribe to libp2p topic: %v", err)
	}
}
