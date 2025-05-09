package uniswap_v3_solver

import (
	"encoding/json"
)

type Operation struct {
	SolverMetadata json.RawMessage `json:"solverMetadata"`
	Result         string          `json:"result"`
}

type Intent struct {
	Operations []Operation `json:"operations"`
	Identity   string      `json:"identity"`
}

type LPMetadata struct {
	Action    string `json:"action"`    // "mint" or "exit"
	Pool      string `json:"pool"`      // pool address
	TokenA    string `json:"tokenA"`    // first token address
	TokenB    string `json:"tokenB"`    // second token address
	AmountA   string `json:"amountA"`   // amount of first token
	AmountB   string `json:"amountB"`   // amount of second token
	Fee       uint32 `json:"fee"`       // pool fee tier (500 = 0.05%, 3000 = 0.3%, 10000 = 1%)
	TickLower int    `json:"tickLower"` // lower price bound
	TickUpper int    `json:"tickUpper"` // upper price bound
	TokenId   uint   `json:"tokenId"`   // NFT token ID
}

// LPOutput represents the result of a Uniswap V3 operation
type LPOutput struct {
	TxHash    string `json:"txHash"`
	TokenId   uint   `json:"tokenId,omitempty"`
	Liquidity string `json:"liquidity,omitempty"`
	AmountA   string `json:"amountA,omitempty"`
	AmountB   string `json:"amountB,omitempty"`
}
