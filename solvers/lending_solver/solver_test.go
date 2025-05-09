package lending_solver

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testRPCURL      = "" // MODIFY ME: RPC URL
	testChainID     = 44331
	testLendingPool = "0x4a3f5a545210cD900d7AA2330816Bc945BDe7567"
	testToken       = "0x3Ea4F15Bab4FcCeA81fD8287d5FD3C80b458d379"
	testAmount      = "1000000000000000000" // 1 token in wei
	privateKeyHex   = ""                    // MODIFY ME: Set this to test private key for integration tests
)

func setupTestSolver(t *testing.T) *LendingSolver {
	solver, err := NewLendingSolver(testRPCURL, testChainID, testLendingPool)
	require.NoError(t, err)
	require.NotNil(t, solver)
	return solver
}

func createTestIntent(action string) Intent {
	metadata := LendingMetadata{
		Action: action,
		Token:  testToken,
		Amount: uint256{Int: testAmount},
	}
	metadataBytes, _ := json.Marshal(metadata)

	return Intent{
		Operations: []Operation{
			{
				SolverMetadata: metadataBytes,
			},
		},
	}
}

func TestConstruct(t *testing.T) {
	tests := []struct {
		name    string
		action  string
		wantErr bool
	}{
		{
			name:    "Supply",
			action:  "supply",
			wantErr: false,
		},
		{
			name:    "Borrow",
			action:  "borrow",
			wantErr: false,
		},
		{
			name:    "Repay",
			action:  "repay",
			wantErr: false,
		},
		{
			name:    "Withdraw",
			action:  "withdraw",
			wantErr: false,
		},
		{
			name:    "Invalid Action",
			action:  "invalid",
			wantErr: true,
		},
	}

	solver := setupTestSolver(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intent := createTestIntent(tt.action)
			hash, err := solver.Construct(intent, 0)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, hash)
			// Transaction hash should be 0x + 64 hex characters (32 bytes)
			assert.True(t, strings.HasPrefix(hash, "0x"))
			assert.Len(t, hash[2:], 64)
			_, err = hex.DecodeString(hash[2:])
			assert.NoError(t, err)
		})
	}
}

func TestSupplyFlow(t *testing.T) {
	if privateKeyHex == "" {
		return
	}
	// Setup private key from hex string
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)
	testRecipientAddress := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	solver := setupTestSolver(t)
	intent := createTestIntent("supply")
	intent.Identity = testRecipientAddress

	dataToSign, err := solver.Construct(intent, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, dataToSign)

	txHashBytes, err := hex.DecodeString(strings.TrimPrefix(dataToSign, "0x"))
	require.NoError(t, err)
	signature, err := crypto.Sign(txHashBytes, privateKey)
	require.NoError(t, err)
	signatureHex := "0x" + hex.EncodeToString(signature)

	txHash, err := solver.Solve(intent, 0, signatureHex)
	require.NoError(t, err)
	assert.NotEmpty(t, txHash)

	intent.Operations[0].Result = txHash

	// Test Status with polling
	maxAttempts := 20
	for i := 0; i < maxAttempts; i++ {
		status, err := solver.Status(intent, 0)
		require.NoError(t, err)

		if status == "success" || status == "failed" {
			break
		}

		time.Sleep(2 * time.Second)
	}

	// Test GetOutput
	output, err := solver.GetOutput(intent, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, output)

	// Decode and verify output
	var lendingOutput LendingOutput
	err = json.Unmarshal([]byte(output), &lendingOutput)
	require.NoError(t, err)
	assert.NotEmpty(t, lendingOutput.TxHash)
	assert.NotEmpty(t, lendingOutput.Amount)
	assert.NotEmpty(t, lendingOutput.Token)
}

func TestBorrowFlow(t *testing.T) {
	if privateKeyHex == "" {
		return
	}
	// Setup private key from hex string
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)
	testRecipientAddress := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	solver := setupTestSolver(t)
	intent := createTestIntent("borrow")
	intent.Identity = testRecipientAddress

	dataToSign, err := solver.Construct(intent, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, dataToSign)

	txHashBytes, err := hex.DecodeString(strings.TrimPrefix(dataToSign, "0x"))
	require.NoError(t, err)
	signature, err := crypto.Sign(txHashBytes, privateKey)
	require.NoError(t, err)
	signatureHex := "0x" + hex.EncodeToString(signature)

	txHash, err := solver.Solve(intent, 0, signatureHex)
	require.NoError(t, err)
	assert.NotEmpty(t, txHash)

	intent.Operations[0].Result = txHash

	// Test Status with polling
	maxAttempts := 20
	for i := 0; i < maxAttempts; i++ {
		status, err := solver.Status(intent, 0)
		require.NoError(t, err)

		if status == "success" || status == "failed" {
			break
		}

		time.Sleep(2 * time.Second)
	}

	// Test GetOutput
	output, err := solver.GetOutput(intent, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, output)

	// Decode and verify output
	var lendingOutput LendingOutput
	err = json.Unmarshal([]byte(output), &lendingOutput)
	require.NoError(t, err)
	assert.NotEmpty(t, lendingOutput.TxHash)
	assert.NotEmpty(t, lendingOutput.Amount)
	assert.NotEmpty(t, lendingOutput.RemainingDebt)
	assert.NotEmpty(t, lendingOutput.HealthFactor)
}

func TestRepayFlow(t *testing.T) {
	if privateKeyHex == "" {
		return
	}
	// Setup private key from hex string
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)
	testRecipientAddress := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	solver := setupTestSolver(t)
	intent := createTestIntent("repay")
	intent.Identity = testRecipientAddress

	dataToSign, err := solver.Construct(intent, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, dataToSign)

	txHashBytes, err := hex.DecodeString(strings.TrimPrefix(dataToSign, "0x"))
	require.NoError(t, err)
	signature, err := crypto.Sign(txHashBytes, privateKey)
	require.NoError(t, err)
	signatureHex := "0x" + hex.EncodeToString(signature)

	txHash, err := solver.Solve(intent, 0, signatureHex)
	require.NoError(t, err)
	assert.NotEmpty(t, txHash)

	intent.Operations[0].Result = txHash

	// Test Status with polling
	maxAttempts := 20
	for i := 0; i < maxAttempts; i++ {
		status, err := solver.Status(intent, 0)
		require.NoError(t, err)

		if status == "success" || status == "failed" {
			break
		}

		time.Sleep(2 * time.Second)
	}

	// Test GetOutput
	output, err := solver.GetOutput(intent, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, output)

	// Decode and verify output
	var lendingOutput LendingOutput
	err = json.Unmarshal([]byte(output), &lendingOutput)
	require.NoError(t, err)
	assert.NotEmpty(t, lendingOutput.TxHash)
	assert.NotEmpty(t, lendingOutput.Amount)
	assert.NotEmpty(t, lendingOutput.RemainingDebt)
	assert.NotEmpty(t, lendingOutput.HealthFactor)
}

func TestWithdrawFlow(t *testing.T) {
	if privateKeyHex == "" {
		return
	}
	// Setup private key from hex string
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)
	testRecipientAddress := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	solver := setupTestSolver(t)
	intent := createTestIntent("withdraw")
	intent.Identity = testRecipientAddress

	dataToSign, err := solver.Construct(intent, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, dataToSign)

	txHashBytes, err := hex.DecodeString(strings.TrimPrefix(dataToSign, "0x"))
	require.NoError(t, err)
	signature, err := crypto.Sign(txHashBytes, privateKey)
	require.NoError(t, err)
	signatureHex := "0x" + hex.EncodeToString(signature)

	txHash, err := solver.Solve(intent, 0, signatureHex)
	require.NoError(t, err)
	assert.NotEmpty(t, txHash)

	intent.Operations[0].Result = txHash

	// Test Status with polling
	maxAttempts := 20
	for i := 0; i < maxAttempts; i++ {
		status, err := solver.Status(intent, 0)
		require.NoError(t, err)

		if status == "success" || status == "failed" {
			break
		}

		time.Sleep(2 * time.Second)
	}

	// Test GetOutput
	output, err := solver.GetOutput(intent, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, output)

	// Decode and verify output
	var lendingOutput LendingOutput
	err = json.Unmarshal([]byte(output), &lendingOutput)
	require.NoError(t, err)
	assert.NotEmpty(t, lendingOutput.TxHash)
	assert.NotEmpty(t, lendingOutput.Amount)
	assert.NotEmpty(t, lendingOutput.Token)
}
