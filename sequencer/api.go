package sequencer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Operation struct {
	ID             int64  `json:"id"`
	SerializedTxn  string `json:"serializedTxn"`
	DataToSign     string `json:"dataToSign"`
	ChainId        string `json:"chainId"`
	KeyCurve       string `json:"keyCurve"`
	Status         string `json:"status"`
	Result         string `json:"result"`
	Type           string `json:"type"`
	Solver         string `json:"solver"`
	SolverMetadata string `json:"solverMetadata"`
}

type Intent struct {
	ID            int64       `json:"id"`
	Operations    []Operation `json:"operations"`
	Signature     string      `json:"signature"`
	Identity      string      `json:"identity"`
	IdentityCurve string      `json:"identityCurve"`
	Status        string      `json:"status"`
	Expiry        uint64      `json:"expiry"`
}

const (
	INTENT_STATUS_PROCESSING = "processing"
	INTENT_STATUS_COMPLETED  = "completed"
	INTENT_STATUS_FAILED     = "failed"
	INTENT_STATUS_EXPIRED    = "expired"
)

const (
	OPERATION_STATUS_PENDING   = "pending"
	OPERATION_STATUS_WAITING   = "waiting"
	OPERATION_STATUS_COMPLETED = "completed"
	OPERATION_STATUS_FAILED    = "failed"
)

const (
	OPERATION_TYPE_TRANSACTION = "transaction"
	OPERATION_TYPE_SOLVER      = "solver"
)

type CreateWalletRequest struct {
	Identity      string   `json:"identity"`
	IdentityCurve string   `json:"identityCurve"`
	KeyCurve      string   `json:"keyCurve"`
	Signers       []string `json:"signers"`
}

type GetAddressResponse struct {
	Address string `json:"address"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func startHTTPServer(port string) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/createWallet", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		// select a list of nodes.
		// If length of selected nodes is more than maximum nodes then use maximum nodes length as signers.
		// If length of selected nodes is less than maximum nodes then use all nodes as signers.

		signers := SignersList()

		if len(signers) > MaximumSigners {
			// select random number of max signers
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(signers), func(i, j int) { signers[i], signers[j] = signers[j], signers[i] })
			signers = signers[:MaximumSigners]
		}

		signersPublicKeyList := make([]string, len(signers))
		for i, signer := range signers {
			signersPublicKeyList[i] = signer.PublicKey
		}

		// then store the wallet and it's list of signers in the db
		identity := r.URL.Query().Get("identity")
		identityCurve := r.URL.Query().Get("identityCurve")

		createWallet := false

		_, err := GetWallet(identity, identityCurve)
		if err != nil {

			if err.Error() == "pg: no rows in result set" {
				createWallet = true
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if !createWallet {
			fmt.Fprintf(w, "wallet already exists")
			return
		}

		// now create the wallet here
		createWalletRequest := CreateWalletRequest{
			Identity:      identity,
			IdentityCurve: identityCurve,
			KeyCurve:      "eddsa",
			Signers:       signersPublicKeyList,
		}
		marshalled, err := json.Marshal(createWalletRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("GET", signers[0].URL+"/keygen", bytes.NewReader(marshalled))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		client := http.Client{Timeout: 3 * time.Minute}
		_, err = client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		createWalletRequest = CreateWalletRequest{
			Identity:      identity,
			IdentityCurve: identityCurve,
			KeyCurve:      "ecdsa",
			Signers:       signersPublicKeyList,
		}
		marshalled, err = json.Marshal(createWalletRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req, err = http.NewRequest("GET", signers[0].URL+"/keygen", bytes.NewReader(marshalled))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		client = http.Client{Timeout: 3 * time.Minute}
		_, err = client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//get the wallets addresses
		resp, err := http.Get(signers[0].URL + "/address?identity=" + identity + "&identityCurve=" + identityCurve + "&keyCurve=eddsa")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// bodyStr := strings.ReplaceAll(string(body), "\"", "")
		// body = []byte(bodyStr[1 : len(bodyStr)-1])

		// fmt.Println(string(body))

		var getAddressResponse GetAddressResponse
		err = json.Unmarshal(body, &getAddressResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// fmt.Println("reached here")

		eddsaAddress := getAddressResponse.Address

		resp, err = http.Get(signers[0].URL + "/address?identity=" + identity + "&identityCurve=" + identityCurve + "&keyCurve=ecdsa")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}

		// fmt.Println(string(body))

		err = json.Unmarshal(body, &getAddressResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ecdsaAddress := getAddressResponse.Address

		wallet := WalletSchema{
			Identity:       r.URL.Query().Get("identity"),
			IdentityCurve:  r.URL.Query().Get("identityCurve"),
			Signers:        strings.Join(signersPublicKeyList, ","),
			EDDSAPublicKey: eddsaAddress,
			ECDSAPublicKey: ecdsaAddress,
		}

		_, err = AddWallet(&wallet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/getWallet", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		identity := r.URL.Query().Get("identity")
		identityCurve := r.URL.Query().Get("identityCurve")

		wallet, err := GetWallet(identity, identityCurve)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(wallet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/createIntent", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

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

		go ProcessIntent(id)

		fmt.Fprintf(w, "{\"id\": %d}", id)
	})

	http.HandleFunc("/getIntent", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

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
		enableCors(&w)

		limit := r.URL.Query().Get("limit")
		skip := r.URL.Query().Get("skip")

		l, _ := strconv.Atoi(limit)
		s, _ := strconv.Atoi(skip)

		intents, err := GetIntentsWithPagination(l, s)
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
		enableCors(&w)

		fmt.Fprintf(w, "OK")
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
