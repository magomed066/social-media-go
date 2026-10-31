package comments

import (
	"errors"
	"net/http"
	"social-media-app/internal/json"
	commentsDTO "social-media-app/internal/shared/dto/comments"
	mainErrors "social-media-app/internal/shared/errors/main"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service service
}

func NewHandler(service service) Handler {
	return Handler{service: service}
}

// GetByPostID godoc
//
//	@Summary		List comments for a post
//	@Description	Returns the comments for the specified post, including author details.
//	@Tags			Comments
//	@Produce		json
//	@Param			postID	path		int	true	"Post ID"	minimum(1)
//	@Success		200		{object}	json.envelope{data=[]Comment}
//	@Failure		400		{object}	json.envelope	"Post ID must be a positive integer"
//	@Failure		401		{object}	json.envelope	"Missing, invalid, or expired access token"
//	@Failure		404		{object}	json.envelope	"Post not found"
//	@Failure		500		{object}	json.envelope	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID}/comments [get]
func (h *Handler) GetByPostID(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)
	if err != nil || postID <= 0 {
		mainErrors.BadRequestResponse(w, r, errors.New("postID must be a positive integer"))
		return
	}

	comments, err := h.service.GetByPostID(r.Context(), postID)
	if err != nil {
		mainErrors.InternalServerError(w, r, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, comments)
}

// Create godoc
//
//	@Summary		Create a comment
//	@Description	Creates a comment on the specified post. The current implementation assigns user ID 55 as the author instead of the authenticated user.
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int									true	"Post ID"	minimum(1)
//	@Param			payload	body		commentsDTO.CreateCommentPayload	true	"Comment details"
//	@Success		201		{object}	json.envelope{data=Comment}
//	@Failure		400		{object}	json.envelope	"Invalid post ID, request body, or comment content"
//	@Failure		401		{object}	json.envelope	"Missing, invalid, or expired access token"
//	@Failure		404		{object}	json.envelope	"Post not found"
//	@Failure		500		{object}	json.envelope	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID}/comments [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")

	parsedPostID, err := strconv.ParseInt(postID, 10, 64)
	if err != nil || parsedPostID <= 0 {
		mainErrors.BadRequestResponse(w, r, errors.New("postID must be a positive integer"))
		return
	}

	var payload commentsDTO.CreateCommentPayload

	if err := json.ReadJson(w, r, &payload); err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	// Temporary until authentication provides the current user.
	userID := int64(55)

	comment, err := h.service.Create(r.Context(), parsedPostID, userID, payload)
	if err != nil {
		mainErrors.InternalServerError(w, r, err)
		return
	}

	json.WriteJSON(w, http.StatusCreated, comment)
}
