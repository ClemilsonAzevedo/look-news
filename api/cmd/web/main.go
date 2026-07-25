package main

import (
	"os"

	"github.com/clemilsonazevedo/look-news/cmd/api"
)

func main() {
	if err := api.InitServer(); nil != err {
		os.Exit(1)
	}
}
