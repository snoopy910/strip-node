package signer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/mr-tron/base58"
)

var messageChan = make(map[string]chan (Message))

func generateKeygenMessage(networkId string) {
	message := Message{
		Type:      MESSAGE_TYPE_GENERATE_START_KEYGEN,
		NetworkId: networkId,
	}

	broadcast(message)
}

func generateSignatureMessage(networkId string, msg string) {
	message := Message{
		Type:      MESSAGE_TYPE_START_SIGN,
		Hash:      []byte(msg),
		NetworkId: networkId,
	}

	broadcast(message)
}

func startHTTPServer(port string) {
	http.HandleFunc("/keygen", func(w http.ResponseWriter, r *http.Request) {
		networkId := r.URL.Query().Get("networkId")
		go generateKeygenMessage(networkId)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/address", func(w http.ResponseWriter, r *http.Request) {
		networkId := r.URL.Query().Get("networkId")

		if networks[networkId].Key == nil {
			return
		}

		pk := edwards.PublicKey{
			Curve: tss.Edwards(),
			X:     networks[networkId].Key.EDDSAPub.X(),
			Y:     networks[networkId].Key.EDDSAPub.Y(),
		}

		publicKeyStr := base58.Encode(pk.Serialize())

		fmt.Fprintf(w, "%s", publicKeyStr)
	})

	http.HandleFunc("/signature", func(w http.ResponseWriter, r *http.Request) {
		hash := r.URL.Query().Get("hash")
		networkId := r.URL.Query().Get("networkId")
		go generateSignatureMessage(networkId, hash)

		messageChan[hash] = make(chan Message)

		sig := <-messageChan[hash]
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode("{\"signature\":\"" + string(sig.Message) + "\",\"address\":\"" + sig.Address + "\"}")
		if err != nil {
			http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
		}

		delete(messageChan, hash)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
