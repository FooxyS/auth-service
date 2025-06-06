package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/FooxyS/auth-service/pkg/apperrors"
	"github.com/FooxyS/auth-service/pkg/consts"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func New() domain.TokenService {
	return &JWTService{}
}

type MyCustomClaims struct {
	UserID string
	PairID string
	jwt.RegisteredClaims
}

type JWTService struct {
}

func (j *JWTService) GenerateAccessToken(id string, pairID string) (string, error) {
	newClaims := MyCustomClaims{
		UserID: id,
		PairID: pairID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	secretString := os.Getenv(consts.JWT_KEY)

	newAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, newClaims).SignedString([]byte(secretString))
	if err != nil {
		return "", err
	}
	return newAccessToken, nil
}

func (j *JWTService) GenerateRefreshToken() (string, string, error) {
	b := make([]byte, 32)
	if _, errRead := rand.Read(b); errRead != nil {
		return "", "", errRead
	}

	refresh := base64.URLEncoding.EncodeToString(b)

	refreshHash, errHash := bcrypt.GenerateFromPassword([]byte(refresh), bcrypt.DefaultCost)
	if errHash != nil {
		return "", "", errHash
	}
	return refresh, string(refreshHash), nil
}

func (j *JWTService) GeneratePairID() (string, error) {
	return uuid.New().String(), nil
}

func (j *JWTService) ValidateAccessToken(access string) (string, string, error) {
	newClaims := new(MyCustomClaims)

	s := strings.Split(access, " ")
	if len(s) != 2 {
		return "", "", apperrors.ErrBearer
	}

	_, errParse := jwt.ParseWithClaims(s[1], newClaims, func(t *jwt.Token) (interface{}, error) {
		secret := os.Getenv(consts.JWT_KEY)
		return []byte(secret), nil
	})
	if errParse != nil && !errors.Is(errParse, jwt.ErrTokenExpired) {
		return "", "", errParse
	}

	return newClaims.UserID, newClaims.PairID, nil
}
