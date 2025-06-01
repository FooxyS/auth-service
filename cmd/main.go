package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/FooxyS/auth-service/internal/adapters/http/router"
	"github.com/FooxyS/auth-service/internal/infrastructure/hasher"
	"github.com/FooxyS/auth-service/internal/infrastructure/postgres"
	"github.com/FooxyS/auth-service/internal/infrastructure/tokens"
	"github.com/FooxyS/auth-service/pkg/consts"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDBURL() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	dburl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, name)

	os.Setenv("DATABASE_URL", dburl)
}

func main() {

	GetDBURL()

	pgpool, errDB := pgxpool.New(context.Background(), os.Getenv(consts.DATABASE_URL))
	if errDB != nil {
		log.Fatalf("error with connecting to db: %v", errDB)
	}

	router := router.SetupRouter(postgres.NewUserRepo(pgpool), postgres.NewSessionRepo(pgpool), tokens.New(), hasher.New())

	if errServer := http.ListenAndServe("0.0.0.0:8080", router); errServer != nil {
		log.Fatalf("errors with launching the server: %v", errServer)
	}
}
