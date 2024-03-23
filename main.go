package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	bootnode "github.com/Silent-Protocol/go-sio/bootnode"
	intentoperatorsregistry "github.com/Silent-Protocol/go-sio/intentOperatorsRegistry"
	"github.com/Silent-Protocol/go-sio/sequencer"
	signer "github.com/Silent-Protocol/go-sio/signer"
)

func main() {
	isSolanaTest := flag.Bool("isSolanaTest", LookupEnvOrBool("IS_SOLANA_TEST", false), "is the process a signer")

	isDeployIntentOperatorsRegistry := flag.Bool("isDeployIntentOperatorsRegistry", LookupEnvOrBool("IS_DEPLOY_SIGNER_HUB", false), "deploy IntentOperatorsRegistry contract")
	isAddSigner := flag.Bool("isAddsigner", LookupEnvOrBool("IS_ADD_SIGNER", false), "add signer to IntentOperatorsRegistry")
	privateKey := flag.String("privateKey", LookupEnvOrString("PRIVATE_KEY", ""), "private key of account to execute ethereum transactions")
	isBootstrap := flag.Bool("isBootstrap", LookupEnvOrBool("IS_BOOTSTRAP", false), "is the process a signer")
	isSequencer := flag.Bool("isSequencer", LookupEnvOrBool("IS_SEQUENCER", false), "is the process a sequencer")
	listenHost := flag.String("host", LookupEnvOrString("LISTEN_HOST", "0.0.0.0"), "The bootstrap node host listen address\n")
	port := flag.Int("port", LookupEnvOrInt("PORT", 4001), "The bootstrap node listen port")
	bootnodeURL := flag.String("bootnode", LookupEnvOrString("BOOTNODE_URL", ""), "is the process a signer")
	httpPort := flag.String("httpPort", LookupEnvOrString("HTTP_PORT", "8080"), "http API port")
	signerPublicKey := flag.String("signerPublicKey", LookupEnvOrString("SIGNER_PUBLIC_KEY", ""), "public key of the signer nodes")
	signerPrivateKey := flag.String("signerPrivateKey", LookupEnvOrString("SIGNER_PRIVATE_KEY", ""), "private key of the signer nodes")

	//specific to network
	intentOperatorsRegistryContractAddress := flag.String("intentOperatorsRegistryAddress", LookupEnvOrString("SIGNER_HUB_CONTRACT_ADDRESS", "0x716A4f850809d929F85BF1C589c24FB25F884674"), "address of IntentOperatorsRegistry contract")
	rpcURL := flag.String("rpcURL", LookupEnvOrString("RPC_URL", "http://localhost:8545"), "ethereum node RPC URL")
	maximumSigners := flag.Int("maximumSigners", LookupEnvOrInt("MAXIMUM_SIGNERS", 3), "maximum number of signers for an account")

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

	if *isDeployIntentOperatorsRegistry {
		intentoperatorsregistry.DeployIntentOperatorsRegistryContract(*rpcURL, *privateKey)
	} else if *isAddSigner {
		intentoperatorsregistry.AddSignerToHub(*rpcURL, *intentOperatorsRegistryContractAddress, *privateKey, *signerPublicKey)
	} else if *isBootstrap {
		bootnode.Start(*listenHost, *port, *path)
	} else if *isSequencer {
		sequencer.InitialiseDB(*postgresHost, *postgresDB, *postgresUser, *postgresPassword)
		sequencer.StartSequencer(*httpPort)
	} else if *isSolanaTest {
		sequencer.TestBuildSolana()
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
			*maximumSigners,
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
