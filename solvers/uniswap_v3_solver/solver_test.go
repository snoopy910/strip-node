package uniswap_v3_solver

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testRPCURL           = "" // RPC URL
	testChainID          = 44331
	testFactoryAddress   = "0xb1a101860602D32A50E0e426CB827ce2121f12D2"
	testNPMAddress       = "0x782Ed0e82F04fBcF8F6De1F609215A6CeD0EdB85"
	testTokenAAddress    = "0xb228f5F03B05137b38C248D73bA591133128faDB" // StripUSD
	testTokenBAddress    = "0xe109F006D577251340F103d3aCE72B56Fdc6E172" // TestTokenA
	testPoolAddress      = "0x1234567890123456789012345678901234567890" // Example pool address
	testAmount           = "1000000000000000000"                        // 1 token in wei
	testRecipientAddress = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
)

func setupTestSolver(t *testing.T) *UniswapV3Solver {
	solver, err := NewUniswapV3Solver(testRPCURL, testChainID, testFactoryAddress, testNPMAddress)
	require.NoError(t, err)
	require.NotNil(t, solver)
	return solver
}

func createTestIntent(action string) Intent {
	metadata := LPMetadata{
		Action:    action,
		Pool:      testPoolAddress,
		TokenA:    testTokenAAddress,
		TokenB:    testTokenBAddress,
		AmountA:   testAmount,
		AmountB:   testAmount,
		Fee:       3000,
		TickLower: -100,
		TickUpper: 100,
		TokenId:   1,
	}
	metadataBytes, _ := json.Marshal(metadata)

	return Intent{
		Operations: []Operation{
			{
				SolverMetadata: metadataBytes,
			},
		},
		Identity: testRecipientAddress,
	}
}

func TestConstruct(t *testing.T) {
	tests := []struct {
		name    string
		action  string
		wantErr bool
	}{
		{
			name:    "mint position",
			action:  "mint",
			wantErr: false,
		},
		{
			name:    "exit position",
			action:  "exit",
			wantErr: false,
		},
		{
			name:    "invalid action",
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
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
			}
		})
	}
}

func TestSolve(t *testing.T) {
	solver := setupTestSolver(t)
	intent := createTestIntent("mint")
	signature := "0x1234567890..." // Replace with valid signature

	txHash, err := solver.Solve(intent, 0, signature)
	require.NoError(t, err)
	assert.NotEmpty(t, txHash)
}

func TestStatus(t *testing.T) {
	solver := setupTestSolver(t)
	intent := createTestIntent("mint")

	status, err := solver.Status(intent, 0)
	require.NoError(t, err)
	assert.Contains(t, []string{"pending", "success", "failed"}, status)
}

func TestGetOutput(t *testing.T) {
	solver := setupTestSolver(t)
	intent := createTestIntent("mint")

	output, err := solver.GetOutput(intent, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, output)
}

func TestConstructMint(t *testing.T) {
	solver := setupTestSolver(t)
	intent := createTestIntent("mint")

	hash, err := solver.Construct(intent, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestConstructExit(t *testing.T) {
	solver := setupTestSolver(t)
	intent := createTestIntent("exit")

	hash, err := solver.Construct(intent, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}
