package lending_solver

import "encoding/json"

// uint256 equivalent in Go
type uint256 struct {
	Int string `json:"int"`
}

// Intent represents a user's intention to perform lending operations
type Intent struct {
	Operations []Operation `json:"operations"`
	Identity   string      `json:"identity"` // caller's address
}

// Operation represents a single lending operation
type Operation struct {
	ID             int64           `json:"id"`
	Type           string          `json:"type"`
	ChainId        string          `json:"chainId"`
	SolverMetadata json.RawMessage `json:"solverMetadata"`
	Result         string          `json:"result"`
}

// LendingMetadata represents the metadata for lending operations
type LendingMetadata struct {
	Action       string  `json:"action"`       // supply, borrow, repay, withdraw
	Token        string  `json:"token"`        // token address for supply/withdraw
	Amount       uint256 `json:"amount"`       // amount to supply/borrow/repay/withdraw
	IsCollateral bool    `json:"isCollateral"` // whether token should be used as collateral
}

// LendingOutput represents the result of a lending operation
type LendingOutput struct {
	TxHash        string  `json:"txHash"`
	Token         string  `json:"token,omitempty"`
	Amount        uint256 `json:"amount,omitempty"`
	HealthFactor  uint256 `json:"healthFactor,omitempty"`
	Rate          string  `json:"rate,omitempty"`
	RemainingDebt string  `json:"remainingDebt,omitempty"`
}
