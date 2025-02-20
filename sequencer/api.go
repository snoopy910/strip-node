package sequencer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/StripChain/strip-node/algorand"
	"github.com/StripChain/strip-node/aptos"
	solversRegistry "github.com/StripChain/strip-node/solversRegistry"
	"github.com/StripChain/strip-node/stellar"
)

type Operation struct {
	ID               int64  `json:"id"`
	SerializedTxn    string `json:"serializedTxn"`
	DataToSign       string `json:"dataToSign"`
	ChainId          string `json:"chainId"`
	GenesisHash      string `json:"genesisHash"`
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
	OPERATION_TYPE_SWAP           = "swap"
	OPERATION_TYPE_BURN           = "burn"
	OPERATION_TYPE_WITHDRAW       = "withdraw"
)

const (
	GET_WALLET_ERROR              = "failed to get wallet"
	CREATE_WALLET_ERROR           = "failed to create wallet"
	ENCODE_ERROR                  = "failed to encode response"
	DECODE_ERROR                  = "request decode error"
	ADD_INTENT_ERROR              = "failed to add intent"
	GET_INTENT_ERROR              = "failed to get intent"
	GET_SOLVER_STATS_ERROR        = "failed to get solver stats"
	GET_OPERATION_ERROR           = "failed to get operation"
	OPERATION_NOT_COMPLETED_ERROR = "operation not completed"
	GET_TRANSFERS_ERROR           = "failed to get transfers"
	PARSE_INT_ERROR               = "int parsing error"
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

type GetBitcoinAddressesResponse struct {
	MainnetAddress string `json:"mainnetAddress"`
	TestnetAddress string `json:"testnetAddress"`
	RegtestAddress string `json:"regtestAddress"`
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
	// Root endpoint - Health check
	// Method: GET
	// Response: Plain text "OK"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		fmt.Fprintf(w, "OK")
	})

	// CreateWallet endpoint - Creates a new wallet for a given identity
	// Method: GET
	// Query Parameters:
	//   - identity: The unique identifier for the wallet
	//   - identityCurve: The curve type for the identity
	// Response:
	//   - Success: "wallet already exists" if wallet exists
	//   - Success: Empty response if wallet created
	//   - Error: 500 with error message
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
				http.Error(w, GET_WALLET_ERROR, http.StatusInternalServerError)
				return
			}
		}

		if !_createWallet {
			fmt.Fprintf(w, "wallet already exists")
			return
		}

		err = createWallet(identity, identityCurve)
		if err != nil {
			http.Error(w, CREATE_WALLET_ERROR, http.StatusInternalServerError)
			return
		}
	})

	// GetWallet endpoint - Retrieves wallet information
	// Method: GET
	// Query Parameters:
	//   - identity: The unique identifier for the wallet
	//   - identityCurve: The curve type for the identity
	// Response:
	//   - Success: JSON encoded wallet object
	//   - Error: 500 with error message
	http.HandleFunc("/getWallet", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		identity := r.URL.Query().Get("identity")
		identityCurve := r.URL.Query().Get("identityCurve")

		wallet, err := GetWallet(identity, identityCurve)
		if err != nil {
			http.Error(w, GET_WALLET_ERROR, http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(wallet)
		if err != nil {
			http.Error(w, ENCODE_ERROR, http.StatusInternalServerError)
			return
		}
	})

	// GetBridgeAddress endpoint - Retrieves the bridge contract wallet address
	// Method: GET
	// Response:
	//   - Success: JSON encoded wallet object for the bridge contract
	//   - Error: 500 with error message
	http.HandleFunc("/getBridgeAddress", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		wallet, err := GetWallet(BridgeContractAddress, "ecdsa")
		if err != nil {
			http.Error(w, GET_WALLET_ERROR, http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(wallet)
		if err != nil {
			http.Error(w, ENCODE_ERROR, http.StatusInternalServerError)
			return
		}
	})

	// CreateIntent endpoint - Creates a new intent for processing
	// Method: POST
	// Body: JSON encoded Intent object
	// Response:
	//   - Success: JSON with created intent ID {"id": <number>}
	//   - Error: 400 for invalid request body
	//   - Error: 500 for processing errors
	// Notes: Triggers async intent processing after creation
	// TODO: Intent validation
	http.HandleFunc("/createIntent", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		var intent Intent

		err := json.NewDecoder(r.Body).Decode(&intent)
		if err != nil {
			http.Error(w, DECODE_ERROR, http.StatusBadRequest)
			return
		}

		id, err := AddIntent(&intent)

		if err != nil {
			http.Error(w, ADD_INTENT_ERROR, http.StatusInternalServerError)
			return
		}

		go ProcessIntent(id)

		fmt.Fprintf(w, "{\"id\": %d}", id)
	})

	// GetIntent endpoint - Retrieves a specific intent by ID
	// Method: GET
	// Query Parameters:
	//   - id: The intent ID to retrieve
	// Response:
	//   - Success: JSON encoded Intent object
	//   - Error: 500 with error message
	http.HandleFunc("/getIntent", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		intentId := r.URL.Query().Get("id")
		i, err := strconv.ParseInt(intentId, 10, 64)
		if err != nil {
			http.Error(w, PARSE_INT_ERROR, http.StatusInternalServerError)
			return
		}

		intent, err := GetIntent(i)
		if err != nil {
			http.Error(w, GET_INTENT_ERROR, http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(intent)
		if err != nil {
			http.Error(w, ENCODE_ERROR, http.StatusInternalServerError)
			return
		}
	})

	// GetIntents endpoint - Retrieves paginated list of intents
	// Method: GET
	// Query Parameters:
	//   - limit: Maximum number of intents to return
	//   - skip: Number of intents to skip (for pagination)
	// Response:
	//   - Success: JSON encoded IntentsResult {intents: Intent[], total: number}
	//   - Error: 500 with error message
	http.HandleFunc("/getIntents", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		limit := r.URL.Query().Get("limit")
		skip := r.URL.Query().Get("skip")

		l, _ := strconv.Atoi(limit)
		s, _ := strconv.Atoi(skip)

		intents, count, err := GetIntentsWithPagination(l, s)
		if err != nil {
			http.Error(w, GET_INTENT_ERROR, http.StatusInternalServerError)
			return
		}

		intentsResult := IntentsResult{
			Intents: intents,
			Total:   count,
		}

		err = json.NewEncoder(w).Encode(intentsResult)
		if err != nil {
			http.Error(w, ENCODE_ERROR, http.StatusInternalServerError)
			return
		}
	})

	// GetSolverIntents endpoint - Retrieves intents for a specific solver
	// Method: GET
	// Query Parameters:
	//   - solver: Address of the solver
	//   - limit: Maximum number of intents to return
	//   - skip: Number of intents to skip
	// Response:
	//   - Success: JSON encoded IntentsResult
	//   - Error: 500 with error message
	http.HandleFunc("/getSolverIntents", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		limit := r.URL.Query().Get("limit")
		skip := r.URL.Query().Get("skip")
		solver := r.URL.Query().Get("solver")

		l, _ := strconv.Atoi(limit)
		s, _ := strconv.Atoi(skip)

		intents, count, err := GetSolverIntents(solver, l, s)
		if err != nil {
			http.Error(w, GET_INTENT_ERROR, http.StatusInternalServerError)
			return
		}

		intentsResult := IntentsResult{
			Intents: intents,
			Total:   count,
		}

		err = json.NewEncoder(w).Encode(intentsResult)
		if err != nil {
			http.Error(w, ENCODE_ERROR, http.StatusInternalServerError)
			return
		}
	})

	// GetIntentsOfAddress endpoint - Retrieves intents for a specific address
	// Method: GET
	// Query Parameters:
	//   - address: The address to query intents for
	//   - limit: Maximum number of intents to return
	//   - skip: Number of intents to skip
	// Response:
	//   - Success: JSON encoded IntentsResult
	//   - Error: 500 with error message
	http.HandleFunc("/getIntentsOfAddress", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		limit := r.URL.Query().Get("limit")
		skip := r.URL.Query().Get("skip")
		address := r.URL.Query().Get("address")

		l, _ := strconv.Atoi(limit)
		s, _ := strconv.Atoi(skip)

		intents, count, err := GetIntentsOfAddress(address, l, s)
		if err != nil {
			http.Error(w, GET_INTENT_ERROR, http.StatusInternalServerError)
			return
		}

		intentsResult := IntentsResult{
			Intents: intents,
			Total:   count,
		}

		err = json.NewEncoder(w).Encode(intentsResult)
		if err != nil {
			http.Error(w, ENCODE_ERROR, http.StatusInternalServerError)
			return
		}
	})

	// GetStatsOfSolver endpoint - Retrieves solver statistics
	// Method: GET
	// Query Parameters:
	//   - solver: Address of the solver
	// Response:
	//   - Success: JSON encoded SolverStatResult {
	//       isActive: boolean,
	//       activeSince: number,
	//       chains: number[]
	//     }
	//   - Error: 500 with error message
	http.HandleFunc("/getStatsOfSolver", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		solver := r.URL.Query().Get("solver")

		isActive, activeSince, chains, err := solversRegistry.Stats(
			RPC_URL,
			SolversRegistryContractAddress,
			solver,
		)

		if err != nil {
			http.Error(w, GET_SOLVER_STATS_ERROR, http.StatusInternalServerError)
			return
		}

		solverStatResult := SolverStatResult{
			IsActive:    isActive,
			ActiveSince: activeSince,
			Chains:      chains,
		}

		err = json.NewEncoder(w).Encode(solverStatResult)
		if err != nil {
			http.Error(w, ENCODE_ERROR, http.StatusInternalServerError)
			return
		}
	})

	// GetTotalStats endpoint - Retrieves global statistics
	// Method: GET
	// Response:
	//   - Success: JSON encoded TotalStats {
	//       totalSolvers: number,
	//       totalIntents: number
	//     }
	//   - Error: 500 with error message
	http.HandleFunc("/getTotalStats", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		totalSolvers, err := solversRegistry.TotalSolvers(
			RPC_URL,
			SolversRegistryContractAddress,
		)

		if err != nil {
			http.Error(w, GET_SOLVER_STATS_ERROR, http.StatusInternalServerError)
			return

		}

		totalIntents, err := GetTotalIntents()

		if err != nil {
			http.Error(w, GET_INTENT_ERROR, http.StatusInternalServerError)
			return
		}

		totalStats := TotalStats{
			TotalSolvers: totalSolvers,
			TotalIntents: uint(totalIntents),
		}

		err = json.NewEncoder(w).Encode(totalStats)
		if err != nil {
			http.Error(w, ENCODE_ERROR, http.StatusInternalServerError)
			return
		}
	})

	// ParseOperation endpoint - Parses and retrieves operation details
	// Method: GET
	// Query Parameters:
	//   - operationId: ID of the operation to parse
	//   - intentId: ID of the parent intent
	// Response:
	//   - Success: JSON encoded transfers data
	//   - Error: 500 if operation not completed or other errors
	// Notes: Supports both Ethereum and Solana transfers
	http.HandleFunc("/parseOperation", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		operationId := r.URL.Query().Get("operationId")
		intentId := r.URL.Query().Get("intentId")
		i, _ := strconv.ParseInt(operationId, 10, 64)
		j, _ := strconv.ParseInt(intentId, 10, 64)

		operation, err := GetOperation(j, i)
		if err != nil {
			http.Error(w, GET_OPERATION_ERROR, http.StatusInternalServerError)
			return
		}

		intent, err := GetIntent(j)
		if err != nil {
			http.Error(w, GET_INTENT_ERROR, http.StatusInternalServerError)
			return
		}

		wallet, err := GetWallet(intent.Identity, intent.IdentityCurve)

		if err != nil {
			http.Error(w, GET_WALLET_ERROR, http.StatusInternalServerError)
			return
		}

		if operation.Result == "" || operation.Status != OPERATION_STATUS_COMPLETED {
			http.Error(w, OPERATION_NOT_COMPLETED_ERROR, http.StatusInternalServerError)
			return
		}

		if operation.KeyCurve == "ecdsa" {
			transfers, err := GetEthereumTransfers(operation.ChainId, operation.Result, wallet.ECDSAPublicKey)

			if err != nil {
				http.Error(w, GET_TRANSFERS_ERROR, http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(transfers)
			if err != nil {
				http.Error(w, ENCODE_ERROR, http.StatusInternalServerError)
				return
			}
		} else if operation.KeyCurve == "eddsa" {
			transfers, err := GetSolanaTransfers(operation.ChainId, operation.Result, HeliusApiKey)

			if err != nil {
				http.Error(w, GET_TRANSFERS_ERROR, http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(transfers)
			if err != nil {
				http.Error(w, ENCODE_ERROR, http.StatusInternalServerError)
				return
			}
		} else if operation.KeyCurve == "aptos_eddsa" {
			transfers, err := aptos.GetAptosTransfers(operation.ChainId, operation.Result)

			if err != nil {
				http.Error(w, GET_TRANSFERS_ERROR, http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(transfers)
			if err != nil {
				http.Error(w, ENCODE_ERROR, http.StatusInternalServerError)
				return
			}
		} else if operation.KeyCurve == "secp256k1" {
			transfers, _, err := GetBitcoinTransfers(operation.ChainId, operation.Result)

			if err != nil {
				http.Error(w, GET_TRANSFERS_ERROR, http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(transfers)
			if err != nil {
				http.Error(w, ENCODE_ERROR, http.StatusInternalServerError)
				return
			}
		} else if operation.KeyCurve == "stellar_eddsa" {
			transfers, err := stellar.GetStellarTransfers(operation.ChainId, operation.Result)

			if err != nil {
				http.Error(w, GET_TRANSFERS_ERROR, http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(transfers)
			if err != nil {
				http.Error(w, ENCODE_ERROR, http.StatusInternalServerError)
				return
			}
		} else if operation.KeyCurve == "algorand_eddsa" {
			transfers, err := algorand.GetAlgorandTransfers(operation.GenesisHash, operation.Result)
			if err != nil {
				http.Error(w, GET_TRANSFERS_ERROR, http.StatusInternalServerError)
				return
			}
			err = json.NewEncoder(w).Encode(transfers)
			if err != nil {
				http.Error(w, ENCODE_ERROR, http.StatusInternalServerError)
				return
			}
		}
	})

	// Status endpoint - Service health check
	// Method: GET
	// Response: Plain text "OK"
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		fmt.Fprintf(w, "OK")
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
