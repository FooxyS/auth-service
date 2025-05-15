package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"strings"

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

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, errGenRandStr := rand.Read(b)
	if errGenRandStr != nil {
		return "", errGenRandStr
	}
	result := base64.URLEncoding.EncodeToString(b)
	return result, nil
}

func ParseTokenFromHeader(s string) (string, error) {
	substr := strings.Split(s, " ")
	if len(substr) < 2 {
		return "", errors.New("massive is too short. Out of range")
	}

	return substr[1], nil
}
