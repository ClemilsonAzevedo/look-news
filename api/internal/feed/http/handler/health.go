package handler

import (
  "net/http"
  "time"

  "github.com/clemilsonazevedo/look-news/internal/feed"
  "github.com/clemilsonazevedo/look-news/internal/feed/http/helpers"
)

func HealthHandler(cache *feed.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count, newest := cache.Stats()
		helpers.WriteJSON(w, map[string]any{
			"status":   "ok",
			"articles": count,
			"newest":   newest.Format(time.RFC3339),
		})
	}
}