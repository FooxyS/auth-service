package usecase

import (
	"context"
	"errors"
	"github.com/FooxyS/auth-service/internal/domain"
)

var (
	ErrExists         = errors.New("error while checking existence of user")
	ErrSave           = errors.New("error while saving")
	ErrDelete         = errors.New("error while deleting")
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
	FindByEmailFail  bool
	CalledSlice      *[]string
}

func (m *MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	if m.ExistsFail {
		return false, ErrExists
	}

	if email == m.existingUser.Email {
		return true, nil
	}

	*m.CalledSlice = append(*m.CalledSlice, "Exists")

	return false, nil
}

func (m *MockUserRepository) Save(ctx context.Context, user domain.User) error {
	if m.SaveFail {
		return ErrSave
	}
	m.savedUser = user

	*m.CalledSlice = append(*m.CalledSlice, "Save")

	return nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	if m.FindByEmailFail {
		return domain.User{}, ErrFind
	}

	*m.CalledSlice = append(*m.CalledSlice, "FindByEmail")

	return m.existingUser, nil
}

func (m *MockUserRepository) FindByUserID(ctx context.Context, id string) (domain.User, error) {
	if m.FindByUserIDFail {
		return domain.User{}, ErrFind
	}

	*m.CalledSlice = append(*m.CalledSlice, "FindByUserID")

	return m.user, nil
}

type MockSessionRepository struct {
	SavedSession     domain.Session
	DeletedSession   domain.Session
	SessionForDelete domain.Session
	DeleteFail       bool
	SaveFail         bool
	CalledSlice      *[]string
}

func (m *MockSessionRepository) Save(ctx context.Context, session domain.Session) error {
	if m.SaveFail {
		return ErrSave
	}
	m.SavedSession = session

	*m.CalledSlice = append(*m.CalledSlice, "Save")

	return nil
}

func (m *MockSessionRepository) Delete(ctx context.Context, pairID string) error {
	if m.DeleteFail {
		return ErrDelete
	}
	m.DeletedSession = m.SessionForDelete

	*m.CalledSlice = append(*m.CalledSlice, "Delete")

	return nil
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
	CalledSlice              *[]string
}

func (m MockTokenService) GenerateAccessToken(id string, pairID string) (string, error) {
	if m.GenerateAccessTokenFail {
		return "", ErrGenAccess
	}

	*m.CalledSlice = append(*m.CalledSlice, "GenerateAccessToken")

	return "some access token", nil
}

func (m MockTokenService) GenerateRefreshToken() (string, string, error) {
	if m.GenerateRefreshTokenFail {
		return "", "", ErrGenRefresh
	}

	*m.CalledSlice = append(*m.CalledSlice, "GenerateRefreshToken")

	return "some refresh token", "some refresh token hash", nil
}

func (m MockTokenService) GeneratePairID() (string, error) {
	if m.GeneratePairIDFail {
		return "", ErrGenPairID
	}

	*m.CalledSlice = append(*m.CalledSlice, "GeneratePairID")

	return "some new pairID", nil
}

func (m MockTokenService) ValidateAccessToken(access string) (string, string, error) {
	if m.ValidateAccessTokenFail {
		return "", "", ErrValidateAccess
	}

	*m.CalledSlice = append(*m.CalledSlice, "ValidateAccessToken")

	return m.userID, m.PairID, nil
}

type MockPasswordHasher struct {
	CompareFail bool
	CalledSlice *[]string
}

func (m MockPasswordHasher) Compare(hash, password string) error {
	if m.CompareFail {
		return ErrCompare
	}

	*m.CalledSlice = append(*m.CalledSlice, "Compare")

	return nil
}
