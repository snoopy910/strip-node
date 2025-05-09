package sequencer

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/StripChain/strip-node/ERC20"
	"github.com/StripChain/strip-node/bridge"
	"github.com/StripChain/strip-node/libs"
	"github.com/StripChain/strip-node/libs/blockchains"
	db "github.com/StripChain/strip-node/libs/database"
	pb "github.com/StripChain/strip-node/libs/proto"
	"github.com/StripChain/strip-node/solver"
	solversregistry "github.com/StripChain/strip-node/solversRegistry"
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
			logger.Sugar().Infow("intent expired", "intent", intent)
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
				logger.Sugar().Errorw("error getting blockchain", "error", err)
				break ProcessLoop
			}

			// Create context with timeout for each operation using chain's OpTimeout: Timeout logic needs to be reviewed TO DO later
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

			// then get the data signed
			if operation.Type != libs.OperationTypeBridgeDeposit && operation.Type != libs.OperationTypeBurnSynthetic && operation.Type != libs.OperationTypeWithdraw { // for bridge deposit the dataToSign is set later
				signature, err = getSignature(intent, i)
				if err != nil {
					logger.Sugar().Errorw("error getting signature", "error", err)
					break
				}
			}

			wallet, err = db.GetWallet(intent.Identity, intent.BlockchainID)
			if err != nil {
				logger.Sugar().Errorw("error getting wallet", "error", err)
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
						logger.Sugar().Errorw("error getting signature", "error", err)
						break
					}
					if operation.SerializedTxn == nil {
						logger.Sugar().Errorw("serialized txn is nil", "operation", operation)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					if operation.Type == libs.OperationTypeSendToBridge {
						// Get bridge wallet for the chain
						bridgeWallet, err := db.GetWallet(BridgeContractAddress, blockchains.Ethereum)
						if err != nil {
							logger.Sugar().Errorw("Failed to get bridge wallet", "error", err)
							db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
							db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
							break
						}
						isVerified, err := verifyDestinationAddress(bridgeWallet, &operation)
						if err != nil {
							logger.Sugar().Errorw("Failed to verify destination address", "error", err)
							db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
							db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
							break
						}
						if !isVerified {
							logger.Sugar().Errorw("Bridge address not verified", "error", err)
							db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
							db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
							break
						}
					}

					txHash, err := opBlockchain.BroadcastTransaction(*operation.SerializedTxn, signature, &publicKey)
					if err != nil {
						logger.Sugar().Errorw("error broadcasting transaction", "error", err)
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

						db.UpdateOperationResult(operation.ID, libs.OperationStatusCompleted, txHash)
					} else {
						db.UpdateOperationResult(operation.ID, libs.OperationStatusWaiting, txHash)
					}
				case libs.OperationTypeSolver:
					solverExists, err := solversregistry.SolverExistsAndWhitelisted(RPC_URL, SolversRegistryContractAddress, operation.Solver)
					if err != nil {
						logger.Sugar().Errorw("error checking if solver exists", "error", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					if !solverExists {
						logger.Sugar().Errorw("solver is not registered or not whitelisted", "solver", operation.Solver)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break OperationLoop
					}

					chainID := opBlockchain.ChainID()
					if chainID == nil {
						logger.Sugar().Errorw("chainID is nil", "blockchainID", operation.BlockchainID, "networkType", operation.NetworkType)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}
					validChain, err := solversregistry.ValidateChain(RPC_URL, SolversRegistryContractAddress, operation.Solver, *chainID)
					if err != nil {
						logger.Sugar().Errorw("error validating chain", "error", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}
					if !validChain {
						logger.Sugar().Errorw("chain is not valid", "chain", *chainID)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break OperationLoop
					}

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
					depositOpBlockchain, err := blockchains.GetBlockchain(depositOperation.BlockchainID, depositOperation.NetworkType)
					if err != nil {
						logger.Sugar().Errorw("error getting blockchain", "error", err)
						break
					}
					// TODO: This code is not correct, but swapping is being worked on
					transfers, err := depositOpBlockchain.GetTransfers(depositOperation.Result, &publicKey)
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

					chainID := depositOpBlockchain.ChainID()
					if chainID == nil {
						logger.Sugar().Errorw("chainID is nil", "blockchainID", depositOperation.BlockchainID, "networkType", depositOperation.NetworkType)
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
						logger.Sugar().Errorw("Token does not exist", "srcAddress", srcAddress, "blockchainId", depositOperation.BlockchainID, "networkType", depositOperation.NetworkType)

						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					wallet, err := db.GetWallet(intent.Identity, intent.BlockchainID)
					if err != nil {
						logger.Sugar().Errorw("error getting wallet", "error", err)
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

					wallet, err := db.GetWallet(intent.Identity, blockchains.Ethereum)
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

					var result string
					if swapMetadata.Multiple {
						logger.Sugar().Infow("Using multiple pools swap function", "path", swapMetadata.Path)
						result, err = swapMultiplePoolsBridge(
							wallet.ECDSAPublicKey,
							tokenIn,
							swapMetadata.Path,
							amountIn,
							deadline,
							signature,
						)
					} else {
						logger.Sugar().Infow("Using single pool swap function")
						result, err = swapBridge(
							wallet.ECDSAPublicKey,
							tokenIn,
							tokenOut,
							amountIn,
							deadline,
							signature,
						)
					}

					if err != nil {
						logger.Sugar().Errorw("error swapping bridge", "error", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					db.UpdateOperationResult(operation.ID, libs.OperationStatusWaiting, result)
					// Also update the SolverOutput with the amount - this is critical for the burn operation
					db.UpdateOperationSolverOutput(operation.ID, amountIn)

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

					// Validate burnAmount is not empty
					if burnAmount == "" {
						logger.Sugar().Errorw("Empty burn amount in solver output",
							"bridgeSwapID", bridgeSwap.ID,
							"result", bridgeSwap.Result,
							"solverOutput", bridgeSwap.SolverOutput)

						// Try to extract the amount from the transaction receipt
						if bridgeSwap.Result != "" {
							// Get the actual output amount from the transaction
							swapOutput, err := bridge.GetSwapOutput(
								RPC_URL,
								bridgeSwap.Result,
							)

							if err == nil && swapOutput != "" {
								logger.Sugar().Infow("Successfully extracted amount from swap output",
									"amount", swapOutput)
								burnAmount = swapOutput
								// Update the swap operation with the correct output amount
								db.UpdateOperationSolverOutput(bridgeSwap.ID, swapOutput)
							} else {
								logger.Sugar().Warnw("Failed to extract amount from swap output",
									"error", err)
							}
						}

						// If still empty after fallback, fail the operation
						if burnAmount == "" {
							logger.Sugar().Errorw("Cannot proceed with burn operation: no valid amount available")
							db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
							db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
							break OperationLoop
						}
					}

					wallet, err := db.GetWallet(intent.Identity, blockchains.Ethereum)
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
					logger.Sugar().Infow("BURN_SYNTHETIC validation passed: followed by WITHDRAW",
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

					wallet, err := db.GetWallet(intent.Identity, intent.BlockchainID)
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

					fmt.Println("DEBUG: BURN_SYNTHETIC dataToSign format check:")
					fmt.Printf("- Raw value: %q\n", dataToSign)
					fmt.Printf("- Has '0x' prefix: %v\n", strings.HasPrefix(dataToSign, "0x"))
					fmt.Printf("- Length: %d\n", len(dataToSign))
					fmt.Printf("- First 20 bytes: %v\n", dataToSign[:min(len(dataToSign), 20)])

					// Remove the code that adds 0x prefix
					// if !strings.HasPrefix(dataToSign, "0x") {
					// 	dataToSign = "0x" + dataToSign
					// 	logger.Sugar().Infow("BURN_SYNTHETIC added 0x prefix to dataToSign",
					// 		"operationId", operation.ID,
					// 		"originalLength", len(dataToSign)-2,
					// 		"newLength", len(dataToSign))
					// }

					logger.Sugar().Infow("BURN_SYNTHETIC data generated successfully",
						"operationId", operation.ID,
						"dataToSignLength", len(dataToSign),
						"dataToSignPrefix", truncateString(dataToSign, 20),
						"dataToSignHasPrefix", strings.HasPrefix(dataToSign, "0x"),
						"dataToSignFormat", fmt.Sprintf("%T", dataToSign))

					// Update operation with data to sign and wait for signature
					db.UpdateOperationSolverDataToSign(operation.ID, dataToSign)
					intent.Operations[i].SolverDataToSign = dataToSign

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

					// Save the transaction hash and amount as the solver output
					// This amount is needed by the subsequent WITHDRAW operation
					db.UpdateOperationResult(operation.ID, libs.OperationStatusWaiting, result)
					db.UpdateOperationSolverOutput(operation.ID, burnSyntheticMetadata.Amount)

					logger.Sugar().Infow("BURN_SYNTHETIC saved amount as solver output",
						"operationId", operation.ID,
						"amount", burnSyntheticMetadata.Amount)

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

					// Set default value for unlock if not specified (empty JSON would set it to false)
					// Check if the original metadata JSON contains the "unlock" field
					var metadataMap map[string]interface{}
					err := json.Unmarshal([]byte(operation.SolverMetadata), &metadataMap)
					if err != nil {
						logger.Sugar().Errorw("Failed to unmarshal metadata for withdraw operation", "error", err)
					}
					if _, hasUnlock := metadataMap["unlock"]; !hasUnlock {
						// If "unlock" wasn't specified in the JSON, default to true
						withdrawMetadata.Unlock = true
						logger.Sugar().Infow("Setting default unlock=true for withdrawal",
							"operationID", operation.ID,
							"tokenAddress", withdrawMetadata.Token)
					}

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
					} else if operation.Type == OPERATION_TYPE_SIGN_MESSAGE {
						signature, err := getSignature(intent, i)
						if err != nil {
							fmt.Println("Message signing error:", err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
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

					bridgeWalletPublicKey := getBridgeWalletPublicKey(&operation, bridgeWallet)

					// Convert burn.SolverOutput from numeric string to proper JSON format
					// This fixes the "cannot unmarshal number into Go value of type map[string]interface{}" error
					burnSolverOutputJSONString := fmt.Sprintf(`{"amount": "%s"}`, burn.SolverOutput)
					logger.Sugar().Infow("Created JSON format for burn output",
						"originalOutput", burn.SolverOutput,
						"jsonFormatted", burnSolverOutputJSONString)

					tx, dataToSign, err := opBlockchain.BuildWithdrawTx(bridgeWalletPublicKey, burnSolverOutputJSONString, publicKey, &tokenToWithdraw)
					if err != nil {
						logger.Sugar().Errorw("error building withdraw transaction", "error", err)
						break
					}
					db.UpdateOperationSolverDataToSign(operation.ID, dataToSign)
					intent.Operations[i].SolverDataToSign = dataToSign

					withdrawSignature, err := getSignature(intent, i)
					if err != nil {
						logger.Sugar().Errorw("error getting signature", "error", err)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					result, err := opBlockchain.BroadcastTransaction(
						tx,
						withdrawSignature,
						&bridgeWalletPublicKey,
					)

					if err != nil {
						logger.Sugar().Errorw("error withdrawing tokens", "error", err)
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
				case libs.OperationTypeTransaction, libs.OperationTypeSendToBridge, libs.OperationTypeBridgeDeposit:
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
				case libs.OperationTypeSwap:
					confirmed, err := opBlockchain.IsTransactionBroadcastedAndConfirmed(operation.Result)
					if err != nil {
						logger.Sugar().Errorw("error checking swap transaction", "error", err)
						break
					}

					if !confirmed {
						break
					}

					// Extract the actual output amount from the swap transaction
					swapOutput, err := bridge.GetSwapOutput(
						RPC_URL,
						operation.Result,
					)

					if err != nil {
						logger.Sugar().Errorw("error getting swap output", "error", err,
							"txHash", operation.Result)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break
					}

					logger.Sugar().Infow("Successfully extracted swap output for burn operation",
						"txHash", operation.Result,
						"swapOutput", swapOutput)

					// Update the operation status and solver output with the actual amount
					db.UpdateOperationStatus(operation.ID, libs.OperationStatusCompleted)
					db.UpdateOperationSolverOutput(operation.ID, swapOutput)

					if i+1 == len(intent.Operations) {
						// update the intent status to completed
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusCompleted)
					}

					break OperationLoop
				case libs.OperationTypeBurn, libs.OperationTypeBurnSynthetic:
					confirmed, err := opBlockchain.IsTransactionBroadcastedAndConfirmed(operation.Result)
					if err != nil {
						logger.Sugar().Errorw("error checking burn transaction", "error", err)
						break
					}

					if !confirmed {
						break
					}

					// Extract the actual output amount from the burn transaction
					burnOutput, err := bridge.GetBurnOutput(
						RPC_URL,
						operation.Result,
					)

					if err != nil {
						logger.Sugar().Errorw("error getting burn output", "error", err,
							"txHash", operation.Result)
						db.UpdateOperationStatus(operation.ID, libs.OperationStatusFailed)
						db.UpdateIntentStatus(intent.ID, libs.IntentStatusFailed)
						break ProcessLoop
					}

					logger.Sugar().Infow("Successfully extracted burn output for withdraw operation",
						"txHash", operation.Result,
						"burnOutput", burnOutput)

					// Update the operation status and solver output with the actual amount
					db.UpdateOperationStatus(operation.ID, libs.OperationStatusCompleted)
					db.UpdateOperationSolverOutput(operation.ID, burnOutput)

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

					// Set default value for unlock if not specified (empty JSON would set it to false)
					// Check if the original metadata JSON contains the "unlock" field
					var metadataMap map[string]interface{}
					json.Unmarshal([]byte(operation.SolverMetadata), &metadataMap)
					if _, hasUnlock := metadataMap["unlock"]; !hasUnlock {
						// If "unlock" wasn't specified in the JSON, default to true
						withdrawMetadata.Unlock = true
						logger.Sugar().Infow("Setting default unlock=true for withdrawal",
							"operationID", operation.ID,
							"tokenAddress", withdrawMetadata.Token)
					}

					lockSchema, err := db.GetLock(intent.Identity, intent.BlockchainID)
					if err != nil {
						logger.Sugar().Errorw("error getting lock", "error", err)
						break
					}

					if withdrawMetadata.Unlock {
						// TODO: proper unlock handling
						// Check if i-4 is a valid index before accessing it
						if i >= 4 {
							depositOperation := intent.Operations[i-4]
							// check for confirmations
							confirmed, err = opBlockchain.IsTransactionBroadcastedAndConfirmed(depositOperation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking transaction", "error", err)
								break
							}
						} else {
							// Log that we couldn't find an expected deposit operation
							logger.Sugar().Warnw("no deposit operation found 4 positions before withdraw",
								"operationId", operation.ID,
								"intentId", intent.ID,
								"currentIndex", i,
								"unlockRequested", withdrawMetadata.Unlock)
							// Proceed without checking deposit confirmation
							confirmed = true
						}
					}
					if confirmed {
						err := db.UnlockIdentity(lockSchema.Id)
						if err != nil {
							logger.Sugar().Errorw("error unlocking identity", "error", err)
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
	signature, err := getSignatureEx(intent, operationIndex)
	if err != nil {
		return "", err
	}
	return signature, nil
}

func getSignatureEx(intent *libs.Intent, operationIndex int) (string, error) {
	// get wallet
	transactionType := intent.Operations[operationIndex].Type
	var wallet *db.WalletSchema
	var err error
	if transactionType == libs.OperationTypeWithdraw {
		wallet, err = db.GetWallet(BridgeContractAddress, blockchains.Ethereum)
		if err != nil {
			return "", fmt.Errorf("error getting wallet: %v", err)
		}
	} else {
		wallet, err = db.GetWallet(intent.Identity, intent.BlockchainID)
		if err != nil {
			return "", fmt.Errorf("error getting wallet: %v", err)
		}
	}

	// get the signer
	signers := wallet.Signers
	signer, err := GetSigner(signers[0])

	if err != nil {
		return "", fmt.Errorf("error getting signer: %v", err)
	}

	operation := intent.Operations[operationIndex]
	fmt.Printf("DEBUG: Requesting signature for operation type: %s\n", operation.Type)

	// For BURN_SYNTHETIC operations, add extra debug details
	if operation.Type == libs.OperationTypeBurnSynthetic {
		fmt.Println("DEBUG: BURN_SYNTHETIC operation signature request details:")
		fmt.Printf("- SolverDataToSign format: %T\n", operation.SolverDataToSign)
		fmt.Printf("- SolverDataToSign has '0x' prefix: %v\n", strings.HasPrefix(operation.SolverDataToSign, "0x"))
		fmt.Printf("- SolverDataToSign length: %d\n", len(operation.SolverDataToSign))
		fmt.Printf("- SolverDataToSign first 20 chars: %v\n", operation.SolverDataToSign[:min(len(operation.SolverDataToSign), 20)])

		// Create a modified intent for debugging to see if validator unmarshals properly
		debugIntent := *intent
		debugOps := make([]libs.Operation, len(intent.Operations))
		copy(debugOps, intent.Operations)
		debugIntent.Operations = debugOps

		// Attempt manual JSON marshal of intent to check format
		debugBytes, debugErr := json.Marshal(debugIntent)
		if debugErr == nil {
			fmt.Printf("DEBUG: Intent JSON starts with: %s\n", string(debugBytes)[:min(len(string(debugBytes)), 100)])
		}
	}

	intentBytes, err := json.Marshal(intent)
	if err != nil {
		return "", fmt.Errorf("error marshalling intent: %v", err)
	}

	operationIndexStr := strconv.FormatUint(uint64(operationIndex), 10)

	// Log the request details for debugging
	logger.Sugar().Infow("Requesting signature from validator",
		"url", signer.URL+"/signature?operationIndex="+operationIndexStr,
		"intentID", intent.ID,
		"operationIndex", operationIndex,
		"operationType", operation.Type,
		"requestBodyLength", len(intentBytes),
		"solverDataToSignLength", len(operation.SolverDataToSign))

	client, err := validatorClientManager.GetClient(signer.URL)
	if err != nil {
		return "", fmt.Errorf("error getting validator client: %v", err)
	}

	protoIntent, err := libs.IntentToProto(intent)
	if err != nil {
		return "", fmt.Errorf("error converting intent to proto: %v", err)
	}
	resp, err := client.SignIntentOperation(context.Background(), &pb.SignIntentOperationRequest{
		Intent:         protoIntent,
		OperationIndex: uint32(operationIndex),
	})
	if err != nil {
		return "", fmt.Errorf("error getting signature: %v", err)
	}

	// Validate signature response
	if len(resp.Signature) == 0 {
		return "", fmt.Errorf("empty signature in response")
	}

	return resp.Signature, nil
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

func verifyDestinationAddress(bridgeWallet *db.WalletSchema, operation *libs.Operation) (bool, error) {
	var destAddress string
	var expectedAddress string
	var tokenAddress string
	opBlockchain, err := blockchains.GetBlockchain(operation.BlockchainID, operation.NetworkType)
	if err != nil {
		return false, fmt.Errorf("error getting blockchain: %v", err)
	}
	chainID := opBlockchain.ChainID()
	if chainID == nil {
		return false, fmt.Errorf("chainID is nil ")
	}
	destAddress, tokenAddress, _ = opBlockchain.ExtractDestinationAddress(*operation.SerializedTxn)
	switch operation.BlockchainID {
	// Extract destination address from serialized transaction
	case blockchains.Bitcoin:
		expectedAddress = bridgeWallet.BitcoinMainnetPublicKey
	case blockchains.Dogecoin:
		expectedAddress = bridgeWallet.DogecoinMainnetPublicKey
	// Extract destination address from serialized transaction based on chain type
	case blockchains.Solana:
		// Get the actual account address from the message accounts
		expectedAddress = bridgeWallet.SolanaPublicKey
	case blockchains.Aptos:
		expectedAddress = bridgeWallet.AptosEDDSAPublicKey
	case blockchains.Stellar:
		expectedAddress = bridgeWallet.StellarPublicKey
	case blockchains.Algorand:
		expectedAddress = bridgeWallet.AlgorandEDDSAPublicKey
	case blockchains.Ripple:
		expectedAddress = bridgeWallet.RippleEDDSAPublicKey
	case blockchains.Cardano:
		expectedAddress = bridgeWallet.CardanoPublicKey
	case blockchains.Sui:
		expectedAddress = bridgeWallet.SuiPublicKey
	default:
		expectedAddress = bridgeWallet.EthereumPublicKey
	}

	if tokenAddress != "" {
		exists, peggedToken, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, *chainID, tokenAddress)
		if err != nil {
			return false, fmt.Errorf("error checking token existence in bridge", err)
		}
		if !exists {
			return false, fmt.Errorf("ERC20 token not registered in bridge", tokenAddress, chainID)
		}
		logger.Sugar().Infow("ERC20 token exists in bridge", "token", tokenAddress, "peggedToken", peggedToken)
	}

	// Verify the extracted destination matches the bridge wallet
	if destAddress == "" {
		return false, fmt.Errorf("Failed to extract destination address from %s transaction", chainID)
	}

	if !strings.EqualFold(destAddress, expectedAddress) {
		return false, fmt.Errorf("Invalid bridge destination address for %s transaction", chainID)
	}

	return true, nil
}

func getBridgeWalletPublicKey(operation *libs.Operation, bridgeWallet *db.WalletSchema) string {
	var bridgeWalletPublicKey string
	switch operation.BlockchainID {
	case blockchains.Bitcoin:
		bridgeWalletPublicKey = bridgeWallet.BitcoinMainnetPublicKey
	case blockchains.Dogecoin:
		bridgeWalletPublicKey = bridgeWallet.DogecoinMainnetPublicKey
	case blockchains.Solana:
		// Get the actual account address from the message accounts
		bridgeWalletPublicKey = bridgeWallet.SolanaPublicKey
	case blockchains.Aptos:
		bridgeWalletPublicKey = bridgeWallet.AptosEDDSAPublicKey
	case blockchains.Stellar:
		bridgeWalletPublicKey = bridgeWallet.StellarPublicKey
	case blockchains.Algorand:
		bridgeWalletPublicKey = bridgeWallet.AlgorandEDDSAPublicKey
	case blockchains.Ripple:
		bridgeWalletPublicKey = bridgeWallet.RippleEDDSAPublicKey
	case blockchains.Cardano:
		bridgeWalletPublicKey = bridgeWallet.CardanoPublicKey
	case blockchains.Sui:
		bridgeWalletPublicKey = bridgeWallet.SuiPublicKey
	default:
		bridgeWalletPublicKey = bridgeWallet.EthereumPublicKey
	}

	return bridgeWalletPublicKey
}
