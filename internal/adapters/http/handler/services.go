package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, code int, message string) {
	resp := MessageResponse{
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("WriteJSON error: %v", err)
	}
}

func WriteLoginJSON(w http.ResponseWriter, code int, access, refresh string) {
	resp := LoginResponse{
		Access:  access,
		Refresh: refresh,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("WriteLoginJSON error: %v", err)
	}
}
