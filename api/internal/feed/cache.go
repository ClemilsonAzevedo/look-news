package feed

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	mu       sync.RWMutex
	articles []Article
	seen     map[string]struct{}
	urls     []string
	ttl      time.Duration
	interval time.Duration
}

func NewCache(urls []string, ttl, interval time.Duration) *Cache {
	return &Cache{
		urls:     urls,
		ttl:      ttl,
		interval: interval,
		seen:     make(map[string]struct{}),
	}
}

func (c *Cache) Start(ctx context.Context) {
	c.refresh()

	go func() {
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.refresh()
			}
		}
	}()
}

func (c *Cache) Articles() []Article {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Article, len(c.articles))
	copy(out, c.articles)
	return out
}

func (c *Cache) Stats() (count int, newest time.Time) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, a := range c.articles {
		if a.Date.After(newest) {
			newest = a.Date
		}
	}
	return len(c.articles), newest
}

func (c *Cache) refresh() {
	fmt.Printf("[cache %s] refresh -> %d fontes /n", time.Now().Format("15:04:05"), len(c.urls))

	results := FetchFromURLs(c.urls)

	added := 0
	c.mu.Lock()

	for _, r := range results {
		if r.Err != nil {
			continue
		}

		arts, err := ParseFeed(r.Body)
		if err != nil {
			fmt.Printf("[cache]  parse error %s: %v\n", r.URL, err)
			continue
		}

		for _, a := range arts {
			key := cache(a)
			if _, exists := c.seen[key]; !exists {
				c.seen[key] = struct{}{}
				c.articles = append(c.articles, a)
				added++
			}
		}
	}

	c.prune()
	total := len(c.articles)
	c.mu.Unlock()

	fmt.Printf("[cache] %d novos | total: %d artigos\n", added, total)
}

func (c *Cache) prune() {
	cutoff := time.Now().Add(-c.ttl)
	fresh := make([]Article, 0, len(c.articles))

	for _, a := range c.articles {
		if a.Date.IsZero() || a.Date.After(cutoff) {
			fresh = append(fresh, a)
		}
	}

	if pruned := len(c.articles) - len(fresh); pruned > 0 {
		fmt.Printf("[cache] %d artigos expirados\n", pruned)
		newSeen := make(map[string]struct{}, len(fresh))

		for _, a := range fresh {
			newSeen[cache(a)] = struct{}{}
		}
		c.articles = fresh
		c.seen = newSeen
	}
}

func cache(a Article) string {
	if a.Link != "" {
		return a.Link
	}
	return a.Source + "|" + a.Title + "|" + a.Published
}