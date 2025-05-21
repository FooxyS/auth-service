package usecase

import (
	"context"
	"github.com/FooxyS/auth-service/internal/domain"
)

type MeUseCase struct {
	Tokens   domain.TokenService
	UserRepo domain.UserRepository
}

func (uc MeUseCase) Execute(ctx context.Context, access string) (domain.User, error) {
	userID, _, errValidate := uc.Tokens.ValidateAccessToken(access)
	if errValidate != nil {
		return domain.User{}, errValidate
	}
	user, err := uc.UserRepo.FindByUserID(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
