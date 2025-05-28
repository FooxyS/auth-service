package postgres

import (
	"context"

	"github.com/FooxyS/auth-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewSessionRepo(DB *pgxpool.Pool) domain.SessionRepository {

	return &SessionPostgres{DB: DB}
}

type SessionPostgres struct {
	DB *pgxpool.Pool
}

func (s *SessionPostgres) Save(ctx context.Context, session domain.Session) error {
	_, err := s.DB.Exec(ctx, "INSERT INTO session_table (user_id, ip_address, refresh_token, pair_id, useragent) VALUES ($1, $2, $3, $4, $5)", session.UserID, session.IPAddress, session.RefreshHash, session.PairID, session.UserAgent)
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionPostgres) Delete(ctx context.Context, pairID string) error {
	_, err := s.DB.Exec(ctx, "DELETE FROM session_table WHERE pair_id=$1", pairID)
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionPostgres) UpdateSession(ctx context.Context, oldPair, pair, refreshHash string) error {
	_, err := s.DB.Exec(ctx, "UPDATE session_table SET pair_id=$1, refresh_token=$2 WHERE pair_id=$3", pair, refreshHash, oldPair)
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionPostgres) FindByPairID(ctx context.Context, pairID string) (domain.Session, error) {
	var session domain.Session

	row := s.DB.QueryRow(ctx, "SELECT user_id, ip_address, refresh_token, pair_id, useragent FROM session_table WHERE pair_id=$1", pairID)
	err := row.Scan(&session.UserID, &session.IPAddress, &session.RefreshHash, &session.PairID, &session.UserAgent)
	if err != nil {
		return domain.Session{}, err
	}
	return session, nil
}
