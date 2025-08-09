package store

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/lib/pq"
	"github.com/mafi020/social/internal/errs"
)

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, userID, followerID int64) error {
	query := `
		INSERT INTO followers(user_id, follower_id)
		VALUES ($1, $2)
	`
	_, err := s.db.ExecContext(ctx, query, userID, followerID)

	if err != nil {
		log.Printf("Error %d\n", err)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // 23505 is unique_violation
			return errs.ErrDuplicateEntry
		}
		return err

	}
	return nil
}
func (s *FollowerStore) UnFollow(ctx context.Context, userIDToUnfollow, followerID int64) error {
	query := `
		DELETE FROM followers
		WHERE user_id = $1 AND follower_id = $2
	`
	_, err := s.db.ExecContext(ctx, query, userIDToUnfollow, followerID)

	if err != nil {
		return err
	}
	return nil
}
