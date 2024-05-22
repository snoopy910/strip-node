package sequencer

import "fmt"

type Chain struct {
	ChainId   string
	ChainType string
	ChainUrl  string
	KeyCurve  string
}

var Chains = []Chain{
	{
		ChainId:   "1337",
		ChainType: "evm",
		ChainUrl:  "http://ganache:8545",
		KeyCurve:  "ecdsa",
	},
	{
		ChainId:   "11155111",
		ChainType: "evm",
		ChainUrl:  "https://ethereum-sepolia-rpc.publicnode.com",
		KeyCurve:  "ecdsa",
	},
	{
		ChainId:   "901",
		ChainType: "solana",
		ChainUrl:  "https://api.devnet.solana.com",
		KeyCurve:  "eddsa",
	},
	{
		ChainId:   "1",
		ChainType: "evm",
		ChainUrl:  "https://ethereum-rpc.publicnode.com",
		KeyCurve:  "ecdsa",
	},
	{
		ChainId:   "900",
		ChainType: "solana",
		ChainUrl:  "https://api.solana.com",
		KeyCurve:  "eddsa",
	},
	{
		ChainId:   "137",
		ChainType: "evm",
		ChainUrl:  "https://polygon-pokt.nodies.app",
		KeyCurve:  "ecdsa",
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
