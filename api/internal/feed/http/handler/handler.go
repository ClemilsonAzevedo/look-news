package handler

import (
	"net/http"
	"reflect"

	"github.com/clemilsonazevedo/look-news/internal/feed"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
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
		option.Summary("Get news from RSS sources"),
		option.Description("Accepts a search criterion and a list of RSS sources. Uses Groq AI to filter and return the most relevant news."),
		// option.Tags("news"),
		option.Header("Content-Type", "application/json"),

		// query params
		option.Query(
			"criterion",
			"Search criterion / prompt for Groq AI",
			param.Required(),
			param.Example("example", "latest news about artificial intelligence"),
		),
		option.QueryArray(
			"sources",
			"Sources to filter by (e.g. domains)",
			reflect.String,
			param.Required(),
			param.Example("sources", []any{"https://techcrunch.com/feed"}),
		),

		// errors
		option.AddResponse(http.StatusOK, "List of filtered news", fuego.Response{
			Type: NewsRes{},
		}),

		option.AddResponse(http.StatusBadRequest, "Missing or invalid parameters", fuego.Response{
			Type: []fuego.HTTPError{{
				Title:  "Param required",
				Detail: "This query parameter is required",
			}},
		}),

		option.AddResponse(http.StatusInternalServerError, "Internal server error", fuego.Response{
			Type: fuego.HTTPError{
				Title:  "Internal server error",
				Detail: "Something went wrong",
			},
		}),

		option.AddResponse(http.StatusServiceUnavailable, "Service unavailable", fuego.Response{
			Type: fuego.HTTPError{
				Title:  "Service unavailable",
				Detail: "Something went wrong",
			},
		}),

		option.AddResponse(http.StatusTooManyRequests, "Rate limit exceeded", fuego.Response{
			Type: fuego.HTTPError{
				Title:  "Rate limit exceeded",
				Detail: "You have exceeded the rate limit",
			},
		}),
	)

}
