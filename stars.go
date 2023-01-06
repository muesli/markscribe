package main

import (
	"context"

	"github.com/shurcooL/githubv4"
)

var recentStarsQuery struct {
	User struct {
		Login githubv4.String
		Stars struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor    githubv4.String
				StarredAt githubv4.DateTime
				Node      qlRepository
			}
		} `graphql:"starredRepositories(first: $count, after:$after, orderBy: {field: STARRED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

func recentStars(count int) []Star {
	var starredRepos []Star
	var after *githubv4.String

	for {
		if len(starredRepos) >= count {
			break
		}
		variables := map[string]interface{}{
			"username": githubv4.String(username),
			"count":    githubv4.Int(count),
			"after":    after,
		}
		err := gitHubClient.Query(context.Background(), &recentStarsQuery, variables)
		if err != nil {
			panic(err)
		}

		for _, v := range recentStarsQuery.User.Stars.Edges {
			if v.Node.IsPrivate {
				continue
			}
			if len(starredRepos) >= count {
				break
			}
			starredRepos = append(starredRepos, Star{
				StarredAt: v.StarredAt.Time,
				Repo:      repoFromQL(v.Node),
			})
			after = githubv4.NewString(v.Cursor)
		}
	}

	return starredRepos
}

/*
{
	viewer {
		login
		starredRepositories(first: 3, orderBy: {field: STARRED_AT, direction: DESC}) {
			totalCount
			edges {
				cursor
				starredAt
				node {
					nameWithOwner
					url
					description
				}
			}
		}
	}
}
*/
