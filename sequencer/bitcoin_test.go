package sequencer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBitcoinTransfers(t *testing.T) {
	tests := []struct {
		Name           string
		ChainID        string
		TxHash         string
		ExpectedFrom   string
		ExpectedTo     string
		ExpectedAmount string
		Token          string
		IsNative       bool
		TokenAddress   string
		ScaledAmount   string
	}{
		{
			Name:           "Basic BTC Transfer",
			ChainID:        "1000", // Example chain ID for Bitcoin mainnet
			TxHash:         "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			ExpectedFrom:   "3DGxAYYUA61WrrdbBac8Ra9eA9peAQwTJF", // Address from inputs
			ExpectedTo:     "1MR5zo89V2ygCZm6AiVsVQ2vKVk1Tjmp7i", // Address from outputs
			ExpectedAmount: "0.00169806",                         // Value in BTC (169806 satoshis / 1e8)
			Token:          "BTC",
			IsNative:       true,
			TokenAddress:   BTC_ZERO_ADDRESS, // Replace with the constant for zero address
			ScaledAmount:   "169806",         // Amount in satoshis
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			transfers, err := GetBitcoinTransfers(tt.ChainID, tt.TxHash)
			if err != nil {
				t.Errorf("GetBitcoinTransfers() error: %v", err)
				return
			}

			if len(transfers) == 0 {
				t.Fatalf("No transfers found for TxHash: %v", tt.TxHash)
			}

			require.Equal(t, tt.ExpectedFrom, transfers[0].From)
			require.Equal(t, tt.ExpectedTo, transfers[0].To)
			require.Equal(t, tt.ExpectedAmount, transfers[0].Amount)
			require.Equal(t, tt.Token, transfers[0].Token)
			require.Equal(t, tt.IsNative, transfers[0].IsNative)
			require.Equal(t, tt.TokenAddress, transfers[0].TokenAddress)
			require.Equal(t, tt.ScaledAmount, transfers[0].ScaledAmount)
		})
	}
}
