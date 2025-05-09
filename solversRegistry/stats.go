package solversregistry

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Stats(rpcURL string, contractAddress string, solverDomain string) (bool, uint, []uint, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return false, 0, nil, err
	}

	instance, err := NewSolversRegistry(common.HexToAddress(contractAddress), client)

	if err != nil {
		return false, 0, nil, err
	}

	data, err := instance.Solvers(nil, solverDomain)

	if err != nil {
		return false, 0, nil, err
	}

	chains, err := instance.GetChains(nil, solverDomain)

	if err != nil {
		return false, 0, nil, err
	}

	_chains := make([]uint, len(chains))

	for i := 0; i < len(chains); i++ {
		_chains[i] = uint(chains[i].Int64())
	}

	return data.Whitelisted, uint(data.LastUpdated.Int64()), _chains, nil
}

func TotalSolvers(rpcURL string, contractAddress string) (uint, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return 0, err
	}

	instance, err := NewSolversRegistry(common.HexToAddress(contractAddress), client)

	if err != nil {
		return 0, err
	}

	total, err := instance.TotalSolvers(nil)

	if err != nil {
		return 0, err
	}

	return uint(total.Int64()), nil
}

func SolverExistsAndWhitelisted(rpcURL string, contractAddress string, solverDomain string) (bool, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return false, fmt.Errorf("failed to dial ethclient: %w", err)
	}

	instance, err := NewSolversRegistry(common.HexToAddress(contractAddress), client)

	if err != nil {
		return false, fmt.Errorf("failed to create new solvers registry: %w", err)
	}

	exists, err := instance.Solvers(nil, solverDomain)

	if err != nil {
		return false, fmt.Errorf("failed to get solver: %w", err)
	}

	return exists.Whitelisted, nil
}

func ValidateChain(rpcURL string, contractAddress string, solverDomain string, chainID string) (bool, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return false, fmt.Errorf("failed to dial ethclient: %w", err)
	}

	instance, err := NewSolversRegistry(common.HexToAddress(contractAddress), client)

	if err != nil {
		return false, fmt.Errorf("failed to create new solvers registry: %w", err)
	}

	chains, err := instance.GetChains(nil, solverDomain)

	if err != nil {
		return false, fmt.Errorf("failed to get chains: %w", err)
	}

	chainIDInt, err := strconv.ParseInt(chainID, 10, 64)
	if err != nil {
		return false, fmt.Errorf("failed to parse chainID: %w", err)
	}

	for _, chain := range chains {
		if chain.Int64() == chainIDInt {
			return true, nil
		}
	}

	return false, nil
}
