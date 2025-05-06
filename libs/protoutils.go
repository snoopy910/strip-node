package libs

import (
	"errors"
	"fmt"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/libs/blockchains"
	pb "github.com/StripChain/strip-node/libs/proto"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var curveToProto = map[common.Curve]pb.Curve{
	common.CurveEcdsa: pb.Curve_CURVE_ECDSA,
	common.CurveEddsa: pb.Curve_CURVE_EDDSA,
}

var protoCurveToCommon = map[pb.Curve]common.Curve{
	pb.Curve_CURVE_ECDSA: common.CurveEcdsa,
	pb.Curve_CURVE_EDDSA: common.CurveEddsa,
}

func ProtoToCommonCurve(c pb.Curve) (common.Curve, error) {
	if val, ok := protoCurveToCommon[c]; ok {
		return val, nil
	}
	return "", fmt.Errorf("unsupported or unspecified proto curve: %s", c.String())
}

func CommonCurveToProto(c common.Curve) (pb.Curve, error) {
	if val, ok := curveToProto[c]; ok {
		return val, nil
	}
	return pb.Curve_CURVE_UNSPECIFIED, fmt.Errorf("unsupported common curve: %s", c)
}

var blockchainsIDToProtoMap = map[blockchains.BlockchainID]pb.BlockchainID{
	blockchains.Bitcoin:    pb.BlockchainID_BITCOIN,
	blockchains.Ethereum:   pb.BlockchainID_ETHEREUM,
	blockchains.Solana:     pb.BlockchainID_SOLANA,
	blockchains.Dogecoin:   pb.BlockchainID_DOGECOIN,
	blockchains.Sui:        pb.BlockchainID_SUI,
	blockchains.Aptos:      pb.BlockchainID_APTOS,
	blockchains.Stellar:    pb.BlockchainID_STELLAR,
	blockchains.Algorand:   pb.BlockchainID_ALGORAND,
	blockchains.Ripple:     pb.BlockchainID_RIPPLE,
	blockchains.Cardano:    pb.BlockchainID_CARDANO,
	blockchains.StripChain: pb.BlockchainID_STRIPCHAIN,
	blockchains.Arbitrum:   pb.BlockchainID_ARBITRUM,
	blockchains.Sonic:      pb.BlockchainID_SONIC,
	blockchains.Berachain:  pb.BlockchainID_BERACHAIN,
}

var protoToBlockchainsIDMap = map[pb.BlockchainID]blockchains.BlockchainID{
	pb.BlockchainID_BITCOIN:    blockchains.Bitcoin,
	pb.BlockchainID_ETHEREUM:   blockchains.Ethereum,
	pb.BlockchainID_SOLANA:     blockchains.Solana,
	pb.BlockchainID_DOGECOIN:   blockchains.Dogecoin,
	pb.BlockchainID_SUI:        blockchains.Sui,
	pb.BlockchainID_APTOS:      blockchains.Aptos,
	pb.BlockchainID_STELLAR:    blockchains.Stellar,
	pb.BlockchainID_ALGORAND:   blockchains.Algorand,
	pb.BlockchainID_RIPPLE:     blockchains.Ripple,
	pb.BlockchainID_CARDANO:    blockchains.Cardano,
	pb.BlockchainID_STRIPCHAIN: blockchains.StripChain,
	pb.BlockchainID_ARBITRUM:   blockchains.Arbitrum,
	pb.BlockchainID_SONIC:      blockchains.Sonic,
	pb.BlockchainID_BERACHAIN:  blockchains.Berachain,
}

func ProtoToBlockchainsID(b pb.BlockchainID) (blockchains.BlockchainID, error) {
	if val, ok := protoToBlockchainsIDMap[b]; ok {
		return val, nil
	}
	return "", fmt.Errorf("unsupported or unspecified proto blockchain ID: %s", b.String())
}

func BlockchainsIDToProto(b blockchains.BlockchainID) (pb.BlockchainID, error) {
	if val, ok := blockchainsIDToProtoMap[b]; ok {
		return val, nil
	}
	return pb.BlockchainID_BLOCKCHAIN_ID_UNSPECIFIED, fmt.Errorf("unsupported blockchain ID: %s", b)
}

var networkTypeToProtoMap = map[blockchains.NetworkType]pb.NetworkType{
	blockchains.Mainnet: pb.NetworkType_MAINNET,
	blockchains.Testnet: pb.NetworkType_TESTNET,
	blockchains.Devnet:  pb.NetworkType_DEVNET,
	blockchains.Regnet:  pb.NetworkType_REGNET,
}

var protoToNetworkTypeMap = map[pb.NetworkType]blockchains.NetworkType{
	pb.NetworkType_MAINNET: blockchains.Mainnet,
	pb.NetworkType_TESTNET: blockchains.Testnet,
	pb.NetworkType_DEVNET:  blockchains.Devnet,
	pb.NetworkType_REGNET:  blockchains.Regnet,
}

func protoToNetworkType(n pb.NetworkType) (blockchains.NetworkType, error) {
	if val, ok := protoToNetworkTypeMap[n]; ok {
		return val, nil
	}
	return "", fmt.Errorf("unsupported or unspecified proto network type: %s", n.String())
}

func networkTypeToProto(n blockchains.NetworkType) (pb.NetworkType, error) {
	if val, ok := networkTypeToProtoMap[n]; ok {
		return val, nil
	}
	return pb.NetworkType_NETWORK_TYPE_UNSPECIFIED, fmt.Errorf("unsupported network type: %s", n)
}

var operationTypeToProtoMap = map[OperationType]pb.OperationType{
	OperationTypeTransaction:   pb.OperationType_TRANSACTION,
	OperationTypeSolver:        pb.OperationType_SOLVER,
	OperationTypeBridgeDeposit: pb.OperationType_BRIDGE_DEPOSIT,
	OperationTypeSwap:          pb.OperationType_SWAP,
	OperationTypeBurn:          pb.OperationType_BURN,
	OperationTypeBurnSynthetic: pb.OperationType_BURN_SYNTHETIC,
	OperationTypeWithdraw:      pb.OperationType_WITHDRAW,
	OperationTypeSendToBridge:  pb.OperationType_SEND_TO_BRIDGE,
}

var protoToOperationTypeMap = map[pb.OperationType]OperationType{
	pb.OperationType_TRANSACTION:    OperationTypeTransaction,
	pb.OperationType_SOLVER:         OperationTypeSolver,
	pb.OperationType_BRIDGE_DEPOSIT: OperationTypeBridgeDeposit,
	pb.OperationType_SWAP:           OperationTypeSwap,
	pb.OperationType_BURN:           OperationTypeBurn,
	pb.OperationType_BURN_SYNTHETIC: OperationTypeBurnSynthetic,
	pb.OperationType_WITHDRAW:       OperationTypeWithdraw,
	pb.OperationType_SEND_TO_BRIDGE: OperationTypeSendToBridge,
}

func protoToOperationType(o pb.OperationType) (OperationType, error) {
	if val, ok := protoToOperationTypeMap[o]; ok {
		return val, nil
	}
	return "", fmt.Errorf("unsupported or unspecified proto operation type: %s", o.String())
}

func operationTypeToProto(o OperationType) (pb.OperationType, error) {
	if val, ok := operationTypeToProtoMap[o]; ok {
		return val, nil
	}
	return pb.OperationType_OPERATION_TYPE_UNSPECIFIED, fmt.Errorf("unsupported operation type: %s", o)
}

var protoToIntentStatusMap = map[pb.IntentStatus]IntentStatus{
	pb.IntentStatus_INTENT_STATUS_PROCESSING: IntentStatusProcessing,
	pb.IntentStatus_INTENT_STATUS_COMPLETED:  IntentStatusCompleted,
	pb.IntentStatus_INTENT_STATUS_FAILED:     IntentStatusFailed,
	pb.IntentStatus_INTENT_STATUS_EXPIRED:    IntentStatusExpired,
}

var intentStatusToProtoMap = map[IntentStatus]pb.IntentStatus{
	IntentStatusProcessing: pb.IntentStatus_INTENT_STATUS_PROCESSING,
	IntentStatusCompleted:  pb.IntentStatus_INTENT_STATUS_COMPLETED,
	IntentStatusFailed:     pb.IntentStatus_INTENT_STATUS_FAILED,
	IntentStatusExpired:    pb.IntentStatus_INTENT_STATUS_EXPIRED,
}

func protoToIntentStatus(s pb.IntentStatus) (IntentStatus, error) {
	if val, ok := protoToIntentStatusMap[s]; ok {
		return val, nil
	}
	return "", fmt.Errorf("unsupported or unspecified proto intent status: %s", s.String())
}

func intentStatusToProto(s IntentStatus) (pb.IntentStatus, error) {
	if val, ok := intentStatusToProtoMap[s]; ok {
		return val, nil
	}
	return pb.IntentStatus_INTENT_STATUS_UNSPECIFIED, fmt.Errorf("unsupported intent status: %s", s)
}

var protoToOperationStatusMap = map[pb.OperationStatus]OperationStatus{
	pb.OperationStatus_OPERATION_STATUS_PENDING:   OperationStatusPending,
	pb.OperationStatus_OPERATION_STATUS_WAITING:   OperationStatusWaiting,
	pb.OperationStatus_OPERATION_STATUS_COMPLETED: OperationStatusCompleted,
	pb.OperationStatus_OPERATION_STATUS_FAILED:    OperationStatusFailed,
	pb.OperationStatus_OPERATION_STATUS_EXPIRED:   OperationStatusExpired,
}

var operationStatusToProtoMap = map[OperationStatus]pb.OperationStatus{
	OperationStatusPending:   pb.OperationStatus_OPERATION_STATUS_PENDING,
	OperationStatusWaiting:   pb.OperationStatus_OPERATION_STATUS_WAITING,
	OperationStatusCompleted: pb.OperationStatus_OPERATION_STATUS_COMPLETED,
	OperationStatusFailed:    pb.OperationStatus_OPERATION_STATUS_FAILED,
	OperationStatusExpired:   pb.OperationStatus_OPERATION_STATUS_EXPIRED,
}

func protoToOperationStatus(s pb.OperationStatus) (OperationStatus, error) {
	if val, ok := protoToOperationStatusMap[s]; ok {
		return val, nil
	}
	return "", fmt.Errorf("unsupported or unspecified proto operation status: %s", s.String())
}

func operationStatusToProto(s OperationStatus) (pb.OperationStatus, error) {
	if val, ok := operationStatusToProtoMap[s]; ok {
		return val, nil
	}
	return pb.OperationStatus_OPERATION_STATUS_UNSPECIFIED, fmt.Errorf("unsupported operation status: %s", s)
}

// ProtoToLibsIntent maps pb.Intent to Intent
func ProtoToLibsIntent(pi *pb.Intent) (*Intent, error) {
	if pi == nil {
		return nil, errors.New("nil proto intent provided")
	}

	intentID, err := uuid.Parse(pi.ID)
	if err != nil {
		// Allow empty ID string from proto? Or require valid UUID? Requiring for now.
		return nil, fmt.Errorf("invalid intent ID format (expected UUID string): %w", err)
	}

	blockchainID, err := ProtoToBlockchainsID(pi.BlockchainId)
	if err != nil {
		return nil, err
	}
	networkType, err := protoToNetworkType(pi.NetworkType)
	if err != nil {
		return nil, err
	}

	ops := make([]Operation, 0, len(pi.Operations)) // Use append for safety
	for i, pOp := range pi.Operations {
		libOp, err := ProtoToLibsOperation(pOp) // Needs implementation
		if err != nil {
			return nil, fmt.Errorf("failed converting operation %d: %w", i, err)
		}
		if libOp != nil { // Append only if mapping is successful
			ops = append(ops, *libOp)
		} else {
			return nil, fmt.Errorf("nil operation resulted from mapping operation %d", i)
		}
	}

	expiry := time.Time{}
	if pi.Expiry != nil {
		if err := pi.Expiry.CheckValid(); err != nil {
			return nil, fmt.Errorf("invalid expiry timestamp: %w", err)
		}
		expiry = pi.Expiry.AsTime()
	} else {
		// Default expiry if not set? Or error? Let's default to zero time.
		logger.Sugar().Warnw("Proto intent missing expiry timestamp", "intentID", pi.ID)
	}

	status, err := protoToIntentStatus(pi.Status)
	if err != nil {
		return nil, err
	}
	createdAt := time.Time{}
	if pi.CreatedAt != nil {
		createdAt = pi.CreatedAt.AsTime()
	} else {
		logger.Sugar().Warnw("Proto intent missing created at timestamp", "intentID", pi.ID)
	}

	return &Intent{
		ID:           intentID,
		Identity:     pi.Identity,
		BlockchainID: blockchainID,
		NetworkType:  networkType,
		Operations:   ops,
		Expiry:       expiry,
		Signature:    pi.Signature,
		Status:       status,
		CreatedAt:    createdAt,
	}, nil
}

// ProtoToLibsOperation maps pb.Operation to Operation
func ProtoToLibsOperation(po *pb.Operation) (*Operation, error) {
	if po == nil {
		return nil, errors.New("nil proto operation provided")
	}

	opType, err := protoToOperationType(po.Type)
	if err != nil {
		return nil, err
	}
	blockchainID, err := ProtoToBlockchainsID(po.BlockchainId)
	if err != nil {
		return nil, err
	}
	networkType, err := protoToNetworkType(po.NetworkType)
	if err != nil {
		return nil, err
	}

	var dataToSign *string
	if po.DataToSign != nil {
		ds := *po.DataToSign
		dataToSign = &ds
	}

	var serializedTxn *string
	if po.SerializedTxn != "" {
		st := po.SerializedTxn
		serializedTxn = &st
	}

	status, err := protoToOperationStatus(pb.OperationStatus(po.Status))
	if err != nil {
		return nil, err
	}

	return &Operation{
		ID:               int64(po.ID),
		Type:             opType,
		BlockchainID:     blockchainID,
		NetworkType:      networkType,
		Result:           po.Result,
		SerializedTxn:    serializedTxn,
		Solver:           po.Solver.Domain,
		SolverMetadata:   po.Solver.Metadata,
		SolverDataToSign: po.Solver.DataToSign,
		DataToSign:       dataToSign,
		GenesisHash:      po.GenesisHash,
		Status:           status,
		SolverOutput:     po.Solver.Output,
		CreatedAt:        po.CreatedAt.AsTime(),
	}, nil
}

func IntentToProto(intent *Intent) (*pb.Intent, error) {
	// Convert blockchain ID to proto enum
	blockchainID, err := BlockchainsIDToProto(intent.BlockchainID)
	if err != nil {
		blockchainID = pb.BlockchainID_BLOCKCHAIN_ID_UNSPECIFIED
	}

	// Convert network type to proto enum
	networkType, err := networkTypeToProto(intent.NetworkType)
	if err != nil {
		networkType = pb.NetworkType_NETWORK_TYPE_UNSPECIFIED
	}

	// Convert intent status to proto enum
	intentStatus, err := intentStatusToProto(intent.Status)
	if err != nil {
		intentStatus = pb.IntentStatus_INTENT_STATUS_UNSPECIFIED
	}

	// Convert operations to proto operations
	protoOperations := make([]*pb.Operation, 0, len(intent.Operations))
	for _, op := range intent.Operations {
		protoOp, err := operationToProto(&op)
		if err != nil {
			logger.Sugar().Errorw("failed to convert operation to proto", "error", err)
			return nil, fmt.Errorf("failed to convert operation to proto: %v", err)
		}
		protoOperations = append(protoOperations, protoOp)
	}

	return &pb.Intent{
		ID:           intent.ID.String(),
		Identity:     intent.Identity,
		BlockchainId: blockchainID,
		NetworkType:  networkType,
		Operations:   protoOperations,
		Expiry:       timestamppb.New(intent.Expiry),
		Signature:    intent.Signature,
		Status:       intentStatus,
		CreatedAt:    timestamppb.New(intent.CreatedAt),
	}, nil
}

func operationToProto(op *Operation) (*pb.Operation, error) {
	opType, err := operationTypeToProto(op.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to convert operation type to proto: %v", err)
	}

	opStatus, err := operationStatusToProto(op.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to convert operation status to proto: %v", err)
	}

	opBlockchainID, err := BlockchainsIDToProto(op.BlockchainID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert operation blockchain ID to proto: %v", err)
	}

	opNetworkType, err := networkTypeToProto(op.NetworkType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert operation network type to proto: %v", err)
	}

	protoOp := &pb.Operation{
		ID:           int32(op.ID),
		Type:         opType,
		BlockchainId: opBlockchainID,
		NetworkType:  opNetworkType,
		Result:       op.Result,
		Status:       opStatus,
		GenesisHash:  op.GenesisHash,
		DataToSign:   op.DataToSign,
		CreatedAt:    timestamppb.New(op.CreatedAt),
		Solver: &pb.Solver{
			Domain:     op.Solver,
			Metadata:   op.SolverMetadata,
			DataToSign: op.SolverDataToSign,
			Output:     op.SolverOutput,
		},
	}

	if op.SerializedTxn != nil {
		protoOp.SerializedTxn = *op.SerializedTxn
	}

	return protoOp, nil
}
