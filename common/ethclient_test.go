package common

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEthClient struct {
	mock.Mock
}

func (m *MockEthClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	args := m.Called(ctx, call)
	return args.Get(0).(uint64), args.Error(1)
}

func TestEstimateTransactionGasNotAdjusted(t *testing.T) {
	client := new(MockEthClient)
	from := common.HexToAddress("0xaaaa111111111111111111111111111111111111")
	to := common.HexToAddress("0xf222222222222222222222222222222222222222")
	data := []byte("testFunctionData")
	gasPrice := new(big.Int).SetUint64(100000)
	gasMultiplier := 0.9

	expectedMsg := ethereum.CallMsg{
		From:     from,
		To:       &to,
		Value:    big.NewInt(0),
		GasPrice: gasPrice,
		Data:     data,
	}
	expectedGas := uint64(21000)
	client.On("EstimateGas", mock.Anything, expectedMsg).Return(expectedGas, nil)

	gas, err := EstimateTransactionGas(from, &to, 0, gasPrice, nil, nil, data, client, gasMultiplier)
	assert.NoError(t, err)
	assert.Equal(t, expectedGas, gas)

	client.AssertExpectations(t)
}

func TestEstimateTransactionGasAdjusted(t *testing.T) {
	client := new(MockEthClient)
	from := common.HexToAddress("0xaaaa111111111111111111111111111111111111")
	to := common.HexToAddress("0xf222222222222222222222222222222222222222")
	data := []byte("testFunctionData")
	gasPrice := new(big.Int).SetUint64(100000)
	gasMultiplier := 1.2

	expectedMsg := ethereum.CallMsg{
		From:     from,
		To:       &to,
		Value:    big.NewInt(0),
		GasPrice: gasPrice,
		Data:     data,
	}
	expectedGas := uint64(21000)
	client.On("EstimateGas", mock.Anything, expectedMsg).Return(expectedGas, nil)

	gas, err := EstimateTransactionGas(from, &to, 0, gasPrice, nil, nil, data, client, gasMultiplier)
	assert.NoError(t, err)
	assert.Equal(t, uint64(float64(expectedGas)*gasMultiplier), gas)

	client.AssertExpectations(t)
}

func TestEstimateTransactionGasError(t *testing.T) {
	client := new(MockEthClient)
	from := common.HexToAddress("0xaaaa111111111111111111111111111111111111")
	to := common.HexToAddress("0xf222222222222222222222222222222222222222")
	data := []byte("testFunctionData")
	gasPrice := new(big.Int).SetUint64(100000)
	gasMultiplier := 1.2

	expectedMsg := ethereum.CallMsg{
		From:     from,
		To:       &to,
		Value:    big.NewInt(0),
		GasPrice: gasPrice,
		Data:     data,
	}

	mockError := errors.New("failed to estimate gas")
	client.On("EstimateGas", mock.Anything, expectedMsg).Return(uint64(0), mockError)

	_, err := EstimateTransactionGas(from, &to, 0, gasPrice, nil, nil, data, client, gasMultiplier)
	assert.EqualError(t, err, mockError.Error())

	client.AssertExpectations(t)
}
