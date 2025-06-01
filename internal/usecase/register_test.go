package usecase

import (
	"context"
	"testing"

	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/FooxyS/auth-service/pkg/apperrors"
	"github.com/stretchr/testify/assert"
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

	expectedUser := domain.User{
		UserID:       "some new userID",
		Email:        "Tim@gmail.com",
		PasswordHash: "some password hash",
	}

	CalledSlice := new([]string)

	expectedSlice := []string{"Exists", "Hash", "GenerateUserID", "Save"}

	tables := []struct {
		Name       string
		UserRepo   *MockUserRepository
		Hasher     *MockPasswordHasher
		WantedUser domain.User
		WantErr    error
	}{
		{
			Name: "Exists fails",
			UserRepo: &MockUserRepository{
				existingUser: userInRepo,
				ExistsFail:   true,
				CalledNeed:   false,
			},
			Hasher:  &MockPasswordHasher{},
			WantErr: ErrExists,
		},
		{
			Name: "Hash fails",
			UserRepo: &MockUserRepository{
				existingUser: domain.User{},
			},
			Hasher: &MockPasswordHasher{
				HashFail: true,
			},
			WantErr: ErrHash,
		},
		{
			Name: "user already exists",
			UserRepo: &MockUserRepository{
				existingUser: userInRepo,
				CalledNeed:   false,
			},
			Hasher:     &MockPasswordHasher{},
			WantedUser: domain.User{},
			WantErr:    apperrors.ErrUserExists,
		},
		{
			Name: "Save fails",
			UserRepo: &MockUserRepository{
				existingUser: domain.User{},
				SaveFail:     true,
				CalledNeed:   false,
			},
			Hasher:     &MockPasswordHasher{},
			WantedUser: domain.User{},
			WantErr:    ErrSave,
		},
		{
			Name: "happy path",
			UserRepo: &MockUserRepository{
				CalledNeed:  true,
				CalledSlice: CalledSlice,
			},
			Hasher: &MockPasswordHasher{
				CalledNeed:  true,
				CalledSlice: CalledSlice,
			},
			WantedUser: expectedUser,
			WantErr:    nil,
		},
	}

	for _, table := range tables {
		t.Run(table.Name, func(t *testing.T) {
			useCase := RegisterUseCase{
				UserRepo: table.UserRepo,
				Hasher:   table.Hasher,
			}

			err := useCase.Execute(context.Background(), input.user.Email, input.user.PasswordHash)

			assert.ErrorIs(t, err, table.WantErr)

			if table.WantErr == nil {
				assert.Equal(t, table.WantedUser, table.UserRepo.savedUser)

				assert.Equal(t, expectedSlice, *CalledSlice)
			}
		})
	}
}
