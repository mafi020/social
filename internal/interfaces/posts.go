package interfaces

import (
	"context"

	"github.com/mafi020/social/internal/dto"
)

type PostsInterface interface {
	Create(context.Context, *dto.Post) error
	GetByID(context.Context, int64) (*dto.Post, error)
	Delete(context.Context, int64) error
	Update(context.Context, *dto.Post) error
	Feed(context.Context, int64, dto.FeedQueryParams) ([]dto.Feed, int, error)
}
