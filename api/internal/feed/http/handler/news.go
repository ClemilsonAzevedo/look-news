package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/clemilsonazevedo/look-news/internal/feed"
	"github.com/go-fuego/fuego"
)

func (h *Handler) HandleNews(c fuego.ContextNoBody) (NewsRes, error) {
	criterion := c.QueryParam("criterion")
	sources := c.QueryParamArr("sources")

	if len(sources) == 0 {
		slog.Error("no sources provided")
		return NewsRes{}, fuego.BadRequestError{
			Title: "no sources provided",
			Err:   fmt.Errorf("no sources provided"),
		}
	}

	var articles []feed.Article
	var hashes []string
	var failures []string

	for _, url := range sources {
		arts, err := h.refresher.GetOrRefresh(url, criterion)

		if err != nil {
			slog.Warn("fonte falhou, seguindo com o que tinha em cache",
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
		slog.Error("todas as fontes falharam sem cache disponível",
			"fontes", failures,
		)
		return NewsRes{}, fuego.HTTPError{
			Title:  "sources unavailable",
			Detail: "não foi possível obter conteúdo de nenhuma fonte e não há cache disponível",
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
