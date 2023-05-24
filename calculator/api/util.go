package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}

func WrapHandlers(fn HandleCalculatorFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			log.Println(err)

			if apiErr, ok := err.(Error); ok {
				_ = WriteJSON(w, apiErr.Code, apiErr.Message)
				return
			}
			_ = WriteJSON(w, http.StatusInternalServerError, err)
		}
	}
}
