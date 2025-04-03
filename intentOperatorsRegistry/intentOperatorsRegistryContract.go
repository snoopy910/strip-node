package intentoperatorsregistry

import (
	"fmt"

	"github.com/StripChain/strip-node/util/logger"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetIntentOperatorsRegistryContract(
	rpcURL string,
	intentOperatorsRegistryContractAddress string,
) (*IntentOperatorsRegistry, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		logger.Sugar().Errorw("failed to dial ethclient", "error", err)
		return nil, fmt.Errorf("failed to dial ethclient: %w", err)
	}

	instance, err := NewIntentOperatorsRegistry(ethCommon.HexToAddress(intentOperatorsRegistryContractAddress), client)

	if err != nil {
		logger.Sugar().Errorw("failed to create intent operators registry contract", "error", err)
		return nil, fmt.Errorf("failed to create intent operators registry contract: %w", err)
	}

	return instance, nil
}
