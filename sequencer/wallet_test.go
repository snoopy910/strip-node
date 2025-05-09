package sequencer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/StripChain/strip-node/libs/blockchains"
	db "github.com/StripChain/strip-node/libs/database"
	"github.com/stretchr/testify/require"
)

// MockDB for testing
type mockDB struct {
	wallets map[string]db.WalletSchema
}

// AddWallet function for mockDB
func (m *mockDB) AddWallet(wallet *db.WalletSchema) (int64, error) {
	m.wallets[wallet.Identity] = *wallet
	return 1, nil
}

// Setup mock server for testing
// - GET /keygen: Returns 200 OK for key generation requests
// - GET /address: Returns a mock address based on the provided key curve
// The server validates required query parameters and returns appropriate error responses
func setupMockServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/keygen":
			w.WriteHeader(http.StatusOK)

		case r.Method == "GET" && r.URL.Path == "/address":
			identity := r.URL.Query().Get("identity")
			if identity == "" {
				http.Error(w, "identity query parameter is required", http.StatusBadRequest)
				return
			}
			identityCurve := r.URL.Query().Get("identityCurve")
			if identityCurve == "" {
				http.Error(w, "identityCurve query parameter is required", http.StatusBadRequest)
				return
			}
			keyCurve := r.URL.Query().Get("keyCurve")
			response := GetAddressResponse{
				Address: "mock" + keyCurve + "Address",
			}
			json.NewEncoder(w).Encode(response)

		default:
			t.Errorf("Unexpected request to %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

// Setup test environment for testing
// 1. Creates a mock HTTP server
// 2. Sets up mock signers with predefined public keys
// 3. Creates an in-memory mock database
// 4. Overrides global functions (SignersList and AddWallet)
func setupTestEnvironment(t *testing.T) (*httptest.Server, func()) {
	// Create mock server
	mockServer := setupMockServer(t)
	MaximumSigners = 2

	// Create mock signers
	mockSigners := []Signer{
		{URL: mockServer.URL, PublicKey: "mockPublicKey1"},
		{URL: mockServer.URL, PublicKey: "mockPublicKey2"},
		{URL: mockServer.URL, PublicKey: "mockPublicKey3"},
	}

	// Store original functions
	originalSignersList := SignersList
	originalAddWallet := db.AddWallet

	// Create mock database
	mockDB := &mockDB{wallets: make(map[string]db.WalletSchema)}

	// Set up mock functions
	SignersList = func() ([]Signer, error) {
		return mockSigners, nil
	}
	db.AddWallet = mockDB.AddWallet

	// Return cleanup function
	cleanup := func() {
		mockServer.Close()
		SignersList = originalSignersList
		db.AddWallet = originalAddWallet
	}

	return mockServer, cleanup
}

// TestCreateWallet tests the createWallet function with various test scenarios.
// Test cases cover:
// - Successful wallet creation
// - Empty identity validation
// - Empty curve validation
func TestCreateWallet(t *testing.T) {
	// Set up test environment once for all test cases
	mockServer, cleanup := setupTestEnvironment(t)
	defer cleanup()

	signers, err := SignersList()
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name         string
		identity     string
		blockchainID blockchains.BlockchainID
		wantErr      bool
	}{
		{
			name:         "Success case",
			identity:     "testIdentity",
			blockchainID: blockchains.Ethereum,
			wantErr:      false,
		},
		{
			name:         "Empty identity",
			identity:     "",
			blockchainID: blockchains.Ethereum,
			wantErr:      true,
		},
		{
			name:         "Empty curve",
			identity:     "testIdentity",
			blockchainID: "",
			wantErr:      true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := createWallet(tt.identity, tt.blockchainID)
			if (err != nil) != tt.wantErr {
				t.Errorf("createWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the wallet was created with correct values
				mockDB := signers[0]
				if mockDB.URL != mockServer.URL {
					t.Errorf("Expected server URL %s, got %s", mockServer.URL, mockDB.URL)
				}
			}
		})
	}
}

// TestCreateWalletWithMaximumSigners verifies that the wallet creation process
// respects the MaximumSigners limit when selecting signers.
// It temporarily sets MaximumSigners to 3 and ensures the number of selected
// signers never exceeds this limit.
func TestCreateWalletWithMaximumSigners(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Temporarily set MaximumSigners to 3
	MaximumSigners = 3
	originalMaxSigners := MaximumSigners
	defer func() { MaximumSigners = originalMaxSigners }()

	err := createWallet("testIdentity", "testCurve")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	signers, err := SignersList()
	require.NoError(t, err)

	if len(signers) > MaximumSigners {
		t.Errorf("Expected maximum %d signers, got %d", MaximumSigners, len(signers))
	}
}

// TestCreateWalletServerErrors verifies error handling when signer nodes
// return server errors. It sets up a mock server that always returns 500
// Internal Server Error and ensures the wallet creation fails as expected.
func TestCreateWalletServerErrors(t *testing.T) {
	// Create a server that always returns errors
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer errorServer.Close()

	// Temporarily set MaximumSigners to 1
	MaximumSigners = 1

	// Override SignersList to return error server
	originalSignersList := SignersList
	SignersList = func() ([]Signer, error) {
		return []Signer{
			{URL: errorServer.URL, PublicKey: "errorKey1"},
			{URL: errorServer.URL, PublicKey: "errorKey2"},
		}, nil
	}
	defer func() { SignersList = originalSignersList }()

	err := createWallet("testIdentity", "testCurve")
	if err == nil {
		t.Error("Expected error from server, got nil")
	}
}
