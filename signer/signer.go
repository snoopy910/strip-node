package signer

import (
	intentoperatorsregistry "github.com/StripChain/strip-node/intentOperatorsRegistry"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/multiformats/go-multiaddr"
)

var RPC_URL, IntentOperatorsRegistryContractAddress, SolversRegistryContractAddress, BridgeContractAddress, NodePrivateKey, NodePublicKey string
var MaximumSigners int
var HeliusApiKey string

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
	solversRegistryContractAddress string,
	maximumSigners int,
	heliusApiKey string,
	bridgeContractAddress string,
) {
	HeliusApiKey = heliusApiKey
	RPC_URL = rpcURL
	IntentOperatorsRegistryContractAddress = intentOperatorsRegistryContractAddress
	SolversRegistryContractAddress = solversRegistryContractAddress
	NodePrivateKey = signerPrivateKey
	BridgeContractAddress = bridgeContractAddress

	NodePublicKey = signerPublicKey

	instance := intentoperatorsregistry.GetIntentOperatorsRegistryContract(RPC_URL, IntentOperatorsRegistryContractAddress)
	_maxSigners, err := instance.MAXIMUMSIGNERS(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}

	MaximumSigners = int(_maxSigners.Int64())

	go startHTTPServer(httpPort)

	h, addr, err := createHost(listenHost, port, bootnodeURL)
	if err != nil {
		panic(err)
	}

	go discoverPeers(h, []multiaddr.Multiaddr{addr})
	err = subscribe(h)
	if err != nil {
		panic(err)
	}
}
