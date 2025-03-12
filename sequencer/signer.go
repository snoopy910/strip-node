package sequencer

import (
	"context"
	"encoding/hex"
	"fmt"

	intentoperatorsregistry "github.com/StripChain/strip-node/intentOperatorsRegistry"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Signer struct {
	PublicKey string
	URL       string
}

var SignersList = func() ([]Signer, error) {
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ethclient: %w", err)
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
		return nil, fmt.Errorf("failed to filter logs: %w", err)
	}

	signers := []Signer{}

	for _, log := range logs {
		data, err := instance.ParseSignerUpdated(log)
		if err != nil {
			return nil, fmt.Errorf("failed to parse signer updated: %w", err)
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

	return signers, nil
}

func GetSigner(publicKey string) (Signer, error) {
	signers, err := SignersList()
	if err != nil {
		return Signer{}, fmt.Errorf("failed to get signers: %w", err)
	}

	for _, signer := range signers {
		if signer.PublicKey == publicKey {
			return signer, nil
		}
	}
	return Signer{}, fmt.Errorf("signer not found")
}
