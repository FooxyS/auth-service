package services

import (
	"context"
	"github.com/FooxyS/auth-service/pkg/apperrors"
	"log"
	"net"
	"net/http"

	"github.com/FooxyS/auth-service/auth/models"
	"github.com/FooxyS/auth-service/pkg/consts"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// достаёт пул подключений к БД из контекста
func GetDBPoolFromContext(r *http.Request) *pgxpool.Pool {
	pgpool := r.Context().Value(consts.CTX_KEY_DB).(*pgxpool.Pool)
	return pgpool
}

// возвращает заголовок из запроса по заданному ключу
func GetValueFromHeader(r *http.Request, key string) string {
	value := r.Header.Get(key)
	return value
}

// проверяет является ли переданная строка пустой
func ValidateStringNotEmpty(s string) error {
	if s == "" {
		return apperrors.ErrEmptyString
	}
	return nil
}

// парсит access токен через секретную строку, возвращает его claims
func ParseAccesstoken(access string) (*models.MyCustomClaims, *jwt.Token, error) {
	accessClaims := new(models.MyCustomClaims)

	token, errWithParse := jwt.ParseWithClaims(access, accessClaims, func(t *jwt.Token) (interface{}, error) {
		return GetFromEnv(consts.JWT_KEY)
	})

	if errWithParse != nil {
		return nil, nil, errWithParse
	}

	return accessClaims, token, nil
}

// проверяет истёк ли access токен
func ValidateAccessToken(token *jwt.Token) bool {
	return token.Valid
}

// обращается к базе данных и возвращает структуру с данными о сессии пользователя по user_id
func GetSessionByUserID(ctx context.Context, pgpool *pgxpool.Pool, id string) (*models.Session, error) {
	session := new(models.Session)

	errWithScan := pgpool.QueryRow(ctx, "SELECT * FROM session_table WHERE user_id=$1", id).Scan(&session.ID, &session.IP, &session.RefreshToken, &session.PairID, &session.UserAgent)
	if errWithScan != nil {
		return nil, errWithScan
	}

	return session, nil
}

// принимает remoteAddr и возвращает ip клиента
func GetClientIPFromRemoteAddr(r *http.Request) (string, error) {
	host, _, errSplitHost := net.SplitHostPort(r.RemoteAddr)
	if errSplitHost != nil {
		return "", errSplitHost
	}
	return host, nil
}

// сравнивает ip клиента с тем, что хранится в базе сессий
func CompareIP(host, ip string) error {
	if host != ip {
		return apperrors.ErrNotMatch
	}
	return nil
}

// сравнивает userAgent
func CompareUserAgent(agent, agentDB string) error {
	if agent != agentDB {
		return apperrors.ErrNotMatch
	}
	return nil
}

func UnauthorizedResponseWithReason(w http.ResponseWriter, err error) {
	log.Printf("Unauthorized user: %v\n", err)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
