package main

import (
	"os"

	"github.com/clemilsonazevedo/look-news/cmd/api"
)

var feedURLs = []string{
	"https://www.tabnews.com.br/recentes/rss",
	"https://techcrunch.com/feed/",
	"https://news.bitcoin.com/feed/",
	"https://tecnoblog.net/feed/",
	"https://canaltech.com.br/rss/",
	"https://feeds.feedburner.com/TheHackersNews",
	"https://www.securityweek.com/feed/",
	"https://www.freecodecamp.org/news/rss/",
	"https://www.freecodecamp.org/news/rss/",
	"https://newsletter.pragmaticengineer.com/feed",
	"https://blog.bytebytego.com/feed",
	"https://techcrunch.com/feed/",
	"https://spacetoday.com.br/feed/",
}

func main() {
	if err := api.InitServer(feedURLs); nil != err {
		os.Exit(1)
	}
}