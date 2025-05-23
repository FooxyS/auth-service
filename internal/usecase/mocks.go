package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/FooxyS/auth-service/internal/domain"
)

var (
	ErrExists         = errors.New("error while checking existence of user")
	ErrSave           = errors.New("error while saving")
	ErrGenAccess      = errors.New("error while generating access token")
	ErrValidateAccess = errors.New("error while validating access token")
	ErrGenRefresh     = errors.New("error while generating refresh token")
	ErrGenPairID      = errors.New("error while generating pair id")
	ErrCompare        = errors.New("error while compare")
	ErrFind           = errors.New("error while find")
)

type MockUserRepository struct {
	//для loginUseCase
	user domain.User

	//для registerUseCase
	existingUser     domain.User
	savedUser        domain.User
	ExistsFail       bool
	SaveFail         bool
	FindByUserIDFail bool
}

func (m *MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
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

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	if m.user.Email == email {
		return m.user, nil
	}
	return domain.User{}, nil
}

func (m *MockUserRepository) FindByUserID(ctx context.Context, id string) (domain.User, error) {
	if m.FindByUserIDFail {
		return domain.User{}, ErrFind
	}
	return m.user, nil
}

type MockSessionRepository struct {
	SavedSession domain.Session
	SaveFail     bool
}

func (m *MockSessionRepository) Save(ctx context.Context, session domain.Session) error {
	if m.SaveFail {
		return ErrSave
	}
	m.SavedSession = session
	return nil
}

func (m *MockSessionRepository) Delete(ctx context.Context, pairID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockSessionRepository) UpdateSession(ctx context.Context, oldPair, pair, refreshHash string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockSessionRepository) FindByPairID(ctx context.Context, pairID string) (domain.Session, error) {
	//TODO implement me
	panic("implement me")
}

type MockTokenService struct {
	userID                   string
	PairID                   string
	GenerateAccessTokenFail  bool
	GenerateRefreshTokenFail bool
	GeneratePairIDFail       bool
	ValidateAccessTokenFail  bool
}

func (m MockTokenService) GenerateAccessToken(id string, pairID string) (string, error) {
	if m.GenerateAccessTokenFail {
		return "", ErrGenAccess
	}
	access := fmt.Sprintf("access_token_%s_%s", id, pairID)
	return access, nil
}

func (m MockTokenService) GenerateRefreshToken() (string, string, error) {
	if m.GenerateRefreshTokenFail {
		return "", "", ErrGenRefresh
	}
	return "refresh_token", "refresh_token_hash", nil
}

func (m MockTokenService) GeneratePairID() (string, error) {
	if m.GeneratePairIDFail {
		return "", ErrGenPairID
	}
	return "pair_id", nil
}

func (m MockTokenService) ValidateAccessToken(access string) (string, string, error) {
	if m.ValidateAccessTokenFail {
		return "", "", ErrValidateAccess
	}
	return m.userID, m.PairID, nil
}

type MockPasswordHasher struct {
	CompareFail bool
}

func (m MockPasswordHasher) Compare(hash, password string) error {
	if m.CompareFail {
		return ErrCompare
	}
	return nil
}
