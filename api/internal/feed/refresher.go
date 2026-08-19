package feed

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type sourceCriterion struct {
	source    string
	criterion string
}

type Refresher struct {
	mu      sync.RWMutex
	known   map[string]sourceCriterion
	fetcher *Fetcher
	parser  *Parser
	filter  *Filter
	cache   *Cache
}

func NewRefresher(fetcher *Fetcher, parser *Parser, filter *Filter, cache *Cache) *Refresher {
	return &Refresher{
		known:   make(map[string]sourceCriterion),
		fetcher: fetcher,
		parser:  parser,
		filter:  filter,
		cache:   cache,
	}
}

func (r *Refresher) EnsureSources(sources []string, criterion string) {
	r.mu.Lock()
	var toRefresh []sourceCriterion
	for _, s := range sources {
		key := CacheKey(s, criterion)
		if _, known := r.known[key]; !known {
			sc := sourceCriterion{source: s, criterion: criterion}
			r.known[key] = sc
			toRefresh = append(toRefresh, sc)
		}
	}
	r.mu.Unlock()

	if len(toRefresh) > 0 {
		slog.Info("refresher: new combinations registered", "quantity", len(toRefresh))
	}
	for _, sc := range toRefresh {
		go r.refreshSource(sc.source, sc.criterion)
	}
}

func (r *Refresher) GetOrRefresh(source, criterion string) ([]Article, error) {
	key := CacheKey(source, criterion)

	r.mu.Lock()
	if _, known := r.known[key]; !known {
		r.known[key] = sourceCriterion{
			source:    source,
			criterion: criterion,
		}
	}
	r.mu.Unlock()

	if entry, ok := r.cache.Get(key); ok {
		return entry.Articles, nil
	}

	return r.refreshSource(source, criterion)
}

func (r *Refresher) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			slog.Info("refresher: stopping")
			return
		case <-ticker.C:
			r.refreshAll()
		}
	}
}

func (r *Refresher) refreshAll() {
	r.mu.RLock()
	combos := make([]sourceCriterion, 0, len(r.known))
	for _, sc := range r.known {
		combos = append(combos, sc)
	}
	r.mu.RUnlock()

	slog.Info("refresher: refreshing all", "combinations", len(combos))
	for _, sc := range combos {
		go r.refreshSource(sc.source, sc.criterion)
	}
}

func (r *Refresher) refreshSource(source, criterion string) ([]Article, error) {
	key := CacheKey(source, criterion)

	results := r.fetcher.FetchFromURLs([]string{source})
	if len(results) == 0 {
		err := fmt.Errorf("no content returned")
		slog.Error("refresh: no content returned",
			"source", source,
		)
		return r.fallbackToCache(key), err
	}

	res := results[0]
	if res.Err != nil {
		slog.Error("refresh: fetch failed, will retry in next cycle",
			"source", source,
			"error", res.Err,
		)
		return r.fallbackToCache(key), res.Err
	}

	arts, err := r.parser.ParseFeed(res)
	if err != nil {
		slog.Error("refresh: parse failed, will retry in next cycle",
			"source", source,
			"error", err,
		)
		return r.fallbackToCache(key), err
	}

	if len(arts) == 0 {
		slog.Info("refresh: source has no items at the moment, keeping cache",
			"source", source,
		)
		return r.fallbackToCache(key), nil
	}

	hash := HashArticles(arts)
	if entry, ok := r.cache.Get(key); ok && entry.Hash == hash {
		slog.Info("refresh: no changes, keeping cache",
			"source", source,
		)
		return entry.Articles, nil
	}

	filtered, err := r.filter.ApplyFilter(criterion, arts)
	if err != nil {
		slog.Error("refresh: filter failed, will retry in next cycle",
			"source", source,
			"error", err,
		)
		return r.fallbackToCache(key), err
	}

	r.cache.Set(key, hash, filtered)
	slog.Info("refresh: source updated",
		"source", source,
		"articles", len(filtered),
	)
	return filtered, nil
}

func (r *Refresher) fallbackToCache(key string) []Article {
	if entry, ok := r.cache.Get(key); ok {
		return entry.Articles
	}
	return nil
}
