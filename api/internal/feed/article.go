package feed

import (
	"time"
)

type Article struct {
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Link      string    `json:"link"`
	Date      time.Time `json:"date"`
	Source    string    `json:"source"`
	Author    string    `json:"author"`
	Published string    `json:"published"`
	Terms     []string  `json:"terms"`
}