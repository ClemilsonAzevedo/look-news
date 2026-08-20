package main

import (
	"log/slog"
	"os"

	"github.com/clemilsonazevedo/look-news/cmd/api"
)

func main() {
	if err := api.InitServer(); err != nil {
		slog.Error("Cannot init the server",
			"error", err,
		)
		os.Exit(1)
	}
}
