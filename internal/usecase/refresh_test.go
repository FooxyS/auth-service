package usecase

import (
	"context"
	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/FooxyS/auth-service/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRefreshUseCase_Execute(t *testing.T) {
	input := struct {
		access, refresh, ip, agent string
	}{
		access:  "some access token",
		refresh: "some refresh token",
		ip:      "123.123.123.123",
		agent:   "some user agent",
	}

	expectedUpdateSession := domain.Session{
		UserID:      "123",
		IPAddress:   "123.123.123.123",
		RefreshHash: "some refresh token hash",
		PairID:      "some new pairID",
		UserAgent:   "some user agent",
	}

	SessionInRepo := domain.Session{
		UserID:      "123",
		IPAddress:   "123.123.123.123",
		RefreshHash: "some refresh token",
		PairID:      "some pairID",
		UserAgent:   "some user agent",
	}

	BadIPSessionInRepo := domain.Session{
		UserID:      "123",
		IPAddress:   "100.100.100.100",
		RefreshHash: "some refresh token",
		PairID:      "some new pairID",
		UserAgent:   "some user agent",
	}

	BadAgentSessionInRepo := domain.Session{
		UserID:      "123",
		IPAddress:   "123.123.123.123",
		RefreshHash: "some refresh token",
		PairID:      "some new pairID",
		UserAgent:   "kajbefkjaebfk",
	}

	CalledSlice := new([]string)

	ExpectedCalled := []string{"ValidateAccessToken", "FindByPairID", "Compare", "GeneratePairID", "GenerateAccessToken", "GenerateRefreshToken", "UpdateSession"}

	tables := []struct {
		Name            string
		Tokens          MockTokenService
		SessionRepo     *MockSessionRepository
		Hasher          MockPasswordHasher
		WantedTokenPair domain.TokenPair
		WantedError     error
	}{
		{
			Name: "ValidateAccessToken fails",
			Tokens: MockTokenService{
				ValidateAccessTokenFail: true,
			},
			SessionRepo:     &MockSessionRepository{},
			Hasher:          MockPasswordHasher{},
			WantedTokenPair: domain.TokenPair{},
			WantedError:     ErrValidateAccess,
		},
		{
			Name:   "FindByPairID fails",
			Tokens: MockTokenService{},
			SessionRepo: &MockSessionRepository{
				FindByPairIDFail: true,
			},
			Hasher:          MockPasswordHasher{},
			WantedTokenPair: domain.TokenPair{},
			WantedError:     ErrFind,
		},
		{
			Name:        "Compare fails",
			Tokens:      MockTokenService{},
			SessionRepo: &MockSessionRepository{},
			Hasher: MockPasswordHasher{
				CompareFail: true,
			},
			WantedTokenPair: domain.TokenPair{},
			WantedError:     ErrCompare,
		},
		{
			Name:   "mismatch IP",
			Tokens: MockTokenService{},
			SessionRepo: &MockSessionRepository{
				Session: BadIPSessionInRepo,
			},
			Hasher:          MockPasswordHasher{},
			WantedTokenPair: domain.TokenPair{},
			WantedError:     apperrors.ErrIPMismatch,
		},
		{
			Name:   "mismatch User-Agent, delete fails",
			Tokens: MockTokenService{},
			SessionRepo: &MockSessionRepository{
				Session:    BadAgentSessionInRepo,
				DeleteFail: true,
			},
			Hasher:          MockPasswordHasher{},
			WantedTokenPair: domain.TokenPair{},
			WantedError:     ErrDelete,
		},
		{
			Name:   "mismatch User-Agent",
			Tokens: MockTokenService{},
			SessionRepo: &MockSessionRepository{
				Session: BadAgentSessionInRepo,
			},
			Hasher:          MockPasswordHasher{},
			WantedTokenPair: domain.TokenPair{},
			WantedError:     apperrors.ErrAgentMismatch,
		},
		{
			Name: "GeneratePairID fails",
			Tokens: MockTokenService{
				GeneratePairIDFail: true,
			},
			SessionRepo: &MockSessionRepository{
				Session: SessionInRepo,
			},
			Hasher:          MockPasswordHasher{},
			WantedTokenPair: domain.TokenPair{},
			WantedError:     ErrGenPairID,
		},
		{
			Name: "GenerateAccessToken fails",
			Tokens: MockTokenService{
				GenerateAccessTokenFail: true,
			},
			SessionRepo: &MockSessionRepository{
				Session: SessionInRepo,
			},
			Hasher:          MockPasswordHasher{},
			WantedTokenPair: domain.TokenPair{},
			WantedError:     ErrGenAccess,
		},
		{
			Name: "GenerateRefreshToken fails",
			Tokens: MockTokenService{
				GenerateRefreshTokenFail: true,
			},
			SessionRepo: &MockSessionRepository{
				Session: SessionInRepo,
			},
			Hasher:          MockPasswordHasher{},
			WantedTokenPair: domain.TokenPair{},
			WantedError:     ErrGenRefresh,
		},
		{
			Name:   "UpdateSession fails",
			Tokens: MockTokenService{},
			SessionRepo: &MockSessionRepository{
				Session:           SessionInRepo,
				UpdateSessionFail: true,
			},
			Hasher:          MockPasswordHasher{},
			WantedTokenPair: domain.TokenPair{},
			WantedError:     ErrUpdate,
		},
		{
			Name: "happy path",
			Tokens: MockTokenService{
				CalledNeed:  true,
				CalledSlice: CalledSlice,
			},
			SessionRepo: &MockSessionRepository{
				Session:     SessionInRepo,
				CalledNeed:  true,
				CalledSlice: CalledSlice,
			},
			Hasher: MockPasswordHasher{
				CalledNeed:  true,
				CalledSlice: CalledSlice,
			},
			WantedTokenPair: domain.TokenPair{
				AccessToken:  "some access token",
				RefreshToken: "some refresh token",
			},
			WantedError: nil,
		},
	}

	for _, table := range tables {
		t.Run(table.Name, func(t *testing.T) {
			useCase := RefreshUseCase{
				Tokens:      table.Tokens,
				SessionRepo: table.SessionRepo,
				Hasher:      table.Hasher,
			}

			tokenPair, err := useCase.Execute(context.Background(), input.access, input.refresh, input.ip, input.agent)

			assert.ErrorIs(t, err, table.WantedError)

			if table.WantedError == nil {
				assert.Equal(t, table.WantedTokenPair, tokenPair)

				assert.Equal(t, expectedUpdateSession, table.SessionRepo.UpdatedSession)

				assert.Equal(t, ExpectedCalled, *CalledSlice)
			}
		})
	}
}
