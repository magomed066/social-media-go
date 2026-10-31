package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"social-media-app/internal/auth"
	"social-media-app/internal/json"
	"social-media-app/internal/mailer"
	usersDTO "social-media-app/internal/shared/dto/users"
	mainErrors "social-media-app/internal/shared/errors/main"
	usersErrors "social-media-app/internal/shared/errors/users"
	"strconv"
	"time"

	"charm.land/log/v2"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	service     service
	mailer      *mailer.WnessenMailer
	frontendUrl string
}

func NewHandler(service service, mailer *mailer.WnessenMailer, frontendUrl string) Handler {
	return Handler{
		service:     service,
		mailer:      mailer,
		frontendUrl: frontendUrl,
	}
}

// Get User By ID godoc
//
//	@Summary		Fetch a user
//	@Description	Fetches a user by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	json.envelope{data=User}
//	@Failure		400		{object}	json.envelope
//	@Failure		404		{object}	json.envelope
//	@Failure		500		{object}	json.envelope
//	@Failure		401		{object}	json.envelope	"Missing, invalid, or expired access token"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/ [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())

	json.WriteJSON(w, http.StatusOK, user)
}

// Follow godoc
//
//	@Summary		Follow a user
//	@Description	Adds the user identified by user_id in the request body as a follower of the user in the path.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int								true	"User ID to follow"
//	@Param			payload	body	usersDTO.FollowerUserPayload	true	"Follower's user ID"
//	@Success		204		"No Content"
//	@Failure		400		{object}	json.envelope
//	@Failure		404		{object}	json.envelope
//	@Failure		500		{object}	json.envelope
//	@Failure		401		{object}	json.envelope	"Missing, invalid, or expired access token"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (h *Handler) Follow(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())

	// Revert back to userID from ctx
	var payload usersDTO.FollowerUserPayload
	if err := json.ReadJson(w, r, &payload); err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	followUserPayload := usersDTO.FollowPayload{
		UserID:     user.ID,
		FollowerID: payload.UserID,
	}

	if err := h.service.Follow(r.Context(), followUserPayload); err != nil {
		mainErrors.InternalServerError(w, r, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusNoContent, nil); err != nil {
		mainErrors.InternalServerError(w, r, err)
	}
}

// UnFollow godoc
//
//	@Summary		Unfollow a user
//	@Description	Removes the user identified by user_id in the request body as a follower of the user in the path.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int								true	"User ID to unfollow"
//	@Param			payload	body	usersDTO.FollowerUserPayload	true	"Follower's user ID"
//	@Success		204		"No Content"
//	@Failure		400		{object}	json.envelope
//	@Failure		404		{object}	json.envelope
//	@Failure		500		{object}	json.envelope
//	@Failure		401		{object}	json.envelope	"Missing, invalid, or expired access token"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (h *Handler) UnFollow(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())

	// Revert back to userID from ctx
	var payload usersDTO.FollowerUserPayload
	if err := json.ReadJson(w, r, &payload); err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	unFollowUserPayload := usersDTO.FollowPayload{
		UserID:     user.ID,
		FollowerID: payload.UserID,
	}

	if err := h.service.UnFollow(r.Context(), unFollowUserPayload); err != nil {
		mainErrors.InternalServerError(w, r, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusNoContent, nil); err != nil {
		mainErrors.InternalServerError(w, r, err)
	}
}

// Register godoc
//
//	@Summary		Register a user
//	@Description	Creates an inactive user and an invitation for account activation.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		usersDTO.RegisterUserPayload	true	"User registration details"
//	@Success		201		{object}	json.envelope{data=User}
//	@Failure		400		{object}	json.envelope	"Invalid registration details or duplicate email or username"
//	@Failure		500		{object}	json.envelope
//	@Router			/auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var payload usersDTO.RegisterUserPayload

	if err := json.ReadJson(w, r, &payload); err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	user, err := h.service.Register(r.Context(), payload)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			mainErrors.BadRequestResponse(w, r, usersErrors.ErrPasswordTooLong)
			return
		}

		if err == usersErrors.ErrDuplicateEmail {
			mainErrors.BadRequestResponse(w, r, usersErrors.ErrDuplicateEmail)
			return
		}

		if err == usersErrors.ErrDuplicateUsername {
			mainErrors.BadRequestResponse(w, r, usersErrors.ErrDuplicateUsername)
			return
		}

		mainErrors.InternalServerError(w, r, err)
		return
	}

	// Send email
	activationURL := fmt.Sprintf("%s/confirm/%s", h.frontendUrl, user.Token)

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	err = h.mailer.Send(
		user.Email,
		mailer.UserWelcomeTemplate,
		vars,
	)
	if err != nil {
		log.Error("error sending welcome email")

		cleanupCtx, cancel := context.WithTimeout(
			context.WithoutCancel(r.Context()),
			time.Second*10,
		)
		defer cancel()

		// rollback user creation if email fails (SAGA pattern)
		if deleteErr := h.service.Delete(cleanupCtx, user.ID); deleteErr != nil {
			log.Error("error deleting user",
				"user_id", user.ID,
				"error", deleteErr,
			)
		}
		mainErrors.InternalServerError(w, r, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusCreated, user.User); err != nil {
		mainErrors.InternalServerError(w, r, err)
	}
}

// ActivateUser godoc
//
//	@Summary		Activate a user
//	@Description	Activates a user using an unexpired invitation token.
//	@Tags			Auth
//	@Produce		json
//	@Param			token	path	string	true	"Invitation token"
//	@Success		204		"No Content"
//	@Failure		404		{object}	json.envelope	"Invalid or expired invitation token"
//	@Failure		500		{object}	json.envelope
//	@Router			/auth/activate/{token} [put]
func (h *Handler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := h.service.ActivateUser(r.Context(), token)
	if err != nil {
		if errors.Is(err, mainErrors.ErrNotFound) {
			mainErrors.NotFoundResponse(w, r, err)
			return
		}

		mainErrors.InternalServerError(w, r, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusNoContent, ""); err != nil {
		mainErrors.InternalServerError(w, r, err)
	}
}

// Delete godoc
//
//	@Summary		Delete a user
//	@Description	Deletes a user by ID.
//	@Tags			Users
//	@Produce		json
//	@Param			userID	path	int	true	"User ID"
//	@Success		204		"No Content"
//	@Failure		400		{object}	json.envelope	"Invalid user ID"
//	@Failure		401		{object}	json.envelope	"Missing, invalid, or expired access token"
//	@Failure		404		{object}	json.envelope	"User not found"
//	@Failure		500		{object}	json.envelope
//	@Security		ApiKeyAuth
//	@Router			/users/delete/{userID} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	parsedUserID, err := strconv.Atoi(userID)
	if err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	if err := h.service.Delete(r.Context(), int64(parsedUserID)); err != nil {
		if errors.Is(err, mainErrors.ErrNotFound) {
			mainErrors.NotFoundResponse(w, r, err)
			return
		}

		mainErrors.InternalServerError(w, r, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusNoContent, ""); err != nil {
		mainErrors.InternalServerError(w, r, err)
		return
	}
}

// Login godoc
//
//	@Summary		Log in
//	@Description	Authenticates an activated account and returns the user and a JWT access token. Use Bearer <accessToken> in the Authorization header for protected endpoints.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		usersDTO.LoginUserPayload	true	"Login credentials"
//	@Success		200		{object}	json.envelope{data=UserWithAuthToken}
//	@Failure		400		{object}	json.envelope	"Invalid request body or missing credentials"
//	@Failure		401		{object}	json.envelope	"Invalid email or password"
//	@Failure		403		{object}	json.envelope	"Account is not activated"
//	@Failure		500		{object}	json.envelope
//	@Router			/auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var payload usersDTO.LoginUserPayload

	if err := json.ReadJson(w, r, &payload); err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	user, err := h.service.Login(r.Context(), payload)
	if err != nil {
		if errors.Is(err, mainErrors.ErrNotFound) {
			mainErrors.UnauthorizedErrorResponse(
				w, r, usersErrors.ErrInvalidCredentials,
			)
			return
		}

		mainErrors.InternalServerError(w, r, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(payload.Password),
	); err != nil {
		mainErrors.UnauthorizedErrorResponse(
			w, r, usersErrors.ErrInvalidCredentials,
		)
		return
	}

	if !user.IsActive {
		mainErrors.ForbiddenResponse(
			w, r, usersErrors.ErrAccountNotActivated,
		)
		return
	}

	token, err := auth.GenerateAuthToken(user.ID, user.Username)
	if err != nil {
		mainErrors.InternalServerError(w, r, err)
		return
	}

	data := UserWithAuthToken{
		User:        user,
		AccessToken: token,
	}

	if err := json.WriteJSON(w, http.StatusOK, data); err != nil {
		mainErrors.InternalServerError(w, r, err)
		return
	}
}
