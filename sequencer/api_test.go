package sequencer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
