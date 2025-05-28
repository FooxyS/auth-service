package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/FooxyS/auth-service/internal/adapters/http/router"
	"github.com/FooxyS/auth-service/internal/infrastructure/hasher"
	"github.com/FooxyS/auth-service/internal/infrastructure/postgres"
	"github.com/FooxyS/auth-service/internal/infrastructure/tokens"
	"github.com/FooxyS/auth-service/pkg/consts"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if errEnv := godotenv.Load(); errEnv != nil {
		log.Fatalf("error with loading .env: %v", errEnv)
	}

	pgpool, errDB := pgxpool.New(context.Background(), os.Getenv(consts.DATABASE_URL))
	if errDB != nil {
		log.Fatalf("error with connecting to db: %v", errDB)
	}

	router := router.SetupRouter(postgres.NewUserRepo(pgpool), postgres.NewSessionRepo(pgpool), tokens.New(), hasher.New())

	if errServer := http.ListenAndServe("localhost:8080", router); errServer != nil {
		log.Fatalf("errors with launching the server: %v", errServer)
	}
}
