package bitcoin

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/stretchr/testify/require"
)

// mockBlockCypherServer creates a mock server for testing
func mockBlockCypherServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

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

// TestGetBitcoinTransfersMainnet tests Bitcoin transfers on mainnet
func TestGetBitcoinTransfersMainnet(t *testing.T) {
	tests := []struct {
		name          string
		txHash        string
		mockResponse  BlockCypherTransaction
		wantTransfers int
		wantFee       int64
		wantInputs    int64
		wantOutputs   int64
	}{
		{
			name:   "Single Input Multiple Output",
			txHash: "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: BlockCypherTransaction{
				Hash: "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
				Inputs: []BlockCypherInput{
					{
						Addresses:   []string{"3DGxAYYUA61WrrdbBac8Ra9eA9peAQwTJF"},
						OutputValue: 1991454,
					},
				},
				Outputs: []BlockCypherOutput{
					{
						Addresses: []string{"1MR5zo89V2ygCZm6AiVsVQ2vKVk1Tjmp7i"},
						Value:     169806,
					},
					{
						Addresses: []string{"3DGxAYYUA61WrrdbBac8Ra9eA9peAQwTJF"},
						Value:     1799400,
					},
				},
				Fees: 22248,
			},
			wantTransfers: 2,
			wantFee:       22248,
			wantInputs:    1991454,
			wantOutputs:   1969206,
		},
		{
			name:   "Transaction with OP_RETURN",
			txHash: "605dc2ce2c1c9d95f3f83ab2b146ef97fde3b4df15b0990b38eb06edf41fabb0",
			mockResponse: BlockCypherTransaction{
				Hash: "605dc2ce2c1c9d95f3f83ab2b146ef97fde3b4df15b0990b38eb06edf41fabb0",
				Inputs: []BlockCypherInput{
					{
						Addresses:   []string{"1NDxDDSHVHvSv48vd27NNHkXHYZjDdVLss"},
						OutputValue: 131814113,
					},
				},
				Outputs: []BlockCypherOutput{
					{
						Value:     0,
						Script:    "6a4c5058325b766621fb519aa4eceffdebdf0646dd40de581a325701e5bb46fd86ac19d5195e7b1d68b469d24a7e634eef1217a1a04d4ff937d0277eb3c05666a90394e4e100000d696e0076000d3ad500be6b",
						Addresses: nil,
					},
					{
						Addresses: []string{"bc1qs0kkdpsrzh3ngqgth7mkavlwlzr7lms2zv3wxe"},
						Value:     150000,
					},
					{
						Addresses: []string{"bc1qxf5ephyanvpxe593a3kg36cx0k92dq3yua46n2"},
						Value:     150000,
					},
					{
						Addresses: []string{"1NDxDDSHVHvSv48vd27NNHkXHYZjDdVLss"},
						Value:     131506892,
					},
				},
				Fees: 7221,
			},
			wantTransfers: 3,
			wantFee:       7221,
			wantInputs:    131814113,
			wantOutputs:   131806892,
		},
		{
			name:   "Multiple inputs from same address with OP_RETURN",
			txHash: "c505340c5b8b36b02a036f6b33dff3e30494abd5772d138abb0f73a358ddc71c",
			mockResponse: BlockCypherTransaction{
				Hash: "c505340c5b8b36b02a036f6b33dff3e30494abd5772d138abb0f73a358ddc71c",
				Inputs: []BlockCypherInput{
					{
						Addresses:   []string{"bc1q67zeuej5a20hje2f8mcxmnczvyrx66pzn6n9pk"},
						OutputValue: 39677,
					},
					{
						Addresses:   []string{"bc1qvhhgq66mwjugagus9eql05cask4xrkdjh9uh6n"},
						OutputValue: 200000000,
					},
				},
				Outputs: []BlockCypherOutput{
					{
						Value:     0,
						Script:    "6a4c50582a5b7666b8d9b5a2c0a8d0e5d2c1b8d4b5a2c0a8d0e5d2c1b8d4b5a2c0a8d0e5d2c1b8d4",
						Addresses: nil,
					},
					{
						Addresses: []string{"bc1qqu0y5z2gu87trnegj5zgz4l5udn6q9jvkzq6qh"},
						Value:     36899,
					},
					{
						Addresses: []string{"bc1q02kcjhqwg6agyhyxg74veucf63n8l5yjamrvn6"},
						Value:     11504780,
					},
					{
						Addresses: []string{"bc1qc9zpn3j4guy4fktq8qt2h86uprrjmm2f58ddah"},
						Value:     188495608,
					},
				},
				Fees: 2390,
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
			server := mockBlockCypherServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(tt.mockResponse)
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
		server := mockBlockCypherServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			response := BlockCypherTransaction{
				Hash: "d7a9ea7629ab6183a5f9b01a445830dbcd9b1998c7efd18373e67dc27917d96b",
				Inputs: []BlockCypherInput{
					{
						Addresses:   []string{"muvVhrJivbcYwe9Bs4cecDBU3eEc8KFhzu"},
						OutputValue: 596626934,
					},
				},
				Outputs: []BlockCypherOutput{
					{
						Addresses: []string{"muvVhrJivbcYwe9Bs4cecDBU3eEc8KFhzu"},
						Value:     594061934,
					},
					{
						Value:     0,
						Script:    "6a4c5048454d4901006da1c6001fa3062e5cf548dca050e2f3ab849eabfd362cef095a3b5f2d14796fe7644b4f38523fcbdb187a71b9cccbdb9d53f3af1b318796a35e14e42441cf92eba28253682c38d2e2dc",
						Addresses: nil,
					},
				},
				Fees: 2565000,
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
				http.Error(w, "Transaction not found", http.StatusNotFound)
			},
			expectError: true,
		},
		{
			name:    "API Server Error",
			chainId: "1000",
			txHash:  "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{"inputs": [], "outputs": []}`)
			},
			expectError: true,
		},
		{
			name:    "Missing Output Addresses",
			chainId: "1000",
			txHash:  "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				response := BlockCypherTransaction{
					Inputs: []BlockCypherInput{
						{
							Addresses: []string{"address1"},
						},
					},
					Outputs: []BlockCypherOutput{
						{
							Value:     100000,
							Addresses: []string{},
						},
					},
				}
				json.NewEncoder(w).Encode(response)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server if chainId is valid
			if tt.chainId == "1000" || tt.chainId == "1001" {
				server := mockBlockCypherServer(t, func(w http.ResponseWriter, r *http.Request) {
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

// TestFetchUTXOValue tests the FetchUTXOValue function
func TestFetchUTXOValue(t *testing.T) {
	tests := []struct {
		name          string
		chainUrl      string
		txHash        string
		mockResponse  func(w http.ResponseWriter)
		expectedValue int64
		expectError   bool
	}{
		{
			name:     "Valid Bitcoin Transaction",
			chainUrl: "https://api.blockcypher.com/v1/btc/main",
			txHash:   "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				response := BlockCypherTransaction{
					Outputs: []BlockCypherOutput{
						{
							Value: 1000000,
						},
					},
				}
				json.NewEncoder(w).Encode(response)
			},
			expectedValue: 1000000,
			expectError:   false,
		},
		{
			name:     "Invalid Transaction Hash",
			chainUrl: "https://api.blockcypher.com/v1/btc/main",
			txHash:   "invalid_hash",
			mockResponse: func(w http.ResponseWriter) {
				http.Error(w, "Transaction not found", http.StatusNotFound)
			},
			expectedValue: 0,
			expectError:   true,
		},
		{
			name:     "Invalid Chain URL",
			chainUrl: "invalid_url",
			txHash:   "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				http.Error(w, "Not Found", http.StatusNotFound)
			},
			expectedValue: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server if chainUrl is valid
			if tt.chainUrl != "invalid_url" {
				server := mockBlockCypherServer(t, func(w http.ResponseWriter, r *http.Request) {
					tt.mockResponse(w)
				})
				defer server.Close()
				tt.chainUrl = server.URL
			}

			// Call the function
			value, err := FetchUTXOValue(tt.chainUrl, tt.txHash)

			// Check error
			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, int64(0), value)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedValue, value)
			}
		})
	}
}

// TestFetchUTXOValueErrors tests error cases for UTXO value fetching
func TestFetchUTXOValueErrors(t *testing.T) {
	tests := []struct {
		name          string
		chainUrl      string
		txHash        string
		mockResponse  func(w http.ResponseWriter)
		expectedError string
	}{
		{
			name:     "Invalid Chain URL",
			chainUrl: "invalid_url",
			txHash:   "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				http.Error(w, "Not Found", http.StatusNotFound)
			},
			expectedError: "failed to fetch transaction",
		},
		{
			name:     "Invalid Transaction Hash",
			chainUrl: "https://api.blockcypher.com/v1/btc/main",
			txHash:   "invalid_hash",
			mockResponse: func(w http.ResponseWriter) {
				http.Error(w, "Transaction not found", http.StatusNotFound)
			},
			expectedError: "unexpected status code: 404",
		},
		{
			name:     "API Server Error",
			chainUrl: "https://api.blockcypher.com/v1/btc/main",
			txHash:   "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			},
			expectedError: "unexpected status code: 500",
		},
		{
			name:     "Malformed JSON Response",
			chainUrl: "https://api.blockcypher.com/v1/btc/main",
			txHash:   "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{"invalid_json":`)
			},
			expectedError: "failed to decode transaction response",
		},
		{
			name:     "Network Timeout",
			chainUrl: "https://api.blockcypher.com/v1/btc/main",
			txHash:   "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				time.Sleep(100 * time.Millisecond) // Simulate delay
				http.Error(w, "Timeout", http.StatusGatewayTimeout)
			},
			expectedError: "unexpected status code: 504",
		},
		{
			name:     "Empty Response",
			chainUrl: "https://api.blockcypher.com/v1/btc/main",
			txHash:   "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, `{"inputs": [], "outputs": []}`)
			},
			expectedError: "transaction has no outputs",
		},
		{
			name:     "No Outputs in Response",
			chainUrl: "https://api.blockcypher.com/v1/btc/main",
			txHash:   "dbe01947bffa898a0ed281c29227f5a810bc43775076412918ec9519a70789ae",
			mockResponse: func(w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				response := BlockCypherTransaction{
					Inputs: []BlockCypherInput{
						{
							Addresses: []string{"address1"},
						},
					},
					Outputs: []BlockCypherOutput{},
				}
				json.NewEncoder(w).Encode(response)
			},
			expectedError: "transaction has no outputs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := mockBlockCypherServer(t, func(w http.ResponseWriter, r *http.Request) {
				tt.mockResponse(w)
			})
			defer server.Close()

			// Override the chain URL if it's a valid one
			if tt.chainUrl != "invalid_url" {
				tt.chainUrl = server.URL
			}

			// Call the function
			value, err := FetchUTXOValue(tt.chainUrl, tt.txHash)

			// Verify error
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedError)
			require.Equal(t, int64(0), value)
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

func TestBitcoinPsbtSerilizedTxn(t *testing.T) {
	serializedTxn := "cHNidP8BAHQCAAAAAQrNIZka3cgdzRlhZXX1YwG5zwZCdmtwg/M3eiBjuJJHAAAAAAD/////AmQAAAAAAAAAFgAUhHfs0h8qvtrZ9eCDwcAAk7NR8j2SyZo7AAAAABl2qRTupX02brO38k+CEHq86KDosx1PUoisAAAAAAABASIAypo7AAAAABl2qRTupX02brO38k+CEHq86KDosx1PUoisAAAA"

	unsignedTx, err := ParseSerializedTransaction(serializedTxn)
	require.NoError(t, err)
	require.NotNil(t, unsignedTx)

	// Verify transaction details
	require.Equal(t, 2, int(unsignedTx.Version))
	require.Equal(t, 1, len(unsignedTx.TxIn))
	require.Equal(t, 2, len(unsignedTx.TxOut))

	// Verify input details
	txIn := unsignedTx.TxIn[0]
	require.Equal(t, "4792b863207a37f383706b764206cfb90163f575656119cd1dc8dd1a9921cd0a", txIn.PreviousOutPoint.Hash.String())
	require.Equal(t, uint32(0), txIn.PreviousOutPoint.Index)
	require.Equal(t, uint32(4294967295), txIn.Sequence)

	// Verify output details
	require.Equal(t, int64(100), unsignedTx.TxOut[0].Value)
	require.Equal(t, 22, len(unsignedTx.TxOut[0].PkScript))
	require.Equal(t, "00148477ecd21f2abedad9f5e083c1c00093b351f23d", hex.EncodeToString(unsignedTx.TxOut[0].PkScript))

	require.Equal(t, int64(999999890), unsignedTx.TxOut[1].Value)
	require.Equal(t, 25, len(unsignedTx.TxOut[1].PkScript))
	require.Equal(t, "76a914eea57d366eb3b7f24f82107abce8a0e8b31d4f5288ac", hex.EncodeToString(unsignedTx.TxOut[1].PkScript))
}

func TestSendBitcoinTransaction(t *testing.T) {
	serializedTxn := "cHNidP8BAHQCAAAAAbzOlSXVYBO8x+11Btz6/aQStOQn7AcGl+MPj6wYWEyqAAAAAAD/////AugDAAAAAAAAFgAUSGWbePkEw7w7okG9RjjsaTwOKVskxJo7AAAAABl2qRT2cW9Eqq1ybINmC64ZyBJ+IgabJ4isAAAAAAABASIAypo7AAAAABl2qRT2cW9Eqq1ybINmC64ZyBJ+IgabJ4isAAAA"
	chainId := "1002"
	keyCurve := "secp256k1"
	dataToSign := "7bb3a60651a0c7e5a39ac42d308bc6b35fa66cd799414f796c1145fd326292f4"
	signatureHex := "d8949b22405f76ef5331cd6c0219b6b214e177ad8e92491904999ee5db2a150a2354612ca392ddebabf66c4fb79754ecf1c322ce047d28b6f89fa2ce677a8cdc01"

	rlt, err := SendBitcoinTransaction(serializedTxn, chainId, keyCurve, dataToSign, signatureHex)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rlt)
}
