package handler

import (
	"net/http"
	"strconv"

	"github.com/clemilsonazevedo/look-news/internal/feed"
	"github.com/clemilsonazevedo/look-news/internal/feed/http/helpers"
)

func NewsHandler(cache *feed.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		arts := cache.Articles()

		total := len(arts)
		w.Header().Set("X-Total-Count", strconv.Itoa(total))

		limit := 100
		if l := r.URL.Query().Get("limit"); l != "" {
			if n, err := strconv.Atoi(l); err == nil && n > 0 {
				limit = n
			}
		}

		if len(arts) > limit {
			arts = arts[:limit]
		}

		helpers.WriteJSON(w, arts)
	}
}
