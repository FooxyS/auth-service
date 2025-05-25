package usecase

import (
	"context"
	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoginUseCase_Execute(t *testing.T) {
	input := struct {
		Email    string
		Password string
		Ip       string
		Agent    string
	}{
		Email:    "Tima@gmail.com",
		Password: "some password",
		Ip:       "123.123.123.123",
		Agent:    "some user agent",
	}

	expectedTokenPair := domain.TokenPair{
		AccessToken:  "some access token",
		RefreshToken: "some refresh token",
	}

	expectedSession := domain.Session{
		UserID:      "123",
		IPAddress:   "123.123.123.123",
		RefreshHash: "some refresh token hash",
		PairID:      "some new pairID",
		UserAgent:   "some user agent",
	}

	userInRepo := domain.User{
		UserID:       "123",
		Email:        "Tima@gmail.com",
		PasswordHash: "some password hash",
	}

	CalledSlice := new([]string)

	expectedCalledSlice := []string{"FindByEmail", "Compare", "GeneratePairID", "GenerateAccessToken", "GenerateRefreshToken", "Save"}

	tables := []struct {
		Name            string
		UserRepo        *MockUserRepository
		SessionRepo     *MockSessionRepository
		Tokens          MockTokenService
		Hasher          MockPasswordHasher
		WantedTokenPair domain.TokenPair
		WantedErr       error
	}{
		{
			Name: "FindByEmail fails",
			UserRepo: &MockUserRepository{
				FindByEmailFail: true,
			},
			SessionRepo:     &MockSessionRepository{},
			Tokens:          MockTokenService{},
			Hasher:          MockPasswordHasher{},
			WantedTokenPair: domain.TokenPair{},
			WantedErr:       ErrFind,
		},
		{
			Name: "Compare fails",
			UserRepo: &MockUserRepository{
				existingUser: userInRepo,
			},
			SessionRepo: &MockSessionRepository{},
			Tokens:      MockTokenService{},
			Hasher: MockPasswordHasher{
				CompareFail: true,
			},
			WantedTokenPair: domain.TokenPair{},
			WantedErr:       ErrCompare,
		},
		{
			Name: "GeneratePairID fails",
			UserRepo: &MockUserRepository{
				existingUser: userInRepo,
			},
			SessionRepo: &MockSessionRepository{},
			Tokens: MockTokenService{
				GeneratePairIDFail: true,
			},
			Hasher:          MockPasswordHasher{},
			WantedTokenPair: domain.TokenPair{},
			WantedErr:       ErrGenPairID,
		},
		{
			Name: "GenerateAccessToken fails",
			UserRepo: &MockUserRepository{
				existingUser: userInRepo,
			},
			SessionRepo: &MockSessionRepository{},
			Tokens: MockTokenService{
				GenerateAccessTokenFail: true,
			},
			Hasher:          MockPasswordHasher{},
			WantedTokenPair: domain.TokenPair{},
			WantedErr:       ErrGenAccess,
		},
		{
			Name: "GenerateRefreshToken fails",
			UserRepo: &MockUserRepository{
				existingUser: userInRepo,
			},
			SessionRepo: &MockSessionRepository{},
			Tokens: MockTokenService{
				GenerateRefreshTokenFail: true,
			},
			Hasher:          MockPasswordHasher{},
			WantedTokenPair: domain.TokenPair{},
			WantedErr:       ErrGenRefresh,
		},
		{
			Name: "Save fails",
			UserRepo: &MockUserRepository{
				existingUser: userInRepo,
			},
			SessionRepo: &MockSessionRepository{
				SaveFail: true,
			},
			Tokens:          MockTokenService{},
			Hasher:          MockPasswordHasher{},
			WantedTokenPair: domain.TokenPair{},
			WantedErr:       ErrSave,
		},
		{
			Name: "happy path",
			UserRepo: &MockUserRepository{
				existingUser: userInRepo,
				CalledNeed:   true,
				CalledSlice:  CalledSlice,
			},
			SessionRepo: &MockSessionRepository{
				CalledNeed:  true,
				CalledSlice: CalledSlice,
			},
			Tokens: MockTokenService{
				CalledNeed:  true,
				CalledSlice: CalledSlice,
			},
			Hasher: MockPasswordHasher{
				CalledNeed:  true,
				CalledSlice: CalledSlice,
			},
			WantedTokenPair: domain.TokenPair{},
			WantedErr:       nil,
		},
	}

	for _, table := range tables {
		t.Run(table.Name, func(t *testing.T) {
			useCase := LoginUseCase{
				UserRepo:    table.UserRepo,
				SessionRepo: table.SessionRepo,
				Tokens:      table.Tokens,
				Hasher:      table.Hasher,
			}

			tokenPair, err := useCase.Execute(context.Background(), input.Email, input.Password, input.Ip, input.Agent)

			assert.ErrorIs(t, err, table.WantedErr)

			if table.WantedErr == nil {
				assert.Equal(t, expectedTokenPair, tokenPair)

				assert.Equal(t, expectedSession, table.SessionRepo.SavedSession)

				assert.Equal(t, expectedCalledSlice, *CalledSlice)
			}
		})
	}
}
