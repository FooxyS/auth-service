package database

import (
	"context"
	"time"

	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/FooxyS/auth-service/pkg/consts"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPostgres struct {
}

func (u *UserPostgres) Exists(ctx context.Context, email string) (bool, error) {
	pgpool := ctx.Value(consts.CTX_KEY_DB).(*pgxpool.Pool)

	var userID string

	if err := pgpool.QueryRow(ctx, "SELECT user_id FROM users WHERE email=$1", email).Scan(&userID); err != nil {
		return false, err
	}

	return true, nil
}

func (u *UserPostgres) Save(ctx context.Context, user domain.User) error {
	pgpool := ctx.Value(consts.CTX_KEY_DB).(*pgxpool.Pool)

	_, err := pgpool.Exec(ctx, "INSERT INTO users (user_id, email, password, creation_date) VALUES ($1, $2, $3, $4)", user.UserID, user.Email, user.PasswordHash, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (u *UserPostgres) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	pgpool := ctx.Value(consts.CTX_KEY_DB).(*pgxpool.Pool)

	var user domain.User

	err := pgpool.QueryRow(ctx, "SELECT * FROM users WHERE email=$1", email).Scan(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserPostgres) FindByUserID(ctx context.Context, id string) (domain.User, error) {
	pgpool := ctx.Value(consts.CTX_KEY_DB).(*pgxpool.Pool)

	var user domain.User

	err := pgpool.QueryRow(ctx, "SELECT * FROM users WHERE email=$1", id).Scan(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserPostgres) GenerateUserID() (string, error) {
	return uuid.New().String(), nil
}
