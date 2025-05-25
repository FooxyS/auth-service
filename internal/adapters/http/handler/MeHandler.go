package handler

import (
	"log"
	"net/http"

	"github.com/FooxyS/auth-service/internal/usecase"
)

type MeHandler struct {
	UseCase usecase.MeUseCase
}

func (h *MeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	access := r.Header.Get("Authorization")

	user, errMeExecute := h.UseCase.Execute(r.Context(), access)
	if errMeExecute != nil {
		log.Printf("Execute error: %v", errMeExecute)
		WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	WriteMeJSON(w, http.StatusOK, user)
}
