package solver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func StartTestSolver(httpPort string) {
	keepAlive := make(chan string)
	go startHTTPServer(httpPort)

	<-keepAlive
}

func startHTTPServer(port string) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/construct", func(w http.ResponseWriter, r *http.Request) {
		res := ConstructResponse{
			DataToSign: "406cf191c468eb76e34aec5e570c51d975bedd36208f9a80212a09a0f48015e0",
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/solve", func(w http.ResponseWriter, r *http.Request) {
		res := SolveResponse{
			Result: "0x65bd0e037280d8584202ce028f57359e11219df51ba8da4ad4acf5ee5f340420",
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		res := StatusResponse{
			Status: SOLVER_OPERATION_STATUS_SUCCESS,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/output", func(w http.ResponseWriter, r *http.Request) {
		res := OutputResponse{
			Output: "87789783467823",
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
