package httph

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type httpCoder interface {
	error
	HTTPStatus() int
}

func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("httph: failed to encode response: %v", err)
	}
}

func SendEmpty(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func sendError(w http.ResponseWriter, status int, encErr error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": encErr.Error()}); err != nil {
		log.Printf("httph: failed to encode response: %v", err)
	}
}

func HandleError(w http.ResponseWriter, err error) {
	var hc httpCoder
	if errors.As(err, &hc) {
		sendError(w, hc.HTTPStatus(), hc)
		return
	}
	sendError(w, http.StatusInternalServerError, err)
}
