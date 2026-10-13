package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
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
