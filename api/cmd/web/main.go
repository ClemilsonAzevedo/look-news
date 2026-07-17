package main

import (
	"os"

	"github.com/clemilsonazevedo/look-news/cmd/api"
)

var feedURLs = []string{
	"https://www.tabnews.com.br/recentes/rss",
	"https://tecnoblog.net/feed/",
	"https://canaltech.com.br/rss/",
	"https://techcrunch.com/feed/",
	"https://news.bitcoin.com/feed/",
	"https://www.securityweek.com/feed/",
	"https://news.google.com/rss?hl=en-US&gl=US&ceid=US:en",
	"https://feeds.feedburner.com/TheHackersNews",
	"https://www.freecodecamp.org/news/rss/",
	"https://newsletter.pragmaticengineer.com/feed",
	"https://blog.bytebytego.com/feed",
	"https://medium.com/feed/netflix-techblog",
}

func main() {
	if err := api.InitServer(feedURLs); nil != err {
		os.Exit(1)
	}
}
