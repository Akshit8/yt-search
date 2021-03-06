package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse represents a response containing an error message.
type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

// generic error handler for REST API
func renderErrorResponse(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	resp := ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Error:      msg,
	}
	status := http.StatusInternalServerError

	log.Println("REST ERROR: ", err.Error())

	renderResponse(w, resp, status)
}

// generic response handler for REST API/
func renderResponse(w http.ResponseWriter, res interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")

	content, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	_, err = w.Write(content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
