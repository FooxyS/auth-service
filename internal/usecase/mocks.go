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
	ErrGenUserID      = errors.New("error while generating user id")
	ErrCompare        = errors.New("error while compare")
	ErrFind           = errors.New("error while find")
	ErrUpdate         = errors.New("error while updating")
	ErrHash           = errors.New("error while hashing password")
)

type MockUserRepository struct {
	//для loginUseCase
	user domain.User

	//для registerUseCase
	existingUser       domain.User
	savedUser          domain.User
	ExistsFail         bool
	SaveFail           bool
	FindByUserIDFail   bool
	FindByEmailFail    bool
	CalledNeed         bool
	GenerateUserIDFail bool
	CalledSlice        *[]string
}

func (m *MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	if m.ExistsFail {
		return false, ErrExists
	}

	if email == m.existingUser.Email {
		return true, nil
	}

	if m.CalledNeed && m.CalledSlice != nil {
		*m.CalledSlice = append(*m.CalledSlice, "Exists")
	}

	return false, nil
}

func (m *MockUserRepository) Save(ctx context.Context, user domain.User) error {
	if m.SaveFail {
		m.savedUser = domain.User{}
		return ErrSave
	}
	m.savedUser = user

	if m.CalledNeed && m.CalledSlice != nil {
		*m.CalledSlice = append(*m.CalledSlice, "Save")
	}

	return nil

}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	if m.FindByEmailFail {
		return domain.User{}, ErrFind
	}

	if m.CalledNeed && m.CalledSlice != nil {
		*m.CalledSlice = append(*m.CalledSlice, "FindByEmail")
	}

	return m.existingUser, nil

}

func (m *MockUserRepository) FindByUserID(ctx context.Context, id string) (domain.User, error) {
	if m.FindByUserIDFail {
		return domain.User{}, ErrFind
	}

	if m.CalledNeed && m.CalledSlice != nil {
		*m.CalledSlice = append(*m.CalledSlice, "FindByUserID")
	}

	return m.user, nil
}

func (m *MockUserRepository) GenerateUserID() (string, error) {
	if m.GenerateUserIDFail {
		return "", ErrGenUserID
	}

	if m.CalledNeed && m.CalledSlice != nil {
		*m.CalledSlice = append(*m.CalledSlice, "GenerateUserID")
	}

	return "some new userID", nil
}

type MockSessionRepository struct {
	SavedSession      domain.Session
	UpdatedSession    domain.Session
	Session           domain.Session
	DeletedSession    domain.Session
	SessionForDelete  domain.Session
	UpdateSessionFail bool
	DeleteFail        bool
	SaveFail          bool
	FindByPairIDFail  bool
	CalledNeed        bool
	CalledSlice       *[]string
}

func (m *MockSessionRepository) Save(ctx context.Context, session domain.Session) error {
	if m.SaveFail {
		return ErrSave
	}
	m.SavedSession = session

	if m.CalledNeed && m.CalledSlice != nil {
		*m.CalledSlice = append(*m.CalledSlice, "Save")
	}

	return nil

}

func (m *MockSessionRepository) Delete(ctx context.Context, pairID string) error {
	if m.DeleteFail {
		return ErrDelete
	}
	m.DeletedSession = m.SessionForDelete

	if m.CalledNeed && m.CalledSlice != nil {
		*m.CalledSlice = append(*m.CalledSlice, "Delete")
	}

	return nil

}

func (m *MockSessionRepository) UpdateSession(ctx context.Context, oldPair, pair, refreshHash string) error {
	if m.UpdateSessionFail {
		return ErrUpdate
	}
	m.UpdatedSession = m.Session

	m.UpdatedSession.PairID = pair
	m.UpdatedSession.RefreshHash = refreshHash

	if m.CalledNeed && m.CalledSlice != nil {
		*m.CalledSlice = append(*m.CalledSlice, "UpdateSession")
	}

	return nil
}

func (m *MockSessionRepository) FindByPairID(ctx context.Context, pairID string) (domain.Session, error) {
	if m.FindByPairIDFail {
		return domain.Session{}, ErrFind
	}

	if m.CalledNeed && m.CalledSlice != nil {
		*m.CalledSlice = append(*m.CalledSlice, "FindByPairID")
	}

	return m.Session, nil
}

type MockTokenService struct {
	userID                   string
	PairID                   string
	GenerateAccessTokenFail  bool
	GenerateRefreshTokenFail bool
	GeneratePairIDFail       bool
	ValidateAccessTokenFail  bool
	CalledNeed               bool
	CalledSlice              *[]string
}

func (m MockTokenService) GenerateAccessToken(id string, pairID string) (string, error) {
	if m.GenerateAccessTokenFail {
		return "", ErrGenAccess
	}

	if m.CalledNeed && m.CalledSlice != nil {
		*m.CalledSlice = append(*m.CalledSlice, "GenerateAccessToken")
	}

	return "some access token", nil
}

func (m MockTokenService) GenerateRefreshToken() (string, string, error) {
	if m.GenerateRefreshTokenFail {
		return "", "", ErrGenRefresh
	}

	if m.CalledNeed && m.CalledSlice != nil {
		*m.CalledSlice = append(*m.CalledSlice, "GenerateRefreshToken")
	}

	return "some refresh token", "some refresh token hash", nil
}

func (m MockTokenService) GeneratePairID() (string, error) {
	if m.GeneratePairIDFail {
		return "", ErrGenPairID
	}

	if m.CalledNeed && m.CalledSlice != nil {
		*m.CalledSlice = append(*m.CalledSlice, "GeneratePairID")
	}

	return "some new pairID", nil
}

func (m MockTokenService) ValidateAccessToken(access string) (string, string, error) {
	if m.ValidateAccessTokenFail {
		return "", "", ErrValidateAccess
	}

	if m.CalledNeed && m.CalledSlice != nil {
		*m.CalledSlice = append(*m.CalledSlice, "ValidateAccessToken")
	}

	return m.userID, m.PairID, nil
}

type MockPasswordHasher struct {
	CompareFail bool
	HashFail    bool
	CalledNeed  bool
	CalledSlice *[]string
}

func (m MockPasswordHasher) Compare(hash, password string) error {
	if m.CompareFail {
		return ErrCompare
	}

	if m.CalledNeed && m.CalledSlice != nil {
		*m.CalledSlice = append(*m.CalledSlice, "Compare")
	}

	return nil
}

func (m MockPasswordHasher) Hash(password string) (string, error) {
	if m.HashFail {
		return "", ErrHash
	}

	if m.CalledNeed && m.CalledSlice != nil {
		*m.CalledSlice = append(*m.CalledSlice, "Hash")
	}

	return "some password hash", nil
}
