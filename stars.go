package main

import (
	"context"

	"github.com/shurcooL/githubv4"
)

var recentStarsQuery struct {
	User struct {
		Login githubv4.String
		Stars struct {
			Nodes []QLRepository
		} `graphql:"starredRepositories(first: $count, orderBy: {field: STARRED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

func recentStars(count int) []Star {
	var stars []Star
	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"count":    githubv4.Int(count),
	}
	err := gitHubClient.Query(context.Background(), &recentStarsQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentStarsQuery.User.Stars.Nodes {
		s := Star{
			Repo: RepoFromQL(v),
		}
		stars = append(stars, s)
	}

	return stars
}

/*
{
  viewer {
    login
    starredRepositories(first: 3, orderBy: {field: STARRED_AT, direction: DESC}) {
      nodes {
        nameWithOwner
        url
        description
        stargazers {
          totalCount
        }
      }
    }
  }
}
*/
