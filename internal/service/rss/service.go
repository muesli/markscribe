package rss

import (
	"hufschlaeger.net/markscribe/internal/domain"
	"hufschlaeger.net/markscribe/internal/usecase/ports"
)

// Service wraps the RssFeedPort and contains app-level logic for rss feed features.
type Service struct {
	rss ports.RssFeedPort
}

func New(rss ports.RssFeedPort) *Service { return &Service{rss: rss} }

func (s *Service) LastFeedEntries(url string, count int) []domain.RSSEntry {
	entries, err := s.rss.RecentFeedEntries(url, count)
	if err != nil {
		panic(err)
	}
	return entries
}
