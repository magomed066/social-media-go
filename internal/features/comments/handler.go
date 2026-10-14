package comments

import (
	"errors"
	"net/http"
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
