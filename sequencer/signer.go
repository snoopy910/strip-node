package sequencer

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	intentoperatorsregistry "github.com/Silent-Protocol/go-sio/intentOperatorsRegistry"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Signer struct {
	PublicKey string
	URL       string
}

// TODO: This list will be fetched from SC by the sequencer
// var Signers = []Signer{
// 	{
// 		PublicKey: "0x0226d1556a83c01a9d2b1cce29b32cb520238efc602f86481d2d0b9af8a2fff0cf",
// 		URL:       "http://localhost:8080",
// 	},
// 	{
// 		PublicKey: "0x0354455a1f7f4244ef645ac62baa8bd90af0cc18cdb0eae369766b7b58134edf35",
// 		URL:       "http://localhost:8081",
// 	},
// }

func SignersList() []Signer {
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Fatal(err)
	}

	instance := intentoperatorsregistry.GetIntentOperatorsRegistryContract(RPC_URL, IntentOperatorsRegistryContractAddress)

	eventSignature := []byte("SignerUpdated(bytes32,string,bool)")
	hashEvent := crypto.Keccak256Hash(eventSignature)

	query := ethereum.FilterQuery{
		FromBlock: nil, // Start from the genesis block
		ToBlock:   nil, // Filter until the latest block
		Addresses: []common.Address{
			common.HexToAddress(IntentOperatorsRegistryContractAddress),
		},
		Topics: [][]common.Hash{
			{
				common.HexToHash(hashEvent.Hex()),
			},
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	signers := []Signer{}

	for _, log := range logs {
		data, err := instance.ParseSignerUpdated(log)
		if err != nil {
			panic(err)
		}

		fmt.Println(data.Url)

		if data.Added {
			signers = append(signers, Signer{
				PublicKey: hex.EncodeToString(data.Publickey[:]),
				URL:       data.Url,
			})
		} else {
			for i, signer := range signers {
				if signer.PublicKey == hex.EncodeToString(data.Publickey[:]) {
					signers = append(signers[:i], signers[i+1:]...)
				}
			}
		}
	}

	return signers
}

func GetSigner(publicKey string) (Signer, error) {
	for _, signer := range SignersList() {
		if signer.PublicKey == publicKey {
			return signer, nil
		}
	}
	return Signer{}, fmt.Errorf("signer not found")
}
