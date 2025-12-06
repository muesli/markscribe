package goodreadsadapter

import (
	"context"

	"github.com/KyleBanks/goodreads"
	"github.com/KyleBanks/goodreads/responses"
)

// Adapter implements ports.GoodReadsPort using the KyleBanks goodreads client.
type Adapter struct {
	client *goodreads.Client
	userID string
}

func New(client *goodreads.Client, userID string) *Adapter {
	return &Adapter{client: client, userID: userID}
}

// Reviews returns finished reviews from the "read" shelf.
func (a *Adapter) Reviews(ctx context.Context, count int) ([]responses.Review, error) {
	// API does not take ctx, but keep signature for port compatibility.
	reviews, err := a.client.ReviewList(a.userID, "read", "date_read", "", "d", 1, count)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

// CurrentlyReading returns reviews from the "currently-reading" shelf.
func (a *Adapter) CurrentlyReading(ctx context.Context, count int) ([]responses.Review, error) {
	reviews, err := a.client.ReviewList(a.userID, "currently-reading", "date_updated", "", "d", 1, count)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
