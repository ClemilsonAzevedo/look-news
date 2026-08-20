package api

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clemilsonazevedo/look-news/internal/feed"
	"github.com/clemilsonazevedo/look-news/internal/feed/http/handler"
	"github.com/clemilsonazevedo/look-news/internal/feed/http/middlewares"
	"github.com/go-fuego/fuego"
	"github.com/joho/godotenv"
)

func InitServer() error {
	_ = godotenv.Load()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := fuego.NewServer(
		fuego.WithAddr(":"+PORT),
		fuego.WithGlobalMiddlewares(
			middlewares.Cors,
			middlewares.RateLimit,
		),
		fuego.WithTimeouts(fuego.ServerTimeouts{
			Read:  20 * time.Second,
			Write: 10 * time.Second,
			Idle:  1 * time.Minute,
		}),
	)

	cache := feed.NewCache()
	refresher := feed.NewRefresher(
		feed.NewFetcher(),
		feed.NewParser(),
		feed.NewFilter(),
		cache,
	)

	go refresher.Start(ctx, time.Hour)

	handler.NewHandler(srv, cache, refresher).BindRoutes()

	if err := srv.RunContext(ctx); err != nil {
		slog.Error("failed to run server",
			"error", err,
		)
		return err
	}

	slog.Info("server shut down successfully")
	return nil
}
