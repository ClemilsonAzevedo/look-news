package handler

import (
	"reflect"

	"github.com/clemilsonazevedo/look-news/internal/feed"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

type NewsRes struct {
	Articles []feed.Article `json:"articles"`
	Total    int            `json:"total"`
}

type Handler struct {
	srv       *fuego.Server
	cache     *feed.Cache
	refresher *feed.Refresher
}

func NewHandler(srv *fuego.Server, cache *feed.Cache, refresher *feed.Refresher) *Handler {
	return &Handler{srv: srv, cache: cache, refresher: refresher}
}

func (h *Handler) BindRoutes() {
	fuego.Get(h.srv, "/news", h.HandleNews,
		option.Summary("get news on your sources"),
		option.Description("A route to send your sources and get news filtered with groq AI"),
		option.Query(
			"criterion",
			"user criterion",
			fuego.ParamRequired(),
		),
		option.QueryArray(
			"sources",
			"an sources array to filter news",
			reflect.String,
			fuego.ParamRequired(),
		),
	)
}
