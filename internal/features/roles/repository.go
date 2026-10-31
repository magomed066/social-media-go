package roles

import (
	"context"
	"database/sql"
	"errors"
	mainErrors "social-media-app/internal/shared/errors/main"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) repository {
	return repository{
		db: db,
	}
}

func (r *repository) GetByName(ctx context.Context, roleName string) (Role, error) {
	query := `
		SELECT id, name, description, level 
		FROM roles WHERE name = $1
	`

	role := Role{}

	err := r.db.QueryRowContext(ctx, query, roleName).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.Level,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Role{}, mainErrors.ErrNotFound
		}
		return Role{}, nil
	}

	return role, nil
}
