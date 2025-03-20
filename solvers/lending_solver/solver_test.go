package lending_solver

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testRPCURL      = "http://localhost:8545"
	testChainID     = 1337
	testLendingPool = "0x1234567890123456789012345678901234567890"
	testStripUSD    = "0x0987654321098765432109876543210987654321"
	testToken       = "0xabcdef0123456789abcdef0123456789abcdef01"
	testAmount      = "1000000000000000000" // 1 token in wei
)

func setupTestSolver(t *testing.T) *LendingSolver {
	solver, err := NewLendingSolver(testRPCURL, testChainID, testLendingPool, testStripUSD)
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
			assert.True(t, common.IsHexAddress(hash[2:]))
		})
	}
}

func TestSolve(t *testing.T) {
	solver := setupTestSolver(t)
	intent := createTestIntent("supply")

	// First get the hash to sign
	hash, err := solver.Construct(intent, 0)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	// Create a mock signature (65 bytes)
	mockSig := make([]byte, 65)
	mockSig[64] = 27 // v value
	mockSigHex := "0x" + common.Bytes2Hex(mockSig)

	// Test Solve
	txHash, err := solver.Solve(intent, 0, mockSigHex)
	require.NoError(t, err)
	require.NotEmpty(t, txHash)
	assert.True(t, common.IsHexAddress(txHash[2:]))
}

func TestStatus(t *testing.T) {
	solver := setupTestSolver(t)
	intent := createTestIntent("supply")

	// Test initial status (no transaction)
	status, err := solver.Status(intent, 0)
	require.NoError(t, err)
	assert.Equal(t, "pending", status)

	// Add a mock transaction
	mockTxHash := "0x1234567890"
	intent.Operations[0].Result = mockTxHash

	// Test pending status
	solver.txStatus.Store(mockTxHash, &TransactionStatus{
		TxHash: mockTxHash,
		Status: "pending",
	})
	status, err = solver.Status(intent, 0)
	require.NoError(t, err)
	assert.Equal(t, "pending", status)

	// Test success status
	solver.txStatus.Store(mockTxHash, &TransactionStatus{
		TxHash:  mockTxHash,
		Status:  "success",
		Receipt: &types.Receipt{Status: 1},
	})
	status, err = solver.Status(intent, 0)
	require.NoError(t, err)
	assert.Equal(t, "success", status)

	// Test failure status
	solver.txStatus.Store(mockTxHash, &TransactionStatus{
		TxHash:  mockTxHash,
		Status:  "failure",
		Receipt: &types.Receipt{Status: 0},
		Error:   assert.AnError,
	})
	status, err = solver.Status(intent, 0)
	require.Error(t, err)
	assert.Equal(t, "failure", status)
}

func TestGetOutput(t *testing.T) {
	solver := setupTestSolver(t)
	intent := createTestIntent("supply")

	// Add a mock successful transaction
	mockTxHash := "0x1234567890"
	intent.Operations[0].Result = mockTxHash
	solver.txStatus.Store(mockTxHash, &TransactionStatus{
		TxHash:  mockTxHash,
		Status:  "success",
		Receipt: &types.Receipt{Status: 1},
	})

	// Test output
	output, err := solver.GetOutput(intent, 0)
	require.NoError(t, err)
	require.NotEmpty(t, output)

	// Verify output format
	var lendingOutput LendingOutput
	err = json.Unmarshal([]byte(output), &lendingOutput)
	require.NoError(t, err)
	assert.Equal(t, mockTxHash, lendingOutput.TxHash)
	assert.Equal(t, testAmount, lendingOutput.Amount.Int)
}

func TestConstructSupply(t *testing.T) {
	solver := setupTestSolver(t)
	metadata := LendingMetadata{
		Action: "supply",
		Token:  testToken,
		Amount: uint256{Int: testAmount},
	}

	data, err := solver.constructSupply(metadata)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	assert.True(t, len(data) > 2) // More than just "0x"
}

func TestConstructBorrow(t *testing.T) {
	solver := setupTestSolver(t)
	metadata := LendingMetadata{
		Action: "borrow",
		Amount: uint256{Int: testAmount},
	}

	data, err := solver.constructBorrow(metadata)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	assert.True(t, len(data) > 2)
}

func TestConstructRepay(t *testing.T) {
	solver := setupTestSolver(t)
	metadata := LendingMetadata{
		Action: "repay",
		Amount: uint256{Int: testAmount},
	}

	data, err := solver.constructRepay(metadata)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	assert.True(t, len(data) > 2)
}

func TestConstructWithdraw(t *testing.T) {
	solver := setupTestSolver(t)
	metadata := LendingMetadata{
		Action: "withdraw",
		Token:  testToken,
		Amount: uint256{Int: testAmount},
	}

	data, err := solver.constructWithdraw(metadata)
	require.NoError(t, err)
	require.NotEmpty(t, data)
	assert.True(t, len(data) > 2)
}
