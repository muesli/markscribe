package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/shurcooL/githubv4"
)

type Repo struct {
	Name        string
	URL         string
	Description string
	Stargazers  int
	LastRelease Release
}

type Release struct {
	Name        string
	TagName     string
	PublishedAt time.Time
	URL         string
}

type Contribution struct {
	OccurredAt time.Time
	Repo       Repo
}

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
				Repository struct {
					NameWithOwner githubv4.String
					URL           githubv4.String
					Description   githubv4.String
					IsPrivate     githubv4.Boolean
				}
			} `graphql:"commitContributionsByRepository(maxRepositories: 100)"`
		}
	} `graphql:"user(login:$username)"`
}

var recentReposQuery struct {
	User struct {
		Login        githubv4.String
		Repositories struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   struct {
					NameWithOwner githubv4.String
					URL           githubv4.String
					Description   githubv4.String
				}
			}
		} `graphql:"repositories(first: $count, privacy: PUBLIC, isFork: false, ownerAffiliations: OWNER, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

var recentReleasesQuery struct {
	User struct {
		Login                     githubv4.String
		RepositoriesContributedTo struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   struct {
					NameWithOwner githubv4.String
					URL           githubv4.String
					Description   githubv4.String
					Stargazers    struct {
						TotalCount githubv4.Int
					}
					Releases struct {
						Nodes []struct {
							Name         githubv4.String
							TagName      githubv4.String
							PublishedAt  githubv4.DateTime
							URL          githubv4.String
							IsPrerelease githubv4.Boolean
							IsDraft      githubv4.Boolean
						}
					} `graphql:"releases(first: 10, orderBy: {field: CREATED_AT, direction: DESC})"`
				}
			}
		} `graphql:"repositoriesContributedTo(first: 100, after:$after includeUserRepositories: true, contributionTypes: COMMIT, privacy: PUBLIC)"`
	} `graphql:"user(login:$username)"`
}

func recentContributions(count int) []Contribution {
	// fmt.Printf("Finding recent contributions...\n")

	var contributions []Contribution
	variables := map[string]interface{}{
		"username": githubv4.String(username),
	}
	err := client.Query(context.Background(), &recentContributionsQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentContributionsQuery.User.ContributionsCollection.CommitContributionsByRepository {
		// ignore meta-repo
		if string(v.Repository.NameWithOwner) == fmt.Sprintf("%s/%s", username, username) {
			continue
		}
		if v.Repository.IsPrivate {
			continue
		}

		c := Contribution{
			Repo: Repo{
				Name:        string(v.Repository.NameWithOwner),
				URL:         string(v.Repository.URL),
				Description: string(v.Repository.Description),
			},
			OccurredAt: v.Contributions.Edges[0].Node.OccurredAt.Time,
		}

		contributions = append(contributions, c)
	}

	sort.Slice(contributions, func(i, j int) bool {
		return contributions[i].OccurredAt.After(contributions[j].OccurredAt)
	})

	// fmt.Printf("Found %d contributions!\n", len(repos))
	if len(contributions) > count {
		return contributions[:count]
	}
	return contributions
}

func recentRepos(count int) []Repo {
	// fmt.Printf("Finding recently created repos...\n")

	var repos []Repo
	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"count":    githubv4.Int(count + 1), // +1 in case we encounter the meta-repo itself
	}
	err := client.Query(context.Background(), &recentReposQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentReposQuery.User.Repositories.Edges {
		// ignore meta-repo
		if string(v.Node.NameWithOwner) == fmt.Sprintf("%s/%s", username, username) {
			continue
		}

		r := Repo{
			Name:        string(v.Node.NameWithOwner),
			URL:         string(v.Node.URL),
			Description: string(v.Node.Description),
		}

		repos = append(repos, r)
		if len(repos) == count {
			break
		}
	}

	// fmt.Printf("Found %d repos!\n", len(repos))
	return repos
}

func recentReleases(count int) []Repo {
	// fmt.Printf("Finding recent releases...\n")

	var after *githubv4.String
	var repos []Repo

	for {
		variables := map[string]interface{}{
			"username": githubv4.String(username),
			"after":    after,
		}
		err := client.Query(context.Background(), &recentReleasesQuery, variables)
		if err != nil {
			panic(err)
		}

		// fmt.Printf("%+v\n", query)
		if len(recentReleasesQuery.User.RepositoriesContributedTo.Edges) == 0 {
			break
		}

		for _, v := range recentReleasesQuery.User.RepositoriesContributedTo.Edges {
			r := Repo{
				Name:        string(v.Node.NameWithOwner),
				URL:         string(v.Node.URL),
				Description: string(v.Node.Description),
				Stargazers:  int(v.Node.Stargazers.TotalCount),
			}

			for _, rel := range v.Node.Releases.Nodes {
				if rel.IsPrerelease || rel.IsDraft {
					continue
				}
				if v.Node.Releases.Nodes[0].TagName == "" ||
					v.Node.Releases.Nodes[0].PublishedAt.Time.IsZero() {
					continue
				}
				r.LastRelease = Release{
					Name:        string(v.Node.Releases.Nodes[0].Name),
					TagName:     string(v.Node.Releases.Nodes[0].TagName),
					PublishedAt: v.Node.Releases.Nodes[0].PublishedAt.Time,
					URL:         string(v.Node.Releases.Nodes[0].URL),
				}
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

	// fmt.Printf("Found %d repos!\n", len(repos))
	if len(repos) > count {
		return repos[:count]
	}
	return repos
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
