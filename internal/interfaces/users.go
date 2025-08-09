package interfaces

import (
	"context"

	"github.com/mafi020/social/internal/dto"
)

type UsersInterface interface {
	Create(context.Context, *dto.User) error
	GetByEmail(context.Context, string) (*dto.User, error)
	GetByUsername(context.Context, string) (*dto.User, error)
	IsUserUnique(context.Context, string, string) (map[string]string, error)
	GetById(context.Context, int64) (*dto.User, error)
	Delete(context.Context, int64) error
}
