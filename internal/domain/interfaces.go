package domain

import "context"

type UserRepository interface {
	Exists(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
}

type SessionRepository interface {
	Save(ctx context.Context, session Session) error
	Delete(ctx context.Context, id string) error
	UpdateToken(ctx context.Context, pair, refreshHash string) error
}

type TokenService interface {
	GenerateAccessToken(id string, pairID string) (string, error) //return access, err
	GenerateRefreshToken() (string, string, error)                //refresh, refreshHash, err
	GeneratePairID() (string, error)
	CheckPair(pair TokenPair) error
	ValidateAccessToken(access string) (string, error)
}
