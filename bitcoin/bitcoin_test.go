package bitcoin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
)

// mockGetChain creates a mock GetChain function for testing
func mockGetChain(chainId string, chainUrl string) GetChainFunc {
	return func(id string) (common.Chain, error) {
		if id != chainId {
			return common.Chain{}, fmt.Errorf("chain not found")
		}
		return common.Chain{
			ChainId:     id,
			ChainType:   "bitcoin",
			ChainUrl:    chainUrl,
			KeyCurve:    "ecdsa",
			TokenSymbol: "BTC",
		}, nil
	}
}

// mockBitcoinRPCServer creates a test server that mocks Bitcoin RPC responses
func mockBitcoinRPCServer(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

// BitcoinRPCResponse represents a Bitcoin RPC response
type BitcoinRPCResponse struct {
	Result json.RawMessage `json:"result"`
	Error  interface{}     `json:"error"`
	ID     string          `json:"id"`
}

// TestGetBitcoinTransfersMainnet tests Bitcoin transfers on mainnet
func TestGetBitcoinTransfersMainnet(t *testing.T) {
	tests := []struct {
		name          string
		txHash        string
		mockResponses map[string]interface{} // Map of method to response
		wantTransfers int
		wantFee       int64
		wantInputs    int64
		wantOutputs   int64
	}{
		{
			name:   "Single Input Multiple Output",
			txHash: "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponses: map[string]interface{}{
				"getrawtransaction": map[string]interface{}{
					"txid": "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
					"vin": []map[string]interface{}{
						{
							"txid": "previous_tx_id",
							"vout": 0,
						},
					},
					"vout": []map[string]interface{}{
						{
							"value": 0.00169806,
							"scriptPubKey": map[string]interface{}{
								"address": "1MR5zo89V2ygCZm6AiVsVQ2vKVk1Tjmp7i",
							},
						},
						{
							"value": 0.01799400,
							"scriptPubKey": map[string]interface{}{
								"address": "3DGxAYYUA61WrrdbBac8Ra9eA9peAQwTJF",
							},
						},
					},
				},
				"getrawtransaction_prev": map[string]interface{}{
					"vout": []map[string]interface{}{
						{
							"value": 0.01991454,
							"scriptPubKey": map[string]interface{}{
								"address": "3DGxAYYUA61WrrdbBac8Ra9eA9peAQwTJF",
							},
						},
					},
				},
			},
			wantTransfers: 2,
			wantFee:       22248,
			wantInputs:    1991454,
			wantOutputs:   1969206,
		},
		{
			name:   "Transaction with OP_RETURN",
			txHash: "605dc2ce2c1c9d95f3f83ab2b146ef97fde3b4df15b0990b38eb06edf41fabb0",
			mockResponses: map[string]interface{}{
				"getrawtransaction": map[string]interface{}{
					"txid": "605dc2ce2c1c9d95f3f83ab2b146ef97fde3b4df15b0990b38eb06edf41fabb0",
					"vin": []map[string]interface{}{
						{
							"txid": "previous_tx_id_2",
							"vout": 0,
						},
					},
					"vout": []map[string]interface{}{
						{
							"value": 0,
							"scriptPubKey": map[string]interface{}{
								"type": "nulldata",
								"asm":  "OP_RETURN 58325b766621fb519aa4eceffdebdf0646dd40de581a325701e5bb46fd86ac19d5195e7b1d68b469d24a7e634eef1217a1a04d4ff937d0277eb3c05666a90394e4e100000d696e0076000d3ad500be6b",
							},
						},
						{
							"value": 0.00150000,
							"scriptPubKey": map[string]interface{}{
								"address": "bc1qs0kkdpsrzh3ngqgth7mkavlwlzr7lms2zv3wxe",
							},
						},
						{
							"value": 0.00150000,
							"scriptPubKey": map[string]interface{}{
								"address": "bc1qxf5ephyanvpxe593a3kg36cx0k92dq3yua46n2",
							},
						},
						{
							"value": 1.31506892,
							"scriptPubKey": map[string]interface{}{
								"address": "1NDxDDSHVHvSv48vd27NNHkXHYZjDdVLss",
							},
						},
					},
				},
				"getrawtransaction_prev": map[string]interface{}{
					"vout": []map[string]interface{}{
						{
							"value": 1.31814113,
							"scriptPubKey": map[string]interface{}{
								"address": "1NDxDDSHVHvSv48vd27NNHkXHYZjDdVLss",
							},
						},
					},
				},
			},
			wantTransfers: 3,
			wantFee:       7221,
			wantInputs:    131814113,
			wantOutputs:   131806892,
		},
		{
			name:   "Multiple inputs from same address with OP_RETURN",
			txHash: "c505340c5b8b36b02a036f6b33dff3e30494abd5772d138abb0f73a358ddc71c",
			mockResponses: map[string]interface{}{
				"getrawtransaction": map[string]interface{}{
					"txid": "c505340c5b8b36b02a036f6b33dff3e30494abd5772d138abb0f73a358ddc71c",
					"vin": []map[string]interface{}{
						{
							"txid": "previous_tx_id_3",
							"vout": 0,
						},
						{
							"txid": "previous_tx_id_4",
							"vout": 0,
						},
					},
					"vout": []map[string]interface{}{
						{
							"value": 0,
							"scriptPubKey": map[string]interface{}{
								"type": "nulldata",
								"asm":  "OP_RETURN 6a4c50582a5b7666b8d9b5a2c0a8d0e5d2c1b8d4b5a2c0a8d0e5d2c1b8d4",
							},
						},
						{
							"value": 0.00036899,
							"scriptPubKey": map[string]interface{}{
								"address": "bc1qqu0y5z2gu87trnegj5zgz4l5udn6q9jvkzq6qh",
							},
						},
						{
							"value": 1.1504780,
							"scriptPubKey": map[string]interface{}{
								"address": "bc1q02kcjhqwg6agyhyxg74veucf63n8l5yjamrvn6",
							},
						},
						{
							"value": 18.8495608,
							"scriptPubKey": map[string]interface{}{
								"address": "bc1qc9zpn3j4guy4fktq8qt2h86uprrjmm2f58ddah",
							},
						},
					},
				},
				"getrawtransaction_prev": map[string]interface{}{
					"vout": []map[string]interface{}{
						{
							"value": 0.00039677,
							"scriptPubKey": map[string]interface{}{
								"address": "bc1q67zeuej5a20hje2f8mcxmnczvyrx66pzn6n9pk",
							},
						},
						{
							"value": 20.00000000,
							"scriptPubKey": map[string]interface{}{
								"address": "bc1qvhhgq66mwjugagus9eql05cask4xrkdjh9uh6n",
							},
						},
					},
				},
			},
			wantTransfers: 6,
			wantFee:       2390,
			wantInputs:    200039677,
			wantOutputs:   200037287,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original defaultGetChain
			origGetChain := defaultGetChain
			defer func() {
				defaultGetChain = origGetChain
			}()

			// Create mock server
			server := mockBitcoinRPCServer(t, func(w http.ResponseWriter, r *http.Request) {
				var req map[string]interface{}
				json.NewDecoder(r.Body).Decode(&req)
				method := req["method"].(string)

				response := BitcoinRPCResponse{
					ID: req["id"].(string),
				}

				if mockResp, ok := tt.mockResponses[method]; ok {
					respBytes, _ := json.Marshal(mockResp)
					response.Result = respBytes
				} else if method == "getrawtransaction" {
					params := req["params"].([]interface{})
					if txid, ok := params[0].(string); ok && txid == "previous_tx_id" {
						respBytes, _ := json.Marshal(tt.mockResponses["getrawtransaction_prev"])
						response.Result = respBytes
					} else if txid == "previous_tx_id_2" {
						respBytes, _ := json.Marshal(tt.mockResponses["getrawtransaction_prev"])
						response.Result = respBytes
					} else if txid == "previous_tx_id_3" {
						respBytes, _ := json.Marshal(tt.mockResponses["getrawtransaction_prev"])
						response.Result = respBytes
					} else if txid == "previous_tx_id_4" {
						respBytes, _ := json.Marshal(tt.mockResponses["getrawtransaction_prev"])
						response.Result = respBytes
					} else {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			})
			defer server.Close()

			// Override defaultGetChain for testing
			defaultGetChain = mockGetChain("1000", server.URL)

			// Test GetBitcoinTransfers
			transfers, feeDetails, err := GetBitcoinTransfers("1000", tt.txHash)
			require.NoError(t, err)
			require.NotNil(t, transfers)
			require.NotNil(t, feeDetails)

			// Verify number of transfers
			require.Len(t, transfers, tt.wantTransfers)

			// Verify fee details
			require.Equal(t, tt.wantFee, feeDetails.FeeAmount)
			require.Equal(t, tt.wantInputs, feeDetails.TotalInputs)
			require.Equal(t, tt.wantOutputs, feeDetails.TotalOutputs)

			// Verify all transfers have required fields
			for _, transfer := range transfers {
				require.NotEmpty(t, transfer.From)
				require.NotEmpty(t, transfer.To)
				require.NotEmpty(t, transfer.Amount)
				require.NotEmpty(t, transfer.ScaledAmount)
				require.Equal(t, BTC_TOKEN_SYMBOL, transfer.Token)
				require.True(t, transfer.IsNative)
			}
		})
	}
}

// TestGetBitcoinTransfersTestnet tests Bitcoin transfers on testnet
func TestGetBitcoinTransfersTestnet(t *testing.T) {
	t.Run("Basic BTC Transfer", func(t *testing.T) {
		// Save original defaultGetChain
		origGetChain := defaultGetChain
		defer func() {
			defaultGetChain = origGetChain
		}()

		// Create mock server
		server := mockBitcoinRPCServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			var req map[string]interface{}
			json.NewDecoder(r.Body).Decode(&req)
			method := req["method"].(string)

			response := BitcoinRPCResponse{
				ID: req["id"].(string),
			}

			if method == "getrawtransaction" {
				params := req["params"].([]interface{})
				txid := params[0].(string)
				if txid == "d7a9ea7629ab6183a5f9b01a445830dbcd9b1998c7efd18373e67dc27917d96b" {
					mockResp := map[string]interface{}{
						"txid": txid,
						"vin": []map[string]interface{}{
							{
								"txid": "previous_tx_id",
								"vout": 0,
							},
						},
						"vout": []map[string]interface{}{
							{
								"value": 5.94061934,
								"scriptPubKey": map[string]interface{}{
									"address": "muvVhrJivbcYwe9Bs4cecDBU3eEc8KFhzu",
								},
							},
							{
								"value": 0,
								"scriptPubKey": map[string]interface{}{
									"type": "nulldata",
									"asm":  "OP_RETURN 6a4c5048454d4901006da1c6001fa3062e5cf548dca050e2f3ab849eabfd362cef095a3b5f2d14796fe7644b4f38523fcbdb187a71b9cccbdb9d53f3af1b318796a35e14e42441cf92eba28253682c38d2e2dc",
								},
							},
						},
					}
					respBytes, _ := json.Marshal(mockResp)
					response.Result = respBytes
				} else if txid == "previous_tx_id" {
					mockResp := map[string]interface{}{
						"vout": []map[string]interface{}{
							{
								"value": 5.96626934,
								"scriptPubKey": map[string]interface{}{
									"address": "muvVhrJivbcYwe9Bs4cecDBU3eEc8KFhzu",
								},
							},
						},
					}
					respBytes, _ := json.Marshal(mockResp)
					response.Result = respBytes
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(response)
		})
		defer server.Close()

		// Override defaultGetChain for testing
		defaultGetChain = mockGetChain("1001", server.URL)

		// Test GetBitcoinTransfers
		transfers, feeDetails, err := GetBitcoinTransfers("1001", "d7a9ea7629ab6183a5f9b01a445830dbcd9b1998c7efd18373e67dc27917d96b")
		require.NoError(t, err)
		require.NotNil(t, transfers)
		require.NotNil(t, feeDetails)

		// Verify transfer details (self-transfer)
		require.Len(t, transfers, 1)
		require.Equal(t, "muvVhrJivbcYwe9Bs4cecDBU3eEc8KFhzu", transfers[0].From)
		require.Equal(t, "muvVhrJivbcYwe9Bs4cecDBU3eEc8KFhzu", transfers[0].To)
		require.Equal(t, "5.94061934", transfers[0].Amount)
		require.Equal(t, "594061934", transfers[0].ScaledAmount)
		require.Equal(t, BTC_TOKEN_SYMBOL, transfers[0].Token)
		require.True(t, transfers[0].IsNative)

		// Verify fee details
		require.Equal(t, int64(2565000), feeDetails.FeeAmount)
		require.Equal(t, int64(596626934), feeDetails.TotalInputs)
		require.Equal(t, int64(594061934), feeDetails.TotalOutputs)
	})
}

// TestGetBitcoinTransfersErrors tests error cases for GetBitcoinTransfers
func TestGetBitcoinTransfersErrors(t *testing.T) {
	// Save original defaultGetChain
	origGetChain := defaultGetChain
	defer func() {
		defaultGetChain = origGetChain
	}()

	tests := []struct {
		name         string
		chainId      string
		txHash       string
		mockResponse func(w http.ResponseWriter)
		expectError  bool
	}{
		{
			name:    "Invalid Chain ID",
			chainId: "invalid",
			txHash:  "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				http.Error(w, "Chain not found", http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name:    "Invalid Transaction Hash",
			chainId: "1000",
			txHash:  "invalid_hash",
			mockResponse: func(w http.ResponseWriter) {
				response := BitcoinRPCResponse{
					Error: map[string]interface{}{
						"code":    -8,
						"message": "parameter 1 must be hexadecimal string",
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			},
			expectError: true,
		},
		{
			name:    "API Server Error",
			chainId: "1000",
			txHash:  "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				response := BitcoinRPCResponse{
					Error: map[string]interface{}{
						"code":    -32603,
						"message": "Internal error",
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			},
			expectError: true,
		},
		{
			name:    "Malformed JSON Response",
			chainId: "1000",
			txHash:  "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{"invalid_json":`)
			},
			expectError: true,
		},
		{
			name:    "Network Timeout",
			chainId: "1000",
			txHash:  "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				time.Sleep(100 * time.Millisecond) // Simulate delay
				http.Error(w, "Timeout", http.StatusGatewayTimeout)
			},
			expectError: true,
		},
		{
			name:    "Empty Response",
			chainId: "1000",
			txHash:  "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				response := BitcoinRPCResponse{
					Result: json.RawMessage(`{"vin": [], "vout": []}`),
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			},
			expectError: true,
		},
		{
			name:    "Missing Output Addresses",
			chainId: "1000",
			txHash:  "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				response := BitcoinRPCResponse{
					Result: json.RawMessage(`{
						"vin": [{"txid": "prev_tx", "vout": 0}],
						"vout": [{
							"value": 1.0,
							"scriptPubKey": {
								"type": "pubkeyhash"
							}
						}]
					}`),
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server if chainId is valid
			if tt.chainId == "1000" || tt.chainId == "1001" {
				server := mockBitcoinRPCServer(t, func(w http.ResponseWriter, r *http.Request) {
					tt.mockResponse(w)
				})
				defer server.Close()
				defaultGetChain = mockGetChain(tt.chainId, server.URL)
			}

			// Call the function
			transfers, feeDetails, err := GetBitcoinTransfers(tt.chainId, tt.txHash)

			// Check error
			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, transfers)
				require.Nil(t, feeDetails)
			} else {
				require.NoError(t, err)
				require.NotNil(t, transfers)
				require.NotNil(t, feeDetails)
			}
		})
	}
}

// TestIsValidBitcoinAddress tests the isValidBitcoinAddress function
func TestIsValidBitcoinAddress(t *testing.T) {
	t.Run("Valid BTC Address", func(t *testing.T) {
		address := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
		require.True(t, isValidBitcoinAddress(address))
	})

	t.Run("Invalid BTC Address", func(t *testing.T) {
		address := "InvalidAddress123"
		require.False(t, isValidBitcoinAddress(address))
	})
}

func TestSendBitcoinTransaction(t *testing.T) {
	// Create a valid Simnet address
	pubKeyHash := []byte{118, 169, 82, 84, 36, 23, 237, 52, 71, 4, 218, 23, 74, 247, 134, 103, 127, 119, 163, 82}
	_, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &chaincfg.SimNetParams)
	require.NoError(t, err)

	// Test transaction data
	serializedTxn := "01000000018be5c497555fd91ac3987e7eea01c89f7f7a9f7d1a1f8033f10e780e9fee2d820000000000ffffffff021027000000000000160014bd1979541c3abe7557f0c2a28bb7607cdc62b2dffca09a3b000000001600148d7a0a3461e3891dbdf560f0f45c18fe5dd2274d00000000"
	chainId := "1003" // Simnet
	keyCurve := "bitcoin_ecdsa"
	dataToSign := "badce2c5be9dba537fe0e370681699d029e777eb2abb185aa5ae244ea78012f4"
	address := "021b43d4eda394393e130a333af4ac6d553c2f34b3aeed2dcaa2d6c7bb6139bbae"
	signatureHex := "d18297fe4cf25a05077be893bc14fefed531a93d1abd5269f01ffc235539240a11b721b18225ff9597fdd99af4f80720df7846257e74408d02f888abff45c59f"

	// Setup mock Simnet server
	server, chain := MockSimnetServer(t)
	defer server.Close()

	// Setup test chain
	common.Chains = []common.Chain{chain}

	rlt, err := SendBitcoinTransaction(serializedTxn, chainId, keyCurve, dataToSign, address, signatureHex)
	require.NoError(t, err)
	require.NotEmpty(t, rlt)
	t.Log(rlt)
}

func TestSendBitcoinTransaction1(t *testing.T) {
	// Test transaction data
	serializedTxn := "01000000019ec24b70e3d51c41fe4934aa01eb4f74a96d7ce77c8d2e6bdf85fbd9ffb2515e0100000000ffffffff0210270000000000001600146c94a6480381743e80688254aebc309b109b2737fca09a3b000000001976a914967582f2f3d2650f849ba35abb9ba6066087331c88ac00000000"
	chainId := "1002" // Simnet
	keyCurve := "bitcoin_ecdsa"
	dataToSign := "afdb7db5c03e4d70d7f889d415fb6303eb7828e04bbbc24213d3b92f45fc90a0"
	address := "0395d93e206424c4f7b784771197dbf261104215970723e2f0de5a4762bec5d49c"
	signatureHex := "304402200c082668f732dfbbce9605091617b91046df70a5b81d03e453c45fa86a6f8d4002203356d9601a05492ed2407f5e249339d29ae6bb9ff69621bcc9a855563f57f78e01"

	// Convert raw signature to DER format
	derSignature, err := derEncode(signatureHex)
	require.NoError(t, err)

	// Verify signature before sending transaction
	isValid := VerifyECDSASignature(dataToSign, derSignature, address)
	require.True(t, isValid, "Signature verification failed")

	rlt, err := SendBitcoinTransaction(serializedTxn, chainId, keyCurve, dataToSign, address, derSignature)
	require.NoError(t, err)
	require.NotEmpty(t, rlt)
	t.Log(rlt)
}

func TestSendBitcoinTransaction2(t *testing.T) {
	// Test transaction data
	serializedTxn := "01000000011c876af1d62cc6be3b091df0513f2d44549dced4ccf4522b88ede70e10eb18260100000000ffffffff0210270000000000001600149543d35e1aaacd5312bc95bad32b13b5e8efbd1ffca09a3b000000001976a914f13f7b52722dd5889f6958284f3532e7868dc72d88ac00000000"
	chainId := "1002" // Simnet
	keyCurve := "bitcoin_ecdsa"
	dataToSign := "d67eacda0cc3734ef46667a273ad450ab34a97f05109b567f011989e0939a664"
	address := "02fb4aa16cbe8c63051880899ca3a423f3bd4e51d6a5f8dd5be7019cefdd2ab83d"
	signatureHex := "e7a474d7ad4f4ed6fcd86cbf0710e24141ae59d9b33981af3f5fc604b94beced4887836d0d83901d0c83cc64d76d5b702f239801fda22adc4ccd5717854df456"

	// Convert raw signature to DER format
	derSignature, err := derEncode(signatureHex)
	require.NoError(t, err)

	// Verify signature before sending transaction
	isValid := VerifyECDSASignature(dataToSign, derSignature, address)
	require.True(t, isValid, "Signature verification failed")

	rlt, err := SendBitcoinTransaction(serializedTxn, chainId, keyCurve, dataToSign, address, derSignature)
	require.NoError(t, err)
	require.NotEmpty(t, rlt)
	t.Log(rlt)
}

func TestGetBitcoinTransfersRegtest(t *testing.T) {
	chainId := "1002"
	txHash := "ff90c3b8d2bcbf7c216a4d7fc67aad2e7fe48c7857bacf0b8814b3acef38d43c"
	rlt, _, err := GetBitcoinTransfers(chainId, txHash)
	require.NoError(t, err)
	require.NotEmpty(t, rlt)
	t.Log(rlt)

	isConfirmed, err := CheckBitcoinTransactionConfirmed(chainId, txHash)
	require.NoError(t, err)
	require.True(t, isConfirmed)
}

func TestDerEncode(t *testing.T) {
	signatureHex := "4e48cf9a2f08be3e29a29b66c56a079535f09b0a4d22a05eecc85bc65a6a5c987a15b6e7f942f8b4b0a3ac09a3f5da0ed5d8687b4f2ac47cfe3cd170b01e98ab"
	expected := "304402204e48cf9a2f08be3e29a29b66c56a079535f09b0a4d22a05eecc85bc65a6a5c9802207a15b6e7f942f8b4b0a3ac09a3f5da0ed5d8687b4f2ac47cfe3cd170b01e98ab01"

	encoded, err := derEncode(signatureHex)
	if err != nil {
		t.Fatal(err)
	}

	if encoded != expected {
		t.Errorf("DER encoding mismatch:\nwant: %s\ngot:  %s", expected, encoded)
	}
}
