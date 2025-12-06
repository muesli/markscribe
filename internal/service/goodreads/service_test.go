package goodreads

import (
	"context"
	"errors"
	"testing"

	"github.com/KyleBanks/goodreads/responses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGoodReadsPort is a mock implementation of ports.GoodReadsPort
type MockGoodReadsPort struct {
	mock.Mock
}

func (m *MockGoodReadsPort) Reviews(ctx context.Context, count int) ([]responses.Review, error) {
	args := m.Called(ctx, count)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]responses.Review), args.Error(1)
}

func (m *MockGoodReadsPort) CurrentlyReading(ctx context.Context, count int) ([]responses.Review, error) {
	args := m.Called(ctx, count)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]responses.Review), args.Error(1)
}

func TestNew(t *testing.T) {
	mockGR := new(MockGoodReadsPort)

	svc := New(mockGR)

	assert.NotNil(t, svc)
	assert.Equal(t, mockGR, svc.gr)
}

func TestService_Reviews(t *testing.T) {
	tests := []struct {
		name          string
		count         int
		mockReviews   []responses.Review
		mockError     error
		expectedPanic bool
	}{
		{
			name:  "successful retrieval with multiple reviews",
			count: 3,
			mockReviews: []responses.Review{
				{
					Book: responses.AuthorBook{
						Title: "The Great Gatsby",
						ID:    "123",
					},
					Rating: 5,
				},
				{
					Book: responses.AuthorBook{
						Title: "1984",
						ID:    "456",
					},
					Rating: 4,
				},
				{
					Book: responses.AuthorBook{
						Title: "To Kill a Mockingbird",
						ID:    "789",
					},
					Rating: 5,
				},
			},
			expectedPanic: false,
		},
		{
			name:          "successful retrieval with empty reviews",
			count:         5,
			mockReviews:   []responses.Review{},
			expectedPanic: false,
		},
		{
			name:          "successful retrieval with single review",
			count:         1,
			mockReviews:   []responses.Review{{Book: responses.AuthorBook{Title: "Test Book"}, Rating: 3}},
			expectedPanic: false,
		},
		{
			name:          "panics on error",
			count:         5,
			mockReviews:   nil,
			mockError:     errors.New("api error"),
			expectedPanic: true,
		},
		{
			name:          "panics on network error",
			count:         3,
			mockReviews:   nil,
			mockError:     errors.New("network timeout"),
			expectedPanic: true,
		},
		{
			name:          "panics on authentication error",
			count:         2,
			mockReviews:   nil,
			mockError:     errors.New("unauthorized"),
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGR := new(MockGoodReadsPort)
			mockGR.On("Reviews", mock.Anything, tt.count).
				Return(tt.mockReviews, tt.mockError)

			svc := New(mockGR)

			if tt.expectedPanic {
				assert.Panics(t, func() {
					svc.Reviews(tt.count)
				}, "Expected Reviews to panic on error")
				mockGR.AssertExpectations(t)
				return
			}

			result := svc.Reviews(tt.count)

			assert.NotNil(t, result)
			assert.Equal(t, len(tt.mockReviews), len(result))
			assert.Equal(t, tt.mockReviews, result)
			mockGR.AssertExpectations(t)
		})
	}
}

func TestService_CurrentlyReading(t *testing.T) {
	tests := []struct {
		name          string
		count         int
		mockReviews   []responses.Review
		mockError     error
		expectedPanic bool
	}{
		{
			name:  "successful retrieval with multiple currently reading books",
			count: 2,
			mockReviews: []responses.Review{
				{
					Book: responses.AuthorBook{
						Title: "Clean Code",
						ID:    "111",
					},
					Rating: 0, // Not rated yet
				},
				{
					Book: responses.AuthorBook{
						Title: "Design Patterns",
						ID:    "222",
					},
					Rating: 0,
				},
			},
			expectedPanic: false,
		},
		{
			name:          "successful retrieval with empty currently reading list",
			count:         5,
			mockReviews:   []responses.Review{},
			expectedPanic: false,
		},
		{
			name:  "successful retrieval with single currently reading book",
			count: 1,
			mockReviews: []responses.Review{
				{
					Book: responses.AuthorBook{
						Title: "The Pragmatic Programmer",
						ID:    "333",
					},
					Rating: 0,
				},
			},
			expectedPanic: false,
		},
		{
			name:          "panics on error",
			count:         3,
			mockReviews:   nil,
			mockError:     errors.New("api error"),
			expectedPanic: true,
		},
		{
			name:          "panics on network error",
			count:         2,
			mockReviews:   nil,
			mockError:     errors.New("connection refused"),
			expectedPanic: true,
		},
		{
			name:          "panics on rate limit error",
			count:         5,
			mockReviews:   nil,
			mockError:     errors.New("rate limit exceeded"),
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGR := new(MockGoodReadsPort)
			mockGR.On("CurrentlyReading", mock.Anything, tt.count).
				Return(tt.mockReviews, tt.mockError)

			svc := New(mockGR)

			if tt.expectedPanic {
				assert.Panics(t, func() {
					svc.CurrentlyReading(tt.count)
				}, "Expected CurrentlyReading to panic on error")
				mockGR.AssertExpectations(t)
				return
			}

			result := svc.CurrentlyReading(tt.count)

			assert.NotNil(t, result)
			assert.Equal(t, len(tt.mockReviews), len(result))
			assert.Equal(t, tt.mockReviews, result)
			mockGR.AssertExpectations(t)
		})
	}
}

func TestService_Reviews_ContextPassed(t *testing.T) {
	mockGR := new(MockGoodReadsPort)
	count := 5

	mockGR.On("Reviews", mock.MatchedBy(func(ctx context.Context) bool {
		return ctx != nil
	}), count).Return([]responses.Review{}, nil)

	svc := New(mockGR)
	result := svc.Reviews(count)

	assert.NotNil(t, result)
	mockGR.AssertExpectations(t)
}

func TestService_CurrentlyReading_ContextPassed(t *testing.T) {
	mockGR := new(MockGoodReadsPort)
	count := 3

	mockGR.On("CurrentlyReading", mock.MatchedBy(func(ctx context.Context) bool {
		return ctx != nil
	}), count).Return([]responses.Review{}, nil)

	svc := New(mockGR)
	result := svc.CurrentlyReading(count)

	assert.NotNil(t, result)
	mockGR.AssertExpectations(t)
}

func TestService_Reviews_WithZeroCount(t *testing.T) {
	mockGR := new(MockGoodReadsPort)
	count := 0

	mockGR.On("Reviews", mock.Anything, count).
		Return([]responses.Review{}, nil)

	svc := New(mockGR)
	result := svc.Reviews(count)

	assert.NotNil(t, result)
	assert.Empty(t, result)
	mockGR.AssertExpectations(t)
}

func TestService_CurrentlyReading_WithZeroCount(t *testing.T) {
	mockGR := new(MockGoodReadsPort)
	count := 0

	mockGR.On("CurrentlyReading", mock.Anything, count).
		Return([]responses.Review{}, nil)

	svc := New(mockGR)
	result := svc.CurrentlyReading(count)

	assert.NotNil(t, result)
	assert.Empty(t, result)
	mockGR.AssertExpectations(t)
}

func TestService_Reviews_WithLargeCount(t *testing.T) {
	mockGR := new(MockGoodReadsPort)
	count := 100

	// Create a large slice of reviews
	reviews := make([]responses.Review, count)
	for i := 0; i < count; i++ {
		reviews[i] = responses.Review{
			Book: responses.AuthorBook{
				Title: "Book " + string(rune(i)),
				ID:    string(rune(i)),
			},
			Rating: i % 5,
		}
	}

	mockGR.On("Reviews", mock.Anything, count).
		Return(reviews, nil)

	svc := New(mockGR)
	result := svc.Reviews(count)

	assert.NotNil(t, result)
	assert.Len(t, result, count)
	mockGR.AssertExpectations(t)
}

func TestService_CurrentlyReading_WithLargeCount(t *testing.T) {
	mockGR := new(MockGoodReadsPort)
	count := 50

	// Create a large slice of currently reading reviews
	reviews := make([]responses.Review, count)
	for i := 0; i < count; i++ {
		reviews[i] = responses.Review{
			Book: responses.AuthorBook{
				Title: "Reading Book " + string(rune(i)),
				ID:    string(rune(i)),
			},
			Rating: 0,
		}
	}

	mockGR.On("CurrentlyReading", mock.Anything, count).
		Return(reviews, nil)

	svc := New(mockGR)
	result := svc.CurrentlyReading(count)

	assert.NotNil(t, result)
	assert.Len(t, result, count)
	mockGR.AssertExpectations(t)
}
