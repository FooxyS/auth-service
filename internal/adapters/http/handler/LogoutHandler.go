package handler

import (
	"log"
	"net/http"

	"github.com/FooxyS/auth-service/internal/usecase"
)

type LogoutHandler struct {
	UseCase usecase.LogoutUseCase
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	access := r.Header.Get("Authorization")

	if errLogoutExecute := h.UseCase.Execute(r.Context(), access); errLogoutExecute != nil {
		log.Printf("Execute error: %v", errLogoutExecute)
		WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	WriteJSON(w, http.StatusOK, "пользователь успешно деавторизован")
}
