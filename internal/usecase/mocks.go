package usecase

import (
	"context"
	"errors"
	"github.com/FooxyS/auth-service/internal/domain"
)

var (
	ErrExists = errors.New("error while checking existence of user")
	ErrSave   = errors.New("error while writing of user")
)

type MockUserRepository struct {
	existingUser domain.User
	savedUser    domain.User
	ExistsFail   bool
	SaveFail     bool
}

func (m MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	if m.ExistsFail {
		return false, ErrExists
	}

	if email == m.existingUser.Email {
		return true, nil
	}

	return false, nil
}

func (m *MockUserRepository) Save(ctx context.Context, user domain.User) error {
	if m.SaveFail {
		return ErrSave
	}
	m.savedUser = user
	return nil
}

func (m MockUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockUserRepository) FindByUserID(ctx context.Context, id string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

type MockSessionRepository struct {
}

func (m MockSessionRepository) Save(ctx context.Context, session domain.Session) error {
	//TODO implement me
	panic("implement me")
}

func (m MockSessionRepository) Delete(ctx context.Context, pairID string) error {
	//TODO implement me
	panic("implement me")
}

func (m MockSessionRepository) UpdateSession(ctx context.Context, oldPair, pair, refreshHash string) error {
	//TODO implement me
	panic("implement me")
}

func (m MockSessionRepository) FindByPairID(ctx context.Context, pairID string) (domain.Session, error) {
	//TODO implement me
	panic("implement me")
}

type MockTokenService struct {
}

func (m MockTokenService) GenerateAccessToken(id string, pairID string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockTokenService) GenerateRefreshToken() (string, string, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockTokenService) GeneratePairID() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockTokenService) ValidateAccessToken(access string) (string, string, error) {
	//TODO implement me
	panic("implement me")
}

type MockPasswordHasher struct {
}

func (m MockPasswordHasher) Compare(hash, password string) error {
	//TODO implement me
	panic("implement me")
}
