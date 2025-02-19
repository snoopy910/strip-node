package stellar

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	ErrInvalidSignature = fmt.Errorf("error decoding signature: encoding/hex: invalid byte: U+0069 'i'")
	ErrBadAuth          = fmt.Errorf("transaction submission failed: Transaction Failed (result codes: map[transaction:tx_bad_auth])")
)

func TestCheckStellarTransactionConfirmed(t *testing.T) {
	tests := []struct {
		name         string
		chainId      string
		txHash       string
		wantSuccess  bool
		expectError  bool
		errorMessage error
	}{
		{
			name:        "Successful transaction on mainnet",
			chainId:     "mainnet",
			txHash:      "a23c036f896bd2b74e51abd638ebda1577076b8cbff1bd28fb52e652939d38dc",
			wantSuccess: true,
		},
		{
			name:        "Successful transaction on testnet",
			chainId:     "testnet",
			txHash:      "391e7dbd7140eb9a1c4496bf53dbbce3fd4613328dc53e6bba3c6c0c2d7707e4",
			wantSuccess: true,
		},
		{
			name:        "Transaction not found on testnet",
			chainId:     "testnet",
			txHash:      "a23c036f896bd2b74e51abd638ebda1577076b8cbff1bd28fb52e652939d38dc",
			expectError: false,
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success, err := CheckStellarTransactionConfirmed(tt.chainId, tt.txHash)
			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, tt.errorMessage, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantSuccess, success)
		})
	}
}

func TestGetStellarTransfers(t *testing.T) {
	tests := []struct {
		name         string
		chainId      string
		txHash       string
		from         string
		to           string
		amount       string
		token        string
		isNative     bool
		tokenAddress string
		scaledAmount string
		expectError  bool
	}{
		{
			name:         "Native XLM transfer on mainnet",
			chainId:      "mainnet",
			txHash:       "ee4167a3d3bd4ade5264c00b0aa1a5e0ac9b883a07ae5045888591a3019baa3e",
			from:         "GCIS55JJ7VHNBLVTOZKRZGPQDWOACTSXBEYDX4YKU3HVIJQA6MURTX4E",
			to:           "GABFQIK63R2NETJM7T673EAMZN4RJLLGP3OFUEJU5SZVTGWUKULZJNL6",
			amount:       "3219.7100000",
			token:        "XLM",
			isNative:     true,
			tokenAddress: "XLM",
			scaledAmount: "3219.7100000",
		},
		{
			name:         "Non-native token transfer on mainnet",
			chainId:      "mainnet",
			txHash:       "7c3ea735dda80280d9d3d7c3076c6b5ad9edea3165e52dfdeb1de6b0a0a30395",
			from:         "GARWVZ4TGEXRFY47HLRVJL5YLLDO5FNOAPCUBYYA7RXLS3KVBA5YZ3DA",
			to:           "GAUA7XL5K54CC2DDGP77FJ2YBHRJLT36CPZDXWPM6MP7MANOGG77PNJU",
			amount:       "500.0000000",
			token:        "USDC",
			isNative:     false,
			tokenAddress: "USDC:GA5ZSEJYB37JRC5AVCIA5MOP4RHTM335X2KGX3IHOJAPP5RE34K4KZVN",
			scaledAmount: "500.0000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transfers, err := GetStellarTransfers(tt.chainId, tt.txHash)
			if tt.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, transfers)

			transfer := transfers[0]
			require.Equal(t, tt.from, transfer.From)
			require.Equal(t, tt.to, transfer.To)
			require.Equal(t, tt.amount, transfer.Amount)
			require.Equal(t, tt.token, transfer.Token)
			require.Equal(t, tt.isNative, transfer.IsNative)
			require.Equal(t, tt.tokenAddress, transfer.TokenAddress)
			require.Equal(t, tt.scaledAmount, transfer.ScaledAmount)
		})
	}
}

func TestSendStellarTxn(t *testing.T) {
	tests := []struct {
		name         string
		chainId      string
		serializedTx string
		keyCurve     string
		dataToSign   string
		signature    string
		expectError  bool
		errorMessage error
	}{
		{
			name:         "Valid transaction on testnet",
			chainId:      "testnet",
			serializedTx: "AAAAAgAAAAA6BrT9L+GG9XS6ZknRuYtcX6Zt6Rsav0PPYxJuSrI0WAAAAGQAEptTAAAAAwAAAAEAAAAAAAAAAAAAAABntyJuAAAAAAAAAAEAAAAAAAAAAQAAAAB/X2atFcr1OOrJ35T6ESjTHI8H9EQpV56/HGy4H1vA3gAAAAAAAAAAAJiWgAAAAAAAAAAA",
			signature:    "d917d6d4334a54c8ed0d069be153de4477efaf6b297cc609626c1b676475db6e0c2406a2d1f7186429ab6f9bdf9678a20723ca7aded6552b685e44c443c9f604",
			expectError:  false,
			errorMessage: nil,
		},
		{
			name:         "Invalid signature",
			chainId:      "testnet",
			serializedTx: "AAAAAgAAAAA6BrT9L+GG9XS6ZknRuYtcX6Zt6Rsav0PPYxJuSrI0WAAAAGQAEptTAAAAAwAAAAEAAAAAAAAAAAAAAABntyJuAAAAAAAAAAEAAAAAAAAAAQAAAAB/X2atFcr1OOrJ35T6ESjTHI8H9EQpV56/HGy4H1vA3gAAAAAAAAAAAJiWgAAAAAAAAAAA",
			signature:    "invalid_signature",
			expectError:  true,
			errorMessage: ErrInvalidSignature,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txHash, err := SendStellarTxn(
				tt.serializedTx,
				tt.chainId,
				tt.keyCurve,
				tt.dataToSign,
				tt.signature,
			)
			if tt.expectError {
				require.Error(t, err)
				if tt.errorMessage != nil {
					require.Equal(t, tt.errorMessage.Error(), err.Error())
				}
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, txHash)
		})
	}
}
