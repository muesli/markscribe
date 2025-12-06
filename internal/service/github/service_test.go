package github

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"hufschlaeger.net/markscribe/internal/domain"
)

// MockGithubPort is a mock implementation of ports.GithubPort
type MockGithubPort struct {
	mock.Mock
}

func (m *MockGithubPort) RecentRepos(ctx context.Context, username string, count int, forks bool) ([]domain.Repo, error) {
	args := m.Called(ctx, username, count, forks)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Repo), args.Error(1)
}

func (m *MockGithubPort) Repo(ctx context.Context, owner, name string) (domain.Repo, error) {
	args := m.Called(ctx, owner, name)
	return args.Get(0).(domain.Repo), args.Error(1)
}

func (m *MockGithubPort) Followers(ctx context.Context, username string, count int) ([]domain.User, error) {
	args := m.Called(ctx, username, count)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockGithubPort) RecentPullRequests(ctx context.Context, username string, count int) ([]domain.PullRequest, error) {
	args := m.Called(ctx, username, count)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PullRequest), args.Error(1)
}

func (m *MockGithubPort) RecentReleases(ctx context.Context, username string, count int) ([]domain.Repo, error) {
	args := m.Called(ctx, username, count)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Repo), args.Error(1)
}

func (m *MockGithubPort) RecentContributions(ctx context.Context, username string, count int) ([]domain.Contribution, error) {
	args := m.Called(ctx, username, count)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Contribution), args.Error(1)
}

func (m *MockGithubPort) Gists(ctx context.Context, username string, count int) ([]domain.Gist, error) {
	args := m.Called(ctx, username, count)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Gist), args.Error(1)
}

func (m *MockGithubPort) RecentStars(ctx context.Context, username string, count int) ([]domain.Star, error) {
	args := m.Called(ctx, username, count)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Star), args.Error(1)
}

func (m *MockGithubPort) RecentIssues(ctx context.Context, username string, count int) ([]domain.Issue, error) {
	args := m.Called(ctx, username, count)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Issue), args.Error(1)
}

func (m *MockGithubPort) Sponsors(ctx context.Context, username string, count int) ([]domain.Sponsor, error) {
	args := m.Called(ctx, username, count)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Sponsor), args.Error(1)
}

func (m *MockGithubPort) ViewerLogin(ctx context.Context) (string, error) {
	return "", nil
}

func TestNew(t *testing.T) {
	mockGH := new(MockGithubPort)
	username := "testuser"

	svc := New(mockGH, username)

	assert.NotNil(t, svc)
	assert.Equal(t, username, svc.username)
	assert.Equal(t, mockGH, svc.gh)
}

func TestService_RecentRepos(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		count          int
		mockRepos      []domain.Repo
		mockError      error
		expectedLen    int
		expectedPanic  bool
		expectedResult []domain.Repo
	}{
		{
			name:     "successful retrieval with meta repo filtered",
			username: "testuser",
			count:    2,
			mockRepos: []domain.Repo{
				{Name: "testuser/repo1"},
				{Name: "testuser/testuser"}, // meta repo - should be filtered
				{Name: "testuser/repo2"},
			},
			expectedLen: 2,
			expectedResult: []domain.Repo{
				{Name: "testuser/repo1"},
				{Name: "testuser/repo2"},
			},
		},
		{
			name:     "returns exact count",
			username: "testuser",
			count:    1,
			mockRepos: []domain.Repo{
				{Name: "testuser/repo1"},
				{Name: "testuser/repo2"},
			},
			expectedLen: 1,
			expectedResult: []domain.Repo{
				{Name: "testuser/repo1"},
			},
		},
		{
			name:          "panics on error",
			username:      "testuser",
			count:         2,
			mockError:     errors.New("api error"),
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGH := new(MockGithubPort)
			mockGH.On("RecentRepos", mock.Anything, tt.username, tt.count+1, false).
				Return(tt.mockRepos, tt.mockError)

			svc := New(mockGH, tt.username)

			if tt.expectedPanic {
				assert.Panics(t, func() {
					svc.RecentRepos(tt.count)
				})
				return
			}

			result := svc.RecentRepos(tt.count)

			assert.Len(t, result, tt.expectedLen)
			assert.Equal(t, tt.expectedResult, result)
			mockGH.AssertExpectations(t)
		})
	}
}

func TestService_RecentForks(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		count          int
		mockRepos      []domain.Repo
		mockError      error
		expectedLen    int
		expectedPanic  bool
		expectedResult []domain.Repo
	}{
		{
			name:     "successful retrieval with meta repo filtered",
			username: "testuser",
			count:    2,
			mockRepos: []domain.Repo{
				{Name: "testuser/fork1"},
				{Name: "testuser/testuser"}, // meta repo - should be filtered
				{Name: "testuser/fork2"},
			},
			expectedLen: 2,
			expectedResult: []domain.Repo{
				{Name: "testuser/fork1"},
				{Name: "testuser/fork2"},
			},
		},
		{
			name:          "panics on error",
			username:      "testuser",
			count:         2,
			mockError:     errors.New("api error"),
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGH := new(MockGithubPort)
			mockGH.On("RecentRepos", mock.Anything, tt.username, tt.count+1, true).
				Return(tt.mockRepos, tt.mockError)

			svc := New(mockGH, tt.username)

			if tt.expectedPanic {
				assert.Panics(t, func() {
					svc.RecentForks(tt.count)
				})
				return
			}

			result := svc.RecentForks(tt.count)

			assert.Len(t, result, tt.expectedLen)
			assert.Equal(t, tt.expectedResult, result)
			mockGH.AssertExpectations(t)
		})
	}
}

func TestService_Repo(t *testing.T) {
	tests := []struct {
		name          string
		owner         string
		repoName      string
		mockRepo      domain.Repo
		mockError     error
		expectedPanic bool
	}{
		{
			name:     "successful retrieval",
			owner:    "testowner",
			repoName: "testrepo",
			mockRepo: domain.Repo{Name: "testowner/testrepo"},
		},
		{
			name:          "panics on error",
			owner:         "testowner",
			repoName:      "testrepo",
			mockError:     errors.New("api error"),
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGH := new(MockGithubPort)
			mockGH.On("Repo", mock.Anything, tt.owner, tt.repoName).
				Return(tt.mockRepo, tt.mockError)

			svc := New(mockGH, "testuser")

			if tt.expectedPanic {
				assert.Panics(t, func() {
					svc.Repo(tt.owner, tt.repoName)
				})
				return
			}

			result := svc.Repo(tt.owner, tt.repoName)

			assert.Equal(t, tt.mockRepo, result)
			mockGH.AssertExpectations(t)
		})
	}
}

func TestService_Followers(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		count         int
		mockUsers     []domain.User
		mockError     error
		expectedPanic bool
	}{
		{
			name:     "successful retrieval",
			username: "testuser",
			count:    2,
			mockUsers: []domain.User{
				{Login: "follower1"},
				{Login: "follower2"},
			},
		},
		{
			name:          "panics on error",
			username:      "testuser",
			count:         2,
			mockError:     errors.New("api error"),
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGH := new(MockGithubPort)
			mockGH.On("Followers", mock.Anything, tt.username, tt.count).
				Return(tt.mockUsers, tt.mockError)

			svc := New(mockGH, tt.username)

			if tt.expectedPanic {
				assert.Panics(t, func() {
					svc.Followers(tt.count)
				})
				return
			}

			result := svc.Followers(tt.count)

			assert.Equal(t, tt.mockUsers, result)
			mockGH.AssertExpectations(t)
		})
	}
}

func TestService_RecentPullRequests(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		count          int
		mockPRs        []domain.PullRequest
		mockError      error
		expectedLen    int
		expectedPanic  bool
		expectedResult []domain.PullRequest
	}{
		{
			name:     "filters meta repo and private repos",
			username: "testuser",
			count:    2,
			mockPRs: []domain.PullRequest{
				{Repo: domain.Repo{Name: "testuser/repo1", IsPrivate: false}},
				{Repo: domain.Repo{Name: "testuser/testuser", IsPrivate: false}}, // meta - filtered
				{Repo: domain.Repo{Name: "testuser/repo2", IsPrivate: true}},     // private - filtered
				{Repo: domain.Repo{Name: "testuser/repo3", IsPrivate: false}},
			},
			expectedLen: 2,
			expectedResult: []domain.PullRequest{
				{Repo: domain.Repo{Name: "testuser/repo1", IsPrivate: false}},
				{Repo: domain.Repo{Name: "testuser/repo3", IsPrivate: false}},
			},
		},
		{
			name:          "panics on error",
			username:      "testuser",
			count:         2,
			mockError:     errors.New("api error"),
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGH := new(MockGithubPort)
			mockGH.On("RecentPullRequests", mock.Anything, tt.username, tt.count+1).
				Return(tt.mockPRs, tt.mockError)

			svc := New(mockGH, tt.username)

			if tt.expectedPanic {
				assert.Panics(t, func() {
					svc.RecentPullRequests(tt.count)
				})
				return
			}

			result := svc.RecentPullRequests(tt.count)

			assert.Len(t, result, tt.expectedLen)
			assert.Equal(t, tt.expectedResult, result)
			mockGH.AssertExpectations(t)
		})
	}
}

func TestService_RecentReleases(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-1 * time.Hour)

	tests := []struct {
		name           string
		username       string
		count          int
		mockRepos      []domain.Repo
		mockError      error
		expectedLen    int
		expectedPanic  bool
		expectedResult []domain.Repo
	}{
		{
			name:     "sorts by PublishedAt desc, then Stargazers",
			username: "testuser",
			count:    3,
			mockRepos: []domain.Repo{
				{Name: "repo1", LastRelease: domain.Release{PublishedAt: earlier}, Stargazers: 10},
				{Name: "repo2", LastRelease: domain.Release{PublishedAt: now}, Stargazers: 5},
				{Name: "repo3", LastRelease: domain.Release{PublishedAt: now}, Stargazers: 15},
			},
			expectedLen: 3,
			expectedResult: []domain.Repo{
				{Name: "repo3", LastRelease: domain.Release{PublishedAt: now}, Stargazers: 15},
				{Name: "repo2", LastRelease: domain.Release{PublishedAt: now}, Stargazers: 5},
				{Name: "repo1", LastRelease: domain.Release{PublishedAt: earlier}, Stargazers: 10},
			},
		},
		{
			name:     "limits to count",
			username: "testuser",
			count:    2,
			mockRepos: []domain.Repo{
				{Name: "repo1", LastRelease: domain.Release{PublishedAt: now}},
				{Name: "repo2", LastRelease: domain.Release{PublishedAt: now.Add(-1 * time.Hour)}},
				{Name: "repo3", LastRelease: domain.Release{PublishedAt: now.Add(-2 * time.Hour)}},
			},
			expectedLen: 2,
		},
		{
			name:          "panics on error",
			username:      "testuser",
			count:         2,
			mockError:     errors.New("api error"),
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGH := new(MockGithubPort)
			mockGH.On("RecentReleases", mock.Anything, tt.username, tt.count).
				Return(tt.mockRepos, tt.mockError)

			svc := New(mockGH, tt.username)

			if tt.expectedPanic {
				assert.Panics(t, func() {
					svc.RecentReleases(tt.count)
				})
				return
			}

			result := svc.RecentReleases(tt.count)

			assert.Len(t, result, tt.expectedLen)
			if tt.expectedResult != nil {
				assert.Equal(t, tt.expectedResult, result)
			}
			mockGH.AssertExpectations(t)
		})
	}
}

func TestService_RecentContributions(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-1 * time.Hour)

	tests := []struct {
		name              string
		username          string
		count             int
		mockContributions []domain.Contribution
		mockError         error
		expectedLen       int
		expectedPanic     bool
		expectedResult    []domain.Contribution
	}{
		{
			name:     "filters meta repo and private repos, sorts by time",
			username: "testuser",
			count:    2,
			mockContributions: []domain.Contribution{
				{Repo: domain.Repo{Name: "testuser/repo1", IsPrivate: false}, OccurredAt: earlier},
				{Repo: domain.Repo{Name: "testuser/testuser", IsPrivate: false}, OccurredAt: now}, // meta - filtered
				{Repo: domain.Repo{Name: "testuser/repo2", IsPrivate: true}, OccurredAt: now},     // private - filtered
				{Repo: domain.Repo{Name: "testuser/repo3", IsPrivate: false}, OccurredAt: now},
			},
			expectedLen: 2,
			expectedResult: []domain.Contribution{
				{Repo: domain.Repo{Name: "testuser/repo3", IsPrivate: false}, OccurredAt: now},
				{Repo: domain.Repo{Name: "testuser/repo1", IsPrivate: false}, OccurredAt: earlier},
			},
		},
		{
			name:          "panics on error",
			username:      "testuser",
			count:         2,
			mockError:     errors.New("api error"),
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGH := new(MockGithubPort)
			mockGH.On("RecentContributions", mock.Anything, tt.username, tt.count+10).
				Return(tt.mockContributions, tt.mockError)

			svc := New(mockGH, tt.username)

			if tt.expectedPanic {
				assert.Panics(t, func() {
					svc.RecentContributions(tt.count)
				})
				return
			}

			result := svc.RecentContributions(tt.count)

			assert.Len(t, result, tt.expectedLen)
			if tt.expectedResult != nil {
				assert.Equal(t, tt.expectedResult, result)
			}
			mockGH.AssertExpectations(t)
		})
	}
}

func TestService_Gists(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		count         int
		mockGists     []domain.Gist
		mockError     error
		expectedPanic bool
	}{
		{
			name:     "successful retrieval",
			username: "testuser",
			count:    2,
			mockGists: []domain.Gist{
				{Description: "gist1"},
				{Description: "gist2"},
			},
		},
		{
			name:          "panics on error",
			username:      "testuser",
			count:         2,
			mockError:     errors.New("api error"),
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGH := new(MockGithubPort)
			mockGH.On("Gists", mock.Anything, tt.username, tt.count).
				Return(tt.mockGists, tt.mockError)

			svc := New(mockGH, tt.username)

			if tt.expectedPanic {
				assert.Panics(t, func() {
					svc.Gists(tt.count)
				})
				return
			}

			result := svc.Gists(tt.count)

			assert.Equal(t, tt.mockGists, result)
			mockGH.AssertExpectations(t)
		})
	}
}

func TestService_RecentStars(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		count         int
		mockStars     []domain.Star
		mockError     error
		expectedPanic bool
	}{
		{
			name:     "successful retrieval",
			username: "testuser",
			count:    2,
			mockStars: []domain.Star{
				{Repo: domain.Repo{Name: "repo1"}},
				{Repo: domain.Repo{Name: "repo2"}},
			},
		},
		{
			name:          "panics on error",
			username:      "testuser",
			count:         2,
			mockError:     errors.New("api error"),
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGH := new(MockGithubPort)
			mockGH.On("RecentStars", mock.Anything, tt.username, tt.count).
				Return(tt.mockStars, tt.mockError)

			svc := New(mockGH, tt.username)

			if tt.expectedPanic {
				assert.Panics(t, func() {
					svc.RecentStars(tt.count)
				})
				return
			}

			result := svc.RecentStars(tt.count)

			assert.Equal(t, tt.mockStars, result)
			mockGH.AssertExpectations(t)
		})
	}
}

func TestService_RecentIssues(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-1 * time.Hour)

	tests := []struct {
		name           string
		username       string
		count          int
		mockIssues     []domain.Issue
		mockError      error
		expectedLen    int
		expectedPanic  bool
		expectedResult []domain.Issue
	}{
		{
			name:     "filters meta repo and private repos, sorts by time",
			username: "testuser",
			count:    2,
			mockIssues: []domain.Issue{
				{Repo: domain.Repo{Name: "testuser/repo1", IsPrivate: false}, OccurredAt: earlier},
				{Repo: domain.Repo{Name: "testuser/testuser", IsPrivate: false}, OccurredAt: now}, // meta - filtered
				{Repo: domain.Repo{Name: "testuser/repo2", IsPrivate: true}, OccurredAt: now},     // private - filtered
				{Repo: domain.Repo{Name: "testuser/repo3", IsPrivate: false}, OccurredAt: now},
			},
			expectedLen: 2,
			expectedResult: []domain.Issue{
				{Repo: domain.Repo{Name: "testuser/repo3", IsPrivate: false}, OccurredAt: now},
				{Repo: domain.Repo{Name: "testuser/repo1", IsPrivate: false}, OccurredAt: earlier},
			},
		},
		{
			name:          "panics on error",
			username:      "testuser",
			count:         2,
			mockError:     errors.New("api error"),
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGH := new(MockGithubPort)
			mockGH.On("RecentIssues", mock.Anything, tt.username, tt.count+10).
				Return(tt.mockIssues, tt.mockError)

			svc := New(mockGH, tt.username)

			if tt.expectedPanic {
				assert.Panics(t, func() {
					svc.RecentIssues(tt.count)
				})
				return
			}

			result := svc.RecentIssues(tt.count)

			assert.Len(t, result, tt.expectedLen)
			if tt.expectedResult != nil {
				assert.Equal(t, tt.expectedResult, result)
			}
			mockGH.AssertExpectations(t)
		})
	}
}

func TestService_Sponsors(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		count          int
		mockSponsors   []domain.Sponsor
		mockError      error
		expectedLen    int
		expectedPanic  bool
		expectedResult []domain.Sponsor
	}{
		{
			name:     "successful retrieval",
			username: "testuser",
			count:    2,
			mockSponsors: []domain.Sponsor{
				{User: domain.User{Login: "sponsor1"}},
				{User: domain.User{Login: "sponsor2"}},
			},
			expectedLen: 2,
			expectedResult: []domain.Sponsor{
				{User: domain.User{Login: "sponsor1"}},
				{User: domain.User{Login: "sponsor2"}},
			},
		},
		{
			name:     "limits to count when more returned",
			username: "testuser",
			count:    2,
			mockSponsors: []domain.Sponsor{
				{User: domain.User{Login: "sponsor1"}},
				{User: domain.User{Login: "sponsor2"}},
				{User: domain.User{Login: "sponsor3"}},
			},
			expectedLen: 2,
			expectedResult: []domain.Sponsor{
				{User: domain.User{Login: "sponsor1"}},
				{User: domain.User{Login: "sponsor2"}},
			},
		},
		{
			name:          "panics on error",
			username:      "testuser",
			count:         2,
			mockError:     errors.New("api error"),
			expectedPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGH := new(MockGithubPort)
			mockGH.On("Sponsors", mock.Anything, tt.username, tt.count).
				Return(tt.mockSponsors, tt.mockError)

			svc := New(mockGH, tt.username)

			if tt.expectedPanic {
				assert.Panics(t, func() {
					svc.Sponsors(tt.count)
				})
				return
			}

			result := svc.Sponsors(tt.count)

			assert.Len(t, result, tt.expectedLen)
			assert.Equal(t, tt.expectedResult, result)
			mockGH.AssertExpectations(t)
		})
	}
}
