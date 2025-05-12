package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	ID           string `json:"id"`
	IP           string `json:"ip"`
	RefreshToken string `json:"refreshtoken"`
	PairID       string `json:"pairid"`
	UserAgent    string `json:"useragent"`
}

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
	//достаём пул подключений из контекста
	pgpool, ok := r.Context().Value("postgres").(*pgxpool.Pool)
	if !ok {
		log.Println("value not found in context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	idFromURL := r.URL.Query().Get("userid")
	if idFromURL == "" {
		http.Error(w, "query param is empty", http.StatusBadRequest)
		return
	}

	var ExistedUserID string

	//проверка, есть ли пользователь в БД
	errQueryRow := pgpool.QueryRow(r.Context(), "select user_id from session_table where user_id=$1", idFromURL).Scan(&ExistedUserID)
	if errQueryRow == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("пользователь уже зарегистрирован"))
		return
	}

	//структура с данными сессии пользователя
	session := new(Session)

	//добавление GUID в информацию о сессии
	session.ID = idFromURL

	//создание id пары токенов, по которому мы сможем определить были ли они выданы вместе
	pairid := uuid.New().String()

	jwtkey, errGotEnv := GetFromEnv("JWT_KEY")
	if errGotEnv != nil {
		log.Printf("error with env: %v", errGotEnv)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//добавление pairid в информацию о сессии
	session.PairID = pairid

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
	accessToken, errGenAccess := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims).SignedString([]byte(jwtkey))
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

	cookie := http.Cookie{
		Name:     "refreshtoken",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

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

	//добавление refresh токена в БД
	hashedRefresh, errHashRefresh := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if errHashRefresh != nil {
		log.Printf("error with hashing refresh token: %v", errHashRefresh)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//добавление hashedrefresh в информацию о сессии
	session.RefreshToken = string(hashedRefresh)

	//добавление useragent в информацию о сессии
	agent := r.Header.Get("User-Agent")
	session.UserAgent = agent

	//добавление IP в информацию о сессии
	ip, _, errSplitHostPost := net.SplitHostPort(r.RemoteAddr)
	if errSplitHostPost != nil {
		log.Printf("error with spliting remoteaddr: %v\n", errSplitHostPost)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	session.IP = ip

	//добавляем информацию о сессии в БД
	_, errWithExec := pgpool.Exec(r.Context(), "INSERT INTO session_table (user_id, ip_address, refresh_token, pair_id, useragent) VALUES ($1, $2, $3, $4, $5)", session.ID, session.IP, session.RefreshToken, session.PairID, session.UserAgent)
	if errWithExec != nil {
		log.Printf("error with Exec: %v\n", errWithExec)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
