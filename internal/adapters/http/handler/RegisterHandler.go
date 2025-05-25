package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/FooxyS/auth-service/internal/usecase"
	"github.com/FooxyS/auth-service/pkg/apperrors"
)

type RegisterHandler struct {
	UseCase usecase.RegisterUseCase
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userInfo := new(RegisterRequest)

	if errDecode := json.NewDecoder(r.Body).Decode(userInfo); errDecode != nil {
		log.Printf("NewDecoder error: %v", errDecode)
		WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	errExecute := h.UseCase.Execute(r.Context(), userInfo.Email, userInfo.Password)
	if errors.Is(errExecute, apperrors.ErrUserExists) {
		WriteJSON(w, http.StatusConflict, "User Already Exists")
		return
	}
	if errExecute != nil {
		log.Printf("Execute error: %v", errExecute)
		WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	WriteJSON(w, http.StatusCreated, "User Successfully Created")
}
