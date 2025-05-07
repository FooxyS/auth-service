package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type MyCustomClaims struct {
	UserID string
	PairID string
	jwt.RegisteredClaims
}

type AccessTokenJson struct {
	Access string `json:"access"`
}

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

func InitHandler(w http.ResponseWriter, r *http.Request) {
	idFromURL := r.URL.Query().Get("userid")
	if idFromURL == "" {
		http.Error(w, "query param is empty", http.StatusBadRequest)
		return
	}

	//создание id пары токенов, по которому мы сможем определить были ли они выданы вместе
	pairid := uuid.New().String()

	jwtkey, errGotEnv := GetFromEnv("JWT_KEY")
	if errGotEnv != nil {
		log.Printf("error with env: %v", errGotEnv)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//определение параметров для создания JWT токена
	accessClaims := MyCustomClaims{
		UserID: idFromURL,
		PairID: pairid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	//создание access и refresh токенов
	accessToken, errGenAccess := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims).SignedString(jwtkey)
	if errGenAccess != nil {
		log.Printf("error with creating jwt: %v", errGenAccess)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	refreshToken, errGenRandStr := GenerateRefreshToken()
	if errGenRandStr != nil {
		log.Printf("error with creating refresh string: %v", errGenRandStr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//отправка токенов
	accessResp := AccessTokenJson{
		Access: accessToken,
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	errJsonResp := json.NewEncoder(w).Encode(accessResp)
	if errJsonResp != nil {
		log.Printf("error with encoding json resp: %v", errJsonResp)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "refresh-token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)

	//добавление refresh токена в БД
	hashedRefresh, errHashRefresh := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if errHashRefresh != nil {
		log.Printf("error with hashing refresh token: %v", errHashRefresh)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("логика добавления хэшированного токена в БД: %v", hashedRefresh)
	//нужно добавить в БД userID, hashedrefreshtoken, User-Agent, IP, pairID
}
