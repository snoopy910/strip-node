package signer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	identityVerification "github.com/StripChain/strip-node/identity"
	"github.com/StripChain/strip-node/sequencer"
	"github.com/ethereum/go-ethereum/crypto"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type ITopic interface {
	Publish(ctx context.Context, data []byte, opts ...pubsub.PubOpt) error
}

type Topic struct {
	p     *pubsub.PubSub
	topic string
}

func (t *Topic) Publish(ctx context.Context, data []byte, opts ...pubsub.PubOpt) error {
	fmt.Println("### Publish fake called")
	return nil
}

var topicMock *Topic

func broadcastWithMockPublish(message Message) {

	messageBytes, err := json.Marshal(message)

	if err != nil {
		fmt.Println("### Marshal error")
		fmt.Println(err)
		panic(err)
	}

	hash := crypto.Keccak256Hash(messageBytes)

	cleanedPrivateKey := strings.Replace(NodePrivateKey, "0x", "", 1)
	privateKey, err := crypto.HexToECDSA(cleanedPrivateKey)
	if err != nil {
		fmt.Println("### HexToECDSA error")
		log.Fatal(err)
	}

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		fmt.Println("### Sign error")
		log.Fatal(err)
	}

	message.Signature = signature

	out, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	ctx := context.Background()
	fmt.Println(ctx)
	if err := topicMock.Publish(ctx, out); err != nil {
		fmt.Println("### Publish error")
		fmt.Println(err)
		panic(err)
	}
}

func TestKeygenEndpoint(t *testing.T) {

	oldNodePublicKey := NodePublicKey
	oldTopic := topic
	defer func() { NodePublicKey = oldNodePublicKey }()
	defer func() { topic = oldTopic }()

	NodePublicKey = "0x04dae88f5367ea7f086f2a680f212034b73f4c26438b174bd093a71d9b78904eb3b20bc874ef4a04a2ac9531ce29cffacdacd347edbc25133dddf07fac07feedca"
	topicMock = &Topic{}

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

func broadcastTest(message interface{}, mockPublish bool) {
	_, err := json.Marshal(message)

	if err != nil {
		fmt.Println("### Marshal error")
		fmt.Println(err)
		panic(err)
	}
	if !mockPublish {
		broadcast(message.(Message))
	} else {
		broadcastWithMockPublish(message.(Message))
	}
}

func TestBroadcastMessage(t *testing.T) {

	oldNodePublicKey := NodePublicKey
	oldNodePrivateKey := NodePrivateKey
	defer func() { NodePublicKey = oldNodePublicKey }()
	defer func() { NodePrivateKey = oldNodePrivateKey }()

	identity := "0x7d2e55B99cA3bd06977fB499d70B896A90b2A19e"
	identityCurve := "ecdsa"
	keyCurve := "ecdsa"
	signers := []string{"0x04dae88f5367ea7f086f2a680f212034b73f4c26438b174bd093a71d9b78904eb3b20bc874ef4a04a2ac9531ce29cffacdacd347edbc25133dddf07fac07feedca"}
	message := Message{
		Type:          MESSAGE_TYPE_GENERATE_START_KEYGEN,
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      keyCurve,
		Signers:       signers,
	}

	wrapped := struct {
		Message
		Extra interface{} `json:"extra"`
	}{
		Message: message,
		Extra:   func() {}, // functions are not JSON marshallable
	}

	tests := []struct {
		name        string
		message     interface{}
		nodePubKey  string
		nodePrivKey string
		signers     []string
		panic       bool
		fatal       bool
		mockPublish bool
		expected    string
	}{
		{
			name:        "Valid data and mock publish",
			message:     message,
			nodePubKey:  "0x04934172634cf8f04e50697b53a6dae3708560c0620137f4fbc638d5d34bc38998915cec982f4f0f3644c68fac09a8190a07bd09491dbdf68eb3e52395aed51dbc",
			nodePrivKey: "0xd4f4347d1d4db7064945267eb8bfbd0145d322ffac9320b5de854e2b54508296",
			panic:       false,
			fatal:       false,
			mockPublish: true,
		},
		{
			name:        "Valid data but panic on publish (topic nil in unit test)",
			message:     message,
			nodePubKey:  "0x04934172634cf8f04e50697b53a6dae3708560c0620137f4fbc638d5d34bc38998915cec982f4f0f3644c68fac09a8190a07bd09491dbdf68eb3e52395aed51dbc",
			nodePrivKey: "0xd4f4347d1d4db7064945267eb8bfbd0145d322ffac9320b5de854e2b54508296",
			panic:       true,
			fatal:       false,
			mockPublish: false,
			expected:    "runtime error: invalid memory address or nil pointer dereference",
		},
		{
			name:        "Invalid node private key",
			message:     message,
			nodePubKey:  "0x04934172634cf8f04e50697b53a6dae3708560c0620137f4fbc638d5d34",
			nodePrivKey: "invalid-key",
			panic:       false,
			fatal:       true,
			mockPublish: false,
			expected:    "invalid hex character 'i' in private key",
		},
		{
			name:        "Invalid message",
			message:     "invalid-message",
			nodePubKey:  "0x04934172634cf8f04e50697b53a6dae3708560c0620137f4fbc638d5d34bc38998915cec982f4f0f3644c68fac09a8190a07bd09491dbdf68eb3e52395aed51dbc",
			nodePrivKey: "0xd4f4347d1d4db7064945267eb8bfbd0145d322ffac9320b5de854e2b54508296",
			panic:       true,
			fatal:       false,
			mockPublish: false,
			expected:    "interface conversion: interface {} is string, not signer.Message",
		},
		{
			name:        "Invalid json message",
			message:     wrapped,
			nodePubKey:  "0x04934172634cf8f04e50697b53a6dae3708560c0620137f4fbc638d5d34bc38998915cec982f4f0f3644c68fac09a8190a07bd09491dbdf68eb3e52395aed51dbc",
			nodePrivKey: "0xd4f4347d1d4db7064945267eb8bfbd0145d322ffac9320b5de854e2b54508296",
			panic:       true,
			fatal:       false,
			mockPublish: false,
			expected:    "json: unsupported type: func()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NodePrivateKey = tt.nodePrivKey
			NodePublicKey = tt.nodePubKey
			if tt.mockPublish {
				broadcastTest(tt.message, true)
			}
			if tt.panic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic, but function did not panic")
					} else {
						expected := tt.expected
						err, _ := r.(error)
						if err.Error() != expected {
							t.Errorf("Expected panic %q, got %q", expected, r)
						}
					}
				}()
				broadcastTest(tt.message, false)
			} else if tt.fatal {
				if os.Getenv("TEST_BROADCAST_FATAL") == "1" {
					broadcastTest(tt.message, false)
				}

				// Otherwise, spawn a subprocess that will execute the fatal path.
				cmd := exec.Command(os.Args[0], "-test.run=TestBroadcastMessage")
				// Pass an environment variable so that the subprocess knows to run the fatal code.
				cmd.Env = append(os.Environ(), "TEST_BROADCAST_FATAL=1")
				// err = cmd.Run()
				stderr, err := cmd.CombinedOutput()
				output := string(stderr)
				fmt.Println(output)

				if err == nil {
					t.Fatalf("Expected broadcast to call log.Fatal and exit with non-zero status, but it did not")
				}
				// If err is an ExitError, check that it indicates a non-zero exit code.
				if exitErr, ok := err.(*exec.ExitError); ok {
					if exitErr.Success() {
						t.Fatalf("Expected a non-zero exit status, but got success")
					}
					// The test passed: the subprocess exited as expected.
				} else {
					t.Fatalf("Expected an ExitError, got: %v", err)
				}
				if !strings.Contains(string(stderr), tt.expected) {
					t.Errorf("Expected output to contain %q, got %q", tt.expected, string(stderr))
				}

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
		SerializedTxn:  "eb808477359400825208945e721f69f4c3c91befeb94b1e068d2e64a82a7f488016345785d8a000080808080",
		DataToSign:     "bc0efa2d6c1a0fcb888e82d400a8273e88ee641b7b615e071dafdc0f4b44c91f",
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
		Signature:     "0xa86c8a070b328afe5d56e340a8b1eebbd58e8948b583718b6138a633e6f066c62c5cc1877afb78d956913d37921015722c5eaf8256bf82da7588442320026b741b",
		Expiry:        uint64(1741196126),
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

					// operationIndex, _ := strconv.Atoi(r.URL.Query().Get("operationIndex"))
					// operationIndexInt := uint(operationIndex)

					// msg := intent.Operations[operationIndexInt].DataToSign

					// identity := intent.Identity
					// identityCurve := intent.IdentityCurve
					// keyCurve := intent.Operations[operationIndexInt].KeyCurve

					// verify signature
					intentStr, err := identityVerification.SanitiseIntent(intent)
					if err != nil {
						http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
						return
					}

					fmt.Println("intentStr: ", intentStr)

					verified, err := identityVerification.VerifySignature(
						intent.Identity,
						intent.IdentityCurve,
						intentStr,
						intent.Signature,
					)

					fmt.Println("verified: ", verified)

					if err != nil {
						http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
						return
					}

					if !verified {
						http.Error(w, "signature verification failed", http.StatusBadRequest)
						return
					}

					// go generateSignatureMessage(identity, identityCurve, keyCurve, []byte(msg))

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
