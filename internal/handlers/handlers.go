package handlers

import (
	"database/sql"
	"social-media-app/internal/features/comments"
	"social-media-app/internal/features/posts"
	"social-media-app/internal/features/users"
	"social-media-app/internal/mailer"
)

type Handlers struct {
	Posts    posts.Handler
	Comments comments.Handler
	Users    users.Handler
}

func NewHandlers(db *sql.DB, mailer *mailer.WnessenMailer, isSandbox bool, frontendUrl string) Handlers {

	// Comments
	commentsRepository := comments.NewRepository(db)
	commentsService := comments.NewService(commentsRepository)
	commentsHandler := comments.NewHandler(commentsService)

	// Posts
	postsRepository := posts.NewRepository(db)
	postsService := posts.NewService(postsRepository)
	postsHandler := posts.NewHandler(postsService, &commentsService)

	// Users
	usersRepository := users.NewRepository(db)
	usersService := users.NewService(usersRepository)
	usersHandler := users.NewHandler(usersService, mailer, isSandbox, frontendUrl)

	return Handlers{
		Posts:    postsHandler,
		Comments: commentsHandler,
		Users:    usersHandler,
	}
}
