package router

import (
	"net/http"

	"github.com/FooxyS/auth-service/internal/adapters/http/handler"
	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/FooxyS/auth-service/internal/usecase"
)

func SetupRouter(UserRepo domain.UserRepository, SessionRepo domain.SessionRepository, Tokens domain.TokenService, Hasher domain.PasswordHasher) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/register", &handler.RegisterHandler{UseCase: usecase.RegisterUseCase{UserRepo: UserRepo, Hasher: Hasher}})
	mux.Handle("/login", &handler.LoginHandler{UseCase: usecase.LoginUseCase{UserRepo: UserRepo, SessionRepo: SessionRepo, Tokens: Tokens, Hasher: Hasher}})
	mux.Handle("/logout", &handler.LogoutHandler{UseCase: usecase.LogoutUseCase{SessionRepo: SessionRepo, Tokens: Tokens}})
	mux.Handle("/refresh", &handler.RefreshHandler{UseCase: usecase.RefreshUseCase{Tokens: Tokens, SessionRepo: SessionRepo, Hasher: Hasher}})
	mux.Handle("/me", &handler.MeHandler{UseCase: usecase.MeUseCase{Tokens: Tokens, UserRepo: UserRepo}})

	return mux
}
