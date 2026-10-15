package comments

import (
	"context"
	"database/sql"
	"social-media-app/internal/consts"
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

func (r *repository) Create(ctx context.Context, comment *Comment) error {
	query := `
		WITH created_comment AS (
			INSERT INTO comments (content, post_id, user_id)
			VALUES ($1, $2, $3)
			RETURNING id, created_at, user_id
		)
		SELECT c.id, c.created_at, u.id, u.username
		FROM created_comment c
		JOIN users u ON u.id = c.user_id
	`

	ctx, cancel := context.WithTimeout(ctx, consts.QueryTimeoutDuration)
	defer cancel()

	err := r.db.QueryRowContext(
		ctx,
		query,
		comment.Content,
		comment.PostID,
		comment.UserID,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.User.ID,
		&comment.User.Username,
	)
	if err != nil {
		return err
	}

	return nil
}
