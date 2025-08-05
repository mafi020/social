package store

import (
	"context"
	"database/sql"

	"github.com/mafi020/social/internal/models"
)

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context, comment *models.Comment) error {
	query := `
		INSERT INTO comments(post_id, user_id, content)
		VALUES($1, $2, $3) RETURNING id, post_id, user_id, content, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		&comment.PostID,
		&comment.UserID,
		&comment.Content,
	).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *CommentStore) GetCommentsByPostID(ctx context.Context, postID int64) ([]models.Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at, u.id, u.username FROM comments c
		JOIN users u
		ON u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC 
	`

	rows, err := s.db.QueryContext(ctx, query, postID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var comment models.Comment
		comment.User = models.CommentUser{}

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.User.ID,
			&comment.User.UserName,
		)

		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
