package cardano

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	ErrInvalidAmount = "amount not found in solver output"
	ErrInvalidAcc    = "failed to parse bridge address"
)

func TestWithdrawCardanoNativeGetSignature(t *testing.T) {
	tests := []struct {
		name         string
		rpcURL       string
		bridgeAddr   string
		solverOutput string
		userAddr     string
		expectError  bool
		errorMessage string
	}{
		{
			name:         "Valid native ADA withdrawal on mainnet",
			rpcURL:       "https://cardano-mainnet.blockfrost.io/api/v0",
			bridgeAddr:   "addr1qxck9k6y05qernnz4c9kx3rh6qphq8dj8h48wwjhkz9j3vmg8k89l5v879tjxt5uzqq4myqga0uhjwypj42qhw2t5mpswz4yxt",
			solverOutput: `{"amount":"100.000000"}`,
			userAddr:     "addr1q8j5m99vrkv9t4wqq9pqunx5w88ej7zs0qxk0rk4fn2mgmg8k89l5v879tjxt5uzqq4myqga0uhjwypj42qhw2t5mpsxjy8vm",
			expectError:  false,
		},
		{
			name:         "Invalid solver output",
			rpcURL:       "https://cardano-mainnet.blockfrost.io/api/v0",
			bridgeAddr:   "addr1qxck9k6y05qernnz4c9kx3rh6qphq8dj8h48wwjhkz9j3vmg8k89l5v879tjxt5uzqq4myqga0uhjwypj42qhw2t5mpswz4yxt",
			solverOutput: `{"invalid":"json"}`,
			userAddr:     "addr1q8j5m99vrkv9t4wqq9pqunx5w88ej7zs0qxk0rk4fn2mgmg8k89l5v879tjxt5uzqq4myqga0uhjwypj42qhw2t5mpsxjy8vm",
			expectError:  true,
			errorMessage: ErrInvalidAmount,
		},
		// {
		// 	name:         "Invalid bridge address",
		// 	rpcURL:       "https://cardano-mainnet.blockfrost.io/api/v0",
		// 	bridgeAddr:   "invalid_address",
		// 	solverOutput: `{"amount":"100.000000"}`,
		// 	userAddr:     "addr1q8j5m99vrkv9t4wqq9pqunx5w88ej7zs0qxk0rk4fn2mgmg8k89l5v879tjxt5uzqq4myqga0uhjwypj42qhw2t5mpsxjy8vm",
		// 	expectError:  true,
		// 	errorMessage: ErrInvalidAcc,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn, dataToSign, err := WithdrawCardanoNativeGetSignature(
				tt.rpcURL,
				tt.bridgeAddr,
				tt.solverOutput,
				tt.userAddr,
			)
			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMessage)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, txn)
			require.NotEmpty(t, dataToSign)
		})
	}
}

func TestWithdrawCardanoTokenGetSignature(t *testing.T) {
	tests := []struct {
		name         string
		rpcURL       string
		bridgeAddr   string
		solverOutput string
		userAddr     string
		policyID     string
		expectError  bool
		errorMessage string
	}{
		{
			name:         "Valid token withdrawal on mainnet",
			rpcURL:       "https://cardano-mainnet.blockfrost.io/api/v0",
			bridgeAddr:   "addr1qxck9k6y05qernnz4c9kx3rh6qphq8dj8h48wwjhkz9j3vmg8k89l5v879tjxt5uzqq4myqga0uhjwypj42qhw2t5mpswz4yxt",
			solverOutput: `{"amount":"100.000000"}`,
			userAddr:     "addr1q8j5m99vrkv9t4wqq9pqunx5w88ej7zs0qxk0rk4fn2mgmg8k89l5v879tjxt5uzqq4myqga0uhjwypj42qhw2t5mpsxjy8vm",
			policyID:     "2aa9c1557fcf8e7caa049fa0911a8724a1cdaf8037fe0b431c6ac664",
			expectError:  false,
		},
		{
			name:         "Invalid solver output",
			rpcURL:       "https://cardano-mainnet.blockfrost.io/api/v0",
			bridgeAddr:   "addr1qxck9k6y05qernnz4c9kx3rh6qphq8dj8h48wwjhkz9j3vmg8k89l5v879tjxt5uzqq4myqga0uhjwypj42qhw2t5mpswz4yxt",
			solverOutput: `{"invalid":"json"}`,
			userAddr:     "addr1q8j5m99vrkv9t4wqq9pqunx5w88ej7zs0qxk0rk4fn2mgmg8k89l5v879tjxt5uzqq4myqga0uhjwypj42qhw2t5mpsxjy8vm",
			policyID:     "2aa9c1557fcf8e7caa049fa0911a8724a1cdaf8037fe0b431c6ac664",
			expectError:  true,
			errorMessage: ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn, dataToSign, err := WithdrawCardanoTokenGetSignature(
				tt.rpcURL,
				tt.bridgeAddr,
				tt.solverOutput,
				tt.userAddr,
				tt.policyID,
			)
			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMessage)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, txn)
			require.NotEmpty(t, dataToSign)
		})
	}
}
