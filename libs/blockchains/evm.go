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

	newEVMBlockchain, err := NewEVMBlockchain(Ethereum, network, "hex", 18, time.Second*10, &chainID, "ETH")
	if err != nil {
		logger.Sugar().Errorw("failed to create new EVMBlockchain", "error", err)
		return nil, err
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

	newEVMBlockchain, err := NewEVMBlockchain(StripChain, network, "hex", 18, time.Second*10, &chainID, "ETH")
	if err != nil {
		logger.Sugar().Errorw("failed to create new EVMBlockchain", "error", err)
		return nil, fmt.Errorf("failed to create new EVMBlockchain: %w", err)
	}
	return &newEVMBlockchain, nil
}
