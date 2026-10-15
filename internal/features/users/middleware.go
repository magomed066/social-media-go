package users

import (
	"context"
	"errors"
	"net/http"
	mainErrors "social-media-app/internal/shared/errors/main"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userContextKey struct{}

func (h *Handler) UserContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")

		parsedUserID, err := strconv.ParseInt(userID, 10, 64)
		if err != nil || parsedUserID <= 0 {
			mainErrors.BadRequestResponse(w, r, errors.New("userID must be a positive integer"))
			return
		}

		user, err := h.service.GetByID(r.Context(), parsedUserID)
		if err != nil {
			if errors.Is(err, mainErrors.ErrNotFound) {
				mainErrors.NotFoundResponse(w, r, err)
				return
			}

			mainErrors.InternalServerError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey{}, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(userContextKey{}).(*User)
	if !ok {
		return &User{}
	}
	return user
}
