package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/FooxyS/auth-service/handlers/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func TestRefreshHandler(t *testing.T) {
	//подключаемся к тестовой БД
	godotenv.Load()
	dburl := os.Getenv("DATABASE_URL_TEST")
	pgpool, errWithConn := pgxpool.New(context.Background(), dburl)
	if errWithConn != nil {
		t.Errorf("error with connecting to db: %v\n", errWithConn)
		return
	}
	defer pgpool.Close()

	//информация для записи в БД
	newsession := auth.Session{
		ID:           "72257344-9a79-4f52-9108-527fbaa73bb6",
		IP:           "192.103.32.12",
		RefreshToken: "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		PairID:       "562093e3-cfcb-4f7e-931e-a13816ceead9",
		UserAgent:    "google chrome",
	}

	_, errWithExec := pgpool.Exec(context.Background(), "INSERT INTO session_table (user_id, ip_address, refresh_token, pair_id, useragent) VALUES ($1, $2, $3, $4, $5)", newsession.ID, newsession.IP, newsession.RefreshToken, newsession.PairID, newsession.UserAgent)
	if errWithExec != nil {
		t.Errorf("error with pushing the test data to db: %v\n", errWithExec)
		return
	}

	secret := os.Getenv("JWT_KEY")

	newClaims := auth.MyCustomClaims{
		UserID: newsession.ID,
		PairID: newsession.PairID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	access, errWithNewAccess := jwt.NewWithClaims(jwt.SigningMethodHS512, newClaims).SignedString([]byte(secret))
	if errWithNewAccess != nil {
		t.Errorf("error with creating new access token: %v\n", errWithNewAccess)
		return
	}

	//формирование тел запроса и ответа
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)

	req.Header.Set("Authorization", access)

	ctx := context.WithValue(context.Background(), "postgres", pgpool)
	reqWithCtx := req.WithContext(ctx)

	auth.RefreshHandler(resp, reqWithCtx)

	sessionFromDB := new(auth.Session)

	errWithQueryRow := pgpool.QueryRow(context.Background(), "SELECT * FROM session_table WHERE user_id=$1", newsession.ID).Scan(&sessionFromDB.ID, &sessionFromDB.IP, &sessionFromDB.PairID, &sessionFromDB.RefreshToken, &sessionFromDB.UserAgent)
	if errWithQueryRow != nil {
		t.Errorf("error with QueryRow: %v\n", errWithQueryRow)
		return
	}

	if sessionFromDB.PairID == newsession.PairID {
		t.Error("pairid wasn't changed after RefreshHandler")
		return
	}

	errCompHash := bcrypt.CompareHashAndPassword([]byte(sessionFromDB.RefreshToken), []byte(newsession.RefreshToken))
	if errCompHash == nil {
		t.Error("hashed refresh token wasn't changed after RefreshHandler")
		return
	}

	_, errTrunTable := pgpool.Exec(context.Background(), "TRUNCATE TABLE session_table RESTART IDENTITY")
	if errTrunTable != nil {
		t.Errorf("error with clearing the data from db: %v\n", errTrunTable)
		return
	}
}
