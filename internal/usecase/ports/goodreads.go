package ports

import (
	"context"

	"github.com/KyleBanks/goodreads/responses"
)

// GoodReadsPort defines the minimal GoodReads operations used by the app.
type GoodReadsPort interface {
	// Reviews returns the latest finished reviews (shelf: read)
	Reviews(ctx context.Context, count int) ([]responses.Review, error)
	// CurrentlyReading returns the latest currently reading reviews
	CurrentlyReading(ctx context.Context, count int) ([]responses.Review, error)
}
