package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/FooxyS/auth-service/internal/usecase"
	"github.com/FooxyS/auth-service/pkg/apperrors"
)

type LogoutHandler struct {
	UseCase usecase.LogoutUseCase
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	access := r.Header.Get("Authorization")

	if access == "" {
		WriteJSON(w, http.StatusBadRequest, "Bad Request")
		return
	}

	errLogoutExecute := h.UseCase.Execute(r.Context(), access)
	if errors.Is(errLogoutExecute, apperrors.ErrSessionNotFound) {
		WriteJSON(w, http.StatusBadRequest, "session was not deleted")
		return
	}
	if errLogoutExecute != nil {
		log.Printf("Execute error: %v", errLogoutExecute)
		WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	WriteJSON(w, http.StatusOK, "пользователь успешно деавторизован")
}
