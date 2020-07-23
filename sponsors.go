package main

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

type Sponsor struct {
	User      User
	CreatedAt time.Time
}

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
						SponsorUser struct {
							Login     githubv4.String
							Name      githubv4.String
							AvatarURL githubv4.String
							URL       githubv4.String
						} `graphql:"... on User"`
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
	err := client.Query(context.Background(), &sponsorsQuery, variables)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%+v\n", query)

	for _, v := range sponsorsQuery.User.SponsorshipsAsMaintainer.Edges {
		s := Sponsor{
			User: User{
				Login:     string(v.Node.SponsorEntity.SponsorUser.Login),
				Name:      string(v.Node.SponsorEntity.SponsorUser.Name),
				AvatarURL: string(v.Node.SponsorEntity.SponsorUser.AvatarURL),
				URL:       string(v.Node.SponsorEntity.SponsorUser.URL),
			},
			CreatedAt: v.Node.CreatedAt.Time,
		}
		sponsors = append(sponsors, s)
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
            ... on User {
			  login
			  name
			  avatar
			  url
            }
          }
        }
      }
    }
  }
}
*/
