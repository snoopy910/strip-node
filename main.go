package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	bootnode "github.com/Silent-Protocol/go-sio/bootnode"
	signer "github.com/Silent-Protocol/go-sio/signer"
	signerhub "github.com/Silent-Protocol/go-sio/signerhub"
	solana "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

func main() {
	isSigningMessageGenerate := flag.Bool("isSigningMessageGenerate", LookupEnvOrBool("IS_SIGNING_MESSAGE_GENERATE", false), "generate signature")
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
	verifyHash := flag.Bool("verifyHash", LookupEnvOrBool("VERIFY_HASH", false), "URL of the paymaster service")

	//specific to network
	signerHubContractAddresses := flag.String("signerHubAddress", LookupEnvOrString("SIGNER_HUB_CONTRACT_ADDRESSES", "0x716A4f850809d929F85BF1C589c24FB25F884674"), "address of SignerHub contract")
	paymasterURLs := flag.String("paymasterURLs", LookupEnvOrString("PAYMASTER_URLS", "http://localhost:80"), "URL of the paymaster service")
	networkIds := flag.String("networkIds", LookupEnvOrString("NETWORK_IDS", "1"), "ethereum network id")
	rpcURLs := flag.String("rpcURLs", LookupEnvOrString("RPC_URLS", "http://localhost:8545"), "ethereum node URL")

	defaultPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := flag.String("keyPath", LookupEnvOrString("KEY_PATH", defaultPath+"/keys"), "path to store keygen")

	flag.Parse()

	if *isDeploySignerHub {
		networks := strings.Split(*networkIds, ",")
		rpcs := strings.Split(*rpcURLs, ",")
		for i := 0; i < len(networks); i++ {
			fmt.Println("Network: ", networks[i])
			signerhub.DeploySignerHubContract(rpcs[i], *privateKey)
		}
	} else if *isSigningMessageGenerate {
		c := rpc.New("https://api.devnet.solana.com")
		// pubKeyByte, _ := base58.Decode("4dEgqPG9FtjCiwW1HUdReraBozwq6qUCcDFXD8BnUn9Z")
		accountFrom := solana.MustPublicKeyFromBase58("4dEgqPG9FtjCiwW1HUdReraBozwq6qUCcDFXD8BnUn9Z")
		accountTo := solana.MustPublicKeyFromBase58("5oNDL3swdJJF1g9DzJiZ4ynHXgszjAEpUkxVYejchzrY")
		amount := uint64(1)

		// _, err := c.RequestAirdrop(context.Background(), accountFrom, solana.LAMPORTS_PER_SOL, rpc.CommitmentConfirmed)

		recentHash, err := c.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
		if err != nil {
			panic(err)
		}

		tx, err := solana.NewTransaction(
			[]solana.Instruction{
				system.NewTransferInstruction(
					amount,
					accountFrom,
					accountTo,
				).Build(),
			},
			solana.MustHashFromBase58(recentHash.Value.Blockhash.String()),
			solana.TransactionPayer(accountFrom),
		)

		if err != nil {
			panic(err)
		}

		msg, err := tx.Message.MarshalBinary()
		if err != nil {
			panic(err)
		}

		//0100010335db7213e45b7498a259b23793399c2821d56df35856f436d87c932ae396034c474f7335d5399e496566fcfe89b06ddf2f9df31fab601aafbf9afe5574a596ad000000000000000000000000000000000000000000000000000000000000000028ef4fdffb91cdfc80b5c397a086cd81d6a171dfd0ade48fd448243d8bc3686801020200010c020000000100000000000000
		// fmt.Println(hex.EncodeToString(msg))

		//[1 0 1 3 53 219 114 19 228 91 116 152 162 89 178 55 147 57 156 40 33 213 109 243 88 86 244 54 216 124 147 42 227 150 3 76 71 79 115 53 213 57 158 73 101 102 252 254 137 176 109 223 47 157 243 31 171 96 26 175 191 154 254 85 116 165 150 173 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 200 196 8 222 236 199 174 150 78 150 17 5 177 231 191 139 183 168 69 6
		fmt.Println(msg)

		// use it as original message
		// msgBigInt := (&big.Int{}).SetBytes(msg)
	} else if *isAddSigner {
		rpcs := strings.Split(*rpcURLs, ",")
		networks := strings.Split(*networkIds, ",")
		contractAddresses := strings.Split(*signerHubContractAddresses, ",")
		for i := 0; i < len(networks); i++ {
			fmt.Println("Network: ", networks[i])
			signerhub.AddSignerToHub(rpcs[i], contractAddresses[i], *privateKey, *signerPublicKey)
		}
	} else if *isBootstrap {
		bootnode.Start(*listenHost, *port, *path)
	} else {
		signer.Start(
			*signerPrivateKey,
			*signerPublicKey,
			*bootnodeURL,
			*path,
			*httpPort,
			*listenHost,
			*port,
			*privateKey,
			*verifyHash,
			*networkIds,
			*rpcURLs,
			*signerHubContractAddresses,
			*paymasterURLs,
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
