package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/StripChain/strip-node/ERC20"
	bootnode "github.com/StripChain/strip-node/bootnode"
	"github.com/StripChain/strip-node/bridge"
	"github.com/StripChain/strip-node/bridgeTokenMock"
	intentoperatorsregistry "github.com/StripChain/strip-node/intentOperatorsRegistry"
	"github.com/StripChain/strip-node/sequencer"
	signer "github.com/StripChain/strip-node/signer"
	"github.com/StripChain/strip-node/solver"
	solversregistry "github.com/StripChain/strip-node/solversRegistry"
	"github.com/StripChain/strip-node/util/logger"
)

func main() {
	isSolanaTest := flag.Bool("isSolanaTest", LookupEnvOrBool("IS_SOLANA_TEST", false), "is the process a signer")
	isEthereumTest := flag.Bool("isEthereumTest", LookupEnvOrBool("IS_SOLANA_TEST", false), "is the process a signer")

	isDeployIntentOperatorsRegistry := flag.Bool("isDeployIntentOperatorsRegistry", LookupEnvOrBool("IS_DEPLOY_SIGNER_HUB", false), "deploy IntentOperatorsRegistry contract")
	isDeploySolversRegistry := flag.Bool("isDeploySolversRegistry", LookupEnvOrBool("IS_DEPLOY_SOLVERS_REGISTRY", false), "deploy SolversRegistry contract")
	isDeployBridge := flag.Bool("isDeployBridge", LookupEnvOrBool("IS_DEPLOY_BRIDGE", false), "deploy Bridge contract")
	isDeployBridgeToken := flag.Bool("isDeployBridgeToken", LookupEnvOrBool("IS_DEPLOY_BRIDGE_TOKEN", false), "deploy BridgeToken contract")
	isDeployERC20Token := flag.Bool("isDeployERC20Token", LookupEnvOrBool("IS_DEPLOY_ERC20_TOKEN", false), "deploy ERC20 token contract")
	isAddSigner := flag.Bool("isAddsigner", LookupEnvOrBool("IS_ADD_SIGNER", false), "add signer to IntentOperatorsRegistry")
	isAddSolver := flag.Bool("isAddSolver", LookupEnvOrBool("IS_ADD_SOLVER", false), "add solver to SolversRegistry")
	isAddToken := flag.Bool("isAddToken", LookupEnvOrBool("IS_ADD_TOKEN", false), "add token to Bridge")
	isSetSwapRouter := flag.Bool("isSetSwapRouter", LookupEnvOrBool("IS_SET_SWAP_ROUTER", false), "set swap router in Bridge")
	privateKey := flag.String("privateKey", LookupEnvOrString("PRIVATE_KEY", ""), "private key of account to execute ethereum transactions")
	isBootstrap := flag.Bool("isBootstrap", LookupEnvOrBool("IS_BOOTSTRAP", false), "is the process a signer")
	isSequencer := flag.Bool("isSequencer", LookupEnvOrBool("IS_SEQUENCER", false), "is the process a sequencer")
	isTestSolver := flag.Bool("isTestSolver", LookupEnvOrBool("IS_TEST_SOLVER", false), "is the process a solver")
	listenHost := flag.String("host", LookupEnvOrString("LISTEN_HOST", "0.0.0.0"), "The bootstrap node host listen address\n")
	port := flag.Int("port", LookupEnvOrInt("PORT", 4001), "The bootstrap node listen port")
	bootnodeURL := flag.String("bootnode", LookupEnvOrString("BOOTNODE_URL", ""), "is the process a signer")
	httpPort := flag.String("httpPort", LookupEnvOrString("HTTP_PORT", "8080"), "http API port")
	signerPublicKey := flag.String("signerPublicKey", LookupEnvOrString("SIGNER_PUBLIC_KEY", ""), "public key of the signer nodes")
	signerPrivateKey := flag.String("signerPrivateKey", LookupEnvOrString("SIGNER_PRIVATE_KEY", ""), "private key of the signer nodes")
	signerNodeURL := flag.String("signerNodeURL", LookupEnvOrString("SIGNER_NODE_URL", ""), "URL of the signer node")
	solverDomain := flag.String("solverDomain", LookupEnvOrString("SOLVER_DOMAIN", ""), "domain of the solver")
	heliusApiKey := flag.String("heliusApiKey", LookupEnvOrString("HELIUS_API_KEY", "6ccb4a2e-a0e6-4af3-afd0-1e06e1439547"), "helius API key")

	intentOperatorsRegistryContractAddress := flag.String("intentOperatorsRegistryAddress", LookupEnvOrString("SIGNER_HUB_CONTRACT_ADDRESS", "0x716A4f850809d929F85BF1C589c24FB25F884674"), "address of IntentOperatorsRegistry contract")
	solversRegistryContractAddress := flag.String("solversRegistryAddress", LookupEnvOrString("SOLVERS_REGISTRY_CONTRACT_ADDRESS", "0x56A9bCddF533Af1859842074B46B0daD07b7686a"), "address of SolversRegistry contract")
	bridgeContractAddress := flag.String("bridgeContractAddress", LookupEnvOrString("BRIDGE_CONTRACT_ADDRESS", "0x79E3A2B39e77dfB5C9C6a370D4a8a4fa42c482c0"), "address of Bridge contract")
	rpcURL := flag.String("rpcURL", LookupEnvOrString("RPC_URL", "http://localhost:8545"), "ethereum node RPC URL")
	maximumSigners := flag.Int("maximumSigners", LookupEnvOrInt("MAXIMUM_SIGNERS", 3), "maximum number of signers for an account")
	tokenName := flag.String("tokenName", LookupEnvOrString("TOKEN_NAME", "Strip"), "name of the token")
	tokenSymbol := flag.String("tokenSymbol", LookupEnvOrString("TOKEN_SYMBOL", "STRP"), "symbol of the token")
	tokenDecimals := flag.Int("tokenDecimals", LookupEnvOrInt("TOKEN_DECIMALS", 18), "decimals of the token")
	chainId := flag.String("chainId", LookupEnvOrString("CHAIN_ID", "1337"), "chain id of the token")
	tokenAddress := flag.String("tokenAddress", LookupEnvOrString("TOKEN_ADDRESS", "0x0000000000000000000000000000000000000000"), "address of the token")
	peggedToken := flag.String("peggedToken", LookupEnvOrString("PEGGED_TOKEN", ""), "address of the pegged token")
	swapRouter := flag.String("swapRouter", LookupEnvOrString("SWAP_ROUTER", "0x3466c635Bdf084DA32CD5bc16c00C1CA1A459011"), "address of the swap router")

	// postgres
	postgresHost := flag.String("postgresHost", LookupEnvOrString("POSTGRES_HOST", "localhost:5432"), "postgres host")
	postgresDB := flag.String("postgresDB", LookupEnvOrString("POSTGRES_DB", "postgres"), "postgres db name")
	postgresUser := flag.String("postgresUser", LookupEnvOrString("POSTGRES_USER", "postgres"), "postgres user")
	postgresPassword := flag.String("postgresPassword", LookupEnvOrString("POSTGRES_PASSWORD", "password"), "postgres password")

	defaultPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := flag.String("keyPath", LookupEnvOrString("KEY_PATH", defaultPath+"/keys"), "path to store keygen")

	flag.Parse()

	if err := logger.Init(); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	if *isDeployIntentOperatorsRegistry {
		intentoperatorsregistry.DeployIntentOperatorsRegistryContract(*rpcURL, *privateKey)
	} else if *isDeploySolversRegistry {
		solversregistry.DeploySolversRegistryContract(*rpcURL, *privateKey)
	} else if *isAddSigner {
		intentoperatorsregistry.AddSignerToHub(*rpcURL, *intentOperatorsRegistryContractAddress, *privateKey, *signerPublicKey, *signerNodeURL)
	} else if *isAddSolver {
		solversregistry.AddSolver(*rpcURL, *solversRegistryContractAddress, *privateKey, *solverDomain)
	} else if *isAddToken {
		bridge.AddToken(*rpcURL, *bridgeContractAddress, *privateKey, *chainId, *tokenAddress, *peggedToken)
	} else if *isDeployBridge {
		bridge.DeployBridgeContract(*rpcURL, *privateKey)
	} else if *isDeployBridgeToken {
		bridgeTokenMock.DeployBridgeTokenContract(*rpcURL, *privateKey, *tokenName, *tokenSymbol, uint(*tokenDecimals), *bridgeContractAddress)
	} else if *isDeployERC20Token {
		ERC20.DeployERC20Token(*rpcURL, *privateKey, *tokenName, *tokenSymbol, uint(*tokenDecimals))
	} else if *isSetSwapRouter {
		bridge.SetSwapRouter(*rpcURL, *privateKey, *bridgeContractAddress, *swapRouter)
	} else if *isBootstrap {
		bootnode.Start(*listenHost, *port, *path)
	} else if *isSequencer {
		sequencer.InitialiseDB(*postgresHost, *postgresDB, *postgresUser, *postgresPassword)
		sequencer.StartSequencer(
			*httpPort,
			*rpcURL,
			*intentOperatorsRegistryContractAddress,
			*solversRegistryContractAddress,
			*heliusApiKey,
			*bridgeContractAddress,
			*privateKey,
		)
	} else if *isTestSolver {
		solver.StartTestSolver(*httpPort)
	} else if *isSolanaTest {
		sequencer.GetSolanaTransfers(
			"901",
			"243hStsqpngr2Dv4ktE9wCW6CZTbFYaRiXo1QAK1PTtyLZBL2xz17XTAuud3HN8YmpYhdSRJmP3Rx3pMHdu6Pxqi",
			*heliusApiKey,
		)
		// identity.VerifySignature(
		// 	"GScvaHyfG3NMNm8AdPjjZt3xRiNtAwHy5z5yY1oaQA4Q",
		// 	"eddsa",
		// 	"The quick brown fox jumps over the lazy dog",
		// 	"3XdzeBWMBCAhuTZ7237A6GZRW2N9ge5LjszPBycvkFdaspSwN8vn5kMN8cW9dvqJtur34Wef5rdn6uMFyUXBwsVJ",
		// )

		// identity.VerifySignature(
		// 	"0x805B25e9246e1D80c399f05C4B515a0C22097457",
		// 	"ecdsa",
		// 	"The quick brown fox jumps over the lazy dog",
		// 	"0x3835f5c4f8ccf5ab0c1b3d827ca72dd1953409fa971ca68c6f6dda7905acf59f1411f1e1d40495bf224c4744dcf9d63032ca96e5685f435b4781aabf685fd88a1c",
		// )

		// sequencer.TestBuildSolana()
	} else if *isEthereumTest {
		sequencer.GetEthereumTransfers("1", "0x1e96c4f5dc65ba33b4ea2a50e350f119d133d2b4c9f36ac79152198382a16375", "0x06Cd69B61900B426499ef0319Fae5CEC2acca4DE")
	} else {
		signer.InitialiseDB(*postgresHost, *postgresDB, *postgresUser, *postgresPassword)
		signer.Start(
			*signerPrivateKey,
			*signerPublicKey,
			*bootnodeURL,
			*httpPort,
			*listenHost,
			*port,
			*rpcURL,
			*intentOperatorsRegistryContractAddress,
			*solversRegistryContractAddress,
			*maximumSigners,
			*bridgeContractAddress,
		)
	}
}

func LookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("LookupEnvOrInt[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func LookupEnvOrBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.ParseBool(val)
		if err != nil {
			log.Fatalf("LookupEnvOrInt[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}
