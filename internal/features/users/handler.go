package users

import (
	"net/http"
	"social-media-app/internal/json"
	usersDTO "social-media-app/internal/shared/dto/users"
	mainErrors "social-media-app/internal/shared/errors/main"
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
