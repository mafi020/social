package store

import (
	"database/sql"

	"github.com/mafi020/social/internal/interfaces"
)

type Storage struct {
	Posts    interfaces.PostsInterface
	Users    interfaces.UsersInterface
	Comments interfaces.CommentsInterface
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
