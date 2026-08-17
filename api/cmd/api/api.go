package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clemilsonazevedo/look-news/internal/feed"
	"github.com/clemilsonazevedo/look-news/internal/feed/http/handler"
	"github.com/go-fuego/fuego"
	"github.com/joho/godotenv"
)

func InitServer() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

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
			Read:  20 * time.Second,
			Write: 10 * time.Second,
			Idle:  1 * time.Minute,
		}),
	)

	ctx := context.Background()
	cache := feed.NewCache()
	refresher := feed.NewRefresher(
		feed.NewFetcher(),
		feed.NewParser(),
		feed.NewFilter(),
		cache,
	)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	go refresher.Start(ctx, time.Hour)

	handler.NewHandler(srv, cache, refresher).BindRoutes()

	if err := srv.Run(); err != nil {
		slog.Error("cannot init the server",
			"error", err,
		)
		return err
	}

	return nil
}
