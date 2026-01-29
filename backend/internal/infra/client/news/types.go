package news

import "time"

type Item struct {
	Title       string
	Link        string
	PublishedAt time.Time
	Source      string
}
