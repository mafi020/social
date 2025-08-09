package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mafi020/social/internal/dto"
	"github.com/mafi020/social/internal/errs"
)

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context, comment *dto.Comment) error {
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
func (s *CommentStore) GetCommentsByPostID(ctx context.Context, postID int64) ([]dto.Comment, error) {
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

	comments := []dto.Comment{}
	for rows.Next() {
		var comment dto.Comment
		comment.User = dto.CommentUser{}

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
func (s *CommentStore) GetByID(ctx context.Context, commentID int64) (*dto.Comment, error) {
	query := `
		SELECT id, post_id, user_id, content, created_at, updated_at
		FROM comments
		WHERE id=$1
	`

	comment := &dto.Comment{}

	err := s.db.QueryRowContext(ctx, query, commentID).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errs.ErrNotFound
		default:
			return nil, err
		}
	}
	return comment, nil
}
func (s *CommentStore) Update(ctx context.Context, comment *dto.Comment) error {
	query := `
		UPDATE comments
		SET content=$1
		WHERE id=$2
		RETURNING id, post_id, user_id, content, created_at, updated_at
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.Content,
		comment.ID,
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
func (s *CommentStore) Delete(ctx context.Context, commentID int64) error {
	query := `
		DELETE FROM comments
		WHERE id=$1
	`

	res, err := s.db.ExecContext(ctx, query, commentID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}
