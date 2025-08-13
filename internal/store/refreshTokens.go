package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/mafi020/social/internal/dto"
)

type RefreshTokensStore struct {
	db *sql.DB
}

// Create stores a new refresh token (hash) for a user.
func (s *RefreshTokensStore) Create(ctx context.Context, userID int64, tokenHash, userAgent, ip string, expiresAt time.Time) (int64, error) {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id int64
	err := s.db.QueryRowContext(ctx, query, userID, tokenHash, userAgent, (ip), expiresAt).Scan(&id)
	return id, err
}

// GetByHash returns a token row by its hash (only if not revoked).
func (s *RefreshTokensStore) GetByHash(ctx context.Context, tokenHash string) (*dto.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, user_agent, COALESCE(ip_address::text, ''), expires_at, revoked_at, created_at, updated_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`
	var rt dto.RefreshToken
	err := s.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.TokenHash,
		&rt.UserAgent,
		&rt.IPAddress,
		&rt.ExpiresAt,
		&rt.RevokedAt,
		&rt.CreatedAt,
		&rt.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

// Revoke marks a token as revoked now.
func (s *RefreshTokensStore) Revoke(ctx context.Context, tokenHash string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE token_hash = $1 AND revoked_at IS NULL
	`
	_, err := s.db.ExecContext(ctx, query, tokenHash)
	return err
}

// RevokeAllForUser revokes all active refresh tokens for a user.
func (s *RefreshTokensStore) RevokeAllForUser(ctx context.Context, userID int64) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`,
		userID,
	)
	return err
}

// CleanupExpired removes (or revokes) expired tokens.
// You can call periodically via a cron/worker; here we hard-delete.
func (s *RefreshTokensStore) CleanupExpired(ctx context.Context) (int64, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE expires_at < NOW()`)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}
