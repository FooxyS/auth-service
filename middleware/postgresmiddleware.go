package middleware

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func PostgresWithContext(pgpool *pgxpool.Pool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "postgres", pgpool)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
