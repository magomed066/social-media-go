package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"social-media-app/internal/auth"
	"social-media-app/internal/features/posts"
	"social-media-app/internal/features/users"
	mainErrors "social-media-app/internal/shared/errors/main"
	"strings"
)

func (app *application) RequireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				mainErrors.UnauthorizedErrorResponse(w, r, fmt.Errorf("Authorization header is required"))
				return
			}

			// Handle both "Bearer token" and direct token formats
			tokenString := strings.TrimSpace(authHeader)
			if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
				tokenString = strings.TrimSpace(tokenString[7:])
			}

			err := auth.VerifyToken(tokenString)
			if err != nil {
				mainErrors.UnauthorizedErrorResponse(
					w, r, fmt.Errorf("Invalid or expired token"),
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := users.GetUserFromContext(r.Context())
		post, _ := posts.PostFromContext(r.Context())

		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := app.checkRole(r.Context(), user, requiredRole)
		if err != nil {
			if errors.Is(err, mainErrors.ErrNotFound) {
				mainErrors.NotFoundResponse(w, r, err)
				return
			}

			mainErrors.InternalServerError(w, r, err)
			return
		}

		if !allowed {
			mainErrors.ForbiddenResponse(w, r, errors.New("forbiden"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRole(ctx context.Context, user *users.User, roleName string) (bool, error) {
	role, err := app.handlers.Roles.FindByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil

}
