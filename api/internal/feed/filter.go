package feed

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/clemilsonazevedo/look-news/pkg/groq"
)

type Filter struct{}

func NewFilter() *Filter {
	return &Filter{}
}

type filterResponse struct {
	Relevant []string `json:"relevant"`
}

func (f *Filter) ApplyFilter(userCriterion string, articles []Article) ([]Article, error) {
	client := groq.NewClient()

	today := time.Now()
	userPrompt := buildUserPrompt(userCriterion, today, articles)

	resp, err := client.ChatCompletion(
		[]groq.Message{
			{Role: "system", Content: groq.SystemPrompt},
			{Role: "user", Content: userPrompt},
		},
	)
	if err != nil {
		slog.Error("error generating response", "err", err)
		return nil, fmt.Errorf("error filtering articles: %w", err)
	}

	llmResponse := ""
	for _, c := range resp.Choices {
		fmt.Println(c.Message.Content)
		llmResponse = c.Message.Content
	}

	var result filterResponse
	if err := json.Unmarshal([]byte(llmResponse), &result); err != nil {
		slog.Error("error parsing LLM response", "err", err, "raw", llmResponse)
		return nil, fmt.Errorf("error parsing filtered articles: %w", err)
	}

	relevantLinks := make(map[string]struct{}, len(result.Relevant))
	for _, link := range result.Relevant {
		relevantLinks[link] = struct{}{}
	}

	filtered := make([]Article, 0, len(result.Relevant))
	for _, a := range articles {
		if _, ok := relevantLinks[a.Link]; ok {
			filtered = append(filtered, a)
		}
	}

	return filtered, nil
}

func buildUserPrompt(criterion string, today time.Time, articles []Article) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Today's date: %s\n", today.Format("2006-01-02"))
	fmt.Fprintf(&b, "User's interests: %s\n\n", criterion)
	b.WriteString("Articles:\n")
	for i, a := range articles {
		fmt.Fprintf(&b, "%d. title: %s\n", i+1, a.Title)
		fmt.Fprintf(&b, "   link: %s\n", a.Link)
	}
	return b.String()
}

func startOfWeek(ref time.Time) time.Time {
	startOfDay := time.Date(ref.Year(), ref.Month(), ref.Day(), 0, 0, 0, 0, ref.Location())
	return startOfDay.AddDate(0, 0, -int(ref.Weekday()))
}

func FilterCurrentWeek(articles []Article, ref time.Time) []Article {
	weekStart := startOfWeek(ref)

	filtered := make([]Article, 0, len(articles))
	for _, a := range articles {
		if a.Date.IsZero() {
			continue
		}
		if a.Date.Before(weekStart) || a.Date.After(ref) {
			continue
		}
		filtered = append(filtered, a)
	}

	return filtered
}
