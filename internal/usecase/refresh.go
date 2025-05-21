package usecase

import (
	"context"
	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/FooxyS/auth-service/pkg/apperrors"
	"log"
)

type RefreshUseCase struct {
	Tokens      domain.TokenService
	SessionRepo domain.SessionRepository
	Hasher      domain.PasswordHasher
}

func (ru RefreshUseCase) Execute(ctx context.Context, access, refresh, ip, agent string) (domain.TokenPair, error) {
	//проверить текущие токены
	userID, pairID, errValidate := ru.Tokens.ValidateAccessToken(access)
	if errValidate != nil {
		return domain.TokenPair{}, errValidate
	}
	session, errSearch := ru.SessionRepo.FindByPairID(ctx, pairID)
	if errSearch != nil {
		return domain.TokenPair{}, errSearch
	}
	if err := ru.Hasher.Compare(session.RefreshHash, refresh); err != nil {
		return domain.TokenPair{}, err
	}
	if ip != session.IPAddress {
		return domain.TokenPair{}, apperrors.ErrIPMismatch
	}
	if agent != session.UserAgent {
		if err := ru.SessionRepo.Delete(ctx, pairID); err != nil {
			log.Printf("error with deleting session: %v", err)
			return domain.TokenPair{}, err
		}
		return domain.TokenPair{}, apperrors.ErrAgentMismatch
	}

	//генерация новых токенов
	newPairID, errPairID := ru.Tokens.GeneratePairID()
	if errPairID != nil {
		return domain.TokenPair{}, errPairID
	}

	newAccess, errAccess := ru.Tokens.GenerateAccessToken(userID, newPairID)
	if errAccess != nil {
		return domain.TokenPair{}, errAccess
	}

	newRefresh, newRefreshHash, errRefresh := ru.Tokens.GenerateRefreshToken()
	if errRefresh != nil {
		return domain.TokenPair{}, errRefresh
	}

	//обновление базы
	if err := ru.SessionRepo.UpdateSession(ctx, pairID, newPairID, newRefreshHash); err != nil {
		return domain.TokenPair{}, err
	}

	//возвращение пары
	return domain.TokenPair{AccessToken: newAccess, RefreshToken: newRefresh}, nil
}
