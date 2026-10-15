package users

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
	return Handler{
		service: service,
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	parsedUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil || parsedUserID <= 0 {
		mainErrors.BadRequestResponse(w, r, errors.New("userID must be a positive integer"))
		return
	}

	user, err := h.service.GetByID(r.Context(), parsedUserID)
	if err != nil {
		if errors.Is(err, mainErrors.ErrNotFound) {
			mainErrors.NotFoundResponse(w, r, err)
			return
		}

		mainErrors.InternalServerError(w, r, err)
		return
	}

	json.WriteJSON(w, http.StatusOK, user)
}
