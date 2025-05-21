package usecase

import (
	"context"
	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/FooxyS/auth-service/pkg/apperrors"
)

type LoginUseCase struct {
	UserRepo    domain.UserRepository
	SessionRepo domain.SessionRepository
	tokens      domain.TokenService
}

func (uc *LoginUseCase) Execute(ctx context.Context, email, password, ip, agent string) (*domain.TokenPair, error) {
	foundUser, errFindByEmail := uc.UserRepo.FindByEmail(ctx, email)
	if errFindByEmail != nil {
		return nil, errFindByEmail
	}
	rightPass := foundUser.CheckPassword(password)
	if !rightPass {
		return nil, apperrors.ErrPasswordMismatch
	}

	pairID, errGeneratePairID := uc.tokens.GeneratePairID()
	if errGeneratePairID != nil {
		return nil, errGeneratePairID
	}

	access, errGenerateAccess := uc.tokens.GenerateAccessToken(foundUser.UserID, pairID)
	if errGenerateAccess != nil {
		return nil, errGenerateAccess
	}
	refresh, refreshHash, errGenerateRefresh := uc.tokens.GenerateRefreshToken()
	if errGenerateRefresh != nil {
		return nil, errGenerateRefresh
	}

	pair := &domain.TokenPair{access, refresh}

	session := domain.Session{UserID: foundUser.UserID, IPAddress: ip, RefreshHash: refreshHash, PairID: pairID, UserAgent: agent}

	errSave := uc.SessionRepo.Save(ctx, session)
	if errSave != nil {
		return nil, errSave
	}
	return pair, nil
}
