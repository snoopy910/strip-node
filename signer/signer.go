package signer

import (
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/multiformats/go-multiaddr"
)

var RPC_URL, SignerHubContractAddress, NodePrivateKey, NodePublicKey string
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
	signerHubContractAddress string,
	maximumSigners int,
) {
	RPC_URL = rpcURL
	SignerHubContractAddress = signerHubContractAddress
	NodePrivateKey = signerPrivateKey
	NodePublicKey = signerPublicKey

	// this should come from contract
	MaximumSigners = maximumSigners

	go startHTTPServer(httpPort)

	h, addr := createHost(listenHost, port, bootnodeURL)
	go discoverPeers(h, []multiaddr.Multiaddr{addr})
	subscribe(h)
}
