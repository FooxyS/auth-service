package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewUserRepo(DB *pgxpool.Pool) domain.UserRepository {

	return &UserPostgres{DB: DB}
}

type UserPostgres struct {
	DB *pgxpool.Pool
}

func (u *UserPostgres) Exists(ctx context.Context, email string) (bool, error) {
	var userID string

	err := u.DB.QueryRow(ctx, "SELECT user_id FROM users WHERE email=$1", email).Scan(&userID)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (u *UserPostgres) Save(ctx context.Context, user domain.User) error {
	_, err := u.DB.Exec(ctx, "INSERT INTO users (user_id, email, password, creation_date) VALUES ($1, $2, $3, $4)", user.UserID, user.Email, user.PasswordHash, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (u *UserPostgres) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User

	err := u.DB.QueryRow(ctx, "SELECT user_id, email, password FROM users WHERE email=$1", email).Scan(&user.UserID, &user.Email, &user.PasswordHash)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserPostgres) FindByUserID(ctx context.Context, id string) (domain.User, error) {
	var user domain.User

	err := u.DB.QueryRow(ctx, "SELECT * FROM users WHERE user_id=$1", id).Scan(&user.UserID, &user.Email, &user.PasswordHash)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserPostgres) GenerateUserID() (string, error) {
	return uuid.New().String(), nil
}
