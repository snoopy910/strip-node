package signer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Silent-Protocol/go-sio/common"
	"github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/multiformats/go-multiaddr"
)

var ethPrivateKey string
var verifyhash bool

type Network struct {
	PaymasterURL             string
	RPC_URL                  string
	SignerHubContractAddress string
	KeyFilePath              string
	NetworkDataFilePath      string
	StartKeyInt              int
	Key                      *keygen.LocalPartySaveData
	Index                    int
	TotalSigners             int
	Threshold                int
}

var networks map[string]Network = make(map[string]Network)

type PartyProcess struct {
	Party  *tss.Party
	Exists bool
}

var partyProcesses = make(map[string]map[string]PartyProcess)

func Start(
	signerPrivateKey string,
	signerPublicKey string,
	bootnodeURL string,
	path string,
	httpPort string,
	listenHost string,
	port int,
	_ethPrivateKey string,
	_verifyhash bool,

	//specific network config array
	networkIds string,
	rpcURLs string,
	signerHubContractAddresses string,
	paymasterURLs string,
) {
	networkIdsArray := strings.Split(networkIds, ",")
	rpcURLsArray := strings.Split(rpcURLs, ",")
	signerHubContractAddressesArray := strings.Split(signerHubContractAddresses, ",")
	paymasterURLsArray := strings.Split(paymasterURLs, ",")

	for i := 0; i < len(networkIdsArray); i++ {
		_paymasterURL := ""
		if _verifyhash {
			_paymasterURL = paymasterURLsArray[i]
		}
		network := Network{
			PaymasterURL:             _paymasterURL,
			RPC_URL:                  rpcURLsArray[i],
			SignerHubContractAddress: signerHubContractAddressesArray[i],
		}
		networks[networkIdsArray[i]] = network
	}

	ethPrivateKey = _ethPrivateKey
	verifyhash = _verifyhash

	loadKey(path, signerPrivateKey, signerPublicKey)

	for i := 0; i < len(networkIdsArray); i++ {
		instance := getSignerHubContract(networks[networkIdsArray[i]].RPC_URL, networks[networkIdsArray[i]].SignerHubContractAddress)

		t, err := instance.CurrentThreshold(&bind.CallOpts{})
		if err != nil {
			log.Fatal(err)
		}
		threshold := int(t.Int64())

		ts, err := instance.NextIndex(&bind.CallOpts{})
		if err != nil {
			log.Fatal(err)
		}
		totalSigners := int(ts.Int64())

		_i, err := instance.Signers(&bind.CallOpts{}, common.PublicKeyStrToBytes32(nodeKey.PublicKey))
		if err != nil {
			log.Fatal(err)
		}
		index := int(_i.Index.Int64())

		networkId := networkIdsArray[i]
		network := networks[networkId]

		network.Threshold = threshold
		network.TotalSigners = totalSigners
		network.Index = index

		filePath := path + "/" + networkIdsArray[i] + "_key.json"
		filePathNetworkData := path + "/" + networkIdsArray[i] + "_network_data.json"

		network.KeyFilePath = filePath
		network.NetworkDataFilePath = filePathNetworkData

		exists, err := common.FileExists(filePath)
		if err != nil {
			panic(err)
		}

		networks[networkId] = network

		if exists {
			jsonFile, err := os.Open(filePath)
			if err != nil {
				panic(err)
			}

			defer jsonFile.Close()

			byteValue, _ := ioutil.ReadAll(jsonFile)
			json.Unmarshal(byteValue, &network.Key)

			content, err := ioutil.ReadFile(filePathNetworkData)

			if err != nil {
				panic(err)
			}

			networkDataStruct := NetworkData{}
			json.Unmarshal(content, &networkDataStruct)

			network.StartKeyInt = int(networkDataStruct.StartKey)
			network.TotalSigners = int(networkDataStruct.TotalSigners)
			network.Threshold = int(networkDataStruct.Threshold)
			network.Index = int(networkDataStruct.Index)

			networks[networkId] = network
		}

		partyProcesses[networkId] = make(map[string]PartyProcess)
	}

	go startHTTPServer(httpPort)

	h, addr := createHost(listenHost, port, bootnodeURL)
	go discoverPeers(h, []multiaddr.Multiaddr{addr})
	subscribe(h)
}
