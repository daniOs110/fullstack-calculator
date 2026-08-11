package httpapi

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

// writeResult guards the two float64 values that must never reach the client:
// NaN and infinities, which encoding/json refuses to marshal, and negative
// zero, which is valid IEEE-754 but would render as "-0".
func writeResult(w http.ResponseWriter, result float64) {
	if math.IsNaN(result) || math.IsInf(result, 0) {
		writeError(w, http.StatusBadRequest, "the result is not a finite number")
		return
	}

	if result == 0 {
		result = 0
	}

	writeJSON(w, http.StatusOK, resultResponse{Result: result})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("httpapi: encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
