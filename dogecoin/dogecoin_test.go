package dogecoin

import (
	"testing"

	"github.com/test-go/testify/require"
)

func TestIsValidDogecoinAddress(t *testing.T) {
	address := "D1aUkq8VYXNqBwZwJHxMzJv2yf6jg5F7p9"
	require.True(t, ValidateDogeAddress(address))
}

// test to get transaction by hash
func TestGetTransactionByHash(t *testing.T) {
	tests := []struct {
		Name        string
		ChainID     string
		TxHash      string
		ExpectError bool
	}{
		{
			Name:    "Transfer on Mainnet",
			ChainID: "2000",
			TxHash:  "0abfeecf6099d1cfbd93c1258b6248280da029cd4fa8d2d86c1536ff41a51820",
		},
		{
			Name:        "Invalid TxHash",
			ChainID:     "2000",
			TxHash:      "0abfeecf6099d1cfbd93c1258b6248280da019cd4fa8d2d86c1536ff41a51820",
			ExpectError: true,
		},
		{
			Name:        "Invalid chainID",
			ChainID:     "2",
			TxHash:      "0abfeecf6099d1cfbd93c1258b6248280da029cd4fa8d2d86c1536ff41a51820",
			ExpectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			transfers, err := GetDogeTransfers(tt.ChainID, tt.TxHash)
			if tt.ExpectError {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			require.NotNil(t, transfers)
		})
	}
}

// test to check transaction confirmed
func TestCheckTransactionConfirmed(t *testing.T) {
	tests := []struct {
		Name        string
		ChainID     string
		TxHash      string
		ExpectError bool
	}{
		{
			Name:    "Transfer on Mainnet",
			ChainID: "2000",
			TxHash:  "0abfeecf6099d1cfbd93c1258b6248280da029cd4fa8d2d86c1536ff41a51820",
		},
		{
			Name:        "Invalid TxHash",
			ChainID:     "2000",
			TxHash:      "0abfeecf6099d1cfbd93c1258b6248280da019cd4fa8d2d86c1536ff41a51820",
			ExpectError: true,
		},
		{
			Name:        "Invalid chainID",
			ChainID:     "2",
			TxHash:      "0abfeecf6099d1cfbd93c1258b6248280da029cd4fa8d2d86c1536ff41a51820",
			ExpectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			confirmed, err := CheckDogeTransactionConfirmed(tt.ChainID, tt.TxHash)
			if tt.ExpectError {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			require.NotNil(t, confirmed)
			require.True(t, confirmed)
		})
	}
}
