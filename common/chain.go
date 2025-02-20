package common

import "fmt"

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
}

var Chains = []Chain{
	{
		ChainId:     "1337",
		ChainType:   "evm",
		ChainUrl:    "http://ganache:8545",
		KeyCurve:    "ecdsa",
		TokenSymbol: "ETH",
	},
	{
		ChainId:     "11155111",
		ChainType:   "evm",
		ChainUrl:    "https://ethereum-sepolia-rpc.publicnode.com",
		KeyCurve:    "ecdsa",
		TokenSymbol: "ETH",
	},
	{
		ChainId:     "901",
		ChainType:   "solana",
		ChainUrl:    "https://api.devnet.solana.com",
		KeyCurve:    "eddsa",
		TokenSymbol: "SOL",
	},
	{
		ChainId:     "1",
		ChainType:   "evm",
		ChainUrl:    "https://ethereum-rpc.publicnode.com",
		KeyCurve:    "ecdsa",
		TokenSymbol: "ETH",
	},
	{
		ChainId:     "900",
		ChainType:   "solana",
		ChainUrl:    "https://api.solana.com",
		KeyCurve:    "eddsa",
		TokenSymbol: "SOL",
	},
	{
		ChainId:     "137",
		ChainType:   "evm",
		ChainUrl:    "https://polygon-pokt.nodies.app",
		KeyCurve:    "ecdsa",
		TokenSymbol: "MATIC",
	},
	{
		ChainId:     "11",
		ChainType:   "aptos",
		ChainUrl:    "https://fullnode.mainnet.aptoslabs.com",
		KeyCurve:    "eddsa",
		TokenSymbol: "APT",
	},
	{
		ChainId:     "167",
		ChainType:   "aptos",
		ChainUrl:    "https://fullnode.devnet.aptoslabs.com",
		KeyCurve:    "eddsa",
		TokenSymbol: "APT",
	},
	{
		ChainId:     "1000",
		ChainType:   "bitcoin",
		ChainUrl:    "https://api.blockcypher.com/v1/btc/main",
		KeyCurve:    "ecdsa",
		TokenSymbol: "BTC",
	},
	{
		ChainId:     "1001",
		ChainType:   "bitcoin",
		ChainUrl:    "https://api.blockcypher.com/v1/btc/test3",
		KeyCurve:    "ecdsa",
		TokenSymbol: "BTC",
	},
	{
		ChainId:     "1002",
		ChainType:   "bitcoin",
		ChainUrl:    "http://127.0.0.1:18443",
		KeyCurve:    "ecdsa",
		TokenSymbol: "BTC",
	},
	{
		GenesisId:   "mainnet-v1.0",
		GenesisHash: "wGHE2Pwdvd7S12BL5FaOP20EGYesN73ktiC1qzkkit8=",
		ChainType:   "algorand",
		ChainUrl:    "https://mainnet-api.4160.nodely.dev",
		IndexerUrl:  "https://mainnet-idx.4160.nodely.dev",
		KeyCurve:    "ed25519",
		TokenSymbol: "ALGO",
	},
	{
		GenesisId:   "testnet-v1.0",
		GenesisHash: "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
		ChainType:   "algorand",
		ChainUrl:    "https://testnet-api.4160.nodely.dev",
		IndexerUrl:  "https://testnet-idx.4160.nodely.dev",
		KeyCurve:    "ed25519",
		TokenSymbol: "ALGO",
	},
	{
		GenesisId:   "betanet-v1.0",
		GenesisHash: "mFgazF+2uRS1tMiL9dsj01hJGySEmPN28B/TjjvpVW0=",
		ChainType:   "algorand",
		ChainUrl:    "https://betanet-api.4160.nodely.dev",
		IndexerUrl:  "https://betanet-idx.4160.nodely.dev", // port 443?
		KeyCurve:    "ed25519",
		TokenSymbol: "ALGO",
	},
	{
		ChainId:     "testnet",
		ChainType:   "stellar",
		ChainUrl:    "https://horizon-testnet.stellar.org",
		KeyCurve:    "eddsa",
		TokenSymbol: "XLM",
	},
	{
		ChainId:     "mainnet",
		ChainType:   "stellar",
		ChainUrl:    "https://horizon.stellar.org",
		KeyCurve:    "eddsa",
		TokenSymbol: "XLM",
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
