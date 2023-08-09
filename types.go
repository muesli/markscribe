package main

import (
	"time"

	"github.com/shurcooL/githubv4"
)

// Contribution represents a contribution to a repo.
type Contribution struct {
	OccurredAt time.Time
	Repo       Repo
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
	Name         string
	TagName      string
	PublishedAt  time.Time
	CreatedAt    time.Time
	URL          string
	IsLatest     bool
	IsPreRelease bool
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

type qlGist struct {
	Name        githubv4.String
	Description githubv4.String
	URL         githubv4.String
	CreatedAt   githubv4.DateTime
}

type qlPullRequest struct {
	URL        githubv4.String
	Title      githubv4.String
	State      githubv4.PullRequestState
	CreatedAt  githubv4.DateTime
	Repository qlRepository
}

type qlRelease struct {
	Nodes []struct {
		Name         githubv4.String
		TagName      githubv4.String
		PublishedAt  githubv4.DateTime
		CreatedAt    githubv4.DateTime
		URL          githubv4.String
		IsLatest     githubv4.Boolean
		IsPrerelease githubv4.Boolean
		IsDraft      githubv4.Boolean
	}
}

type qlRepository struct {
	NameWithOwner githubv4.String
	URL           githubv4.String
	Description   githubv4.String
	IsPrivate     githubv4.Boolean
	Stargazers    struct {
		TotalCount githubv4.Int
	}
}

type qlUser struct {
	Login     githubv4.String
	Name      githubv4.String
	AvatarURL githubv4.String
	URL       githubv4.String
}

func gistFromQL(gist qlGist) Gist {
	return Gist{
		Name:        string(gist.Name),
		Description: string(gist.Description),
		URL:         string(gist.URL),
		CreatedAt:   gist.CreatedAt.Time,
	}
}

func pullRequestFromQL(pullRequest qlPullRequest) PullRequest {
	return PullRequest{
		Title:     string(pullRequest.Title),
		URL:       string(pullRequest.URL),
		State:     string(pullRequest.State),
		CreatedAt: pullRequest.CreatedAt.Time,
		Repo:      repoFromQL(pullRequest.Repository),
	}
}

func releaseFromQL(release qlRelease) Release {
	return Release{
		Name:        string(release.Nodes[0].Name),
		TagName:     string(release.Nodes[0].TagName),
		PublishedAt: release.Nodes[0].PublishedAt.Time,
		URL:         string(release.Nodes[0].URL),
	}
}

func repoFromQL(repo qlRepository) Repo {
	return Repo{
		Name:        string(repo.NameWithOwner),
		URL:         string(repo.URL),
		Description: string(repo.Description),
		Stargazers:  int(repo.Stargazers.TotalCount),
		IsPrivate:   bool(repo.IsPrivate),
	}
}

func userFromQL(user qlUser) User {
	return User{
		Login:     string(user.Login),
		Name:      string(user.Name),
		AvatarURL: string(user.AvatarURL),
		URL:       string(user.URL),
	}
}
