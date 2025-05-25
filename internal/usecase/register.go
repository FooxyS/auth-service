package usecase

import (
	"context"
	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/FooxyS/auth-service/pkg/apperrors"
)

type RegisterUseCase struct {
	UserRepo domain.UserRepository
	Hasher   domain.PasswordHasher
}

func (ru RegisterUseCase) Execute(ctx context.Context, email, password string) error {
	exist, errExists := ru.UserRepo.Exists(ctx, email)
	if errExists != nil {
		return errExists
	}
	if exist {
		return apperrors.ErrUserExists
	}

	passHash, errHash := ru.Hasher.Hash(password)
	if errHash != nil {
		return errHash
	}

	newUserID, errGenerateID := ru.UserRepo.GenerateUserID()
	if errGenerateID != nil {
		return errGenerateID
	}

	user := domain.User{
		UserID:       newUserID,
		Email:        email,
		PasswordHash: string(passHash),
	}

	errSave := ru.UserRepo.Save(ctx, user)
	if errSave != nil {
		return errSave
	}
	return nil
}
