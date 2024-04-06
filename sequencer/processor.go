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

	"github.com/Silent-Protocol/go-sio/solver"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
)

func ProcessIntent(intentId int64) {
	for {
		intent, err := GetIntent(intentId)
		if err != nil {
			log.Println(err)
			return
		}

		if intent.Status != INTENT_STATUS_PROCESSING {
			log.Println("intent processed")
			return
		}

		if err != nil {
			log.Println(err)
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
					// get signature
					// then pass the solver info to solver moduler with signature to handle the rest
					signature, err := getSignature(intent, i)
					if err != nil {
						fmt.Println(err)
						break
					}

					// pass the solver info to solver moduler with signature to handle the rest
					id, err := solver.Solve(
						operation.Solver,
						operation.SolverMetadata,
						operation.DataToSign,
						signature,
					)

					if err != nil {
						fmt.Println(err)
						UpdateOperationStatus(operation.ID, OPERATION_STATUS_FAILED)
						UpdateIntentStatus(intent.ID, INTENT_STATUS_FAILED)
						break
					}

					UpdateOperationResult(operation.ID, OPERATION_STATUS_WAITING, id)
				}

				break
			}

			if operation.Status == OPERATION_STATUS_WAITING {
				// check for confirmations and update the status to completed
				confirmed := false
				if operation.Type == OPERATION_TYPE_TRANSACTION {
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
				} else if operation.Type == OPERATION_TYPE_SOLVER {
					status, err := solver.CheckStatus(
						operation.Solver,
						operation.SolverMetadata,
						operation.Result,
					)

					if err != nil {
						fmt.Println(err)
						break
					}

					if status == "completed" {
						confirmed = true
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
