package sequencer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Operation struct {
	SerializedTxn string `json:"serializedTxn"`
	DataToSign    string `json:"dataToSign"`
	ChainId       string `json:"chainId"`
	KeyCurve      string `json:"keyCurve"`
	Status        string `json:"status"`
	TxnHash       string `json:"txnHash"`
}

type Intent struct {
	Operations    []Operation `json:"operations"`
	Signature     string      `json:"signature"`
	Identity      string      `json:"identity"`
	IdentityCurve string      `json:"identityCurve"`
	Status        string      `json:"status"`
}

func startHTTPServer(port string) {
	http.HandleFunc("/createIntent", func(w http.ResponseWriter, r *http.Request) {
		var intent Intent

		err := json.NewDecoder(r.Body).Decode(&intent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := AddIntent(&intent)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "{\"id\": %d}", id)
	})

	http.HandleFunc("/getIntent", func(w http.ResponseWriter, r *http.Request) {
		intentId := r.URL.Query().Get("id")
		i, _ := strconv.ParseInt(intentId, 10, 64)

		intent, err := GetIntent(i)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(intent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/getIntents", func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")

		intents, err := GetIntents(status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(intents)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
