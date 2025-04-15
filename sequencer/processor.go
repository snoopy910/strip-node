package sequencer

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/StripChain/strip-node/ERC20"
	"github.com/StripChain/strip-node/algorand"
	"github.com/StripChain/strip-node/aptos"
	"github.com/StripChain/strip-node/bitcoin"
	"github.com/StripChain/strip-node/bridge"
	"github.com/StripChain/strip-node/cardano"
	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/dogecoin"
	"github.com/StripChain/strip-node/evm"
	"github.com/StripChain/strip-node/libs"
	db "github.com/StripChain/strip-node/libs/database"
	"github.com/StripChain/strip-node/ripple"
	"github.com/StripChain/strip-node/solana"
	"github.com/StripChain/strip-node/solver"
	"github.com/StripChain/strip-node/stellar"
	"github.com/StripChain/strip-node/sui"
	"github.com/StripChain/strip-node/util"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	algorandTypes "github.com/algorand/go-algorand-sdk/types"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/coming-chat/go-sui/v2/sui_types"
	cardanolib "github.com/echovl/cardano-go"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	bin "github.com/gagliardetto/binary"
	solanasdk "github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
	"github.com/rubblelabs/ripple/data"
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

func ProcessIntent(intentId int64) {
	for {
		intent, err := db.GetIntent(intentId)
		if err != nil {
			logger.Sugar().Errorw("error getting intent", "error", err)
			return
		}

		intentBytes, err := json.Marshal(intent)
		if err != nil {
			logger.Sugar().Errorw("error marshalling intent", "error", err)
			return
		}

		if intent.Status != db.INTENT_STATUS_PROCESSING {
			logger.Sugar().Infow("intent processed", "intent", intent)
			return
		}

		if intent.Expiry < uint64(time.Now().Unix()) {
			db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_EXPIRED)
			return
		}

		// now process the operations of the intent
		for i, operation := range intent.Operations {
			if operation.Status == db.OPERATION_STATUS_COMPLETED || operation.Status == db.OPERATION_STATUS_FAILED {
				continue
			}

			if operation.Status == db.OPERATION_STATUS_PENDING {
				// sign and send the txn. Change status to waiting

				if operation.Type == db.OPERATION_TYPE_TRANSACTION {
					lockSchema, _, err := verifyIdentityLockSchema(intent, &operation)
					if lockSchema == nil || err != nil {
						logger.Sugar().Errorw("error verifying identity lock", "error", err)
						break
					}

					if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "bitcoin_ecdsa" || operation.KeyCurve == "dogecoin_ecdsa" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							logger.Sugar().Errorw("error getting chain", "error", err)
							break
						}

						var txnHash string

						if chain.ChainType == "bitcoin" {
							signature, bitcoinPubkey, err := getSignatureEx(intent, i)
							if err != nil {
								logger.Sugar().Errorw("error getting signature", "error", err)
								break
							}
							txnHash, err = bitcoin.SendBitcoinTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, bitcoinPubkey, signature)

							if err != nil {
								logger.Sugar().Errorw("error sending bitcoin transaction", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}
						} else if chain.ChainType == "dogecoin" {
							signature, dogecoinPubKey, err := getSignatureEx(intent, i)
							if err != nil {
								logger.Sugar().Errorw("error getting signature", "error", err)
								break
							}
							txnHash, err = dogecoin.SendDogeTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, dogecoinPubKey, signature)
							fmt.Println(txnHash)

							if err != nil {
								logger.Sugar().Errorw("error sending dogecoin transaction", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}
						} else {
							signature, err := getSignature(intent, i)
							if err != nil {
								logger.Sugar().Errorw("error getting signature", "error", err)
								break
							}

							txnHash, err = evm.SendEVMTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)

							// @TODO: For our infra errors, don't mark the intent and operation as failed
							if err != nil {
								logger.Sugar().Errorw("error sending evm transaction", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}
						}

						var lockMetadata LockMetadata
						json.Unmarshal([]byte(operation.SolverMetadata), &lockMetadata)

						if lockMetadata.Lock {
							err := db.LockIdentity(lockSchema.Id)
							if err != nil {
								logger.Sugar().Errorw("error locking identity", "error", err)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_COMPLETED, txnHash)
						} else {
							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, txnHash)
						}
					} else if operation.KeyCurve == "eddsa" || operation.KeyCurve == "aptos_eddsa" || operation.KeyCurve == "stellar_eddsa" || operation.KeyCurve == "algorand_eddsa" || operation.KeyCurve == "ripple_eddsa" || operation.KeyCurve == "cardano_eddsa" || operation.KeyCurve == "sui_eddsa" {
						chId := operation.ChainId
						if chId == "" {
							chId = operation.GenesisHash
						}
						chain, err := common.GetChain(chId)
						if err != nil {
							logger.Sugar().Errorw("error getting chain", "error", err)
							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
							break
						}

						signature, err := getSignature(intent, i)

						if err != nil {
							logger.Sugar().Errorw("error getting signature", "error", err)
							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
							break
						}

						var txnHash string

						if chain.ChainType == "solana" {
							txnHash, err = solana.SendSolanaTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
							if err != nil {
								logger.Sugar().Errorw("error sending solana transaction", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}
						}

						if chain.ChainType == "aptos" {
							// Convert public key
							wallet, err := db.GetWallet(intent.Identity, intent.IdentityCurve)
							if err != nil {
								logger.Sugar().Errorw("error getting public key", "error", err)
								break
							}
							txnHash, err = aptos.SendAptosTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, wallet.AptosEDDSAPublicKey, signature)
							if err != nil {
								logger.Sugar().Errorw("error sending aptos transaction", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}
						}

						if chain.ChainType == "algorand" {
							txnHash, err = algorand.SendAlgorandTransaction(operation.SerializedTxn, operation.GenesisHash, signature)
							if err != nil {
								logger.Sugar().Errorw("error sending algorand transaction", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

						}
						if chain.ChainType == "stellar" {
							// Send Stellar transaction
							txnHash, err = stellar.SendStellarTxn(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
							if err != nil {
								logger.Sugar().Errorw("error sending Stellar transaction", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}
						}

						if chain.ChainType == "ripple" {
							// Convert public key
							wallet, err := db.GetWallet(intent.Identity, intent.IdentityCurve)
							if err != nil {
								logger.Sugar().Errorw("error getting public key", "error", err)
								break
							}

							txnHash, err = ripple.SendRippleTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, wallet.RippleEDDSAPublicKey, signature)
							if err != nil {
								logger.Sugar().Errorw("error sending Ripple transaction", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}
						}

						if chain.ChainType == "cardano" {
							wallet, err := db.GetWallet(intent.Identity, intent.IdentityCurve)
							if err != nil {
								logger.Sugar().Errorw("error getting public key", "error", err)
								break
							}

							txnHash, err = cardano.SendCardanoTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, wallet.CardanoPublicKey, signature)
							if err != nil {
								logger.Sugar().Errorw("error sending Cardano transaction", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}
						}

						if chain.ChainType == "sui" {
							wallet, err := db.GetWallet(intent.Identity, intent.IdentityCurve)
							if err != nil {
								logger.Sugar().Errorw("error getting public key", "error", err)
								break
							}

							txnHash, err = sui.SendSuiTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, wallet.SuiPublicKey, signature)
							if err != nil {
								logger.Sugar().Errorw("error sending Sui transaction", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}
						}

						var lockMetadata LockMetadata
						json.Unmarshal([]byte(operation.SolverMetadata), &lockMetadata)

						if lockMetadata.Lock {
							err := db.LockIdentity(lockSchema.Id)
							if err != nil {
								logger.Sugar().Errorw("error locking identity", "error", err)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_COMPLETED, txnHash)
						} else {
							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, txnHash)
						}

					}
				} else if operation.Type == db.OPERATION_TYPE_SOLVER {
					lockSchema, _, err := verifyIdentityLockSchema(intent, &operation)
					if lockSchema == nil || err != nil {
						logger.Sugar().Errorw("error verifying identity lock", "error", err)
						break
					}

					// get data to sign from solver
					dataToSign, err := solver.Construct(operation.Solver, &intentBytes, i)

					if err != nil {
						logger.Sugar().Errorw("error constructing solver data to sign", "error", err)
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					db.UpdateOperationSolverDataToSign(operation.ID, dataToSign)

					// then get the data signed
					signature, err := getSignature(intent, i)
					if err != nil {
						logger.Sugar().Errorw("error getting signature", "error", err)
						break
					}

					// then send the signature to solver
					result, err := solver.Solve(
						operation.Solver, &intentBytes,
						i,
						signature,
					)

					if err != nil {
						logger.Sugar().Errorw("error solving", "error", err)
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
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
						db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_COMPLETED, result)
					} else {
						db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
					}
				} else if operation.Type == db.OPERATION_TYPE_SEND_TO_BRIDGE {
					// Get bridge wallet for the chain
					bridgeWallet, err := db.GetWallet(BridgeContractAddress, operation.KeyCurve)
					if err != nil {
						logger.Sugar().Errorw("Failed to get bridge wallet", "error", err)
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					// Process transaction based on key curve and chain type
					lockSchema, _, err := verifyIdentityLockSchema(intent, &operation)
					if lockSchema == nil || err != nil {
						logger.Sugar().Errorw("error verifying identity lock", "error", err)
						break
					}

					if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "bitcoin_ecdsa" || operation.KeyCurve == "dogecoin_ecdsa" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							logger.Sugar().Errorw("error getting chain", "error", err)
							break
						}

						// Extract destination address from serialized transaction
						var destAddress string
						if chain.ChainType == "bitcoin" {
							// For Bitcoin, decode the serialized transaction to get output address
							var tx wire.MsgTx
							txBytes, err := hex.DecodeString(operation.SerializedTxn)
							if err != nil {
								logger.Sugar().Errorw("error decoding bitcoin&dogecoin transaction", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}
							if err := tx.Deserialize(bytes.NewReader(txBytes)); err != nil {
								logger.Sugar().Errorw("error deserializing bitcoin&dogecoin transaction", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}
							// Get the first output's address (assuming it's the bridge address)
							if len(tx.TxOut) > 0 {
								_, addrs, _, err := txscript.ExtractPkScriptAddrs(tx.TxOut[0].PkScript, nil)
								if err != nil || len(addrs) == 0 {
									logger.Sugar().Errorw("error extracting bitcoin&dogecoin address", "error", err)
									db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
									db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
									break
								}
								destAddress = addrs[0].String()
							}
						} else {
							// For EVM chains, decode the transaction to get the 'to' address
							txBytes, err := hex.DecodeString(operation.SerializedTxn)
							if err != nil {
								logger.Sugar().Errorw("error decoding EVM transaction", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}
							tx := new(types.Transaction)
							if err := rlp.DecodeBytes(txBytes, tx); err != nil {
								logger.Sugar().Errorw("error deserializing EVM transaction", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}
							if tx.To() == nil {
								logger.Sugar().Errorw("EVM transaction has nil To address")
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}
							if len(tx.Data()) >= 4 && bytes.Equal(tx.Data()[:4], []byte{0xa9, 0x05, 0x9c, 0xbb}) {
								// ERC20 transfer detected, extract recipient from call data
								if len(tx.Data()) < 36 {
									logger.Sugar().Errorw("ERC20 transfer data too short to extract destination")
									db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
									db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
									break
								}
								destAddress = ethCommon.BytesToAddress(tx.Data()[4:36]).Hex()

								// For ERC20 transfers, verify the token exists in the bridge contract
								tokenAddress := tx.To().Hex()
								exists, peggedToken, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, operation.ChainId, tokenAddress)
								if err != nil {
									logger.Sugar().Errorw("error checking token existence in bridge", "error", err, "token", tokenAddress)
									db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
									db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
									break
								}
								if !exists {
									logger.Sugar().Errorw("ERC20 token not registered in bridge", "token", tokenAddress, "chain", operation.ChainId)
									db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
									db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
									break
								}
								logger.Sugar().Infow("ERC20 token exists in bridge", "token", tokenAddress, "peggedToken", peggedToken)
							} else {
								destAddress = tx.To().Hex()
							}
						}

						// Verify destination address matches bridge wallet
						var expectedAddress string
						if chain.ChainType == "bitcoin" {
							expectedAddress = bridgeWallet.BitcoinMainnetPublicKey
						} else if chain.ChainType == "dogecoin" {
							expectedAddress = bridgeWallet.DogecoinMainnetPublicKey
						} else {
							expectedAddress = bridgeWallet.ECDSAPublicKey
						}

						if !strings.EqualFold(destAddress, expectedAddress) {
							logger.Sugar().Errorw("Invalid bridge destination address", "expected", expectedAddress, "got", destAddress)
							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
							break
						}

						signature, err := getSignature(intent, i)
						if err != nil {
							logger.Sugar().Errorw("error getting signature", "error", err)
							break
						}

						var txnHash string
						switch chain.ChainType {
						case "bitcoin":
							signature, bitcoinPubkey, err_ := getSignatureEx(intent, i)
							if err_ != nil {
								logger.Sugar().Errorw("error getting signature", "error", err_)
								break
							}
							txnHash, err = bitcoin.SendBitcoinTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, bitcoinPubkey, signature)
						case "dogecoin":
							signature, dogecoinPubkey, err_ := getSignatureEx(intent, i)
							fmt.Println("Signature:", signature)
							fmt.Println("Dogecoin Public key:", dogecoinPubkey)
							if err_ != nil {
								fmt.Printf("error getting signature: %+v\n", err_)
								break
							}
							txnHash, err = dogecoin.SendDogeTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, dogecoinPubkey, signature)
						default: // EVM chains
							txnHash, err = evm.SendEVMTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
						}

						if err != nil {
							logger.Sugar().Errorw("error sending transaction", "error", err)
							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
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

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_COMPLETED, txnHash)
						} else {
							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, txnHash)
						}
					} else if operation.KeyCurve == "eddsa" || operation.KeyCurve == "aptos_eddsa" || operation.KeyCurve == "stellar_eddsa" || operation.KeyCurve == "algorand_eddsa" || operation.KeyCurve == "ripple_eddsa" || operation.KeyCurve == "cardano_eddsa" || operation.KeyCurve == "sui_eddsa" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							logger.Sugar().Errorw("error getting chain", "error", err)
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
								logger.Sugar().Errorw("error decoding Solana transaction", "error", err)
								break
							}
							tx, err := solanasdk.TransactionFromDecoder(bin.NewBinDecoder(decodedTxn))
							if err != nil || len(tx.Message.Instructions) == 0 {
								logger.Sugar().Errorw("error deserializing Solana transaction", "error", err)
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
								logger.Sugar().Errorw("error parsing Aptos transaction", "error", err)
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
								logger.Sugar().Errorw("error parsing Stellar transaction", "error", err)
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
								logger.Sugar().Errorw("failed to decode serialized transaction", "error", err)
								break
							}
							var txn algorandTypes.Transaction
							err = msgpack.Decode(txnBytes, &txn)
							if err != nil {
								logger.Sugar().Errorw("failed to deserialize transaction", "error", err)
								break
							}
							if txn.Type == algorandTypes.PaymentTx {
								destAddress = txn.PaymentTxnFields.Receiver.String()
							} else if txn.Type == algorandTypes.AssetTransferTx {
								destAddress = txn.AssetTransferTxnFields.AssetReceiver.String()
							} else {
								logger.Sugar().Errorw("Unknown transaction type", "type", txn.Type)
								break
							}
						case "ripple":
							// For Ripple, the destination is in the transaction payload
							// Decode the serialized transaction
							txBytes, err := hex.DecodeString(strings.TrimPrefix(operation.SerializedTxn, "0x"))
							if err != nil {
								logger.Sugar().Errorw("error decoding transaction", "error", err)
								break
							}

							// Parse the transaction
							var tx data.Payment
							err = json.Unmarshal(txBytes, &tx)
							if err != nil {
								logger.Sugar().Errorw("error unmarshalling transaction", "error", err)
								break
							}
							destAddress = tx.Destination.String()
						case "cardano":
							var tx cardanolib.Tx
							txBytes, err := hex.DecodeString(operation.SerializedTxn)
							if err != nil {
								logger.Sugar().Errorw("error decoding Cardano transaction", "error", err)
								break
							}
							if err := json.Unmarshal(txBytes, &tx); err != nil {
								logger.Sugar().Errorw("error parsing Cardano transaction", "error", err)
								break
							}
							destAddress = tx.Body.Outputs[0].Address.String()
						case "sui":
							var tx sui_types.TransactionData
							txBytes, err := base64.StdEncoding.DecodeString(operation.SerializedTxn)
							if err != nil {
								logger.Sugar().Errorw("error decoding Sui transaction", "error", err)
								break
							}
							if err := json.Unmarshal(txBytes, &tx); err != nil {
								logger.Sugar().Errorw("error parsing Sui transaction", "error", err)
								break
							}
							if len(tx.V1.Kind.ProgrammableTransaction.Inputs) < 1 {
								logger.Sugar().Errorw("wrong format sui transaction", "error", err)
								break
							}
							destAddress = string(*tx.V1.Kind.ProgrammableTransaction.Inputs[0].Pure)
						}

						// Verify the extracted destination matches the bridge wallet
						if destAddress == "" {
							logger.Sugar().Errorw("Failed to extract destination address from %s transaction", chain.ChainType)
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
							case "ripple":
								validDestination = strings.EqualFold(destAddress, bridgeWallet.RippleEDDSAPublicKey)
							case "cardano":
								validDestination = strings.EqualFold(destAddress, bridgeWallet.CardanoPublicKey)
							case "sui":
								validDestination = strings.EqualFold(destAddress, bridgeWallet.SuiPublicKey)
							}
						}

						if !validDestination {
							logger.Sugar().Errorw("Invalid bridge destination address for", "chain", chain.ChainType)
							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
							break
						}

						signature, err := getSignature(intent, i)
						if err != nil {
							logger.Sugar().Errorw("error getting signature", "error", err)
							break
						}

						var txnHash string
						switch chain.ChainType {
						case "solana":
							txnHash, err = solana.SendSolanaTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
						case "aptos":
							txnHash, err = aptos.SendAptosTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
						case "stellar":
							txnHash, err = stellar.SendStellarTxn(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
						case "algorand":
							txnHash, err = algorand.SendAlgorandTransaction(operation.SerializedTxn, operation.GenesisHash, signature)
						}

						if err != nil {
							logger.Sugar().Errorw("error sending transaction", "error", err)
							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
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
							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_COMPLETED, txnHash)
						} else {
							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, txnHash)
						}
					}

				} else if operation.Type == db.OPERATION_TYPE_BRIDGE_DEPOSIT {
					lockSchema, _, err := verifyIdentityLockSchema(intent, &operation)
					if lockSchema == nil || err != nil {
						logger.Sugar().Errorw("error verifying identity lock", "error", err)
						break
					}

					depositOperation := intent.Operations[i-1]

					if i == 0 || !(depositOperation.Type == db.OPERATION_TYPE_SEND_TO_BRIDGE) {
						logger.Sugar().Errorw("Invalid operation type for bridge deposit")
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					if depositOperation.KeyCurve == "ecdsa" {
						// find token transfer events and check if first transfer is a valid token
						transfers, err := evm.GetEthereumTransfers(depositOperation.ChainId, depositOperation.Result, intent.Identity)
						if err != nil {
							logger.Sugar().Errorw("error getting transfers", "error", err)
							break
						}

						if len(transfers) == 0 {
							logger.Sugar().Errorw("No transfers found", "result", depositOperation.Result, "identity", intent.Identity)
							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
							break
						}

						// check if the token exists
						transfer := transfers[0]
						srcAddress := transfer.TokenAddress
						amount := transfer.ScaledAmount

						exists, destAddress, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, depositOperation.ChainId, srcAddress)

						if err != nil {
							logger.Sugar().Errorw("error checking token existence", "error", err)
							break
						}

						if !exists {
							logger.Sugar().Errorw("Token does not exist", "srcAddress", srcAddress, "chainId", depositOperation.ChainId)

							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
							break
						}

						wallet, err := db.GetWallet(intent.Identity, intent.IdentityCurve)
						if err != nil {
							logger.Sugar().Errorw("error getting wallet", "error", err)
							break
						}

						dataToSign, err := bridge.BridgeDepositDataToSign(RPC_URL, BridgeContractAddress, amount, wallet.ECDSAPublicKey, destAddress)
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

						logger.Sugar().Infow("Minting bridge %s %s %s %s", amount, wallet.ECDSAPublicKey, destAddress, signature)

						result, err := mintBridge(
							amount, wallet.ECDSAPublicKey, destAddress, signature)

						if err != nil {
							logger.Sugar().Errorw("error minting bridge", "error", err)
							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
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

						db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)

					} else if depositOperation.KeyCurve == "eddsa" || depositOperation.KeyCurve == "aptos_eddsa" || depositOperation.KeyCurve == "sui_eddsa" ||
						depositOperation.KeyCurve == "bitcoin_ecdsa" || depositOperation.KeyCurve == "dogecoin_ecdsa" || depositOperation.KeyCurve == "stellar_eddsa" ||
						depositOperation.KeyCurve == "algorand_eddsa" || depositOperation.KeyCurve == "ripple_eddsa" || depositOperation.KeyCurve == "cardano_eddsa" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							logger.Sugar().Errorw("error getting chain", "error", err)
							break
						}

						var transfers []common.Transfer

						if chain.ChainType == "solana" {
							transfers, err = solana.GetSolanaTransfers(depositOperation.ChainId, depositOperation.Result, HeliusApiKey)
							if err != nil {
								logger.Sugar().Errorw("error getting solana transfers", "error", err)
								break
							}
						}

						if chain.ChainType == "dogecoin" {
							transfers, err = dogecoin.GetDogeTransfers(depositOperation.ChainId, depositOperation.Result)
							if err != nil {
								logger.Sugar().Errorw("error getting dogecoin transfers", "error", err)
								break
							}
						}

						if chain.ChainType == "aptos" {
							transfers, err = aptos.GetAptosTransfers(depositOperation.ChainId, depositOperation.Result)
							if err != nil {
								logger.Sugar().Errorw("error getting aptos transfers", "error", err)
								break
							}
						}

						if chain.ChainType == "bitcoin" {
							transfers, _, err = bitcoin.GetBitcoinTransfers(depositOperation.ChainId, depositOperation.Result)
							if err != nil {
								logger.Sugar().Errorw("error getting bitcoin transfers", "error", err)
								break
							}
						} else if chain.ChainType == "sui" {
							transfers, err = sui.GetSuiTransfers(depositOperation.ChainId, depositOperation.Result)
							if err != nil {
								logger.Sugar().Errorw("error getting sui transfers", "error", err)
								break
							}
						}

						if chain.ChainType == "algorand" {
							transfers, err = algorand.GetAlgorandTransfers(depositOperation.GenesisHash, depositOperation.Result)
							if err != nil {
								logger.Sugar().Errorw("error getting algorand transfers", "error", err)
								break
							}
						}
						if chain.ChainType == "stellar" {
							transfers, err = stellar.GetStellarTransfers(depositOperation.ChainId, depositOperation.Result)
							if err != nil {
								logger.Sugar().Errorw("error getting stellar transfers", "error", err)
								break
							}
						}

						if chain.ChainType == "ripple" {
							transfers, err = ripple.GetRippleTransfers(depositOperation.ChainId, depositOperation.Result)
							if err != nil {
								logger.Sugar().Errorw("error getting ripple transfers", "error", err)
								break
							}
						}

						if chain.ChainType == "cardano" {
							transfers, err = cardano.GetCardanoTransfers(depositOperation.ChainId, depositOperation.Result)
							if err != nil {
								logger.Sugar().Errorw("error getting cardano transfers", "error", err)
								break
							}
						}

						if len(transfers) == 0 {
							logger.Sugar().Errorw("No transfers found", "result", depositOperation.Result, "identity", intent.Identity)
							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
							break
						}

						// check if the token exists
						transfer := transfers[0]
						srcAddress := transfer.TokenAddress
						amount := transfer.ScaledAmount

						exists, destAddress, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, depositOperation.ChainId, srcAddress)

						if err != nil {
							logger.Sugar().Errorw("error checking token existence", "error", err)
							break
						}

						if !exists {
							logger.Sugar().Errorw("Token does not exist", "srcAddress", srcAddress, "chainId", depositOperation.ChainId)

							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
							break
						}

						wallet, err := db.GetWallet(intent.Identity, "ecdsa")
						if err != nil {
							logger.Sugar().Errorw("error getting wallet", "error", err)
							break
						}

						dataToSign, err := bridge.BridgeDepositDataToSign(RPC_URL, BridgeContractAddress, amount, wallet.ECDSAPublicKey, destAddress)
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

						logger.Sugar().Infow("Minting bridge %s %s %s %s", amount, wallet.ECDSAPublicKey, destAddress, signature)

						result, err := mintBridge(
							amount, wallet.ECDSAPublicKey, destAddress, signature)

						if err != nil {
							logger.Sugar().Errorw("error minting bridge", "error", err)
							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
							break
						}

						db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)

					}
				} else if operation.Type == db.OPERATION_TYPE_SWAP {
					lockSchema, _, err := verifyIdentityLockSchema(intent, &operation)
					if lockSchema == nil || err != nil {
						logger.Sugar().Errorw("error verifying identity lock", "error", err)
						break
					}
					bridgeDeposit := intent.Operations[i-1]

					if i == 0 || !(bridgeDeposit.Type == db.OPERATION_TYPE_BRIDGE_DEPOSIT) {
						logger.Sugar().Errorw("Invalid operation type for swap")
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					// Get the deposit operation details to find the actual source token address
					depositOperation := intent.Operations[i-2] // The operation before bridge deposit is send-to-bridge

					// Get the actual token used in the deposit from the transfer events
					transfers, err := evm.GetEthereumTransfers(depositOperation.ChainId, depositOperation.Result, intent.Identity)
					if err != nil {
						logger.Sugar().Errorw("error getting transfers for swap tokenIn", "error", err)
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					if len(transfers) == 0 {
						logger.Sugar().Errorw("No transfers found for swap tokenIn", "result", depositOperation.Result)
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					// Use the actual token address from the transfer event
					transfer := transfers[0]
					tokenIn := transfer.TokenAddress
					logger.Sugar().Infow("Using token from transfer event", "tokenIn", tokenIn)

					var bridgeDepositData MintOutput
					var swapMetadata SwapMetadata
					json.Unmarshal([]byte(bridgeDeposit.SolverOutput), &bridgeDepositData)
					json.Unmarshal([]byte(operation.SolverMetadata), &swapMetadata)

					// Use the token from metadata for tokenOut
					tokenOut := swapMetadata.Token
					amountIn := bridgeDepositData.Amount
					deadline := time.Now().Add(time.Hour).Unix()

					wallet, err := db.GetWallet(intent.Identity, "ecdsa")
					if err != nil {
						logger.Sugar().Errorw("error getting wallet", "error", err)
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

					logger.Sugar().Infow("Swapping bridge", "wallet", wallet.ECDSAPublicKey, "tokenIn", tokenIn, "tokenOut", tokenOut, "amountIn", amountIn, "deadline", deadline, "signature", signature)

					result, err := swapBridge(
						wallet.ECDSAPublicKey,
						tokenIn,
						tokenOut,
						amountIn,
						deadline,
						signature,
					)

					if err != nil {
						logger.Sugar().Errorw("error swapping bridge", "error", err)
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)

					break
				} else if operation.Type == db.OPERATION_TYPE_BURN {
					lockSchema, _, err := verifyIdentityLockSchema(intent, &operation)
					if lockSchema == nil || err != nil {
						logger.Sugar().Errorw("error verifying identity lock", "error", err)
						break
					}

					bridgeSwap := intent.Operations[i-1]

					if i+1 >= len(intent.Operations) || intent.Operations[i+1].Type != db.OPERATION_TYPE_WITHDRAW {
						fmt.Println("BURN operation must be followed by a WITHDRAW operation")
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					if i == 0 || !(bridgeSwap.Type == db.OPERATION_TYPE_SWAP) {
						logger.Sugar().Errorw("Invalid operation type for swap")
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

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
						wallet.ECDSAPublicKey,
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

					logger.Sugar().Infow("Burn tokens", "wallet", wallet.ECDSAPublicKey, "burnAmount", burnAmount, "token", burnMetadata.Token, "signature", signature)

					result, err := burnTokens(
						wallet.ECDSAPublicKey,
						burnAmount,
						burnMetadata.Token,
						signature,
					)

					if err != nil {
						logger.Sugar().Errorw("error burning tokens", "error", err)
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
					break
				} else if operation.Type == db.OPERATION_TYPE_BURN_SYNTHETIC {
					lockSchema, _, err := verifyIdentityLockSchema(intent, &operation)
					if lockSchema == nil || err != nil {
						logger.Sugar().Errorw("error verifying identity lock", "error", err)
						break
					}

					// This operation allows direct burning of ERC20 tokens from the wallet
					// without requiring a prior swap operation
					burnSyntheticMetadata := BurnSyntheticMetadata{}

					// Verify that this operation is followed by a withdraw operation
					if i+1 >= len(intent.Operations) || intent.Operations[i+1].Type != db.OPERATION_TYPE_WITHDRAW {
						logger.Sugar().Errorw("BURN_SYNTHETIC validation failed: must be followed by WITHDRAW",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"operationIndex", i,
							"totalOperations", len(intent.Operations),
							"nextOperationType", getNextOperationType(intent, i))
						fmt.Println("BURN_SYNTHETIC operation must be followed by a WITHDRAW operation")
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					logger.Sugar().Infow("BURN_SYNTHETIC validation passed: followed by WITHDRAW",
						"operationId", operation.ID,
						"nextOperationId", intent.Operations[i+1].ID,
						"nextOperationType", intent.Operations[i+1].Type)

					err := json.Unmarshal([]byte(operation.SolverMetadata), &burnSyntheticMetadata)
					if err != nil {
						logger.Sugar().Errorw("BURN_SYNTHETIC metadata parsing failed",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"error", err,
							"rawMetadata", operation.SolverMetadata)
						fmt.Println("Error unmarshalling burn synthetic metadata:", err)
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
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
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
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
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					logger.Sugar().Infow("BURN_SYNTHETIC wallet retrieved successfully",
						"operationId", operation.ID,
						"publicKey", wallet.ECDSAPublicKey)

					// Verify the user has sufficient token balance
					balance, err := ERC20.GetBalance(RPC_URL, burnSyntheticMetadata.Token, wallet.ECDSAPublicKey)
					if err != nil {
						logger.Sugar().Errorw("BURN_SYNTHETIC balance check failed",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"token", burnSyntheticMetadata.Token,
							"account", wallet.ECDSAPublicKey,
							"error", err)
						fmt.Println("Error getting token balance:", err)
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					logger.Sugar().Infow("BURN_SYNTHETIC balance retrieved successfully",
						"operationId", operation.ID,
						"token", burnSyntheticMetadata.Token,
						"account", wallet.ECDSAPublicKey,
						"balance", balance)

					balanceBig, ok := new(big.Int).SetString(balance, 10)
					if !ok {
						logger.Sugar().Errorw("BURN_SYNTHETIC balance parsing failed",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"balance", balance)
						fmt.Println("Error parsing balance")
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					amountBig, ok := new(big.Int).SetString(burnSyntheticMetadata.Amount, 10)
					if !ok {
						logger.Sugar().Errorw("BURN_SYNTHETIC amount parsing failed",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"amount", burnSyntheticMetadata.Amount)
						fmt.Println("Error parsing amount")
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					// Log balance check details
					logBurnSyntheticBalanceCheck(balanceBig, amountBig, burnSyntheticMetadata.Token)

					if balanceBig.Cmp(amountBig) < 0 {
						logger.Sugar().Errorw("BURN_SYNTHETIC insufficient balance",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"token", burnSyntheticMetadata.Token,
							"account", wallet.ECDSAPublicKey,
							"balance", balanceBig.String(),
							"requiredAmount", amountBig.String())
						fmt.Println("Insufficient token balance")
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
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
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					if exists {
						logger.Sugar().Errorw("BURN_SYNTHETIC invalid token: token exists on bridge",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"token", burnSyntheticMetadata.Token,
							"peggedToken", destAddress)
						fmt.Println("Invalid token: token exists on bridge, use BURN instead of BURN_SYNTHETIC")
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					logger.Sugar().Infow("BURN_SYNTHETIC token validated: not a bridged token",
						"operationId", operation.ID,
						"token", burnSyntheticMetadata.Token)

					// Generate data to sign for burning tokens
					dataToSign, err := bridge.BridgeBurnDataToSign(
						RPC_URL,
						BridgeContractAddress,
						wallet.ECDSAPublicKey,
						burnSyntheticMetadata.Amount,
						burnSyntheticMetadata.Token,
					)

					if err != nil {
						logger.Sugar().Errorw("BURN_SYNTHETIC data generation failed",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"token", burnSyntheticMetadata.Token,
							"account", wallet.ECDSAPublicKey,
							"bridgeContract", BridgeContractAddress,
							"error", err)
						fmt.Println("Error generating burn data to sign:", err)
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
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
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					// Log signature details
					logBurnSyntheticSignature(signature, dataToSign)

					// Verify signature locally before submitting
					isValidSignature, verifyErr := verifyBurnSyntheticSignature(dataToSign, signature, wallet.ECDSAPublicKey)
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
						"account", wallet.ECDSAPublicKey,
						"amount", burnSyntheticMetadata.Amount,
						"token", burnSyntheticMetadata.Token,
						"signatureLength", len(signature))

					fmt.Println("Burning synthetic tokens", wallet.ECDSAPublicKey, burnSyntheticMetadata.Amount, burnSyntheticMetadata.Token, signature)

					// Execute the burn transaction
					result, err := burnTokens(
						wallet.ECDSAPublicKey,
						burnSyntheticMetadata.Amount,
						burnSyntheticMetadata.Token,
						signature,
					)

					if err != nil {
						logger.Sugar().Errorw("BURN_SYNTHETIC transaction failed",
							"operationId", operation.ID,
							"intentId", intent.ID,
							"token", burnSyntheticMetadata.Token,
							"account", wallet.ECDSAPublicKey,
							"error", err)
						fmt.Println("Error burning tokens:", err)
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					logger.Sugar().Infow("BURN_SYNTHETIC transaction submitted successfully",
						"operationId", operation.ID,
						"intentId", intent.ID,
						"token", burnSyntheticMetadata.Token,
						"account", wallet.ECDSAPublicKey,
						"transactionHash", result)

					db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)

					// Log integration with next withdraw operation
					if i+1 < len(intent.Operations) && intent.Operations[i+1].Type == db.OPERATION_TYPE_WITHDRAW {
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

					break
				} else if operation.Type == db.OPERATION_TYPE_WITHDRAW {
					lockSchema, _, err := verifyIdentityLockSchema(intent, &operation)
					if lockSchema == nil || err != nil {
						logger.Sugar().Errorw("error verifying identity lock", "error", err)
						break
					}

					burn := intent.Operations[i-1]

					if i == 0 || !(burn.Type == db.OPERATION_TYPE_BURN || burn.Type == db.OPERATION_TYPE_BURN_SYNTHETIC) {
						logger.Sugar().Errorw("Invalid operation type for withdraw after burn")
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					var withdrawMetadata WithdrawMetadata
					json.Unmarshal([]byte(operation.SolverMetadata), &withdrawMetadata)

					// Handle different burn operation types
					var tokenToWithdraw string
					var burnTokenAddress string
					if burn.Type == db.OPERATION_TYPE_BURN {
						var burnMetadata BurnMetadata
						json.Unmarshal([]byte(burn.SolverMetadata), &burnMetadata)
						tokenToWithdraw = withdrawMetadata.Token
						burnTokenAddress = burnMetadata.Token
					} else if burn.Type == db.OPERATION_TYPE_BURN_SYNTHETIC {
						var burnSyntheticMetadata BurnSyntheticMetadata
						json.Unmarshal([]byte(burn.SolverMetadata), &burnSyntheticMetadata)
						tokenToWithdraw = withdrawMetadata.Token
						burnTokenAddress = burnSyntheticMetadata.Token
					}

					// verify these fields
					exists, destAddress, err := bridge.TokenExists(RPC_URL, BridgeContractAddress, operation.ChainId, tokenToWithdraw)

					if err != nil {
						logger.Sugar().Errorw("error checking token existence", "error", err)
						break
					}

					if !exists {
						logger.Sugar().Errorw("Token does not exist", "token", tokenToWithdraw, "chainId", operation.ChainId)

						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					if destAddress != burnTokenAddress {
						logger.Sugar().Errorw("Token mismatch", "destAddress", destAddress, "token", burnTokenAddress)

						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					withdrawalChain, err := common.GetChain(operation.ChainId)

					if err != nil {
						logger.Sugar().Errorw("error getting chain", "error", err)
						break
					}

					bridgeWallet, err := db.GetWallet(BridgeContractAddress, "ecdsa")
					if err != nil {
						logger.Sugar().Errorw("error getting bridge wallet", "error", err)
						break
					}

					user, err := db.GetWallet(intent.Identity, intent.IdentityCurve)
					if err != nil {
						logger.Sugar().Errorw("error getting user wallet", "error", err)
						break
					}

					if withdrawalChain.KeyCurve == "ecdsa" {
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

						result, err := withdrawEVMTxn(
							withdrawalChain.ChainUrl,
							signature,
							tx,
							operation.ChainId,
						)

						if err != nil {
							logger.Sugar().Errorw("error withdrawing ERC20", "error", err)
							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
							break
						}

						db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
						break
					} else if withdrawalChain.KeyCurve == "bitcoin_ecdsa" {
						bridgeWalletBitcoinAddress, err := readBitcoinAddress(bridgeWallet, withdrawalChain.ChainId)
						if err != nil {
							logger.Sugar().Errorw("error reading bitcoin address", "error", err)
							break
						}

						userBitcoinAddress, err := readBitcoinAddress(user, withdrawalChain.ChainId)
						if err != nil {
							logger.Sugar().Errorw("error reading bitcoin address", "error", err)
							break
						}

						// handle bitcoin withdrawal
						dataToSign, err := bitcoin.WithdrawBitcoinGetSignature(
							withdrawalChain.ChainId,
							bridgeWalletBitcoinAddress,
							burn.SolverOutput,
							userBitcoinAddress,
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

						result, err := bitcoin.WithdrawBitcoinTxn(
							withdrawalChain.ChainId,
							dataToSign,
							signature,
						)

						if err != nil {
							logger.Sugar().Errorw("error withdrawing bitcoin", "error", err)
							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
							break
						}

						db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
						break
					} else if withdrawalChain.KeyCurve == "dogecoin_ecdsa" {
						// handle dogecoin withdrawal
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

						// Get appropriate Dogecoin addresses based on network
						var userAddress, bridgeAddress string
						if withdrawalChain.ChainId == "2000" {
							userAddress = user.DogecoinMainnetPublicKey
							bridgeAddress = bridgeWallet.DogecoinMainnetPublicKey
						} else if withdrawalChain.ChainId == "2001" {
							userAddress = user.DogecoinTestnetPublicKey
							bridgeAddress = bridgeWallet.DogecoinTestnetPublicKey
						} else {
							fmt.Println("Invalid dogecoin chainID")
						}

						// Validate that we have the Dogecoin addresses
						if userAddress == "" || bridgeAddress == "" {
							fmt.Println("Dogecoin addresses not found in wallet")
							break
						}

						txn, dataToSign, err := dogecoin.WithdrawDogeNativeGetSignature(
							withdrawalChain.ChainUrl,
							bridgeAddress,
							amount,
							userAddress,
						)

						if err != nil {
							fmt.Println(err)
							break
						}

						db.UpdateOperationSolverDataToSign(operation.ID, dataToSign)
						intent.Operations[i].SolverDataToSign = dataToSign

						signature, err := getSignature(intent, i)
						if err != nil {
							fmt.Println(err)
							break
						}

						// Use the same Dogecoin address we used for signing
						result, err := dogecoin.WithdrawDogeTxn(
							withdrawalChain.ChainId,
							txn,         // Use the serialized transaction instead of dataToSign
							userAddress, // Use Dogecoin address instead of ECDSA key
							signature,
						)

						if err != nil {
							fmt.Println(err)
							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
							break
						}

						db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
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

							result, err := withdrawSolanaTxn(
								withdrawalChain.ChainUrl,
								transaction,
								signature,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing solana native", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
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

							result, err := withdrawSolanaTxn(
								withdrawalChain.ChainUrl,
								transaction,
								signature,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing solana SPL", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
						}
						break
					} else if withdrawalChain.KeyCurve == "stellar_eddsa" {
						wallet, err := db.GetWallet(intent.Identity, intent.IdentityCurve)
						if err != nil {
							logger.Sugar().Errorw("error getting wallet", "error", err)
							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
							break
						}

						if wallet.StellarPublicKey == "" {
							logger.Sugar().Errorw("error: no Stellar public key found in wallet")
							db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
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
								logger.Sugar().Errorw("error withdrawing native XLM", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationSolverDataToSign(operation.ID, dataToSign)
							intent.Operations[i].SolverDataToSign = dataToSign

							signature, err := getSignature(intent, i)
							if err != nil {
								logger.Sugar().Errorw("error getting signature", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							result, err := stellar.WithdrawStellarTxn(
								client,
								txn,
								signature,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing Stellar", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
						} else {
							// Handle non-native Stellar asset transfer
							assetParts := strings.Split(tokenToWithdraw, ":")
							if len(assetParts) != 2 {
								logger.Sugar().Errorw("invalid asset format", "asset", tokenToWithdraw)
								break
							}

							assetCode := assetParts[0]
							assetIssuer := assetParts[1]

							txn, dataToSign, err := stellar.WithdrawStellarAssetGetSignature(
								client,
								bridgeWallet.StellarPublicKey,
								burn.SolverOutput,
								wallet.StellarPublicKey, // Use the wallet's Stellar public key
								assetCode,
								assetIssuer,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing Stellar asset", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationSolverDataToSign(operation.ID, dataToSign)
							intent.Operations[i].SolverDataToSign = dataToSign

							signature, err := getSignature(intent, i)
							if err != nil {
								logger.Sugar().Errorw("error getting signature", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							result, err := stellar.WithdrawStellarTxn(
								client,
								txn,
								signature,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing Stellar", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
						}
						break
					} else if withdrawalChain.KeyCurve == "aptos_eddsa" {
						wallet, err := db.GetWallet(intent.Identity, intent.IdentityCurve)
						if err != nil {
							logger.Sugar().Errorw("error getting public key", "error", err)
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

							result, err := aptos.WithdrawAptosTxn(
								withdrawalChain.ChainUrl,
								transaction,
								wallet.AptosEDDSAPublicKey,
								signature,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing aptos", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
						} else {
							transaction, dataToSign, err := aptos.WithdrawAptosTokenGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWallet.AptosEDDSAPublicKey,
								burn.SolverOutput,
								user.AptosEDDSAPublicKey,
								tokenToWithdraw,
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

							result, err := aptos.WithdrawAptosTxn(
								withdrawalChain.ChainUrl,
								transaction,
								wallet.AptosEDDSAPublicKey,
								signature,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing aptos", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
						}
						break
					} else if withdrawalChain.KeyCurve == "sui_eddsa" {
						wallet, err := db.GetWallet(intent.Identity, intent.IdentityCurve)
						if err != nil {
							logger.Sugar().Errorw("error getting public key", "error", err)
							break
						}

						if tokenToWithdraw == util.ZERO_ADDRESS {
							// Handle native SUI token withdrawal
							transaction, dataToSign, err := sui.WithdrawSuiNativeGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWallet.SuiPublicKey,
								burn.SolverOutput,
								user.SuiPublicKey,
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

							result, err := sui.WithdrawSuiTxn(
								withdrawalChain.ChainUrl,
								transaction,
								wallet.SuiPublicKey,
								signature,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing sui", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
						} else {
							// Handle Sui token withdrawal
							transaction, dataToSign, err := sui.WithdrawSuiTokenGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWallet.SuiPublicKey,
								burn.SolverOutput,
								user.SuiPublicKey,
								tokenToWithdraw,
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

							result, err := sui.WithdrawSuiTxn(
								withdrawalChain.ChainUrl,
								transaction,
								wallet.SuiPublicKey,
								signature,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing sui", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
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
								logger.Sugar().Errorw("error getting data to sign", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationSolverDataToSign(operation.ID, dataToSign)
							intent.Operations[i].SolverDataToSign = dataToSign

							signature, err := getSignature(intent, i)
							if err != nil {
								logger.Sugar().Errorw("error getting signature", "error", err)
								break
							}

							result, err := algorand.WithdrawAlgorandTxn(
								withdrawalChain.ChainUrl,
								signature,
								tx,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing algorand", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
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
								logger.Sugar().Errorw("error getting data to sign", "error", err)
								break
							}

							db.UpdateOperationSolverDataToSign(operation.ID, dataToSign)
							intent.Operations[i].SolverDataToSign = dataToSign

							signature, err := getSignature(intent, i)
							if err != nil {
								logger.Sugar().Errorw("error getting signature", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							result, err := algorand.WithdrawAlgorandTxn(
								withdrawalChain.ChainUrl,
								signature,
								tx,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing algorand", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
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

							result, err := ripple.SendRippleTransaction(
								transaction,
								withdrawalChain.ChainId,
								withdrawalChain.KeyCurve,
								user.RippleEDDSAPublicKey,
								signature,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing ripple", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
						} else {

							transaction, dataToSign, err := ripple.WithdrawRippleTokenGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWallet.RippleEDDSAPublicKey,
								burn.SolverOutput,
								user.RippleEDDSAPublicKey,
								tokenToWithdraw,
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

							result, err := ripple.SendRippleTransaction(
								transaction,
								withdrawalChain.ChainId,
								withdrawalChain.KeyCurve,
								user.RippleEDDSAPublicKey,
								signature,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing ripple", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
						}
						break
					} else if withdrawalChain.KeyCurve == "cardano_eddsa" {
						bridgeWalletAddress := bridgeWallet.CardanoPublicKey
						userAddress := user.CardanoPublicKey
						if tokenToWithdraw == util.ZERO_ADDRESS {
							// handle native token
							transaction, dataToSign, err := cardano.WithdrawCardanoNativeGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWalletAddress,
								burn.SolverOutput,
								userAddress,
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

							result, err := cardano.SendCardanoTransaction(
								transaction,
								withdrawalChain.ChainId,
								withdrawalChain.KeyCurve,
								userAddress,
								signature,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing cardano", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
						} else {

							transaction, dataToSign, err := cardano.WithdrawCardanoTokenGetSignature(
								withdrawalChain.ChainUrl,
								bridgeWalletAddress,
								burn.SolverOutput,
								userAddress,
								tokenToWithdraw,
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

							result, err := cardano.SendCardanoTransaction(
								transaction,
								withdrawalChain.ChainId,
								withdrawalChain.KeyCurve,
								userAddress,
								signature,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing cardano", "error", err)
								db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
								db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
								break
							}

							db.UpdateOperationResult(operation.ID, db.OPERATION_STATUS_WAITING, result)
						}
						break
					}
				}

				break
			}

			if operation.Status == db.OPERATION_STATUS_WAITING {
				// check for confirmations and update the status to completed
				if operation.Type == db.OPERATION_TYPE_TRANSACTION {
					confirmed := false
					if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "bitcoin_ecdsa" || operation.KeyCurve == "dogecoin_ecdsa" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							logger.Sugar().Errorw("error getting chain", "error", err)
							break
						}

						if chain.ChainType == "bitcoin" {
							confirmed, err = bitcoin.CheckBitcoinTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking bitcoin transaction", "error", err)
								break
							}
						} else if chain.ChainType == "dogecoin" {
							confirmed, err = dogecoin.CheckDogeTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking dogecoin transaction", "error", err)
								break
							}
						} else {
							confirmed, err = evm.CheckEVMTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking evm transaction", "error", err)
								break
							}
						}

					} else if operation.KeyCurve == "eddsa" ||
						operation.KeyCurve == "aptos_eddsa" ||
						operation.KeyCurve == "stellar_eddsa" ||
						operation.KeyCurve == "algorand_eddsa" ||
						operation.KeyCurve == "sui_eddsa" ||
						operation.KeyCurve == "ripple_eddsa" ||
						operation.KeyCurve == "cardano_eddsa" {
						chId := operation.ChainId
						if chId == "" {
							chId = operation.GenesisHash
						}
						chain, err := common.GetChain(chId)
						if err != nil {
							logger.Sugar().Errorw("error getting chain", "error", err)
							break
						}

						if chain.ChainType == "solana" {
							confirmed, err = solana.CheckSolanaTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking solana transaction", "error", err)
								break
							}
						}

						if chain.ChainType == "aptos" {
							confirmed, err = aptos.CheckAptosTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking aptos transaction", "error", err)
								break
							}
						}

						if chain.ChainType == "sui" {
							confirmed, err = sui.CheckSuiTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking sui transaction", "error", err)
								break
							}
						}
						if chain.ChainType == "algorand" {
							confirmed, err = algorand.CheckAlgorandTransactionConfirmed(operation.GenesisHash, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking algorand transaction", "error", err)
								break
							}
						}
						if chain.ChainType == "stellar" {
							confirmed, err = stellar.CheckStellarTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking Stellar transaction", "error", err)
								break
							}
						}

						if chain.ChainType == "ripple" {
							confirmed, err = ripple.CheckRippleTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking Ripple transaction", "error", err)
								break
							}
						}

						if chain.ChainType == "cardano" {
							confirmed, err = cardano.CheckCardanoTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking Cardano transaction", "error", err)
								break
							}
						}
					}

					if !confirmed {
						break
					}

					db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_COMPLETED)

					if i+1 == len(intent.Operations) {
						// update the intent status to completed
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_COMPLETED)
					}

					break
				} else if operation.Type == db.OPERATION_TYPE_SEND_TO_BRIDGE {
					confirmed := false
					chain, err := common.GetChain(operation.ChainId)
					if err != nil {
						logger.Sugar().Errorw("error getting chain", "error", err)
						break
					}

					switch chain.ChainType {
					case "bitcoin":
						confirmed, err = bitcoin.CheckBitcoinTransactionConfirmed(operation.ChainId, operation.Result)
					case "dogecoin":
						confirmed, err = dogecoin.CheckDogeTransactionConfirmed(operation.ChainId, operation.Result)
					case "solana":
						confirmed, err = solana.CheckSolanaTransactionConfirmed(operation.ChainId, operation.Result)
					case "aptos":
						confirmed, err = aptos.CheckAptosTransactionConfirmed(operation.ChainId, operation.Result)
					case "stellar":
						confirmed, err = stellar.CheckStellarTransactionConfirmed(operation.ChainId, operation.Result)
					case "algorand":
						confirmed, err = algorand.CheckAlgorandTransactionConfirmed(operation.GenesisHash, operation.Result)
					default: // EVM chains
						confirmed, err = evm.CheckEVMTransactionConfirmed(operation.ChainId, operation.Result)
					}

					if err != nil {
						logger.Sugar().Errorw("error checking transaction", "chain", chain.ChainType, "error", err)
						break
					}

					if !confirmed {
						break
					}

					db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_COMPLETED)

					if i+1 == len(intent.Operations) {
						// update the intent status to completed
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_COMPLETED)
					}

					break
				} else if operation.Type == db.OPERATION_TYPE_BRIDGE_DEPOSIT {
					confirmed := false
					if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "bitcoin_ecdsa" || operation.KeyCurve == "dogecoin_ecdsa" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							logger.Sugar().Errorw("error getting chain", "error", err)
							break
						}
						if chain.ChainType == "bitcoin" {
							confirmed, err = bitcoin.CheckBitcoinTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking bitcoin transaction", "error", err)
								break
							}
						} else if chain.ChainType == "dogecoin" {
							confirmed, err = dogecoin.CheckDogeTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking dogecoin transaction", "error", err)
								break
							}
						} else {
							confirmed, err = evm.CheckEVMTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking evm transaction", "error", err)
								break
							}
						}
					} else if operation.KeyCurve == "eddsa" ||
						operation.KeyCurve == "aptos_eddsa" ||
						operation.KeyCurve == "stellar_eddsa" ||
						operation.KeyCurve == "algorand_eddsa" ||
						operation.KeyCurve == "sui_eddsa" ||
						operation.KeyCurve == "ripple_eddsa" ||
						operation.KeyCurve == "cardano_eddsa" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							logger.Sugar().Errorw("error getting chain", "error", err)
							break
						}

						if chain.ChainType == "solana" {
							confirmed, err = solana.CheckSolanaTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking solana transaction", "error", err)
								break
							}
						}

						if chain.ChainType == "aptos" {
							confirmed, err = aptos.CheckAptosTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking aptos transaction", "error", err)
								break
							}
						}

						if chain.ChainType == "sui" {
							confirmed, err = sui.CheckSuiTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking sui transaction", "error", err)
								break
							}
						}
						if chain.ChainType == "stellar" {
							confirmed, err = stellar.CheckStellarTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking Stellar transaction", "error", err)
								break
							}
						}

						if chain.ChainType == "algorand" {
							confirmed, err = algorand.CheckAlgorandTransactionConfirmed(operation.GenesisHash, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking algorand transaction", "error", err)
								break
							}
						}

						if chain.ChainType == "ripple" {
							confirmed, err = ripple.CheckRippleTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking Ripple transaction", "error", err)
								break
							}
						}

						if chain.ChainType == "cardano" {
							confirmed, err = cardano.CheckCardanoTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking Cardano transaction", "error", err)
								break
							}
						}
					}

					if !confirmed {
						break
					}

					db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_COMPLETED)

					if i+1 == len(intent.Operations) {
						// update the intent status to completed
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_COMPLETED)
					}

					break
				} else if operation.Type == db.OPERATION_TYPE_SOLVER {
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

						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_COMPLETED)
						db.UpdateOperationSolverOutput(operation.ID, output)

						if i+1 == len(intent.Operations) {
							// update the intent status to completed
							db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_COMPLETED)
						}
					}

					if status == solver.SOLVER_OPERATION_STATUS_FAILURE {
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
					}

					break
				} else if operation.Type == db.OPERATION_TYPE_SWAP {
					confirmed, err := evm.CheckEVMTransactionConfirmed(operation.ChainId, operation.Result)
					if err != nil {
						logger.Sugar().Errorw("error checking evm transaction", "error", err)
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
						logger.Sugar().Errorw("error getting swap output", "error", err)
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_COMPLETED)
					db.UpdateOperationSolverOutput(operation.ID, swapOutput)

					if i+1 == len(intent.Operations) {
						// update the intent status to completed
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_COMPLETED)
					}

					break
				} else if operation.Type == db.OPERATION_TYPE_BURN {
					confirmed, err := evm.CheckEVMTransactionConfirmed(operation.ChainId, operation.Result)
					if err != nil {
						logger.Sugar().Errorw("error checking evm transaction", "error", err)
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
						logger.Sugar().Errorw("error getting burn output", "error", err)
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_COMPLETED)
					db.UpdateOperationSolverOutput(operation.ID, swapOutput)

					if i+1 == len(intent.Operations) {
						// update the intent status to completed
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_COMPLETED)
					}

					break
				} else if operation.Type == db.OPERATION_TYPE_BURN_SYNTHETIC {
					confirmed, err := evm.CheckEVMTransactionConfirmed(operation.ChainId, operation.Result)
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
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
						break
					}

					db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_COMPLETED)
					db.UpdateOperationSolverOutput(operation.ID, swapOutput)

					if i+1 == len(intent.Operations) {
						// update the intent status to completed
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_COMPLETED)
					}

					break
				} else if operation.Type == db.OPERATION_TYPE_WITHDRAW {
					confirmed := false
					if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "bitcoin_ecdsa" || operation.KeyCurve == "dogecoin_ecdsa" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							logger.Sugar().Errorw("error getting chain", "error", err)
							break
						}
						if chain.ChainType == "bitcoin" {
							confirmed, err = bitcoin.CheckBitcoinTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking bitcoin transaction", "error", err)
								break
							}
						} else if chain.ChainType == "dogecoin" {
							confirmed, err = dogecoin.CheckDogeTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking dogecoin transaction", "error", err)
								break
							}
						} else {
							confirmed, err = evm.CheckEVMTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking evm transaction", "error", err)
								break
							}
						}
					} else if operation.KeyCurve == "eddsa" ||
						operation.KeyCurve == "aptos_eddsa" ||
						operation.KeyCurve == "stellar_eddsa" ||
						operation.KeyCurve == "algorand_eddsa" ||
						operation.KeyCurve == "sui_eddsa" ||
						operation.KeyCurve == "ripple_eddsa" ||
						operation.KeyCurve == "cardano_eddsa" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							logger.Sugar().Errorw("error getting chain", "error", err)
							break
						}

						if chain.ChainType == "solana" {
							confirmed, err = solana.CheckSolanaTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking solana transaction", "error", err)
								break
							}
						}

						if chain.ChainType == "aptos" {
							confirmed, err = aptos.CheckAptosTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking aptos transaction", "error", err)
								break
							}
						}

						if chain.ChainType == "sui" {
							confirmed, err = sui.CheckSuiTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking sui transaction", "error", err)
								break
							}
						}
						if chain.ChainType == "algorand" {
							confirmed, err = algorand.CheckAlgorandTransactionConfirmed(operation.GenesisHash, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking Algorand transaction", "error", err)
								break
							}
						}

						if chain.ChainType == "stellar" {
							confirmed, err = stellar.CheckStellarTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking Stellar transaction", "error", err)
								break
							}
						}

						if chain.ChainType == "ripple" {
							confirmed, err = ripple.CheckRippleTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking Ripple transaction", "error", err)
								break
							}
						}

						if chain.ChainType == "cardano" {
							confirmed, err = cardano.CheckCardanoTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking Cardano transaction", "error", err)
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

					lockSchema, err := db.GetLock(intent.Identity, intent.IdentityCurve)
					if err != nil {
						logger.Sugar().Errorw("error getting lock", "error", err)
						break
					}

					if withdrawMetadata.Unlock {
						depositOperation := intent.Operations[i-4]
						// check for confirmations
						confirmed = false
						if depositOperation.KeyCurve == "ecdsa" || depositOperation.KeyCurve == "bitcoin_ecdsa" || operation.KeyCurve == "dogecoin_ecdsa" {
							chain, err := common.GetChain(depositOperation.ChainId)
							if err != nil {
								logger.Sugar().Errorw("error getting chain", "error", err)
								break
							}

							if chain.ChainType == "bitcoin" {
								txnConfirmed, err := bitcoin.CheckBitcoinTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									logger.Sugar().Errorw("error checking bitcoin transaction", "error", err)
									break
								}

								if txnConfirmed {
									confirmed = true
									err := db.UnlockIdentity(lockSchema.Id)
									if err != nil {
										logger.Sugar().Errorw("error unlocking identity", "error", err)
										break
									}
								}
							} else if chain.ChainType == "dogecoin" {
								txnConfirmed, err := dogecoin.CheckDogeTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									logger.Sugar().Errorw("error checking dogecoin transaction", "error", err)
									break
								}

								if txnConfirmed {
									confirmed = true
									err := db.UnlockIdentity(lockSchema.Id)
									if err != nil {
										logger.Sugar().Errorw("error unlocking identity", "error", err)
										break
									}
								}
							} else {
								txnConfirmed, err := evm.CheckEVMTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									logger.Sugar().Errorw("error checking evm transaction", "error", err)
									break
								}

								if txnConfirmed {
									confirmed = true
									err := db.UnlockIdentity(lockSchema.Id)
									if err != nil {
										logger.Sugar().Errorw("error unlocking identity", "error", err)
										break
									}
								}
							}
						} else if depositOperation.KeyCurve == "eddsa" ||
							depositOperation.KeyCurve == "aptos_eddsa" ||
							depositOperation.KeyCurve == "stellar_eddsa" ||
							depositOperation.KeyCurve == "algorand_eddsa" ||
							depositOperation.KeyCurve == "sui_eddsa" ||
							depositOperation.KeyCurve == "ripple_eddsa" ||
							depositOperation.KeyCurve == "cardano_eddsa" {
							chain, err := common.GetChain(depositOperation.ChainId)
							if err != nil {
								logger.Sugar().Errorw("error getting chain", "error", err)
								break
							}

							if chain.ChainType == "solana" {
								txnConfirmed, err := solana.CheckSolanaTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									logger.Sugar().Errorw("error checking solana transaction", "error", err)
									break
								}

								if txnConfirmed {
									confirmed = true
									err := db.UnlockIdentity(lockSchema.Id)
									if err != nil {
										logger.Sugar().Errorw("error unlocking identity", "error", err)
										break
									}
								}
							}

							if chain.ChainType == "aptos" {
								txnConfirmed, err := aptos.CheckAptosTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									logger.Sugar().Errorw("error checking aptos transaction", "error", err)
									break
								}
								if txnConfirmed {
									confirmed = true
									err := db.UnlockIdentity(lockSchema.Id)
									if err != nil {
										logger.Sugar().Errorw("error unlocking identity", "error", err)
										break
									}
								}
							}

							if chain.ChainType == "algorand" {
								txnConfirmed, err := algorand.CheckAlgorandTransactionConfirmed(depositOperation.GenesisHash, depositOperation.Result)
								if err != nil {
									logger.Sugar().Errorw("error checking algorand transaction", "error", err)
									break
								}
								if txnConfirmed {
									confirmed = true
									err := db.UnlockIdentity(lockSchema.Id)
									if err != nil {
										logger.Sugar().Errorw("error unlocking identity", "error", err)
										break
									}
								}
							}
							if chain.ChainType == "stellar" {
								txnConfirmed, err := stellar.CheckStellarTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									logger.Sugar().Errorw("error checking Stellar transaction", "error", err)
									break
								}
								if txnConfirmed {
									confirmed = true
									err := db.UnlockIdentity(lockSchema.Id)
									if err != nil {
										logger.Sugar().Errorw("error unlocking identity", "error", err)
										break
									}
								}
							}

							if chain.ChainType == "sui" {
								txnConfirmed, err := sui.CheckSuiTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									logger.Sugar().Errorw("error checking sui transaction", "error", err)
									break
								}

								if txnConfirmed {
									confirmed = true
									err := db.UnlockIdentity(lockSchema.Id)
									if err != nil {
										logger.Sugar().Errorw("error unlocking identity", "error", err)
										break
									}
								}
							}

							if chain.ChainType == "ripple" {
								txnConfirmed, err := ripple.CheckRippleTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									logger.Sugar().Errorw("error checking Ripple transaction", "error", err)
									break
								}

								if txnConfirmed {
									confirmed = true
									err := db.UnlockIdentity(lockSchema.Id)
									if err != nil {
										logger.Sugar().Errorw("error unlocking identity", "error", err)
										break
									}
								}
							}

							if chain.ChainType == "cardano" {
								txnConfirmed, err := cardano.CheckCardanoTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									logger.Sugar().Errorw("error checking Cardano transaction", "error", err)
									break
								}

								if txnConfirmed {
									confirmed = true
									err := db.UnlockIdentity(lockSchema.Id)
									if err != nil {
										logger.Sugar().Errorw("error unlocking identity", "error", err)
										break
									}

								}
							}
						}
					}

					if confirmed {
						db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_COMPLETED)
					}

					if i+1 == len(intent.Operations) {
						// update the intent status to completed
						db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_COMPLETED)
					}

					break
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func verifyIdentityLockSchema(intent *libs.Intent, operation *libs.Operation) (*db.LockSchema, bool, error) {
	lockSchema, err := db.GetLock(intent.Identity, intent.IdentityCurve)
	if err != nil {
		if err.Error() == "pg: no rows in result set" {
			_, err := db.AddLock(intent.Identity, intent.IdentityCurve)

			if err != nil {
				logger.Sugar().Errorw("error adding lock", "error", err)
				return nil, false, err
			}

			lockSchema, err = db.GetLock(intent.Identity, intent.IdentityCurve)

			if err != nil {
				logger.Sugar().Errorw("error getting lock after adding", "error", err)
				return nil, false, err
			}
		} else {
			logger.Sugar().Errorw("error getting lock", "error", err)
			return nil, false, err
		}
	}

	if lockSchema.Locked {
		db.UpdateOperationStatus(operation.ID, db.OPERATION_STATUS_FAILED)
		db.UpdateIntentStatus(intent.ID, db.INTENT_STATUS_FAILED)
		return nil, false, fmt.Errorf("identity is locked")
	}

	return lockSchema, true, nil
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
	wallet, err := db.GetWallet(intent.Identity, intent.IdentityCurve)
	if err != nil {
		return "", "", fmt.Errorf("error getting wallet: %v", err)
	}

	// get the signer
	signers := strings.Split(wallet.Signers, ",")
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

func getNextOperationType(intent *libs.Intent, operationIndex int) string {
	if operationIndex+1 < len(intent.Operations) {
		return intent.Operations[operationIndex+1].Type
	}
	return "none"
}
