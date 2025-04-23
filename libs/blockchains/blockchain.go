package blockchains

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/StripChain/strip-node/common"
)

type IBlockchain interface {
	ChainName() BlockchainID
	KeyCurve() common.Curve
	Decimals() uint
	SigningEncoding() string
	OpTimeout() time.Duration
	ChainID() *string
	TokenSymbol() string
	// Replacing for Send*Transaction
	BroadcastTransaction(txn string, signedHash string, publicKey *string) (string, error)
	GetTransfers(txHash string, address *string) ([]common.Transfer, error)
	// Replacing for Check*TransactionConfirmed
	IsTransactionBroadcastedAndConfirmed(txHash string) (bool, error)
	// Replacing for Withdraw*GetSignature
	BuildWithdrawTx(bridgeAddress string,
		solverOutput string,
		userAddress string,
		tokenAddress *string,
	) (string, string, error)

	RawPublicKeyBytesToAddress(pkBytes []byte, networkType NetworkType) (string, error)
	RawPublicKeyToPublicKeyStr(pkBytes []byte) (string, error)
	ExtractDestinationAddress(serializedTxn string) (string, string, error)
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
	Mainnet NetworkType = "MAINNET"
	Testnet NetworkType = "TESTNET"
	Devnet  NetworkType = "DEVNET"
	Regnet  NetworkType = "REGNET"
)

type BlockchainID string

const (
	Ethereum   BlockchainID = "ETHEREUM"
	Dogecoin   BlockchainID = "DOGECOIN"
	Stellar    BlockchainID = "STELLAR"
	Cardano    BlockchainID = "CARDANO"
	Algorand   BlockchainID = "ALGORAND"
	Ripple     BlockchainID = "RIPPLE"
	Sui        BlockchainID = "SUI"
	Solana     BlockchainID = "SOLANA"
	Bitcoin    BlockchainID = "BITCOIN"
	Aptos      BlockchainID = "APTOS"
	StripChain BlockchainID = "STRIPCHAIN"
	Arbitrum   BlockchainID = "ARBITRUM"
)

// BlockchainFactory creates blockchain instances for specific networks
type BlockchainFactory func(networkType NetworkType) (IBlockchain, error)

// Registry manages blockchain implementations
type Registry struct {
	factories map[BlockchainID]BlockchainFactory
	instances map[BlockchainID]map[NetworkType]IBlockchain // chainID -> networkType -> instance
}

// singleton instance
var (
	registryInstance *Registry
	registryOnce     sync.Once
)

// GetRegistry returns the singleton instance of the registry
func GetRegistry() *Registry {
	registryOnce.Do(func() {
		registryInstance = newRegistry()
	})
	return registryInstance
}

// newRegistry creates a new blockchain registry (private constructor)
func newRegistry() *Registry {
	return &Registry{
		factories: make(map[BlockchainID]BlockchainFactory),
		instances: make(map[BlockchainID]map[NetworkType]IBlockchain),
	}
}

// Register adds a blockchain factory to the registry
func (r *Registry) Register(chainID BlockchainID, factory BlockchainFactory) {
	r.factories[chainID] = factory
}

// GetBlockchain returns a blockchain instance for the given chain and network
func (r *Registry) GetBlockchain(blockchainID BlockchainID, networkType NetworkType) (IBlockchain, error) {
	// Check if instance already exists
	if networks, ok := r.instances[blockchainID]; ok {
		if blockchain, ok := networks[networkType]; ok {
			return blockchain, nil
		}
	}

	// Instance doesn't exist, try to create it
	factory, ok := r.factories[blockchainID]
	if !ok {
		return nil, fmt.Errorf("no blockchain factory registered for blockchain id: %s", blockchainID)
	}

	// Create instance
	blockchain, err := factory(networkType)
	if err != nil {
		return nil, err
	}

	// Store instance
	if _, ok := r.instances[blockchainID]; !ok {
		r.instances[blockchainID] = make(map[NetworkType]IBlockchain)
	}
	r.instances[blockchainID][networkType] = blockchain

	return blockchain, nil
}

func (r *Registry) GetRegisteredBlockchains() []BlockchainID {
	blockchains := make([]BlockchainID, 0, len(r.factories))
	for blockchainID := range r.factories {
		blockchains = append(blockchains, blockchainID)
	}
	return blockchains
}

// BaseBlockchain provides common functionality for blockchain implementations
type BaseBlockchain struct {
	chainName       BlockchainID
	network         Network
	keyCurve        common.Curve
	signingEncoding string
	decimals        uint
	opTimeout       time.Duration
	chainID         *string
	tokenSymbol     string
}

func (b *BaseBlockchain) ChainName() BlockchainID {
	return b.chainName
}

func (b *BaseBlockchain) KeyCurve() common.Curve {
	return b.keyCurve
}

func (b *BaseBlockchain) SigningEncoding() string {
	return b.signingEncoding
}

func (b *BaseBlockchain) Decimals() uint {
	return b.decimals
}

func (b *BaseBlockchain) OpTimeout() time.Duration {
	return b.opTimeout
}

func (b *BaseBlockchain) ChainID() *string {
	return b.chainID
}

func (b *BaseBlockchain) TokenSymbol() string {
	return b.tokenSymbol
}

// Default implementations - can be overridden by specific blockchains
func (b *BaseBlockchain) BroadcastTransaction(txn string, signedHash string, publicKey *string) (string, error) {
	return "", errors.New("BroadcastTransaction not implemented")
}

func (b *BaseBlockchain) GetTransfers(txHash string, address *string) ([]common.Transfer, error) {
	return nil, errors.New("GetTransfers not implemented")
}

func (b *BaseBlockchain) IsTransactionBroadcastedAndConfirmed(txHash string) (bool, error) {
	return false, errors.New("IsTransactionBroadcastedAndConfirmed not implemented")
}

func (b *BaseBlockchain) BuildWithdrawTx(bridgeAddress string,
	solverOutput string,
	userAddress string,
	tokenAddress *string,
) (string, string, error) {
	return "", "", errors.New("BuildWithdrawTx not implemented")
}

func (b *BaseBlockchain) RawPublicKeyBytesToAddress(pkBytes []byte, networkType NetworkType) (string, error) {
	return "", errors.New("RawPublicKeyBytesToAddress not implemented")
}

func (b *BaseBlockchain) RawPublicKeyToPublicKeyStr(pkBytes []byte) (string, error) {
	return "", errors.New("RawPublicKeyToPublicKeyStr not implemented")
}

func (b *BaseBlockchain) ExtractDestinationAddress(serializedTxn string) (string, string, error) {
	return "", "", errors.New("ExtractDestinationAddress not implemented")
}

// TODO: This needs improvement
func ParseBlockchainID(blockchainID string) (BlockchainID, error) {
	return BlockchainID(blockchainID), nil
}

func ParseNetworkType(networkType string) (NetworkType, error) {
	return NetworkType(networkType), nil
}

// InitBlockchainRegistry sets up the blockchain registry with all supported chains
func InitBlockchainRegistry() *Registry {
	registry := GetRegistry()

	chains := map[BlockchainID]BlockchainFactory{
		Algorand:   NewAlgorandBlockchain,
		Arbitrum:   NewArbitrumBlockchain,
		Cardano:    NewCardanoBlockchain,
		Ethereum:   NewEthereumBlockchain,
		Stellar:    NewStellarBlockchain,
		Ripple:     NewRippleBlockchain,
		Sui:        NewSuiBlockchain,
		Solana:     NewSolanaBlockchain,
		Bitcoin:    NewBitcoinBlockchain,
		Aptos:      NewAptosBlockchain,
		Dogecoin:   NewDogecoinBlockchain,
		StripChain: NewStripChainBlockchain,
	}
	for chain, factory := range chains {
		registry.Register(chain, factory)
	}

	return registry
}

func GetBlockchain(blockchainID BlockchainID, networkType NetworkType) (IBlockchain, error) {
	registry := GetRegistry()
	return registry.GetBlockchain(blockchainID, networkType)
}

func GetRegisteredBlockchains() []BlockchainID {
	registry := GetRegistry()
	return registry.GetRegisteredBlockchains()
}

var _ IBlockchain = &BaseBlockchain{}

// Verify BaseBlockchain implements IBlockchain
var _ IBlockchain = (*BaseBlockchain)(nil)
