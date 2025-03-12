package sequencer

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/StripChain/strip-node/algorand"
	"github.com/StripChain/strip-node/aptos"
	"github.com/StripChain/strip-node/bitcoin"
	"github.com/StripChain/strip-node/bridge"
	"github.com/StripChain/strip-node/cardano"
	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/dogecoin"
	"github.com/StripChain/strip-node/ripple"
	"github.com/StripChain/strip-node/solver"
	"github.com/StripChain/strip-node/stellar"
	"github.com/StripChain/strip-node/sui"
	"github.com/StripChain/strip-node/util"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	algorandTypes "github.com/algorand/go-algorand-sdk/types"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	cardanolib "github.com/echovl/cardano-go"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
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
			logger.Sugar().Errorw("error getting intent", "error", err)
			return
		}

		intentBytes, err := json.Marshal(intent)
		if err != nil {
			logger.Sugar().Errorw("error marshalling intent", "error", err)
			return
		}

		if intent.Status != INTENT_STATUS_PROCESSING {
			logger.Sugar().Infow("intent processed", "intent", intent)
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
								logger.Sugar().Errorw("error adding lock", "error", err)
								break
							}

							lockSchema, err = GetLock(intent.Identity, intent.IdentityCurve)

							if err != nil {
								logger.Sugar().Errorw("error getting lock after adding", "error", err)
								break
							}
						} else {
							logger.Sugar().Errorw("error getting lock", "error", err)
							break
						}
					}

					if lockSchema.Locked {
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "bitcoin_ecdsa" {
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
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
						} else if chain.ChainType == "dogecoin" {
							signature, err := getSignature(intent, i)
							if err != nil {
								logger.Sugar().Errorw("error getting signature", "error", err)
								break
							}
							txnHash, err = dogecoin.SendDogeTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)

							if err != nil {
								logger.Sugar().Errorw("error sending dogecoin transaction", "error", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
						} else {
							signature, err := getSignature(intent, i)
							if err != nil {
								logger.Sugar().Errorw("error getting signature", "error", err)
								break
							}

							txnHash, err = sendEVMTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)

							// @TODO: For our infra errors, don't mark the intent and operation as failed
							if err != nil {
								logger.Sugar().Errorw("error sending evm transaction", "error", err)
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
								logger.Sugar().Errorw("error locking identity", "error", err)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_COMPLETED, txnHash)
						} else {
							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, txnHash)
						}
					} else if operation.KeyCurve == "sui_eddsa" {
						signature, err := getSignature(intent, i)
						if err != nil {
							logger.Sugar().Errorw("error getting signature", "error", err)
							break
						}

						txnHash, err := sui.SendSuiTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
						if err != nil {
							logger.Sugar().Errorw("error sending sui transaction", "error", err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						var lockMetadata LockMetadata
						json.Unmarshal([]byte(operation.SolverMetadata), &lockMetadata)

						if lockMetadata.Lock {
							err := LockIdentity(lockSchema.Id)
							if err != nil {
								logger.Sugar().Errorw("error locking identity", "error", err)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_COMPLETED, txnHash)
						} else {
							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, txnHash)
						}
					} else if operation.KeyCurve == "eddsa" || operation.KeyCurve == "aptos_eddsa" || operation.KeyCurve == "stellar_eddsa" || operation.KeyCurve == "algorand_eddsa" || operation.KeyCurve == "ripple_eddsa" || operation.KeyCurve == "cardano_eddsa" {
						chId := operation.ChainId
						if chId == "" {
							chId = operation.GenesisHash
						}
						chain, err := common.GetChain(chId)
						if err != nil {
							logger.Sugar().Errorw("error getting chain", "error", err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						signature, err := getSignature(intent, i)

						if err != nil {
							logger.Sugar().Errorw("error getting signature", "error", err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						var txnHash string

						if chain.ChainType == "solana" {
							txnHash, err = sendSolanaTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
							if err != nil {
								logger.Sugar().Errorw("error sending solana transaction", "error", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
						}

						if chain.ChainType == "aptos" {
							// Convert public key
							wallet, err := GetWallet(intent.Identity, intent.IdentityCurve)
							if err != nil {
								logger.Sugar().Errorw("error getting public key", "error", err)
								break
							}
							txnHash, err = aptos.SendAptosTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, wallet.AptosEDDSAPublicKey, signature)
							if err != nil {
								logger.Sugar().Errorw("error sending aptos transaction", "error", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
						}

						if chain.ChainType == "algorand" {
							txnHash, err = algorand.SendAlgorandTransaction(operation.SerializedTxn, operation.GenesisHash, signature)
							if err != nil {
								logger.Sugar().Errorw("error sending algorand transaction", "error", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

						}
						if chain.ChainType == "stellar" {
							// Send Stellar transaction
							txnHash, err = stellar.SendStellarTxn(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
							if err != nil {
								logger.Sugar().Errorw("error sending Stellar transaction", "error", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
						}

						if chain.ChainType == "ripple" {
							// Convert public key
							wallet, err := GetWallet(intent.Identity, intent.IdentityCurve)
							if err != nil {
								logger.Sugar().Errorw("error getting public key", "error", err)
								break
							}

							txnHash, err = ripple.SendRippleTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, wallet.RippleEDDSAPublicKey, signature)
							if err != nil {
								logger.Sugar().Errorw("error sending Ripple transaction", "error", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
						}

						if chain.ChainType == "cardano" {
							wallet, err := GetWallet(intent.Identity, intent.IdentityCurve)
							if err != nil {
								logger.Sugar().Errorw("error getting public key", "error", err)
								break
							}

							txnHash, err = cardano.SendCardanoTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, wallet.CardanoPublicKey, signature)
							if err != nil {
								logger.Sugar().Errorw("error sending Cardano transaction", "error", err)
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
								logger.Sugar().Errorw("error locking identity", "error", err)
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
								logger.Sugar().Errorw("error adding lock", "error", err)
								break
							}

							lockSchema, err = GetLock(intent.Identity, intent.IdentityCurve)

							if err != nil {
								logger.Sugar().Errorw("error getting lock after adding", "error", err)
								break
							}
						} else {
							logger.Sugar().Errorw("error getting lock", "error", err)
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
						logger.Sugar().Errorw("error constructing solver data to sign", "error", err)
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					UpdateOperationSolverDataToSign(operation.ID, dataToSign)

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
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					var lockMetadata LockMetadata
					json.Unmarshal([]byte(operation.SolverMetadata), &lockMetadata)

					if lockMetadata.Lock {
						err := LockIdentity(lockSchema.Id)
						if err != nil {
							logger.Sugar().Errorw("error locking identity", "error", err)
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
						logger.Sugar().Errorw("Failed to get bridge wallet", "error", err)
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
								logger.Sugar().Errorw("error adding lock", "error", err)
								break
							}
							lockSchema, err = GetLock(intent.Identity, intent.IdentityCurve)
							if err != nil {
								logger.Sugar().Errorw("error getting lock after adding", "error", err)
								break
							}
						} else {
							logger.Sugar().Errorw("error getting lock", "error", err)
							break
						}
					}

					if lockSchema.Locked {
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "bitcoin_ecdsa" {
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
								logger.Sugar().Errorw("error decoding bitcoin transaction", "error", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
							if err := tx.Deserialize(bytes.NewReader(txBytes)); err != nil {
								logger.Sugar().Errorw("error deserializing bitcoin transaction", "error", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
							// Get the first output's address (assuming it's the bridge address)
							if len(tx.TxOut) > 0 {
								_, addrs, _, err := txscript.ExtractPkScriptAddrs(tx.TxOut[0].PkScript, nil)
								if err != nil || len(addrs) == 0 {
									logger.Sugar().Errorw("error extracting bitcoin address", "error", err)
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
								logger.Sugar().Errorw("error decoding EVM transaction", "error", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}
							tx := new(types.Transaction)
							if err := rlp.DecodeBytes(txBytes, tx); err != nil {
								logger.Sugar().Errorw("error deserializing EVM transaction", "error", err)
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
							logger.Sugar().Errorw("Invalid bridge destination address", "expected", expectedAddress, "got", destAddress)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
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
						default: // EVM chains
							txnHash, err = sendEVMTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
						}

						if err != nil {
							logger.Sugar().Errorw("error sending transaction", "error", err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						var lockMetadata LockMetadata
						json.Unmarshal([]byte(operation.SolverMetadata), &lockMetadata)

						if lockMetadata.Lock {
							err := LockIdentity(lockSchema.Id)
							if err != nil {
								logger.Sugar().Errorw("error locking identity", "error", err)
								break
							}
							UpdateOperationResult(operation.ID, OPERATION_STATUS_COMPLETED, txnHash)
						} else {
							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, txnHash)
						}
					} else if operation.KeyCurve == "eddsa" || operation.KeyCurve == "aptos_eddsa" || operation.KeyCurve == "stellar_eddsa" || operation.KeyCurve == "algorand_eddsa" || operation.KeyCurve == "ripple_eddsa" || operation.KeyCurve == "cardano_eddsa" {
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
							tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(decodedTxn))
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
							}
						}

						if !validDestination {
							logger.Sugar().Errorw("Invalid bridge destination address for", "chain", chain.ChainType)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
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
							txnHash, err = sendSolanaTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
						case "aptos":
							txnHash, err = aptos.SendAptosTransaction(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
						case "stellar":
							txnHash, err = stellar.SendStellarTxn(operation.SerializedTxn, operation.ChainId, operation.KeyCurve, operation.DataToSign, signature)
						case "algorand":
							txnHash, err = algorand.SendAlgorandTransaction(operation.SerializedTxn, operation.GenesisHash, signature)
						}

						if err != nil {
							logger.Sugar().Errorw("error sending transaction", "error", err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						var lockMetadata LockMetadata
						json.Unmarshal([]byte(operation.SolverMetadata), &lockMetadata)

						if lockMetadata.Lock {
							err := LockIdentity(lockSchema.Id)
							if err != nil {
								logger.Sugar().Errorw("error locking identity", "error", err)
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
						logger.Sugar().Errorw("Invalid operation type for bridge deposit")
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					if depositOperation.KeyCurve == "ecdsa" {
						// find token transfer events and check if first transfer is a valid token
						transfers, err := GetEthereumTransfers(depositOperation.ChainId, depositOperation.Result, intent.Identity)
						if err != nil {
							logger.Sugar().Errorw("error getting transfers", "error", err)
							break
						}

						if len(transfers) == 0 {
							logger.Sugar().Errorw("No transfers found", "result", depositOperation.Result, "identity", intent.Identity)
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
							logger.Sugar().Errorw("error checking token existence", "error", err)
							break
						}

						if !exists {
							logger.Sugar().Errorw("Token does not exist", "srcAddress", srcAddress, "chainId", depositOperation.ChainId)

							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						wallet, err := GetWallet(intent.Identity, "ecdsa")
						if err != nil {
							logger.Sugar().Errorw("error getting wallet", "error", err)
							break
						}

						dataToSign, err := bridge.BridgeDepositDataToSign(RPC_URL, BridgeContractAddress, amount, wallet.ECDSAPublicKey, destAddress)
						if err != nil {
							logger.Sugar().Errorw("error getting data to sign", "error", err)
							break
						}

						UpdateOperationSolverDataToSign(operation.ID, dataToSign)
						intent.Operations[i].SolverDataToSign = dataToSign

						signature, err := getSignature(intent, i)
						if err != nil {
							logger.Sugar().Errorw("error getting signature", "error", err)
							break
						}

						logger.Sugar().Infof("Minting bridge %s %s %s %s", amount, wallet.ECDSAPublicKey, destAddress, signature)

						result, err := mintBridge(
							amount, wallet.ECDSAPublicKey, destAddress, signature)

						if err != nil {
							logger.Sugar().Errorw("error minting bridge", "error", err)
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
							logger.Sugar().Errorw("error marshalling mint output", "error", err)
							break
						}

						UpdateOperationSolverOutput(operation.ID, string(mintOutputBytes))

						UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)

					} else if depositOperation.KeyCurve == "eddsa" || depositOperation.KeyCurve == "aptos_eddsa" || depositOperation.KeyCurve == "sui_eddsa" ||
						depositOperation.KeyCurve == "bitcoin_ecdsa" || depositOperation.KeyCurve == "secp256k1" || depositOperation.KeyCurve == "stellar_eddsa" ||
						depositOperation.KeyCurve == "algorand_eddsa" || depositOperation.KeyCurve == "ripple_eddsa" || depositOperation.KeyCurve == "cardano_eddsa" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							logger.Sugar().Errorw("error getting chain", "error", err)
							break
						}

						var transfers []common.Transfer

						if chain.ChainType == "solana" {
							transfers, err = GetSolanaTransfers(depositOperation.ChainId, depositOperation.Result, HeliusApiKey)
							if err != nil {
								logger.Sugar().Errorw("error getting solana transfers", "error", err)
								break
							}
						} else if chain.ChainType == "dogecoin" {
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
							logger.Sugar().Errorw("error checking token existence", "error", err)
							break
						}

						if !exists {
							logger.Sugar().Errorw("Token does not exist", "srcAddress", srcAddress, "chainId", depositOperation.ChainId)

							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						wallet, err := GetWallet(intent.Identity, "ecdsa")
						if err != nil {
							logger.Sugar().Errorw("error getting wallet", "error", err)
							break
						}

						dataToSign, err := bridge.BridgeDepositDataToSign(RPC_URL, BridgeContractAddress, amount, wallet.ECDSAPublicKey, destAddress)
						if err != nil {
							logger.Sugar().Errorw("error getting data to sign", "error", err)
							break
						}

						UpdateOperationSolverDataToSign(operation.ID, dataToSign)
						intent.Operations[i].SolverDataToSign = dataToSign

						signature, err := getSignature(intent, i)
						if err != nil {
							logger.Sugar().Errorw("error getting signature", "error", err)
							break
						}

						result, err := mintBridge(
							amount, wallet.ECDSAPublicKey, destAddress, signature)

						if err != nil {
							logger.Sugar().Errorw("error minting bridge", "error", err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)

					}
				} else if operation.Type == OPERATION_TYPE_SWAP {
					bridgeDeposit := intent.Operations[i-1]

					if i == 0 || !(bridgeDeposit.Type == OPERATION_TYPE_BRIDGE_DEPOSIT) {
						logger.Sugar().Errorw("Invalid operation type for swap")
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

					UpdateOperationSolverDataToSign(operation.ID, dataToSign)
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
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)

					break
				} else if operation.Type == OPERATION_TYPE_BURN {
					bridgeSwap := intent.Operations[i-1]

					if i == 0 || !(bridgeSwap.Type == OPERATION_TYPE_SWAP) {
						logger.Sugar().Errorw("Invalid operation type for swap")
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					logger.Sugar().Infow("Burning tokens", "bridgeSwap", bridgeSwap)

					burnAmount := bridgeSwap.SolverOutput
					burnMetadata := BurnMetadata{}

					json.Unmarshal([]byte(operation.SolverMetadata), &burnMetadata)

					wallet, err := GetWallet(intent.Identity, "ecdsa")
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

					UpdateOperationSolverDataToSign(operation.ID, dataToSign)
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
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
					break
				} else if operation.Type == OPERATION_TYPE_WITHDRAW {
					burn := intent.Operations[i-1]

					if i == 0 || !(burn.Type == OPERATION_TYPE_BURN) {
						logger.Sugar().Errorw("Invalid operation type for withdraw after burn")
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
						logger.Sugar().Errorw("error checking token existence", "error", err)
						break
					}

					if !exists {
						logger.Sugar().Errorw("Token does not exist", "token", tokenToWithdraw, "chainId", operation.ChainId)

						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					if destAddress != burnMetadata.Token {
						logger.Sugar().Errorw("Token mismatch", "destAddress", destAddress, "token", burnMetadata.Token)

						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					withdrawalChain, err := common.GetChain(operation.ChainId)

					if err != nil {
						logger.Sugar().Errorw("error getting chain", "error", err)
						break
					}

					bridgeWallet, err := GetWallet(BridgeContractAddress, "ecdsa")
					if err != nil {
						logger.Sugar().Errorw("error getting bridge wallet", "error", err)
						break
					}

					user, err := GetWallet(intent.Identity, intent.IdentityCurve)
					if err != nil {
						logger.Sugar().Errorw("error getting user wallet", "error", err)
						break
					}

					if withdrawalChain.KeyCurve == "ecdsa" || withdrawalChain.KeyCurve == "secp256k1" {
						if withdrawalChain.ChainType == "dogecoin" {
							// handle dogecoin withdrawal
							var solverData map[string]interface{}
							if err := json.Unmarshal([]byte(burn.SolverOutput), &solverData); err != nil {
								logger.Sugar().Errorw("failed to parse solver output", "error", err)
								break
							}

							amount, ok := solverData["amount"].(string)
							if !ok {
								logger.Sugar().Errorw("amount not found in solver output")
								break
							}

							// Get appropriate Dogecoin addresses based on network
							var userAddress, bridgeAddress string
							userAddress = user.DogecoinMainnetPublicKey
							bridgeAddress = bridgeWallet.DogecoinMainnetPublicKey

							// Validate that we have the Dogecoin addresses
							if userAddress == "" || bridgeAddress == "" {
								logger.Sugar().Errorw("Dogecoin addresses not found in wallet")
								break
							}

							txn, dataToSign, err := dogecoin.WithdrawDogeNativeGetSignature(
								withdrawalChain.ChainUrl,
								bridgeAddress,
								amount,
								userAddress,
							)

							if err != nil {
								logger.Sugar().Errorw("error getting data to sign", "error", err)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
							intent.Operations[i].SolverDataToSign = dataToSign

							signature, err := getSignature(intent, i)
							if err != nil {
								logger.Sugar().Errorw("error getting signature", "error", err)
								break
							}

							// Use the same Dogecoin address we used for signing
							result, err := dogecoin.WithdrawDogeTxn(
								withdrawalChain.ChainUrl,
								txn,         // Use the serialized transaction instead of dataToSign
								userAddress, // Use Dogecoin address instead of ECDSA key
								signature,
							)

							if err != nil {
								logger.Sugar().Errorw("error withdrawing dogecoin", "error", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
							break
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
								logger.Sugar().Errorw("error getting data to sign", "error", err)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
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
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
						}
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

						UpdateOperationSolverDataToSign(operation.ID, dataToSign)
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
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
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

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
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
								logger.Sugar().Errorw("error getting data to sign", "error", err)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
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
							logger.Sugar().Errorw("error getting wallet", "error", err)
							UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
							UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
							break
						}

						if wallet.StellarPublicKey == "" {
							logger.Sugar().Errorw("error: no Stellar public key found in wallet")
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
								logger.Sugar().Errorw("error withdrawing native XLM", "error", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
							intent.Operations[i].SolverDataToSign = dataToSign

							signature, err := getSignature(intent, i)
							if err != nil {
								logger.Sugar().Errorw("error getting signature", "error", err)
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
								logger.Sugar().Errorw("error withdrawing Stellar", "error", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
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
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
							intent.Operations[i].SolverDataToSign = dataToSign

							signature, err := getSignature(intent, i)
							if err != nil {
								logger.Sugar().Errorw("error getting signature", "error", err)
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
								logger.Sugar().Errorw("error withdrawing Stellar", "error", err)
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

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
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
								logger.Sugar().Errorw("error getting data to sign", "error", err)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
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
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
						}
						break
					} else if withdrawalChain.KeyCurve == "sui_eddsa" {
						wallet, err := GetWallet(intent.Identity, intent.IdentityCurve)
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

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
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
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
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

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
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
								logger.Sugar().Errorw("error getting data to sign", "error", err)
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
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
								logger.Sugar().Errorw("error getting data to sign", "error", err)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
							intent.Operations[i].SolverDataToSign = dataToSign

							signature, err := getSignature(intent, i)
							if err != nil {
								logger.Sugar().Errorw("error getting signature", "error", err)
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
								logger.Sugar().Errorw("error withdrawing algorand", "error", err)
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
								logger.Sugar().Errorw("error getting data to sign", "error", err)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
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
								logger.Sugar().Errorw("error getting data to sign", "error", err)
								break
							}

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
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
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
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

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
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
								UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
								UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
								break
							}

							UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, result)
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

							UpdateOperationSolverDataToSign(operation.ID, dataToSign)
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
					if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "bitcoin_ecdsa" {
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
							confirmed, err = checkEVMTransactionConfirmed(operation.ChainId, operation.Result)
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
							confirmed, err = checkSolanaTransactionConfirmed(operation.ChainId, operation.Result)
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
						logger.Sugar().Errorw("error getting chain", "error", err)
						break
					}

					switch chain.ChainType {
					case "bitcoin":
						confirmed, err = bitcoin.CheckBitcoinTransactionConfirmed(operation.ChainId, operation.Result)
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
						logger.Sugar().Errorw("error checking transaction", "chain", chain.ChainType, "error", err)
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
					if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "bitcoin_ecdsa" {
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
							confirmed, err = checkEVMTransactionConfirmed(operation.ChainId, operation.Result)
							if err != nil {
								logger.Sugar().Errorw("error checking evm transaction", "error", err)
								break
							}
						}
					} else if operation.KeyCurve == "eddsa" ||
						operation.KeyCurve == "aptos_eddsa" ||
						operation.KeyCurve == "stellar_eddsa" ||
						operation.KeyCurve == "sui_eddsa" ||
						operation.KeyCurve == "algorand_eddsa" ||
						operation.KeyCurve == "ripple_eddsa" ||
						operation.KeyCurve == "cardano_eddsa" {
						chain, err := common.GetChain(operation.ChainId)
						if err != nil {
							logger.Sugar().Errorw("error getting chain", "error", err)
							break
						}

						if chain.ChainType == "solana" {
							confirmed, err = checkSolanaTransactionConfirmed(operation.ChainId, operation.Result)
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
						logger.Sugar().Errorw("error checking solver status", "error", err)
						break
					}

					if status == solver.SOLVER_OPERATION_STATUS_SUCCESS {
						output, err := solver.GetOutput(operation.Solver, &intentBytes, i)

						if err != nil {
							logger.Sugar().Errorw("error getting solver output", "error", err)
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
					if operation.KeyCurve == "ecdsa" || operation.KeyCurve == "bitcoin_ecdsa" {
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
							confirmed, err = checkEVMTransactionConfirmed(operation.ChainId, operation.Result)
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
							confirmed, err = checkSolanaTransactionConfirmed(operation.ChainId, operation.Result)
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

					lockSchema, err := GetLock(intent.Identity, intent.IdentityCurve)
					if err != nil {
						logger.Sugar().Errorw("error getting lock", "error", err)
						break
					}

					if withdrawMetadata.Unlock {
						depositOperation := intent.Operations[i-4]
						// check for confirmations
						confirmed = false
						if depositOperation.KeyCurve == "ecdsa" || depositOperation.KeyCurve == "bitcoin_ecdsa" {
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
									err := UnlockIdentity(lockSchema.Id)
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
									err := UnlockIdentity(lockSchema.Id)
									if err != nil {
										logger.Sugar().Errorw("error unlocking identity", "error", err)
										break
									}
								}
							} else {
								txnConfirmed, err := checkEVMTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									logger.Sugar().Errorw("error checking evm transaction", "error", err)
									break
								}

								if txnConfirmed {
									confirmed = true
									err := UnlockIdentity(lockSchema.Id)
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
								txnConfirmed, err := checkSolanaTransactionConfirmed(depositOperation.ChainId, depositOperation.Result)
								if err != nil {
									logger.Sugar().Errorw("error checking solana transaction", "error", err)
									break
								}

								if txnConfirmed {
									confirmed = true
									err := UnlockIdentity(lockSchema.Id)
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
									err := UnlockIdentity(lockSchema.Id)
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
									err := UnlockIdentity(lockSchema.Id)
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
									err := UnlockIdentity(lockSchema.Id)
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
									err := UnlockIdentity(lockSchema.Id)
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
									err := UnlockIdentity(lockSchema.Id)
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
									err := UnlockIdentity(lockSchema.Id)
									if err != nil {
										logger.Sugar().Errorw("error unlocking identity", "error", err)
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
	signature, _, err := getSignatureEx(intent, operationIndex)
	if err != nil {
		return "", err
	}
	return signature, nil
}

func getSignatureEx(intent *Intent, operationIndex int) (string, string, error) {
	// get wallet
	wallet, err := GetWallet(intent.Identity, intent.IdentityCurve)
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

	req, err := http.NewRequest("POST", signer.URL+"/signature?operationIndex="+operationIndexStr, bytes.NewBuffer(intentBytes))

	if err != nil {
		return "", "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", "", fmt.Errorf("error sending request: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response body: %v", err)
	}

	var signatureResponse SignatureResponse
	err = json.Unmarshal(body, &signatureResponse)
	if err != nil {
		return "", "", fmt.Errorf("error unmarshalling response body: %v", err)
	}

	return signatureResponse.Signature, signatureResponse.Address, nil
}

func checkEVMTransactionConfirmed(chainId string, txnHash string) (bool, error) {
	chain, err := common.GetChain(chainId)
	if err != nil {
		return false, err
	}

	client, err := ethclient.Dial(chain.ChainUrl)
	if err != nil {
		return false, fmt.Errorf("failed to dial EVM client: %v", err)
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
		return "", fmt.Errorf("failed to decode transaction data: %v", err)
	}

	// Deserialize the binary data into a Solana transaction
	// This reconstructs the transaction object with all its instructions
	_tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(decodedTransactionData))
	if err != nil {
		return "", fmt.Errorf("failed to deserialize transaction data: %v", err)
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
		return "", fmt.Errorf("failed to verify signatures: %v", err)
	}

	// Submit the transaction to the Solana network
	// The returned hash can be used to track the transaction status
	hash, err := c.SendTransaction(context.Background(), _tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	// Return the transaction hash as a string
	return hash.String(), nil
}

// readBitcoinAddress returns the appropriate Bitcoin public key based on the chain configuration
func readBitcoinAddress(wallet *WalletSchema, chainId string) (string, error) {
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
