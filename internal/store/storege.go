package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Posts interface {
		Create(context context.Context, post *Post) error
	}
	Users interface {
		Create(context context.Context, user *User) error
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostsStore{db},
		Users: &UsersStore{db},
	}
}
