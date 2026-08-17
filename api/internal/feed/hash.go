package feed

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

func HashArticles(articles []Article) string {
	h := sha256.New()
	for _, a := range articles {
		h.Write([]byte(a.Link))
		h.Write([]byte(a.Published))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func CombineHashes(hashes []string) string {
	sorted := make([]string, len(hashes))
	copy(sorted, hashes)
	sort.Strings(sorted)

	h := sha256.New()
	for _, hh := range sorted {
		h.Write([]byte(hh))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func CacheKey(source, criterion string) string {
	h := sha256.New()
	h.Write([]byte(source))
	h.Write([]byte{0})
	h.Write([]byte(criterion))
	return hex.EncodeToString(h.Sum(nil))
}
