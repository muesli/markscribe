package github

import (
	"context"
	"fmt"
	"sort"

	domain "hufschlaeger.net/markscribe/internal/domain"
	"hufschlaeger.net/markscribe/internal/usecase/ports"
)

// Service wraps the GithubPort and contains app-level logic for GitHub features.
type Service struct {
	gh       ports.GithubPort
	username string
}

func New(gh ports.GithubPort, username string) *Service {
	return &Service{gh: gh, username: username}
}

// RecentRepos returns the most recent non-fork repositories owned by the user,
// excluding the meta repo "username/username".
func (s *Service) RecentRepos(count int) []domain.Repo {
	repos, err := s.gh.RecentRepos(context.Background(), s.username, count+1, false)
	if err != nil {
		panic(err)
	}
	var out []domain.Repo
	for _, r := range repos {
		if r.Name == fmt.Sprintf("%s/%s", s.username, s.username) {
			continue
		}
		out = append(out, r)
		if len(out) == count {
			break
		}
	}
	return out
}

// RecentForks returns the most recent forked repositories for the user,
// excluding the meta repo "username/username".
func (s *Service) RecentForks(count int) []domain.Repo {
	repos, err := s.gh.RecentRepos(context.Background(), s.username, count+1, true)
	if err != nil {
		panic(err)
	}
	var out []domain.Repo
	for _, r := range repos {
		if r.Name == fmt.Sprintf("%s/%s", s.username, s.username) {
			continue
		}
		out = append(out, r)
		if len(out) == count {
			break
		}
	}
	return out
}

// Repo returns details for a repository.
func (s *Service) Repo(owner, name string) domain.Repo {
	r, err := s.gh.Repo(context.Background(), owner, name)
	if err != nil {
		panic(err)
	}
	return r
}

// Followers returns a list of followers for the configured user.
func (s *Service) Followers(count int) []domain.User {
	users, err := s.gh.Followers(context.Background(), s.username, count)
	if err != nil {
		panic(err)
	}
	return users
}

// RecentPullRequests returns recent pull requests created by the user,
// excluding the meta repo "username/username" and private repositories.
func (s *Service) RecentPullRequests(count int) []domain.PullRequest {
	prs, err := s.gh.RecentPullRequests(context.Background(), s.username, count+1)
	if err != nil {
		panic(err)
	}
	var out []domain.PullRequest
	meta := fmt.Sprintf("%s/%s", s.username, s.username)
	for _, pr := range prs {
		if pr.Repo.Name == meta {
			continue
		}
		if pr.Repo.IsPrivate {
			continue
		}
		out = append(out, pr)
		if len(out) == count {
			break
		}
	}
	return out
}

// RecentReleases returns repositories with the most recent valid releases,
// sorted by PublishedAt desc, then Stargazers desc, limited to count.
func (s *Service) RecentReleases(count int) []domain.Repo {
	repos, err := s.gh.RecentReleases(context.Background(), s.username, count)
	if err != nil {
		panic(err)
	}
	// sort as in legacy implementation
	sort.Slice(repos, func(i, j int) bool {
		if repos[i].LastRelease.PublishedAt.Equal(repos[j].LastRelease.PublishedAt) {
			return repos[i].Stargazers > repos[j].Stargazers
		}
		return repos[i].LastRelease.PublishedAt.After(repos[j].LastRelease.PublishedAt)
	})
	if len(repos) > count {
		return repos[:count]
	}
	return repos
}

// RecentContributions returns recent commit contributions by repository for the user,
// excluding the meta repo and private repositories, sorted by time desc and limited to count.
func (s *Service) RecentContributions(count int) []domain.Contribution {
	cons, err := s.gh.RecentContributions(context.Background(), s.username, count+10) // fetch a few extra for filtering
	if err != nil {
		panic(err)
	}
	meta := fmt.Sprintf("%s/%s", s.username, s.username)
	var out []domain.Contribution
	for _, c := range cons {
		if c.Repo.Name == meta {
			continue
		}
		if c.Repo.IsPrivate {
			continue
		}
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].OccurredAt.After(out[j].OccurredAt) })
	if len(out) > count {
		out = out[:count]
	}
	return out
}

// Gists returns user's gists ordered by creation date desc limited by count.
func (s *Service) Gists(count int) []domain.Gist {
	gists, err := s.gh.Gists(context.Background(), s.username, count)
	if err != nil {
		panic(err)
	}
	return gists
}

// RecentStars returns recently starred public repositories.
func (s *Service) RecentStars(count int) []domain.Star {
	stars, err := s.gh.RecentStars(context.Background(), s.username, count)
	if err != nil {
		panic(err)
	}
	return stars
}

// RecentIssues returns recent issue contributions grouped by repository,
// excluding the meta repo and private repositories, sorted by time desc and limited to count.
func (s *Service) RecentIssues(count int) []domain.Issue {
	issues, err := s.gh.RecentIssues(context.Background(), s.username, count+10)
	if err != nil {
		panic(err)
	}
	meta := fmt.Sprintf("%s/%s", s.username, s.username)
	var out []domain.Issue
	for _, is := range issues {
		if is.Repo.Name == meta {
			continue
		}
		if is.Repo.IsPrivate {
			continue
		}
		out = append(out, is)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].OccurredAt.After(out[j].OccurredAt) })
	if len(out) > count {
		out = out[:count]
	}
	return out
}

// Sponsors returns the most recent sponsors up to count.
func (s *Service) Sponsors(count int) []domain.Sponsor {
	sponsors, err := s.gh.Sponsors(context.Background(), s.username, count)
	if err != nil {
		panic(err)
	}
	if len(sponsors) > count {
		sponsors = sponsors[:count]
	}
	return sponsors
}
