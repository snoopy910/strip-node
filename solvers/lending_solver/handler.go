package lending_solver

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Server struct {
	solver *LendingSolver
}

func NewServer(rpcURL string, chainId int64, lendingPool, stripUSD string) (*Server, error) {
	solver, err := NewLendingSolver(rpcURL, chainId, lendingPool, stripUSD)
	if err != nil {
		return nil, err
	}
	return &Server{solver: solver}, nil
}

func (s *Server) Serve() error {
	http.HandleFunc("/construct", s.handleConstruct)
	http.HandleFunc("/solve", s.handleSolve)
	http.HandleFunc("/status", s.handleStatus)
	http.HandleFunc("/output", s.handleOutput)

	return http.ListenAndServe(":8080", nil)
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

	json.NewEncoder(w).Encode(map[string]string{
		"output": output,
	})
}
