package interfaces

import (
	"context"
)

type FollowersInterface interface {
	Follow(ctx context.Context, userID, followerID int64) error
	UnFollow(ctx context.Context, userIDToUnfollow, followerID int64) error
}
