package bridge

import (
	"github.com/StripChain/strip-node/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TokenExists(rpcURL string, bridgeContractAddress string, chainId string, srcToken string) (bool, string, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return false, "", err
	}

	instance, err := NewBridge(common.HexToAddress(bridgeContractAddress), client)
	if err != nil {
		return false, "", err
	}

	peggedToken, err := instance.PeggedTokens(&bind.CallOpts{}, chainId, srcToken)

	if err != nil {
		return false, "", err
	}

	if peggedToken != common.HexToAddress(util.ZERO_ADDRESS) {
		return true, peggedToken.Hex(), nil
	}

	return false, "", nil
}
