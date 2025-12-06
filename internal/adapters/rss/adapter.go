package rss

import (
	"github.com/mmcdole/gofeed"
	"hufschlaeger.net/markscribe/internal/domain"
)

type Adapter struct {
}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) RecentFeedEntries(url string, count int) (
	[]domain.RSSEntry, error) {
	var parser = gofeed.NewParser()
	var r []domain.RSSEntry

	feed, err := parser.ParseURL(url)
	if err != nil {
		panic(err)
	}

	for _, v := range feed.Items {
		r = append(r, domain.RSSEntry{
			Title:       v.Title,
			URL:         v.Link,
			PublishedAt: *v.PublishedParsed,
		})
		if len(r) == count {
			break
		}
	}

	return r, nil
}
