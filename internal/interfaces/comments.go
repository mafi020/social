package interfaces

import (
	"context"

	"github.com/mafi020/social/internal/models"
)

type CommentsInterface interface {
	Create(context.Context, *models.Comment) error
	GetCommentsByPostID(context.Context, int64) ([]models.Comment, error)
	GetByID(context.Context, int64) (*models.Comment, error)
	Update(context.Context, *models.Comment) error
	Delete(context.Context, int64) error
}
