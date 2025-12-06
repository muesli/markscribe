package ports

import (
	"context"

	domain "hufschlaeger.net/markscribe/internal/domain"
)

// GithubPort defines the minimal set of GitHub operations used by the application.
// This is intentionally small for the first incremental extraction.
type GithubPort interface {
	RecentRepos(ctx context.Context, username string, count int, isFork bool) ([]domain.Repo, error)
	Repo(ctx context.Context, owner, name string) (domain.Repo, error)
	ViewerLogin(ctx context.Context) (string, error)
	Followers(ctx context.Context, username string, count int) ([]domain.User, error)
	RecentPullRequests(ctx context.Context, username string, count int) ([]domain.PullRequest, error)
	RecentReleases(ctx context.Context, username string, count int) ([]domain.Repo, error)
	RecentContributions(ctx context.Context, username string, count int) ([]domain.Contribution, error)
	Gists(ctx context.Context, username string, count int) ([]domain.Gist, error)
	RecentStars(ctx context.Context, username string, count int) ([]domain.Star, error)
	RecentIssues(ctx context.Context, username string, count int) ([]domain.Issue, error)
	Sponsors(ctx context.Context, username string, count int) ([]domain.Sponsor, error)
}
