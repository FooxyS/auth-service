package services

import (
	"os"

	"github.com/joho/godotenv"
)

func GetFromEnv(key string) (string, error) {
	errGotEnv := godotenv.Load()
	if errGotEnv != nil {
		return "", errGotEnv
	}
	val := os.Getenv(key)
	return val, nil
}
