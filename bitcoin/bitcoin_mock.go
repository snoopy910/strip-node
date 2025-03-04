package bitcoin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/StripChain/strip-node/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockSimnetResponse creates a mock Bitcoin Simnet server for testing
// It returns:
// - The mock server that should be closed after use
// - The chain configuration for easy test setup
type MockSimnetResponse struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
	ID     int         `json:"id"`
}

func MockSimnetServer(t *testing.T) (*httptest.Server, common.Chain) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		assert.Equal(t, "POST", r.Method)

		// Verify auth header
		username, password, ok := r.BasicAuth()
		require.True(t, ok)
		assert.Equal(t, "testuser", username)
		assert.Equal(t, "testpass", password)

		// Parse request
		var request struct {
			Method string        `json:"method"`
			Params []interface{} `json:"params"`
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		require.NoError(t, err)

		// Handle different RPC methods
		var response MockSimnetResponse
		switch request.Method {
		case "sendrawtransaction":
			// Simulate successful transaction broadcast
			response = MockSimnetResponse{
				Result: "2c0c07c0a5aa99d54f4c6e5cdb599a2b3e42abb5c8dcf02e8ee0a49c77e9c561",
				Error:  nil,
				ID:     1,
			}
		case "getrawtransaction":
			// Simulate raw transaction fetch
			response = MockSimnetResponse{
				Result: "020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff00ffffffff0118ddf505000000001600148d7a0a3461e3891dbdf560f0f45c18fe5dd2274d00000000",
				Error:  nil,
				ID:     1,
			}
		case "gettxout":
			// Simulate UTXO fetch
			response = MockSimnetResponse{
				Result: map[string]interface{}{
					"value":         0.01,
					"scriptPubKey": "00148d7a0a3461e3891dbdf560f0f45c18fe5dd2274d",
				},
				Error: nil,
				ID:    1,
			}
		default:
			// Unknown method
			response = MockSimnetResponse{
				Result: nil,
				Error: map[string]interface{}{
					"code":    -32601,
					"message": "Method not found",
				},
				ID: 1,
			}
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))

	// Create chain config
	chain := common.Chain{
		ChainId:     "1003", // Simnet
		ChainUrl:    server.URL,
		RpcUsername: "testuser",
		RpcPassword: "testpass",
		ChainType:   "bitcoin",
		KeyCurve:    "bitcoin_ecdsa",
	}

	return server, chain
}