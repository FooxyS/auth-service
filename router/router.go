package router

import (
	"net/http"

	auth "github.com/FooxyS/auth-service/handlers/auth"
)

func SetupRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/init", auth.InitHandler)
	mux.HandleFunc("/auth/logout", auth.LogoutHandler)
	mux.HandleFunc("/auth/me", auth.MeHandler)
	mux.HandleFunc("/auth/refresh", auth.RefreshHandler)
	return mux
}
