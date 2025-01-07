package sequencer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test function for Bitcoin mainnet transfers
func TestGetBitcoinTransfersMainnet(t *testing.T) {
	// Define test cases for Bitcoin mainnet transactions
	tests := []struct {
		Name           string // Test case name
		ChainID        string // Chain ID for the network (e.g., Bitcoin mainnet)
		TxHash         string // Transaction hash to fetch details for
		ExpectedFrom   string // Expected "from" address in the transfer
		ExpectedTo     string // Expected "to" address in the transfer
		ExpectedAmount string // Expected amount in Bitcoin (formatted)
		Token          string // Expected token symbol (BTC)
		IsNative       bool   // Whether it's a native token transfer (true for Bitcoin)
		TokenAddress   string // Token address (BTC_ZERO_ADDRESS in this case)
		ScaledAmount   string // Amount in satoshis (scaled version of the amount)
	}{
		{
			Name:           "Basic BTC Transfer",                                               // Test case name
			ChainID:        "1000",                                                             // Example chain ID for Bitcoin mainnet
			TxHash:         "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae", // Transaction hash
			ExpectedFrom:   "3DGxAYYUA61WrrdbBac8Ra9eA9peAQwTJF",                               // Expected sender address
			ExpectedTo:     "1MR5zo89V2ygCZm6AiVsVQ2vKVk1Tjmp7i",                               // Expected recipient address
			ExpectedAmount: "0.00169806",                                                       // Expected amount in Bitcoin (formatted)
			Token:          "BTC",                                                              // Token symbol (Bitcoin)
			IsNative:       true,                                                               // Native token (Bitcoin)
			TokenAddress:   BTC_ZERO_ADDRESS,                                                   // Zero address for Bitcoin (no specific token address)
			ScaledAmount:   "169806",                                                           // Amount in satoshis (scaled version)
		},
	}

	// Loop through all test cases
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Call GetBitcoinTransfers to fetch the transfers
			transfers, err := GetBitcoinTransfers(tt.ChainID, tt.TxHash)
			if err != nil {
				// If an error occurs, fail the test
				t.Errorf("GetBitcoinTransfers() error: %v", err)
				return
			}

			// Check if no transfers were found, which is unexpected
			if len(transfers) == 0 {
				t.Fatalf("No transfers found for TxHash: %v", tt.TxHash)
			}

			// Compare expected values with the actual transfer details
			require.Equal(t, tt.ExpectedFrom, transfers[0].From)         // Check "from" address
			require.Equal(t, tt.ExpectedTo, transfers[0].To)             // Check "to" address
			require.Equal(t, tt.ExpectedAmount, transfers[0].Amount)     // Check formatted amount
			require.Equal(t, tt.Token, transfers[0].Token)               // Check token symbol
			require.Equal(t, tt.IsNative, transfers[0].IsNative)         // Check if it's a native transfer
			require.Equal(t, tt.TokenAddress, transfers[0].TokenAddress) // Check token address
			require.Equal(t, tt.ScaledAmount, transfers[0].ScaledAmount) // Check scaled amount in satoshis
		})
	}
}

// Test function for Bitcoin testnet transfers
func TestGetBitcoinTransfersTestnet(t *testing.T) {
	// Define test cases for Bitcoin testnet transactions
	tests := []struct {
		Name           string // Test case name
		ChainID        string // Chain ID for the network (e.g., Bitcoin testnet)
		TxHash         string // Transaction hash to fetch details for
		ExpectedFrom   string // Expected "from" address in the transfer
		ExpectedTo     string // Expected "to" address in the transfer
		ExpectedAmount string // Expected amount in Bitcoin (formatted)
		Token          string // Expected token symbol (BTC)
		IsNative       bool   // Whether it's a native token transfer (true for Bitcoin)
		TokenAddress   string // Token address (BTC_ZERO_ADDRESS in this case)
		ScaledAmount   string // Amount in satoshis (scaled version of the amount)
	}{
		{
			Name:           "Basic BTC Transfer",                                               // Test case name
			ChainID:        "1001",                                                             // Example chain ID for Bitcoin testnet
			TxHash:         "d7a9ea7629ab6183a5f9b01a445830dbcd9b1998c7efd18373e67dc27917d96b", // Transaction hash
			ExpectedFrom:   "muvVhrJivbcYwe9Bs4cecDBU3eEc8KFhzu",                               // Expected sender address
			ExpectedTo:     "muvVhrJivbcYwe9Bs4cecDBU3eEc8KFhzu",                               // Expected recipient address
			ExpectedAmount: "5.94061934",                                                       // Expected amount in Bitcoin (formatted)
			Token:          "BTC",                                                              // Token symbol (Bitcoin)
			IsNative:       true,                                                               // Native token (Bitcoin)
			TokenAddress:   BTC_ZERO_ADDRESS,                                                   // Zero address for Bitcoin (no specific token address)
			ScaledAmount:   "594061934",                                                        // Amount in satoshis (scaled version)
		},
	}

	// Loop through all test cases
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Call GetBitcoinTransfers to fetch the transfers
			transfers, err := GetBitcoinTransfers(tt.ChainID, tt.TxHash)
			if err != nil {
				// If an error occurs, fail the test
				t.Errorf("GetBitcoinTransfers() error: %v", err)
				return
			}

			// Check if no transfers were found, which is unexpected
			if len(transfers) == 0 {
				t.Fatalf("No transfers found for TxHash: %v", tt.TxHash)
			}

			// Compare expected values with the actual transfer details
			require.Equal(t, tt.ExpectedFrom, transfers[0].From)         // Check "from" address
			require.Equal(t, tt.ExpectedTo, transfers[0].To)             // Check "to" address
			require.Equal(t, tt.ExpectedAmount, transfers[0].Amount)     // Check formatted amount
			require.Equal(t, tt.Token, transfers[0].Token)               // Check token symbol
			require.Equal(t, tt.IsNative, transfers[0].IsNative)         // Check if it's a native transfer
			require.Equal(t, tt.TokenAddress, transfers[0].TokenAddress) // Check token address
			require.Equal(t, tt.ScaledAmount, transfers[0].ScaledAmount) // Check scaled amount in satoshis
		})
	}
}
