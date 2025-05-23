package usecase

import (
	"context"
	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMeUseCase_Execute(t *testing.T) {
	expectedUser := domain.User{
		UserID:       "12345",
		Email:        "Tim@gmail.com",
		PasswordHash: "TheBestProgrammer",
	}

	tables := []struct {
		Name        string
		Input       string
		Tokens      domain.TokenService
		UserRepo    domain.UserRepository
		WantedUser  domain.User
		WantedError error
	}{
		{
			Name:  "ValidateAccess() fails",
			Input: "some access token",
			Tokens: MockTokenService{
				userID:                  "12345",
				PairID:                  "76913857",
				ValidateAccessTokenFail: true,
			},
			UserRepo:    &MockUserRepository{},
			WantedUser:  domain.User{},
			WantedError: ErrValidateAccess,
		},
		{
			Name:  "FindByUserID() fails",
			Input: "some access token",
			Tokens: MockTokenService{
				userID: "12345",
				PairID: "76913857",
			},
			UserRepo: &MockUserRepository{
				user:             domain.User{},
				FindByUserIDFail: true,
			},
			WantedUser:  domain.User{},
			WantedError: ErrFind,
		},
		{
			Name:  "success",
			Input: "some access token",
			Tokens: MockTokenService{
				userID: "12345",
				PairID: "76913857",
			},
			UserRepo: &MockUserRepository{
				user: expectedUser,
			},
			WantedUser:  expectedUser,
			WantedError: nil,
		},
	}

	for _, table := range tables {
		t.Run(table.Name, func(t *testing.T) {
			useCase := MeUseCase{Tokens: table.Tokens, UserRepo: table.UserRepo}

			user, err := useCase.Execute(context.Background(), table.Input)

			assert.ErrorIs(t, err, table.WantedError)

			assert.Equal(t, table.WantedUser, user)
		})
	}
}
