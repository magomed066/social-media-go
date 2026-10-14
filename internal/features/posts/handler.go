package posts

import (
	"net/http"
	postsDTO "social-media-app/internal/dto/posts"
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var payload postsDTO.CreatePostPayload

	err := json.ReadJson(w, r, &payload)
	if err != nil {
		json.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.service.Create(r.Context(), payload)
	if err != nil {
		json.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.WriteJSON(w, http.StatusCreated, post)
}
