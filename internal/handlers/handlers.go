package handlers

import (
	"database/sql"
	"social-media-app/internal/features/posts"
)

type Handlers struct {
	Posts posts.Handler
}

func NewHandlers(db *sql.DB) Handlers {
	// Posts
	postsRepository := posts.NewRepository(db)
	postsService := posts.NewService(postsRepository)
	postsHandler := posts.NewHandler(postsService)

	return Handlers{
		Posts: postsHandler,
	}
}
