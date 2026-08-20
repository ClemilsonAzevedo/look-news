package handler

import (
	"log/slog"
	"net/http"

	"github.com/clemilsonazevedo/look-news/internal/feed"
	"github.com/go-fuego/fuego"
)

func (h *Handler) HandleNews(c fuego.ContextNoBody) (NewsRes, error) {
	criterion := c.QueryParam("criterion")
	sources := c.QueryParamArr("sources")

	if criterion == "" {
		slog.Error("no criterion provided")
		return NewsRes{}, fuego.BadRequestError{
			Title:  "Missing required param",
			Detail: "criterion is required",
			Errors: []fuego.ErrorItem{
				{
					Name:   "criterion",
					Reason: "This query parameter is required",
				},
			},
		}
	}

	if len(sources) == 0 {
		slog.Error("no sources provided")
		return NewsRes{}, fuego.BadRequestError{
			Title:  "Missing required param",
			Detail: "sources is required",
			Errors: []fuego.ErrorItem{
				{
					Name:   "sources",
					Reason: "This query parameter is required",
				},
			},
		}
	}

	var articles []feed.Article
	var hashes []string
	var failures []string

	for _, url := range sources {
		arts, err := h.refresher.GetOrRefresh(url, criterion)

		if err != nil {
			slog.Warn("source failed, falling back to cache",
				"source", url,
				"error", err,
			)
			if len(arts) == 0 {
				failures = append(failures, url)
			}
		}

		if len(arts) > 0 {
			articles = append(articles, arts...)
		}

		key := feed.CacheKey(url, criterion)
		if entry, ok := h.cache.Get(key); ok {
			hashes = append(hashes, entry.Hash)
		}
	}

	if len(failures) == len(sources) {
		slog.Error("all sources failed, no cache available",
			"sources", failures,
		)
		return NewsRes{}, fuego.HTTPError{
			Title:  "sources unavailable",
			Detail: "all sources failed, no cache available",
			Status: http.StatusServiceUnavailable,
		}
	}

	etag := `"` + feed.CombineHashes(hashes) + `"`
	w := c.Response()

	if c.Request().Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return NewsRes{}, nil
	}

	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "max-age=60, stale-while-revalidate=300")

	return NewsRes{
		Articles: articles,
		Total:    len(articles),
	}, nil
}
