package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"time"

	"github.com/FooxyS/auth-service/pkg/apperrors"
	"github.com/FooxyS/auth-service/pkg/consts"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

	newAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, newClaims).SignedString(os.Getenv(consts.JWT_KEY))
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

	token, errParse := jwt.ParseWithClaims(access, newClaims, func(t *jwt.Token) (interface{}, error) {
		secret := os.Getenv(consts.JWT_KEY)
		return []byte(secret), nil
	})
	if errParse != nil {
		return "", "", errParse
	}

	if !token.Valid {
		return newClaims.UserID, newClaims.PairID, apperrors.ErrNotValid
	}

	return newClaims.UserID, newClaims.PairID, nil
}
