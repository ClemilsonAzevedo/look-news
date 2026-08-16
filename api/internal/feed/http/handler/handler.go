package handler

import (
	"github.com/go-fuego/fuego"
)

type Handler struct {
	srv *fuego.Server
}

func NewHandler(srv *fuego.Server) *Handler {
	return &Handler{
		srv: srv,
	}
}

func (h *Handler) BindRoutes() {
	fuego.Post(h.srv, "/news", h.HandleNews)
}
