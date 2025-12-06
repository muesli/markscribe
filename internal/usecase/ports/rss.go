package ports

import (
	domain "hufschlaeger.net/markscribe/internal/domain"
)

// RssFeedPort defines operations we use from literal.club integration.
type RssFeedPort interface {
	RecentFeedEntries(url string, count int) ([]domain.RSSEntry, error)
}
