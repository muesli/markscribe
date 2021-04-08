package main

import (
	"time"

	"github.com/shurcooL/githubv4"
)

type Contribution struct {
	OccurredAt time.Time
	Repo       Repo
}

type Gist struct {
	Name        string
	Description string
	URL         string
	CreatedAt   time.Time
}

type PullRequest struct {
	Title     string
	URL       string
	State     string
	CreatedAt time.Time
	Repo      Repo
}

type Release struct {
	Name        string
	TagName     string
	PublishedAt time.Time
	URL         string
}

type Repo struct {
	Name        string
	URL         string
	Description string
	Stargazers  int
	LastRelease Release
}

type Sponsor struct {
	User      User
	CreatedAt time.Time
}

type User struct {
	Login     string
	Name      string
	AvatarURL string
	URL       string
}

type Star struct {
	Name string
	URL  string
}

type QLGist struct {
	Name        githubv4.String
	Description githubv4.String
	URL         githubv4.String
	CreatedAt   githubv4.DateTime
}

type QLPullRequest struct {
	URL        githubv4.String
	Title      githubv4.String
	State      githubv4.PullRequestState
	CreatedAt  githubv4.DateTime
	Repository QLRepository
}

type QLRelease struct {
	Nodes []struct {
		Name         githubv4.String
		TagName      githubv4.String
		PublishedAt  githubv4.DateTime
		URL          githubv4.String
		IsPrerelease githubv4.Boolean
		IsDraft      githubv4.Boolean
	}
}

type QLRepository struct {
	NameWithOwner githubv4.String
	URL           githubv4.String
	Description   githubv4.String
	IsPrivate     githubv4.Boolean
	Stargazers    struct {
		TotalCount githubv4.Int
	}
}

type QLUser struct {
	Login     githubv4.String
	Name      githubv4.String
	AvatarURL githubv4.String
	URL       githubv4.String
}

type QLStar struct {
	NameWithOwner githubv4.String
	URL           githubv4.String
}

func GistFromQL(gist QLGist) Gist {
	return Gist{
		Name:        string(gist.Name),
		Description: string(gist.Description),
		URL:         string(gist.URL),
		CreatedAt:   gist.CreatedAt.Time,
	}
}

func PullRequestFromQL(pullRequest QLPullRequest) PullRequest {
	return PullRequest{
		Title:     string(pullRequest.Title),
		URL:       string(pullRequest.URL),
		State:     string(pullRequest.State),
		CreatedAt: pullRequest.CreatedAt.Time,
		Repo:      RepoFromQL(pullRequest.Repository),
	}
}

func ReleaseFromQL(release QLRelease) Release {
	return Release{
		Name:        string(release.Nodes[0].Name),
		TagName:     string(release.Nodes[0].TagName),
		PublishedAt: release.Nodes[0].PublishedAt.Time,
		URL:         string(release.Nodes[0].URL),
	}
}

func RepoFromQL(repo QLRepository) Repo {
	return Repo{
		Name:        string(repo.NameWithOwner),
		URL:         string(repo.URL),
		Description: string(repo.Description),
		Stargazers:  int(repo.Stargazers.TotalCount),
	}
}

func UserFromQL(user QLUser) User {
	return User{
		Login:     string(user.Login),
		Name:      string(user.Name),
		AvatarURL: string(user.AvatarURL),
		URL:       string(user.URL),
	}
}

func StarFromQL(star QLStar) Star {
	return Star{
		Name: string(star.NameWithOwner),
		URL:  string(star.URL),
	}
}
