package router

import (
	"net/http"

	"github.com/FooxyS/auth-service/internal/adapters/http/handler"
	"github.com/FooxyS/auth-service/internal/domain"
)

func SetupRouter(UserRepo domain.UserRepository, SessionRepo domain.SessionRepository, Tokens domain.TokenService, Hasher domain.PasswordHasher) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/register", &handler.RegisterHandler{})
	mux.Handle("/login", &handler.LoginHandler{})
	mux.Handle("/logout", &handler.LogoutHandler{})
	mux.Handle("/refresh", &handler.RefreshHandler{})
	mux.Handle("/me", &handler.MeHandler{})

	return mux
}
