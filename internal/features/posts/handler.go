package posts

import (
	"errors"
	"net/http"
	postsDTO "social-media-app/internal/dto/posts"
	mainErrors "social-media-app/internal/errors/main"
	"social-media-app/internal/json"
)

type Handler struct {
	service         service
	commentsService CommentsService
}

func NewHandler(service service, commentsService CommentsService) Handler {
	return Handler{
		service:         service,
		commentsService: commentsService,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var payload postsDTO.CreatePostPayload

	err := json.ReadJson(w, r, &payload)
	if err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
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
	post, ok := PostFromContext(r.Context())
	if !ok {
		mainErrors.InternalServerError(
			w, r, errors.New("post missing from request context"),
		)
		return
	}

	comments, err := h.commentsService.GetByPostID(r.Context(), post.ID)
	if err != nil {
		mainErrors.InternalServerError(w, r, err)
		return
	}

	post.Comments = comments

	json.WriteJSON(w, http.StatusOK, post)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	post, ok := PostFromContext(r.Context())
	if !ok {
		mainErrors.InternalServerError(
			w, r, errors.New("post missing from request context"),
		)
		return
	}

	var payload postsDTO.UpdatePostPayload
	if err := json.ReadJson(w, r, &payload); err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		mainErrors.BadRequestResponse(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}

	err := h.service.Update(r.Context(), post)
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

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	post, ok := PostFromContext(r.Context())
	if !ok {
		mainErrors.InternalServerError(
			w, r, errors.New("post missing from request context"),
		)
		return
	}

	err := h.service.Delete(r.Context(), post.ID)
	if err != nil {
		if errors.Is(err, mainErrors.ErrNotFound) {
			mainErrors.NotFoundResponse(w, r, err)
			return
		}

		mainErrors.InternalServerError(w, r, err)
		return
	}

	json.WriteJSON(w, http.StatusNoContent, nil)
}
