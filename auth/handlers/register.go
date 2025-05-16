package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/FooxyS/auth-service/auth/models"
	"github.com/FooxyS/auth-service/auth/services"
	"github.com/FooxyS/auth-service/pkg/consts"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// эндпоинт, который должен регистрировать пользователя: принимать POST-запрос с данными пользователя, регистрировать, если ещё нет в базе
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	pgpool := r.Context().Value(consts.CTX_KEY_DB).(*pgxpool.Pool)

	registrationData := new(models.UserData)

	errWithDecode := json.NewDecoder(r.Body).Decode(registrationData)
	if errWithDecode != nil {
		log.Printf("error with decoding json data: %v\n", errWithDecode)
		http.Error(w, "error with json", http.StatusBadRequest)
		return
	}

	if registrationData.Email == "" || registrationData.Password == "" {
		http.Error(w, "json fields are empty", http.StatusBadRequest)
		return
	}

	UserForCheck := new(models.UserData)
	errQueryRow := pgpool.QueryRow(r.Context(), "select user_id from users where email=$1", registrationData.Email).Scan(&UserForCheck.UserID)
	if errQueryRow == nil {
		w.WriteHeader(http.StatusOK)
		message := fmt.Sprintf("Пользователь %v уже зарегистрирован!\n", UserForCheck.UserID)
		w.Write([]byte(message))
		return
	}

	newIDForUser := uuid.New()

	hashedPassword, errWithHashing := bcrypt.GenerateFromPassword([]byte(registrationData.Password), bcrypt.DefaultCost)
	if errWithHashing != nil {
		log.Printf("error with hashing: %v\n", errWithHashing)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	creationDate := time.Now()

	_, errExec := pgpool.Exec(r.Context(), "insert into users (user_id, email, password, creation_date) VALUES ($1, $2, $3, $4)", newIDForUser.String(), registrationData.Email, hashedPassword, creationDate)
	if errExec != nil {
		log.Printf("error with Exec: %v\n", errExec)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	newPairID := uuid.New()

	accessdata := models.MyCustomClaims{
		UserID: newIDForUser.String(),
		PairID: newPairID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	secStr, errWithEnv := services.GetFromEnv(consts.JWT_KEY)
	if errWithEnv != nil {
		log.Printf("error with GetFromEnv(): %v\n", errWithEnv)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	access, errWithGenAccess := jwt.NewWithClaims(jwt.SigningMethodHS512, accessdata).SignedString([]byte(secStr))
	if errWithGenAccess != nil {
		log.Printf("error with NewWithClaims(): %v\n", errWithGenAccess)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//дописать логику отправки токенов + собрать нужную инфу и загрузить в БД.
	//также написать функцию, которая будет возвращать структуру с нужными данными
	//создать в consts для Env ключей

	refresh, errWithGenRefresh := services.GenerateRefreshToken()
	if errWithGenRefresh != nil {
		log.Printf("error with GenerateRefreshToken(): %v\n", errWithGenRefresh)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	hashedRefresh, errWithHashing := bcrypt.GenerateFromPassword([]byte(refresh), bcrypt.DefaultCost)
	if errWithHashing != nil {
		log.Printf("error with hashing: %v\n", errWithHashing)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "refreshtoken",
		Value:    refresh,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	accessToken := models.AccessTokenJson{
		Access: access,
	}

	host, _, errSplitHostPost := net.SplitHostPort(r.RemoteAddr)
	if errSplitHostPost != nil {
		log.Printf("error with SplitHostPort(): %v\n", errSplitHostPost)
	}

	agent := r.Header.Get("User-Agent")

	_, errInsertIntoSession := pgpool.Exec(r.Context(), "INSERT INTO session_table (user_id, ip_address, refresh_token, pair_id, useragent) VALUES ($1, $2, $3, $4, $5)", newIDForUser, host, hashedRefresh, newPairID, agent)
	if errInsertIntoSession != nil {
		log.Printf("error with Exec(): %v\n", errInsertIntoSession)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	errEncodeResp := json.NewEncoder(w).Encode(accessToken)
	if errEncodeResp != nil {
		log.Printf("error with encoding json resonse: %v\n", errEncodeResp)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
