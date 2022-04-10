package main

import (
	"github.com/KyleBanks/goodreads/responses"
)

func goodReadsReviews(count int) []responses.Review {
	reviews, err := goodReadsClient.ReviewList(goodReadsID, "read", "date_read", "", "d", 1, count)
	if err != nil {
		panic(err)
	}
	return reviews
}

func goodReadsCurrentlyReading(count int) []responses.Review {
	reviews, err := goodReadsClient.ReviewList(goodReadsID, "currently-reading", "date_updated", "", "d", 1, count)
	if err != nil {
		panic(err)
	}
	return reviews
}
