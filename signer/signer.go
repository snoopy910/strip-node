package signer

import (
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/multiformats/go-multiaddr"
)

var RPC_URL, IntentOperatorsRegistryContractAddress, NodePrivateKey, NodePublicKey string
var MaximumSigners int

type PartyProcess struct {
	Party  *tss.Party
	Exists bool
}

var partyProcesses = make(map[string]PartyProcess)

func Start(
	signerPrivateKey string,
	signerPublicKey string,
	bootnodeURL string,
	httpPort string,
	listenHost string,
	port int,
	rpcURL string,
	intentOperatorsRegistryContractAddress string,
	maximumSigners int,
) {
	RPC_URL = rpcURL
	IntentOperatorsRegistryContractAddress = intentOperatorsRegistryContractAddress
	NodePrivateKey = signerPrivateKey
	NodePublicKey = signerPublicKey

	instance := getIntentOperatorsRegistryContract(RPC_URL, IntentOperatorsRegistryContractAddress)
	_maxSigners, err := instance.MAXIMUMSIGNERS(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}

	MaximumSigners = int(_maxSigners.Int64())

	go startHTTPServer(httpPort)

	h, addr := createHost(listenHost, port, bootnodeURL)
	go discoverPeers(h, []multiaddr.Multiaddr{addr})
	subscribe(h)
}
