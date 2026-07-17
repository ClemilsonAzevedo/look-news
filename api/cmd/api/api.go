package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clemilsonazevedo/look-news/internal/feed"
	"github.com/clemilsonazevedo/look-news/internal/feed/http/handler"
	"github.com/clemilsonazevedo/look-news/internal/feed/http/middlawares"
)

func InitServer(feedURLs []string) error {
	f, err := feed.Start(
		os.Getenv("FILTRO_PYTHON"),
		os.Getenv("FILTRO_SCRIPT"),
		os.Getenv("FILTRO_QUERY"),
	)
	if err != nil {
		return fmt.Errorf("subindo o filtro: %w", err)
	}
	defer func(f *feed.Filter) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	cache, err := feed.NewCache(feedURLs, 24*time.Hour, 1*time.Hour, f)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Carregando Feeds iniciais...")
	cache.Start(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/news", middlawares.Cors(handler.NewsHandler(cache)))
	mux.HandleFunc("/health", middlawares.Cors(handler.HealthHandler(cache)))

	srv := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		fmt.Println("Desligando...")
		cancel()
		err := srv.Shutdown(context.Background())
		if err != nil {
			return
		}
	}()

	count, newest := cache.Stats()
	fmt.Printf(
		"HTTP SERVER RUNNING 🔥| %d artigos | mais recente: %s\n",
		count, newest.Format("02 Jan 15:04"),
	)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
