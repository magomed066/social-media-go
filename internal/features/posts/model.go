package posts

import "social-media-app/internal/features/comments"

type Post struct {
	ID        int64              `json:"id"`
	Content   string             `json:"content"`
	Title     string             `json:"title"`
	UserID    int64              `json:"user_id"`
	Tags      []string           `json:"tags"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
	Comments  []comments.Comment `json:"comments"`
}
