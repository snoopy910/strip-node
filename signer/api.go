package signer

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	identityVerification "github.com/Silent-Protocol/go-sio/identity"
	"github.com/Silent-Protocol/go-sio/sequencer"
	ecdsaKeygen "github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/mr-tron/base58"
)

var messageChan = make(map[string]chan (Message))
var keygenGeneratedChan = make(map[string]chan (string))

var (
	ECDSA_CURVE = "ecdsa"
	EDDSA_CURVE = "eddsa"
)

func generateKeygenMessage(identity string, identityCurve string, keyCurve string, signers []string) {
	message := Message{
		Type:          MESSAGE_TYPE_GENERATE_START_KEYGEN,
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      keyCurve,
		Signers:       signers,
	}

	broadcast(message)
}

func generateSignatureMessage(identity string, identityCurve string, keyCurve string, msg []byte) {
	message := Message{
		Type:          MESSAGE_TYPE_START_SIGN,
		Hash:          msg,
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      keyCurve,
	}

	broadcast(message)
}

type CreateWallet struct {
	Identity      string   `json:"identity"`
	IdentityCurve string   `json:"identityCurve"`
	KeyCurve      string   `json:"keyCurve"`
	Signers       []string `json:"signers"`
}

type SignMessage struct {
	Message       string `json:"message"`
	Identity      string `json:"identity"`
	IdentityCurve string `json:"identityCurve"`
	KeyCurve      string `json:"keyCurve"`
}

func startHTTPServer(port string) {
	http.HandleFunc("/keygen", func(w http.ResponseWriter, r *http.Request) {
		var createWallet CreateWallet

		err := json.NewDecoder(r.Body).Decode(&createWallet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		key := createWallet.Identity + "_" + createWallet.IdentityCurve + "_" + createWallet.KeyCurve

		keygenGeneratedChan[key] = make(chan string)

		go generateKeygenMessage(createWallet.Identity, createWallet.IdentityCurve, createWallet.KeyCurve, createWallet.Signers)

		<-keygenGeneratedChan[key]
		delete(keygenGeneratedChan, key)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/address", func(w http.ResponseWriter, r *http.Request) {
		identity := r.URL.Query().Get("identity")
		identityCurve := r.URL.Query().Get("identityCurve")
		keyCurve := r.URL.Query().Get("keyCurve")

		keyShare, err := GetKeyShare(identity, identityCurve, keyCurve)

		if err != nil {
			http.Error(w, "error from postgres", http.StatusBadRequest)
			return
		}

		if keyShare == "" {
			http.Error(w, "key share not found.", http.StatusBadRequest)
			return
		}

		var rawKeyEddsa *eddsaKeygen.LocalPartySaveData
		var rawKeyEcdsa *ecdsaKeygen.LocalPartySaveData

		if keyCurve == EDDSA_CURVE {
			json.Unmarshal([]byte(keyShare), &rawKeyEddsa)

			pk := edwards.PublicKey{
				Curve: tss.Edwards(),
				X:     rawKeyEddsa.EDDSAPub.X(),
				Y:     rawKeyEddsa.EDDSAPub.Y(),
			}

			publicKeyStr := base58.Encode(pk.Serialize())

			getAddressResponse := GetAddressResponse{
				Address: publicKeyStr,
			}
			err := json.NewEncoder(w).Encode(getAddressResponse)
			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			}
		} else {
			json.Unmarshal([]byte(keyShare), &rawKeyEcdsa)

			x := toHexInt(rawKeyEcdsa.ECDSAPub.X())
			y := toHexInt(rawKeyEcdsa.ECDSAPub.Y())

			publicKeyStr := "04" + x + y
			publicKeyBytes, _ := hex.DecodeString(publicKeyStr)
			address := publicKeyToAddress(publicKeyBytes)

			getAddressResponse := GetAddressResponse{
				Address: address,
			}
			err := json.NewEncoder(w).Encode(getAddressResponse)
			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			}

			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			}
		}
	})

	http.HandleFunc("/signature", func(w http.ResponseWriter, r *http.Request) {
		// the owner of the wallet must have created an intent and signed it.
		// we generate signature for an intent operation

		var intent sequencer.Intent

		err := json.NewDecoder(r.Body).Decode(&intent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		operationIndex, _ := strconv.Atoi(r.URL.Query().Get("operationIndex"))
		operationIndexInt := uint(operationIndex)

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

		if keyCurve == EDDSA_CURVE {
			msgBytes, err := base58.Decode(msg)
			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
				return
			}

			go generateSignatureMessage(identity, identityCurve, keyCurve, msgBytes)
		} else if keyCurve == ECDSA_CURVE {
			go generateSignatureMessage(identity, identityCurve, keyCurve, []byte(msg))
		} else {
			http.Error(w, "invalid key curve", http.StatusBadRequest)
			return
		}

		messageChan[msg] = make(chan Message)

		sig := <-messageChan[msg]

		w.Header().Set("Content-Type", "application/json")

		signatureResponse := SignatureReponse{}

		if keyCurve == ECDSA_CURVE {
			signatureResponse.Signature = string(sig.Message)
			signatureResponse.Address = sig.Address
		} else {
			signatureResponse.Signature = base58.Encode(sig.Message)
			signatureResponse.Address = sig.Address
		}

		err = json.NewEncoder(w).Encode(signatureResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		}

		delete(messageChan, msg)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

type SignatureReponse struct {
	Signature string `json:"signature"`
	Address   string `json:"address"`
}

type GetAddressResponse struct {
	Address string `json:"address"`
}
