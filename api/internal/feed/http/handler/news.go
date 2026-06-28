package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/clemilsonazevedo/look-news/internal/feed"
	"github.com/clemilsonazevedo/look-news/internal/feed/http/helpers"
)

func NewsHandler(cache *feed.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		arts := cache.Articles()

		var filters []feed.Filter
		if term := r.URL.Query().Get("term"); term != "" {
			filters = append(filters, feed.HasTerm(term))
		}

		if src := r.URL.Query().Get("source"); src != "" {
			filters = append(filters, feed.FromSource(src))
		}

		if since := r.URL.Query().Get("since"); since != "" {
			if d, err := helpers.ParseDuration(since); err == nil {
				filters = append(filters, feed.Since(time.Now().Add(-d)))
			}
		}

		arts = feed.Apply(arts, filters...)

		desc := r.URL.Query().Get("sort") != "asc"
		feed.ByDate{Desc: desc}.Rank(arts)

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