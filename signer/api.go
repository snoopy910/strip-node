package signer

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Silent-Protocol/go-sio/db"
	ecdsaKeygen "github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/mr-tron/base58"
)

var messageChan = make(map[string]chan (Message))

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

		go generateKeygenMessage(createWallet.Identity, createWallet.IdentityCurve, createWallet.KeyCurve, createWallet.Signers)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/address", func(w http.ResponseWriter, r *http.Request) {
		identity := r.URL.Query().Get("identity")
		identityCurve := r.URL.Query().Get("identityCurve")
		keyCurve := r.URL.Query().Get("keyCurve")

		keyShare, err := db.GetKeyShare(identity, identityCurve, keyCurve)

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

			err := json.NewEncoder(w).Encode("{\"address\":\"" + publicKeyStr + "\"}")
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

			err := json.NewEncoder(w).Encode("{\"address\":\"" + address + "\"}")
			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
			}
		}

	})

	http.HandleFunc("/signature", func(w http.ResponseWriter, r *http.Request) {
		msg := r.URL.Query().Get("message")
		identity := r.URL.Query().Get("identity")
		identityCurve := r.URL.Query().Get("identityCurve")
		keyCurve := r.URL.Query().Get("keyCurve")

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
		err := json.NewEncoder(w).Encode("{\"signature\":\"" + base58.Encode(sig.Message) + "\",\"address\":\"" + sig.Address + "\"}")
		if err != nil {
			http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		}

		delete(messageChan, msg)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
