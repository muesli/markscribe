package gh

import (
	"context"
	"net/http"
	"os"

	"github.com/shurcooL/githubv4"

	"golang.org/x/oauth2"
)

type Client struct {
	client   *githubv4.Client
	username githubv4.String
}

type Map map[string]interface{}

// New creates and returns a new GitHub client.
func New() *Client {
	gitHubToken := os.Getenv("GITHUB_TOKEN")
	httpClient := http.DefaultClient

	if len(gitHubToken) > 0 {
		httpClient = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: gitHubToken},
		))
	}

	return &Client{
		client: githubv4.NewClient(httpClient),
	}
}
