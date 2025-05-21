package usecase

import (
	"context"
	"github.com/FooxyS/auth-service/internal/domain"
)

type LogoutUseCase struct {
	SessionRepo domain.SessionRepository
	Tokens      domain.TokenService
}

func (uc *LogoutUseCase) Execute(ctx context.Context, access string) error {
	_, pairID, errValidate := uc.Tokens.ValidateAccessToken(access)
	if errValidate != nil {
		return errValidate
	}
	if err := uc.SessionRepo.Delete(ctx, pairID); err != nil {
		return err
	}
	return nil
}
