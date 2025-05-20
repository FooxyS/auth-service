package models

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/FooxyS/auth-service/auth/apperrors"
	"github.com/FooxyS/auth-service/pkg/consts"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// используется для работы с таблицой сессий пользователей в postgres
type Session struct {
	ID           string `json:"id"`
	IP           string `json:"ip"`
	RefreshToken string `json:"refreshtoken"`
	PairID       string `json:"pairid"`
	UserAgent    string `json:"useragent"`
}

// используется для работы с данными пользователей в postgres
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	CreationDate time.Time `json:"creation_date"`
}

// используется для генерации access токена
type MyCustomClaims struct {
	UserID string
	PairID string
	jwt.RegisteredClaims
}

// используется для отправки access токена в формате json пользователю
type AccessTokenJson struct {
	Access string `json:"access"`
}

// используется для отправки пользоватею его GUID из обработчика /auth/me и /auth/register
type UserJsonID struct {
	UserID string `json:"userid"`
}

// используется для отправки сообщения на webhook
type WebhookJson struct {
	Message string `json:"message"`
}

// клиент Postgres
type Postgres struct {
	pgpool *pgxpool.Pool
}

// сохраняет пользователя в БД
func (pg Postgres) SaveUser(ctx context.Context, user UserData) error {
	creation := time.Now()

	hashedPass, errGenHash := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if errGenHash != nil {
		return errGenHash
	}

	_, errWithExec := pg.pgpool.Exec(ctx, "INSERT INTO users (user_id, email, password, creation_date) VALUES ($1, $2, $3, $4)", user.UserID, user.Email, hashedPass, creation)
	if errWithExec != nil {
		return errWithExec
	}

	return nil
}

// удаляет пользователя по user_id
func (pg Postgres) DeleteByID(ctx context.Context, id string, tablename string) error {
	_, errDelete := pg.pgpool.Exec(ctx, "DELETE FROM $1 WHERE user_id=$2", tablename, id)
	if errDelete != nil {
		return errDelete
	}
	return nil
}

// ищет пользователя в базе данных по заданному email. Если такой есть, возвращает user_id
func (pg Postgres) FindUser(ctx context.Context, email string) (string, error) {
	var id string

	err := pg.pgpool.QueryRow(ctx, "SELECT user_id FROM users WHERE email=$1", email).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// достаёт клиент БД из контекста
func (pg *Postgres) FromContext(ctx context.Context) *Postgres {
	pg.pgpool = ctx.Value(consts.CTX_KEY_DB).(*pgxpool.Pool)
	return pg
}

// возвращает данные о пользователе из БД по ID
func (pg *Postgres) GetUserByID(ctx context.Context, id string) (*User, error) {
	user := new(User)

	errScan := pg.pgpool.QueryRow(ctx, "SELECT * FROM users WHERE user_id=$1", id).Scan(&user.ID, &user.Email, &user.Password, &user.CreationDate)
	if errScan != nil {
		return nil, errScan
	}

	return user, nil
}

type DBClient interface {
	SaveUser(ctx context.Context, user UserData) error
	DeleteByID(ctx context.Context, id string, tablename string) error
	FindUser(ctx context.Context, email string) (string, error)
	FromContext(ctx context.Context) *Postgres
	GetUserByID(ctx context.Context, id string) (*User, error)
}

// структура данных пользователя для регистрации
type UserData struct {
	UserID   string `json:"userid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ud UserData) Validate() error {
	if ud.Email == "" || ud.Password == "" {
		return apperrors.ErrEmptyField
	}
	return nil
}

func (ud UserData) Save() UserData {
	return ud
}

func (ud *UserData) ParseJson(body io.Reader) error {

	errDecode := json.NewDecoder(body).Decode(ud)
	if errDecode != nil {
		return errDecode
	}

	return nil
}

func (ud UserData) GetEmail() string {
	return ud.Email
}

// сравнивает пароли
func (ud UserData) ComparePassword(user User) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ud.Password))
	if err != nil {
		return err
	}
	return nil
}

func (ud *UserData) WriteID(id string) error {
	if id == "" {
		return apperrors.ErrEmptyString
	}
	ud.UserID = id
	return nil
}

type NewUser interface {
	GetEmail() string
	ParseJson(body io.Reader) error
	Save() UserData
	Validate() error
	ComparePassword(User) error
	WriteID(id string) error
}

type tokens struct {
	access  string
	Refresh string
}

func (t *tokens) CreateTokens(id string) error {
	claims := MyCustomClaims{
		UserID: id,
		PairID: uuid.New().String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	secret := os.Getenv(consts.JWT_KEY)
	access, errCreateAccess := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(secret))
	if errCreateAccess != nil {
		return errCreateAccess
	}
	t.access = access

	b := make([]byte, 32)
	_, errRead := rand.Read(b)
	if errRead != nil {
		return errRead
	}
	refresh := base64.URLEncoding.EncodeToString(b)
	t.Refresh = refresh
	return nil
}

func (t tokens) SaveTokens(db DBClient) error {

}

func (t tokens) Sendtokens() {

}

type GenerateTokens interface {
	CreateTokens(id string) error
	SaveTokens() error
	Sendtokens()
}
