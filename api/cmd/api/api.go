package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/clemilsonazevedo/look-news/internal/feed"
	"github.com/clemilsonazevedo/look-news/internal/feed/http/handler"
	"github.com/go-fuego/fuego"
)

func InitServer() error {
	srv := fuego.NewServer(
		fuego.WithAddr("localhost:8080"),
		fuego.WithGlobalMiddlewares(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

				w.Header().Set("Content-Type", "application/json")
				next.ServeHTTP(w, r)
			})
		}),
		fuego.WithTimeouts(fuego.ServerTimeouts{
			Read:  10 * time.Second,
			Write: 10 * time.Second,
			Idle:  1 * time.Minute,
		}),
	)

	fetcher := feed.NewFetcher()
	parser := feed.NewParser()
	filter := feed.NewFilter()
	handler.NewHandler(srv, fetcher, parser, filter).BindRoutes()

	if err := srv.Run(); err != nil {
		slog.Error("cannot init the server",
			"error", err,
		)
		return err
	}

	return nil
}
