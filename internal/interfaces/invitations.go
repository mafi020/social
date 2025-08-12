package interfaces

import (
	"context"
	"time"

	"github.com/mafi020/social/internal/dto"
)

type InvitationInterface interface {
	Create(context.Context, *dto.Invitation) error
	GetByID(context.Context, int64) (*dto.Invitation, error)
	GetByEmail(context.Context, string) (*dto.Invitation, error)
	GetByToken(context.Context, string) (*dto.Invitation, error)
	Update(context.Context, *dto.Invitation) error
	UpdateStatus(context.Context, int64, string) error
	UpdateEmailStatus(context.Context, int64, *time.Time) error
}
