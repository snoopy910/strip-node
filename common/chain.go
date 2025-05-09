package common

import (
	"fmt"
	"time"
)

type Chain struct {
	ChainId     string
	GenesisId   string
	GenesisHash string
	IndexerUrl  string
	ChainApiKey string
	ChainType   string
	ChainUrl    string
	KeyCurve    string
	TokenSymbol string
	RpcUsername string
	RpcPassword string
	OpTimeout   time.Duration
}

var Chains = []Chain{
	{
		ChainId:     "2000",
		ChainType:   "dogecoin",
		ChainUrl:    "https://doge-mainnet.gateway.tatum.io/",
		KeyCurve:    "secp256k1",
		TokenSymbol: "DOGE",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "2001",
		ChainType:   "dogecoin",
		ChainUrl:    "https://doge-testnet.gateway.tatum.io/",
		KeyCurve:    "secp256k1",
		TokenSymbol: "DOGE",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "1337",
		ChainType:   "evm",
		ChainUrl:    "http://ganache:8545",
		KeyCurve:    "ecdsa",
		TokenSymbol: "ETH",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "44331",
		ChainType:   "evm",
		ChainUrl:    "https://rpc-stripsepolia-5w8r5b9f7b.t.conduit.xyz",
		KeyCurve:    "ecdsa",
		TokenSymbol: "ETH",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "421614",
		ChainType:   "evm",
		ChainUrl:    "https://arbitrum-sepolia-rpc.publicnode.com",
		KeyCurve:    "ecdsa",
		TokenSymbol: "ETH",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "11155111",
		ChainType:   "evm",
		ChainUrl:    "https://ethereum-sepolia-rpc.publicnode.com",
		KeyCurve:    "ecdsa",
		TokenSymbol: "ETH",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "901",
		ChainType:   "solana",
		ChainUrl:    "https://api.devnet.solana.com",
		KeyCurve:    "eddsa",
		TokenSymbol: "SOL",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "1",
		ChainType:   "evm",
		ChainUrl:    "https://ethereum-rpc.publicnode.com",
		KeyCurve:    "ecdsa",
		TokenSymbol: "ETH",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "900",
		ChainType:   "solana",
		ChainUrl:    "https://api.solana.com",
		KeyCurve:    "eddsa",
		TokenSymbol: "SOL",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "137",
		ChainType:   "evm",
		ChainUrl:    "https://polygon-pokt.nodies.app",
		KeyCurve:    "ecdsa",
		TokenSymbol: "MATIC",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "11",
		ChainType:   "aptos",
		ChainUrl:    "https://fullnode.mainnet.aptoslabs.com",
		KeyCurve:    "eddsa",
		TokenSymbol: "APT",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "167",
		ChainType:   "aptos",
		ChainUrl:    "https://fullnode.devnet.aptoslabs.com",
		KeyCurve:    "eddsa",
		TokenSymbol: "APT",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "1000",
		ChainType:   "bitcoin",
		ChainUrl:    "http://172.17.0.1:8332", // Local Bitcoin Core RPC, calling from docker process
		KeyCurve:    "ecdsa",
		TokenSymbol: "BTC",
		RpcUsername: "your_rpc_user",
		RpcPassword: "your_rpc_password",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "1001",
		ChainType:   "bitcoin",
		ChainUrl:    "http://172.17.0.1:18332", // Local Bitcoin Core RPC, calling from docker process
		KeyCurve:    "ecdsa",
		TokenSymbol: "BTC",
		RpcUsername: "your_rpc_user",
		RpcPassword: "your_rpc_password",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "1002",
		ChainType:   "bitcoin",
		ChainUrl:    "http://bitcoind:8332", // Local Bitcoin Core RPC, calling from docker process
		KeyCurve:    "ecdsa",
		TokenSymbol: "BTC",
		RpcUsername: "bitcoin",
		RpcPassword: "bitcoin",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "3001",
		ChainType:   "sui",
		ChainUrl:    "https://fullnode.mainnet.sui.io:443",
		KeyCurve:    "eddsa",
		TokenSymbol: "SUI",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "3002",
		ChainType:   "sui",
		ChainUrl:    "https://fullnode.devnet.sui.io:443",
		KeyCurve:    "eddsa",
		TokenSymbol: "SUI",
		OpTimeout:   time.Second * 10,
	},
	{
		GenesisId:   "mainnet-v1.0",
		GenesisHash: "wGHE2Pwdvd7S12BL5FaOP20EGYesN73ktiC1qzkkit8=",
		ChainType:   "algorand",
		ChainUrl:    "https://mainnet-api.4160.nodely.dev",
		IndexerUrl:  "https://mainnet-idx.4160.nodely.dev",
		KeyCurve:    "ed25519",
		TokenSymbol: "ALGO",
		OpTimeout:   time.Second * 10,
	},
	{
		GenesisId:   "testnet-v1.0",
		GenesisHash: "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
		ChainType:   "algorand",
		ChainUrl:    "https://testnet-api.4160.nodely.dev",
		IndexerUrl:  "https://testnet-idx.4160.nodely.dev",
		KeyCurve:    "ed25519",
		TokenSymbol: "ALGO",
		OpTimeout:   time.Second * 10,
	},
	{
		GenesisId:   "betanet-v1.0",
		GenesisHash: "mFgazF+2uRS1tMiL9dsj01hJGySEmPN28B/TjjvpVW0=",
		ChainType:   "algorand",
		ChainUrl:    "https://betanet-api.4160.nodely.dev",
		IndexerUrl:  "https://betanet-idx.4160.nodely.dev", // port 443?
		KeyCurve:    "ed25519",
		TokenSymbol: "ALGO",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "testnet",
		ChainType:   "stellar",
		ChainUrl:    "https://horizon-testnet.stellar.org",
		KeyCurve:    "eddsa",
		TokenSymbol: "XLM",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "mainnet",
		ChainType:   "stellar",
		ChainUrl:    "https://horizon.stellar.org",
		KeyCurve:    "eddsa",
		TokenSymbol: "XLM",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "1003", // Made up id for ripple mainnet
		ChainType:   "ripple",
		ChainUrl:    "wss://s1.ripple.com:51233",
		KeyCurve:    "eddsa",
		TokenSymbol: "XRP",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "1004", // Made up id for ripple testnet
		ChainType:   "ripple",
		ChainUrl:    "wss://s.altnet.rippletest.net:51233",
		KeyCurve:    "eddsa",
		TokenSymbol: "XRP",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "1005", // Made up id for cardano mainnet
		ChainType:   "cardano",
		ChainUrl:    "https://cardano-mainnet.blockfrost.io/api/v0",
		KeyCurve:    "eddsa",
		TokenSymbol: "ADA",
		OpTimeout:   time.Second * 10,
	},
	{
		ChainId:     "1006", // Made up id for cardano testnet
		ChainType:   "cardano",
		ChainUrl:    "https://cardano-preprod.blockfrost.io/api/v0",
		KeyCurve:    "eddsa",
		TokenSymbol: "ADA",
		OpTimeout:   time.Second * 10,
	},
}

func GetChain(chainId string) (Chain, error) {
	for _, chain := range Chains {
		if chain.ChainId == chainId || chain.GenesisHash == chainId {
			return chain, nil
		}
	}
	return Chain{}, fmt.Errorf("chain not found")
}
