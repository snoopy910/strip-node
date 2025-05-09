// Package uniswap_v3_solver provides HTTP endpoints for Uniswap V3 operations
package uniswap_v3_solver

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Server struct {
	solver *UniswapV3Solver
}

func NewServer(solver *UniswapV3Solver) *Server {
	return &Server{
		solver: solver,
	}
}

func startHTTPServer(solver *UniswapV3Solver, port string) {
	server := NewServer(solver)

	http.HandleFunc("/health", server.handleHealth)
	http.HandleFunc("/construct", server.handleConstruct)
	http.HandleFunc("/solve", server.handleSolve)
	http.HandleFunc("/status", server.handleStatus)
	http.HandleFunc("/output", server.handleOutput)

	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		panic(err)
	}
}

func (s *Server) handleConstruct(w http.ResponseWriter, r *http.Request) {
	var intent Intent
	if err := json.NewDecoder(r.Body).Decode(&intent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	opIndex, _ := strconv.Atoi(r.URL.Query().Get("operationIndex"))

	dataToSign, err := s.solver.Construct(intent, opIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"dataToSign": dataToSign,
	})
}

func (s *Server) handleSolve(w http.ResponseWriter, r *http.Request) {
	var intent Intent
	if err := json.NewDecoder(r.Body).Decode(&intent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	opIndex, _ := strconv.Atoi(r.URL.Query().Get("operationIndex"))
	signature := r.URL.Query().Get("signature")

	result, err := s.solver.Solve(intent, opIndex, signature)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"result": result,
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	var intent Intent
	if err := json.NewDecoder(r.Body).Decode(&intent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	opIndex, _ := strconv.Atoi(r.URL.Query().Get("operationIndex"))

	status, err := s.solver.Status(intent, opIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": status,
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleOutput(w http.ResponseWriter, r *http.Request) {
	var intent Intent
	if err := json.NewDecoder(r.Body).Decode(&intent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	opIndex, _ := strconv.Atoi(r.URL.Query().Get("operationIndex"))

	output, err := s.solver.GetOutput(intent, opIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(output))
}
