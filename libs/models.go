package libs

import (
	"time"

	"github.com/StripChain/strip-node/libs/blockchains"
	"github.com/google/uuid"
)

type Operation struct {
	ID               int64                    `json:"id"`
	SerializedTxn    string                   `json:"serializedTxn"`
	DataToSign       string                   `json:"dataToSign"`
	BlockchainID     blockchains.BlockchainID `json:"blockchainID"`
	NetworkType      blockchains.NetworkType  `json:"networkType"`
	GenesisHash      string                   `json:"genesisHash"`
	Status           OperationStatus          `json:"status"`
	Result           string                   `json:"result"`
	Type             OperationType            `json:"type"`
	Solver           string                   `json:"solver"`
	SolverMetadata   string                   `json:"solverMetadata"`
	SolverDataToSign string                   `json:"solverDataToSign"`
	SolverOutput     string                   `json:"solverOutput"`
	CreatedAt        time.Time                `json:"createdAt"`
}

type Intent struct {
	ID           uuid.UUID                `json:"id"`
	Operations   []Operation              `json:"operations"`
	Signature    string                   `json:"signature"`
	Identity     string                   `json:"identity"`
	BlockchainID blockchains.BlockchainID `json:"blockchainID"`
	NetworkType  blockchains.NetworkType  `json:"networkType"`
	Status       IntentStatus             `json:"status"`
	Expiry       time.Time                `json:"expiry"`
	CreatedAt    time.Time                `json:"createdAt"`
}

type OperationType string

const (
	OperationTypeTransaction   OperationType = "TRANSACTION"
	OperationTypeSolver        OperationType = "SOLVER"
	OperationTypeBridgeDeposit OperationType = "BRIDGE_DEPOSIT"
	OperationTypeSwap          OperationType = "SWAP"
	OperationTypeBurn          OperationType = "BURN"
	OperationTypeBurnSynthetic OperationType = "BURN_SYNTHETIC"
	OperationTypeWithdraw      OperationType = "WITHDRAW"
	OperationTypeSendToBridge  OperationType = "SEND_TO_BRIDGE"
)

type IntentStatus string

const (
	IntentStatusProcessing IntentStatus = "PROCESSING"
	IntentStatusCompleted  IntentStatus = "COMPLETED"
	IntentStatusFailed     IntentStatus = "FAILED"
	IntentStatusExpired    IntentStatus = "EXPIRED"
)

type OperationStatus string

const (
	OperationStatusPending   OperationStatus = "PENDING"
	OperationStatusWaiting   OperationStatus = "WAITING"
	OperationStatusCompleted OperationStatus = "COMPLETED"
	OperationStatusFailed    OperationStatus = "FAILED"
	OperationStatusExpired   OperationStatus = "EXPIRED"
)

// type Curve string

// const (
// 	CurveECDSA Curve = "ECDSA"
// 	CurveEDDSA Curve = "EDDSA"
// )

// func (c Curve) String() string {
// 	return string(c)
// }
