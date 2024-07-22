package sequencer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	solversRegistry "github.com/StripChain/strip-node/solversRegistry"
)

type Operation struct {
	ID               int64  `json:"id"`
	SerializedTxn    string `json:"serializedTxn"`
	DataToSign       string `json:"dataToSign"`
	ChainId          string `json:"chainId"`
	KeyCurve         string `json:"keyCurve"`
	Status           string `json:"status"`
	Result           string `json:"result"`
	Type             string `json:"type"`
	Solver           string `json:"solver"`
	SolverMetadata   string `json:"solverMetadata"`
	SolverDataToSign string `json:"solverDataToSign"`
	SolverOutput     string `json:"solverOutput"`
}

type Intent struct {
	ID            int64       `json:"id"`
	Operations    []Operation `json:"operations"`
	Signature     string      `json:"signature"`
	Identity      string      `json:"identity"`
	IdentityCurve string      `json:"identityCurve"`
	Status        string      `json:"status"`
	Expiry        uint64      `json:"expiry"`
	CreatedAt     uint64      `json:"createdAt"`
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
	OPERATION_TYPE_TRANSACTION    = "transaction"
	OPERATION_TYPE_SOLVER         = "solver"
	OPERATION_TYPE_BRIDGE_DEPOSIT = "bridgeDeposit"
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

type IntentsResult struct {
	Intents []*Intent `json:"intents"`
	Total   int       `json:"total"`
}

type SolverStatResult struct {
	IsActive    bool   `json:"isActive"`
	ActiveSince uint   `json:"activeSince"`
	Chains      []uint `json:"chains"`
}

type TotalStats struct {
	TotalSolvers uint `json:"totalSolvers"`
	TotalIntents uint `json:"totalIntents"`
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

		// then store the wallet and it's list of signers in the db
		identity := r.URL.Query().Get("identity")
		identityCurve := r.URL.Query().Get("identityCurve")

		_createWallet := false

		_, err := GetWallet(identity, identityCurve)
		if err != nil {

			if err.Error() == "pg: no rows in result set" {
				_createWallet = true
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if !_createWallet {
			fmt.Fprintf(w, "wallet already exists")
			return
		}

		err = createWallet(identity, identityCurve)
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

	http.HandleFunc("/getBridgeAddress", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		wallet, err := GetWallet(BridgeContractAddress, "ecdsa")
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

		intents, count, err := GetIntentsWithPagination(l, s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		intentsResult := IntentsResult{
			Intents: intents,
			Total:   count,
		}

		err = json.NewEncoder(w).Encode(intentsResult)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/getSolverIntents", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		limit := r.URL.Query().Get("limit")
		skip := r.URL.Query().Get("skip")
		solver := r.URL.Query().Get("solver")

		l, _ := strconv.Atoi(limit)
		s, _ := strconv.Atoi(skip)

		intents, count, err := GetSolverIntents(solver, l, s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		intentsResult := IntentsResult{
			Intents: intents,
			Total:   count,
		}

		err = json.NewEncoder(w).Encode(intentsResult)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/getIntentsOfAddress", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		limit := r.URL.Query().Get("limit")
		skip := r.URL.Query().Get("skip")
		address := r.URL.Query().Get("address")

		l, _ := strconv.Atoi(limit)
		s, _ := strconv.Atoi(skip)

		intents, count, err := GetIntentsOfAddress(address, l, s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		intentsResult := IntentsResult{
			Intents: intents,
			Total:   count,
		}

		err = json.NewEncoder(w).Encode(intentsResult)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/getStatsOfSolver", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		solver := r.URL.Query().Get("solver")

		isActive, activeSince, chains, err := solversRegistry.Stats(
			RPC_URL,
			SolversRegistryContractAddress,
			solver,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		solverStatResult := SolverStatResult{
			IsActive:    isActive,
			ActiveSince: activeSince,
			Chains:      chains,
		}

		err = json.NewEncoder(w).Encode(solverStatResult)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/getTotalStats", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		totalSolvers, err := solversRegistry.TotalSolvers(
			RPC_URL,
			SolversRegistryContractAddress,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}

		totalIntents, err := GetTotalIntents()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		totalStats := TotalStats{
			TotalSolvers: totalSolvers,
			TotalIntents: uint(totalIntents),
		}

		err = json.NewEncoder(w).Encode(totalStats)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/parseOperation", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		operationId := r.URL.Query().Get("operationId")
		intentId := r.URL.Query().Get("intentId")
		i, _ := strconv.ParseInt(operationId, 10, 64)
		j, _ := strconv.ParseInt(intentId, 10, 64)

		operation, err := GetOperation(j, i)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		intent, err := GetIntent(j)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wallet, err := GetWallet(intent.Identity, intent.IdentityCurve)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if operation.Result == "" || operation.Status != OPERATION_STATUS_COMPLETED {
			http.Error(w, "operation not completed", http.StatusInternalServerError)
			return
		}

		if operation.KeyCurve == "ecdsa" {
			transfers, err := GetEthereumTransfers(operation.ChainId, operation.Result, wallet.ECDSAPublicKey)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(transfers)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			transfers, err := GetSolanaTransfers(operation.ChainId, operation.Result, HeliusApiKey)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(transfers)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		fmt.Fprintf(w, "OK")
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
