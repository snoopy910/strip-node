package common

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type EthClient interface {
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error)
}

func EstimateTransactionGas(fromAddress common.Address, toAddress *common.Address, value int64, gasPrice *big.Int, tipCap *big.Int, feeCap *big.Int, data []byte, client EthClient, gasMultiplier float64) (uint64, error) {
	factor := gasMultiplier
	if factor < 1 {
		factor = 1
	}
	msg := ethereum.CallMsg{
		From:      fromAddress,
		To:        toAddress,
		Value:     big.NewInt(value),
		GasPrice:  gasPrice,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		Data:      data,
	}

	gas, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		return 0, err
	}
	gas = uint64(float64(gas) * factor)
	return gas, nil
}
