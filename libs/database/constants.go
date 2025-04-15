package database

const (
	INTENT_STATUS_PROCESSING = "processing"
	INTENT_STATUS_COMPLETED  = "completed"
	INTENT_STATUS_FAILED     = "failed"
	INTENT_STATUS_EXPIRED    = "expired"
)

const (
	OPERATION_STATUS_PENDING   = "pending"
	OPERATION_STATUS_WAITING   = "waiting"
	OPERATION_STATUS_COMPLETED = "completed"
	OPERATION_STATUS_FAILED    = "failed"
)

const (
	OPERATION_TYPE_TRANSACTION    = "transaction"
	OPERATION_TYPE_SOLVER         = "solver"
	OPERATION_TYPE_BRIDGE_DEPOSIT = "bridgeDeposit"
	OPERATION_TYPE_SWAP           = "swap"
	OPERATION_TYPE_BURN           = "burn"
	OPERATION_TYPE_BURN_SYNTHETIC = "burn_synthetic"
	OPERATION_TYPE_WITHDRAW       = "withdraw"
	OPERATION_TYPE_SEND_TO_BRIDGE = "sendToBridge"
)
