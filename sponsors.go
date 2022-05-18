package main

import (
	"context"

	"github.com/shurcooL/githubv4"
)

var sponsorsQuery struct {
	User struct {
		Login                    githubv4.String
		SponsorshipsAsMaintainer struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   struct {
					CreatedAt     githubv4.DateTime
					SponsorEntity struct {
						Typename     githubv4.String `graphql:"__typename"`
						User         qlUser          `graphql:"... on User"`
						Organization qlUser          `graphql:"... on Organization"`
					}
				}
			}
		} `graphql:"sponsorshipsAsMaintainer(first: $count, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

func sponsors(count int) []Sponsor {
	// fmt.Printf("Finding sponsors...\n")

	var sponsors []Sponsor
	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"count":    githubv4.Int(count),
	}
	err := gitHubClient.Query(context.Background(), &sponsorsQuery, variables)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%+v\n", query)

	for _, v := range sponsorsQuery.User.SponsorshipsAsMaintainer.Edges {
		switch v.Node.SponsorEntity.Typename {
		case "User":
			sponsors = append(sponsors, Sponsor{
				User:      userFromQL(v.Node.SponsorEntity.User),
				CreatedAt: v.Node.CreatedAt.Time,
			})
		case "Organization":
			sponsors = append(sponsors, Sponsor{
				User:      userFromQL(v.Node.SponsorEntity.Organization),
				CreatedAt: v.Node.CreatedAt.Time,
			})
		}
	}

	// fmt.Printf("Found %d sponsors!\n", len(users))
	return sponsors
}

/*
{
  user(login: "muesli") {
    login
    sponsorshipsAsMaintainer(first: 100) {
      totalCount
      edges {
        cursor
        node {
          createdAt
          sponsorEntity {
            __typename
            ... on User {
              login
              name
              avatarUrl
              url
            }
            ... on Organization {
              login
              name
              avatarUrl
              url
            }
          }
        }
      }
    }
  }
}
*/
