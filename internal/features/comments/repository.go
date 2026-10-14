package comments

import (
	"context"
	"database/sql"
	"time"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) repository {
	return repository{
		db: db,
	}
}

func (r *repository) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id 
		FROM comments c
		JOIN users ON users.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}

	for rows.Next() {
		var c Comment

		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.Username, &c.User.ID)
		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
