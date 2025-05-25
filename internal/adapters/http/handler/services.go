package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FooxyS/auth-service/internal/domain"
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

func WriteLoginJSON(w http.ResponseWriter, code int, access string) {
	resp := TokenResponse{
		Access: access,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("WriteLoginJSON error: %v", err)
	}
}

func WriteMeJSON(w http.ResponseWriter, code int, user domain.User) {
	resp := MeResponse{
		UserID: user.UserID,
		Email:  user.Email,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("WriteLoginJSON error: %v", err)
	}
}
