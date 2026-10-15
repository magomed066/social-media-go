package comments

import "social-media-app/internal/features/users"

type Comment struct {
	ID        int64        `json:"id"`
	Content   string       `json:"content"`
	UserID    int64        `json:"user_id"`
	PostID    int64        `json:"post_id"`
	CreatedAt string       `json:"created_at"`
	User      users.Author `json:"user"`
}
