package httpx

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, payload any, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, msg string, code int) error {
	return WriteJSON(w, Response{StatusError, code, msg}, code)
}

func WriteSuccess(w http.ResponseWriter, msg string, code int) error {
	return WriteJSON(w, Response{StatusSuccess, code, msg}, code)
}
