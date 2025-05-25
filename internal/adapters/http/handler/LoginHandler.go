package handler

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/FooxyS/auth-service/internal/usecase"
	"github.com/FooxyS/auth-service/pkg/consts"
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

	if userInfo.Email == "" || userInfo.Password == "" || host == "" || agent == "" {
		WriteJSON(w, http.StatusBadRequest, "Bad Request")
		return
	}

	tokens, errExecute := h.UseCase.Execute(r.Context(), userInfo.Email, userInfo.Password, host, agent)

	if errExecute != nil {
		log.Printf("Execute error: %v", errExecute)
		WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	cookie := &http.Cookie{
		Name:     consts.RefreshCookieName,
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)

	WriteLoginJSON(w, http.StatusOK, tokens.AccessToken)
}
