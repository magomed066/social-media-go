package main

import (
	"net/http"
	"social-media-app/docs"
	"social-media-app/internal/handlers"
	"time"

	"charm.land/log/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// httpSwagger "github.com/swaggo/http-swagger"
	// _ "github.com/swaggo/http-swagger/example/go-chi/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	config   config
	handlers handlers.Handlers
}

type config struct {
	addr   string
	db     dbConfig
	env    string
	apiURL string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr) // pick one ClientIPFrom* based on your infra, see below
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Dependencies

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("server has started"))
		})

		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("doc.json"),
		))

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.handlers.Posts.Create)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.handlers.Posts.PostContextMiddleware)

				r.Get("/", app.handlers.Posts.GetByID)
				r.Patch("/", app.handlers.Posts.Update)
				r.Delete("/", app.handlers.Posts.Delete)

				r.Get("/comments", app.handlers.Comments.GetByPostID)
				r.Post("/comments", app.handlers.Comments.Create)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.handlers.Users.ActivateUser)

			r.Delete("/delete/{userID}", app.handlers.Users.Delete)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.handlers.Users.UserContextMiddleware)

				r.Get("/", app.handlers.Users.GetByID)
				r.Put("/follow", app.handlers.Users.Follow)
				r.Put("/unfollow", app.handlers.Users.UnFollow)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.handlers.Posts.GetUserFeed)
			})

			r.Route("/auth", func(r chi.Router) {
				r.Post("/register", app.handlers.Users.Register)
			})
		})
	})

	return r
}

func (app *application) run(handler http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Infof("server has started at: %s", app.config.addr)

	return srv.ListenAndServe()
}
