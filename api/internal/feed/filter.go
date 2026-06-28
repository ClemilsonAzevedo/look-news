package feed

import (
	"sort"
	"strings"
	"time"
)

type Filter func(Article) bool

func Apply(arts []Article, filters ...Filter) []Article {
	if len(filters) == 0 {
		return arts
	}

	out := make([]Article, 0, len(arts))
	for _, a := range arts {
		ok := true
		for _, f := range filters {
			if !f(a) {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, a)
		}
	}
	return out
}

func Since(t time.Time) Filter {
	return func(a Article) bool { return a.Date.After(t) }
}

func FromSource(s string) Filter {
	s = strings.ToLower(s)
	return func(a Article) bool { return strings.ToLower(a.Source) == s }
}

func HasTerm(term string) Filter {
	term = strings.ToLower(term)
	return func(a Article) bool {
		for _, t := range a.Terms {
			if strings.ToLower(t) == term {
				return true
			}
		}
		return false
	}
}

type Ranker interface {
	Rank([]Article) []Article
}

type ByDate struct{ Desc bool }

func (r ByDate) Rank(arts []Article) []Article {
	sort.SliceStable(arts, func(i, j int) bool {
		if r.Desc {
			return arts[i].Date.After(arts[j].Date)
		}
		return arts[i].Date.Before(arts[j].Date)
	})
	return arts
}