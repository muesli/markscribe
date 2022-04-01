package gh

import (
	"context"

	"github.com/muesli/markscribe/pkg/types"

	"github.com/shurcooL/githubv4"
)

// GetRecentGists retrieves 'n' recent gists.
func (c *Client) GetRecentGists(n int) []types.Gist {
	gists := []types.Gist{}
	variables := map[string]interface{}{
		"count":    githubv4.Int(n),
		"username": c.username,
	}
	err := c.client.Query(context.Background(), &gistsQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range gistsQuery.User.Gists.Edges {
		gists = append(gists, types.GistFromQL(v.Node))
	}

	return gists
}
