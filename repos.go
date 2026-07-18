package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/shurcooL/githubv4"
)

var recentContributionsQuery struct {
	User struct {
		Login                   githubv4.String
		ContributionsCollection struct {
			CommitContributionsByRepository []struct {
				Contributions struct {
					Edges []struct {
						Cursor githubv4.String
						Node   struct {
							OccurredAt githubv4.DateTime
						}
					}
				} `graphql:"contributions(first: 1)"`
				Repository qlRepository
			} `graphql:"commitContributionsByRepository(maxRepositories: 100)"`
		} `graphql:"contributionsCollection(from: $from, to: $to)"`
	} `graphql:"user(login:$username)"`
}

var recentPullRequestsQuery struct {
	User struct {
		Login        githubv4.String
		PullRequests struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   qlPullRequest
			}
		} `graphql:"pullRequests(first: $count, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

var recentReposQuery struct {
	User struct {
		Login        githubv4.String
		Repositories struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   qlRepository
			}
		} `graphql:"repositories(first: $count, privacy: PUBLIC, isFork: $isFork, ownerAffiliations: OWNER, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

var repoReleasesQuery struct {
	Repository struct {
		qlRepository
		Releases qlRelease `graphql:"releases(first: 10, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"repository(owner:$owner, name:$name)"`
}

var repoQuery struct {
	Repository struct {
		Description   githubv4.String
		NameWithOwner githubv4.String
		IsPrivate     githubv4.Boolean
		URL           githubv4.String
		Stargazers    struct {
			TotalCount githubv4.Int
		}
		Releases qlRelease `graphql:"releases(last: 1)"`
	} `graphql:"repository(owner:$owner, name:$name)"`
}

type repoContribution struct {
	Repository qlRepository
	OccurredAt time.Time
}

// contributedRepos returns the repositories the user committed to over the
// past year, most recent contribution first. The GitHub GraphQL API enforces
// a resource limit on contributionsCollection that the default one-year
// window can exceed, so query in six-month windows and merge the results.
func contributedRepos() []repoContribution {
	var contributions []repoContribution
	latestByRepo := make(map[string]int)
	now := time.Now()
	for _, monthsAgo := range []int{6, 12} {
		variables := map[string]interface{}{
			"username": githubv4.String(username),
			"from":     githubv4.DateTime{Time: now.AddDate(0, -monthsAgo, 0)},
			"to":       githubv4.DateTime{Time: now.AddDate(0, -monthsAgo+6, 0)},
		}
		err := gitHubClient.Query(context.Background(), &recentContributionsQuery, variables)
		if err != nil {
			panic(err)
		}

		for _, v := range recentContributionsQuery.User.ContributionsCollection.CommitContributionsByRepository {
			if len(v.Contributions.Edges) == 0 {
				continue
			}

			c := repoContribution{
				Repository: v.Repository,
				OccurredAt: v.Contributions.Edges[0].Node.OccurredAt.Time,
			}

			if i, ok := latestByRepo[string(c.Repository.NameWithOwner)]; ok {
				if c.OccurredAt.After(contributions[i].OccurredAt) {
					contributions[i] = c
				}
				continue
			}
			latestByRepo[string(c.Repository.NameWithOwner)] = len(contributions)
			contributions = append(contributions, c)
		}
	}

	sort.Slice(contributions, func(i, j int) bool {
		return contributions[i].OccurredAt.After(contributions[j].OccurredAt)
	})

	return contributions
}

func recentContributions(count int) []Contribution {
	// fmt.Printf("Finding recent contributions...\n")

	var contributions []Contribution
	for _, v := range contributedRepos() {
		// ignore meta-repo
		if string(v.Repository.NameWithOwner) == fmt.Sprintf("%s/%s", username, username) {
			continue
		}
		if v.Repository.IsPrivate {
			continue
		}

		contributions = append(contributions, Contribution{
			Repo:       repoFromQL(v.Repository),
			OccurredAt: v.OccurredAt,
		})
	}

	// fmt.Printf("Found %d contributions!\n", len(repos))
	if len(contributions) > count {
		return contributions[:count]
	}
	return contributions
}

func recentPullRequests(count int) []PullRequest {
	// fmt.Printf("Finding recently created pullRequests...\n")

	var pullRequests []PullRequest
	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"count":    githubv4.Int(count + 1), // +1 in case we encounter the meta-repo itself
	}
	err := gitHubClient.Query(context.Background(), &recentPullRequestsQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentPullRequestsQuery.User.PullRequests.Edges {
		// ignore meta-repo
		if string(v.Node.Repository.NameWithOwner) == fmt.Sprintf("%s/%s", username, username) {
			continue
		}
		if v.Node.Repository.IsPrivate {
			continue
		}

		pullRequests = append(pullRequests, pullRequestFromQL(v.Node))
		if len(pullRequests) == count {
			break
		}
	}

	// fmt.Printf("Found %d pullRequests!\n", len(pullRequests))
	return pullRequests
}

func recentRepos(count int) []Repo {
	// fmt.Printf("Finding recently created repos...\n")

	var repos []Repo
	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"count":    githubv4.Int(count + 1), // +1 in case we encounter the meta-repo itself
		"isFork":   githubv4.Boolean(false),
	}
	err := gitHubClient.Query(context.Background(), &recentReposQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentReposQuery.User.Repositories.Edges {
		// ignore meta-repo
		if string(v.Node.NameWithOwner) == fmt.Sprintf("%s/%s", username, username) {
			continue
		}

		repos = append(repos, repoFromQL(v.Node))
		if len(repos) == count {
			break
		}
	}

	// fmt.Printf("Found %d repos!\n", len(repos))
	return repos
}

func recentForks(count int) []Repo {
	// fmt.Printf("Finding recently created repos...\n")

	var repos []Repo
	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"count":    githubv4.Int(count + 1), // +1 in case we encounter the meta-repo itself
		"isFork":   githubv4.Boolean(true),
	}
	err := gitHubClient.Query(context.Background(), &recentReposQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentReposQuery.User.Repositories.Edges {
		// ignore meta-repo
		if string(v.Node.NameWithOwner) == fmt.Sprintf("%s/%s", username, username) {
			continue
		}

		repos = append(repos, repoFromQL(v.Node))
		if len(repos) == count {
			break
		}
	}

	// fmt.Printf("Found %d repos!\n", len(repos))
	return repos
}

func recentReleases(count int) []Repo {
	// fmt.Printf("Finding recent releases...\n")

	var repos []Repo

	for _, c := range contributedRepos() {
		if c.Repository.IsPrivate {
			continue
		}

		owner, name, ok := strings.Cut(string(c.Repository.NameWithOwner), "/")
		if !ok {
			continue
		}
		variables := map[string]interface{}{
			"owner": githubv4.String(owner),
			"name":  githubv4.String(name),
		}
		err := gitHubClient.Query(context.Background(), &repoReleasesQuery, variables)
		if err != nil {
			panic(err)
		}

		r := repoFromQL(repoReleasesQuery.Repository.qlRepository)

		for _, rel := range repoReleasesQuery.Repository.Releases.Nodes {
			if rel.IsPrerelease || rel.IsDraft {
				continue
			}
			if repoReleasesQuery.Repository.Releases.Nodes[0].TagName == "" ||
				repoReleasesQuery.Repository.Releases.Nodes[0].PublishedAt.Time.IsZero() {
				continue
			}
			r.LastRelease = releaseFromQL(repoReleasesQuery.Repository.Releases)
			break
		}

		if !r.LastRelease.PublishedAt.IsZero() {
			repos = append(repos, r)
		}
	}

	sort.Slice(repos, func(i, j int) bool {
		if repos[i].LastRelease.PublishedAt.Equal(repos[j].LastRelease.PublishedAt) {
			return repos[i].Stargazers > repos[j].Stargazers
		}
		return repos[i].LastRelease.PublishedAt.After(repos[j].LastRelease.PublishedAt)
	})

	// fmt.Printf("Found %d repos!\n", len(repos))
	if len(repos) > count {
		return repos[:count]
	}
	return repos
}

func repo(owner, name string) Repo {
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}
	err := gitHubClient.Query(context.Background(), &repoQuery, variables)
	if err != nil {
		panic(err)
	}
	repo := repoQuery.Repository
	return Repo{
		Name:        string(repo.NameWithOwner),
		URL:         string(repo.URL),
		Description: string(repo.Description),
		Stargazers:  int(repo.Stargazers.TotalCount),
		IsPrivate:   bool(repo.IsPrivate),
		LastRelease: releaseFromQL(repo.Releases),
	}
}

/*
{
  user(login: "muesli") {
    login
    repositoriesContributedTo(first: 100, includeUserRepositories: true, contributionTypes: COMMIT) {
      totalCount
      edges {
        cursor
        node {
          id
          nameWithOwner
        }
      }
    }
  }
}

{
  user(login: "muesli") {
    login
    repositoriesContributedTo(first: 100, includeUserRepositories: true, contributionTypes: COMMIT) {
      totalCount
      edges {
        cursor
        node {
          id
          nameWithOwner
		  releases(first: 3, orderBy: {field: CREATED_AT, direction: DESC}) {
          	nodes {
          	  name
              PublishedAt
			  url
			  isPrerelease
			  isDraft
            }
          }
        }
      }
    }
  }
}

{
  user(login: "muesli") {
    login
    repositories(first: 10, privacy: PUBLIC, isFork: false, ownerAffiliations: OWNER, orderBy: {field: CREATED_AT, direction: DESC}) {
      totalCount
      edges {
        cursor
        node {
          id
          nameWithOwner
        }
      }
    }
  }
}

{
  user(login: "muesli") {
    login
    contributionsCollection {
      commitContributionsByRepository {
        contributions(first: 1) {
          edges {
            cursor
            node {
              occurredAt
            }
          }
        }
        repository {
          id
		  nameWithOwner
		  url
		  description
        }
      }
    }
  }
}
*/
