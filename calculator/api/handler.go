package api

import (
	"encoding/json"
	"fmt"
	. "net/http"
)

type HandleCalculatorFunc func(ResponseWriter, *Request) error

func NewHandleCalculator() Handler {
	return wrapHandler(handlePostCalculation)
}

func handlePostCalculation(w ResponseWriter, r *Request) error {
	if r.Method != "POST" {
		return WriteJSON(w, StatusBadRequest, fmt.Errorf("invalid HTTP method %s", r.Method))
	}
	var resp map[string]float32
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return WriteJSON(w, StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var total float32
	for _, v := range resp {
		total += v
	}

	resp["Total"] = total
	return WriteJSON(w, StatusOK, resp)
}

func wrapHandler(fn HandleCalculatorFunc) Handler {
	return HandlerFunc(func(w ResponseWriter, r *Request) {
		if err := fn(w, r); err != nil {
			_ = WriteJSON(w, StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	})
}
