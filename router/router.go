package router

import (
	"net/http"

	auth "github.com/FooxyS/auth-service/auth/handlers"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/FooxyS/auth-service/middleware"
)

func SetupRouter(db *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/init", auth.InitHandler)
	mux.HandleFunc("/auth/logout", auth.LogoutHandler)
	mux.Handle("/auth/me", middleware.AuthMiddleware(http.HandlerFunc(auth.MeHandler)))
	mux.HandleFunc("/auth/refresh", auth.RefreshHandler)
	wrappedmux := middleware.PostgresWithContext(db, mux)

	return wrappedmux
}
