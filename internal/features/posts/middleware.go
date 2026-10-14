package posts

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	mainErrors "social-media-app/internal/errors/main"

	"github.com/go-chi/chi/v5"
)

type postContextKey struct{}

func (h *Handler) PostContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "postID")

		id, err := strconv.ParseInt(postID, 10, 64)
		if err != nil {
			mainErrors.BadRequestResponse(w, r, err)
			return
		}

		post, err := h.service.GetByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, mainErrors.ErrNotFound) {
				mainErrors.NotFoundResponse(w, r, err)
				return
			}

			mainErrors.InternalServerError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), postContextKey{}, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func PostFromContext(ctx context.Context) (*Post, bool) {
	post, ok := ctx.Value(postContextKey{}).(*Post)
	return post, ok && post != nil
}
