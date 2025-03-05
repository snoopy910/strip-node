package signer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	identityVerification "github.com/StripChain/strip-node/identity"
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

// add test for generateKeygenMessage

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
			name:          "Invalid identity",
			identity:      "",
			identityCurve: "ecdsa",
			keyCurve:      "ecdsa",
			statusCode:    http.StatusBadRequest,
		},
		{
			name:          "Invalid identityCurve",
			identity:      "0x7d2e55B99cA3bd06977fB499d70B896A90b2A19e",
			identityCurve: "",
			keyCurve:      "ecdsa",
			statusCode:    http.StatusBadRequest,
		},
		{
			name:          "Invalid keyCurve",
			identity:      "0x7d2e55B99cA3bd06977fB499d70B896A90b2A19e",
			identityCurve: "ecdsa",
			keyCurve:      "",
			statusCode:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/address?identity="+tt.identity+"&identityCurve="+tt.identityCurve+"&keyCurve="+tt.keyCurve, nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if tt.identity == "" || tt.identityCurve == "" || tt.keyCurve == "" {
						http.Error(w, "some parameters are missing", http.StatusBadRequest)
						return
					}
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

	operationA := sequencer.Operation{
		SerializedTxn:  "eb018477359400825208941e79929f2a49c27f820340d83b9d0dec35b2137788016345785d8a000080808080",
		DataToSign:     "3225907716b53ed69dc89682cb8633ee2c9be5c66698f27064ff30b367b86e5c",
		ChainId:        "1337",
		GenesisHash:    "",
		KeyCurve:       "ecdsa",
		Type:           sequencer.OPERATION_TYPE_TRANSACTION,
		Solver:         "",
		SolverMetadata: "",
	}

	operationB := sequencer.Operation{
		SerializedTxn: "eb018477359400825208941e79929f2a49c27f820340d83b9d0dec35b2137788016345785d8a000080808080",
		DataToSign:    "2a007c39f2eb743d9afde9e7043d4a22bd262f324fce1c4afc0eb200483eb959",
		ChainId:       "1337",
		GenesisHash:   "",
		KeyCurve:      "ecdsa",
		Status:        sequencer.OPERATION_STATUS_PENDING,
		Type:          sequencer.OPERATION_TYPE_TRANSACTION,
	}

	intentA := sequencer.Intent{
		Operations:    []sequencer.Operation{operationA},
		Identity:      "0xD99eb497608046d3C97B30E62b872daADF6f7dCF",
		IdentityCurve: "ecdsa",
		Signature:     "0x813470a587a320d5b55871944685d61a625355d4aac34117756c062374ad51ae05313cccc0c5b27ffa47392e2a4883cc7c6ef03b509bb232ce944c01dacd17d01b",
		Expiry:        uint64(1741178924),
	}

	intentB := sequencer.Intent{
		ID:            1,
		Operations:    []sequencer.Operation{operationB},
		Signature:     "0x334de74b38aee47bc8e6adc006a6d03cd4a6612b566520027f07423f7dce0f3d52f64cc6fef4ea439f03d1712d802bffbad334d29cbf03c482ac0c224408c5461c",
		Identity:      "0xD99eb497608046d3C97B30E62b872daADF6f7dCF",
		IdentityCurve: "ecdsa",
		Status:        "valid",
		Expiry:        uint64(47547547848),
		CreatedAt:     uint64(0),
	}

	intentC := InvalidIntent{
		Data: []byte{1},
	}

	tests := []struct {
		name       string
		body       interface{}
		opIndex    uint
		signers    []string
		statusCode int
	}{
		{
			name:       "Valid Intent",
			body:       intentA,
			opIndex:    0,
			statusCode: http.StatusOK,
		},
		{
			name:       "Invalid Signature Intent",
			body:       intentB,
			opIndex:    0,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Invalid Intent Object",
			body:       intentC,
			opIndex:    0,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intent, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatalf("Failed to marshal test wallet: %v", err)
			}
			req := httptest.NewRequest("GET", "/signature", bytes.NewBuffer(intent))
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					var intent sequencer.Intent

					err := json.NewDecoder(r.Body).Decode(&intent)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					if reflect.DeepEqual(intent, sequencer.Intent{}) {
						http.Error(w, "Invalid intent", http.StatusBadRequest)
						return
					}

					operationIndex, _ := strconv.Atoi(r.URL.Query().Get("operationIndex"))
					operationIndexInt := uint(operationIndex)

					if intent.Expiry < uint64(time.Now().Unix()) {
						http.Error(w, "Intent has expired", http.StatusBadRequest)
						return
					}

					msg := intent.Operations[operationIndexInt].DataToSign

					identity := intent.Identity
					identityCurve := intent.IdentityCurve
					keyCurve := intent.Operations[operationIndexInt].KeyCurve

					// verify signature
					intentStr, err := identityVerification.SanitiseIntent(intent)
					if err != nil {
						http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
						return
					}

					verified, err := identityVerification.VerifySignature(
						intent.Identity,
						intent.IdentityCurve,
						intentStr,
						intent.Signature,
					)

					if err != nil {
						http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
						return
					}

					if !verified {
						http.Error(w, "signature verification failed", http.StatusBadRequest)
						return
					}

					go generateSignatureMessage(identity, identityCurve, keyCurve, []byte(msg))

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
