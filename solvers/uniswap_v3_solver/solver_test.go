package uniswap_v3_solver

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
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
	testPoolAddress      = "0x8Fc6e3247c13C56747cE4ae4E7621013763d98BC" // Example pool address
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
		Fee:       500,
		TickLower: -180,
		TickUpper: 180,
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

func TestMintFlow(t *testing.T) {
	// Setup private key from hex string
	privateKeyHex := "" // Replace with actual private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)
	testRecipientAddress := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	solver := setupTestSolver(t)
	intent := createTestIntent("mint")
	intent.Identity = testRecipientAddress

	// Test Construct
	txData, err := solver.Construct(intent, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, txData)

	// Sign the transaction data
	txDataBytes := []byte(txData)
	signature, err := crypto.Sign(crypto.Keccak256(txDataBytes), privateKey)
	require.NoError(t, err)
	signatureHex := "0x" + hex.EncodeToString(signature)

	// Test Solve with actual signature
	txHash, err := solver.Solve(intent, 0, signatureHex)
	require.NoError(t, err)
	assert.NotEmpty(t, txHash)

	// Test Status with polling
	maxAttempts := 20
	for i := 0; i < maxAttempts; i++ {
		status, err := solver.Status(intent, 0)
		require.NoError(t, err)

		if status == "success" || status == "failed" {
			break
		}
		if i == maxAttempts-1 {
			t.Fatalf("Transaction did not complete within %d attempts", maxAttempts)
		}

		time.Sleep(2 * time.Second)
	}

	// Test GetOutput
	output, err := solver.GetOutput(intent, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, output)

	// Decode and verify output
	var lpOutput LPOutput
	err = json.Unmarshal([]byte(output), &lpOutput)
	require.NoError(t, err)
	assert.NotEmpty(t, lpOutput.TxHash)
	assert.NotZero(t, lpOutput.TokenId)
	assert.NotEmpty(t, lpOutput.AmountA)
	assert.NotEmpty(t, lpOutput.AmountB)
}

func TestExitFlow(t *testing.T) {
	solver := setupTestSolver(t)
	intent := createTestIntent("exit")

	// Test Construct
	txData, err := solver.Construct(intent, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, txData)

	// Test Solve
	txHash, err := solver.Solve(intent, 0, "0x1234") // Mock signature
	require.NoError(t, err)
	assert.NotEmpty(t, txHash)

	// Test Status
	status, err := solver.Status(intent, 0)
	require.NoError(t, err)
	assert.Equal(t, "pending", status)

	// Test GetOutput
	output, err := solver.GetOutput(intent, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, output)

	// Decode and verify output
	var lpOutput LPOutput
	err = json.Unmarshal([]byte(output), &lpOutput)
	require.NoError(t, err)
	assert.NotEmpty(t, lpOutput.TxHash)
	assert.NotZero(t, lpOutput.TokenId)
	assert.NotEmpty(t, lpOutput.AmountA)
	assert.NotEmpty(t, lpOutput.AmountB)
}
