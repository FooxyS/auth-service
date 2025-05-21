package usecase

import (
	"context"
	"github.com/FooxyS/auth-service/internal/domain"
)

type LoginUseCase struct {
	UserRepo    domain.UserRepository
	SessionRepo domain.SessionRepository
	tokens      domain.TokenService
	hasher      domain.PasswordHasher
}

func (uc *LoginUseCase) Execute(ctx context.Context, email, password, ip, agent string) (domain.TokenPair, error) {
	foundUser, errFindByEmail := uc.UserRepo.FindByEmail(ctx, email)
	if errFindByEmail != nil {
		return domain.TokenPair{}, errFindByEmail
	}

	if err := uc.hasher.Compare(foundUser.PasswordHash, password); err != nil {
		return domain.TokenPair{}, err
	}

	pairID, errGeneratePairID := uc.tokens.GeneratePairID()
	if errGeneratePairID != nil {
		return domain.TokenPair{}, errGeneratePairID
	}

	access, errGenerateAccess := uc.tokens.GenerateAccessToken(foundUser.UserID, pairID)
	if errGenerateAccess != nil {
		return domain.TokenPair{}, errGenerateAccess
	}
	refresh, refreshHash, errGenerateRefresh := uc.tokens.GenerateRefreshToken()
	if errGenerateRefresh != nil {
		return domain.TokenPair{}, errGenerateRefresh
	}

	session := domain.Session{UserID: foundUser.UserID, IPAddress: ip, RefreshHash: refreshHash, PairID: pairID, UserAgent: agent}

	if err := uc.SessionRepo.Save(ctx, session); err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}
