package rss

import (
	"github.com/muesli/markscribe/pkg/types"

	"github.com/mmcdole/gofeed"
)

// Feed retrieves 'n' RSS feed entries from 'url'.
func Feed(url string, n int) []types.RSSEntry {
	entries := []types.RSSEntry{}

	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)
	if err != nil {
		panic(err)
	}

	for _, v := range feed.Items {
		entries = append(entries, types.RSSEntry{
			Title:       v.Title,
			URL:         v.Link,
			PublishedAt: *v.PublishedParsed,
		})

		if len(entries) == n {
			break
		}
	}

	return entries
}
