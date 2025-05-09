package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/FooxyS/auth-service/router"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	//подключение к БД
	errLoadEnv := godotenv.Load()
	if errLoadEnv != nil {
		log.Fatalf("error with loading env file: %v\n", errLoadEnv)
	}

	dburl := os.Getenv("DATABASE_URL")

	pgpool, errConnPostgres := pgxpool.New(context.Background(), dburl)
	if errConnPostgres != nil {
		log.Fatalf("error with connection to postgres: %v\n", errConnPostgres)
	}

	//запуск сервера
	router := router.SetupRouter(pgpool)
	errStart := http.ListenAndServe("localhost:8080", router)
	if errStart != nil {
		log.Fatalf("error with starting the server: %v\n", errStart)
	}
}
