package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/mafi020/social/internal/dto"
	"github.com/mafi020/social/internal/errs"
)

type InvitationStore struct {
	db *sql.DB
}

func (s *InvitationStore) Create(ctx context.Context, inv *dto.Invitation) error {
	query := `
		INSERT INTO invitations (inviter_id, email, token, status, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		inv.InviterID,
		inv.Email,
		inv.Token,
		inv.Status,
		inv.ExpiresAt,
	).Scan(&inv.ID, &inv.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (s *InvitationStore) GetByID(ctx context.Context, id int64) (*dto.Invitation, error) {
	query := `
		SELECT id, email, status, token, expires_at, created_at, updated_at
		FROM invitations
		WHERE id = $1
	`
	var inv dto.Invitation
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&inv.ID,
		&inv.Email,
		&inv.Status,
		&inv.Token,
		&inv.ExpiresAt,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return &inv, nil
}

func (s *InvitationStore) GetByEmail(ctx context.Context, email string) (*dto.Invitation, error) {
	query := `
		SELECT id, email, status, token, expires_at, created_at, updated_at
		FROM invitations
		WHERE email = $1
	`
	var inv dto.Invitation
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&inv.ID,
		&inv.Email,
		&inv.Status,
		&inv.Token,
		&inv.ExpiresAt,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return &inv, nil
}

func (s *InvitationStore) GetByToken(ctx context.Context, token string) (*dto.Invitation, error) {
	query := `
		SELECT id, inviter_id, email, token, status, created_at, expires_at
		FROM invitations
		WHERE token = $1
	`
	inv := &dto.Invitation{}
	err := s.db.QueryRowContext(ctx, query, token).Scan(
		&inv.ID,
		&inv.InviterID,
		&inv.Email,
		&inv.Token,
		&inv.Status,
		&inv.CreatedAt,
		&inv.ExpiresAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errs.ErrNotFound
		default:
			return nil, err
		}
	}
	return inv, nil
}

func (s *InvitationStore) UpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE invitations SET status=$1 WHERE id=$2`, status, id)
	return err
}

func (s *InvitationStore) UpdateEmailStatus(ctx context.Context, id int64, emailSentAt *time.Time) error {
	_, err := s.db.ExecContext(ctx, `UPDATE invitations SET email_sent_at=$1 WHERE id=$2`, emailSentAt, id)
	return err
}

func (s *InvitationStore) Update(ctx context.Context, inv *dto.Invitation) error {
	query := `
		UPDATE invitations
		SET status = $1, expires_at = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, email, token, status, expires_at, created_at, updated_at
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		inv.Status,
		inv.ExpiresAt,
		inv.ID,
	).Scan(
		&inv.ID,
		&inv.Email,
		&inv.Token,
		&inv.Status,
		&inv.ExpiresAt,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}
