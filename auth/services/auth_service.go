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
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/FooxyS/auth-service/auth/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// используется для получения переменной из env
func GetFromEnv(key string) (string, error) {
	errGotEnv := godotenv.Load()
	if errGotEnv != nil {
		return "", errGotEnv
	}
	val := os.Getenv(key)
	return val, nil
}

// используется для генерации refresh токена
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, errGenRandStr := rand.Read(b)
	if errGenRandStr != nil {
		return "", errGenRandStr
	}
	result := base64.URLEncoding.EncodeToString(b)
	return result, nil
}

// используется для отделения токена от bearer
func ParseBearerToToken(s string) (string, error) {
	substr := strings.Split(s, " ")
	if len(substr) < 2 {
		return "", errors.New("massive is too short. Out of range")
	}

	return substr[1], nil
}

// используется для отправки вэбхука
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

// get cookie with refresh and check it with the hashed refresh in db (also check other fields of cookie structure)
func GetRefreshFromCookie(resp *httptest.ResponseRecorder, t *testing.T) string {
	t.Helper()

	refcookie := new(http.Cookie)
	cookies := resp.Result().Cookies()

	found := false
	for _, cookie := range cookies {
		if cookie.Name == "refreshtoken" {
			refcookie = cookie
			found = true
			break
		}
	}
	if !found {
		t.Errorf("there is no expected cookie")
		return ""
	}

	if !refcookie.HttpOnly {
		t.Error("wrong httponly field: want true, got false")
		return ""
	}

	return refcookie.Value
}

func DeleteUserByID(pgpool *pgxpool.Pool, r *http.Request, id string) error {
	_, errExec := pgpool.Exec(r.Context(), "DELETE FROM session_table WHERE user_id=$1", id)
	if errExec != nil {
		log.Printf("error with Exec(): %v\n", errExec)
		return errExec
	}
	return nil
}
