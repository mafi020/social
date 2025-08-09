package store

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/lib/pq"
	"github.com/mafi020/social/internal/dto"
	"github.com/mafi020/social/internal/errs"
)

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *dto.Post) error {
	query := `
		INSERT INTO posts (title, content, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, title, content, user_id, tags, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
func (s *PostStore) GetByID(ctx context.Context, postID int64) (*dto.Post, error) {
	query := `
		SELECT id, title, content, user_id, tags, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	post := &dto.Post{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		postID,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errs.ErrNotFound
		default:
			return nil, err
		}
	}

	return post, nil
}
func (s *PostStore) GetAll(ctx context.Context) ([]dto.Post, error) {
	query := `
		SELECT id, title, content, user_id, tags, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []dto.Post

	for rows.Next() {
		var post dto.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.UserID,
			pq.Array(&post.Tags),
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
func (s *PostStore) Delete(ctx context.Context, postId int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	res, err := s.db.ExecContext(ctx, query, postId)
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
func (s *PostStore) Update(ctx context.Context, post *dto.Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, tags = $3
		WHERE id = $4
		RETURNING id, title, content, tags, user_id, created_at, updated_at
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		post.ID,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.UserID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}
func (s *PostStore) Feed(ctx context.Context, userID int64, params dto.FeedQueryParams) ([]dto.Feed, int, error) {
	log.Printf("Params %v", params)
	countQuery := `
		SELECT COUNT(DISTINCT p.id)
		FROM posts p
		JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE (f.user_id = $1 OR p.user_id = $1)
			AND (
				$2 = '' OR
				p.title ILIKE '%' || $2 || '%' OR
				p.content ILIKE '%' || $2 || '%'
			)
			AND (
				cardinality($3::text[]) = 0 OR
				EXISTS (
					SELECT 1 FROM unnest($3::text[]) AS tag 
					WHERE tag = ANY(p.tags)
				)
			)


	`

	var totalCount int
	var tagsParam any
	if len(params.Tags) == 0 {
		tagsParam = pq.Array([]string{})
	} else {
		tagsParam = pq.Array(params.Tags)
	}
	err := s.db.QueryRowContext(ctx, countQuery, userID, params.Search, tagsParam).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit

	query := `
        SELECT
            p.id, p.user_id, p.title, p.content, p.tags, p.created_at, p.updated_at,
            COUNT(c.id) AS comments_count,
            u.id, u.username
        FROM posts p
        LEFT JOIN comments c ON c.post_id = p.id
        LEFT JOIN users u ON u.id = p.user_id
        JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $3
			WHERE 
			(f.user_id = $3 OR p.user_id = $3)
			AND (
				$4 = '' OR
				p.title ILIKE '%' || $4 || '%' OR
				p.content ILIKE '%' || $4 || '%'
			)
			AND (
				cardinality($5::text[]) = 0 OR
				EXISTS (
					SELECT 1 FROM unnest($5::text[]) AS tag 
					WHERE tag = ANY(p.tags)
				)
			)
        GROUP BY p.id, u.id, u.username
        ORDER BY p.created_at DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := s.db.QueryContext(ctx, query, params.Limit, offset, userID, params.Search, tagsParam)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var feed []dto.Feed
	for rows.Next() {
		var f dto.Feed
		err := rows.Scan(
			&f.ID,
			&f.UserID,
			&f.Title,
			&f.Content,
			pq.Array(&f.Tags),
			&f.CreatedAt,
			&f.UpdatedAt,
			&f.CommentsCount,
			&f.User.ID,
			&f.User.UserName,
		)
		if err != nil {
			return nil, 0, err
		}
		f.Comments = []dto.Comment{}
		feed = append(feed, f)
	}

	return feed, totalCount, nil
}
