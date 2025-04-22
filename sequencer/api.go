package sequencer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/StripChain/strip-node/libs"
	"github.com/StripChain/strip-node/libs/blockchains"
	db "github.com/StripChain/strip-node/libs/database"
	solversRegistry "github.com/StripChain/strip-node/solversRegistry"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/google/uuid"
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
	GET_BLOCKCHAIN_ERROR          = "failed to get blockchain"
	PARSE_UUID_ID_ERROR           = "failed to parse intent id"
)

type GetAddressResponse struct {
	Address string `json:"address"`
}

type GetBitcoinAddressesResponse struct {
	MainnetAddress string `json:"mainnetAddress"`
	TestnetAddress string `json:"testnetAddress"`
	RegtestAddress string `json:"regtestAddress"`
}

type GetDogecoinAddressesResponse struct {
	MainnetAddress string `json:"mainnetAddress"`
	TestnetAddress string `json:"testnetAddress"`
}

type IntentsResult struct {
	Intents []*libs.Intent `json:"intents"`
	Total   int            `json:"total"`
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

type MPCKey struct {
	Title     string `json:"title"`
	Network   string `json:"network"`
	PublicKey string `json:"publicKey"`
}

type WalletResponse struct {
	db.WalletSchema
	MPCKeys []MPCKey `json:"mpcKeys"`
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
	//   - Success: 201 if wallet created
	//   - Error: 400 if identityCurve is invalid or identity is empty
	//   - Error: 409 if wallet already exists
	//   - Error: 500 with error message
	http.HandleFunc("/createWallet", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		// then store the wallet and it's list of signers in the db
		identity := r.URL.Query().Get("identity")
		blockchainIDStr := r.URL.Query().Get("blockchain")

		_createWallet := false

		if identity == "" {
			http.Error(w, "identity required", http.StatusBadRequest)
			return
		}

		blockchainID, err := blockchains.ParseBlockchainID(blockchainIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.GetWallet(identity, blockchainID)
		if err != nil {

			if err.Error() == "pg: no rows in result set" {
				_createWallet = true
			} else {
				http.Error(w, GET_WALLET_ERROR, http.StatusInternalServerError)
				return
			}
		}

		if !_createWallet {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "wallet already exists")
			return
		}

		err = createWallet(identity, blockchainID)
		if err != nil {
			http.Error(w, CREATE_WALLET_ERROR, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	// GetWallet endpoint - Retrieves wallet information
	// Method: GET
	// Query Parameters:
	//   - identity: The unique identifier for the wallet
	//   - identityCurve: The curve type for the identity
	// Response:
	//   - Success: 200 JSON encoded wallet object
	//   - Error: 400 if identityCurve is invalid or identity is empty
	//   - Error: 404 if wallet not found
	//   - Error: 500 with error message
	http.HandleFunc("/getWallet", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		identity := r.URL.Query().Get("identity")
		blockchainIDStr := r.URL.Query().Get("blockchain")
		if identity == "" {
			http.Error(w, "identity required", http.StatusBadRequest)
			return
		}

		blockchainID, err := blockchains.ParseBlockchainID(blockchainIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		wallet, err := db.GetWallet(identity, blockchainID)
		if err != nil {
			if err.Error() == "pg: no rows in result set" {
				http.Error(w, "wallet not found", http.StatusNotFound)
				return
			}
			http.Error(w, GET_WALLET_ERROR, http.StatusInternalServerError)
			return
		}

		// Create response with extended fields
		response := WalletResponse{
			WalletSchema: *wallet,
			MPCKeys:      []MPCKey{}, // Initialize with empty array
		}

		// Helper function to add keys to MPCKeys array
		addKey := func(title, network, publicKey string) {
			if publicKey != "" {
				response.MPCKeys = append(response.MPCKeys, MPCKey{
					Title:     title,
					Network:   network,
					PublicKey: publicKey,
				})
			}
		}

		// Add Bitcoin keys
		addKey("Bitcoin", "mainnet", wallet.BitcoinMainnetPublicKey)
		addKey("Bitcoin", "testnet", wallet.BitcoinTestnetPublicKey)
		addKey("Bitcoin Regtest", "testnet", wallet.BitcoinRegtestPublicKey)

		// Add Dogecoin keys
		addKey("Dogecoin", "mainnet", wallet.DogecoinMainnetPublicKey)
		addKey("Dogecoin", "testnet", wallet.DogecoinTestnetPublicKey)

		// Add other blockchain keys (all default to testnet)
		addKey("Stellar", "testnet", wallet.StellarPublicKey)
		addKey("Ripple", "testnet", wallet.RippleEDDSAPublicKey)
		addKey("Sui", "testnet", wallet.SuiPublicKey)
		addKey("Algorand", "testnet", wallet.AlgorandEDDSAPublicKey)
		addKey("Cardano", "testnet", wallet.CardanoPublicKey)
		addKey("Aptos", "testnet", wallet.AptosEDDSAPublicKey)
		addKey("Ethereum", "testnet", wallet.EthereumPublicKey)
		addKey("Solana", "testnet", wallet.SolanaPublicKey)

		// Add general cryptographic keys
		addKey("EDDSA", "testnet", wallet.EDDSAPublicKey)
		addKey("ECDSA", "testnet", wallet.ECDSAPublicKey)

		if err := json.NewEncoder(w).Encode(response); err != nil {
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

		wallet, err := db.GetWallet(BridgeContractAddress, blockchains.Ethereum)
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

		var intent libs.Intent

		err := json.NewDecoder(r.Body).Decode(&intent)
		if err != nil {
			http.Error(w, DECODE_ERROR, http.StatusBadRequest)
			return
		}

		id, err := db.AddIntent(&intent)
		if err != nil {
			logger.Sugar().Errorw("failed to add intent", "error", err)
			http.Error(w, ADD_INTENT_ERROR, http.StatusInternalServerError)
			return
		}

		go ProcessIntent(id)

		fmt.Fprintf(w, "{\"id\": \"%s\"}", id.String())
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
		id, err := uuid.Parse(intentId)
		if err != nil {
			logger.Sugar().Errorw("failed to parse intent id", "error", err)
			http.Error(w, PARSE_UUID_ID_ERROR, http.StatusInternalServerError)
			return
		}

		intent, err := db.GetIntent(id)
		if err != nil {
			logger.Sugar().Errorw("failed to get intent", "error", err)
			http.Error(w, GET_INTENT_ERROR, http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(intent)
		if err != nil {
			logger.Sugar().Errorw("failed to encode intent", "error", err)
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

		intents, count, err := db.GetIntentsWithPagination(l, s)
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

		intents, count, err := db.GetSolverIntents(solver, l, s)
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

		intents, count, err := db.GetIntentsOfAddress(address, l, s)
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

		totalIntents, err := db.GetTotalIntents()

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
		opID, _ := strconv.ParseInt(operationId, 10, 64)
		intentID, _ := uuid.Parse(intentId)

		operation, err := db.GetOperation(intentID, opID)
		if err != nil {
			http.Error(w, GET_OPERATION_ERROR, http.StatusInternalServerError)
			return
		}

		if operation.Result == "" || operation.Status != libs.OperationStatusCompleted {
			http.Error(w, OPERATION_NOT_COMPLETED_ERROR, http.StatusInternalServerError)
			return
		}

		wallet, err := db.GetWallet(BridgeContractAddress, operation.BlockchainID)
		if err != nil {
			http.Error(w, GET_WALLET_ERROR, http.StatusInternalServerError)
			return
		}

		opBlockchain, err := blockchains.GetBlockchain(operation.BlockchainID, operation.NetworkType)
		if err != nil {
			http.Error(w, GET_BLOCKCHAIN_ERROR, http.StatusInternalServerError)
			return
		}

		transfers, err := opBlockchain.GetTransfers(operation.Result, &wallet.EthereumPublicKey)
		if err != nil {
			http.Error(w, GET_TRANSFERS_ERROR, http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(transfers)
		if err != nil {
			http.Error(w, ENCODE_ERROR, http.StatusInternalServerError)
			return
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
