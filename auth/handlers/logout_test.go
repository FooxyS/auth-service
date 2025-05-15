package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	auth "github.com/FooxyS/auth-service/auth/handlers"
	"github.com/FooxyS/auth-service/auth/models"
	"github.com/FooxyS/auth-service/pkg/consts"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func TestLogoutHandler(t *testing.T) {
	errGetEnv := godotenv.Load()
	if errGetEnv != nil {
		t.Errorf("error with getting env variable: %v\n", errGetEnv)
		return
	}

	//добавление в БД информации о сессии
	dburl := os.Getenv("DATABASE_URL_TEST")

	pgpool, errConnDB := pgxpool.New(context.Background(), dburl)
	if errConnDB != nil {
		t.Errorf("error with connecting to postgres: %v\n", errConnDB)
		return
	}
	defer pgpool.Close()

	testClaims := models.MyCustomClaims{
		UserID: "72257344-9a79-4f52-9108-527fbaa73bb6",
		PairID: "562093e3-cfcb-4f7e-931e-a13816ceead9",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	secretString := os.Getenv("JWT_KEY")

	accesstoken, errGenAccess := jwt.NewWithClaims(jwt.SigningMethodHS512, testClaims).SignedString([]byte(secretString))
	if errGenAccess != nil {
		t.Errorf("error with generating new access token: %v\n", errGenAccess)
		return
	}

	newsession := models.Session{
		ID:           "72257344-9a79-4f52-9108-527fbaa73bb6",
		IP:           "192.103.32.12",
		RefreshToken: "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		PairID:       "562093e3-cfcb-4f7e-931e-a13816ceead9",
		UserAgent:    "google chrome",
	}

	_, errWithExec := pgpool.Exec(context.Background(),
		"INSERT INTO session_table (user_id, ip_address, refresh_token, pair_id, useragent) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING",
		newsession.ID, newsession.IP, newsession.RefreshToken, newsession.PairID, newsession.UserAgent)
	if errWithExec != nil {
		t.Errorf("error with Exec(): %v\n", errWithExec)
		return
	}

	//создание тел запроса и ответа
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)

	req.Header.Set("Authorization", "Bearer "+accesstoken)

	ctx := context.WithValue(req.Context(), consts.CTX_KEY_DB, pgpool)
	reqwithctx := req.WithContext(ctx)

	auth.LogoutHandler(resp, reqwithctx)

	session := new(models.Session)
	ErrNoRows := pgpool.QueryRow(context.Background(), "select * from session_table where user_id=$1", newsession.ID).Scan(&session.ID, &session.IP, &session.PairID, &session.RefreshToken, &session.UserAgent)
	if ErrNoRows == nil {
		t.Error("there is row in db: want nil, got row")
		return
	}

	expectedMessage := "Сессия пользователя успешно удалена"
	mess := resp.Body.Bytes()
	if string(mess) != expectedMessage {
		t.Errorf("no match with returned message: want %s, got %s\n", expectedMessage, mess)
	}

	_, errTrunTable := pgpool.Exec(context.Background(), "TRUNCATE TABLE session_table RESTART IDENTITY")
	if errTrunTable != nil {
		t.Errorf("error with truncating the db: %v\n", errTrunTable)
		return
	}
}
