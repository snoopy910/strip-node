package sequencer

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/StripChain/strip-node/bridge"
	"github.com/StripChain/strip-node/solver"
	"github.com/StripChain/strip-node/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
)

type MintOutput struct {
	Token  string `json:"token"`
	Amount string `json:"amount"`
}

type SwapMetadata struct {
	Token string `json:"token"`
}

type BurnMetadata struct {
	Token string `json:"token"`
}

type WithdrawMetadata struct {
	Token string `json:"token"`
}

func ProcessIntent(intentId int64) {
	for {
		intent, err := GetIntent(intentId)
		if err != nil {
			log.Println(err)
			return
		}

		intentBytes, err := json.Marshal(intent)
		if err != nil {
			log.Println(err)
			return
		}

		if intent.Status != INTENT_STATUS_PROCESSING {
			log.Println("intent processed")
			return
		}

		if intent.Expiry < uint64(time.Now().Unix()) {
			UpdateIntentStatus(intent.ID, INTENT_STATUS_EXPIRED)
			return
		}

		// now process the operations of the intent
		for i, operation := range intent.Operations {
			if operation.Status == OPERATION_STATUS_COMPLETED || operation.Status == OPERATION_STATUS_FAILED {
				continue
			}

			if operation.Status == OPERATION_STATUS_PENDING {
				// sign and send the txn. Change status to waiting

				if operation.Type == OPERATION_TYPE_TRANSACTION {
					if operation.KeyCurve == "ecdsa" {
						signature, err := getSignature(intent, i)
						if err != nil {
							fmt.Println(err)
							break
						}

						txnHash, err := sendEVMTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)

						// @TODO: For our infra errors, don't mark the intent and operation as failed
						if err != nil {
							fmt.Println(err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, txnHash)
					} else if operation.KeyCurve == "eddsa" {
						chain, err := GetChain(operation.ChainId)
						if err != nil {
							fmt.Println(err)
							break
						}

						signature, err := getSignature(intent, i)

						if err != nil {
							fmt.Println(err)
							break
						}

						if chain.ChainType == "solana" {
							txnHash, err := sendSolanaTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)

							if err != nil {
								fmt.Println(err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, txnHash)
						}
					}
				} else if operation.Type == OPERATION_TYPE_SOLVER {
					// get data to sign from solver
					dataToSign, err := solver.Construct(operation.Solver, &intentBytes, i)

					if err != nil {
						fmt.Println(err)
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					UpdateOperationSolverDataToSign(operation.ID, dataToSign)

					// then get the data signed
					signature, err := getSignature(intent, i)
					if err != nil {
						fmt.Println(err)
						break
					}

					// then send the signature to solver
					result, err := solver.Solve(
						operation.Solver, &intentBytes,
						i,
						signature,
					)

					if err != nil {
						fmt.Println(err)
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
				} else if operation.Type == OPERATION_TYPE_BRIDGE_DEPOSIT {
					depositOperation := intent.Operations[i-1]

					if i == 0 || !(depositOperation.Type == OPERATION_TYPE_TRANSACTION || depositOperation.Type == OPERATION_TYPE_SOLVER) {
						fmt.Println("Invalid operation type for bridge deposit")
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					if depositOperation.KeyCurve == "ecdsa" {
						// find token transfer events and check if first transfer is a valid token
						transfers, err := GetEthereumTransfers(depositOperation.ChainId, depositOperation.Result, intent.Identity)
						if err != nil {
							fmt.Println(err)
							break
						}

						if len(transfers) == 0 {
							fmt.Println("No transfers found", depositOperation.Result, intent.Identity)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						// check if the token exists
						transfer := transfers[0]
						srcAddress := transfer.TokenAddress
						amount := transfer.ScaledAmount

						exists, destAddress, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, depositOperation.ChainId, srcAddress)

						if err != nil {
							fmt.Println(err)
							break
						}

						if !exists {
							fmt.Println("Token does not exist", srcAddress, depositOperation.ChainId)

							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						wallet, err := GetWallet(intent.Identity, "ecdsa")
						if err != nil {
							fmt.Println(err)
							break
						}

						dataToSign, err := bridge.BridgeDepositDataToSign(RPC_URL, BridgeContractAddress, amount, wallet.ECDSAPublicKey, destAddress)
						if err != nil {
							fmt.Println(err)
							break
						}

						UpdateOperationSolverDataToSign(operation.ID, dataToSign)
						intent.Operations[i].SolverDataToSign = dataToSign

						signature, err := getSignature(intent, i)
						if err != nil {
							fmt.Println(err)
							break
						}

						fmt.Println("Minting bridge", amount, wallet.ECDSAPublicKey, destAddress, signature)

						result, err := mintBridge(
							amount, wallet.ECDSAPublicKey, destAddress, signature)

						if err != nil {
							fmt.Println(err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						mintOutput := MintOutput{
							Token:  destAddress,
							Amount: amount,
						}

						mintOutputBytes, err := json.Marshal(mintOutput)

						if err != nil {
							fmt.Println(err)
							break
						}

						UpdateOperationSolverOutput(operation.ID, string(mintOutputBytes))
						UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)

					} else if depositOperation.KeyCurve == "eddsa" {
						transfers, err := GetSolanaTransfers(depositOperation.ChainId, depositOperation.Result, HeliusApiKey)
						if err != nil {
							fmt.Println(err)
							break
						}

						if len(transfers) == 0 {
							fmt.Println("No transfers found", depositOperation.Result, intent.Identity)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						// check if the token exists
						transfer := transfers[0]
						srcAddress := transfer.TokenAddress
						amount := transfer.ScaledAmount

						exists, destAddress, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, depositOperation.ChainId, srcAddress)

						if err != nil {
							fmt.Println(err)
							break
						}

						if !exists {
							fmt.Println("Token does not exist", srcAddress, depositOperation.ChainId)

							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						wallet, err := GetWallet(intent.Identity, "ecdsa")
						if err != nil {
							fmt.Println(err)
							break
						}

						dataToSign, err := bridge.BridgeDepositDataToSign(RPC_URL, BridgeContractAddress, amount, wallet.ECDSAPublicKey, destAddress)
						if err != nil {
							fmt.Println(err)
							break
						}

						UpdateOperationSolverDataToSign(operation.ID, dataToSign)
						intent.Operations[i].SolverDataToSign = dataToSign

						signature, err := getSignature(intent, i)
						if err != nil {
							fmt.Println(err)
							break
						}

						result, err := mintBridge(
							amount, wallet.ECDSAPublicKey, destAddress, signature)

						if err != nil {
							fmt.Println(err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
					}
				} else if operation.Type == OPERATION_TYPE_SWAP {
					bridgeDeposit := intent.Operations[i-1]

					if i == 0 || !(bridgeDeposit.Type == OPERATION_TYPE_BRIDGE_DEPOSIT) {
						fmt.Println("Invalid operation type for swap")
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					var bridgeDepositData MintOutput
					var swapMetadata SwapMetadata
					json.Unmarshal([]byte(bridgeDeposit.SolverOutput), &bridgeDepositData)
					json.Unmarshal([]byte(operation.SolverMetadata), &swapMetadata)

					tokenIn := bridgeDepositData.Token
					tokenOut := swapMetadata.Token
					amountIn := bridgeDepositData.Amount
					deadline := time.Now().Add(time.Hour).Unix()

					wallet, err := GetWallet(intent.Identity, "ecdsa")
					if err != nil {
						fmt.Println(err)
						break
					}

					dataToSign, err := bridge.BridgeSwapDataToSign(
						RPC_URL,
						BridgeContractAddress,
						wallet.ECDSAPublicKey,
						tokenIn,
						tokenOut,
						amountIn,
						deadline,
					)

					if err != nil {
						fmt.Println(err)
						break
					}

					UpdateOperationSolverDataToSign(operation.ID, dataToSign)
					intent.Operations[i].SolverDataToSign = dataToSign

					signature, err := getSignature(intent, i)
					if err != nil {
						fmt.Println(err)
						break
					}

					fmt.Println("Swapping bridge", wallet.ECDSAPublicKey, tokenIn, tokenOut, amountIn, deadline, signature)

					result, err := swapBridge(
						wallet.ECDSAPublicKey,
						tokenIn,
						tokenOut,
						amountIn,
						deadline,
						signature,
					)

					if err != nil {
						fmt.Println(err)
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)

					break
				} else if operation.Type == OPERATION_TYPE_BURN {
					bridgeSwap := intent.Operations[i-1]

					if i == 0 || !(bridgeSwap.Type == OPERATION_TYPE_SWAP) {
						fmt.Println("Invalid operation type for swap")
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					fmt.Println("Burning tokens", bridgeSwap)

					burnAmount := bridgeSwap.SolverOutput
					burnMetadata := BurnMetadata{}

					json.Unmarshal([]byte(operation.SolverMetadata), &burnMetadata)

					wallet, err := GetWallet(intent.Identity, "ecdsa")
					if err != nil {
						fmt.Println(err)
						break
					}

					dataToSign, err := bridge.BridgeBurnDataToSign(
						RPC_URL,
						BridgeContractAddress,
						wallet.ECDSAPublicKey,
						burnAmount,
						burnMetadata.Token,
					)

					if err != nil {
						fmt.Println(err)
						break
					}

					UpdateOperationSolverDataToSign(operation.ID, dataToSign)
					intent.Operations[i].SolverDataToSign = dataToSign

					signature, err := getSignature(intent, i)
					if err != nil {
						fmt.Println(err)
						break
					}

					fmt.Println("Burn tokens", wallet.ECDSAPublicKey, burnAmount, burnMetadata.Token, signature)

					result, err := burnTokens(
						wallet.ECDSAPublicKey,
						burnAmount,
						burnMetadata.Token,
						signature,
					)

					if err != nil {
						fmt.Println(err)
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
					break
				} else if operation.Type == OPERATION_TYPE_WITHDRAW {
					burn := intent.Operations[i-1]

					if i == 0 || !(burn.Type == OPERATION_TYPE_BURN) {
						fmt.Println("Invalid operation type for withdraw after burn")
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					var withdrawMetadata WithdrawMetadata
					var burnMetadata BurnMetadata
					json.Unmarshal([]byte(operation.SolverMetadata), &withdrawMetadata)
					json.Unmarshal([]byte(burn.SolverMetadata), &burnMetadata)

					tokenToWithdraw := withdrawMetadata.Token

					// verify these fields
					exists, destAddress, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, operation.ChainId, tokenToWithdraw)

					if err != nil {
						fmt.Println(err)
						break
					}

					if !exists {
						fmt.Println("Token does not exist", tokenToWithdraw, operation.ChainId)

						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					if destAddress != burnMetadata.Token {
						fmt.Println("Token mismatch", destAddress, burnMetadata.Token)

						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					withdrawalChain, err := GetChain(operation.ChainId)

					if err != nil {
						fmt.Println(err)
						break
					}

					bridgeWallet, err := GetWallet(BridgeContractAddress, "ecdsa")
					if err != nil {
						fmt.Println(err)
						break
					}

					user, err := GetWallet(intent.Identity, intent.IdentityCurve)
					if err != nil {
						fmt.Println(err)
						break
					}

					if withdrawalChain.KeyCurve == "ecdsa" {
						if tokenToWithdraw == util.ZERO_ADDRESS {
							// handle native token
							dataToSign, tx, err := withdrawEVMNativeGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWallet.ECDSAPublicKey,
								burn.SolverOutput,
								user.ECDSAPublicKey,
								operation.ChainId,
							)

							if err != nil {
								fmt.Println(err)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
							intent.Operations[i].SolverDataToSign = dataToSign

							signature, err := getSignature(intent, i)
							if err != nil {
								fmt.Println(err)
								break
							}

							result, err := withdrawEVMTxn(
								withdrawalChain.ChainUrl,
								signature,
								tx,
								operation.ChainId,
							)

							if err != nil {
								fmt.Println(err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
						} else {
							dataToSign, tx, err := withdrawERC20GetSignature(
								withdrawalChain.ChainUrl,
								bridgeWallet.ECDSAPublicKey,
								burn.SolverOutput,
								user.ECDSAPublicKey,
								operation.ChainId,
								tokenToWithdraw,
							)

							if err != nil {
								fmt.Println(err)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
							intent.Operations[i].SolverDataToSign = dataToSign

							signature, err := getSignature(intent, i)
							if err != nil {
								fmt.Println(err)
								break
							}

							result, err := withdrawEVMTxn(
								withdrawalChain.ChainUrl,
								signature,
								tx,
								operation.ChainId,
							)

							if err != nil {
								fmt.Println(err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
						}

						break
					} else if withdrawalChain.KeyCurve == "eddsa" {
						if tokenToWithdraw == util.ZERO_ADDRESS {
							// handle native token
							transaction, dataToSign, err := withdrawSolanaNativeGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWallet.EDDSAPublicKey,
								burn.SolverOutput,
								user.ECDSAPublicKey,
							)

							if err != nil {
								fmt.Println(err)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
							intent.Operations[i].SolverDataToSign = dataToSign

							signature, err := getSignature(intent, i)
							if err != nil {
								fmt.Println(err)
								break
							}

							result, err := withdrawSolanaTxn(
								withdrawalChain.ChainUrl,
								transaction,
								signature,
							)

							if err != nil {
								fmt.Println(err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
						} else {
							// implement SPL
							transaction, dataToSign, err := withdrawSolanaSPLGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWallet.EDDSAPublicKey,
								burn.SolverOutput,
								user.ECDSAPublicKey,
								tokenToWithdraw,
							)

							if err != nil {
								fmt.Println(err)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
							intent.Operations[i].SolverDataToSign = dataToSign

							signature, err := getSignature(intent, i)
							if err != nil {
								fmt.Println(err)
								break
							}

							result, err := withdrawSolanaTxn(
								withdrawalChain.ChainUrl,
								transaction,
								signature,
							)

							if err != nil {
								fmt.Println(err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
						}
						break
					}
				}

				break
			}

			if operation.Status == OPERATION_STATUS_WAITING {
				// check for confirmations and update the status to completed
				if operation.Type == OPERATION_TYPE_TRANSACTION {
					confirmed := false
					if operation.KeyCurve == "ecdsa" {
						confirmed, err = checkEVMTransactionConfirmed(operation.ChainId, operation.Result)
						if err != nil {
							fmt.Println(err)
							break
						}
					} else if operation.KeyCurve == "eddsa" {
						chain, err := GetChain(operation.ChainId)
						if err != nil {
							fmt.Println(err)
							break
						}

						if chain.ChainType == "solana" {
							confirmed, err = checkSolanaTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						}
					}

					if !confirmed {
						break
					}

					UpdateOperationStatus(operation.ID, OPERATION_STATUS_COMPLETED)

					if i+1 == len(intent.Operations) {
						// update the intent status to completed
						UpdateIntentStatus(intent.ID, INTENT_STATUS_COMPLETED)
					}

					break
				} else if operation.Type == OPERATION_TYPE_BRIDGE_DEPOSIT {
					confirmed := false
					if operation.KeyCurve == "ecdsa" {
						confirmed, err = checkEVMTransactionConfirmed(operation.ChainId, operation.Result)
						if err != nil {
							fmt.Println(err)
							break
						}
					} else if operation.KeyCurve == "eddsa" {
						chain, err := GetChain(operation.ChainId)
						if err != nil {
							fmt.Println(err)
							break
						}

						if chain.ChainType == "solana" {
							confirmed, err = checkSolanaTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						}
					}

					if !confirmed {
						break
					}

					UpdateOperationStatus(operation.ID, OPERATION_STATUS_COMPLETED)

					if i+1 == len(intent.Operations) {
						// update the intent status to completed
						UpdateIntentStatus(intent.ID, INTENT_STATUS_COMPLETED)
					}

					break
				} else if operation.Type == OPERATION_TYPE_SOLVER {
					status, err := solver.CheckStatus(
						operation.Solver, &intentBytes, i,
					)

					if err != nil {
						fmt.Println(err)
						break
					}

					if status == solver.SOLVER_OPERATION_STATUS_SUCCESS {
						output, err := solver.GetOutput(operation.Solver, &intentBytes, i)

						if err != nil {
							fmt.Println(err)
							break
						}

						UpdateOperationStatus(operation.ID, OPERATION_STATUS_COMPLETED)
						UpdateOperationSolverOutput(operation.ID, output)

						if i+1 == len(intent.Operations) {
							// update the intent status to completed
							UpdateIntentStatus(intent.ID, INTENT_STATUS_COMPLETED)
						}
					}

					if status == solver.SOLVER_OPERATION_STATUS_FAILURE {
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
					}

					break
				} else if operation.Type == OPERATION_TYPE_SWAP {
					confirmed, err := checkEVMTransactionConfirmed(operation.ChainId, operation.Result)
					if err != nil {
						fmt.Println(err)
						break
					}

					if !confirmed {
						break
					}

					swapOutput, err := bridge.GetSwapOutput(
						RPC_URL,
						operation.Result,
					)

					if err != nil {
						fmt.Println(err)
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					UpdateOperationStatus(operation.ID, OPERATION_STATUS_COMPLETED)
					UpdateOperationSolverOutput(operation.ID, swapOutput)

					if i+1 == len(intent.Operations) {
						// update the intent status to completed
						UpdateIntentStatus(intent.ID, INTENT_STATUS_COMPLETED)
					}

					break
				} else if operation.Type == OPERATION_TYPE_BURN {
					confirmed, err := checkEVMTransactionConfirmed(operation.ChainId, operation.Result)
					if err != nil {
						fmt.Println(err)
						break
					}

					if !confirmed {
						break
					}

					swapOutput, err := bridge.GetBurnOutput(
						RPC_URL,
						operation.Result,
					)

					if err != nil {
						fmt.Println(err)
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					UpdateOperationStatus(operation.ID, OPERATION_STATUS_COMPLETED)
					UpdateOperationSolverOutput(operation.ID, swapOutput)

					if i+1 == len(intent.Operations) {
						// update the intent status to completed
						UpdateIntentStatus(intent.ID, INTENT_STATUS_COMPLETED)
					}

					break
				} else if operation.Type == OPERATION_TYPE_WITHDRAW {
					confirmed := false
					if operation.KeyCurve == "ecdsa" {
						confirmed, err = checkEVMTransactionConfirmed(operation.ChainId, operation.Result)
						if err != nil {
							fmt.Println(err)
							break
						}
					} else if operation.KeyCurve == "eddsa" {
						chain, err := GetChain(operation.ChainId)
						if err != nil {
							fmt.Println(err)
							break
						}

						if chain.ChainType == "solana" {
							confirmed, err = checkSolanaTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						}
					}

					if !confirmed {
						break
					}

					UpdateOperationStatus(operation.ID, OPERATION_STATUS_COMPLETED)

					if i+1 == len(intent.Operations) {
						// update the intent status to completed
						UpdateIntentStatus(intent.ID, INTENT_STATUS_COMPLETED)
					}

					break
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}

type SignatureResponse struct {
	Signature string `json:"signature"`
	Address   string `json:"address"`
}

func getSignature(intent *Intent, operationIndex int) (string, error) {
	// get wallet
	wallet, err := GetWallet(intent.Identity, intent.IdentityCurve)
	if err != nil {
		return "", err
	}

	// get the signer
	signers := strings.Split(wallet.Signers, ",")
	signer, err := GetSigner(signers[0])

	if err != nil {
		return "", err
	}

	intentBytes, err := json.Marshal(intent)
	if err != nil {
		return "", err
	}

	operationIndexStr := strconv.FormatUint(uint64(operationIndex), 10)

	req, err := http.NewRequest("POST", signer.URL+"/signature?operationIndex="+operationIndexStr, bytes.NewBuffer(intentBytes))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var signatureResponse SignatureResponse
	err = json.Unmarshal(body, &signatureResponse)
	if err != nil {
		return "", err
	}

	return signatureResponse.Signature, nil
}

func checkEVMTransactionConfirmed(chainId string, txnHash string) (bool, error) {
	chain, err := GetChain(chainId)
	if err != nil {
		return false, err
	}

	client, err := ethclient.Dial(chain.ChainUrl)
	if err != nil {
		log.Fatal(err)
	}

	_, isPending, err := client.TransactionByHash(context.Background(), common.HexToHash(txnHash))
	if err != nil {
		return false, err
	}

	return !isPending, nil
}

func checkSolanaTransactionConfirmed(chainId string, txnHash string) (bool, error) {
	chain, err := GetChain(chainId)
	if err != nil {
		return false, err
	}

	c := rpc.New(chain.ChainUrl)

	signature, err := solana.SignatureFromBase58(txnHash)
	if err != nil {
		return false, err
	}

	_, err = c.GetConfirmedTransaction(context.Background(), signature)

	if err != nil {
		return false, err
	}

	return true, nil

}

func sendEVMTransaction(serializedTxn string, chainId string, keyCurve string, dataToSign string, signatureHex string) (string, error) {
	chain, err := GetChain(chainId)
	if err != nil {
		return "", err
	}

	client, err := ethclient.Dial(chain.ChainUrl)
	if err != nil {
		return "", err
	}

	serializedTx, err := hex.DecodeString(serializedTxn)
	if err != nil {
		return "", err
	}

	var tx types.Transaction
	rlp.DecodeBytes(serializedTx, &tx)

	sigData, err := hex.DecodeString(signatureHex)

	if err != nil {
		return "", err
	}

	n, _ := new(big.Int).SetString(chainId, 10)
	_tx, err := tx.WithSignature(types.NewLondonSigner(n), []byte(sigData))

	if err != nil {
		return "", err
	}

	err = client.SendTransaction(context.Background(), _tx)
	if err != nil {
		return "", err
	}

	return _tx.Hash().Hex(), nil
}

func sendSolanaTransaction(serializedTxn string, chainId string, keyCurve string, dataToSign string, signatureBase58 string) (string, error) {
	chain, err := GetChain(chainId)
	if err != nil {
		return "", err
	}

	c := rpc.New(chain.ChainUrl)

	decodedTransactionData, err := base58.Decode(serializedTxn)
	if err != nil {
		fmt.Println("Error decoding transaction data:", err)
		return "", err
	}

	_tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(decodedTransactionData))
	if err != nil {
		return "", err
	}

	sig, _ := base58.Decode(signatureBase58)
	signature := solana.SignatureFromBytes(sig)

	_tx.Signatures = append(_tx.Signatures, signature)

	err = _tx.VerifySignatures()

	if err != nil {
		fmt.Println("error during verification")
		fmt.Println(err)
		return "", err
	}

	hash, err := c.SendTransaction(context.Background(), _tx)
	if err != nil {
		fmt.Println("error during sending transaction")
		return "", err
	}

	return hash.String(), nil
}
