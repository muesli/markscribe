package goodreads

import (
	"context"

	"github.com/KyleBanks/goodreads/responses"
	"hufschlaeger.net/markscribe/internal/usecase/ports"
)

// Service wraps the GoodReadsPort and contains app-level logic for Goodreads features.
type Service struct {
	gr ports.GoodReadsPort
}

func New(gr ports.GoodReadsPort) *Service { return &Service{gr: gr} }

func (s *Service) Reviews(count int) []responses.Review {
	reviews, err := s.gr.Reviews(context.Background(), count)
	if err != nil {
		panic(err)
	}
	return reviews
}

func (s *Service) CurrentlyReading(count int) []responses.Review {
	reviews, err := s.gr.CurrentlyReading(context.Background(), count)
	if err != nil {
		panic(err)
	}
	return reviews
}
