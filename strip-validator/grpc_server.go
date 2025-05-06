package main

import (
	"context"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"time"

	"github.com/StripChain/strip-node/ERC20"
	"github.com/StripChain/strip-node/bitcoin"
	"github.com/StripChain/strip-node/bridge"
	"github.com/StripChain/strip-node/dogecoin"
	identityVerification "github.com/StripChain/strip-node/identity"
	db "github.com/StripChain/strip-node/libs/database"
	pb "github.com/StripChain/strip-node/libs/proto"
	"github.com/StripChain/strip-node/ripple"
	"github.com/StripChain/strip-node/solver"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/stellar/go/strkey"
	"golang.org/x/crypto/blake2b"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/identity"
	"github.com/StripChain/strip-node/libs"
	"github.com/StripChain/strip-node/libs/blockchains"
	"github.com/StripChain/strip-node/util/logger"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	ecdsaKeygen "github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/mr-tron/base58"
)

type validatorServer struct {
	pb.UnimplementedValidatorServiceServer
	host host.Host
}

// NewValidatorServer creates a new server instance
func NewValidatorServer(host host.Host) *validatorServer {
	return &validatorServer{
		host: host,
	}
}

func (s *validatorServer) Keygen(ctx context.Context, req *pb.KeygenRequest) (*pb.KeygenResponse, error) {
	logger.Sugar().Infow("Received gRPC Keygen request",
		"identity", req.Identity,
		"identityCurve", req.IdentityCurve,
		"signersCount", len(req.Signers))

	identityCurve, err := libs.ProtoToCommonCurve(req.IdentityCurve)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid identity curve: %v", err)
	}

	curvesToGenerate := []common.Curve{common.CurveEcdsa, common.CurveEddsa}
	for _, curve := range curvesToGenerate {

		// TODO: Make this asynchronous
		key := req.Identity + "_" + string(identityCurve) + "_" + string(curve)

		opChan := make(chan string)
		errorChan := make(chan error, 1)

		// FIXME: This is not thread safe and we shouldn't use a global map
		keygenGeneratedChan[key] = opChan

		go func() {
			defer func() {
				// Ensure channel is deleted even on panic or completion
				// delete(keygenGeneratedChan, key)
				if r := recover(); r != nil {
					logger.Sugar().Errorw("Panic in gRPC keygen goroutine", "key", key, "error", fmt.Sprintf("%v", r))
					errorChan <- status.Errorf(codes.Internal, "internal server error: panic in keygen operation")
				}
			}()
			logger.Sugar().Infow("Calling generateKeygenMessage for gRPC request", "key", key)
			generateKeygenMessage(req.Identity, identityCurve, curve, req.Signers)
		}()

		logger.Sugar().Infow("Waiting for keygen operation to complete", "key", key)
		select {
		case <-opChan:
			logger.Sugar().Infow("gRPC Keygen operation completed successfully", "curve", curve)
		case err := <-errorChan:
			logger.Sugar().Errorw("gRPC Keygen operation failed", "key", key, "error", err)
			// Propagate gRPC status error if available
			if s, ok := status.FromError(err); ok {
				return nil, s.Err()
			}
			return nil, status.Errorf(codes.Internal, "keygen failed: %+v", err)
		case <-ctx.Done():
			logger.Sugar().Warnw("gRPC Keygen request cancelled by client", "key", key)
			return nil, status.Error(codes.Canceled, "client cancelled request")
		case <-time.After(5 * time.Minute):
			logger.Sugar().Errorw("gRPC Keygen operation timed out", "key", key)
			return nil, status.Error(codes.DeadlineExceeded, "keygen operation timed out")
		}
	}
	return &pb.KeygenResponse{Message: "Keygen operation completed successfully"}, nil
}

func (s *validatorServer) GetAddresses(ctx context.Context, req *pb.GetAddressesRequest) (*pb.GetAddressesResponse, error) {
	logger.Sugar().Infow("GetAddresses requested", "identity", req.Identity, "identityCurve", req.IdentityCurve)

	identityCurveEnum, err := libs.ProtoToCommonCurve(req.IdentityCurve)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid identity curve: %v", err)
	}

	var rawKeyEddsa *eddsaKeygen.LocalPartySaveData
	var rawKeyEcdsa *ecdsaKeygen.LocalPartySaveData

	response := &pb.GetAddressesResponse{
		Addresses: make(map[int32]*pb.BlockchainAddressMap),
	}

	keyCurvesToCheck := []common.Curve{common.CurveEcdsa, common.CurveEddsa}
	foundAnyKey := false

	for _, keyCurve := range keyCurvesToCheck {
		keyShare, err := GetKeyShare(req.Identity, identityCurveEnum, keyCurve)
		if err != nil {
			logger.Sugar().Errorw("Error getting key share from storage", "identity", req.Identity, "curve", keyCurve, "error", err)
			return nil, status.Errorf(codes.Internal, "error retrieving key share for curve %s: %v", keyCurve, err)
		}

		if keyShare == "" {
			logger.Sugar().Warnw("Key share not found in storage", "identity", req.Identity, "curve", keyCurve)
			continue
		}

		foundAnyKey = true

		switch keyCurve {
		case common.CurveEddsa:
			err = json.Unmarshal([]byte(keyShare), &rawKeyEddsa)
			if err != nil {
				logger.Sugar().Errorw("Failed to unmarshal EDDSA key share", "identity", req.Identity, "error", err)
				return nil, status.Errorf(codes.Internal, "failed to process EDDSA key share")
			}
		case common.CurveEcdsa:
			err = json.Unmarshal([]byte(keyShare), &rawKeyEcdsa)
			if err != nil {
				logger.Sugar().Errorw("Failed to unmarshal ECDSA key share", "identity", req.Identity, "error", err)
				return nil, status.Errorf(codes.Internal, "failed to process ECDSA key share")
			}
		default:
			logger.Sugar().Warnw("Unsupported key curve", "identity", req.Identity, "curve", keyCurve)
			return nil, status.Errorf(codes.Internal, "unsupported key curve: %s", keyCurve)
		}
	}

	if !foundAnyKey {
		return nil, status.Errorf(codes.NotFound, "no key shares found for identity %s with identitycurve %s", req.Identity, identityCurveEnum)
	}

	registeredChains := blockchains.GetRegisteredBlockchains()

	for _, blockchainIDGo := range registeredChains {
		opBlockchain, err := blockchains.GetBlockchain(blockchainIDGo, blockchains.Mainnet)
		if err != nil {
			opBlockchain, err = blockchains.GetBlockchain(blockchainIDGo, blockchains.Testnet)
			if err != nil {
				logger.Sugar().Errorw("Cannot get blockchain instance", "blockchainID", blockchainIDGo, "error", err)
				return nil, status.Errorf(codes.Internal, "cannot get blockchain instance: %v", err)
			}
		}

		requiredKeyCurve := opBlockchain.KeyCurve()
		var currentRawKey interface{}

		if requiredKeyCurve == common.CurveEddsa && rawKeyEddsa != nil {
			currentRawKey = rawKeyEddsa
		} else if requiredKeyCurve == common.CurveEcdsa && rawKeyEcdsa != nil {
			currentRawKey = rawKeyEcdsa
		} else {
			logger.Sugar().Errorw("Skipping blockchain, required key curve not loaded/found", "id", blockchainIDGo, "requiredCurve", requiredKeyCurve)
			return nil, status.Errorf(codes.Internal, "required key curve not loaded/found for blockchain %s", blockchainIDGo)
		}

		blockchainIDProto, err := libs.BlockchainsIDToProto(blockchainIDGo)
		if err != nil {
			logger.Sugar().Warnw("Cannot map Go BlockchainID to proto", "blockchainID", blockchainIDGo, "error", err)
			return nil, status.Errorf(codes.Internal, "cannot map Go BlockchainID to proto: %v", err)
		}
		blockchainIDProtoInt := int32(blockchainIDProto)

		if _, ok := response.Addresses[blockchainIDProtoInt]; !ok {
			response.Addresses[blockchainIDProtoInt] = &pb.BlockchainAddressMap{
				NetworkAddresses: make(map[int32]*pb.AddressDetail),
			}
		}

		switch blockchainIDGo {
		case blockchains.Solana:
			key := currentRawKey.(*eddsaKeygen.LocalPartySaveData)
			pk := edwards.PublicKey{Curve: key.EDDSAPub.Curve(), X: key.EDDSAPub.X(), Y: key.EDDSAPub.Y()}
			address := base58.Encode(pk.Serialize())
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_MAINNET, address)
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_TESTNET, address)
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_DEVNET, address)

		case blockchains.Bitcoin:
			key := currentRawKey.(*ecdsaKeygen.LocalPartySaveData)
			publicKeyBytes, err := getCompressedPublicKeyBytes(key)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "bitcoin pubkey generation error: %v", err)
			}
			mainnetAddress, testnetAddress, regtestAddress := bitcoin.PublicKeyToBitcoinAddresses(publicKeyBytes)
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_MAINNET, mainnetAddress)
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_TESTNET, testnetAddress)
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_REGNET, regtestAddress)

		case blockchains.Dogecoin:
			key := currentRawKey.(*ecdsaKeygen.LocalPartySaveData)
			publicKeyHex, err := getUncompressedPublicKeyHex(key)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "dogecoin pubkey generation error: %v", err)
			}

			mainnetAddress, err1 := dogecoin.PublicKeyToAddress(publicKeyHex)
			testnetAddress, err2 := dogecoin.PublicKeyToTestnetAddress(publicKeyHex)
			if err1 != nil || err2 != nil {
				return nil, status.Errorf(codes.Internal, "doge addr gen error: main=%v test=%v", err1, err2)
			}
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_MAINNET, mainnetAddress)
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_TESTNET, testnetAddress)
		case blockchains.Sui:
			key := currentRawKey.(*eddsaKeygen.LocalPartySaveData)
			pk := edwards.PublicKey{Curve: key.EDDSAPub.Curve(), X: key.EDDSAPub.X(), Y: key.EDDSAPub.Y()}
			pkBytes := pk.Serialize()
			flag := byte(0x00) // Ed25519 flag
			hasher, _ := blake2b.New256(nil)
			hasher.Write([]byte{flag})
			hasher.Write(pkBytes)
			address := "0x" + hex.EncodeToString(hasher.Sum(nil))
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_MAINNET, address)
		case blockchains.Aptos:
			key := currentRawKey.(*eddsaKeygen.LocalPartySaveData)
			pk := edwards.PublicKey{Curve: key.EDDSAPub.Curve(), X: key.EDDSAPub.X(), Y: key.EDDSAPub.Y()}
			address := "0x" + hex.EncodeToString(pk.Serialize())
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_MAINNET, address)
		case blockchains.Stellar:
			key := currentRawKey.(*eddsaKeygen.LocalPartySaveData)
			pk := edwards.PublicKey{Curve: key.EDDSAPub.Curve(), X: key.EDDSAPub.X(), Y: key.EDDSAPub.Y()}
			pkBytes := pk.Serialize()
			if len(pkBytes) != 32 {
				return nil, status.Error(codes.Internal, "invalid stellar pubkey length")
			}
			versionByte := strkey.VersionByteAccountID
			address, err := strkey.Encode(versionByte, pkBytes)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "stellar address encoding error: %v", err)
			}
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_MAINNET, address)
		case blockchains.Algorand:
			key := currentRawKey.(*eddsaKeygen.LocalPartySaveData)
			pk := edwards.PublicKey{Curve: key.EDDSAPub.Curve(), X: key.EDDSAPub.X(), Y: key.EDDSAPub.Y()}
			pkBytes := pk.Serialize()
			// Algorand address = base32(pkBytes + checksum(pkBytes))
			hasher := sha512.New512_256()
			hasher.Write(pkBytes)
			checksum := hasher.Sum(nil)[28:] // Last 4 bytes of sha512_256
			addressBytes := append(pkBytes, checksum...)
			address := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(addressBytes)
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_MAINNET, address)
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_TESTNET, address)

		case blockchains.Ripple:
			key := currentRawKey.(*eddsaKeygen.LocalPartySaveData)
			// Assuming ripple pkg exists with PublicKeyToAddress function
			address := ripple.PublicKeyToAddress(key)
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_MAINNET, address)

		case blockchains.Cardano:
			key := currentRawKey.(*eddsaKeygen.LocalPartySaveData)
			// Cardano address generation is complex. api.go just returned the hex pubkey.
			pk := edwards.PublicKey{Curve: key.EDDSAPub.Curve(), X: key.EDDSAPub.X(), Y: key.EDDSAPub.Y()}
			address := hex.EncodeToString(pk.Serialize()) // Returning hex pubkey as per api.go
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_MAINNET, address)
			addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_TESTNET, address)

		default:
			if blockchains.IsEVMBlockchain(blockchainIDGo) {
				if requiredKeyCurve == common.CurveEcdsa && rawKeyEcdsa != nil {
					key := rawKeyEcdsa
					publicKeyBytes, err := getUncompressedPublicKeyBytes(key)
					if err != nil {
						return nil, status.Errorf(codes.Internal, "evm pubkey generation error: %v", err)
					}
					address := publicKeyToAddress(publicKeyBytes)
					addAddressDetail(response, blockchainIDProtoInt, pb.NetworkType_MAINNET, address)
				} else {
					logger.Sugar().Errorw("Unsupported blockchain ID", "blockchainID", blockchainIDGo)
					return nil, status.Errorf(codes.Internal, "unsupported blockchain ID: %s", blockchainIDGo)
				}
			} else {
				logger.Sugar().Errorw("Unsupported blockchain ID", "blockchainID", blockchainIDGo)
				return nil, status.Errorf(codes.Internal, "unsupported blockchain ID: %s", blockchainIDGo)
			}
		}
	}

	pk := edwards.PublicKey{Curve: rawKeyEddsa.EDDSAPub.Curve(), X: rawKeyEddsa.EDDSAPub.X(), Y: rawKeyEddsa.EDDSAPub.Y()}
	response.EddsaAddress = hex.EncodeToString(pk.Serialize())

	publicKeyBytes, err := getUncompressedPublicKeyBytes(rawKeyEcdsa)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed generating ECDSA address: %v", err)
	}
	response.EcdsaAddress = publicKeyToAddress(publicKeyBytes)
	return response, nil
}

func (s *validatorServer) SignIntentOperation(ctx context.Context, req *pb.SignIntentOperationRequest) (*pb.SignIntentOperationResponse, error) {
	if req.Intent == nil {
		return nil, status.Error(codes.InvalidArgument, "missing intent")
	}
	if int(req.OperationIndex) >= len(req.Intent.Operations) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid operation index %d (intent has %d operations)", req.OperationIndex, len(req.Intent.Operations))
	}

	logger.Sugar().Infow("gRPC SignIntentOperation requested",
		"intentID", req.Intent.ID,
		"opIndex", req.OperationIndex)

	intent, err := libs.ProtoToLibsIntent(req.Intent)
	if err != nil {
		logger.Sugar().Errorw("Failed to map proto Intent to libs.Intent", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to process intent structure: %v", err)
	}
	operationIndex := int(req.OperationIndex)

	if operationIndex >= len(intent.Operations) {
		logger.Sugar().Errorw("Mismatch between proto operations and mapped Go operations", "protoCount", len(req.Intent.Operations), "goCount", len(intent.Operations))
		return nil, status.Errorf(codes.Internal, "failed to map operations from intent")
	}
	operation := intent.Operations[operationIndex]

	if intent.Expiry.Before(time.Now().UTC()) {
		logger.Sugar().Errorw("Intent has expired", "intentID", intent.ID, "expiry", intent.Expiry)
		return nil, status.Errorf(codes.FailedPrecondition, "intent has expired")
	}

	intentBlockchain, err := blockchains.GetBlockchain(intent.BlockchainID, intent.NetworkType)
	if err != nil {
		logger.Sugar().Errorw("Error getting intent blockchain", "id", intent.BlockchainID, "net", intent.NetworkType, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get intent blockchain info: %v", err)
	}
	opBlockchain, err := blockchains.GetBlockchain(operation.BlockchainID, operation.NetworkType)
	if err != nil {
		logger.Sugar().Errorw("Error getting operation blockchain", "id", operation.BlockchainID, "net", operation.NetworkType, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get operation blockchain info: %v", err)
	}

	var msg string
	switch operation.Type {
	case libs.OperationTypeTransaction:
		if operation.DataToSign != nil {
			msg = *operation.DataToSign
		}
	case libs.OperationTypeSendToBridge:
		// Verify only operation for bridging
		// Get bridgewallet by calling /getBridgeAddress from sequencer api
		req, err := http.NewRequest("GET", SequencerHost+"/getBridgeAddress", nil)
		if err != nil {
			logger.Sugar().Errorw("error creating request", "error", err)
			return nil, status.Errorf(codes.Internal, "error creating request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Sugar().Errorw("error sending request", "error", err)
			return nil, status.Errorf(codes.Internal, "error sending request: %v", err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Sugar().Errorw("error reading response body", "error", err)
			return nil, status.Errorf(codes.Internal, "error reading response body: %v", err)
		}

		var bridgeWallet db.WalletSchema
		err = json.Unmarshal(body, &bridgeWallet)
		if err != nil {
			logger.Sugar().Errorw("error unmarshalling response body", "error", err)
			return nil, status.Errorf(codes.Internal, "error unmarshalling response body: %v", err)
		}

		// if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "bitcoin_ecdsa" || operation.KeyCurve == "dogecoin_ecdsa" {
		// 	chain, err := common.GetChain(operation.ChainId)
		// 	if err != nil {
		// 		logger.Sugar().Errorw("error getting chain", "error", err)
		// 		return
		// 	}

		// 	// Extract destination address from serialized transaction
		// 	var destAddress string
		// 	if chain.ChainType == "bitcoin" || chain.ChainType == "dogecoin" {
		// 		// For Bitcoin, decode the serialized transaction to get output address
		// 		var tx wire.MsgTx
		// 		txBytes, err := hex.DecodeString(operation.SerializedTxn)
		// 		if err != nil {
		// 			logger.Sugar().Errorw("error decoding bitcoin&dogecoin transaction", "error", err)
		// 			return
		// 		}
		// 		if err := tx.Deserialize(bytes.NewReader(txBytes)); err != nil {
		// 			logger.Sugar().Errorw("error deserializing bitcoin&dogecoin transaction", "error", err)
		// 			return
		// 		}
		// 		// Get the first output's address (assuming it's the bridge address)
		// 		if len(tx.TxOut) > 0 {
		// 			_, addrs, _, err := txscript.ExtractPkScriptAddrs(tx.TxOut[0].PkScript, nil)
		// 			if err != nil || len(addrs) == 0 {
		// 				logger.Sugar().Errorw("error extracting bitcoin&dogecoin address", "error", err)
		// 				return
		// 			}
		// 			destAddress = addrs[0].String()
		// 		}
		// 	} else {
		// 		// For EVM chains, decode the transaction to get the 'to' address
		// 		txBytes, err := hex.DecodeString(operation.SerializedTxn)
		// 		if err != nil {
		// 			logger.Sugar().Errorw("error decoding EVM transaction", "error", err)
		// 			return
		// 		}
		// 		tx := new(types.Transaction)
		// 		if err := rlp.DecodeBytes(txBytes, tx); err != nil {
		// 			logger.Sugar().Errorw("error deserializing EVM transaction", "error", err)
		// 			return
		// 		}
		// 		// Check if this is an ERC20 transfer (function signature: transfer(address,uint256))
		// 		// ERC20 transfer function signature is: 0xa9059cbb
		// 		isERC20Transfer := false
		// 		contractAddress := tx.To().Hex()
		// 		if tx.Data() != nil && len(tx.Data()) >= 4 && hex.EncodeToString(tx.Data()[:4]) == "a9059cbb" {
		// 			// This is an ERC20 transfer - extract the recipient address from the data
		// 			// The recipient address is the first parameter of the transfer function (32-byte padded)
		// 			isERC20Transfer = true
		// 			if len(tx.Data()) >= 36 {
		// 				// Extract the recipient address from position 4:36 (32 bytes)
		// 				recipientBytes := tx.Data()[4:36]
		// 				// Convert to address (take last 20 bytes for proper Ethereum address length)
		// 				destAddress = "0x" + hex.EncodeToString(recipientBytes[12:])
		// 				logger.Sugar().Infow("detected ERC20 transfer", "contract", contractAddress, "recipient", destAddress)
		// 			} else {
		// 				logger.Sugar().Errorw("invalid ERC20 transfer data length", "data", hex.EncodeToString(tx.Data()))
		// 				return
		// 			}
		// 		} else {
		// 			// For regular ETH transfers, use the 'to' address directly
		// 			destAddress = contractAddress
		// 		}

		// 		// Verify destination address matches bridge wallet
		// 		var expectedAddress string
		// 		if chain.ChainType == "bitcoin" {
		// 			expectedAddress = bridgeWallet.BitcoinMainnetPublicKey
		// 		} else if chain.ChainType == "dogecoin" {
		// 			expectedAddress = bridgeWallet.DogecoinMainnetPublicKey
		// 		} else {
		// 			// For Ethereum chains, determine if we need to use the bridge contract address
		// 			if chain.ChainType == "evm" && !isERC20Transfer {
		// 				// For native ETH transfers, get the bridge address from the dedicated endpoint
		// 				bridgeReq, err := http.NewRequest("GET", SequencerHost+"/getBridgeAddress", nil)
		// 				if err != nil {
		// 					logger.Sugar().Errorw("error creating bridge address request", "error", err)
		// 					return
		// 				}

		// 				bridgeReq.Header.Set("Content-Type", "application/json")
		// 				bridgeClient := &http.Client{}
		// 				bridgeResp, err := bridgeClient.Do(bridgeReq)
		// 				if err != nil {
		// 					logger.Sugar().Errorw("error fetching bridge address", "error", err)
		// 					return
		// 				}

		// 				defer bridgeResp.Body.Close()

		// 				bridgeBody, err := io.ReadAll(bridgeResp.Body)
		// 				if err != nil {
		// 					logger.Sugar().Errorw("error reading bridge address response", "error", err)
		// 					return
		// 				}

		// 				var bridgeAddressWallet db.WalletSchema
		// 				err = json.Unmarshal(bridgeBody, &bridgeAddressWallet)
		// 				if err != nil {
		// 					logger.Sugar().Errorw("error unmarshalling bridge address response", "error", err)
		// 					return
		// 				}

		// 				expectedAddress = bridgeAddressWallet.ECDSAPublicKey
		// 				logger.Sugar().Infow("using bridge contract address for native ETH transfer", "address", expectedAddress)
		// 			} else {
		// 				expectedAddress = bridgeWallet.ECDSAPublicKey
		// 			}
		// 		}

		// 		// Verify the extracted destination matches the bridge wallet
		// 		if !strings.EqualFold(destAddress, expectedAddress) {
		// 			logger.Sugar().Errorw("Invalid bridge destination address", "expected", expectedAddress, "got", destAddress)
		// 			return
		// 		}
		// 	}
		// } else if operation.KeyCurve == "eddsa" || operation.KeyCurve == "aptos_eddsa" || operation.KeyCurve == "stellar_eddsa" || operation.KeyCurve == "algorand_eddsa" || operation.KeyCurve == "ripple_eddsa" || operation.KeyCurve == "cardano_eddsa" || operation.KeyCurve == "sui_eddsa" {
		// 	chain, err := common.GetChain(operation.ChainId)
		// 	if err != nil {
		// 		logger.Sugar().Errorw("error getting chain", "error", err)
		// 		return
		// 	}

		// 	// Verify destination address matches bridge wallet based on chain type
		// 	var validDestination bool
		// 	var destAddress string

		// 	// Extract destination address from serialized transaction based on chain type
		// 	switch chain.ChainType {
		// 	case "solana":
		// 		// Decode base58 transaction and extract destination
		// 		decodedTxn, err := base58.Decode(operation.SerializedTxn)
		// 		if err != nil {
		// 			logger.Sugar().Errorw("error decoding Solana transaction", "error", err)
		// 			return
		// 		}
		// 		tx, err := solanasdk.TransactionFromDecoder(bin.NewBinDecoder(decodedTxn))
		// 		if err != nil || len(tx.Message.Instructions) == 0 {
		// 			logger.Sugar().Errorw("error deserializing Solana transaction", "error", err)
		// 			return
		// 		}
		// 		// Get the first instruction's destination account index
		// 		destAccountIndex := tx.Message.Instructions[0].Accounts[1]
		// 		// Get the actual account address from the message accounts
		// 		destAddress = tx.Message.AccountKeys[destAccountIndex].String()
		// 	case "aptos":
		// 		// For Aptos, the destination is in the transaction payload
		// 		var aptosPayload struct {
		// 			Function string   `json:"function"`
		// 			Args     []string `json:"arguments"`
		// 		}
		// 		if err := json.Unmarshal([]byte(operation.SerializedTxn), &aptosPayload); err != nil {
		// 			logger.Sugar().Errorw("error parsing Aptos transaction", "error", err)
		// 			return
		// 		}
		// 		if len(aptosPayload.Args) > 0 {
		// 			destAddress = aptosPayload.Args[0] // First arg is typically the recipient
		// 		}
		// 	case "stellar":
		// 		// For Stellar, parse the XDR transaction envelope
		// 		var txEnv xdr.TransactionEnvelope
		// 		err := xdr.SafeUnmarshalBase64(operation.SerializedTxn, &txEnv)
		// 		if err != nil {
		// 			logger.Sugar().Errorw("error parsing Stellar transaction", "error", err)
		// 			return
		// 		}

		// 		// Get the first operation's destination
		// 		if len(txEnv.Operations()) > 0 {
		// 			if paymentOp, ok := txEnv.Operations()[0].Body.GetPaymentOp(); ok {
		// 				destAddress = paymentOp.Destination.Address()
		// 			}
		// 		}
		// 	case "algorand":
		// 		txnBytes, err := base64.StdEncoding.DecodeString(operation.SerializedTxn)
		// 		if err != nil {
		// 			logger.Sugar().Errorw("failed to decode serialized transaction", "error", err)
		// 			return
		// 		}
		// 		var txn algorandTypes.Transaction
		// 		err = msgpack.Decode(txnBytes, &txn)
		// 		if err != nil {
		// 			logger.Sugar().Errorw("failed to deserialize transaction", "error", err)
		// 			return
		// 		}
		// 		if txn.Type == algorandTypes.PaymentTx {
		// 			destAddress = txn.PaymentTxnFields.Receiver.String()
		// 		} else if txn.Type == algorandTypes.AssetTransferTx {
		// 			destAddress = txn.AssetTransferTxnFields.AssetReceiver.String()
		// 		} else {
		// 			logger.Sugar().Errorw("Unknown transaction type", "type", txn.Type)
		// 			return
		// 		}
		// 	case "ripple":
		// 		// For Ripple, the destination is in the transaction payload
		// 		// Decode the serialized transaction
		// 		txBytes, err := hex.DecodeString(strings.TrimPrefix(operation.SerializedTxn, "0x"))
		// 		if err != nil {
		// 			logger.Sugar().Errorw("error decoding transaction", "error", err)
		// 			return
		// 		}

		// 		// Parse the transaction
		// 		var tx data.Payment
		// 		err = json.Unmarshal(txBytes, &tx)
		// 		if err != nil {
		// 			logger.Sugar().Errorw("error unmarshalling transaction", "error", err)
		// 			return
		// 		}
		// 		destAddress = tx.Destination.String()
		// 	case "cardano":
		// 		var tx cardanolib.Tx
		// 		txBytes, err := hex.DecodeString(operation.SerializedTxn)
		// 		if err != nil {
		// 			logger.Sugar().Errorw("error decoding Cardano transaction", "error", err)
		// 			return
		// 		}
		// 		if err := json.Unmarshal(txBytes, &tx); err != nil {
		// 			logger.Sugar().Errorw("error parsing Cardano transaction", "error", err)
		// 			return
		// 		}
		// 		destAddress = tx.Body.Outputs[0].Address.String()
		// 	case "sui":
		// 		var tx sui_types.TransactionData
		// 		txBytes, err := base64.StdEncoding.DecodeString(operation.SerializedTxn)
		// 		if err != nil {
		// 			logger.Sugar().Errorw("error decoding Sui transaction", "error", err)
		// 			return
		// 		}
		// 		if err := json.Unmarshal(txBytes, &tx); err != nil {
		// 			logger.Sugar().Errorw("error parsing Sui transaction", "error", err)
		// 			return
		// 		}
		// 		if len(tx.V1.Kind.ProgrammableTransaction.Inputs) < 1 {
		// 			logger.Sugar().Errorw("wrong format sui transaction", "error", err)
		// 			return
		// 		}
		// 		destAddress = string(*tx.V1.Kind.ProgrammableTransaction.Inputs[0].Pure)
		// 	}

		// 	// Verify the extracted destination matches the bridge wallet
		// 	if destAddress == "" {
		// 		logger.Sugar().Errorw("Failed to extract destination address from %s transaction", chain.ChainType)
		// 		validDestination = false
		// 	} else {
		// 		switch chain.ChainType {
		// 		case "solana":
		// 			validDestination = strings.EqualFold(destAddress, bridgeWallet.EDDSAPublicKey)
		// 		case "aptos":
		// 			validDestination = strings.EqualFold(destAddress, bridgeWallet.AptosEDDSAPublicKey)
		// 		case "stellar":
		// 			validDestination = strings.EqualFold(destAddress, bridgeWallet.StellarPublicKey)
		// 		case "algorand":
		// 			validDestination = strings.EqualFold(destAddress, bridgeWallet.AlgorandEDDSAPublicKey)
		// 		case "ripple":
		// 			validDestination = strings.EqualFold(destAddress, bridgeWallet.RippleEDDSAPublicKey)
		// 		case "cardano":
		// 			validDestination = strings.EqualFold(destAddress, bridgeWallet.CardanoPublicKey)
		// 		case "sui":
		// 			validDestination = strings.EqualFold(destAddress, bridgeWallet.SuiPublicKey)
		// 		}
		// 	}

		// 	if !validDestination {
		// 		logger.Sugar().Errorw("Invalid bridge destination address for", "chain", chain.ChainType)
		// 		return
		// 	}
		// }

		// Set message
		msg = ""
		if operation.DataToSign != nil {
			msg = *operation.DataToSign
		}
	case libs.OperationTypeSolver:
		intentBytes, err := json.Marshal(intent)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error marshalling intent: %v", err)
		}

		res, err := solver.Construct(operation.Solver, &intentBytes, operationIndex)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error constructing solver operation: %v", err)
		}

		msg = res
	case libs.OperationTypeBridgeDeposit:
		// For bridgeDeposit operations, extract transaction details from metadata
		logger.Sugar().Infow("Processing bridgeDeposit operation",
			"intentID", intent.ID,
			"operationIndex", operationIndex,
			"solverMetadata", operation.SolverMetadata)

		// Create depositOperation variable at the outer scope
		var depositOperation libs.Operation
		var needPreviousOp bool
		var tokenAddress string

		// Check if metadata is empty
		if operation.SolverMetadata == "" {
			logger.Sugar().Warnw("Bridge deposit has empty metadata, falling back to previous operation")
			needPreviousOp = true
		} else {
			// Extract the transaction hash from the operation metadata
			var metadata BridgeDepositMetadata
			err := json.Unmarshal([]byte(operation.SolverMetadata), &metadata)
			if err != nil {
				logger.Sugar().Errorw("Error unmarshalling bridge deposit metadata",
					"error", err,
					"solverMetadata", operation.SolverMetadata)
				return nil, status.Errorf(codes.Internal, "invalid bridge deposit metadata: %v", err)
			}

			logger.Sugar().Infow("Parsed bridge deposit metadata",
				"blockchainID", metadata.BlockchainID,
				"result", metadata.Result,
				"token", metadata.Token)

			// If metadata has no Result field, we need to get transaction details from previous operation
			if metadata.Result == "" {
				logger.Sugar().Warnw("No transaction hash in bridge deposit metadata, falling back to previous operation",
					"token", metadata.Token)
				needPreviousOp = true

				// Save token address for later use if present
				tokenAddress = metadata.Token
			} else {
				// Use the transaction result from metadata
				depositOperation = libs.Operation{
					BlockchainID: metadata.BlockchainID,
					Result:       metadata.Result,
				}

				logger.Sugar().Infow("Extracted bridge deposit transaction details",
					"blockchainID", depositOperation.BlockchainID,
					"txHash", depositOperation.Result)
			}
		}

		// If we need transaction details from previous operation
		if needPreviousOp {
			// If metadata is empty, try to get transaction details from previous operation
			if operationIndex == 0 {
				logger.Sugar().Errorw("No previous operation to get transaction details from")
				return nil, status.Errorf(codes.Internal, "missing transaction information")
			}

			// Get the previous operation which should be sendToBridge
			prevOp := intent.Operations[operationIndex-1]
			if prevOp.Type != libs.OperationTypeSendToBridge {
				logger.Sugar().Errorw("Previous operation is not sendToBridge",
					"prevOpType", prevOp.Type,
					"prevOpIndex", operationIndex-1)
				return nil, status.Errorf(codes.Internal, "previous operation is not sendToBridge")
			}

			// Use transaction details from previous operation
			depositOperation = libs.Operation{
				BlockchainID: prevOp.BlockchainID,
				NetworkType:  prevOp.NetworkType,
				Result:       prevOp.Result,
			}

			logger.Sugar().Infow("Using transaction details from previous operation",
				"blockchainID", depositOperation.BlockchainID,
				"txHash", depositOperation.Result,
				"savedToken", tokenAddress)
		}

		logger.Sugar().Infow("Processing bridgeDeposit for chain",
			"blockchainID", depositOperation.BlockchainID,
			"networkType", depositOperation.NetworkType,
			"txHash", depositOperation.Result)

		depositOpBlockchain, err := blockchains.GetBlockchain(depositOperation.BlockchainID, depositOperation.NetworkType)
		if err != nil {
			logger.Sugar().Errorw("error getting blockchain", "error", err)
			return nil, status.Errorf(codes.Internal, "error getting blockchain: %v", err)
		}

		transfers, err := depositOpBlockchain.GetTransfers(depositOperation.Result, &intent.Identity)
		if err != nil {
			logger.Sugar().Errorw("error getting transfers", "error", err)
			return nil, status.Errorf(codes.Internal, "error getting transfers: %v", err)
		}

		if len(transfers) == 0 {
			// If we have a token address from metadata, create a minimal transfer to proceed
			if tokenAddress != "" {
				logger.Sugar().Warnw("No transfers found but token address provided in metadata, creating minimal transfer",
					"tokenAddress", tokenAddress,
					"txHash", depositOperation.Result)

				// Create a minimal transfer with the token address from metadata
				transfers = append(transfers, common.Transfer{
					TokenAddress: tokenAddress,
					// We don't have amount information, but we need a non-empty array to proceed
					Amount:   "1", // Minimal placeholder amount
					Token:    "",  // Unknown token symbol
					IsNative: false,
					From:     intent.Identity,
					To:       "", // Unknown recipient
				})

				logger.Sugar().Infow("Created minimal transfer from metadata token",
					"tokenAddress", tokenAddress,
					"from", intent.Identity)
			} else {
				logger.Sugar().Errorw("No transfers found for bridge deposit",
					"result", depositOperation.Result,
					"identity", intent.Identity,
					"blockchainID", depositOperation.BlockchainID)
				return nil, status.Errorf(codes.Internal, "no transfers found in transaction")
			}
		}

		// check if the token exists
		transfer := transfers[0]
		srcAddress := transfer.TokenAddress

		logger.Sugar().Infow("Validating token for bridge deposit",
			"tokenAddress", srcAddress,
			"tokenSymbol", transfer.Token,
			"amount", transfer.Amount,
			"isNative", transfer.IsNative)

		chainID := depositOpBlockchain.ChainID()
		if chainID == nil {
			logger.Sugar().Errorw("Chain ID is nil", "blockchainID", depositOperation.BlockchainID)
			return nil, status.Errorf(codes.Internal, "chain ID is nil")
		}

		exists, peggedToken, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, *chainID, srcAddress)
		if err != nil {
			logger.Sugar().Errorw("Error checking token existence",
				"error", err,
				"tokenAddress", srcAddress,
				"blockchainID", depositOperation.BlockchainID)
			return nil, status.Errorf(codes.Internal, "failed to validate token")
		}

		if !exists {
			logger.Sugar().Errorw("Token does not exist for bridge deposit",
				"tokenAddress", srcAddress,
				"blockchainID", depositOperation.BlockchainID)
			return nil, status.Errorf(codes.Internal, "token does not exist")
		}

		logger.Sugar().Infow("Token exists for bridge deposit",
			"tokenAddress", srcAddress,
			"peggedToken", peggedToken)

		// Set message for signing - first try SolverDataToSign
		msg = operation.SolverDataToSign

		dataToSign := ""
		if operation.DataToSign != nil {
			dataToSign = *operation.DataToSign
		}
		// Log detailed info about the message being signed
		logger.Sugar().Infow("Processing bridge deposit signature",
			"solverDataLength", len(operation.SolverDataToSign),
			"dataToSignLength", len(dataToSign))

		// If no SolverDataToSign is provided, use DataToSign as fallback
		if len(msg) == 0 {
			logger.Sugar().Infow("Using DataToSign for bridge deposit operation", "length", len(dataToSign))
			msg = dataToSign
		}

		if len(msg) == 0 {
			logger.Sugar().Errorw("No message data available for signing bridge deposit")
			return nil, status.Errorf(codes.Internal, "no message data available for signing")
		}

		logger.Sugar().Infow("Bridge deposit message prepared for signing",
			"msgLength", len(msg),
			"transferToken", transfer.Token,
			"transferAmount", transfer.Amount)
	case libs.OperationTypeSwap:
		// Validate previous operation
		bridgeDeposit := intent.Operations[operationIndex-1]

		if operationIndex == 0 || !(bridgeDeposit.Type == libs.OperationTypeBridgeDeposit) {
			logger.Sugar().Errorw("Invalid operation type for swap")
			return nil, status.Errorf(codes.Internal, "invalid operation type for swap")
		}

		// Log detailed information about the swap operation
		logger.Sugar().Infow("Processing swap signature request",
			"solverDataToSign", operation.SolverDataToSign,
			"solverDataLength", len(operation.SolverDataToSign),
			"intent_id", intent.ID,
			"operation_id", operation.ID)

		// Set message
		msg = operation.SolverDataToSign
	case libs.OperationTypeBurn:
		// Validate nearby operations
		bridgeSwap := intent.Operations[operationIndex-1]

		if operationIndex+1 >= len(intent.Operations) || intent.Operations[operationIndex+1].Type != libs.OperationTypeWithdraw {
			logger.Sugar().Errorw("BURN operation must be followed by a WITHDRAW operation")
			return nil, status.Errorf(codes.Internal, "burn operation must be followed by a withdraw operation")
		}

		if operationIndex == 0 || !(bridgeSwap.Type == libs.OperationTypeSwap) {
			logger.Sugar().Errorw("Invalid operation type for swap")
			return nil, status.Errorf(codes.Internal, "invalid operation type for swap")
		}

		logger.Sugar().Infow("Burning tokens", "bridgeSwap", bridgeSwap)

		// Set message
		msg = operation.SolverDataToSign
	case libs.OperationTypeBurnSynthetic:
		// This operation allows direct burning of ERC20 tokens from the wallet
		// without requiring a prior swap operation
		burnSyntheticMetadata := BurnSyntheticMetadata{}

		// Verify that this operation is followed by a withdraw operation
		if operationIndex+1 >= len(intent.Operations) || intent.Operations[operationIndex+1].Type != libs.OperationTypeWithdraw {
			logger.Sugar().Errorw("BURN_SYNTHETIC operation must be followed by a WITHDRAW operation")
			return nil, status.Errorf(codes.Internal, "burn_synthetic operation must be followed by a withdraw operation")
		}

		err := json.Unmarshal([]byte(operation.SolverMetadata), &burnSyntheticMetadata)
		if err != nil {
			logger.Sugar().Errorw("Error unmarshalling burn synthetic metadata:", "error", err)
			return nil, status.Errorf(codes.Internal, "error unmarshalling burn synthetic metadata: %v", err)
		}

		// Get bridgewallet by calling /getwallet from sequencer api
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/getWallet?identity=%s&blockchain=%s", SequencerHost, intent.Identity, intent.BlockchainID), nil)
		if err != nil {
			logger.Sugar().Errorw("error creating request", "error", err)
			return nil, status.Errorf(codes.Internal, "error creating request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Sugar().Errorw("error sending request", "error", err)
			return nil, status.Errorf(codes.Internal, "error sending request: %v", err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Sugar().Errorw("error reading response body", "error", err)
			return nil, status.Errorf(codes.Internal, "error reading response body: %v", err)
		}

		var bridgeWallet db.WalletSchema
		err = json.Unmarshal(body, &bridgeWallet)
		if err != nil {
			logger.Sugar().Errorw("error unmarshalling response body", "error", err)
			return nil, status.Errorf(codes.Internal, "error unmarshalling response body: %v", err)
		}

		// Verify the user has sufficient token balance
		balance, err := ERC20.GetBalance(RPC_URL, burnSyntheticMetadata.Token, bridgeWallet.ECDSAPublicKey)
		if err != nil {
			logger.Sugar().Errorw("Error getting token balance:", "error", err)
			return nil, status.Errorf(codes.Internal, "error getting token balance: %v", err)
		}

		balanceBig, ok := new(big.Int).SetString(balance, 10)
		if !ok {
			logger.Sugar().Errorw("Error parsing balance")
			return nil, status.Errorf(codes.Internal, "error parsing balance")
		}

		amountBig, ok := new(big.Int).SetString(burnSyntheticMetadata.Amount, 10)
		if !ok {
			logger.Sugar().Errorw("Error parsing amount")
			return nil, status.Errorf(codes.Internal, "error parsing amount")
		}

		if balanceBig.Cmp(amountBig) < 0 {
			logger.Sugar().Errorw("Insufficient token balance")
			return nil, status.Errorf(codes.Internal, "insufficient token balance")
		}
		msg = operation.SolverDataToSign
	case libs.OperationTypeWithdraw:
		// Verify nearby operations
		burn := intent.Operations[operationIndex-1]

		if operationIndex == 0 || !(burn.Type == libs.OperationTypeBurn || burn.Type == libs.OperationTypeBurnSynthetic) {
			logger.Sugar().Errorw("Invalid operation type for withdraw after burn")
			return nil, status.Errorf(codes.Internal, "invalid operation type for withdraw after burn")
		}

		var withdrawMetadata WithdrawMetadata
		json.Unmarshal([]byte(operation.SolverMetadata), &withdrawMetadata)

		// Handle different burn operation types
		var tokenToWithdraw string
		var burnTokenAddress string
		if burn.Type == libs.OperationTypeBurn {
			var burnMetadata BurnMetadata
			json.Unmarshal([]byte(burn.SolverMetadata), &burnMetadata)
			tokenToWithdraw = withdrawMetadata.Token
			burnTokenAddress = burnMetadata.Token
		} else if burn.Type == libs.OperationTypeBurnSynthetic {
			var burnSyntheticMetadata BurnSyntheticMetadata
			json.Unmarshal([]byte(burn.SolverMetadata), &burnSyntheticMetadata)
			tokenToWithdraw = withdrawMetadata.Token
			burnTokenAddress = burnSyntheticMetadata.Token
		}

		chainID := opBlockchain.ChainID()
		if chainID == nil {
			logger.Sugar().Errorw("Chain ID is nil", "blockchainID", operation.BlockchainID)
			return nil, status.Errorf(codes.Internal, "chain ID is nil")
		}
		// verify these fields
		exists, destAddress, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, *chainID, tokenToWithdraw)

		if err != nil {
			logger.Sugar().Errorw("error checking token existence", "error", err)
			return nil, status.Errorf(codes.Internal, "error checking token existence: %v", err)
		}

		if !exists {
			logger.Sugar().Errorw("Token does not exist", "token", tokenToWithdraw, "blockchainID", operation.BlockchainID)
			return nil, status.Errorf(codes.Internal, "token does not exist")
		}

		if destAddress != burnTokenAddress {
			logger.Sugar().Errorw("Token mismatch", "destAddress", destAddress, "token", burnTokenAddress)
			return nil, status.Errorf(codes.Internal, "token mismatch")
		}

		// Set message
		msg = operation.SolverDataToSign
	default:
		logger.Sugar().Errorw("Unsupported operation type for signing", "type", operation.Type)
		return nil, status.Errorf(codes.InvalidArgument, "unsupported operation type: %s", operation.Type)
	}

	// Basic check if message is empty (may need refinement based on ported logic)
	if msg == "" && operation.DataToSign == nil && operation.SolverDataToSign == "" {
		logger.Sugar().Errorw("No message data could be determined for signing", "type", operation.Type, "intentID", intent.ID)
		return nil, status.Errorf(codes.FailedPrecondition, "missing message data for signing operation type %s", operation.Type)
	}
	logger.Sugar().Debugw("Message to sign determined", "msgLength", len(msg), "type", operation.Type)

	intentStr, err := identity.SanitiseIntent(*intent)
	if err != nil {
		logger.Sugar().Errorw("Failed to sanitise intent for verification", "intentID", intent.ID, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to process intent for verification: %v", err)
	}
	verified, err := identity.VerifySignature(
		intent.Identity,
		intent.BlockchainID,
		intentStr,
		intent.Signature,
	)
	if err != nil {
		logger.Sugar().Errorw("Error during intent signature verification", "intentID", intent.ID, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to verify intent signature: %v", err)
	}
	if !verified {
		logger.Sugar().Errorw("Intent signature verification failed", "intentID", intent.ID, "identity", intent.Identity)
		return nil, status.Errorf(codes.Unauthenticated, "intent signature verification failed")
	}
	logger.Sugar().Info("Intent signature verified successfully", "intentID", intent.ID)

	identity := intent.Identity
	identityCurve := intentBlockchain.KeyCurve()
	keyCurve := opBlockchain.KeyCurve()
	var msgBytes []byte

	switch operation.BlockchainID {
	case blockchains.Solana:
		msgBytes, err := base58.Decode(msg)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid base58 message for Solana: %v", err)
		}
		if operation.Type == libs.OperationTypeSwap ||
			operation.Type == libs.OperationTypeBurn ||
			operation.Type == libs.OperationTypeBurnSynthetic ||
			operation.Type == libs.OperationTypeWithdraw {
			logger.Sugar().Infow("Generating signature message for withdraw on Solana")
			go generateSignatureMessage(BridgeContractAddress, operation.BlockchainID, common.CurveEcdsa, common.CurveEddsa, msgBytes)
		} else {
			logger.Sugar().Infow("Generating signature message for other operations on Solana")
			go generateSignatureMessage(identity, operation.BlockchainID, identityCurve, keyCurve, msgBytes)
		}
	case blockchains.Bitcoin, blockchains.Dogecoin:
		// api.go seems to send raw string bytes for these
		msgBytes = []byte(msg)
	case blockchains.Sui:
		// api.go uses lib.NewBase64Data -> *[]byte
		mBytesData, err := lib.NewBase64Data(msg)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid base64 message for Sui: %v", err)
		}
		if mBytesData == nil {
			return nil, status.Error(codes.InvalidArgument, "decoded Sui message is nil")
		}
		msgBytes = *mBytesData
	case blockchains.Stellar, blockchains.Algorand:
		mBytes, err := base64.StdEncoding.DecodeString(msg)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid base64 message for %s: %v", operation.BlockchainID, err)
		}
		msgBytes = mBytes
	case blockchains.Ripple, blockchains.Cardano, blockchains.Aptos:
		mBytes, err := hex.DecodeString(msg)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid hex message: %v", err)
		}
		msgBytes = mBytes
	default:
		if blockchains.IsEVMBlockchain(operation.BlockchainID) {
			// EVM expects raw bytes of the message string
			msgBytes = []byte(msg)
		} else {
			logger.Sugar().Errorw("Unsupported blockchain for signing", "id", operation.BlockchainID)
			return nil, status.Errorf(codes.InvalidArgument, "unsupported blockchain for signing: %s", operation.BlockchainID)
		}
	}

	sigChan := make(chan Message)

	// FIXME: This is not thread safe and we shouldn't use a global map
	messageChan[msg] = sigChan

	// Determine identity for signing (special case for bridge ops on EVM)
	signingIdentity := identity
	if blockchains.IsEVMBlockchain(operation.BlockchainID) &&
		(operation.Type == libs.OperationTypeBridgeDeposit ||
			operation.Type == libs.OperationTypeSwap ||
			operation.Type == libs.OperationTypeBurn ||
			operation.Type == libs.OperationTypeBurnSynthetic ||
			operation.Type == libs.OperationTypeWithdraw) {
		signingIdentity = BridgeContractAddress
		logger.Sugar().Infow("Using BridgeContractAddress as signing identity", "address", signingIdentity)
	}

	go func() {
		defer func() {
			// Ensure channel is removed from map if goroutine panics
			// TODO: Add locking
			if r := recover(); r != nil {
				logger.Sugar().Errorw("Panic in generateSignatureMessage goroutine", "error", r)
				// We might not be able to signal back easily if the channel is gone
				// Best effort: Log it. The select will time out.
				delete(messageChan, msg)
			}
		}()
		logger.Sugar().Infow("Calling generateSignatureMessage for gRPC request", "msg", msg, "identity", signingIdentity)
		generateSignatureMessage(signingIdentity, operation.BlockchainID, identityCurve, keyCurve, msgBytes)
		// NOTE: We rely on generateSignatureMessage's logic to eventually send the result
		// back via the messageChan[chanKey].
	}()

	logger.Sugar().Infow("Waiting for signature result", "msg", msg)
	select {
	case sigResult := <-sigChan:
		// TODO: Add locking
		delete(messageChan, msg)
		logger.Sugar().Infow("Received signature result via channel", "address", sigResult.Address, "sigLen", len(sigResult.Message))

		signature := ""
		switch operation.BlockchainID {
		case blockchains.Bitcoin, blockchains.Dogecoin, blockchains.Sui:
			signature = string(sigResult.Message)
		case blockchains.Aptos, blockchains.Ripple, blockchains.Cardano, blockchains.Stellar:
			signature = hex.EncodeToString(sigResult.Message)
		case blockchains.Algorand:
			signature = base64.StdEncoding.EncodeToString(sigResult.Message)
			type algodMsg struct {
				IsRealTransaction bool
				Msg               string
			}
			m := algodMsg{IsRealTransaction: sigResult.AlgorandFlags.IsRealTransaction, Msg: msg}
			jsonBytes, err := json.Marshal(m)
			if err != nil {
				logger.Sugar().Errorw("Error marshaling algodMsg to JSON", "error", err)
				return nil, status.Errorf(codes.Internal, "error marshaling algodMsg to JSON: %v", err)
			}
			v, err := identityVerification.VerifySignature(sigResult.Address, blockchains.Algorand, string(jsonBytes), signature)
			if !v {
				logger.Sugar().Errorf("invalid signature %s, err %v", signature, err)
			}
		case blockchains.Solana:
			signature = base58.Encode(sigResult.Message)
		default:
			if blockchains.IsEVMBlockchain(operation.BlockchainID) {
				signature = string(sigResult.Message)
			} else {
				logger.Sugar().Errorw("Unexpected blockchain ID in signature response handling", "id", operation.BlockchainID)
				return nil, status.Error(codes.Internal, "internal error handling signature response")
			}
		}

		if signature == "" {
			logger.Sugar().Errorw("Empty signature received from async process", "intentID", intent.ID, "opIndex", operationIndex)
			return nil, status.Error(codes.Internal, "failed to generate signature (empty result)")
		}

		logger.Sugar().Infow("Successfully generated signature via gRPC", "intentID", intent.ID, "opIndex", operationIndex)
		return &pb.SignIntentOperationResponse{Signature: signature}, nil

	case <-ctx.Done():
		logger.Sugar().Warnw("gRPC SignIntentOperation request cancelled by client", "msg", msg)
		// Cleanup handled by defer in select or goroutine.
		// TODO: Add locking for delete if select wins race against goroutine finishing
		delete(messageChan, msg)
		return nil, status.Error(codes.Canceled, "client cancelled request")

	case <-time.After(5 * time.Minute): // Match HTTP timeout? Or configurable?
		logger.Sugar().Errorw("gRPC SignIntentOperation timed out", "msg", msg)
		// Cleanup handled by defer in select or goroutine.
		// TODO: Add locking for delete if select wins race against goroutine finishing
		delete(messageChan, msg)
		return nil, status.Error(codes.DeadlineExceeded, "signature operation timed out")
	}
}

// AddAddressDetail safely adds an address to the nested map structure
func addAddressDetail(resp *pb.GetAddressesResponse, blockchainID int32, networkType pb.NetworkType, address string) {
	if _, ok := resp.Addresses[blockchainID]; !ok {
		resp.Addresses[blockchainID] = &pb.BlockchainAddressMap{
			NetworkAddresses: make(map[int32]*pb.AddressDetail),
		}
	}
	networkTypeInt := int32(networkType)
	resp.Addresses[blockchainID].NetworkAddresses[networkTypeInt] = &pb.AddressDetail{
		NetworkType: networkType,
		Address:     address,
	}
}

// getCompressedPublicKeyBytes generates compressed ECDSA public key bytes
func getCompressedPublicKeyBytes(key *ecdsaKeygen.LocalPartySaveData) ([]byte, error) {
	if key == nil || key.ECDSAPub == nil {
		return nil, errors.New("missing ECDSA public key data")
	}
	xStr := fmt.Sprintf("%064x", key.ECDSAPub.X())
	prefix := "02"
	if key.ECDSAPub.Y().Bit(0) == 1 {
		prefix = "03"
	}
	publicKeyStr := prefix + xStr
	return hex.DecodeString(publicKeyStr)
}

// getUncompressedPublicKeyHex generates uncompressed ECDSA public key hex string (04 + X + Y)
func getUncompressedPublicKeyHex(key *ecdsaKeygen.LocalPartySaveData) (string, error) {
	if key == nil || key.ECDSAPub == nil {
		return "", errors.New("missing ECDSA public key data")
	}
	x := toHexInt(key.ECDSAPub.X())
	y := toHexInt(key.ECDSAPub.Y())
	return "04" + x + y, nil
}

// getUncompressedPublicKeyBytes generates uncompressed ECDSA public key bytes (04 + X + Y)
func getUncompressedPublicKeyBytes(key *ecdsaKeygen.LocalPartySaveData) ([]byte, error) {
	hexStr, err := getUncompressedPublicKeyHex(key)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(hexStr)
}

func startGRPCServer(port string, host host.Host, serverCertPEM, serverKeyPEM, clientCaPEM string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Sugar().Fatalf("Failed to listen for gRPC on port %s: %v", port, err)
	}

	logger.Sugar().Debug("Configuring gRPC server mTLS...")

	// Load server's certificate and private key
	// serverCert, err := tls.X509KeyPair([]byte(serverCertPEM), []byte(serverKeyPEM))
	// if err != nil {
	// 	logger.Sugar().Fatalf("Failed to load server key pair: %v", err)
	// }

	// Create a certificate pool for client CAs
	// clientCAPool := x509.NewCertPool()
	// if !clientCAPool.AppendCertsFromPEM([]byte(clientCaPEM)) {
	// 	logger.Sugar().Fatal("Failed to append client CA certificate to pool")
	// }

	// tlsConfig := &tls.Config{
	// 	Certificates: []tls.Certificate{serverCert},
	// 	ClientAuth:   tls.RequireAndVerifyClientCert,
	// 	ClientCAs:    clientCAPool,
	// 	MinVersion:   tls.VersionTLS12,
	// }

	// Create gRPC server credentials
	// serverCreds := credentials.NewTLS(tlsConfig)
	// logger.Sugar().Debug("mTLS configuration loaded successfully.")

	healthServer := health.NewServer()
	// An empty service name "" refers to the status of the overall server.
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	serverImplementation := NewValidatorServer(
		host,
	)

	// Create gRPC server with mTLS credentials
	opts := []grpc.ServerOption{
		// grpc.Creds(serverCreds),
		grpc.Creds(insecure.NewCredentials()),
	}
	s := grpc.NewServer(opts...)

	grpc_health_v1.RegisterHealthServer(s, healthServer)
	pb.RegisterValidatorServiceServer(s, serverImplementation)

	logger.Sugar().Infof("gRPC server with mTLS listening at %s", lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		logger.Sugar().Fatalf("Failed to serve gRPC: %v", err)
	}
}
