package usecase

import (
	"context"
	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogoutUseCase_Execute(t *testing.T) {
	input := "some access token"
	sessionInRepo := domain.Session{
		UserID:      "123",
		IPAddress:   "123.123.123.123",
		RefreshHash: "some refresh hash",
		PairID:      "123",
		UserAgent:   "some user agent",
	}

	tables := []struct {
		Name          string
		Tokens        MockTokenService
		SessionRepo   *MockSessionRepository
		WantedSession domain.Session
		WantedErr     error
	}{
		{
			Name: "Validate fails",
			Tokens: MockTokenService{
				PairID:                  "123",
				ValidateAccessTokenFail: true,
			},
			SessionRepo: &MockSessionRepository{
				SessionForDelete: sessionInRepo,
			},
			WantedSession: domain.Session{},
			WantedErr:     ErrValidateAccess,
		},
		{
			Name: "Delete fails",
			Tokens: MockTokenService{
				PairID: "123",
			},
			SessionRepo: &MockSessionRepository{
				SessionForDelete: sessionInRepo,
				DeleteFail:       true,
			},
			WantedSession: domain.Session{},
			WantedErr:     ErrDelete,
		},
		{
			Name: "success",
			Tokens: MockTokenService{
				PairID: "123",
			},
			SessionRepo: &MockSessionRepository{
				SessionForDelete: sessionInRepo,
			},
			WantedSession: sessionInRepo,
			WantedErr:     nil,
		},
	}

	for _, table := range tables {
		t.Run(table.Name, func(t *testing.T) {
			useCase := LogoutUseCase{
				SessionRepo: table.SessionRepo,
				Tokens:      table.Tokens,
			}

			err := useCase.Execute(context.Background(), input)

			assert.ErrorIs(t, err, table.WantedErr)

			if table.WantedErr == nil {
				assert.Equal(t, table.WantedSession, table.SessionRepo.DeletedSession)
			}
		})
	}
}
