package handler

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/FooxyS/auth-service/internal/usecase"
)

type LoginHandler struct {
	UseCase usecase.LoginUseCase
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userInfo := new(LoginRequest)

	if errDecode := json.NewDecoder(r.Body).Decode(userInfo); errDecode != nil {
		log.Printf("NewDecoder error: %v", errDecode)
		WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	host, _, errSplit := net.SplitHostPort(r.RemoteAddr)
	if errSplit != nil {
		log.Printf("SplitHostPort error: %v", errSplit)
		WriteJSON(w, http.StatusBadRequest, "Bad Request")
		return
	}

	agent := r.Header.Get("User-Agent")

	tokens, errExecute := h.UseCase.Execute(r.Context(), userInfo.Email, userInfo.Password, host, agent)

	if errExecute != nil {
		log.Printf("Execute error: %v", errExecute)
		WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	WriteLoginJSON(w, http.StatusOK, tokens.AccessToken, tokens.RefreshToken)
}
