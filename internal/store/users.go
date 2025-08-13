package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mafi020/social/internal/dto"
	"github.com/mafi020/social/internal/errs"
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *dto.User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3) RETURNING id, username, email, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.UserName,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}
func (s *UserStore) GetById(ctx context.Context, userId int64) (*dto.User, error) {
	query := `
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	user := &dto.User{}

	err := s.db.QueryRowContext(ctx, query, userId).Scan(&user.ID, &user.UserName, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errs.ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}
func (s *UserStore) GetByEmail(ctx context.Context, email string) (*dto.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	user := &dto.User{}

	err := s.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errs.ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}
func (s *UserStore) GetByUsername(ctx context.Context, username string) (*dto.User, error) {
	query := `
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	user := &dto.User{}

	err := s.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.UserName, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errs.ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}
func (s *UserStore) IsUserUnique(ctx context.Context, email, username string) (map[string]string, error) {
	query := `
		SELECT email, username
		FROM users
		WHERE email = $1 OR username = $2
	`

	rows, err := s.db.QueryContext(ctx, query, email, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	errors := make(map[string]string)

	for rows.Next() {
		var existingEmail, existingUsername string
		if err := rows.Scan(&existingEmail, &existingUsername); err != nil {
			return nil, err
		}
		if existingEmail == email {
			errors["email"] = "Email already exists"
		}
		if existingUsername == username {
			errors["username"] = "Username already exists"
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// If map is empty, both are unique
	if len(errors) == 0 {
		return nil, nil
	}

	return errors, nil
}
func (s *UserStore) Delete(ctx context.Context, userID int64) error {
	query := `
		DELETE FROM users
		WHERE id=$1
	`

	res, err := s.db.ExecContext(ctx, query, userID)

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
