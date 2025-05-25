package router

import (
	"net/http"

	"github.com/FooxyS/auth-service/internal/adapters/http/handler"
	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/FooxyS/auth-service/internal/usecase"
)

func SetupRouter(userRepo domain.UserRepository, sessionRepo domain.SessionRepository, tokens domain.TokenService, hasher domain.PasswordHasher) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/register", &handler.RegisterHandler{UseCase: usecase.RegisterUseCase{UserRepo: userRepo, Hasher: hasher}})
	mux.Handle("/login", &handler.LoginHandler{UseCase: usecase.LoginUseCase{UserRepo: userRepo, SessionRepo: sessionRepo, Tokens: tokens, Hasher: hasher}})
	mux.Handle("/logout", &handler.LogoutHandler{UseCase: usecase.LogoutUseCase{SessionRepo: sessionRepo, Tokens: tokens}})
	mux.Handle("/refresh", &handler.RefreshHandler{UseCase: usecase.RefreshUseCase{Tokens: tokens, SessionRepo: sessionRepo, Hasher: hasher}})
	mux.Handle("/me", &handler.MeHandler{UseCase: usecase.MeUseCase{Tokens: tokens, UserRepo: userRepo}})

	return mux
}
