package feed

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type sourceCriterion struct {
	source    string
	criterion string
}

type Refresher struct {
	mu    sync.RWMutex
	known map[string]sourceCriterion

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
		slog.Info("refresher: novas combinações registradas", "quantidade", len(toRefresh))
	}
	for _, sc := range toRefresh {
		go r.refreshSource(sc.source, sc.criterion)
	}
}

func (r *Refresher) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("refresher: encerrando")
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

	slog.Info("refresher: iniciando ciclo", "combinações", len(combos))
	for _, sc := range combos {
		go r.refreshSource(sc.source, sc.criterion)
	}
}

func (r *Refresher) refreshSource(source, criterion string) {
	key := CacheKey(source, criterion)

	results := r.fetcher.FetchFromURLs([]string{source})
	if len(results) == 0 {
		slog.Error("refresh: nenhum conteúdo retornado", "source", source)
		return
	}

	res := results[0]
	if res.Err != nil {
		slog.Error("refresh: falha no fetch, tenta de novo no próximo ciclo",
			"source", source, "error", res.Err)
		return
	}

	arts, err := r.parser.ParseFeed(res)
	if err != nil {
		slog.Error("refresh: falha no parse, tenta de novo no próximo ciclo",
			"source", source, "error", err)
		return
	}

	hash := HashArticles(arts)

	if entry, ok := r.cache.Get(key); ok && entry.Hash == hash {
		slog.Info("refresh: sem mudanças, mantém cache", "source", source)
		return
	}

	filtered, err := r.filter.ApplyFilter(criterion, arts)
	if err != nil {
		slog.Error("refresh: falha no filtro, tenta de novo no próximo ciclo",
			"source", source, "error", err)
		return
	}

	r.cache.Set(key, hash, filtered)
	slog.Info("refresh: fonte atualizada", "source", source, "artigos", len(filtered))
}
