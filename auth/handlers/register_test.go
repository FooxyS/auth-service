package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	auth "github.com/FooxyS/auth-service/auth/handlers"
	"github.com/FooxyS/auth-service/auth/models"
	"github.com/FooxyS/auth-service/auth/services"
	"github.com/FooxyS/auth-service/pkg/consts"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterHandler(t *testing.T) {
	dburl, errGetDBKey := services.GetFromEnv(consts.DATABASE_URL_TEST)
	if errGetDBKey != nil {
		t.Errorf("error with getting db url: %v\n", errGetDBKey)
		return
	}
	pgpool, errConnTestDB := pgxpool.New(context.Background(), dburl)
	if errConnTestDB != nil {
		t.Errorf("error with connection to db: %v\n", errConnTestDB)
		return
	}
	defer pgpool.Close()
	defer pgpool.Exec(context.Background(), "TRUNCATE TABLE users RESTART IDENTITY")
	defer pgpool.Exec(context.Background(), "TRUNCATE TABLE session_table RESTART IDENTITY")

	testdata := models.UserData{
		Email:    "example@gmail.com",
		Password: "example123",
	}

	body, errMarshalJson := json.Marshal(testdata)
	if errMarshalJson != nil {
		t.Errorf("error with marshalling test data: %v\n", errMarshalJson)
		return
	}
	buf := bytes.NewBuffer(body)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", buf)

	ctx := context.WithValue(context.Background(), consts.CTX_KEY_DB, pgpool)
	reqCtx := req.WithContext(ctx)

	auth.RegisterHandler(resp, reqCtx)

	//проверить код, заголовки, тело запроса
	expectedType := "application/json"
	if contType := resp.Header().Get("Content-Type"); contType != expectedType {
		t.Errorf("wrong content type: got %s, want %s\n", contType, expectedType)
		return
	}

	expectedCode := http.StatusCreated
	if code := resp.Code; code != expectedCode {
		t.Errorf("wrong status code: got %d, want %d\n", code, expectedCode)
		return
	}

	//получение access, парсинг. Достаём user_id, чтобы делать запросы в БД
	givenAccess := new(models.AccessTokenJson)
	errDecodeData := json.NewDecoder(resp.Body).Decode(givenAccess)
	if errDecodeData != nil {
		t.Errorf("error with decoding json: %v\n", errDecodeData)
		return
	}

	jwtkey, errGetEnv := services.GetFromEnv(consts.JWT_KEY)
	if errGetEnv != nil {
		t.Errorf("error with getting key from env: %v\n", errGetEnv)
		return
	}

	accessClaims := new(models.MyCustomClaims)
	jwt.ParseWithClaims(givenAccess.Access, accessClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtkey), nil
	})

	//запрос к БД с информацией пользователя
	userInfoFromDB := new(models.User)
	errExec := pgpool.QueryRow(context.Background(), "SELECT * FROM users WHERE user_id=$1", accessClaims.UserID).Scan(&userInfoFromDB.ID, &userInfoFromDB.Email, &userInfoFromDB.Password, &userInfoFromDB.CreationDate)
	if errExec != nil {
		t.Errorf("error wtih Exec(): %v\n", errExec)
		return
	}

	//проверить email, сверить пароль
	if userInfoFromDB.Email != testdata.Email {
		t.Errorf("wrong email in db: got %s, want %s\n", userInfoFromDB.Email, testdata.Email)
		return
	}

	errCompPass := bcrypt.CompareHashAndPassword([]byte(userInfoFromDB.Password), []byte(testdata.Password))
	if errCompPass != nil {
		t.Errorf("wrong pass in db: %v\n", errCompPass)
		return
	}

	//запрос к БД с информацией о сессии
	sessionFromDB := new(models.Session)
	errExecSession := pgpool.QueryRow(context.Background(), "SELECT * FROM session_table WHERE user_id=$1", accessClaims.UserID).Scan(&sessionFromDB.ID, &sessionFromDB.IP, &sessionFromDB.RefreshToken, &sessionFromDB.PairID, &sessionFromDB.UserAgent)
	if errExecSession != nil {
		t.Errorf("error wtih Exec(): %v\n", errExecSession)
		return
	}

	//проверить refresh_token, pair_id
	refresh := services.GetRefreshFromCookie(resp, t)

	errCompRefresh := bcrypt.CompareHashAndPassword([]byte(sessionFromDB.RefreshToken), []byte(refresh))
	if errCompRefresh != nil {
		t.Errorf("wrong refresh in db: got %s, want %s\n", sessionFromDB.RefreshToken, refresh)
		return
	}

	if sessionFromDB.PairID != accessClaims.PairID {
		t.Errorf("wrong pairid: got %s, want %s\n", sessionFromDB.PairID, accessClaims.PairID)
		return
	}
}
