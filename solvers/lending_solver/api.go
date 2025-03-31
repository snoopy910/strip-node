// Package lending_solver provides HTTP endpoints for lending operations
package lending_solver

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Server struct {
	solver *LendingSolver
}

func NewServer(solver *LendingSolver) *Server {
	return &Server{
		solver: solver,
	}
}

func startHTTPServer(solver *LendingSolver, port string) {
	server := NewServer(solver)

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
