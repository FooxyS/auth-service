package hasher

import (
	"github.com/FooxyS/auth-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func New() domain.PasswordHasher {
	return &BcryptHasher{}
}

type BcryptHasher struct {
}

func (n BcryptHasher) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (n BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
