package comments

import (
	"errors"
	"net/http"
	commentsDTO "social-media-app/internal/dto/comments"
	mainErrors "social-media-app/internal/errors/main"
	"social-media-app/internal/json"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service service
}

func NewHandler(service service) Handler {
	return Handler{service: service}
}

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
