package interfaces

import (
	"context"

	"github.com/mafi020/social/internal/models"
)

type UsersInterface interface {
	Create(context.Context, *models.User) error
	GetById(context.Context, int64) (*models.User, error)
	GetByEmail(context.Context, string) (*models.User, error)
	GetByUsername(context.Context, string) (*models.User, error)
	IsUserUnique(context.Context, string, string) (map[string]string, error)
}
