package stellar

import (
	"testing"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/require"
)

var (
	ErrInvalidAmount = "amount not found in solver output"
	ErrInvalidAcc    = "problem: https://stellar.org/horizon-errors/bad_request"
)

func TestWithdrawStellarNativeGetSignature(t *testing.T) {
	tests := []struct {
		name         string
		client       *horizonclient.Client
		bridgeAddr   string
		solverOutput string
		userAddr     string
		expectError  bool
		errorMessage string
	}{
		{
			name:         "Valid native XLM withdrawal on mainnet",
			client:       horizonclient.DefaultPublicNetClient,
			bridgeAddr:   "GDQP2KPQGKIHYJGXNUIYOMHARUARCA7DJT5FO2FFOOKY3B2WSQHG4W37",
			solverOutput: `{"amount":"100.0000000"}`,
			userAddr:     "GB665AL5ZCZ3G53EZQJMSHYBXUL2JGV53ZFVDA2ZU3Z2IGEONI4CEFMV",
			expectError:  false,
		},
		{
			name:         "Valid native XLM withdrawal on mainnet",
			client:       horizonclient.DefaultPublicNetClient,
			bridgeAddr:   "GCIS55JJ7VHNBLVTOZKRZGPQDWOACTSXBEYDX4YKU3HVIJQA6MURTX4E",
			solverOutput: `{"amount":"50.0000000"}`,
			userAddr:     "GABFQIK63R2NETJM7T673EAMZN4RJLLGP3OFUEJU5SZVTGWUKULZJNL6",
			expectError:  false,
		},
		{
			name:         "Invalid solver output",
			client:       horizonclient.DefaultTestNetClient,
			bridgeAddr:   "GDQP2KPQGKIHYJGXNUIYOMHARUARCA7DJT5FO2FFOOKY3B2WSQHG4W37",
			solverOutput: `{"invalid":"json"}`,
			userAddr:     "GBZXN7PIRZGNMHGA7MUUUF4GWPY5AYPV6LY4UV2GL6VJGIQRXFDNMADI",
			expectError:  true,
			errorMessage: ErrInvalidAmount,
		},
		{
			name:         "Invalid bridge account",
			client:       horizonclient.DefaultTestNetClient,
			bridgeAddr:   "invalid_account",
			solverOutput: `{"amount":"100.0000000"}`,
			userAddr:     "GBZXN7PIRZGNMHGA7MUUUF4GWPY5AYPV6LY4UV2GL6VJGIQRXFDNMADI",
			expectError:  true,
			errorMessage: ErrInvalidAcc,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn, dataToSign, err := WithdrawStellarNativeGetSignature(
				tt.client,
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

func TestWithdrawStellarAssetGetSignature(t *testing.T) {
	tests := []struct {
		name         string
		client       *horizonclient.Client
		bridgeAddr   string
		solverOutput string
		userAddr     string
		assetCode    string
		assetIssuer  string
		expectError  bool
		errorMessage string
	}{
		{
			name:         "Valid USDC withdrawal on mainnet",
			client:       horizonclient.DefaultPublicNetClient,
			bridgeAddr:   "GDQP2KPQGKIHYJGXNUIYOMHARUARCA7DJT5FO2FFOOKY3B2WSQHG4W37",
			solverOutput: `{"amount":"100.0000000"}`,
			userAddr:     "GABFQIK63R2NETJM7T673EAMZN4RJLLGP3OFUEJU5SZVTGWUKULZJNL6",
			assetCode:    "USDC",
			assetIssuer:  "GA5ZSEJYB37JRC5AVCIA5MOP4RHTM335X2KGX3IHOJAPP5RE34K4KZVN",
			expectError:  false,
		},
		{
			name:         "Invalid solver output",
			client:       horizonclient.DefaultTestNetClient,
			bridgeAddr:   "GDQP2KPQGKIHYJGXNUIYOMHARUARCA7DJT5FO2FFOOKY3B2WSQHG4W37",
			solverOutput: `{"invalid":"json"}`,
			userAddr:     "GBZXN7PIRZGNMHGA7MUUUF4GWPY5AYPV6LY4UV2GL6VJGIQRXFDNMADI",
			assetCode:    "USDC",
			assetIssuer:  "GA5ZSEJYB37JRC5AVCIA5MOP4RHTM335X2KGX3IHOJAPP5RE34K4KZVN",
			expectError:  true,
			errorMessage: ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn, dataToSign, err := WithdrawStellarAssetGetSignature(
				tt.client,
				tt.bridgeAddr,
				tt.solverOutput,
				tt.userAddr,
				tt.assetCode,
				tt.assetIssuer,
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

func TestWithdrawStellarTxn(t *testing.T) {

	accountID := "GDQP2KPQGKIHYJGXNUIYOMHARUARCA7DJT5FO2FFOOKY3B2WSQHG4W37"
	account, err := horizonclient.DefaultPublicNetClient.AccountDetail(horizonclient.AccountRequest{AccountID: accountID})
	require.NoError(t, err)
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &account,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: "GABFQIK63R2NETJM7T673EAMZN4RJLLGP3OFUEJU5SZVTGWUKULZJNL6",
					Amount:      "1",
					Asset:       txnbuild.CreditAsset{Code: "USDC", Issuer: "GA5ZSEJYB37JRC5AVCIA5MOP4RHTM335X2KGX3IHOJAPP5RE34K4KZVN"},
				},
			},
			BaseFee:       txnbuild.MinBaseFee,
			Memo:          nil,
			Preconditions: txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)
	require.NoError(t, err)

	// Get the transaction in XDR format
	txeB64, err := tx.Base64()
	require.NoError(t, err)

	tests := []struct {
		name          string
		client        *horizonclient.Client
		serializedTxn string
		signature     string
		expectError   bool
		errorMessage  string
	}{
		{
			name:          "Valid transaction submission",
			client:        horizonclient.DefaultPublicNetClient,
			serializedTxn: txeB64,
			signature:     "bf37bab19150d32715cb042de8bc43f1c1ce57f7749ff0e26ec1ab99822ab9ccb9f1347f3fec47adf55d6fc0c616f0d5f181c899f27c97d51d21725e35eecf0f",
			expectError:   true,
			errorMessage:  "tx_bad_auth",
		},
		{
			name:          "Invalid transaction envelope",
			client:        horizonclient.DefaultTestNetClient,
			serializedTxn: "invalid_transaction",
			signature:     "AAAAAA==",
			expectError:   true,
			errorMessage:  "error decoding transaction",
		},
		{
			name:          "Invalid signature",
			client:        horizonclient.DefaultTestNetClient,
			serializedTxn: "AAAAAgAAAACRLvUp/U7QrrN2VRyZ8B2cAU5XCTA78wqmz1QmAPMpGQAAAGQDOO0mAAAAGAAAAAEAAAAAAAAAAAAAAABns10uAAAAAAAAAAEAAAAAAAAAAQAAAAACWCFe3HTSTSz8/f2QDMt5FK1mftxaETTss1ma1FUXlAAAAAFVU0RDAAAAADuZETgO/piLoKiQDrHP5E82b32+lGvtB3JA9/Yk3xXFAAAAADuaygAAAAAAAAAAAA==",
			signature:     "invalid_signature",
			expectError:   true,
			errorMessage:  "error decoding signature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txHash, err := WithdrawStellarTxn(
				tt.client,
				tt.serializedTxn,
				tt.signature,
			)
			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMessage)
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, txHash)
		})
	}
}
