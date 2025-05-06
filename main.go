package main

import (
	"log"
	"net/http"

	"github.com/FooxyS/auth-service/router"
)

func main() {
	router := router.SetupRouter()
	errStart := http.ListenAndServe("localhost:8080", router)
	if errStart != nil {
		log.Fatalf("error with starting the server: %v\n", errStart)
	}
}
