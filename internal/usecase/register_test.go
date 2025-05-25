package usecase

import (
	"context"
	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/FooxyS/auth-service/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterUseCase_Execute(t *testing.T) {
	input := struct {
		user domain.User
	}{
		user: domain.User{
			UserID:       "123",
			Email:        "Tim@gmail.com",
			PasswordHash: "some password token hash",
		},
	}

	userInRepo := domain.User{
		UserID:       "123",
		Email:        "Tim@gmail.com",
		PasswordHash: "some password token hash",
	}

	CalledSlice := new([]string)

	expectedSlice := []string{"Exists", "Save"}

	tables := []struct {
		Name       string
		UserRepo   *MockUserRepository
		WantedUser domain.User
		WantErr    error
	}{
		{
			Name: "Exists fails",
			UserRepo: &MockUserRepository{
				existingUser: userInRepo,
				ExistsFail:   true,
			},
			WantErr: ErrExists,
		},
		{
			Name: "user already exists",
			UserRepo: &MockUserRepository{
				existingUser: userInRepo,
			},
			WantedUser: domain.User{},
			WantErr:    apperrors.ErrUserExists,
		},
		{
			Name: "Save fails",
			UserRepo: &MockUserRepository{
				existingUser: domain.User{},
				SaveFail:     true,
			},
			WantedUser: userInRepo,
			WantErr:    ErrSave,
		},
		{
			Name: "happy path",
			UserRepo: &MockUserRepository{
				CalledNeed:  true,
				CalledSlice: CalledSlice,
			},
			WantedUser: userInRepo,
			WantErr:    nil,
		},
	}

	for _, table := range tables {
		t.Run(table.Name, func(t *testing.T) {
			useCase := RegisterUseCase{
				UserRepo: table.UserRepo,
			}

			err := useCase.Execute(context.Background(), input.user)

			assert.ErrorIs(t, err, table.WantErr)

			if table.WantErr == nil {
				assert.Equal(t, table.WantedUser, table.UserRepo.savedUser)

				assert.Equal(t, expectedSlice, *CalledSlice)
			}
		})
	}
}
