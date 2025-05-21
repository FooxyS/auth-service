package usecase

import (
	"context"
	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/FooxyS/auth-service/pkg/apperrors"
)

type RegisterUseCase struct {
	UserRepo domain.UserRepository
}

func (ru RegisterUseCase) Execute(ctx context.Context, user domain.User) error {
	exist, errExists := ru.UserRepo.Exists(ctx, user.Email)
	if errExists != nil {
		return errExists
	}
	if exist {
		return apperrors.ErrUserExists
	}
	errSave := ru.UserRepo.Save(ctx, user)
	if errSave != nil {
		return errSave
	}
	return nil
}
