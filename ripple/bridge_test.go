package ripple

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	ErrAccNotFound   = "failed to get account info"
	ErrInvalidFormat = "invalid token address format"
	ErrInvalidAmount = "amount not found in solver output"
	ErrInvalidAcc    = "failed to parse bridge address"
)

func TestWithdrawRippleNativeGetSignature(t *testing.T) {
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
			name:         "Valid native XRP withdrawal on mainnet",
			rpcURL:       "wss://s1.ripple.com:51233",
			bridgeAddr:   "rhTsmUJFpiju7syo8V5UbCQoaJjKWSvZju",
			solverOutput: `{"amount":"100.0000000"}`,
			userAddr:     "r3ZDszz2ieaVjCPxGjMjJWrjVaxQxHgf1o",
			expectError:  false,
		},
		{
			name:         "Valid native XRP withdrawal on testnet",
			rpcURL:       "wss://s.altnet.rippletest.net:51233",
			bridgeAddr:   "rGwCu8RtkX34roGwMr5cDVUsz8yq7SuVPT",
			solverOutput: `{"amount":"50.0000000"}`,
			userAddr:     "rGpGUdqUAVkNVr4Hfkvay7ffB7vjoA31uT",
			expectError:  false,
		},
		{
			name:         "Invalid solver output",
			rpcURL:       "wss://s.altnet.rippletest.net:51233",
			bridgeAddr:   "rhTsmUJFpiju7syo8V5UbCQoaJjKWSvZju",
			solverOutput: `{"invalid":"json"}`,
			userAddr:     "r3ZDszz2ieaVjCPxGjMjJWrjVaxQxHgf1o",
			expectError:  true,
			errorMessage: ErrInvalidAmount,
		},
		{
			name:         "Invalid bridge account",
			rpcURL:       "wss://s.altnet.rippletest.net:51233",
			bridgeAddr:   "invalid_account",
			solverOutput: `{"amount":"100.0000000"}`,
			userAddr:     "r3ZDszz2ieaVjCPxGjMjJWrjVaxQxHgf1o",
			expectError:  true,
			errorMessage: ErrInvalidAcc,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn, dataToSign, err := WithdrawRippleNativeGetSignature(
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

func TestWithdrawRippleTokenGetSignature(t *testing.T) {
	tests := []struct {
		name         string
		rpcURL       string
		bridgeAddr   string
		solverOutput string
		userAddr     string
		tokenCode    string
		tokenIssuer  string
		expectError  bool
		errorMessage string
	}{
		{
			name:         "Valid USD withdrawal on mainnet",
			rpcURL:       "wss://s1.ripple.com:51233",
			bridgeAddr:   "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
			solverOutput: `{"amount":"100.0000000"}`,
			userAddr:     "rPT1Sjq2YGrBMTttX4GZHjKu9dyfzbpAYe",
			tokenCode:    "USD",
			tokenIssuer:  "rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B",
			expectError:  false,
		},
		{
			name:         "Invalid solver output",
			rpcURL:       "wss://s.altnet.rippletest.net:51233",
			bridgeAddr:   "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
			solverOutput: `{"invalid":"json"}`,
			userAddr:     "rPT1Sjq2YGrBMTttX4GZHjKu9dyfzbpAYe",
			tokenCode:    "USD",
			tokenIssuer:  "rvYAfWj5gh67oV6fW32ZzP3Aw4Eubs59B",
			expectError:  true,
			errorMessage: ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn, dataToSign, err := WithdrawRippleTokenGetSignature(
				tt.rpcURL,
				tt.bridgeAddr,
				tt.solverOutput,
				tt.userAddr,
				// tt.tokenCode,
				tt.tokenIssuer,
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
