package domain

import "time"

// Contribution represents a contribution to a repo.
type Contribution struct {
	OccurredAt time.Time
	Repo       Repo
}

type Issue struct {
	Repo       Repo
	OccurredAt time.Time
	Title      string
}

// Gist represents a gist.
type Gist struct {
	Name        string
	Description string
	URL         string
	CreatedAt   time.Time
}

// Star represents a star/favorite event.
type Star struct {
	StarredAt time.Time
	Repo      Repo
}

// PullRequest represents a pull request.
type PullRequest struct {
	Title     string
	URL       string
	State     string
	CreatedAt time.Time
	Repo      Repo
}

// Release represents a release.
type Release struct {
	Name        string
	TagName     string
	PublishedAt time.Time
	URL         string
}

// Repo represents a git repo.
type Repo struct {
	Name        string
	URL         string
	Description string
	IsPrivate   bool
	Stargazers  int
	LastRelease Release
}

// Sponsor represents a sponsor.
type Sponsor struct {
	User      User
	CreatedAt time.Time
}

// User represents a SCM user.
type User struct {
	Login     string
	Name      string
	AvatarURL string
	URL       string
}

// RSSEntry represents a single RSS entry.
type RSSEntry struct {
	Title       string
	URL         string
	PublishedAt time.Time
}
