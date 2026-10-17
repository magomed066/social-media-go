package users

import (
	"errors"
	"net/http"
	"social-media-app/internal/json"
	usersDTO "social-media-app/internal/shared/dto/users"
	mainErrors "social-media-app/internal/shared/errors/main"
	usersErrors "social-media-app/internal/shared/errors/users"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	service service
}

func NewHandler(service service) Handler {
	return Handler{
		service: service,
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

		mainErrors.InternalServerError(w, r, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusCreated, user); err != nil {
		mainErrors.InternalServerError(w, r, err)
	}

	// hashing password
	// hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	mainErrors.InternalServerError(w, r, err)
	// 	return
	// }

	// user.Password = hah

	// json.WriteJSON(w, http.StatusCreated, payload)
}
