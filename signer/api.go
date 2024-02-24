package signer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mr-tron/base58"
)

var messageChan = make(map[string]chan (Message))

var (
	ECDSA_CURVE = "ecdsa"
	EDDSA_CURVE = "eddsa"
)

func generateKeygenMessage(identity string, identityCurve string, keyCurve string) {
	message := Message{
		Type:          MESSAGE_TYPE_GENERATE_START_KEYGEN,
		Identity:      identity,
		IdentityCurve: identityCurve,
		KeyCurve:      keyCurve,
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

func startHTTPServer(port string) {
	http.HandleFunc("/keygen", func(w http.ResponseWriter, r *http.Request) {
		identity := r.URL.Query().Get("identity")
		identityCurve := r.URL.Query().Get("identityCurve")
		keyCurve := r.URL.Query().Get("keyCurve")
		go generateKeygenMessage(identity, identityCurve, keyCurve)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/address", func(w http.ResponseWriter, r *http.Request) {
		// networkId := r.URL.Query().Get("networkId")

		// if networks[networkId].Key == nil {
		// 	return
		// }

		// pk := edwards.PublicKey{
		// 	Curve: tss.Edwards(),
		// 	X:     networks[networkId].Key.EDDSAPub.X(),
		// 	Y:     networks[networkId].Key.EDDSAPub.Y(),
		// }

		// publicKeyStr := base58.Encode(pk.Serialize())

		// fmt.Fprintf(w, "%s", publicKeyStr)
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
