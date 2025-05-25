package handler

import (
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/FooxyS/auth-service/internal/usecase"
	"github.com/FooxyS/auth-service/pkg/apperrors"
	"github.com/FooxyS/auth-service/pkg/consts"
)

type RefreshHandler struct {
	UseCase usecase.RefreshUseCase
}

func (h *RefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	access := r.Header.Get("Authorization")

	cookie, errGetCookie := r.Cookie(consts.RefreshCookieName)
	if errGetCookie != nil {
		log.Printf("GetCookie error: %v", errGetCookie)
		WriteJSON(w, http.StatusBadRequest, "Bad Request")
		return
	}
	refresh := cookie.Value

	host, _, errSplit := net.SplitHostPort(r.RemoteAddr)
	if errSplit != nil {
		log.Printf("SplitHostPort error: %v", errSplit)
		WriteJSON(w, http.StatusBadRequest, "Bad Request")
		return
	}

	agent := r.Header.Get("User-Agent")

	tokens, errExecute := h.UseCase.Execute(r.Context(), access, refresh, host, agent)

	if errors.Is(errExecute, apperrors.ErrIPMismatch) {
		WriteJSON(w, http.StatusConflict, "IP address mismatch")
		return
	}
	if errors.Is(errExecute, apperrors.ErrAgentMismatch) {
		WriteJSON(w, http.StatusConflict, "User-Agent mismatch")
		return
	}

	if errExecute != nil {
		log.Printf("Execute error: %v", errExecute)
		WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	cookieRefresh := &http.Cookie{
		Name:     consts.RefreshCookieName,
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookieRefresh)

	WriteLoginJSON(w, http.StatusOK, tokens.AccessToken)
}
