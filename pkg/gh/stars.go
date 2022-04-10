package gh

import (
	"context"

	"github.com/muesli/markscribe/pkg/types"

	"github.com/shurcooL/githubv4"
)

// GetRecentStars retrieves 'n' recently starred repos.
func (c *Client) GetRecentStars(n int) []types.Star {
	starredRepos := []types.Star{}
	variables := Map{
		"count":    githubv4.Int(n),
		"username": c.username,
	}

	err := c.client.Query(context.Background(), &recentStarsQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentStarsQuery.User.Stars.Edges {
		starredRepos = append(starredRepos, types.Star{
			StarredAt: v.StarredAt.Time,
			Repo:      types.RepoFromQL(v.Node),
		})
	}

	return starredRepos
}
