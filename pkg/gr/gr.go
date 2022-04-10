package gr

import (
	"os"

	"github.com/KyleBanks/goodreads"
	"github.com/KyleBanks/goodreads/responses"
)

type Client struct {
	client *goodreads.Client
	token  string
	id     string
}

// New creates and returns a new GoodReads client.
func New() *Client {
	token := os.Getenv("GOODREADS_TOKEN")
	id := os.Getenv("GOODREADS_USER_ID")

	return &Client{
		client: goodreads.NewClient(token),
		token:  token,
		id:     id,
	}
}

// GetReviews retrieves 'n' reviews.
func (c *Client) GetReviews(n int) []responses.Review {
	reviews, err := c.client.ReviewList(c.id, "read", "date_read", "", "d", 1, n)
	if err != nil {
		panic(err)
	}
	return reviews
}

// GetCurrentlyReading retrieves 'n' currently read books.
func (c *Client) GetCurrentlyReading(n int) []responses.Review {
	reviews, err := c.client.ReviewList(c.id, "currently-reading", "date_updated", "", "d", 1, n)
	if err != nil {
		panic(err)
	}
	return reviews
}
