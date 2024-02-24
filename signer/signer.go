package signer

import (
	"log"

	"github.com/Silent-Protocol/go-sio/common"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/multiformats/go-multiaddr"
)

// type Network struct {
// 	RPC_URL                  string
// 	SignerHubContractAddress string
// 	Key                      *keygen.LocalPartySaveData
// }

var RPC_URL, SignerHubContractAddress, NodePrivateKey, NodePublicKey string
var Index, Threshold, TotalSigners, StartKey int

type PartyProcess struct {
	Party  *tss.Party
	Exists bool
}

// identity => operationType => PartyData
var partyProcesses = make(map[string]map[string]PartyProcess)

func Start(
	signerPrivateKey string,
	signerPublicKey string,
	bootnodeURL string,
	httpPort string,
	listenHost string,
	port int,
	rpcURL string,
	signerHubContractAddress string,
) {
	RPC_URL = rpcURL
	SignerHubContractAddress = signerHubContractAddress
	NodePrivateKey = signerPrivateKey
	NodePublicKey = signerPublicKey

	instance := getSignerHubContract(RPC_URL, SignerHubContractAddress)

	t, err := instance.CurrentThreshold(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	Threshold = int(t.Int64())

	ts, err := instance.NextIndex(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	TotalSigners = int(ts.Int64())

	_i, err := instance.Signers(&bind.CallOpts{}, common.PublicKeyStrToBytes32(NodePublicKey))
	if err != nil {
		log.Fatal(err)
	}
	Index = int(_i.Index.Int64())

	startKey, err := instance.StartKey(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	StartKey = int(startKey.Int64())

	go startHTTPServer(httpPort)

	h, addr := createHost(listenHost, port, bootnodeURL)
	go discoverPeers(h, []multiaddr.Multiaddr{addr})
	subscribe(h)
}
