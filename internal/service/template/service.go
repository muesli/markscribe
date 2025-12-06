package template

import (
	"fmt"
	"reflect"
	"time"

	"github.com/KyleBanks/goodreads/responses"
	"github.com/dustin/go-humanize"
	literalpkg "hufschlaeger.net/markscribe/internal/adapters/literal"
	domain "hufschlaeger.net/markscribe/internal/domain"
	githubsvc "hufschlaeger.net/markscribe/internal/service/github"
	goodreadssvc "hufschlaeger.net/markscribe/internal/service/goodreads"
	literalsvc "hufschlaeger.net/markscribe/internal/service/literal"
	rsssvc "hufschlaeger.net/markscribe/internal/service/rss"
)

// Service composes per-port services and exposes template-facing API.
type Service struct {
	gh  *githubsvc.Service
	gr  *goodreadssvc.Service
	lit *literalsvc.Service
	rss *rsssvc.Service
}

func New(gh *githubsvc.Service, gr *goodreadssvc.Service, lit *literalsvc.Service, rss *rsssvc.Service) *Service {
	return &Service{gh: gh, gr: gr, lit: lit, rss: rss}
}

// GitHub
func (s *Service) RecentRepos(count int) []domain.Repo { return s.gh.RecentRepos(count) }
func (s *Service) RecentForks(count int) []domain.Repo { return s.gh.RecentForks(count) }
func (s *Service) Repo(owner, name string) domain.Repo { return s.gh.Repo(owner, name) }
func (s *Service) Followers(count int) []domain.User   { return s.gh.Followers(count) }
func (s *Service) RecentPullRequests(count int) []domain.PullRequest {
	return s.gh.RecentPullRequests(count)
}
func (s *Service) RecentReleases(count int) []domain.Repo { return s.gh.RecentReleases(count) }
func (s *Service) RecentContributions(count int) []domain.Contribution {
	return s.gh.RecentContributions(count)
}
func (s *Service) Gists(count int) []domain.Gist         { return s.gh.Gists(count) }
func (s *Service) RecentStars(count int) []domain.Star   { return s.gh.RecentStars(count) }
func (s *Service) RecentIssues(count int) []domain.Issue { return s.gh.RecentIssues(count) }
func (s *Service) Sponsors(count int) []domain.Sponsor   { return s.gh.Sponsors(count) }

// GoodReads
func (s *Service) GoodReadsReviews(count int) []responses.Review { return s.gr.Reviews(count) }
func (s *Service) GoodReadsCurrentlyReading(count int) []responses.Review {
	return s.gr.CurrentlyReading(count)
}

// Literal.club
func (s *Service) LiteralCurrentlyReading(count int) []literalpkg.Book {
	return s.lit.CurrentlyReading(count)
}

// LatestRssFeeds RSS
func (s *Service) LatestRssFeeds(url string, count int) []domain.RSSEntry {
	return s.rss.LastFeedEntries(url, count)
}

// Utils (moved from root template.go to declutter main package)
func (s *Service) Humanize(t interface{}) string {
	switch v := t.(type) {
	case time.Time:
		// flatten time to prevent updating README too often
		v = time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Location())
		if time.Since(v) <= time.Hour*24 {
			return "today"
		}
		return humanize.Time(v)
	default:
		return fmt.Sprintf("%v", t)
	}
}

func (s *Service) Reverse(slc interface{}) interface{} {
	n := reflect.ValueOf(slc).Len()
	swap := reflect.Swapper(slc)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
	return slc
}
