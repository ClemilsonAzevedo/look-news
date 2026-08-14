package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/clemilsonazevedo/look-news/internal/feed"
	"github.com/go-fuego/fuego"
)

type NewsReq struct {
	URLs []string `json:"urls" validate:"url,required,min=1"`
}

type NewsRes struct {
	Articles []feed.Article `json:"articles"`
	Total    int            `json:"total" header:"X-Total-Count"`
}

func (h *Handler) HandleNews(ctx fuego.ContextWithBody[NewsReq]) (NewsRes, error) {
	var req NewsReq
	var articles []feed.Article
	var totalArticles int

	if err := json.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
		slog.Error("invalid or empty body", "error", err)
		return NewsRes{}, fuego.BadRequestError{
			Title: "cannot decode json",
			Err:   err,
		}
	}

	if len(req.URLs) == 0 {
		slog.Error("request body has no urls", "error", "no urls provided")
		return NewsRes{}, fuego.BadRequestError{
			Title: "no urls provided",
			Err:   fmt.Errorf("no urls provided"),
		}
	}

	fetchedContent := h.fetcher.FetchFromURLs(req.URLs)
	if len(fetchedContent) == 0 {
		slog.Error("no content fetched", "error", "no content fetched")
		return NewsRes{}, fmt.Errorf("no content fetched")
	}

	for _, r := range fetchedContent {
		if r.Err != nil {
			return NewsRes{}, fmt.Errorf("fetch: %w", r.Err)
		}

		arts, err := h.parser.ParseFeed(r)
		if err != nil {
			return NewsRes{}, fmt.Errorf("parse: %w", err)
		}

		// todo: Filtrar os artigos usando a GROQ AI

		articles = append(articles, arts...)
		totalArticles = len(articles)
	}

	return NewsRes{
		Articles: articles,
		Total:    totalArticles,
	}, nil
}
