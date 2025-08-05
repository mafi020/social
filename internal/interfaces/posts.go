package interfaces

import (
	"context"

	"github.com/mafi020/social/internal/models"
)

type PostsInterface interface {
	Create(context.Context, *models.Post) error
	GetByID(context.Context, int64) (*models.Post, error)
	Delete(context.Context, int64) error
	Update(context.Context, *models.Post) error
}
