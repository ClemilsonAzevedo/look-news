package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/clemilsonazevedo/look-news/internal/feed"
	"github.com/go-fuego/fuego"
)

func (h *Handler) HandleNews(ctx fuego.ContextNoBody) (NewsRes, error) {
	sources := ctx.QueryParamArr("sources")
	criterion := strings.TrimSpace(ctx.QueryParam("criterion"))

	if len(sources) == 0 {
		return NewsRes{}, fuego.BadRequestError{
			Title: "no sources provided",
			Err:   fmt.Errorf("at least one source is required"),
		}
	}
	if criterion == "" {
		return NewsRes{}, fuego.BadRequestError{
			Title: "no criterion provided",
			Err:   fmt.Errorf("a filter criterion is required"),
		}
	}
	for _, s := range sources {
		if strings.TrimSpace(s) == "" {
			return NewsRes{}, fuego.BadRequestError{
				Title: "invalid source",
				Err:   fmt.Errorf("sources cannot be empty"),
			}
		}
		u, err := url.ParseRequestURI(s)
		if err != nil || u.Host == "" {
			return NewsRes{}, fuego.BadRequestError{
				Title: "invalid source",
				Err:   fmt.Errorf("invalid URL: %s", s),
			}
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			return NewsRes{}, fuego.BadRequestError{
				Title: "invalid source",
				Err:   fmt.Errorf("unsupported URL scheme: %q", u.Scheme),
			}
		}
	}

	h.refresher.EnsureSources(sources, criterion)

	var articles []feed.Article
	var hashes []string

	for _, s := range sources {
		key := feed.CacheKey(s, criterion)
		entry, ok := h.cache.Get(key)
		if !ok {
			continue
		}
		articles = append(articles, entry.Articles...)
		hashes = append(hashes, entry.Hash)
	}

	etag := `"` + feed.CombineHashes(hashes) + `"`
	w := ctx.Response()

	if ctx.Request().Header.Get("If-None-Match") == etag {
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
