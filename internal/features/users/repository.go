package users

import (
	"context"
	"database/sql"
)

type repository struct {
	db *sql.DB
}

func (s *repository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (first_name, last_name, password, email)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Firstname,
		user.Lastname,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
