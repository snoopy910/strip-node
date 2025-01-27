package common

import "fmt"

type Chain struct {
	ChainId     string
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
}

func GetChain(chainId string) (Chain, error) {
	for _, chain := range Chains {
		if chain.ChainId == chainId {
			return chain, nil
		}
	}
	return Chain{}, fmt.Errorf("chain not found")
}
