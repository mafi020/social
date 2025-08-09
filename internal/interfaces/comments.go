package interfaces

import (
	"context"

	"github.com/mafi020/social/internal/dto"
)

type CommentsInterface interface {
	Create(context.Context, *dto.Comment) error
	GetCommentsByPostID(context.Context, int64) ([]dto.Comment, error)
	GetByID(context.Context, int64) (*dto.Comment, error)
	Update(context.Context, *dto.Comment) error
	Delete(context.Context, int64) error
}
