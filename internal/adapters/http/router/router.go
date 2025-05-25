package router

import (
	"github.com/FooxyS/auth-service/internal/adapters/http/handler"
	"github.com/FooxyS/auth-service/internal/domain"
	"net/http"
)

func SetupRouter(repository domain.SessionRepository, sessionRepository domain.SessionRepository, service domain.TokenService, hasher domain.PasswordHasher) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/register", &handler.RegisterHandler{})
	mux.Handle("/login", &handler.LoginHandler{})
	mux.Handle("/logout", &handler.LogoutHandler{})
	mux.Handle("/refresh", &handler.RefreshHandler{})
	mux.Handle("/me", &handler.MeHandler{})

	return mux
}
