package feed

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Cache struct {
	mu       sync.RWMutex
	articles []Article
	seen     map[string]struct{}
	urls     []string
	ttl      time.Duration
	interval time.Duration
	filter   *Filter

	refreshing int32
}

func NewCache(urls []string, ttl, interval time.Duration, filter *Filter) (*Cache, error) {
	return &Cache{
		urls:     urls,
		ttl:      ttl,
		interval: interval,
		filter:   filter,
		seen:     make(map[string]struct{}),
	}, nil
}

func (c *Cache) Start(ctx context.Context) {
	go c.refresh()

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
	if !atomic.CompareAndSwapInt32(&c.refreshing, 0, 1) {
		fmt.Println("[articleKey] refresh já em andamento, ignorando")
		return
	}
	defer atomic.StoreInt32(&c.refreshing, 0)

	fmt.Printf("[articleKey %s] refresh -> %d fontes\n", time.Now().Format("15:04:05"), len(c.urls))

	results := FetchFromURLs(c.urls)
	totalAdded := 0

	for _, r := range results {
		added, err := c.processSource(r)
		if err != nil {
			fmt.Printf("[articleKey] %s: %v\n", r.URL, err)
			continue
		}
		totalAdded += added
	}

	c.mu.Lock()
	c.prune()
	total := len(c.articles)
	c.mu.Unlock()

	fmt.Printf("[articleKey] %d novos | total: %d artigos\n", totalAdded, total)
}

func (c *Cache) processSource(r Result) (int, error) {
	if r.Err != nil {
		return 0, fmt.Errorf("fetch: %w", r.Err)
	}

	arts, err := ParseFeed(r.Body)
	if err != nil {
		return 0, fmt.Errorf("parse: %w", err)
	}

	arts, err = c.filter.ApplyFilter(arts)
	if err != nil {
		return 0, fmt.Errorf("filtro: %w", err)
	}

	added := 0
	c.mu.Lock()

	for _, a := range arts {
		key := articleKey(a)
		if _, exists := c.seen[key]; !exists {
			c.seen[key] = struct{}{}
			c.articles = append(c.articles, a)
			added++
		}
	}
	c.mu.Unlock()

	fmt.Printf("[articleKey] fonte %s -> %d novos artigos relevantes\n", r.URL, added)
	return added, nil
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
		fmt.Printf("[articleKey] %d artigos expirados\n", pruned)
		newSeen := make(map[string]struct{}, len(fresh))

		for _, a := range fresh {
			newSeen[articleKey(a)] = struct{}{}
		}
		c.articles = fresh
		c.seen = newSeen
	}
}

func articleKey(a Article) string {
	if a.Link != "" {
		return a.Link
	}
	return a.Source + "|" + a.Title + "|" + a.Published
}
