package sequencer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
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

				if operation.KeyCurve == "ecdsa" {
					signature, err := getSignature(operation.DataToSign, operation.KeyCurve, intent.Identity, intent.IdentityCurve)
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

					UpdateOperationTxnHash(operation.ID, OPERATION_STATUS_WAITING, txnHash)
				} else if operation.KeyCurve == "eddsa" {
					// @TOTO
				}

				break
			}

			if operation.Status == OPERATION_STATUS_WAITING {
				// check for confirmations and update the status to completed
				confirmed := false
				if operation.KeyCurve == "ecdsa" {
					confirmed, err = checkEVMTransactionConfirmed(operation.ChainId, operation.TxnHash)
				} else if operation.KeyCurve == "eddsa" {
					// @TODO
				}

				if err != nil {
					fmt.Println(err)
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
			}
		}

		time.Sleep(5 * time.Second)
	}
}

type SignatureResponse struct {
	Signature string `json:"signature"`
	Address   string `json:"address"`
}

func getSignature(data string, keyCurve string, identity string, identityCurve string) (string, error) {
	// @TODO: get the signer URL based on the signer who is part of the TSS wallet
	resp, err := http.Get(Signers[0].URL + "/signature?message=" + data + "&identity=" + identity + "&identityCurve=" + identityCurve + "&keyCurve=" + keyCurve)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println(string(body))

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

	// signature, err := hex.DecodeString(signatureHex)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// r := new([32]byte)
	// copy(r[:], signature[0:32])

	// s := new([32]byte)
	// copy(s[:], signature[32:64])

	// v := uint8(signature[64]) + 27 // Adjust v value for Ethereum

	sigData, err := hex.DecodeString(signatureHex)

	_tx, err := tx.WithSignature(types.NewEIP155Signer(nil), []byte(sigData))

	if err != nil {
		return "", err
	}

	err = client.SendTransaction(context.Background(), _tx)
	if err != nil {
		return "", err
	}

	return _tx.Hash().Hex(), nil
}
