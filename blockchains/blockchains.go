package blockchains

import (
	"fmt"

	"github.com/StripChain/strip-node/common"
)

type Blockchain interface {
	ChainName() BlockchainID
	KeyCurve() common.Curve
	SigningEncoding() string
	// Replacing for Send*Transaction
	BroadcastTransaction(txn string, signedHash string, publicKey *string) (string, error)
	GetTransfers(txHash string) ([]common.Transfer, error)
	// Replacing for Check*TransactionConfirmed
	IsTransactionBroadcastedAndConfirmed(txHash string) (bool, error)
	// Replacing for Withdraw*GetSignature
	BuildWithdrawTx(bridgeAddress string,
		solverOutput string,
		userAddress string,
		tokenAddress *string,
	) (string, string, error)

	RawPublicKeyBytesToAddress(pkBytes []byte) (string, error)
	RawPublicKeyToPublicKeyStr(pkBytes []byte) (string, error)
}

type Network struct {
	networkID   string
	networkType NetworkType
	nodeURL     string
	apiKey      *string
}

func (n *Network) GetNetworkID() string {
	return n.networkID
}

func (n *Network) GetNetworkType() NetworkType {
	return n.networkType
}

func (n *Network) GetNodeURL() string {
	return n.nodeURL
}

func (n *Network) GetAPIKey() *string {
	return n.apiKey
}

type NetworkType string

const (
	Mainnet NetworkType = "mainnet"
	Testnet NetworkType = "testnet"
	Devnet  NetworkType = "devnet"
)

type BlockchainID string

const (
	Ethereum BlockchainID = "ethereum"
	Dogecoin BlockchainID = "dogecoin"
	Stellar  BlockchainID = "stellar"
	Cardano  BlockchainID = "cardano"
	// ... and so on
)

// BlockchainFactory creates blockchain instances for specific networks
type BlockchainFactory func(networkType NetworkType) (Blockchain, error)

// Registry manages blockchain implementations
type Registry struct {
	factories map[BlockchainID]BlockchainFactory
	instances map[BlockchainID]map[NetworkType]Blockchain // chainID -> networkType -> instance
}

// NetworkConfig contains configuration for a specific network
type NetworkConfig struct {
	NetworkName string
	NetworkType NetworkType
	NodeURL     string
	APIKey      *string
	// Optional
	// ExtraParams map[string]interface{}
}

// NewRegistry creates a new blockchain registry
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[BlockchainID]BlockchainFactory),
		instances: make(map[BlockchainID]map[NetworkType]Blockchain),
	}
}

// Register adds a blockchain factory to the registry
func (r *Registry) Register(chainID BlockchainID, factory BlockchainFactory) {
	r.factories[chainID] = factory
}

// GetBlockchain returns a blockchain instance for the given chain and network
func (r *Registry) GetBlockchain(chainID BlockchainID, networkType NetworkType) (Blockchain, error) {
	// Check if instance already exists
	if networks, ok := r.instances[chainID]; ok {
		if blockchain, ok := networks[networkType]; ok {
			return blockchain, nil
		}
	}

	// Instance doesn't exist, try to create it
	factory, ok := r.factories[chainID]
	if !ok {
		return nil, fmt.Errorf("no blockchain factory registered for chain: %s", chainID)
	}

	// Create instance
	blockchain, err := factory(networkType)
	if err != nil {
		return nil, err
	}

	// Store instance
	if _, ok := r.instances[chainID]; !ok {
		r.instances[chainID] = make(map[NetworkType]Blockchain)
	}
	r.instances[chainID][networkType] = blockchain

	return blockchain, nil
}

// BaseBlockchain provides common functionality for blockchain implementations
type BaseBlockchain struct {
	chainName       BlockchainID
	network         Network
	keyCurve        common.Curve
	signingEncoding string
}

func (b *BaseBlockchain) ChainName() BlockchainID {
	return b.chainName
}

func (b *BaseBlockchain) Network() Network {
	return b.network
}

func (b *BaseBlockchain) KeyCurve() common.Curve {
	return b.keyCurve
}

func (b *BaseBlockchain) SigningEncoding() string {
	return b.signingEncoding
}

// NewCardanoBlockchain creates a new Cardano blockchain instance
func NewCardanoBlockchain(networkType NetworkType) (Blockchain, error) {
	mainnetAPIKey := "whatever-key"
	testnetAPIKey := "whatever-key-testnet"
	network := Network{
		networkType: networkType,
		nodeURL:     "https://cardano-mainnet.blockfrost.io/api/v0",
		apiKey:      &mainnetAPIKey,
		networkID:   "mainnet",
	}

	if networkType == Testnet {
		network.nodeURL = "https://cardano-preprod.blockfrost.io/api/v0"
		network.apiKey = &testnetAPIKey
		network.networkID = "preprod"
	}

	return &CardanoBlockchain{
		BaseBlockchain: BaseBlockchain{
			chainName:       Cardano,
			network:         network,
			keyCurve:        common.CurveCardano,
			signingEncoding: "hex",
		},
	}, nil
}

// This is a type assertion to ensure that the CardanoBlockchain implements the Blockchain interface
var _ Blockchain = &CardanoBlockchain{}

// CardanoBlockchain implements the Blockchain interface for Cardano
type CardanoBlockchain struct {
	BaseBlockchain
}

func (b *CardanoBlockchain) BroadcastTransaction(txn string, signedHash string, publicKey *string) (string, error) {
	return "", nil
}

func (b *CardanoBlockchain) GetTransfers(txHash string) ([]common.Transfer, error) {
	return nil, nil
}

func (b *CardanoBlockchain) IsTransactionBroadcastedAndConfirmed(txHash string) (bool, error) {
	return false, nil
}

func (b *CardanoBlockchain) BuildWithdrawTx(bridgeAddress string,
	solverOutput string,
	userAddress string,
	tokenAddress *string,
) (string, string, error) {
	return "", "", nil
}

func (b *CardanoBlockchain) RawPublicKeyBytesToAddress(pkBytes []byte) (string, error) {
	return "", nil
}

func (b *CardanoBlockchain) RawPublicKeyToPublicKeyStr(pkBytes []byte) (string, error) {
	return "", nil
}

// InitBlockchainRegistry sets up the blockchain registry with all supported chains
func InitBlockchainRegistry() *Registry {
	registry := NewRegistry()

	registry.Register(Cardano, NewCardanoBlockchain)

	return registry
}
