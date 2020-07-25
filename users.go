package main

import (
	"context"

	"github.com/shurcooL/githubv4"
)

var viewerQuery struct {
	Viewer struct {
		Login githubv4.String
	}
}

var recentFollowersQuery struct {
	User struct {
		Login     githubv4.String
		Followers struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   QLUser
			}
		} `graphql:"followers(first: $count)"`
	} `graphql:"user(login:$username)"`
}

func getUsername() (string, error) {
	err := client.Query(context.Background(), &viewerQuery, nil)
	if err != nil {
		return "", err
	}

	return string(viewerQuery.Viewer.Login), nil
}

func recentFollowers(count int) []User {
	// fmt.Printf("Finding recent followers...\n")

	var users []User
	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"count":    githubv4.Int(count),
	}
	err := client.Query(context.Background(), &recentFollowersQuery, variables)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%+v\n", query)
	for _, v := range recentFollowersQuery.User.Followers.Edges {
		users = append(users, UserFromQL(v.Node))
	}

	// fmt.Printf("Found %d recent followers!\n", len(users))
	return users
}

/*
{
  user(login: "muesli") {
    login
    followers(first: 10) {
      totalCount
      edges {
        cursor
        node {
          id
          avatarUrl
          login
		  name
		  url
        }
      }
    }
  }
}
*/
