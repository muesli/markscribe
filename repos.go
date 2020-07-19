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
					} `graphql:"releases(first: 1, orderBy: {field: CREATED_AT, direction: DESC})"`
				}
			}
		} `graphql:"repositoriesContributedTo(first: 100, after:$after includeUserRepositories: true, contributionTypes: COMMIT, privacy: PUBLIC)"`
	} `graphql:"user(login:$username)"`
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
	sm := make(map[string]Repo)

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
				r.LastRelease = Release{
					Name:        string(v.Node.Releases.Nodes[0].Name),
					TagName:     string(v.Node.Releases.Nodes[0].TagName),
					PublishedAt: v.Node.Releases.Nodes[0].PublishedAt.Time,
					URL:         string(v.Node.Releases.Nodes[0].URL),
				}
				break
			}

			sm[string(v.Node.NameWithOwner)] = r
			after = &v.Cursor
		}
	}

	// sort repos
	type kv struct {
		URL  string
		Repo Repo
	}

	var sl []kv
	for k, v := range sm {
		sl = append(sl, kv{k, v})
	}

	sort.Slice(sl, func(i, j int) bool {
		if sl[i].Repo.LastRelease.PublishedAt.Equal(sl[j].Repo.LastRelease.PublishedAt) {
			return sl[i].Repo.Stargazers > sl[j].Repo.Stargazers
		}
		return sl[i].Repo.LastRelease.PublishedAt.After(sl[j].Repo.LastRelease.PublishedAt)
	})

	var repos []Repo
	for i, kv := range sl {
		repos = append(repos, kv.Repo)
		if i == count-1 {
			break
		}
	}

	// fmt.Printf("Found %d repos!\n", len(repos))
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
*/
