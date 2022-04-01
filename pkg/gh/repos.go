package gh

import (
	"context"
	"fmt"
	"sort"

	"github.com/muesli/markscribe/pkg/types"

	"github.com/shurcooL/githubv4"
)

// GetRecentContributions retrieves 'n' recent contributions.
func (c *Client) GetRecentContributions(n int) []types.Contribution {
	contributions := []types.Contribution{}
	variables := Map{
		"username": c.username,
	}

	err := c.client.Query(context.Background(), &recentContributionsQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentContributionsQuery.User.ContributionsCollection.CommitContributionsByRepository {
		// Ignore meta-repo
		if string(v.Repository.NameWithOwner) == fmt.Sprintf("%s/%s", c.username, c.username) {
			continue
		}

		// Ignore private repo
		if v.Repository.IsPrivate {
			continue
		}

		c := types.Contribution{
			Repo:       types.RepoFromQL(v.Repository),
			OccurredAt: v.Contributions.Edges[0].Node.OccurredAt.Time,
		}

		contributions = append(contributions, c)
	}

	sort.Slice(contributions, func(i, j int) bool {
		return contributions[i].OccurredAt.After(contributions[j].OccurredAt)
	})

	if len(contributions) > n {
		return contributions[:n]
	}
	return contributions
}

// GetRecentPullRequests retrieves 'n' recent pull requests.
func (c *Client) GetRecentPullRequests(n int) []types.PullRequest {
	pullRequests := []types.PullRequest{}
	variables := Map{
		"count":    githubv4.Int(n + 1),
		"username": c.username,
	}

	err := c.client.Query(context.Background(), &recentPullRequestsQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentPullRequestsQuery.User.PullRequests.Edges {
		// Ignore meta-repo
		if string(v.Node.Repository.NameWithOwner) == fmt.Sprintf("%s/%s", c.username, c.username) {
			continue
		}

		// Ignore private repo
		if v.Node.Repository.IsPrivate {
			continue
		}

		pullRequests = append(pullRequests, types.PullRequestFromQL(v.Node))
		if len(pullRequests) == n {
			break
		}
	}

	return pullRequests
}

// GetRecentRepos retrieves 'n' recent repos.
func (c *Client) GetRecentRepos(n int) []types.Repo {
	repos := []types.Repo{}
	variables := Map{
		"count":    githubv4.Int(n + 1),
		"username": c.username,
	}

	err := c.client.Query(context.Background(), &recentReposQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentReposQuery.User.Repositories.Edges {
		// Ignore meta-repo
		if string(v.Node.NameWithOwner) == fmt.Sprintf("%s/%s", c.username, c.username) {
			continue
		}

		repos = append(repos, types.RepoFromQL(v.Node))
		if len(repos) == n {
			break
		}
	}

	return repos
}

// GetRecentReleases retrieves 'n' recent releases.
func (c *Client) GetRecentReleases(n int) []types.Repo {
	var after *githubv4.String
	var repos []types.Repo

	for {
		variables := Map{
			"username": c.username,
			"after":    after,
		}
		err := c.client.Query(context.Background(), &recentReleasesQuery, variables)
		if err != nil {
			panic(err)
		}

		if len(recentReleasesQuery.User.RepositoriesContributedTo.Edges) == 0 {
			break
		}

		for _, v := range recentReleasesQuery.User.RepositoriesContributedTo.Edges {
			r := types.RepoFromQL(v.Node.QLRepository)

			for _, rel := range v.Node.Releases.Nodes {
				if rel.IsPrerelease || rel.IsDraft {
					continue
				}
				if v.Node.Releases.Nodes[0].TagName == "" ||
					v.Node.Releases.Nodes[0].PublishedAt.Time.IsZero() {
					continue
				}
				r.LastRelease = types.ReleaseFromQL(v.Node.Releases)
				break
			}

			if !r.LastRelease.PublishedAt.IsZero() {
				repos = append(repos, r)
			}

			after = &v.Cursor
		}
	}

	sort.Slice(repos, func(i, j int) bool {
		if repos[i].LastRelease.PublishedAt.Equal(repos[j].LastRelease.PublishedAt) {
			return repos[i].Stargazers > repos[j].Stargazers
		}
		return repos[i].LastRelease.PublishedAt.After(repos[j].LastRelease.PublishedAt)
	})

	if len(repos) > n {
		return repos[:n]
	}
	return repos
}
