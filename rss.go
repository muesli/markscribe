package main

import (
	"time"

	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
)

type RSSEntry struct {
	Title       string
	URL         string
	PublishedAt time.Time
	Extensions  ext.Extensions
}

func rssFeed(url string, count int) []RSSEntry {
	var r []RSSEntry

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		panic(err)
	}

	for _, v := range feed.Items {
		// fmt.Printf("%+v\n", v)

		r = append(r, RSSEntry{
			Title:       v.Title,
			URL:         v.Link,
			PublishedAt: *v.PublishedParsed,
			Extensions:  v.Extensions,
		})
		if len(r) == count {
			break
		}
	}

	return r
}
