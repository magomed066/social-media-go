package main

import (
	"fmt"
	"net/http"
	"social-media-app/internal/auth"
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
