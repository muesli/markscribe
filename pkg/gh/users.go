package gh

import (
	"context"

	"github.com/muesli/markscribe/pkg/types"

	"github.com/shurcooL/githubv4"
)

// GetUsername retrieves the GitHub username and saves it for future requests.
func (c *Client) GetUsername() error {
	err := c.client.Query(context.Background(), &viewerQuery, nil)
	if err != nil {
		return err
	}

	c.username = viewerQuery.Viewer.Login
	return nil
}

// GetRecentFollowers returns recent 'n' followers from the user.
func (c *Client) GetRecentFollowers(n int) []types.User {
	users := []types.User{}

	variables := Map{
		"count":    githubv4.Int(n),
		"username": c.username,
	}

	err := c.client.Query(context.Background(), &recentFollowersQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentFollowersQuery.User.Followers.Edges {
		users = append(users, types.UserFromQL(v.Node))
	}

	return users
}
