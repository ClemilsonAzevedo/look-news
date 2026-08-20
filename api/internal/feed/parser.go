package feed

import (
	"fmt"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type Parser struct {
	parser *gofeed.Parser
}

func NewParser() *Parser {
	return &Parser{
		parser: gofeed.NewParser(),
	}
}

func (p *Parser) ParseFeed(data fetchResponse) ([]Article, error) {
	f, err := p.parser.ParseString(string(data.Body))
	if err != nil {
		return nil, fmt.Errorf("parse feed: %w", err)
	}

	articles := make([]Article, 0, len(f.Items))

	for _, item := range f.Items {
		if item == nil {
			continue
		}

		date := time.Time{}
		published := item.Published

		if item.PublishedParsed != nil {
			date = *item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			date = *item.UpdatedParsed
		} else {
			published = item.Updated
		}

		author := ""
		if item.Author != nil {
			author = item.Author.Name
		}

		articles = append(articles, Article{
			Title:     strings.TrimSpace(item.Title),
			Summary:   strings.TrimSpace(firstOf(item.Description, item.Content)),
			Link:      strings.TrimSpace(firstOf(item.Link, item.GUID)),
			Date:      date,
			Source:    strings.TrimSpace(f.Title),
			Author:    strings.TrimSpace(author),
			Published: strings.TrimSpace(published),
			Terms:     trimSlice(item.Categories),
		})
	}

	return articles, nil
}

func firstOf(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}

	return ""
}

func trimSlice(values []string) []string {
	out := make([]string, 0, len(values))

	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			out = append(out, value)
		}
	}

	return out
}
