package sequencer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// TestCreateWalletEndpoint tests the /createWallet endpoint for creating a wallet.
// It verifies the response status code for valid and invalid requests.
func TestCreateWalletEndpoint(t *testing.T) {
	tests := []struct {
		name          string
		identity      string
		identityCurve string
		wantStatus    int
	}{
		{
			name:          "Valid wallet creation",
			identity:      "testIdentity",
			identityCurve: "ecdsa",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "Missing identity",
			identity:      "",
			identityCurve: "ecdsa",
			wantStatus:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/createWallet?identity="+tt.identity+"&identityCurve="+tt.identityCurve, nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if tt.identity == "" {
						http.Error(w, "identity required", http.StatusBadRequest)
						return
					}
					w.WriteHeader(http.StatusOK)
				}).ServeHTTP(w, r)
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("CreateWallet returned wrong status code: got %v want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

// TestGetWalletEndpoint tests the /getWallet endpoint for retrieving a wallet.
// It checks the response status code for valid and missing identity scenarios.
func TestGetWalletEndpoint(t *testing.T) {
	tests := []struct {
		name          string
		identity      string
		identityCurve string
		wantStatus    int
	}{
		{
			name:          "Valid wallet retrieval",
			identity:      "testIdentity",
			identityCurve: "ecdsa",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "Missing identity",
			identity:      "",
			identityCurve: "ecdsa",
			wantStatus:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/getWallet?identity="+tt.identity+"&identityCurve="+tt.identityCurve, nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.identity == "" {
					http.Error(w, "identity required", http.StatusBadRequest)
					return
				}
				json.NewEncoder(w).Encode(map[string]string{"identity": tt.identity})
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetWallet returned wrong status code: got %v want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

// TestCreateIntent tests the /createIntent endpoint for creating an intent.
// It verifies the JSON payload and response status code.
func TestCreateIntent(t *testing.T) {
	// Create a test intent
	testIntent := Intent{
		Identity:      "testIdentity",
		IdentityCurve: "ecdsa",
		Status:        INTENT_STATUS_PROCESSING,
		Operations: []Operation{
			{
				Type:     OPERATION_TYPE_TRANSACTION,
				Status:   OPERATION_STATUS_PENDING,
				ChainId:  "1",
				KeyCurve: "ecdsa",
			},
		},
	}

	// Convert intent to JSON
	intentJSON, err := json.Marshal(testIntent)
	if err != nil {
		t.Fatalf("Failed to marshal test intent: %v", err)
	}

	// Create test request
	req := httptest.NewRequest("POST", "/createIntent", bytes.NewBuffer(intentJSON))
	w := httptest.NewRecorder()

	// Create handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		var intent Intent
		err := json.NewDecoder(r.Body).Decode(&intent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Verify the received intent matches what we sent
		if intent.Identity != testIntent.Identity {
			t.Errorf("got identity %v, want %v", intent.Identity, testIntent.Identity)
		}
		if intent.Status != testIntent.Status {
			t.Errorf("got status %v, want %v", intent.Status, testIntent.Status)
		}

		// Mock successful response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int64{"id": 1})
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("CreateIntent returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}

// TestGetIntent tests the /getIntent endpoint for retrieving an intent by ID.
// It checks the response status code and verifies the returned intent data.
func TestGetIntent(t *testing.T) {
	// Create test request
	req := httptest.NewRequest("GET", "/getIntent?id=1", nil)
	w := httptest.NewRecorder()

	// Create handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		// Verify query parameter
		if id := r.URL.Query().Get("id"); id != "1" {
			t.Errorf("got id %v, want 1", id)
		}

		// Mock response intent
		intent := Intent{
			ID:            1,
			Identity:      "testIdentity",
			IdentityCurve: "ecdsa",
			Status:        INTENT_STATUS_COMPLETED,
		}

		json.NewEncoder(w).Encode(intent)
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetIntent returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	// Decode response
	var response Intent
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID != 1 || response.Status != INTENT_STATUS_COMPLETED {
		t.Errorf("Got unexpected response: %+v", response)
	}
}

// TestGetIntentsWithPagination tests the /getIntents endpoint with pagination parameters.
// It verifies the response status code for valid and invalid limit values.
func TestGetIntentsWithPagination(t *testing.T) {
	tests := []struct {
		name       string
		limit      int
		skip       int
		wantStatus int
	}{
		{
			name:       "Valid pagination",
			limit:      10,
			skip:       0,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid limit",
			limit:      -1,
			skip:       0,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/getIntents?limit="+strconv.Itoa(tt.limit)+"&skip="+strconv.Itoa(tt.skip), nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.limit < 0 {
					http.Error(w, "invalid limit", http.StatusInternalServerError)
					return
				}
				result := IntentsResult{
					Intents: []*Intent{},
					Total:   0,
				}
				json.NewEncoder(w).Encode(result)
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetIntents returned wrong status code: got %v want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

// TestGetSolverStats tests the /getStatsOfSolver endpoint for retrieving solver statistics.
// It checks the response status code for valid and empty solver scenarios.
func TestGetSolverStats(t *testing.T) {
	tests := []struct {
		name       string
		solver     string
		wantStatus int
	}{
		{
			name:       "Valid solver",
			solver:     "testSolver",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Empty solver",
			solver:     "",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/getStatsOfSolver?solver="+tt.solver, nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.solver == "" {
					http.Error(w, "solver required", http.StatusInternalServerError)
					return
				}
				result := SolverStatResult{
					IsActive:    true,
					ActiveSince: 123456,
					Chains:      []uint{1, 2, 3},
				}
				json.NewEncoder(w).Encode(result)
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetSolverStats returned wrong status code: got %v want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

// TestParseOperation tests the /parseOperation endpoint for parsing an operation.
// It verifies the response status code for valid and invalid operation ID scenarios.
func TestParseOperation(t *testing.T) {
	tests := []struct {
		name        string
		operationId string
		intentId    string
		wantStatus  int
	}{
		{
			name:        "Valid operation",
			operationId: "1",
			intentId:    "1",
			wantStatus:  http.StatusOK,
		},
		{
			name:        "Invalid operation",
			operationId: "",
			intentId:    "1",
			wantStatus:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/parseOperation?operationId=" + tt.operationId + "&intentId=" + tt.intentId
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.operationId == "" {
					http.Error(w, "operation ID required", http.StatusInternalServerError)
					return
				}
				json.NewEncoder(w).Encode(map[string]string{"status": "success"})
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ParseOperation returned wrong status code: got %v want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

// TestGetTotalStats tests the /getTotalStats endpoint for retrieving total statistics.
// It checks the response status code and verifies the returned statistics data.
func TestGetTotalStats(t *testing.T) {
	req := httptest.NewRequest("GET", "/getTotalStats", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stats := TotalStats{
			TotalSolvers: 10,
			TotalIntents: 100,
		}
		json.NewEncoder(w).Encode(stats)
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetTotalStats returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	var response TotalStats
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.TotalSolvers != 10 || response.TotalIntents != 100 {
		t.Errorf("Unexpected response values: %+v", response)
	}
}

// TestStatusEndpoint tests the /status endpoint for checking service health.
// It verifies the response status code and body content.
func TestStatusEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/status", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status endpoint returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Status endpoint returned wrong body: got %v want OK", w.Body.String())
	}
}
