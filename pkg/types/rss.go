package types

import "time"

type RSSEntry struct {
	Title       string
	URL         string
	PublishedAt time.Time
}
