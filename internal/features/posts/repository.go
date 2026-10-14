package posts

import (
	"context"
	"database/sql"
	"errors"
	postsDTO "social-media-app/internal/dto/posts"
	mainErrors "social-media-app/internal/errors/main"

	"github.com/lib/pq"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) repository {
	return repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, user_id, title, content, created_at, tags, updated_at
		FROM posts
		WHERE id = $1
	`

	var post Post

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		pq.Array(&post.Tags),
		&post.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, mainErrors.ErrNotFound
		}

		return nil, err
	}

	return &post, nil
}

func (r *repository) Delete(ctx context.Context, postID int64) error {
	query := "DELETE FROM posts WHERE id = $1"

	res, err := r.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return mainErrors.ErrNotFound
	}

	return nil
}

func (r *repository) Update(ctx context.Context, postID int64, data postsDTO.UpdatePostPayload) (Post, error) {
	query := `
		UPDATE posts
		SET title = COALESCE($1, title),
		    content = COALESCE($2, content),
		    updated_at = NOW()
		WHERE id = $3
		RETURNING id, user_id, title, content, created_at, tags, updated_at
	`

	var post Post

	err := r.db.QueryRowContext(ctx, query, data.Title, data.Content, postID).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		pq.Array(&post.Tags),
		&post.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Post{}, mainErrors.ErrNotFound
		}

		return Post{}, err
	}

	return post, nil
}

func postsContextMiddleware() {

}
