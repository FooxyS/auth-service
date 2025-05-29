package hasher

import (
	"errors"

	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/FooxyS/auth-service/pkg/apperrors"
	"golang.org/x/crypto/bcrypt"
)

func New() domain.PasswordHasher {
	return &BcryptHasher{}
}

type BcryptHasher struct {
}

func (n BcryptHasher) Compare(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return apperrors.ErrPasswordMismatch
	} else {
		return err
	}
}

func (n BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
