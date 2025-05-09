package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/StripChain/strip-node/ERC20"
	"github.com/StripChain/strip-node/bridge"
	"github.com/StripChain/strip-node/bridgeTokenMock"
	"github.com/StripChain/strip-node/evm"
	intentoperatorsregistry "github.com/StripChain/strip-node/intentOperatorsRegistry"
	"github.com/StripChain/strip-node/libs/blockchains"
	"github.com/StripChain/strip-node/libs/database"
	"github.com/StripChain/strip-node/sequencer"
	"github.com/StripChain/strip-node/solana"
	"github.com/StripChain/strip-node/solver"
	lendingsolver "github.com/StripChain/strip-node/solvers/lending_solver"
	uniswapv3solver "github.com/StripChain/strip-node/solvers/uniswap_v3_solver"
	solversregistry "github.com/StripChain/strip-node/solversRegistry"
	"github.com/StripChain/strip-node/util"
	"github.com/StripChain/strip-node/util/logger"
)

func main() {
	isSolanaTest := flag.Bool("isSolanaTest", util.LookupEnvOrBool("IS_SOLANA_TEST", false), "is the process a signer")
	isEthereumTest := flag.Bool("isEthereumTest", util.LookupEnvOrBool("IS_SOLANA_TEST", false), "is the process a signer")
	isLendingSolver := flag.Bool("isLendingSolver", util.LookupEnvOrBool("IS_LENDING_SOLVER", false), "start lending solver")
	isSwapSolver := flag.Bool("isSwapSolver", util.LookupEnvOrBool("IS_SWAP_SOLVER", false), "start swap solver")

	isDeployIntentOperatorsRegistry := flag.Bool("isDeployIntentOperatorsRegistry", util.LookupEnvOrBool("IS_DEPLOY_SIGNER_HUB", false), "deploy IntentOperatorsRegistry contract")
	isDeploySolversRegistry := flag.Bool("isDeploySolversRegistry", util.LookupEnvOrBool("IS_DEPLOY_SOLVERS_REGISTRY", false), "deploy SolversRegistry contract")
	isDeployBridge := flag.Bool("isDeployBridge", util.LookupEnvOrBool("IS_DEPLOY_BRIDGE", false), "deploy Bridge contract")
	isDeployBridgeToken := flag.Bool("isDeployBridgeToken", util.LookupEnvOrBool("IS_DEPLOY_BRIDGE_TOKEN", false), "deploy BridgeToken contract")
	isDeployERC20Token := flag.Bool("isDeployERC20Token", util.LookupEnvOrBool("IS_DEPLOY_ERC20_TOKEN", false), "deploy ERC20 token contract")
	isAddSigner := flag.Bool("isAddsigner", util.LookupEnvOrBool("IS_ADD_SIGNER", false), "add signer to IntentOperatorsRegistry")
	isAddSolver := flag.Bool("isAddSolver", util.LookupEnvOrBool("IS_ADD_SOLVER", false), "add solver to SolversRegistry")
	isAddToken := flag.Bool("isAddToken", util.LookupEnvOrBool("IS_ADD_TOKEN", false), "add token to Bridge")
	isSetSwapRouter := flag.Bool("isSetSwapRouter", util.LookupEnvOrBool("IS_SET_SWAP_ROUTER", false), "set swap router in Bridge")
	privateKey := flag.String("privateKey", util.LookupEnvOrString("PRIVATE_KEY", ""), "private key of account to execute ethereum transactions")
	isSequencer := flag.Bool("isSequencer", util.LookupEnvOrBool("IS_SEQUENCER", false), "is the process a sequencer")
	isTestSolver := flag.Bool("isTestSolver", util.LookupEnvOrBool("IS_TEST_SOLVER", false), "is the process a solver")
	httpPort := flag.String("httpPort", util.LookupEnvOrString("HTTP_PORT", "8080"), "http API port")
	validatorPublicKey := flag.String("validatorPublicKey", util.LookupEnvOrString("VALIDATOR_PUBLIC_KEY", ""), "public key of the signer nodes")
	validatorNodeURL := flag.String("validatorNodeURL", util.LookupEnvOrString("VALIDATOR_NODE_URL", ""), "URL of the signer node")
	solverDomain := flag.String("solverDomain", util.LookupEnvOrString("SOLVER_DOMAIN", ""), "domain of the solver")
	heliusApiKey := flag.String("heliusApiKey", util.LookupEnvOrString("HELIUS_API_KEY", "6ccb4a2e-a0e6-4af3-afd0-1e06e1439547"), "helius API key")

	intentOperatorsRegistryContractAddress := flag.String("intentOperatorsRegistryAddress", util.LookupEnvOrString("SIGNER_HUB_CONTRACT_ADDRESS", "0x716A4f850809d929F85BF1C589c24FB25F884674"), "address of IntentOperatorsRegistry contract")
	solversRegistryContractAddress := flag.String("solversRegistryAddress", util.LookupEnvOrString("SOLVERS_REGISTRY_CONTRACT_ADDRESS", "0x56A9bCddF533Af1859842074B46B0daD07b7686a"), "address of SolversRegistry contract")
	bridgeContractAddress := flag.String("bridgeContractAddress", util.LookupEnvOrString("BRIDGE_CONTRACT_ADDRESS", "0x79E3A2B39e77dfB5C9C6a370D4a8a4fa42c482c0"), "address of Bridge contract")
	lendingPoolAddress := flag.String("lendingPoolAddress", util.LookupEnvOrString("LENDING_POOL_ADDRESS", ""), "address of lending pool contract")
	uniswapV3FactoryAddress := flag.String("uniswapV3FactoryAddress", util.LookupEnvOrString("UNISWAP_V3_FACTORY_ADDRESS", "0x9af0e87FBA28e20863488bda3CfA012d1c7863d9"), "address of Uniswap V3 factory contract")
	npmAddress := flag.String("npmAddress", util.LookupEnvOrString("NPM_ADDRESS", "0x0c3729964A75870f9c692833A18AFE315be700e1"), "address of Uniswap V3 NonfungiblePositionManager contract")
	rpcURL := flag.String("rpcURL", util.LookupEnvOrString("RPC_URL", "http://localhost:8545"), "ethereum node RPC URL")
	tokenName := flag.String("tokenName", util.LookupEnvOrString("TOKEN_NAME", "Strip"), "name of the token")
	tokenSymbol := flag.String("tokenSymbol", util.LookupEnvOrString("TOKEN_SYMBOL", "STRP"), "symbol of the token")
	tokenDecimals := flag.Int("tokenDecimals", util.LookupEnvOrInt("TOKEN_DECIMALS", 18), "decimals of the token")
	chainId := flag.String("chainId", util.LookupEnvOrString("CHAIN_ID", "1337"), "chain id of the token")
	tokenAddress := flag.String("tokenAddress", util.LookupEnvOrString("TOKEN_ADDRESS", "0x0000000000000000000000000000000000000000"), "address of the token")
	peggedToken := flag.String("peggedToken", util.LookupEnvOrString("PEGGED_TOKEN", ""), "address of the pegged token")
	swapRouter := flag.String("swapRouter", util.LookupEnvOrString("SWAP_ROUTER", ""), "address of the swap router")

	// postgres
	postgresHost := flag.String("postgresHost", util.LookupEnvOrString("POSTGRES_HOST", "localhost:5432"), "postgres host")
	postgresDB := flag.String("postgresDB", util.LookupEnvOrString("POSTGRES_DB", "postgres"), "postgres db name")
	postgresUser := flag.String("postgresUser", util.LookupEnvOrString("POSTGRES_USER", "postgres"), "postgres user")
	postgresPassword := flag.String("postgresPassword", util.LookupEnvOrString("POSTGRES_PASSWORD", "password"), "postgres password")

	// clientCertARN := flag.String("client-cert-arn", util.LookupEnvOrString("CLIENT_CERT_ARN", ""), "ARN for sequencer client cert")
	// clientKeyARN := flag.String("client-key-arn", util.LookupEnvOrString("CLIENT_KEY_ARN", ""), "ARN for sequencer client key")
	// // *** Simplified: Assume ONE CA for all validators for this example ***
	// serverCaARN := flag.String("server-ca-arn", util.LookupEnvOrString("SERVER_CA_ARN", ""), "ARN for the CA signing validator server certs")

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
		intentoperatorsregistry.AddSignerToHub(*rpcURL, *intentOperatorsRegistryContractAddress, *privateKey, *validatorPublicKey, *validatorNodeURL)
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
	} else if *isSequencer {
		// if *clientCertARN == "" || *clientKeyARN == "" || *serverCaARN == "" {
		// 	logger.Sugar().Fatal("Client cert ARN, client key ARN, and server CA ARN must be provided via flags/env")
		// }

		// awsCfg, err := config.LoadDefaultConfig(ctx)
		// if err != nil {
		// 	logger.Sugar().Fatalf("Failed to load AWS config: %v", err)
		// }
		// smClient := secretsmanager.NewFromConfig(awsCfg)

		// logger.Sugar().Info("Fetching TLS credentials from AWS Secrets Manager...")
		// clientCertPEM, err := libs.FetchSecret(ctx, smClient, *clientCertARN)
		// if err != nil {
		// 	logger.Sugar().Fatalf("Failed to fetch client cert: %v", err)
		// }
		// clientKeyPEM, err := libs.FetchSecret(ctx, smClient, *clientKeyARN)
		// if err != nil {
		// 	logger.Sugar().Fatalf("Failed to fetch client key: %v", err)
		// }
		// serverCaPEM, err := libs.FetchSecret(ctx, smClient, *serverCaARN)
		// if err != nil {
		// 	logger.Sugar().Fatalf("Failed to fetch server CA: %v", err)
		// }
		// logger.Sugar().Info("Successfully fetched TLS credentials.")

		blockchains.InitBlockchainRegistry()

		database.InitialiseDB(*postgresHost, *postgresDB, *postgresUser, *postgresPassword)

		validatorClientManager, err := sequencer.NewValidatorClientManager(
		// clientCertPEM,
		// clientKeyPEM,
		// serverCaPEM,
		)
		if err != nil {
			logger.Sugar().Fatalf("Failed to initialize ValidatorClientManager: %v", err)
		}
		sequencer.StartSequencer(
			*httpPort,
			*rpcURL,
			*intentOperatorsRegistryContractAddress,
			*solversRegistryContractAddress,
			*heliusApiKey,
			*bridgeContractAddress,
			*privateKey,
			validatorClientManager,
		)
	} else if *isTestSolver {
		solver.StartTestSolver(*httpPort)
	} else if *isSwapSolver {
		chainID, err := strconv.ParseInt(*chainId, 10, 64)
		if err != nil {
			log.Fatal("Failed to parse chain ID:", err)
		}
		uniswapv3solver.Start(*rpcURL, *httpPort, *uniswapV3FactoryAddress, *npmAddress, chainID)
	} else if *isLendingSolver {
		chainID, err := strconv.ParseInt(*chainId, 10, 64)
		if err != nil {
			log.Fatal("Failed to parse chain ID:", err)
		}
		lendingsolver.Start(*rpcURL, *httpPort, *lendingPoolAddress, chainID)
	} else if *isSolanaTest {
		solana.GetSolanaTransfers(
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
		evm.GetEthereumTransfers("1", "0x1e96c4f5dc65ba33b4ea2a50e350f119d133d2b4c9f36ac79152198382a16375", "0x06Cd69B61900B426499ef0319Fae5CEC2acca4DE")
	}
}
