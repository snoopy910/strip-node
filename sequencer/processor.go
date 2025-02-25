package sequencer

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/StripChain/strip-node/algorand"
	"github.com/StripChain/strip-node/aptos"
	"github.com/StripChain/strip-node/bridge"
	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/ripple"
	"github.com/StripChain/strip-node/solver"
	"github.com/StripChain/strip-node/stellar"
	"github.com/StripChain/strip-node/util"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	algorandTypes "github.com/algorand/go-algorand-sdk/types"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
	"github.com/stellar/go/xdr"
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
	Token  string `json:"token"`
	Unlock bool   `json:"unlock"`
}

type LockMetadata struct {
	Lock bool `json:"lock"`
}

func ProcessIntent(intentId int64) {
	for {
		intent, err := GetIntent(intentId)
		if err != nil {
			log.Printf("error getting intent: %+v\n", err)
			return
		}

		intentBytes, err := json.Marshal(intent)
		if err != nil {
			log.Printf("error marshalling intent: %+v\n", err)
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
					lockSchema, err := GetLock(intent.Identity, intent.IdentityCurve)
					if err != nil {
						if err.Error() == "pg: no rows in result set" {
							_, err := AddLock(intent.Identity, intent.IdentityCurve)

							if err != nil {
								fmt.Printf("error adding lock: %+v\n", err)
								break
							}

							lockSchema, err = GetLock(intent.Identity, intent.IdentityCurve)

							if err != nil {
								fmt.Printf("error getting lock after adding: %+v\n", err)
								break
							}
						} else {
							fmt.Printf("error getting lock: %+v\n", err)
							break
						}
					}

					if lockSchema.Locked {
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "secp256k1" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							fmt.Printf("error getting chain: %+v\n", err)
							break
						}

						signature, err := getSignature(intent, i)
						if err != nil {
							fmt.Printf("error getting signature: %+v\n", err)
							break
						}

						if chain.ChainType == "bitcoin" {
							txnHash, err := sendBitcoinTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)

							if err != nil {
								fmt.Println(err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							var lockMetadata LockMetadata
							json.Unmarshal([]byte(operation.SolverMetadata), &lockMetadata)

							if lockMetadata.Lock {
								err := LockIdentity(lockSchema.Id)
								if err != nil {
									fmt.Println(err)
									break
								}

								UpdateOperationResult(operation.ID, OPERATION_STATUS_COMPLETED, txnHash)
							} else {
								UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, txnHash)
							}
						} else {
							txnHash, err := sendEVMTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)

							// @TODO: For our infra errors, don't mark the intent and operation as failed
							if err != nil {
								fmt.Println(err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							var lockMetadata LockMetadata
							json.Unmarshal([]byte(operation.SolverMetadata), &lockMetadata)

							if lockMetadata.Lock {
								err := LockIdentity(lockSchema.Id)
								if err != nil {
									fmt.Println(err)
									break
								}

								UpdateOperationResult(operation.ID, OPERATION_STATUS_COMPLETED, txnHash)
							} else {
								UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, txnHash)
							}
						}

					} else if operation.KeyCurve == "eddsa" || operation.KeyCurve == "aptos_eddsa" || operation.KeyCurve == "stellar_eddsa" || operation.KeyCurve == "algorand_eddsa" || operation.KeyCurve == "ripple_eddsa" {
						chId := operation.ChainId
						if chId == "" {
							chId = operation.GenesisHash
						}
						chain, err := common.GetChain(chId)
						if err != nil {
							fmt.Printf("error getting chain: %+v\n", err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						signature, err := getSignature(intent, i)

						if err != nil {
							fmt.Printf("error getting signature: %+v\n", err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						var txnHash string

						if chain.ChainType == "solana" {
							txnHash, err = sendSolanaTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
							if err != nil {
								fmt.Println(err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
						}

						if chain.ChainType == "aptos" {
							// Convert public key
							wallet, err := GetWallet(intent.Identity, intent.IdentityCurve)
							if err != nil {
								fmt.Printf("error getting public key: %v", err)
								break
							}
							txnHash, err = aptos.SendAptosTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, wallet.AptosEDDSAPublicKey, signature)
							fmt.Println(txnHash)
							if err != nil {
								fmt.Println(err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
						}

						if chain.ChainType == "algorand" {
							txnHash, err = algorand.SendAlgorandTransaction(operation.SerializedTxn, operation.GenesisHash, signature)
							if err != nil {
								fmt.Println(err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

						}
						if chain.ChainType == "stellar" {
							// Send Stellar transaction
							txnHash, err = stellar.SendStellarTxn(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
							if err != nil {
								fmt.Printf("error sending Stellar transaction: %v", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
						}

						if chain.ChainType == "ripple" {
							// Convert public key
							wallet, err := GetWallet(intent.Identity, intent.IdentityCurve)
							if err != nil {
								fmt.Printf("error getting public key: %v", err)
								break
							}

							txnHash, err = ripple.SendRippleTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, wallet.RippleEDDSAPublicKey, signature)
							if err != nil {
								fmt.Printf("error sending Ripple transaction: %+v", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
						}

						var lockMetadata LockMetadata
						json.Unmarshal([]byte(operation.SolverMetadata), &lockMetadata)

						if lockMetadata.Lock {
							err := LockIdentity(lockSchema.Id)
							if err != nil {
								fmt.Println(err)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_COMPLETED, txnHash)
						} else {
							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, txnHash)
						}

					}
				} else if operation.Type == OPERATION_TYPE_SOLVER {
					lockSchema, err := GetLock(intent.Identity, intent.IdentityCurve)
					if err != nil {
						if err.Error() == "pg: no rows in result set" {
							_, err := AddLock(intent.Identity, intent.IdentityCurve)

							if err != nil {
								fmt.Println(err)
								break
							}

							lockSchema, err = GetLock(intent.Identity, intent.IdentityCurve)

							if err != nil {
								fmt.Println(err)
								break
							}
						} else {
							fmt.Println(err)
							break
						}
					}

					if lockSchema.Locked {
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

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

					var lockMetadata LockMetadata
					json.Unmarshal([]byte(operation.SolverMetadata), &lockMetadata)

					if lockMetadata.Lock {
						err := LockIdentity(lockSchema.Id)
						if err != nil {
							fmt.Println(err)
							break
						}

						UpdateOperationResult(operation.ID, OPERATION_STATUS_COMPLETED, result)
					} else {
						UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
					}
				} else if operation.Type == OPERATION_TYPE_SEND_TO_BRIDGE {
					// Get bridge wallet for the chain
					bridgeWallet, err := GetWallet(BridgeContractAddress, operation.KeyCurve)
					if err != nil {
						fmt.Println("Failed to get bridge wallet:", err)
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					// Process transaction based on key curve and chain type
					lockSchema, err := GetLock(intent.Identity, intent.IdentityCurve)
					if err != nil {
						if err.Error() == "pg: no rows in result set" {
							_, err := AddLock(intent.Identity, intent.IdentityCurve)
							if err != nil {
								fmt.Printf("error adding lock: %+v\n", err)
								break
							}
							lockSchema, err = GetLock(intent.Identity, intent.IdentityCurve)
							if err != nil {
								fmt.Printf("error getting lock after adding: %+v\n", err)
								break
							}
						} else {
							fmt.Printf("error getting lock: %+v\n", err)
							break
						}
					}

					if lockSchema.Locked {
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "secp256k1" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							fmt.Printf("error getting chain: %+v\n", err)
							break
						}

						// Extract destination address from serialized transaction
						var destAddress string
						if chain.ChainType == "bitcoin" {
							// For Bitcoin, decode the serialized transaction to get output address
							var tx wire.MsgTx
							txBytes, err := hex.DecodeString(operation.SerializedTxn)
							if err != nil {
								fmt.Printf("error decoding bitcoin transaction: %+v\n", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
							if err := tx.Deserialize(bytes.NewReader(txBytes)); err != nil {
								fmt.Printf("error deserializing bitcoin transaction: %+v\n", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
							// Get the first output's address (assuming it's the bridge address)
							if len(tx.TxOut) > 0 {
								_, addrs, _, err := txscript.ExtractPkScriptAddrs(tx.TxOut[0].PkScript, nil)
								if err != nil || len(addrs) == 0 {
									fmt.Printf("error extracting bitcoin address: %+v\n", err)
									UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
									UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
									break
								}
								destAddress = addrs[0].String()
							}
						} else {
							// For EVM chains, decode the transaction to get the 'to' address
							txBytes, err := hex.DecodeString(operation.SerializedTxn)
							if err != nil {
								fmt.Printf("error decoding EVM transaction: %+v\n", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
							tx := new(types.Transaction)
							if err := rlp.DecodeBytes(txBytes, tx); err != nil {
								fmt.Printf("error deserializing EVM transaction: %+v\n", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
							destAddress = tx.To().Hex()
						}

						// Verify destination address matches bridge wallet
						var expectedAddress string
						if chain.ChainType == "bitcoin" {
							expectedAddress = bridgeWallet.BitcoinMainnetPublicKey
						} else {
							expectedAddress = bridgeWallet.ECDSAPublicKey
						}

						if !strings.EqualFold(destAddress, expectedAddress) {
							fmt.Printf("Invalid bridge destination address. Expected %s, got %s\n", expectedAddress, destAddress)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						signature, err := getSignature(intent, i)
						if err != nil {
							fmt.Printf("error getting signature: %+v\n", err)
							break
						}

						var txnHash string
						switch chain.ChainType {
						case "bitcoin":
							txnHash, err = sendBitcoinTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
						default: // EVM chains
							txnHash, err = sendEVMTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
						}

						if err != nil {
							fmt.Println(err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						var lockMetadata LockMetadata
						json.Unmarshal([]byte(operation.SolverMetadata), &lockMetadata)

						if lockMetadata.Lock {
							err := LockIdentity(lockSchema.Id)
							if err != nil {
								fmt.Println(err)
								break
							}
							UpdateOperationResult(operation.ID, OPERATION_STATUS_COMPLETED, txnHash)
						} else {
							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, txnHash)
						}
					} else if operation.KeyCurve == "eddsa" || operation.KeyCurve == "aptos_eddsa" || operation.KeyCurve == "stellar_eddsa" || operation.KeyCurve == "algorand_eddsa" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							fmt.Printf("error getting chain: %+v\n", err)
							break
						}

						// Verify destination address matches bridge wallet based on chain type
						var validDestination bool
						var destAddress string

						// Extract destination address from serialized transaction based on chain type
						switch chain.ChainType {
						case "solana":
							// Decode base58 transaction and extract destination
							decodedTxn, err := base58.Decode(operation.SerializedTxn)
							if err != nil {
								fmt.Printf("error decoding Solana transaction: %v\n", err)
								break
							}
							tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(decodedTxn))
							if err != nil || len(tx.Message.Instructions) == 0 {
								fmt.Printf("error deserializing Solana transaction: %v\n", err)
								break
							}
							// Get the first instruction's destination account index
							destAccountIndex := tx.Message.Instructions[0].Accounts[1]
							// Get the actual account address from the message accounts
							destAddress = tx.Message.AccountKeys[destAccountIndex].String()
						case "aptos":
							// For Aptos, the destination is in the transaction payload
							var aptosPayload struct {
								Function string   `json:"function"`
								Args     []string `json:"arguments"`
							}
							if err := json.Unmarshal([]byte(operation.SerializedTxn), &aptosPayload); err != nil {
								fmt.Printf("error parsing Aptos transaction: %v\n", err)
								break
							}
							if len(aptosPayload.Args) > 0 {
								destAddress = aptosPayload.Args[0] // First arg is typically the recipient
							}
						case "stellar":
							// For Stellar, parse the XDR transaction envelope
							var txEnv xdr.TransactionEnvelope
							err := xdr.SafeUnmarshalBase64(operation.SerializedTxn, &txEnv)
							if err != nil {
								fmt.Printf("error parsing Stellar transaction: %v\n", err)
								break
							}

							// Get the first operation's destination
							if len(txEnv.Operations()) > 0 {
								if paymentOp, ok := txEnv.Operations()[0].Body.GetPaymentOp(); ok {
									destAddress = paymentOp.Destination.Address()
								}
							}
						case "algorand":
							txnBytes, err := base64.StdEncoding.DecodeString(operation.SerializedTxn)
							if err != nil {
								fmt.Printf("failed to decode serialized transaction: %v", err)
								break
							}
							var txn algorandTypes.Transaction
							err = msgpack.Decode(txnBytes, &txn)
							if err != nil {
								fmt.Printf("failed to deserialize transaction: %v", err)
								break
							}
							if txn.Type == algorandTypes.PaymentTx {
								destAddress = txn.PaymentTxnFields.Receiver.String()
							} else if txn.Type == algorandTypes.AssetTransferTx {
								destAddress = txn.AssetTransferTxnFields.AssetReceiver.String()
							} else {
								fmt.Printf("Unknown transaction type: %s\n", txn.Type)
								break
							}
						}

						// Verify the extracted destination matches the bridge wallet
						if destAddress == "" {
							fmt.Printf("Failed to extract destination address from %s transaction\n", chain.ChainType)
							validDestination = false
						} else {
							switch chain.ChainType {
							case "solana":
								validDestination = strings.EqualFold(destAddress, bridgeWallet.EDDSAPublicKey)
							case "aptos":
								validDestination = strings.EqualFold(destAddress, bridgeWallet.AptosEDDSAPublicKey)
							case "stellar":
								validDestination = strings.EqualFold(destAddress, bridgeWallet.StellarPublicKey)
							// add algorand case
							case "algorand":
								validDestination = strings.EqualFold(destAddress, bridgeWallet.AlgorandEDDSAPublicKey)
							}
						}

						if !validDestination {
							fmt.Println("Invalid bridge destination address for", chain.ChainType)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						signature, err := getSignature(intent, i)
						if err != nil {
							fmt.Printf("error getting signature: %+v\n", err)
							break
						}

						var txnHash string
						switch chain.ChainType {
						case "solana":
							txnHash, err = sendSolanaTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
						case "aptos":
							txnHash, err = aptos.SendAptosTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
						case "stellar":
							txnHash, err = stellar.SendStellarTxn(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
						case "algorand":
							txnHash, err = algorand.SendAlgorandTransaction(operation.SerializedTxn, operation.GenesisHash, signature)
						}

						if err != nil {
							fmt.Println(err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						var lockMetadata LockMetadata
						json.Unmarshal([]byte(operation.SolverMetadata), &lockMetadata)

						if lockMetadata.Lock {
							err := LockIdentity(lockSchema.Id)
							if err != nil {
								fmt.Println(err)
								break
							}
							UpdateOperationResult(operation.ID, OPERATION_STATUS_COMPLETED, txnHash)
						} else {
							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, txnHash)
						}
					}

				} else if operation.Type == OPERATION_TYPE_BRIDGE_DEPOSIT {
					depositOperation := intent.Operations[i-1]

					if i == 0 || !(depositOperation.Type == OPERATION_TYPE_SEND_TO_BRIDGE) {
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

					} else if depositOperation.KeyCurve == "eddsa" || depositOperation.KeyCurve == "aptos_eddsa" ||
						depositOperation.KeyCurve == "secp256k1" || depositOperation.KeyCurve == "stellar_eddsa" || depositOperation.KeyCurve == "algorand_eddsa" || depositOperation.KeyCurve == "ripple_eddsa" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							fmt.Println(err)
							break
						}

						var transfers []common.Transfer

						if chain.ChainType == "solana" {
							transfers, err = GetSolanaTransfers(depositOperation.ChainId, depositOperation.Result, HeliusApiKey)
							if err != nil {
								fmt.Println(err)
								break
							}
						}

						if chain.ChainType == "aptos" {
							transfers, err = aptos.GetAptosTransfers(depositOperation.ChainId, depositOperation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						}

						if chain.ChainType == "bitcoin" {
							transfers, _, err = GetBitcoinTransfers(depositOperation.ChainId, depositOperation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						}

						if chain.ChainType == "algorand" {
							transfers, err = algorand.GetAlgorandTransfers(depositOperation.GenesisHash, depositOperation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						}
						if chain.ChainType == "stellar" {
							transfers, err = stellar.GetStellarTransfers(depositOperation.ChainId, depositOperation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						}

						if chain.ChainType == "ripple" {
							transfers, err = ripple.GetRippleTransfers(depositOperation.ChainId, depositOperation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
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

					withdrawalChain, err := common.GetChain(operation.ChainId)

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

					if withdrawalChain.KeyCurve == "ecdsa" || withdrawalChain.KeyCurve == "secp256k1" {
						if withdrawalChain.ChainType == "bitcoin" {
							// handle bitcoin withdrawal
							var solverData map[string]interface{}
							if err := json.Unmarshal([]byte(burn.SolverOutput), &solverData); err != nil {
								fmt.Println("failed to parse solver output:", err)
								break
							}

							amount, ok := solverData["amount"].(string)
							if !ok {
								fmt.Println("amount not found in solver output")
								break
							}

							dataToSign, err := withdrawBitcoinGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWallet.BitcoinMainnetPublicKey,
								amount,
								user.BitcoinMainnetPublicKey,
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

							result, err := withdrawBitcoinTxn(
								withdrawalChain.ChainUrl,
								dataToSign,
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
						} else {
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
								// handle ERC20 token
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
					} else if withdrawalChain.KeyCurve == "stellar_eddsa" {
						wallet, err := GetWallet(intent.Identity, intent.IdentityCurve)
						if err != nil {
							fmt.Printf("error getting wallet: %+v\n", err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						if wallet.StellarPublicKey == "" {
							fmt.Printf("error: no Stellar public key found in wallet")
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						// Initialize Horizon client
						client := stellar.GetClient(withdrawalChain.ChainId, withdrawalChain.ChainUrl)
						if tokenToWithdraw == util.ZERO_ADDRESS {
							// Handle native XLM transfer
							txn, dataToSign, err := stellar.WithdrawStellarNativeGetSignature(
								client,
								bridgeWallet.StellarPublicKey,
								burn.SolverOutput,
								wallet.StellarPublicKey, // Use the wallet's Stellar public key
							)

							if err != nil {
								fmt.Printf("error withdrawing native XLM: %+v\n", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
							intent.Operations[i].SolverDataToSign = dataToSign

							signature, err := getSignature(intent, i)
							if err != nil {
								fmt.Printf("error getting signature: %+v\n", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							result, err := stellar.WithdrawStellarTxn(
								client,
								txn,
								signature,
							)

							if err != nil {
								fmt.Printf("error withdrawing Stellar: %+v\n", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
						} else {
							// Handle non-native Stellar asset transfer
							assetParts := strings.Split(tokenToWithdraw, ":")
							if len(assetParts) != 2 {
								fmt.Printf("invalid asset format: %s", tokenToWithdraw)
								break
							}

							assetCode := assetParts[0]
							assetIssuer := assetParts[1]
							fmt.Printf("assetCode: %+v\n", assetCode)
							fmt.Printf("assetIssuer: %+v\n", assetIssuer)

							txn, dataToSign, err := stellar.WithdrawStellarAssetGetSignature(
								client,
								bridgeWallet.StellarPublicKey,
								burn.SolverOutput,
								wallet.StellarPublicKey, // Use the wallet's Stellar public key
								assetCode,
								assetIssuer,
							)

							if err != nil {
								fmt.Printf("error withdrawing Stellar asset: %+v\n", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
							intent.Operations[i].SolverDataToSign = dataToSign

							signature, err := getSignature(intent, i)
							if err != nil {
								fmt.Printf("error getting signature: %+v\n", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							result, err := stellar.WithdrawStellarTxn(
								client,
								txn,
								signature,
							)

							if err != nil {
								fmt.Printf("error withdrawing Stellar: %+v\n", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
						}
						break
					} else if withdrawalChain.KeyCurve == "aptos_eddsa" {
						wallet, err := GetWallet(intent.Identity, intent.IdentityCurve)
						if err != nil {
							fmt.Printf("error getting public key: %v", err)
							break
						}
						if tokenToWithdraw == util.ZERO_ADDRESS {
							// handle native token
							transaction, dataToSign, err := aptos.WithdrawAptosNativeGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWallet.AptosEDDSAPublicKey,
								burn.SolverOutput,
								user.AptosEDDSAPublicKey,
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

							result, err := aptos.WithdrawAptosTxn(
								withdrawalChain.ChainUrl,
								transaction,
								wallet.AptosEDDSAPublicKey,
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
							transaction, dataToSign, err := aptos.WithdrawAptosTokenGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWallet.AptosEDDSAPublicKey,
								burn.SolverOutput,
								user.AptosEDDSAPublicKey,
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

							result, err := aptos.WithdrawAptosTxn(
								withdrawalChain.ChainUrl,
								transaction,
								wallet.AptosEDDSAPublicKey,
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

					} else if withdrawalChain.KeyCurve == "algorand_eddsa" {

						if tokenToWithdraw == util.ZERO_ADDRESS {
							// handle native ALGO token
							dataToSign, tx, err := algorand.WithdrawAlgorandNativeGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWallet.AlgorandEDDSAPublicKey,
								burn.SolverOutput,
								user.AlgorandEDDSAPublicKey,
							)

							if err != nil {
								fmt.Println(err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
							intent.Operations[i].SolverDataToSign = dataToSign

							signature, err := getSignature(intent, i)
							if err != nil {
								fmt.Println(err)
								break
							}

							result, err := algorand.WithdrawAlgorandTxn(
								withdrawalChain.ChainUrl,
								signature,
								tx,
							)

							if err != nil {
								fmt.Println(err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
						} else {
							// handle ASA (Algorand Standard Asset)
							dataToSign, tx, err := algorand.WithdrawAlgorandASAGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWallet.AlgorandEDDSAPublicKey,
								burn.SolverOutput,
								user.AlgorandEDDSAPublicKey,
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
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							result, err := algorand.WithdrawAlgorandTxn(
								withdrawalChain.ChainUrl,
								signature,
								tx,
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

					} else if withdrawalChain.KeyCurve == "ripple_eddsa" {
						if tokenToWithdraw == util.ZERO_ADDRESS {
							// handle native token
							transaction, dataToSign, err := ripple.WithdrawRippleNativeGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWallet.RippleEDDSAPublicKey,
								burn.SolverOutput,
								user.RippleEDDSAPublicKey,
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

							result, err := ripple.SendRippleTransaction(
								transaction,
								withdrawalChain.ChainId,
								withdrawalChain.KeyCurve,
								user.RippleEDDSAPublicKey,
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

							transaction, dataToSign, err := ripple.WithdrawRippleTokenGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWallet.RippleEDDSAPublicKey,
								burn.SolverOutput,
								user.RippleEDDSAPublicKey,
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

							result, err := ripple.SendRippleTransaction(
								transaction,
								withdrawalChain.ChainId,
								withdrawalChain.KeyCurve,
								user.RippleEDDSAPublicKey,
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
					if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "secp256k1" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							fmt.Println(err)
							break
						}

						if chain.ChainType == "bitcoin" {
							confirmed, err = checkBitcoinTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						} else {
							confirmed, err = checkEVMTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						}

					} else if operation.KeyCurve == "eddsa" || operation.KeyCurve == "aptos_eddsa" || operation.KeyCurve == "stellar_eddsa" || operation.KeyCurve == "algorand_eddsa" || operation.KeyCurve == "ripple_eddsa" {
						chId := operation.ChainId
						if chId == "" {
							chId = operation.GenesisHash
						}
						chain, err := common.GetChain(chId)
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

						if chain.ChainType == "aptos" {
							confirmed, err = aptos.CheckAptosTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						}

						if chain.ChainType == "algorand" {
							confirmed, err = algorand.CheckAlgorandTransactionConfirmed(operation.GenesisHash, operation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						}
						if chain.ChainType == "stellar" {
							confirmed, err = stellar.CheckStellarTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Printf("error checking Stellar transaction: %v", err)
								break
							}
						}

						if chain.ChainType == "ripple" {
							confirmed, err = ripple.CheckRippleTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Printf("error checking Ripple transaction: %v", err)
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
				} else if operation.Type == OPERATION_TYPE_SEND_TO_BRIDGE {
					confirmed := false
					chain, err := common.GetChain(operation.ChainId)
					if err != nil {
						fmt.Println(err)
						break
					}

					switch chain.ChainType {
					case "bitcoin":
						confirmed, err = checkBitcoinTransactionConfirmed(operation.ChainId, operation.Result)
					case "solana":
						confirmed, err = checkSolanaTransactionConfirmed(operation.ChainId, operation.Result)
					case "aptos":
						confirmed, err = aptos.CheckAptosTransactionConfirmed(operation.ChainId, operation.Result)
					case "stellar":
						confirmed, err = stellar.CheckStellarTransactionConfirmed(operation.ChainId, operation.Result)
					case "algorand":
						confirmed, err = algorand.CheckAlgorandTransactionConfirmed(operation.GenesisHash, operation.Result)
					default: // EVM chains
						confirmed, err = checkEVMTransactionConfirmed(operation.ChainId, operation.Result)
					}

					if err != nil {
						fmt.Printf("error checking %s transaction: %+v\n", chain.ChainType, err)
						break
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
					if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "secp256k1" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							fmt.Println(err)
							break
						}
						if chain.ChainType == "bitcoin" {
							confirmed, err = checkBitcoinTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						} else {
							confirmed, err = checkEVMTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						}
					} else if operation.KeyCurve == "eddsa" || operation.KeyCurve == "aptos_eddsa" || operation.KeyCurve == "stellar_eddsa" || operation.KeyCurve == "algorand_eddsa" || operation.KeyCurve == "ripple_eddsa" {
						chain, err := common.GetChain(operation.ChainId)
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

						if chain.ChainType == "aptos" {
							confirmed, err = aptos.CheckAptosTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						}

						if chain.ChainType == "stellar" {
							confirmed, err = stellar.CheckStellarTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Printf("error checking Stellar transaction: %+v\n", err)
								break
							}
						}

						if chain.ChainType == "algorand" {
							confirmed, err = algorand.CheckAlgorandTransactionConfirmed(operation.GenesisHash, operation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						}

						if chain.ChainType == "ripple" {
							confirmed, err = ripple.CheckRippleTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Printf("error checking Ripple transaction: %+v\n", err)
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
					if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "secp256k1" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							fmt.Println(err)
							break
						}
						if chain.ChainType == "bitcoin" {
							confirmed, err = checkBitcoinTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						} else {
							confirmed, err = checkEVMTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						}
					} else if operation.KeyCurve == "eddsa" || operation.KeyCurve == "aptos_eddsa" || operation.KeyCurve == "stellar_eddsa" || operation.KeyCurve == "algorand_eddsa" || operation.KeyCurve == "ripple_eddsa" {
						chain, err := common.GetChain(operation.ChainId)
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

						if chain.ChainType == "aptos" {
							confirmed, err = aptos.CheckAptosTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Println(err)
								break
							}
						}

						if chain.ChainType == "algorand" {
							confirmed, err = algorand.CheckAlgorandTransactionConfirmed(operation.GenesisHash, operation.Result)
							if err != nil {
								fmt.Printf("error checking Algorand transaction: %+v\n", err)
								break
							}
						}

						if chain.ChainType == "stellar" {
							confirmed, err = stellar.CheckStellarTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Printf("error checking Stellar transaction: %+v\n", err)
								break
							}
						}

						if chain.ChainType == "ripple" {
							confirmed, err = ripple.CheckRippleTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								fmt.Printf("error checking Ripple transaction: %+v\n", err)
								break
							}
						}

					}

					if !confirmed {
						break
					}

					// now unlock the identity if locked
					var withdrawMetadata WithdrawMetadata
					json.Unmarshal([]byte(operation.SolverMetadata), &withdrawMetadata)

					lockSchema, err := GetLock(intent.Identity, intent.IdentityCurve)
					if err != nil {
						fmt.Println(err)
						break
					}

					if withdrawMetadata.Unlock {
						depositOperation := intent.Operations[i-4]
						// check for confirmations
						confirmed = false
						if depositOperation.KeyCurve == "ecdsa" || depositOperation.KeyCurve == "secp256k1" {
							chain, err := common.GetChain(depositOperation.ChainId)
							if err != nil {
								fmt.Println(err)
								break
							}

							if chain.ChainType == "bitcoin" {
								txnConfirmed, err := checkBitcoinTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									fmt.Println(err)
									break
								}

								if txnConfirmed {
									confirmed = true
									err := UnlockIdentity(lockSchema.Id)
									if err != nil {
										fmt.Println(err)
										break
									}
								}
							} else {
								txnConfirmed, err := checkEVMTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									fmt.Println(err)
									break
								}

								if txnConfirmed {
									confirmed = true
									err := UnlockIdentity(lockSchema.Id)
									if err != nil {
										fmt.Println(err)
										break
									}
								}
							}
						} else if depositOperation.KeyCurve == "eddsa" || depositOperation.KeyCurve == "aptos_eddsa" || depositOperation.KeyCurve == "stellar_eddsa" || depositOperation.KeyCurve == "algorand_eddsa" || depositOperation.KeyCurve == "ripple_eddsa" {
							chain, err := common.GetChain(depositOperation.ChainId)
							if err != nil {
								fmt.Println(err)
								break
							}

							if chain.ChainType == "solana" {
								txnConfirmed, err := checkSolanaTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									fmt.Println(err)
									break
								}

								if txnConfirmed {
									confirmed = true
									err := UnlockIdentity(lockSchema.Id)
									if err != nil {
										fmt.Println(err)
										break
									}
								}
							}

							if chain.ChainType == "aptos" {
								txnConfirmed, err := aptos.CheckAptosTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									fmt.Println(err)
									break
								}
								if txnConfirmed {
									confirmed = true
									err := UnlockIdentity(lockSchema.Id)
									if err != nil {
										fmt.Println(err)
										break
									}
								}
							}

							if chain.ChainType == "algorand" {
								txnConfirmed, err := algorand.CheckAlgorandTransactionConfirmed(depositOperation.GenesisHash, depositOperation.Result)
								if err != nil {
									fmt.Println(err)
									break
								}
								if txnConfirmed {
									confirmed = true
									err := UnlockIdentity(lockSchema.Id)
									if err != nil {
										fmt.Println(err)
										break
									}
								}
							}
							if chain.ChainType == "stellar" {
								txnConfirmed, err := stellar.CheckStellarTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									fmt.Printf("error checking Stellar transaction: %+v\n", err)
									break
								}
								if txnConfirmed {
									confirmed = true
									err := UnlockIdentity(lockSchema.Id)
									if err != nil {
										fmt.Println(err)
										break
									}
								}
							}

							if chain.ChainType == "ripple" {
								txnConfirmed, err := ripple.CheckRippleTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									fmt.Printf("error checking Ripple transaction: %+v\n", err)
									break
								}

								if txnConfirmed {
									confirmed = true
									err := UnlockIdentity(lockSchema.Id)
									if err != nil {
										fmt.Println(err)
										break
									}
								}
							}
						}
					}

					if confirmed {
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_COMPLETED)
					}

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
		fmt.Printf("error creating request: %+v\n", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("error sending request: %+v\n", err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body: %+v\n", err)
		return "", err
	}

	var signatureResponse SignatureResponse
	err = json.Unmarshal(body, &signatureResponse)
	if err != nil {
		fmt.Printf("error unmarshalling response body: %+v\n", err)
		return "", err
	}

	return signatureResponse.Signature, nil
}

func checkEVMTransactionConfirmed(chainId string, txnHash string) (bool, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return false, err
	}

	client, err := ethclient.Dial(chain.ChainUrl)
	if err != nil {
		log.Fatal(err)
	}

	_, isPending, err := client.TransactionByHash(context.Background(), ethCommon.HexToHash(txnHash))
	if err != nil {
		return false, err
	}

	return !isPending, nil
}

func checkSolanaTransactionConfirmed(chainId string, txnHash string) (bool, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return false, err
	}

	c := rpc.New(chain.ChainUrl)

	signature, err := solana.SignatureFromBase58(txnHash)
	if err != nil {
		return false, err
	}

	// Regarding the deprecation of GetConfirmedTransaction in Solana-Core v2, this has been updated to use GetTransaction.
	// https://spl_governance.crates.io/docs/rpc/deprecated/getconfirmedtransaction
	_, err = c.GetTransaction(context.Background(), signature, &rpc.GetTransactionOpts{
		Commitment: rpc.CommitmentConfirmed,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func checkBitcoinTransactionConfirmed(chainId string, txnHash string) (bool, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return false, err
	}

	txn, err := FetchTransaction(chain.ChainUrl, txnHash)
	if err != nil {
		return false, err
	}

	// Assuming a transaction is confirmed if it has at least 3 confirmations
	if txn != nil && txn.Confirmations >= 3 {
		return true, nil
	}

	return false, nil
}

func sendEVMTransaction(serializedTxn string, chainId string, keyCurve string, dataToSign string, signatureHex string) (string, error) {
	chain, err := common.GetChain(chainId)
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

// sendSolanaTransaction submits a signed Solana transaction to the network
func sendSolanaTransaction(serializedTxn string, chainId string, keyCurve string, dataToSign string, signatureBase58 string) (string, error) {
	// Get chain configuration for RPC endpoint
	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", err
	}

	// Initialize Solana RPC client
	c := rpc.New(chain.ChainUrl)

	// Decode the base58-encoded transaction data
	// Solana transactions are serialized using a custom binary format and base58-encoded
	decodedTransactionData, err := base58.Decode(serializedTxn)
	if err != nil {
		fmt.Println("Error decoding transaction data:", err)
		return "", err
	}

	// Deserialize the binary data into a Solana transaction
	// This reconstructs the transaction object with all its instructions
	_tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(decodedTransactionData))
	if err != nil {
		return "", err
	}

	// Decode the base58-encoded signature and convert it to Solana's signature format
	// Solana uses 64-byte Ed25519 signatures
	sig, _ := base58.Decode(signatureBase58)
	signature := solana.SignatureFromBytes(sig)

	// Add the signature to the transaction
	// Solana transactions can have multiple signatures for multi-sig transactions
	_tx.Signatures = append(_tx.Signatures, signature)

	// Verify that all required signatures are present and valid
	// This checks signatures against the transaction data and account permissions
	err = _tx.VerifySignatures()
	if err != nil {
		fmt.Println("error during verification")
		fmt.Println(err)
		return "", err
	}

	// Submit the transaction to the Solana network
	// The returned hash can be used to track the transaction status
	hash, err := c.SendTransaction(context.Background(), _tx)
	if err != nil {
		fmt.Println("error during sending transaction")
		return "", err
	}

	// Return the transaction hash as a string
	return hash.String(), nil
}

func sendBitcoinTransaction(serializedTxn string, chainId string, keyCurve string, dataToSign string, signatureHex string) (string, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return "", err
	}

	// Decode the serialized transaction
	serializedTx, err := hex.DecodeString(serializedTxn)
	if err != nil {
		return "", fmt.Errorf("failed to decode serialized transaction: %v", err)
	}

	// Create a new bitcoin transaction
	var msgTx wire.MsgTx
	err = msgTx.Deserialize(bytes.NewReader(serializedTx))
	if err != nil {
		return "", fmt.Errorf("failed to deserialize transaction: %v", err)
	}

	// Decode the signature
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode signature: %v", err)
	}

	// Create signature script
	builder := txscript.NewScriptBuilder()
	builder.AddData(signature)
	builder.AddData([]byte(dataToSign))
	signatureScript, err := builder.Script()
	if err != nil {
		return "", fmt.Errorf("failed to build signature script: %v", err)
	}

	// Apply signature script to all inputs
	for i := range msgTx.TxIn {
		msgTx.TxIn[i].SignatureScript = signatureScript
	}

	// Serialize the signed transaction
	var signedTxBuffer bytes.Buffer
	if err := msgTx.Serialize(&signedTxBuffer); err != nil {
		return "", fmt.Errorf("failed to serialize signed transaction: %v", err)
	}

	// Convert to hex for broadcasting
	signedTxHex := hex.EncodeToString(signedTxBuffer.Bytes())

	// Prepare RPC request
	rpcURL := chain.ChainUrl
	rpcRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "sendrawtransaction",
		"params":  []interface{}{signedTxHex},
	}

	jsonData, err := json.Marshal(rpcRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal RPC request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", rpcURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	// if chain.ChainApiKey != "" {
	// 	req.Header.Set("Authorization", "Bearer "+chain.ChainApiKey)
	// }

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Parse response
	var rpcResponse struct {
		Result string    `json:"result"`
		Error  *RPCError `json:"error"`
	}

	if err := json.Unmarshal(body, &rpcResponse); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Check for RPC error
	if rpcResponse.Error != nil {
		return "", fmt.Errorf("RPC error: %v", rpcResponse.Error.Message)
	}

	// Return the transaction hash from RPC response
	return rpcResponse.Result, nil
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
