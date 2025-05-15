package middleware

import (
	"context"
	"net/http"

	"github.com/FooxyS/auth-service/pkg/consts"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PostgresWithContext(pgpool *pgxpool.Pool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), consts.CTX_KEY_DB, pgpool)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
