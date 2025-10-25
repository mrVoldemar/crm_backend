package errors

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// ErrorHandler wraps HTTP handlers to provide consistent error handling
func ErrorHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic (in production, use proper logging)
				// log.Printf("Panic occurred: %v", err)

				// Return a generic error response
				WriteError(w, http.StatusInternalServerError, "Internal server error")
			}
		}()

		h(w, r)
	}
}

// WriteError writes a standardized error response
func WriteError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errorResp := ErrorResponse{
		Error: message,
		Code:  code,
	}

	json.NewEncoder(w).Encode(errorResp)
}
