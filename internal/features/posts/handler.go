package posts

import (
	"errors"
	"net/http"
	postsDTO "social-media-app/internal/dto/posts"
	mainErrors "social-media-app/internal/errors/main"
	"social-media-app/internal/json"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service service
}

func NewHandler(service service) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var payload postsDTO.CreatePostPayload

	err := json.ReadJson(w, r, &payload)
	if err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	post, err := h.service.Create(r.Context(), payload)
	if err != nil {
		mainErrors.InternalServerError(w, r, err)
		return
	}

	json.WriteJSON(w, http.StatusCreated, post)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")

	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		mainErrors.InternalServerError(w, r, err)
		return
	}

	post, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, mainErrors.ErrNotFound) {
			mainErrors.NotFoundResponse(w, r, err)
			return
		}

		mainErrors.InternalServerError(w, r, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, post)
}
