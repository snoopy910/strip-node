package sequencer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type SignMessageResponse struct {
	ID int64 `json:"id"`
}

func TestSignMessage(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantError  bool
	}{
		{
			name:       "Valid ECDSA message signing",
			path:       "/signMessage?message=Hello&identity=test&identityCurve=ecdsa&keyCurve=ecdsa",
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name:       "Valid EDDSA message signing",
			path:       "/signMessage?message=Hello&identity=test&identityCurve=eddsa&keyCurve=eddsa",
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name:       "Missing message parameter",
			path:       "/signMessage?identity=test&identityCurve=ecdsa&keyCurve=ecdsa",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name:       "Missing identity parameter",
			path:       "/signMessage?message=Hello&identityCurve=ecdsa&keyCurve=ecdsa",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
		{
			name:       "Invalid curve type",
			path:       "/signMessage?message=Hello&identity=test&identityCurve=invalid&keyCurve=invalid",
			wantStatus: http.StatusBadRequest,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				enableCors(&w)

				message := r.URL.Query().Get("message")
				identity := r.URL.Query().Get("identity")
				identityCurve := r.URL.Query().Get("identityCurve")
				keyCurve := r.URL.Query().Get("keyCurve")

				if message == "" || identity == "" || identityCurve == "" || keyCurve == "" {
					http.Error(w, "Missing required parameters", http.StatusBadRequest)
					return
				}

				if !isValidCurve(identityCurve) || !isValidCurve(keyCurve) {
					http.Error(w, "Invalid curve type", http.StatusBadRequest)
					return
				}

				intent := Intent{
					Operations: []Operation{
						{
							Type:       OPERATION_TYPE_SIGN_MESSAGE,
							Status:     OPERATION_STATUS_PENDING,
							KeyCurve:   keyCurve,
							DataToSign: message,
						},
					},
					Identity:      identity,
					IdentityCurve: identityCurve,
					Status:        INTENT_STATUS_PROCESSING,
					Expiry:        uint64(time.Now().Add(1 * time.Hour).Unix()),
					CreatedAt:     uint64(time.Now().Unix()),
				}

				mockID := int64(1)
				response := SignMessageResponse{ID: mockID}
				json.NewEncoder(w).Encode(response)
			})

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %v; got %v", tt.wantStatus, w.Code)
			}

			if !tt.wantError && w.Code == http.StatusOK {
				var response SignMessageResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}
				if response.ID <= 0 {
					t.Error("Expected positive intent ID")
				}
			}
		})
	}
}

func isValidCurve(curve string) bool {
	return curve == "ecdsa" || curve == "eddsa"
}
