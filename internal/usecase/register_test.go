package usecase

import (
	"context"
	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/FooxyS/auth-service/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterUseCase_Execute(t *testing.T) {
	existUser := domain.User{
		UserID:       "4321",
		Email:        "Arthur@gmail.com",
		PasswordHash: "rdr2",
	}

	notExistUser := domain.User{
		UserID:       "1234",
		Email:        "Tim@gmail.com",
		PasswordHash: "qwerty12345",
	}

	tables := []struct {
		name         string
		mockUserRepo *MockUserRepository
		user         domain.User
		wantErr      error
	}{
		{
			name: "success",
			mockUserRepo: &MockUserRepository{
				existingUser: existUser,
				ExistsFail:   false,
				SaveFail:     false,
			},
			user:    notExistUser,
			wantErr: nil,
		},
		{
			name: "user exists",
			mockUserRepo: &MockUserRepository{
				existingUser: existUser,
				ExistsFail:   false,
				SaveFail:     false,
			},
			user:    existUser,
			wantErr: apperrors.ErrUserExists,
		},
		{
			name: "exists fail",
			mockUserRepo: &MockUserRepository{
				existingUser: existUser,
				ExistsFail:   true,
				SaveFail:     false,
			},
			wantErr: ErrExists,
		},
		{
			name: "save fail",
			mockUserRepo: &MockUserRepository{
				existingUser: existUser,
				ExistsFail:   false,
				SaveFail:     true,
			},
			wantErr: ErrSave,
		},
	}

	for _, table := range tables {
		t.Run(table.name, func(t *testing.T) {
			useCase := RegisterUseCase{UserRepo: table.mockUserRepo}

			err := useCase.Execute(context.Background(), table.user)

			if !assert.ErrorIs(t, err, table.wantErr) {
				t.Errorf("Execute() error = %v, wantErr %v", err, table.wantErr)
			}

			if table.wantErr == nil {
				if !assert.Equal(t, table.user, table.mockUserRepo.savedUser) {
					t.Errorf("got %v, want %v", table.mockUserRepo.savedUser, table.user)
				}
			}
		})
	}

}
