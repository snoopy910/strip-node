package bitcoin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/StripChain/strip-node/common"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithdrawBitcoinGetSignature(t *testing.T) {
	// Setup mock Simnet server
	server, chain := MockSimnetServer(t)
	defer server.Close()

	// Setup test chain
	common.Chains = []common.Chain{chain}

	// Test data
	chainId := "1003" // Simnet
	amount := "0.01"
	fromAddress := "021b43d4eda394393e130a333af4ac6d553c2f34b3aeed2dcaa2d6c7bb6139bbae"
	toAddress := "tb1qw508d6qejxtdg4y5r3zarvary0c5xw7kxpjzsx"

	// Test signature generation
	dataToSign, err := WithdrawBitcoinGetSignature(chainId, fromAddress, amount, toAddress)
	require.NoError(t, err)
	require.NotEmpty(t, dataToSign)
	t.Log(dataToSign)
}

func TestWithdrawBitcoinTxn(t *testing.T) {
	// Mock chain data
	chainId := "1003" // Simnet

	// Create a test transaction first
	pubKeyHash := []byte{118, 169, 82, 84, 36, 23, 237, 52, 71, 4, 218, 23, 74, 247, 134, 103, 127, 119, 163, 82}
	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &chaincfg.SimNetParams)
	require.NoError(t, err)
	simnetAddress := addr.String()

	// Get unsigned transaction
	transaction, err := WithdrawBitcoinGetSignature(chainId, simnetAddress, "0.01", simnetAddress)
	require.NoError(t, err)

	// This is a 32-byte signature in hex format (64 characters)
	signature := "f4d169783bef70e57e26790dce5a70ea4bf38968db8e7a9522bcb681a6e8c4c7c20c4e42e5296f0d24bc9d0972a0dd9f2a1a22c52516765076bc3c1f2886066f"

	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth header
		username, password, ok := r.BasicAuth()
		require.True(t, ok)
		assert.Equal(t, "testuser", username)
		assert.Equal(t, "testpass", password)

		// Verify request body
		var rpcRequest map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&rpcRequest)
		require.NoError(t, err)
		assert.Equal(t, "sendrawtransaction", rpcRequest["method"])

		// Send successful response
		response := map[string]interface{}{
			"result": "2c0c07c0a5aa99d54f4c6e5cdb599a2b3e42abb5c8dcf02e8ee0a49c77e9c561",
			"error":  nil,
			"id":     1,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Setup test chain
	common.Chains = []common.Chain{
		{
			ChainId:     chainId,
			ChainUrl:    server.URL,
			RpcUsername: "testuser",
			RpcPassword: "testpass",
			ChainType:   "bitcoin",
			KeyCurve:    "ecdsa",
		},
	}

	// Test successful case
	txHash, err := WithdrawBitcoinTxn(chainId, transaction, signature)
	require.NoError(t, err)
	assert.NotEmpty(t, txHash)

	// Test invalid chain ID
	_, err = WithdrawBitcoinTxn("invalid", transaction, signature)
	assert.Error(t, err)

	// Test invalid transaction hex
	_, err = WithdrawBitcoinTxn(chainId, "invalid", signature)
	assert.Error(t, err)

	// Test RPC error response
	common.Chains = []common.Chain{
		{
			ChainId:     chainId,
			ChainUrl:    "http://invalid-url",
			RpcUsername: "testuser",
			RpcPassword: "testpass",
			ChainType:   "bitcoin",
			KeyCurve:    "ecdsa",
		},
	}
	_, err = WithdrawBitcoinTxn(chainId, transaction, signature)
	assert.Error(t, err)
}
