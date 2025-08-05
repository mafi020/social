package models

type Comment struct {
	ID        int64       `json:"id"`
	PostID    int64       `json:"post_id"`
	UserID    int64       `json:"user_id"`
	Content   string      `json:"content"`
	User      CommentUser `json:"user"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

type CommentUser struct {
	ID       int64  `json:"id"`
	UserName string `json:"username"`
}
