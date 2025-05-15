package services

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/FooxyS/auth-service/auth/models"
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

func SendWebhook(ip string) error {
	//создание ответа
	text := fmt.Sprintf("Попытка входа с нового IP: %s\n", ip)

	respBody := models.WebhookJson{
		Message: text,
	}

	jsonRespBody, errParseJson := json.Marshal(respBody)
	if errParseJson != nil {
		return errParseJson
	}

	//достаём URL из env
	webhookurl, errGetEnv := GetFromEnv("WEBHOOK_URL")
	if errGetEnv != nil || webhookurl == "" {
		log.Printf("error with get from env: %v\n", errGetEnv)
		return errGetEnv
	}

	//формируем новый запрос
	req, errResp := http.NewRequest(http.MethodPost, webhookurl, bytes.NewBuffer([]byte(jsonRespBody)))
	if errResp != nil {
		log.Printf("error with creating new request: %v\n", errResp)
		return errResp
	}

	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)

	resp, errDoReq := client.Do(req)
	if errDoReq != nil {
		log.Printf("error with sending request: %v\n", errDoReq)
		return errDoReq
	}

	defer resp.Body.Close()

	log.Printf("Webhook отправлен. Статус ответа: %v\n", resp.Status)
	return nil
}
