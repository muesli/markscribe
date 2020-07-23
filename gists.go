package main

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

type Gist struct {
	Name        string
	Description string
	URL         string
	CreatedAt   time.Time
}

var gistsQuery struct {
	User struct {
		Login githubv4.String
		Gists struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   struct {
					Name        githubv4.String
					Description githubv4.String
					URL         githubv4.String
					CreatedAt   githubv4.DateTime
				}
			}
		} `graphql:"gists(first: $count, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

func gists(count int) []Gist {
	// fmt.Printf("Finding gists...\n")

	var gists []Gist
	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"count":    githubv4.Int(count),
	}
	err := client.Query(context.Background(), &gistsQuery, variables)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%+v\n", query)
	for _, v := range gistsQuery.User.Gists.Edges {
		g := Gist{
			Name:        string(v.Node.Name),
			Description: string(v.Node.Description),
			URL:         string(v.Node.URL),
			CreatedAt:   v.Node.CreatedAt.Time,
		}
		gists = append(gists, g)
	}

	// fmt.Printf("Found %d gists!\n", len(gists))
	return gists
}

/*
{
  user(login: "muesli") {
    login
    gists(first: 100) {
      totalCount
      edges {
        cursor
        node {
		  name
		  description
		  url
		  createdAt
        }
      }
    }
  }
}
*/
