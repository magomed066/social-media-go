package users

import (
	"net/http"
	usersDTO "social-media-app/internal/dto/users"
	mainErrors "social-media-app/internal/errors/main"
	"social-media-app/internal/json"
)

type Handler struct {
	service service
}

func NewHandler(service service) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())

	json.WriteJSON(w, http.StatusOK, user)
}

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
