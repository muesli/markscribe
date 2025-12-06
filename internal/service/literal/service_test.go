package literal

import (
	"context"
	"errors"
	"testing"

	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"hufschlaeger.net/markscribe/internal/adapters/literal"
)

// MockLiteralPort is a mock implementation of ports.LiteralPort
type MockLiteralPort struct {
	mock.Mock
}

func (m *MockLiteralPort) CurrentlyReading(ctx context.Context, count int) ([]literaladapter.Book, error) {
	args := m.Called(ctx, count)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]literaladapter.Book), args.Error(1)
}

func TestNew(t *testing.T) {
	mockLit := new(MockLiteralPort)

	svc := New(mockLit)

	assert.NotNil(t, svc)
	assert.Equal(t, mockLit, svc.lit)
}

func TestService_CurrentlyReading(t *testing.T) {
	tests := []struct {
		name          string
		count         int
		mockBooks     []literaladapter.Book
		mockError     error
		expectedPanic bool
	}{
		{
			name:  "successful retrieval with multiple books",
			count: 3,
			mockBooks: []literaladapter.Book{
				{
					Title: "Clean Code",
					Authors: []literaladapter.Author{
						{Name: "Robert C. Martin"},
					},
					Slug: "clean-code",
				},
				{
					Title: "The Pragmatic Programmer",
					Authors: []literaladapter.Author{
						{Name: "Andrew Hunt"},
						{Name: "David Thomas"},
					},
					Slug: "pragmatic-programmer",
				},
				{
					Title: "Design Patterns",
					Authors: []literaladapter.Author{
						{Name: "Erich Gamma"},
						{Name: "Richard Helm"},
						{Name: "Ralph Johnson"},
						{Name: "John Vlissides"},
					},
					Slug: "design-patterns",
				},
			},
			expectedPanic: false,
		},
		{
			name:          "successful retrieval with empty books list",
			count:         5,
			mockBooks:     []literaladapter.Book{},
			expectedPanic: false,
		},
		{
			name:  "successful retrieval with single book",
			count: 1,
			mockBooks: []literaladapter.Book{
				{
					Title:   "Test Driven Development",
					Authors: []literaladapter.Author{{Name: "Kent Beck"}},
					Slug:    "tdd",
				},
			},
			expectedPanic: false,
		},
		{
			name:  "successful retrieval with book without authors",
			count: 1,
			mockBooks: []literaladapter.Book{
				{
					Title:   "Anonymous Book",
					Authors: []literaladapter.Author{},
					Slug:    "anonymous",
				},
			},
			expectedPanic: false,
		},
		{
			name:  "successful retrieval with book with multiple authors",
			count: 1,
			mockBooks: []literaladapter.Book{
				{
					Title: "Collaborative Work",
					Authors: []literaladapter.Author{
						{Name: "Author One"}, {Name: "Author Two"}, {Name: "Author Three"}},
					Slug: "collaborative",
				},
			},
			expectedPanic: false,
		},
		{
			name:          "panics on error",
			count:         5,
			mockBooks:     nil,
			mockError:     errors.New("api error"),
			expectedPanic: true,
		},
		{
			name:          "panics on network error",
			count:         3,
			mockBooks:     nil,
			mockError:     errors.New("network timeout"),
			expectedPanic: true,
		},
		{
			name:          "panics on authentication error",
			count:         2,
			mockBooks:     nil,
			mockError:     errors.New("unauthorized"),
			expectedPanic: true,
		},
		{
			name:          "panics on rate limit error",
			count:         4,
			mockBooks:     nil,
			mockError:     errors.New("rate limit exceeded"),
			expectedPanic: true,
		},
		{
			name:          "panics on invalid response error",
			count:         1,
			mockBooks:     nil,
			mockError:     errors.New("invalid json response"),
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLit := new(MockLiteralPort)
			mockLit.On("CurrentlyReading", mock.Anything, tt.count).
				Return(tt.mockBooks, tt.mockError)

			svc := New(mockLit)

			if tt.expectedPanic {
				assert.Panics(t, func() {
					svc.CurrentlyReading(tt.count)
				}, "Expected CurrentlyReading to panic on error")
				mockLit.AssertExpectations(t)
				return
			}

			result := svc.CurrentlyReading(tt.count)

			assert.NotNil(t, result)
			assert.Equal(t, len(tt.mockBooks), len(result))
			assert.Equal(t, tt.mockBooks, result)
			mockLit.AssertExpectations(t)
		})
	}
}

func TestService_CurrentlyReading_ContextPassed(t *testing.T) {
	mockLit := new(MockLiteralPort)
	count := 5

	mockLit.On("CurrentlyReading", mock.MatchedBy(func(ctx context.Context) bool {
		return ctx != nil
	}), count).Return([]literaladapter.Book{}, nil)

	svc := New(mockLit)
	result := svc.CurrentlyReading(count)

	assert.NotNil(t, result)
	mockLit.AssertExpectations(t)
}

func TestService_CurrentlyReading_WithZeroCount(t *testing.T) {
	mockLit := new(MockLiteralPort)
	count := 0

	mockLit.On("CurrentlyReading", mock.Anything, count).
		Return([]literaladapter.Book{}, nil)

	svc := New(mockLit)
	result := svc.CurrentlyReading(count)

	assert.NotNil(t, result)
	assert.Empty(t, result)
	mockLit.AssertExpectations(t)
}

func TestService_CurrentlyReading_WithLargeCount(t *testing.T) {
	mockLit := new(MockLiteralPort)
	count := 50

	// Create a large slice of books
	books := make([]literaladapter.Book, count)
	for i := 0; i < count; i++ {
		books[i] = literaladapter.Book{
			Title: graphql.String("Book " + string(rune('A'+i))),
			Authors: []literaladapter.Author{
				{Name: graphql.String("Author " + string(rune('A'+i)))},
			},
			Slug: graphql.String("book-" + string(rune('a'+i))),
		}
	}

	mockLit.On("CurrentlyReading", mock.Anything, count).
		Return(books, nil)

	svc := New(mockLit)
	result := svc.CurrentlyReading(count)

	assert.NotNil(t, result)
	assert.Len(t, result, count)
	mockLit.AssertExpectations(t)
}

func TestService_CurrentlyReading_WithCompleteBookData(t *testing.T) {
	mockLit := new(MockLiteralPort)
	count := 2

	books := []literaladapter.Book{
		{
			Title:       "Domain-Driven Design",
			Subtitle:    "Tackling Complexity in the Heart of Software",
			Authors:     []literaladapter.Author{{Name: "Eric Evans"}},
			Description: "A comprehensive guide to domain-driven design",
			Slug:        "domain-driven-design",
		},
		{
			Title:    "Refactoring",
			Subtitle: "Improving the Design of Existing Code",
			Authors: []literaladapter.Author{
				{Name: "Martin Fowler"},
				{Name: "Kent Beck"},
				{Name: "John Brant"},
				{Name: "William Opdyke"},
				{Name: "Don Roberts"},
			},
			Description: "How to improve code without changing functionality",
			Slug:        "refactoring",
		},
	}

	mockLit.On("CurrentlyReading", mock.Anything, count).
		Return(books, nil)

	svc := New(mockLit)
	result := svc.CurrentlyReading(count)

	assert.NotNil(t, result)
	assert.Len(t, result, count)

	// Verify all fields are preserved
	for i, book := range result {
		assert.Equal(t, books[i].Title, book.Title)
		assert.Equal(t, books[i].Subtitle, book.Subtitle)
		assert.Equal(t, books[i].Authors, book.Authors)
		assert.Equal(t, books[i].Description, book.Description)
		assert.Equal(t, books[i].Slug, book.Slug)
	}

	mockLit.AssertExpectations(t)
}

func TestService_CurrentlyReading_PreservesOrder(t *testing.T) {
	mockLit := new(MockLiteralPort)
	count := 3

	books := []literaladapter.Book{
		{Title: "First Book", Authors: []literaladapter.Author{{Name: "Author A"}}, Slug: "first"},
		{Title: "Second Book", Authors: []literaladapter.Author{{Name: "Author B"}, {Name: "Author C"}}, Slug: "second"},
		{Title: "Third Book", Authors: []literaladapter.Author{{Name: "Author D"}}, Slug: "third"},
	}

	mockLit.On("CurrentlyReading", mock.Anything, count).
		Return(books, nil)

	svc := New(mockLit)
	result := svc.CurrentlyReading(count)

	assert.NotNil(t, result)
	assert.Len(t, result, count)

	// Verify order is preserved
	for i, book := range result {
		assert.Equal(t, books[i].Title, book.Title)
		assert.Equal(t, books[i].Authors, book.Authors)
		assert.Equal(t, books[i].Slug, book.Slug)
	}

	mockLit.AssertExpectations(t)
}

func TestService_CurrentlyReading_WithPartialBookData(t *testing.T) {
	mockLit := new(MockLiteralPort)
	count := 2

	books := []literaladapter.Book{
		{
			Title:   "Minimal Book",
			Authors: []literaladapter.Author{{Name: "Unknown"}},
			Slug:    "minimal",
			// Other fields empty/zero values
		},
		{
			Title:   "Another Book",
			Authors: []literaladapter.Author{{}},
			Slug:    "another",
		},
	}

	mockLit.On("CurrentlyReading", mock.Anything, count).
		Return(books, nil)

	svc := New(mockLit)
	result := svc.CurrentlyReading(count)

	assert.NotNil(t, result)
	assert.Len(t, result, count)
	assert.Equal(t, books, result)
	mockLit.AssertExpectations(t)
}

func TestService_CurrentlyReading_NegativeCount(t *testing.T) {
	mockLit := new(MockLiteralPort)
	count := -1

	mockLit.On("CurrentlyReading", mock.Anything, count).
		Return([]literaladapter.Book{}, nil)

	svc := New(mockLit)
	result := svc.CurrentlyReading(count)

	assert.NotNil(t, result)
	assert.Empty(t, result)
	mockLit.AssertExpectations(t)
}

func TestService_CurrentlyReading_MultipleCallsSameService(t *testing.T) {
	mockLit := new(MockLiteralPort)
	svc := New(mockLit)

	// First call
	books1 := []literaladapter.Book{
		{Title: "Book 1", Authors: []literaladapter.Author{{Name: "Author 1"}}, Slug: "book-1"},
	}
	mockLit.On("CurrentlyReading", mock.Anything, 1).
		Return(books1, nil).Once()

	result1 := svc.CurrentlyReading(1)
	assert.Len(t, result1, 1)

	// Second call with different count
	books2 := []literaladapter.Book{
		{Title: "Book 2", Authors: []literaladapter.Author{{Name: "Author 2"}, {Name: "Author 3"}}, Slug: "book-2"},
		{Title: "Book 3", Authors: []literaladapter.Author{{Name: "Author 4"}}, Slug: "book-3"},
	}
	mockLit.On("CurrentlyReading", mock.Anything, 2).
		Return(books2, nil).Once()

	result2 := svc.CurrentlyReading(2)
	assert.Len(t, result2, 2)

	mockLit.AssertExpectations(t)
}

func TestService_CurrentlyReading_WithNilAuthors(t *testing.T) {
	mockLit := new(MockLiteralPort)
	count := 1

	books := []literaladapter.Book{
		{
			Title:   "Book Without Authors",
			Authors: nil,
			Slug:    "no-authors",
		},
	}

	mockLit.On("CurrentlyReading", mock.Anything, count).
		Return(books, nil)

	svc := New(mockLit)
	result := svc.CurrentlyReading(count)

	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Nil(t, result[0].Authors)
	mockLit.AssertExpectations(t)
}
