package posts

import (
	"context"
	"database/sql"
	"errors"
	"social-media-app/internal/consts"
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

	ctx, cancel := context.WithTimeout(ctx, consts.QueryTimeoutDuration)
	defer cancel()

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
		SELECT id, user_id, title, content, created_at, tags, updated_at, version
		FROM posts
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, consts.QueryTimeoutDuration)
	defer cancel()

	var post Post

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		pq.Array(&post.Tags),
		&post.UpdatedAt,
		&post.Version,
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

	ctx, cancel := context.WithTimeout(ctx, consts.QueryTimeoutDuration)
	defer cancel()

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

func (r *repository) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET title = $1,
		    content = $2,
		    updated_at = NOW(),
		    version = version + 1
		WHERE id = $3 AND version = $4
	`

	ctx, cancel := context.WithTimeout(ctx, consts.QueryTimeoutDuration)
	defer cancel()

	res, err := r.db.ExecContext(ctx, query, post.Title, post.Content, post.ID, post.Version)
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

	post.Version++

	return nil
}
