package signer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/StripChain/strip-node/sequencer"
)

func TestKeygenEndpoint(t *testing.T) {

	oldNodePublicKey := NodePublicKey
	defer func() { NodePublicKey = oldNodePublicKey }()

	NodePublicKey = "0x04dae88f5367ea7f086f2a680f212034b73f4c26438b174bd093a71d9b78904eb3b20bc874ef4a04a2ac9531ce29cffacdacd347edbc25133dddf07fac07feedca"

	type InvalidCreateWallet struct {
		Data []byte `json:"data"`
	}

	walletA := CreateWallet{
		Identity:      "0x7d2e55B99cA3bd06977fB499d70B896A90b2A19e",
		IdentityCurve: "ecdsa",
		KeyCurve:      "ecdsa",
		Signers:       []string{"0x04dae88f5367ea7f086f2a680f212034b73f4c26438b174bd093a71d9b78904eb3b20bc874ef4a04a2ac9531ce29cffacdacd347edbc25133dddf07fac07feedca", "0x04dae88f5367ea7f086f2a680f212034b73f4c26438b174bd093a71d9b78904eb3b20bc874ef4a04a2ac9531ce29cffacdacd347edbc25133dddf07fac07feedca"},
	}

	walletB := InvalidCreateWallet{
		Data: []byte{1},
	}

	tests := []struct {
		name       string
		body       interface{}
		signers    []string
		statusCode int
	}{
		{
			name:       "Valid keygen",
			body:       walletA,
			statusCode: http.StatusOK,
		},
		{
			name:       "Invalid keygen",
			body:       walletB,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			walletJSON, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatalf("Failed to marshal test wallet: %v", err)
			}
			req := httptest.NewRequest("GET", "/keygen", bytes.NewBuffer(walletJSON))
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					var createWallet CreateWallet

					err := json.NewDecoder(r.Body).Decode(&createWallet)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					if reflect.DeepEqual(createWallet, CreateWallet{}) {
						http.Error(w, "Invalid wallet", http.StatusBadRequest)
						return
					}
					key := createWallet.Identity + "_" + createWallet.IdentityCurve + "_" + createWallet.KeyCurve

					keygenGeneratedChan[key] = make(chan string)

					go generateKeygenMessage(createWallet.Identity, createWallet.IdentityCurve, createWallet.KeyCurve, createWallet.Signers)

					// <-keygenGeneratedChan[key]
					// fmt.Println(v)
					// delete(keygenGeneratedChan, key)
					w.WriteHeader(http.StatusOK)
				}).ServeHTTP(w, r)
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("CreateWallet returned wrong status code: got %v want %v", w.Code, tt.statusCode)
			}
		})
	}
}

func TestAddressEndpoint(t *testing.T) {

	tests := []struct {
		name          string
		identity      string
		identityCurve string
		keyCurve      string
		signers       []string
		statusCode    int
	}{
		{
			name:          "Valid address",
			identity:      "0x7d2e55B99cA3bd06977fB499d70B896A90b2A19e",
			identityCurve: "ecdsa",
			keyCurve:      "ecdsa",
			statusCode:    http.StatusOK,
		},
		{
			name:          "Invalid address",
			identity:      "",
			identityCurve: "ecdsa",
			keyCurve:      "ecdsa",
			statusCode:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/address?identity="+tt.identity+"&identityCurve="+tt.identityCurve+"&keyCurve="+tt.keyCurve, nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}).ServeHTTP(w, r)
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("Address returned wrong status code: got %v want %v", w.Code, tt.statusCode)
			}
		})
	}
}

func TestSignatureEndpoint(t *testing.T) {

	oldNodePublicKey := NodePublicKey
	defer func() { NodePublicKey = oldNodePublicKey }()

	NodePublicKey = "0x04dae88f5367ea7f086f2a680f212034b73f4c26438b174bd093a71d9b78904eb3b20bc874ef4a04a2ac9531ce29cffacdacd347edbc25133dddf07fac07feedca"

	type InvalidIntent struct {
		Data []byte `json:"data"`
	}

	intentA := sequencer.Intent{
		ID:            1,
		Operations:    []sequencer.Operation{},
		Signature:     "0x04dae88f5367ea7f086f2a680f212034b73f4c26438b174bd093a71d9b78904eb3b20bc874ef4a04a2ac9531ce29cffacdacd347edbc25133dddf07fac07feedca",
		Identity:      "0x7d2e55B99cA3bd06977fB499d70B896A90b2A19e",
		IdentityCurve: "ecdsa",
		Status:        "valid",
		Expiry:        uint64(0),
		CreatedAt:     uint64(0),
	}

	intentB := InvalidIntent{
		Data: []byte{1},
	}

	tests := []struct {
		name       string
		body       interface{}
		signers    []string
		statusCode int
	}{
		{
			name:       "Valid Intent",
			body:       intentA,
			statusCode: http.StatusOK,
		},
		{
			name:       "Invalid Intent",
			body:       intentB,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intentJSON, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatalf("Failed to marshal test wallet: %v", err)
			}
			req := httptest.NewRequest("GET", "/signature", bytes.NewBuffer(intentJSON))
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}).ServeHTTP(w, r)
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("Signature returned wrong status code: got %v want %v", w.Code, tt.statusCode)
			}
		})
	}
}
