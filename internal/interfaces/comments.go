package interfaces

import (
	"context"

	"github.com/mafi020/social/internal/models"
)

type CommentsInterface interface {
	Create(context.Context, *models.Comment) error
	GetCommentsByPostID(context.Context, int64) ([]models.Comment, error)
}
