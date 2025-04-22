package sequencer

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
)

func setupTestTransaction(t *testing.T) *solana.Transaction {
	feePayer := solana.MustPublicKeyFromBase58("DpZqkyDKkVv2S7Lhbd5dUVcVCPJz2Lypr4W5Cru2sHr7")
	recipient := solana.MustPublicKeyFromBase58("5oNDL3swdJJF1g9DzJiZ4ynHXgszjAEpUkxVYejchzrY")

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				1000000000,
				feePayer,
				recipient,
			).Build(),
		},
		solana.MustHashFromBase58("CBLp4VEPu9T9W2uzURoawLGqgAQ65LvmUwDYRHymgwbd"),
		solana.TransactionPayer(feePayer),
	)
	assert.NoError(t, err)
	return tx
}

func TestValidateAndOrderSignatures(t *testing.T) {
	tests := []struct {
		name        string
		setupTx     func() *solana.Transaction
		expectError bool
		errorMsg    string
	}{
		{
			name: "Missing signatures",
			setupTx: func() *solana.Transaction {
				tx := setupTestTransaction(t)
				tx.Signatures = []solana.Signature{}
				return tx
			},
			expectError: true,
			errorMsg:    "signature count mismatch",
		},
		{
			name: "Correct signature count",
			setupTx: func() *solana.Transaction {
				tx := setupTestTransaction(t)
				// it is a dummy signature matching the required count
				tx.Signatures = make([]solana.Signature, tx.Message.Header.NumRequiredSignatures)
				return tx
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tx := tc.setupTx()
			err := validateAndOrderSignatures(tx)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSendSolanaTransaction(t *testing.T) {

	// Chains = []common.Chain{
	// 	{
	// 		ChainId:   "901",
	// 		ChainType: "solana",
	// 		ChainUrl:  "https://api.devnet.solana.com",
	// 		KeyCurve:  "eddsa",
	// 	},
	// }

	validSignature := "5jLFtNTCAnHA9uurWhyNNqzwHLwWCaSNrZBWG48AANMGkreX1DYGbkHL2VWNNt2Kz327QwzzsAacJj2YFdSsfkwN"

	// Create and serialize a valid transaction
	tx := setupTestTransaction(t)
	// Serialize just the message, like the client does
	serializedMsg, err := tx.Message.MarshalBinary()
	assert.NoError(t, err)

	tests := []struct {
		name        string
		tx          string
		sig         string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Invalid signature format",
			tx:          base58.Encode(serializedMsg),
			sig:         "invalid_signature",
			expectError: true,
			errorMsg:    "error decoding signature",
		},
		{
			name:        "Valid transaction and signature format",
			tx:          base58.Encode(serializedMsg),
			sig:         validSignature,
			expectError: true,
			errorMsg:    "signature verification failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := sendSolanaTransactionWithValidation(
				tc.tx,
				"901",
				"eddsa",
				"test_data",
				tc.sig,
			)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {

					if !assert.Contains(t, err.Error(), tc.errorMsg) {
						t.Logf("Got error: %v", err)
						t.Logf("Expected to contain: %v", tc.errorMsg)
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
