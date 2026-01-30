package model

import "time"

type NewsItem struct {
	Title       string
	Link        string
	PublishedAt time.Time
	Source      string
}
