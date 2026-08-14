package handler

import (
	"github.com/clemilsonazevedo/look-news/internal/feed"
	"github.com/go-fuego/fuego"
)

type Handler struct {
	srv     *fuego.Server
	fetcher *feed.Fetcher
	parser  *feed.Parser
	filter  *feed.Filter
}

func NewHandler(srv *fuego.Server, fetcher *feed.Fetcher, parser *feed.Parser, filter *feed.Filter) *Handler {
	return &Handler{
		srv:     srv,
		fetcher: fetcher,
		parser:  parser,
		filter:  filter,
	}
}

func (h *Handler) BindRoutes() {
	fuego.Post(h.srv, "/news", h.HandleNews)
}
