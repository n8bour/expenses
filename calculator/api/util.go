package api

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}
