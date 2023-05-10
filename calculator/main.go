package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type handleFunc func(http.ResponseWriter, *http.Request) error

type Expense struct {
	Type  string
	value float32
}

func main() {

	http.HandleFunc("/", wrapHandlers(handleGetCalculation))

	log.Println("Server is up and running on port: 4000")
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func wrapHandlers(fn handleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			_ = writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}

func handleGetCalculation(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return writeJSON(w, http.StatusBadRequest, fmt.Errorf("invalide HTTP method %s", r.Method))
	}
	var resp map[string]float32
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var total float32
	for _, v := range resp {
		total += v
	}

	resp["Total"] = total
	return writeJSON(w, http.StatusOK, resp)

}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
