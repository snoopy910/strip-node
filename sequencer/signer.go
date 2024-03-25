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

	// prefix public key with 0x
	for i, signer := range signers {
		signers[i].PublicKey = "0x" + signer.PublicKey
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
