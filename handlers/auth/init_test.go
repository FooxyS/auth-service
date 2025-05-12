package auth_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/FooxyS/auth-service/handlers/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func TestInitHandler(t *testing.T) {
	errGetEnv := godotenv.Load()
	if errGetEnv != nil {
		t.Errorf("error with getting env variable: %v\n", errGetEnv)
		return
	}

	//соединене с базой данных
	dburl := os.Getenv("DATABASE_URL_TEST")

	pgpool, errConnDB := pgxpool.New(context.Background(), dburl)
	if errConnDB != nil {
		t.Errorf("error with connecting to postgres: %v\n", errConnDB)
		return
	}
	defer pgpool.Close()

	newUserID := uuid.New()

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/auth/init?userid=%s", newUserID.String()), nil)

	expectedAgent := "Google Chrome"
	req.Header.Set("User-Agent", expectedAgent)

	expectedIP := "123.123.123.123"
	req.RemoteAddr = expectedIP + ":8080"

	ctx := context.WithValue(context.Background(), "postgres", pgpool)
	reqWithCtx := req.WithContext(ctx)

	auth.InitHandler(resp, reqWithCtx)

	//check db
	checksession := new(auth.Session)
	errQueryRow := pgpool.QueryRow(context.Background(),
		"select * from session_table where user_id=$1",
		newUserID.String()).Scan(&checksession.ID, &checksession.IP, &checksession.PairID, &checksession.RefreshToken, &checksession.UserAgent)
	if errQueryRow != nil {
		t.Errorf("there is no row with expected data in db: %v\n", errQueryRow)
		return
	}
	if checksession.IP != expectedIP {
		t.Errorf("wrong ip address: want %s, got %s\n", expectedIP, checksession.IP)
		return
	}
	if checksession.UserAgent != expectedAgent {
		t.Errorf("wrong UserAgent: want %s, got %s\n", expectedAgent, checksession.UserAgent)
		return
	}

	//check json, headers, statuscode
	expectedStatusCode := http.StatusOK
	if resp.Code != expectedStatusCode {
		t.Errorf("wrong status code recieved: want %v, got %v\n", expectedStatusCode, resp.Code)
		return
	}

	expectedHeaderType := "application/json"
	if header := resp.Header().Get("Content-Type"); header != expectedHeaderType {
		t.Errorf("wrong header type: want %s, got %s\n", expectedHeaderType, header)
		return
	}

	//parse access and check pair id, user id
	accesstoken := new(auth.AccessTokenJson)
	errDecodeJson := json.NewDecoder(resp.Body).Decode(accesstoken)
	if errDecodeJson != nil {
		t.Error("error with decoding the json")
		return
	}

	checkclaims := new(auth.MyCustomClaims)
	jwt.ParseWithClaims(accesstoken.Access, checkclaims, func(t *jwt.Token) (interface{}, error) {
		return os.Getenv("JWT_KEY"), nil
	})
	if checkclaims.UserID != newUserID.String() {
		t.Errorf("wrong userid in access: want %s, got %s\n", newUserID.String(), checkclaims.UserID)
		return
	}
	if checkclaims.PairID != checksession.PairID {
		t.Errorf("wrong pairid: got %s, want %s\n", checkclaims.PairID, checksession.PairID)
		return
	}

	//get cookie with refresh and check it with the hashed refresh in db (also check other fields of cookie structure)
	refcookie := new(http.Cookie)
	cookies := resp.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "refreshtoken" {
			refcookie = cookie
		} else {
			t.Error("there is no expected cookie")
			return
		}
	}

	if refcookie.HttpOnly == false {
		t.Error("wrong httponly field: want true, got false")
		return
	}

	refreshtoken := refcookie.Value
	errCompare := bcrypt.CompareHashAndPassword([]byte(checksession.RefreshToken), []byte(refreshtoken))
	if errCompare != nil {
		t.Errorf("no match with hash and token: db %v, got%v\n", checksession.RefreshToken, refreshtoken)
		return
	}
}
