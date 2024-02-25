package signer

import (
	"log"

	"github.com/Silent-Protocol/go-sio/common"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/multiformats/go-multiaddr"
)

var RPC_URL, SignerHubContractAddress, NodePrivateKey, NodePublicKey string
var Index, Threshold, TotalSigners int

type PartyProcess struct {
	Party  *tss.Party
	Exists bool
}

// identity => operationType => PartyData
var partyProcesses = make(map[string]PartyProcess)

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

	go startHTTPServer(httpPort)

	h, addr := createHost(listenHost, port, bootnodeURL)
	go discoverPeers(h, []multiaddr.Multiaddr{addr})
	subscribe(h)
}
