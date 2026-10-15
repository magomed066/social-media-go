package handlers

import (
	"database/sql"
	"social-media-app/internal/features/comments"
	"social-media-app/internal/features/posts"
)

type Handlers struct {
	Posts    posts.Handler
	Comments comments.Handler
}

func NewHandlers(db *sql.DB) Handlers {

	// Comments
	commentsRepository := comments.NewRepository(db)
	commentsService := comments.NewService(commentsRepository)
	commentsHandler := comments.NewHandler(commentsService)

	// Posts
	postsRepository := posts.NewRepository(db)
	postsService := posts.NewService(postsRepository)
	postsHandler := posts.NewHandler(postsService, &commentsService)

	return Handlers{
		Posts:    postsHandler,
		Comments: commentsHandler,
	}
}
