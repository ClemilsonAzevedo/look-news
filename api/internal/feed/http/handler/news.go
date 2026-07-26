package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/clemilsonazevedo/look-news/internal/feed"
	"github.com/clemilsonazevedo/look-news/internal/feed/http/helpers"
)

func NewsHandler(cache *feed.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var urls []string
			if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
				fmt.Println("[news] body inválido ou vazio:", err)
			} else if len(urls) > 0 {
				fmt.Printf("[news] recebidas %d fontes\n", len(urls))
				cache.SetURLs(urls)
			}
		}

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
