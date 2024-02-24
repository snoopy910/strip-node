package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	bootnode "github.com/Silent-Protocol/go-sio/bootnode"
	"github.com/Silent-Protocol/go-sio/db"
	signer "github.com/Silent-Protocol/go-sio/signer"
	signerhub "github.com/Silent-Protocol/go-sio/signerhub"
)

func main() {
	isDeploySignerHub := flag.Bool("isDeploySignerHub", LookupEnvOrBool("IS_DEPLOY_SIGNER_HUB", false), "deploy SignerHub contract")
	isAddSigner := flag.Bool("isAddsigner", LookupEnvOrBool("IS_ADD_SIGNER", false), "add signer to SignerHub")
	privateKey := flag.String("privateKey", LookupEnvOrString("PRIVATE_KEY", ""), "private key of account to execute ethereum transactions")
	isBootstrap := flag.Bool("isBootstrap", LookupEnvOrBool("IS_BOOTSTRAP", false), "is the process a signer")
	listenHost := flag.String("host", LookupEnvOrString("LISTEN_HOST", "0.0.0.0"), "The bootstrap node host listen address\n")
	port := flag.Int("port", LookupEnvOrInt("PORT", 4001), "The bootstrap node listen port")
	bootnodeURL := flag.String("bootnode", LookupEnvOrString("BOOTNODE_URL", ""), "is the process a signer")
	httpPort := flag.String("httpPort", LookupEnvOrString("HTTP_PORT", "8080"), "http API port")
	signerPublicKey := flag.String("signerPublicKey", LookupEnvOrString("SIGNER_PUBLIC_KEY", ""), "public key of the signer nodes")
	signerPrivateKey := flag.String("signerPrivateKey", LookupEnvOrString("SIGNER_PRIVATE_KEY", ""), "private key of the signer nodes")

	//specific to network
	signerHubContractAddress := flag.String("signerHubAddress", LookupEnvOrString("SIGNER_HUB_CONTRACT_ADDRESS", "0x716A4f850809d929F85BF1C589c24FB25F884674"), "address of SignerHub contract")
	rpcURL := flag.String("rpcURL", LookupEnvOrString("RPC_URL", "http://localhost:8545"), "ethereum node RPC URL")

	// redis
	redisHost := flag.String("redisHost", LookupEnvOrString("REDIS_HOST", "localhost:6379"), "redis host")
	redisDB := flag.Int("redisDB", LookupEnvOrInt("REDIS_DB", 0), "redis db")
	redisUsername := flag.String("redisUsername", LookupEnvOrString("REDIS_USERNAME", ""), "redis username")
	redisPassword := flag.String("redisPassword", LookupEnvOrString("REDIS_PASSWORD", ""), "redis password")

	defaultPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := flag.String("keyPath", LookupEnvOrString("KEY_PATH", defaultPath+"/keys"), "path to store keygen")

	flag.Parse()

	if *isDeploySignerHub {
		signerhub.DeploySignerHubContract(*rpcURL, *privateKey)
	} else if *isAddSigner {
		signerhub.AddSignerToHub(*rpcURL, *signerHubContractAddress, *privateKey, *signerPublicKey)
	} else if *isBootstrap {
		bootnode.Start(*listenHost, *port, *path)
	} else {
		db.Initialise(*redisHost, *redisDB, *redisUsername, *redisPassword)
		signer.Start(
			*signerPrivateKey,
			*signerPublicKey,
			*bootnodeURL,
			*httpPort,
			*listenHost,
			*port,
			*rpcURL,
			*signerHubContractAddress,
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
