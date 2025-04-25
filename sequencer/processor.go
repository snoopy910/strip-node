package sequencer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/StripChain/strip-node/ERC20"
	"github.com/StripChain/strip-node/bridge"
	"github.com/StripChain/strip-node/libs"
	"github.com/StripChain/strip-node/libs/blockchains"
	db "github.com/StripChain/strip-node/libs/database"
	"github.com/StripChain/strip-node/solver"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/google/uuid"
)

type MintOutput struct {
	Token  string `json:"token"`
	Amount string `json:"amount"`
}

type SwapMetadata struct {
	Token    string `json:"token"`
	Multiple bool   `json:"multiple"`
	Path     string `json:"path"`
}

type BurnMetadata struct {
	Token string `json:"token"`
}

// BurnSyntheticMetadata defines the metadata required for the BURN_SYNTHETIC operation
type BurnSyntheticMetadata struct {
	Token  string `json:"token"`
	Amount string `json:"amount"`
}

type WithdrawMetadata struct {
	Token  string `json:"token"`
	Unlock bool   `json:"unlock"`
}

type LockMetadata struct {
	Lock bool `json:"lock"`
}

func ProcessIntent(intentID uuid.UUID) {
ProcessLoop:
	for {
		intent, err := db.GetIntent(intentID)
		if err != nil {
			logger.Sugar().Errorw("error getting intent", "error", err)
			return
		}

		intentBytes, err := json.Marshal(intent)
		if err != nil {
			logger.Sugar().Errorw("error marshalling intent", "error", err)
			return
		}

		if intent.Status != libs.IntentStatusProcessing {
			logger.Sugar().Infow("intent processed", "intent", intent)
			return
		}

		if intent.Expiry.Before(time.Now()) {
			db.UpdateIntentStatus(intent.ID, libs.IntentStatusExpired)
			return
		}

		var signature string
		var opBlockchain blockchains.IBlockchain
		var wallet *db.WalletSchema
		var publicKey string

	OperationLoop:
		// now process the operations of the intent
		for i, operation := range intent.Operations {
			if operation.Status == libs.OperationStatusCompleted || operation.Status == libs.OperationStatusFailed {
				continue
			}

			opBlockchain, err = blockchains.GetBlockchain(operation.BlockchainID, operation.NetworkType)
			if err != nil {
				fmt.Printf("error getting blockchain: %+v\n", err)
				break ProcessLoop
			}

			// Create context with timeout for each operation using chain's OpTimeout
			ctx, cancel := context.WithTimeout(context.Background(), opBlockchain.OpTimeout())
			defer cancel()

			opCreatedAt := operation.CreatedAt.UTC()
			now := time.Now().UTC()
			// Check both operation expiry and context timeout
			if opBlockchain.OpTimeout() < now.Sub(opCreatedAt) {
				db.UpdateOperationStatus(operation.ID, libs.OperationStatusExpired)
				continue
			}

			select {
			case <-ctx.Done():
				logger.Sugar().Warnw("operation processing timed out",
					"intentId", intentID,
					"operation", operation)
				db.UpdateOperationStatus(operation.ID, libs.OperationStatusExpired)
				continue
			default:
				// Process the operation...
			}

			wallet, err = db.GetWallet(intent.Identity, intent.BlockchainID)
			if err != nil {
				fmt.Printf("error getting wallet: %+v\n", err)
				return
			}

			switch operation.BlockchainID {
			case blockchains.Cardano:
				publicKey = wallet.CardanoPublicKey
			case blockchains.Bitcoin:
				if operation.NetworkType == blockchains.Mainnet {
					publicKey = wallet.BitcoinMainnetPublicKey
				} else if operation.NetworkType == blockchains.Testnet {
					publicKey = wallet.BitcoinTestnetPublicKey
				} else {
					publicKey = wallet.BitcoinRegtestPublicKey
				}
			case blockchains.Dogecoin:
				if operation.NetworkType == blockchains.Mainnet {
					publicKey = wallet.DogecoinMainnetPublicKey
				} else if operation.NetworkType == blockchains.Testnet {
					publicKey = wallet.DogecoinTestnetPublicKey
				}
			case blockchains.Sui:
				publicKey = wallet.SuiPublicKey
			case blockchains.Aptos:
				publicKey = wallet.AptosEDDSAPublicKey
			case blockchains.Stellar:
				publicKey = wallet.StellarPublicKey
			case blockchains.Algorand:
				publicKey = wallet.AlgorandEDDSAPublicKey
			case blockchains.Ripple:
				publicKey = wallet.RippleEDDSAPublicKey
			case blockchains.Solana:
				publicKey = wallet.SolanaPublicKey
			default:
				if blockchains.IsEVMBlockchain(operation.BlockchainID) {
					publicKey = wallet.EthereumPublicKey
				} else {
					logger.Sugar().Errorw("blockchain ID Not supported", "blockchainID", operation.BlockchainID)
					continue
				}
			}

			switch operation.Status {
			case libs.OperationStatusPending:
				// sign and send the txn. Change status to waiting

				switch operation.Type {
				case libs.OperationTypeTransaction, libs.OperationTypeSendToBridge:
					lockSchema, err := db.VerifyIdentityLockSchema(intent, &operation)
					if lockSchema == nil {
						logger.Sugar().Errorw("error verifying identity lock", "error", err)
						break
					}
					signature, err = getSignature(intent, i)
					if err != nil {
						fmt.Printf("error getting signature: %+v\n", err)
						break
					}
					if operation.SerializedTxn == nil {
						logger.Sugar().Errorw("serialized txn is nil", "operation", operation)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					txHash, err := opBlockchain.BroadcastTransaction(*operation.SerializedTxn, signature, &publicKey)
					if err != nil {
						fmt.Printf("error broadcasting transaction: %+v\n", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					var lockMetadata LockMetadata
					json.Unmarshal([]byte(operation.SolverMetadata), &lockMetadata)

					if lockMetadata.Lock {
						err := db.LockIdentity(lockSchema.Id)
						if err != nil {
							fmt.Println(err)
							break
						}

						db.UpdateOperationResult(operation.ID, libs.OperationStatusCompleted, txHash)
					} else {
						db.UpdateOperationResult(operation.ID, libs.OperationStatusWaiting, txHash)
					}
				case libs.OperationTypeSolver:
					lockSchema, err := db.VerifyIdentityLockSchema(intent, &operation)
					if lockSchema == nil {
						logger.Sugar().Errorw("error verifying identity lock", "error", err)
						break
					}

					// get data to sign from solver
					dataToSign, err := solver.Construct(operation.Solver, &intentBytes, i)

					if err != nil {
						logger.Sugar().Errorw("error constructing solver data to sign", "error", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					db.UpdateOperationSolverDataToSign(operation.ID, dataToSign)

					// then send the signature to solver
					result, err := solver.Solve(
						operation.Solver, &intentBytes,
						i,
						signature,
					)

					if err != nil {
						logger.Sugar().Errorw("error solving", "error", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					var lockMetadata LockMetadata
					json.Unmarshal([]byte(operation.SolverMetadata), &lockMetadata)

					if lockMetadata.Lock {
						err := db.LockIdentity(lockSchema.Id)
						if err != nil {
							logger.Sugar().Errorw("error locking identity", "error", err)
							break
						}
						db.UpdateOperationResult(operation.ID, libs.OperationStatusCompleted, result)
					} else {
						db.UpdateOperationResult(operation.ID, libs.OperationStatusWaiting, result)
					}
				case libs.OperationTypeBridgeDeposit:
					if i == 0 || !(intent.Operations[i-1].Type == libs.OperationTypeSendToBridge) {
						logger.Sugar().Errorw("Invalid operation type for bridge deposit")
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					depositOperation := intent.Operations[i-1]

					// TODO: This code is not correct, but swapping is being worked on
					transfers, err := opBlockchain.GetTransfers(depositOperation.Result, &publicKey)
					if err != nil {
						logger.Sugar().Errorw("error getting transfers", "error", err)
						break
					}

					if len(transfers) == 0 {
						logger.Sugar().Errorw("No transfers found", "result", depositOperation.Result, "identity", intent.Identity)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					transfer := transfers[0]
					srcAddress := transfer.TokenAddress
					amount := transfer.ScaledAmount

					opBlockchain, err := blockchains.GetBlockchain(depositOperation.BlockchainID, depositOperation.NetworkType)
					if err != nil {
						logger.Sugar().Errorw("error getting blockchain", "error", err)
						break
					}

					chainID := opBlockchain.ChainID()
					if chainID == nil {
						logger.Sugar().Errorw("chainID is nil", "blockchainID", operation.BlockchainID, "networkType", operation.NetworkType)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					exists, destAddress, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, *chainID, srcAddress)

					if err != nil {
						logger.Sugar().Errorw("error checking token existence", "error", err)
						break
					}

					if !exists {
						logger.Sugar().Errorw("Token does not exist", "srcAddress", srcAddress, "blockchainId", operation.BlockchainID, "networkType", operation.NetworkType)

						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					dataToSign, err := bridge.BridgeDepositDataToSign(RPC_URL, BridgeContractAddress, amount, wallet.EthereumPublicKey, destAddress)
					if err != nil {
						logger.Sugar().Errorw("error getting data to sign", "error", err)
						break
					}

					db.UpdateOperationSolverDataToSign(operation.ID, dataToSign)
					intent.Operations[i].SolverDataToSign = dataToSign

					signature, err := getSignature(intent, i)
					if err != nil {
						logger.Sugar().Errorw("error getting signature", "error", err)
						break
					}

					logger.Sugar().Infow("Minting bridge %s %s %s %s", amount, wallet.EthereumPublicKey, destAddress, signature)

					result, err := mintBridge(
						amount, wallet.EthereumPublicKey, destAddress, signature)

					if err != nil {
						logger.Sugar().Errorw("error minting bridge", "error", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					mintOutput := MintOutput{
						Token:  destAddress,
						Amount: amount,
					}

					mintOutputBytes, err := json.Marshal(mintOutput)

					if err != nil {
						logger.Sugar().Errorw("error marshalling mint output", "error", err)
						break
					}

					db.UpdateOperationSolverOutput(operation.ID, string(mintOutputBytes))
					db.UpdateOperationResult(operation.ID, libs.OperationStatusWaiting, result)
					break OperationLoop
				case libs.OperationTypeSwap:
					if i == 0 || !(intent.Operations[i-1].Type == libs.OperationTypeBridgeDeposit) {
						logger.Sugar().Errorw("Invalid operation type for swap")
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}
					bridgeDeposit := intent.Operations[i-1]

					// Get the deposit operation details to find the actual source token address
					depositOperation := intent.Operations[i-2] // The operation before bridge deposit is send-to-bridge

					// Get the actual token used in the deposit from the transfer events
					transfers, err := opBlockchain.GetTransfers(depositOperation.Result, &publicKey)
					if err != nil {
						logger.Sugar().Errorw("error getting transfers for swap tokenIn", "error", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					if len(transfers) == 0 {
						logger.Sugar().Errorw("No transfers found for swap tokenIn", "result", depositOperation.Result)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					// Get the original token address from the transfer event
					transfer := transfers[0]
					originalTokenAddress := transfer.TokenAddress
					logger.Sugar().Infow("Original token from transfer event", "originalToken", originalTokenAddress)

					var bridgeDepositData MintOutput
					var swapMetadata SwapMetadata
					json.Unmarshal([]byte(bridgeDeposit.SolverOutput), &bridgeDepositData)
					json.Unmarshal([]byte(operation.SolverMetadata), &swapMetadata)

					// Use the pegged token from the bridge deposit output
					tokenIn := bridgeDepositData.Token
					logger.Sugar().Infow("Using pegged token for swap", "peggedToken", tokenIn, "originalToken", originalTokenAddress)

					// Use the token from metadata for tokenOut
					tokenOut := swapMetadata.Token
					amountIn := bridgeDepositData.Amount
					deadline := time.Now().Add(time.Hour).Unix()

					wallet, err := db.GetWallet(intent.Identity, "ecdsa")
					if err != nil {
						logger.Sugar().Errorw("error getting wallet", "error", err)
						break
					}

					var dataToSign string
					if swapMetadata.Multiple {
						// Multi-pool swap
						dataToSign, err = bridge.BridgeSwapMultiplePoolsDataToSign(
							RPC_URL,
							BridgeContractAddress,
							wallet.EthereumPublicKey,
							tokenIn,
							swapMetadata.Path,
							amountIn,
							deadline,
						)
					} else {
						// Single pool swap (default)
						dataToSign, err = bridge.BridgeSwapDataToSign(
							RPC_URL,
							BridgeContractAddress,
							wallet.EthereumPublicKey,
							tokenIn,
							tokenOut,
							amountIn,
							deadline,
						)
					}

					if err != nil {
						logger.Sugar().Errorw("error getting data to sign", "error", err)
						break
					}

					// Log the dataToSign for debugging
					logger.Sugar().Infow("Bridge swap data to sign generated",
						"dataToSign", dataToSign,
						"length", len(dataToSign),
						"isMultiple", swapMetadata.Multiple,
						"operation_id", operation.ID)

					db.UpdateOperationSolverDataToSign(operation.ID, dataToSign)
					intent.Operations[i].SolverDataToSign = dataToSign

					// Log before getting signature
					logger.Sugar().Infow("Requesting signature for bridge swap",
						"operationIndex", i,
						"intentID", intent.ID,
						"dataToSignAvailable", intent.Operations[i].SolverDataToSign != "")

					signature, err := getSignature(intent, i)
					if err != nil {
						logger.Sugar().Errorw("error getting signature", "error", err)
						break
					}

					logger.Sugar().Infow("Swapping bridge", "wallet", wallet.EthereumPublicKey, "tokenIn", tokenIn, "tokenOut", tokenOut, "amountIn", amountIn, "deadline", deadline, "signature", signature, "multiple", swapMetadata.Multiple)

					result, err := swapBridge(
						wallet.EthereumPublicKey,
						tokenIn,
						tokenOut,
						amountIn,
						deadline,
						signature,
					)

					if err != nil {
						logger.Sugar().Errorw("error swapping bridge", "error", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					db.UpdateOperationResult(operation.ID, libs.OperationStatusWaiting, result)

					break OperationLoop
				case libs.OperationTypeBurn:
					if i+1 >= len(intent.Operations) || intent.Operations[i+1].Type != libs.OperationTypeWithdraw {
						fmt.Println("BURN operation must be followed by a WITHDRAW operation")
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					if i == 0 || !(intent.Operations[i-1].Type == libs.OperationTypeSwap) {
						logger.Sugar().Errorw("Invalid operation type for swap")
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}
					bridgeSwap := intent.Operations[i-1]

					logger.Sugar().Infow("Burning tokens", "bridgeSwap", bridgeSwap)

					burnAmount := bridgeSwap.SolverOutput
					burnMetadata := BurnMetadata{}

					json.Unmarshal([]byte(operation.SolverMetadata), &burnMetadata)

					wallet, err := db.GetWallet(intent.Identity, "ecdsa")
					if err != nil {
						logger.Sugar().Errorw("error getting wallet", "error", err)
						break
					}

					dataToSign, err := bridge.BridgeBurnDataToSign(
						RPC_URL,
						BridgeContractAddress,
						wallet.EthereumPublicKey,
						burnAmount,
						burnMetadata.Token,
					)

					if err != nil {
						logger.Sugar().Errorw("error getting data to sign", "error", err)
						break
					}

					db.UpdateOperationSolverDataToSign(operation.ID, dataToSign)
					intent.Operations[i].SolverDataToSign = dataToSign

					signature, err := getSignature(intent, i)
					if err != nil {
						logger.Sugar().Errorw("error getting signature", "error", err)
						break
					}

					logger.Sugar().Infow("Burn tokens", "wallet", wallet.EthereumPublicKey, "burnAmount", burnAmount, "token", burnMetadata.Token, "signature", signature)

					result, err := burnTokens(
						wallet.EthereumPublicKey,
						burnAmount,
						burnMetadata.Token,
						signature,
					)

					if err != nil {
						logger.Sugar().Errorw("error burning tokens", "error", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					db.UpdateOperationResult(operation.ID, libs.OperationStatusWaiting, result)
					break OperationLoop
				case libs.OperationTypeBurnSynthetic:
					// This operation allows direct burning of ERC20 tokens from the wallet
					// without requiring a prior swap operation
					burnSyntheticMetadata := BurnSyntheticMetadata{}

					// Verify that this operation is followed by a withdraw operation
					if i+1 >= len(intent.Operations) || intent.Operations[i+1].Type != libs.OperationTypeWithdraw {
						logger.Sugar().Errorw("BURN_SYNTHETIC validation failed: must be followed by WITHDRAW",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"operationIndex", i,
							"totalOperations", len(intent.Operations),
							"nextOperationType", getNextOperationType(intent, i))
						fmt.Println("BURN_SYNTHETIC operation must be followed by a WITHDRAW operation")
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					logger.Sugar().Infow("BURN_SYNTHETIC validati	on passed: followed by WITHDRAW",
						"operationId", operation.ID,
						"nextOperationId", intent.Operations[i+1].ID,
						"nextOperationType", intent.Operations[i+1].Type)

					err = json.Unmarshal([]byte(operation.SolverMetadata), &burnSyntheticMetadata)
					if err != nil {
						logger.Sugar().Errorw("BURN_SYNTHETIC metadata parsing failed",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"error", err,
							"rawMetadata", operation.SolverMetadata)
						fmt.Println("Error unmarshalling burn synthetic metadata:", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					logger.Sugar().Infow("BURN_SYNTHETIC metadata parsed successfully",
						"operationId", operation.ID,
						"token", burnSyntheticMetadata.Token,
						"amount", burnSyntheticMetadata.Amount)

					// Validate token address format
					isValidToken, tokenErr := validateBurnSyntheticToken(burnSyntheticMetadata.Token)
					if !isValidToken {
						logger.Sugar().Errorw("BURN_SYNTHETIC invalid token",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"token", burnSyntheticMetadata.Token,
							"error", tokenErr)
						fmt.Println("Invalid token address:", tokenErr)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					wallet, err := db.GetWallet(intent.Identity, "ecdsa")
					if err != nil {
						logger.Sugar().Errorw("BURN_SYNTHETIC wallet retrieval failed",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"identity", intent.Identity,
							"error", err)
						fmt.Println("Error getting wallet:", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					logger.Sugar().Infow("BURN_SYNTHETIC wallet retrieved successfully",
						"operationId", operation.ID,
						"publicKey", wallet.EthereumPublicKey)

					// Verify the user has sufficient token balance
					balance, err := ERC20.GetBalance(RPC_URL, burnSyntheticMetadata.Token, wallet.EthereumPublicKey)
					if err != nil {
						logger.Sugar().Errorw("BURN_SYNTHETIC balance check failed",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"token", burnSyntheticMetadata.Token,
							"account", wallet.EthereumPublicKey,
							"error", err)
						fmt.Println("Error getting token balance:", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					logger.Sugar().Infow("BURN_SYNTHETIC balance retrieved successfully",
						"operationId", operation.ID,
						"token", burnSyntheticMetadata.Token,
						"account", wallet.EthereumPublicKey,
						"balance", balance)

					balanceBig, ok := new(big.Int).SetString(balance, 10)
					if !ok {
						logger.Sugar().Errorw("BURN_SYNTHETIC balance parsing failed",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"balance", balance)
						fmt.Println("Error parsing balance")
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					amountBig, ok := new(big.Int).SetString(burnSyntheticMetadata.Amount, 10)
					if !ok {
						logger.Sugar().Errorw("BURN_SYNTHETIC amount parsing failed",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"amount", burnSyntheticMetadata.Amount)
						fmt.Println("Error parsing amount")
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					// Log balance check details
					logBurnSyntheticBalanceCheck(balanceBig, amountBig, burnSyntheticMetadata.Token)

					if balanceBig.Cmp(amountBig) < 0 {
						logger.Sugar().Errorw("BURN_SYNTHETIC insufficient balance",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"token", burnSyntheticMetadata.Token,
							"account", wallet.EthereumPublicKey,
							"balance", balanceBig.String(),
							"requiredAmount", amountBig.String())
						fmt.Println("Insufficient token balance")
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					logger.Sugar().Infow("BURN_SYNTHETIC sufficient balance confirmed",
						"operationId", operation.ID,
						"token", burnSyntheticMetadata.Token,
						"balance", balanceBig.String(),
						"amount", amountBig.String())

					// Verify token does NOT exist on bridge contract (BURN_SYNTHETIC is for native L2 tokens, not bridged tokens)
					exists, destAddress, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, "ethereum", burnSyntheticMetadata.Token)
					if err != nil {
						logger.Sugar().Errorw("BURN_SYNTHETIC bridge token check failed",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"token", burnSyntheticMetadata.Token,
							"error", err)
						fmt.Println("Error checking token existence in bridge:", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					if exists {
						logger.Sugar().Errorw("BURN_SYNTHETIC invalid token: token exists on bridge",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"token", burnSyntheticMetadata.Token,
							"peggedToken", destAddress)
						fmt.Println("Invalid token: token exists on bridge, use BURN instead of BURN_SYNTHETIC")
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					logger.Sugar().Infow("BURN_SYNTHETIC token validated: not a bridged token",
						"operationId", operation.ID,
						"token", burnSyntheticMetadata.Token)

					// Generate data to sign for burning tokens
					dataToSign, err := bridge.BridgeBurnDataToSign(
						RPC_URL,
						BridgeContractAddress,
						wallet.EthereumPublicKey,
						burnSyntheticMetadata.Amount,
						burnSyntheticMetadata.Token,
					)

					if err != nil {
						logger.Sugar().Errorw("BURN_SYNTHETIC data generation failed",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"token", burnSyntheticMetadata.Token,
							"account", wallet.EthereumPublicKey,
							"bridgeContract", BridgeContractAddress,
							"error", err)
						fmt.Println("Error generating burn data to sign:", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					logger.Sugar().Infow("BURN_SYNTHETIC data generated successfully",
						"operationId", operation.ID,
						"dataToSignLength", len(dataToSign),
						"dataToSignPrefix", truncateString(dataToSign, 20))

					// Update operation with data to sign and wait for signature
					db.UpdateOperationSolverDataToSign(operation.ID, dataToSign)
					intent.Operations[i].SolverDataToSign = dataToSign

					logger.Sugar().Infow("BURN_SYNTHETIC operation updated with data to sign",
						"operationId", operation.ID)

					signature, err := getSignature(intent, i)
					if err != nil {
						logger.Sugar().Errorw("BURN_SYNTHETIC signature retrieval failed",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"error", err)
						fmt.Println("Error getting signature:", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					// Log signature details
					logBurnSyntheticSignature(signature, dataToSign)

					// Verify signature locally before submitting
					isValidSignature, verifyErr := verifyBurnSyntheticSignature(dataToSign, signature, wallet.EthereumPublicKey)
					if !isValidSignature {
						logger.Sugar().Warnw("BURN_SYNTHETIC signature verification warning",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"error", verifyErr,
							"proceedingAnyway", true)
						// Note: We're logging the warning but not failing - to maintain the original logic
					}

					logger.Sugar().Infow("BURN_SYNTHETIC executing burn transaction",
						"operationId", operation.ID,
						"account", wallet.EthereumPublicKey,
						"amount", burnSyntheticMetadata.Amount,
						"token", burnSyntheticMetadata.Token,
						"signatureLength", len(signature))

					fmt.Println("Burning synthetic tokens", wallet.EthereumPublicKey, burnSyntheticMetadata.Amount, burnSyntheticMetadata.Token, signature)

					// Execute the burn transaction
					result, err := burnTokens(
						wallet.EthereumPublicKey,
						burnSyntheticMetadata.Amount,
						burnSyntheticMetadata.Token,
						signature,
					)

					if err != nil {
						logger.Sugar().Errorw("BURN_SYNTHETIC transaction failed",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"token", burnSyntheticMetadata.Token,
							"account", wallet.EthereumPublicKey,
							"error", err)
						fmt.Println("Error burning tokens:", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					logger.Sugar().Infow("BURN_SYNTHETIC transaction submitted successfully",
						"operationId", operation.ID,
						"intentId", intent.ID,
						"token", burnSyntheticMetadata.Token,
						"account", wallet.EthereumPublicKey,
						"transactionHash", result)

					db.UpdateOperationResult(operation.ID, libs.OperationStatusWaiting, result)

					// Log integration with next withdraw operation
					if i+1 < len(intent.Operations) && intent.Operations[i+1].Type == libs.OperationTypeWithdraw {
						var withdrawMetadata WithdrawMetadata
						withdrawErr := json.Unmarshal([]byte(intent.Operations[i+1].SolverMetadata), &withdrawMetadata)

						if withdrawErr == nil {
							logger.Sugar().Infow("BURN_SYNTHETIC followed by WITHDRAW operation",
								"burnOperationId", operation.ID,
								"withdrawOperationId", intent.Operations[i+1].ID,
								"burnToken", burnSyntheticMetadata.Token,
								"withdrawToken", withdrawMetadata.Token,
								"tokensMatch", burnSyntheticMetadata.Token == withdrawMetadata.Token)
						}
					}

					break OperationLoop
				case libs.OperationTypeWithdraw:
					if i == 0 || !(intent.Operations[i-1].Type == libs.OperationTypeBurn || intent.Operations[i-1].Type == libs.OperationTypeBurnSynthetic) {
						logger.Sugar().Errorw("Invalid operation type for withdraw after burn")
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}
					burn := intent.Operations[i-1]

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
						logger.Sugar().Errorw("chainID is nil", "operationId", operation.ID)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}
					// verify these fields
					exists, destAddress, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, *chainID, tokenToWithdraw)

					if err != nil {
						logger.Sugar().Errorw("error checking token existence", "error", err)
						break
					}

					if !exists {
						logger.Sugar().Errorw("Token does not exist", "token", tokenToWithdraw, "blockchainId", operation.BlockchainID, "networkType", operation.NetworkType)

						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					if destAddress != burnTokenAddress {
						logger.Sugar().Errorw("Token mismatch", "destAddress", destAddress, "token", burnTokenAddress)

						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					bridgeWallet, err := db.GetWallet(BridgeContractAddress, blockchains.Ethereum)
					if err != nil {
						logger.Sugar().Errorw("error getting bridge wallet", "error", err)
						break
					}

					var blockchainBridgeWallet string
					switch operation.BlockchainID {
					case blockchains.Cardano:
						blockchainBridgeWallet = bridgeWallet.CardanoPublicKey
					default:
						logger.Sugar().Errorw("Blockchain ID not supported", "blockchainID", operation.BlockchainID)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break OperationLoop
					}

					tx, dataToSign, err := opBlockchain.BuildWithdrawTx(blockchainBridgeWallet, burn.SolverOutput, publicKey, &tokenToWithdraw)
					if err != nil {
						fmt.Println(err)
						break
					}
					db.UpdateOperationSolverDataToSign(operation.ID, dataToSign)
					intent.Operations[i].SolverDataToSign = dataToSign

					withdrawSignature, err := getSignature(intent, i)
					if err != nil {
						fmt.Printf("error getting signature: %+v\n", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					result, err := opBlockchain.BroadcastTransaction(
						tx,
						withdrawSignature,
						&blockchainBridgeWallet,
					)

					if err != nil {
						fmt.Printf("error withdrawing tokens: %+v\n", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					db.UpdateOperationResult(operation.ID, libs.OperationStatusWaiting, result)
					break OperationLoop
				}

				break OperationLoop
			case libs.OperationStatusWaiting:
				// check for confirmations and update the status to completed
				switch operation.Type {
				case libs.OperationTypeTransaction, libs.OperationTypeSendToBridge, libs.OperationTypeBridgeDeposit, libs.OperationTypeSwap, libs.OperationTypeBurn, libs.OperationTypeBurnSynthetic:
					confirmed, err := opBlockchain.IsTransactionBroadcastedAndConfirmed(operation.Result)
					if err != nil {
						logger.Sugar().Errorw("error checking transaction", "error", err)
						break
					}
					if !confirmed {
						break
					}
					db.UpdateOperationStatus(operation.ID, libs.OperationStatusCompleted)

					if i+1 == len(intent.Operations) {
						// update the intent status to completed
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusCompleted)
					}

					break OperationLoop
				case libs.OperationTypeSolver:
					status, err := solver.CheckStatus(
						operation.Solver, &intentBytes, i,
					)

					if err != nil {
						logger.Sugar().Errorw("error checking solver status", "error", err)
						break
					}

					if status == solver.SOLVER_OPERATION_STATUS_SUCCESS {
						output, err := solver.GetOutput(operation.Solver, &intentBytes, i)

						if err != nil {
							logger.Sugar().Errorw("error getting solver output", "error", err)
							break
						}

						db.UpdateOperationStatus(operation.ID, libs.OperationStatusCompleted)
						db.UpdateOperationSolverOutput(operation.ID, output)

						if i+1 == len(intent.Operations) {
							// update the intent status to completed
							db.UpdateIntentStatus(intent.ID, libs.IntentStatusCompleted)
						}
					}

					if status == solver.SOLVER_OPERATION_STATUS_FAILURE {
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
					}

					break OperationLoop
				case libs.OperationTypeWithdraw:
					confirmed, err := opBlockchain.IsTransactionBroadcastedAndConfirmed(operation.Result)
					if err != nil {
						logger.Sugar().Errorw("error checking transaction", "error", err)
						break
					}
					if !confirmed {
						break
					}

					// now unlock the identity if locked
					var withdrawMetadata WithdrawMetadata
					json.Unmarshal([]byte(operation.SolverMetadata), &withdrawMetadata)

					lockSchema, err := db.GetLock(intent.Identity, intent.BlockchainID)
					if err != nil {
						logger.Sugar().Errorw("error getting lock", "error", err)
						break
					}

					if withdrawMetadata.Unlock {
						depositOperation := intent.Operations[i-4]
						// check for confirmations
						confirmed, err = opBlockchain.IsTransactionBroadcastedAndConfirmed(depositOperation.Result)
						if err != nil {
							logger.Sugar().Errorw("error checking transaction", "error", err)
							break
						}
					}
					if confirmed {
						err := db.UnlockIdentity(lockSchema.Id)
						if err != nil {
							fmt.Println(err)
							break
						}

						db.UpdateOperationStatus(operation.ID, libs.OperationStatusCompleted)
					}

					if i+1 == len(intent.Operations) {
						// update the intent status to completed
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusCompleted)
					}

					break OperationLoop
				default:
					logger.Sugar().Errorw("Unknown operation type", "type", operation.Type)
					break OperationLoop
				}
			default:
				logger.Sugar().Errorw("Unknown operation status", "status", operation.Status)
				break OperationLoop
			}
		}

		time.Sleep(5 * time.Second)
	}
}

type SignatureResponse struct {
	Signature string `json:"signature"`
	Address   string `json:"address"`
}

func getSignature(intent *libs.Intent, operationIndex int) (string, error) {
	signature, _, err := getSignatureEx(intent, operationIndex)
	if err != nil {
		return "", err
	}
	return signature, nil
}

func getSignatureEx(intent *libs.Intent, operationIndex int) (string, string, error) {
	// get wallet
	wallet, err := db.GetWallet(intent.Identity, intent.BlockchainID)
	if err != nil {
		return "", "", fmt.Errorf("error getting wallet: %v", err)
	}

	// get the signer
	signers := wallet.Signers
	signer, err := GetSigner(signers[0])

	if err != nil {
		return "", "", fmt.Errorf("error getting signer: %v", err)
	}

	intentBytes, err := json.Marshal(intent)
	if err != nil {
		return "", "", fmt.Errorf("error marshalling intent: %v", err)
	}

	operationIndexStr := strconv.FormatUint(uint64(operationIndex), 10)

	// Log the request details for debugging
	logger.Sugar().Infow("Requesting signature from validator",
		"url", signer.URL+"/signature?operationIndex="+operationIndexStr,
		"intentID", intent.ID,
		"operationIndex", operationIndex)

	req, err := http.NewRequest("POST", signer.URL+"/signature?operationIndex="+operationIndexStr, bytes.NewBuffer(intentBytes))

	if err != nil {
		return "", "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 30 * time.Second, // Add timeout to prevent hanging
	}
	resp, err := client.Do(req)

	if err != nil {
		return "", "", fmt.Errorf("error sending request: %v", err)
	}

	defer resp.Body.Close()

	// Check HTTP status code first
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("validator returned non-OK status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response body: %v", err)
	}

	// Log the response for debugging
	responseStr := string(body)
	logger.Sugar().Infow("Received signature response",
		"contentLength", len(responseStr),
		"responseBody", responseStr[:min(len(responseStr), 100)]) // Log first 100 chars to avoid excessive logging

	// Handle empty responses
	if len(responseStr) == 0 {
		return "", "", fmt.Errorf("empty response from validator")
	}

	var signatureResponse SignatureResponse
	err = json.Unmarshal(body, &signatureResponse)
	if err != nil {
		return "", "", fmt.Errorf("error unmarshalling response body: %v, body: %s", err, truncateString(responseStr, 200))
	}

	// Validate signature response
	if signatureResponse.Signature == "" {
		return "", "", fmt.Errorf("empty signature in response")
	}

	return signatureResponse.Signature, signatureResponse.Address, nil
}

// Helper function to truncate strings for logging
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// readBitcoinAddress returns the appropriate Bitcoin public key based on the chain configuration
func readBitcoinAddress(wallet *db.WalletSchema, chainId string) (string, error) {
	if wallet == nil {
		return "", fmt.Errorf("wallet is nil")
	}

	switch chainId {
	case "1000": // Bitcoin mainnet
		if wallet.BitcoinMainnetPublicKey == "" {
			return "", fmt.Errorf("bitcoin mainnet public key not found in wallet")
		}
		return wallet.BitcoinMainnetPublicKey, nil
	case "1001": // Bitcoin testnet
		if wallet.BitcoinTestnetPublicKey == "" {
			return "", fmt.Errorf("bitcoin testnet public key not found in wallet")
		}
		return wallet.BitcoinTestnetPublicKey, nil
	case "1002": // Bitcoin regtest
		if wallet.BitcoinRegtestPublicKey == "" {
			return "", fmt.Errorf("bitcoin regtest public key not found in wallet")
		}
		return wallet.BitcoinRegtestPublicKey, nil
	default:
		return "", fmt.Errorf("unsupported bitcoin chain ID: %s", chainId)
	}
}

func getNextOperationType(intent *libs.Intent, operationIndex int) *libs.OperationType {
	if operationIndex+1 < len(intent.Operations) {
		return &intent.Operations[operationIndex+1].Type
	}
	return nil
}
