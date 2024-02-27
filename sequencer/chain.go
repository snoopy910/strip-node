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
		ChainId:   "11155111",
		ChainType: "evm",
		ChainUrl:  "https://ethereum-sepolia-rpc.publicnode.com	",
		KeyCurve:  "ecdsa",
	},
	{
		ChainId:   "901",
		ChainType: "solana",
		ChainUrl:  "https://api.mainnet-beta.solana.com",
		KeyCurve:  "eddsa",
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
