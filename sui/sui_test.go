package sui

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckSuiTransactionConfirmed(t *testing.T) {
	tests := []struct {
		name    string
		chainId string
		txHash  string
		want    bool
		wantErr bool
	}{
		{
			name:    "Valid transaction",
			chainId: "3002",
			txHash:  "9K4ab1uerCaTbSBfF2GANgV5nQcTSQn2oCntPvxkSiGL",
			want:    true,
			wantErr: false,
		},
		{
			name:    "Invalid transaction hash",
			chainId: "3002",
			txHash:  "invalid",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckSuiTransactionConfirmed(tt.chainId, tt.txHash)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetSuiTransfersError(t *testing.T) {
	tests := []struct {
		Name        string
		ChainID     string
		TxHash      string
		ExpectError bool
	}{
		{
			Name:        "Valid sui transfer",
			ChainID:     "3002",
			TxHash:      "9K4ab1uerCaTbSBfF2GANgV5nQcTSQn2oCntPvxkSiGL",
			ExpectError: false,
		},
		{
			Name:        "Empty TxHash",
			ChainID:     "11",
			TxHash:      "",
			ExpectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			resp, err := GetSuiTransfers(tt.ChainID, tt.TxHash)
			if tt.ExpectError {
				if err == nil {
					t.Fatalf("getSuiTransfers() expected error %v, but got nil", tt.Name)
				}

				require.Error(t, err)
			} else {
				fmt.Println(resp)
				require.NoError(t, err)
			}
		})
	}
}
