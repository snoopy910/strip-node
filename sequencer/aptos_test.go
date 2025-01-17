package sequencer

import (
	"testing"

	"github.com/StripChain/strip-node/util"
	"github.com/stretchr/testify/require"
)

func TestGetAptosTransfers(t *testing.T) {
	// Test cases
	tests := []struct {
		Name         string
		ChainID      string
		TxHash       string
		From         string
		To           string
		Amount       string
		Token        string
		IsNative     bool
		TokenAddress string
		ScaledAmount string
	}{
		{
			Name:         "Native APT Transfer on Mainnet",
			ChainID:      "11",
			TxHash:       "0x0edbfc90e9a5c2d7b2932335088297f34ea42987113b75c752a26cc8978c2869",
			From:         "0x834d639b10d20dcb894728aa4b9b572b2ea2d97073b10eacb111f338b20ea5d7",
			To:           "0xf6dc6b07e7c49f240d197d1880f2806b7dfa3d838762fb3b9391f707bc6afbf7",
			Amount:       "799.99900000",
			Token:        "APT",
			IsNative:     true,
			TokenAddress: util.ZERO_ADDRESS,
			ScaledAmount: "79999900000",
		},
		{
			Name:         "Aptos Token Transfer on Mainnet",
			ChainID:      "11",
			TxHash:       "0x2a90d9b7f03525120bf2404c332d6531799dd68004ec208b15d1eb6435665176",
			From:         "0xe6f6ecc5eb1d93963c9d9446e5f33d677facd01511cf8bab18290c67641ef083",
			To:           "0x2621fd2b392dc53115dc3edea701ff6e4fe29c56a19388b9f9c455ecb65b1e4c",
			Amount:       "13536.630000",
			Token:        "TOMA",
			IsNative:     false,
			TokenAddress: "0x9d0595765a31f8d56e1d2aafc4d6c76f283c67a074ef8812d8c31bd8252ac2c3",
			ScaledAmount: "13536630000",
		},
		{
			Name:         "Fungible Asset Transfer on Mainnet",
			ChainID:      "11",
			TxHash:       "0xa10016b53c6b37d65e20210e04bfb9f79e09dc7defd74e063bb5f7c436e581f4",
			From:         "0x7c6018b4aa445a63f91954fe36ede071a6289ae5b35cfd117755c2a952713128",
			To:           "0x541e28fb12aa661a30358f2bebcd44460187ec918cb9cee075c2db86ee6aed93",
			Amount:       "1.00000000",
			Token:        "TVS",
			IsNative:     false,
			TokenAddress: "0x43782fca70e1416fc0c75954942dadd4af8d305a608b6153397ad5801b71e72d",
			ScaledAmount: "100000000",
		},
		{
			Name:         "Native APT Transfer on Devnet",
			ChainID:      "167",
			TxHash:       "0xed84af67094ebfe3bf5e531ec9ddda7dc2ebc439ed250e476836ce5a877e1ab6",
			From:         "0x7e9983bf1e8a75e305d081bacc994cea4051ef0d40548f38e0fb2140f20be6d0",
			To:           "0x82c076aaf063b09738b54b7277792252a3abc21a8e78e2422472892dc4c4a788",
			Amount:       "0.00000001",
			Token:        "APT",
			IsNative:     true,
			TokenAddress: util.ZERO_ADDRESS,
			ScaledAmount: "1",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			transfers, err := GetAptosTransfers(tt.ChainID, tt.TxHash)
			if err != nil {
				t.Errorf("getAptosTransfers() error = %v", err)
			}

			if len(transfers) == 0 {
				t.Fatalf("No transfer transaction for txHash, %v", tt.TxHash)
			}

			require.Equal(t, tt.From, transfers[0].From)
			require.Equal(t, tt.To, transfers[0].To)
			require.Equal(t, tt.Amount, transfers[0].Amount)
			require.Equal(t, tt.Token, transfers[0].Token)
			require.Equal(t, tt.IsNative, transfers[0].IsNative)
			require.Equal(t, tt.TokenAddress, transfers[0].TokenAddress)
			require.Equal(t, tt.ScaledAmount, transfers[0].ScaledAmount)
		})
	}
}
