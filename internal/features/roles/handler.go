package roles

import (
	"context"
	"net/http"
)

type Handler struct {
	service service
}

func NewHandler(service service) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) GetByName(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) FindByName(ctx context.Context, name string) (Role, error) {
	return h.service.GetByName(ctx, name)
}
