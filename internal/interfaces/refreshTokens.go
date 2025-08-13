package interfaces

import (
	"context"
	"time"

	"github.com/mafi020/social/internal/dto"
)

type RefreshTokensInterface interface {
	Create(context.Context, int64, string, string, string, time.Time) (int64, error)
	GetByHash(context.Context, string) (*dto.RefreshToken, error)
	Revoke(context.Context, string) error
	RevokeAllForUser(context.Context, int64) error
	CleanupExpired(context.Context) (int64, error)
}
