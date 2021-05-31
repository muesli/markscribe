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
				Node      QLRepository
			}
		} `graphql:"starredRepositories(first: $count, orderBy: {field: STARRED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

func recentStars(count int) []Star {
	var starredRepos []Star
	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"count":    githubv4.Int(count),
	}
	err := gitHubClient.Query(context.Background(), &recentStarsQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentStarsQuery.User.Stars.Edges {
		starredRepos = append(starredRepos, Star{
			StarredAt: v.StarredAt.Time,
			Repo:      RepoFromQL(v.Node),
		})
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
