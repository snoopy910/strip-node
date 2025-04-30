package blockchains

import (
	"fmt"
	"time"

	"github.com/StripChain/strip-node/util/logger"
)

// NewEthereumBlockchain creates a new Ethereum blockchain instance
func NewEthereumBlockchain(networkType NetworkType) (IBlockchain, error) {
	var nodeURL, networkID, chainID string
	switch networkType {
	case Mainnet:
		nodeURL = "https://ethereum-rpc.publicnode.com"
		networkID = "mainnet"
		chainID = "1"
	case Testnet:
		nodeURL = "https://ethereum-sepolia-rpc.publicnode.com"
		networkID = "testnet"
		chainID = "11155111"
	case Regnet:
		nodeURL = "http://ganache:8545"
		networkID = "regnet"
		chainID = "1337"
	}
	network := Network{
		networkType: networkType,
		nodeURL:     nodeURL,
		networkID:   networkID,
	}

	newEVMBlockchain, err := NewEVMBlockchain(Ethereum, network, "hex", 18, time.Minute*2, &chainID, "ETH")
	if err != nil {
		logger.Sugar().Errorw("failed to create new EVMBlockchain", "error", err)
		return nil, err
	}
	return &newEVMBlockchain, nil
}

// NewArbitrumBlockchain creates a new Arbitrum blockchain instance
func NewArbitrumBlockchain(networkType NetworkType) (IBlockchain, error) {
	var nodeURL, networkID, chainID string
	switch networkType {
	case Mainnet:
		nodeURL = "https://arbitrum-one-rpc.publicnode.com"
		networkID = "mainnet"
		chainID = "42161"
	case Testnet:
		nodeURL = "https://arbitrum-sepolia-rpc.publicnode.com"
		networkID = "testnet"
		chainID = "421614"
	default:
		return nil, fmt.Errorf("network type not supported: %s", networkType)
	}
	network := Network{
		networkType: networkType,
		nodeURL:     nodeURL,
		networkID:   networkID,
	}

	newEVMBlockchain, err := NewEVMBlockchain(Arbitrum, network, "hex", 18, time.Minute*2, &chainID, "ETH")
	if err != nil {
		logger.Sugar().Errorw("failed to create new EVMBlockchain", "error", err)
		return nil, fmt.Errorf("failed to create new EVMBlockchain: %w", err)
	}
	return &newEVMBlockchain, nil
}

// NewStripChainBlockchain creates a new Ethereum blockchain instance
func NewStripChainBlockchain(networkType NetworkType) (IBlockchain, error) {
	var nodeURL, networkID, chainID string
	switch networkType {
	case Testnet:
		nodeURL = "https://rpc-stripsepolia-5w8r5b9f7b.t.conduit.xyz"
		networkID = "testnet"
		chainID = "44331"
	default:
		return nil, fmt.Errorf("network type not supported: %s", networkType)
	}
	network := Network{
		networkType: networkType,
		nodeURL:     nodeURL,
		networkID:   networkID,
	}

	newEVMBlockchain, err := NewEVMBlockchain(StripChain, network, "hex", 18, time.Minute*2, &chainID, "ETH")
	if err != nil {
		logger.Sugar().Errorw("failed to create new EVMBlockchain", "error", err)
		return nil, fmt.Errorf("failed to create new EVMBlockchain: %w", err)
	}
	return &newEVMBlockchain, nil
}

func NewSonicBlockchain(networkType NetworkType) (IBlockchain, error) {
	var nodeURL, networkID, chainID string
	switch networkType {
	case Mainnet:
		nodeURL = "https://rpc.soniclabs.com"
		networkID = "mainnet"
		chainID = "146"
	case Testnet:
		nodeURL = "https://rpc.blaze.soniclabs.com"
		networkID = "testnet"
		chainID = "57054"
	default:
		return nil, fmt.Errorf("network type not supported: %s", networkType)
	}
	network := Network{
		networkType: networkType,
		nodeURL:     nodeURL,
		networkID:   networkID,
	}

	newEVMBlockchain, err := NewEVMBlockchain(Sonic, network, "hex", 18, time.Minute*2, &chainID, "S")
	if err != nil {
		logger.Sugar().Errorw("failed to create new EVMBlockchain", "error", err)
		return nil, fmt.Errorf("failed to create new EVMBlockchain: %w", err)
	}
	return &newEVMBlockchain, nil
}

func NewBerachainBlockchain(networkType NetworkType) (IBlockchain, error) {
	var nodeURL, networkID, chainID string
	switch networkType {
	case Mainnet:
		nodeURL = "https://berachain-rpc.publicnode.com"
		networkID = "mainnet"
		chainID = "80094"
	case Testnet:
		nodeURL = "https://bepolia.rpc.berachain.com"
		networkID = "testnet"
		chainID = "80069"
	default:
		return nil, fmt.Errorf("network type not supported: %s", networkType)
	}
	network := Network{
		networkType: networkType,
		nodeURL:     nodeURL,
		networkID:   networkID,
	}

	newEVMBlockchain, err := NewEVMBlockchain(Berachain, network, "hex", 18, time.Minute*2, &chainID, "BERA")
	if err != nil {
		logger.Sugar().Errorw("failed to create new EVMBlockchain", "error", err)
		return nil, fmt.Errorf("failed to create new EVMBlockchain: %w", err)
	}
	return &newEVMBlockchain, nil
}
