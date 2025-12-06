package githubadapter

import (
	"context"

	"github.com/shurcooL/githubv4"

	domain "hufschlaeger.net/markscribe/internal/domain"
)

// Adapter implements ports.GithubPort using the GitHub GraphQL v4 API.
type Adapter struct {
	client *githubv4.Client
}

func New(client *githubv4.Client) *Adapter { // constructor kept simple for now
	return &Adapter{client: client}
}

// GraphQL lightweight types local to the adapter
type qlRelease struct {
	Nodes []struct {
		Name         githubv4.String
		TagName      githubv4.String
		PublishedAt  githubv4.DateTime
		URL          githubv4.String
		IsPrerelease githubv4.Boolean
		IsDraft      githubv4.Boolean
	}
}

type qlRepository struct {
	NameWithOwner githubv4.String
	URL           githubv4.String
	Description   githubv4.String
	IsPrivate     githubv4.Boolean
	Stargazers    struct {
		TotalCount githubv4.Int
	}
	Releases qlRelease `graphql:"releases(last: 1)"`
}

type qlUser struct {
	Login     githubv4.String
	Name      githubv4.String
	AvatarURL githubv4.String
	URL       githubv4.String
}

type qlPullRequest struct {
	URL        githubv4.String
	Title      githubv4.String
	State      githubv4.PullRequestState
	CreatedAt  githubv4.DateTime
	Repository qlRepository
}

type qlGist struct {
	Name        githubv4.String
	Description githubv4.String
	URL         githubv4.String
	CreatedAt   githubv4.DateTime
}

// queries
type recentReposQuery struct {
	User struct {
		Login        githubv4.String
		Repositories struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   qlRepository
			}
		} `graphql:"repositories(first: $count, privacy: PUBLIC, isFork: $isFork, ownerAffiliations: OWNER, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

type repoQuery struct {
	Repository qlRepository `graphql:"repository(owner:$owner, name:$name)"`
}

type viewerQuery struct {
	Viewer struct {
		Login githubv4.String
	}
}

type followersQuery struct {
	User struct {
		Login     githubv4.String
		Followers struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   qlUser
			}
		} `graphql:"followers(first: $count)"`
	} `graphql:"user(login:$username)"`
}

type recentPullRequestsQuery struct {
	User struct {
		Login        githubv4.String
		PullRequests struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   qlPullRequest
			}
		} `graphql:"pullRequests(first: $count, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

type recentReleasesQuery struct {
	User struct {
		Login                     githubv4.String
		RepositoriesContributedTo struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   struct {
					qlRepository
					Releases qlRelease `graphql:"releases(first: 10, orderBy: {field: CREATED_AT, direction: DESC})"`
				}
			}
		} `graphql:"repositoriesContributedTo(first: 100, after:$after, includeUserRepositories: true, contributionTypes: COMMIT, privacy: PUBLIC)"`
	} `graphql:"user(login:$username)"`
}

type gistsQuery struct {
	User struct {
		Login githubv4.String
		Gists struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   qlGist
			}
		} `graphql:"gists(first: $count, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

type recentStarsQuery struct {
	User struct {
		Login githubv4.String
		Stars struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor    githubv4.String
				StarredAt githubv4.DateTime
				Node      qlRepository
			}
		} `graphql:"starredRepositories(first: $count, after:$after, orderBy: {field: STARRED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

type recentContributionsQuery struct {
	User struct {
		Login                   githubv4.String
		ContributionsCollection struct {
			CommitContributionsByRepository []struct {
				Contributions struct {
					Edges []struct {
						Cursor githubv4.String
						Node   struct {
							OccurredAt githubv4.DateTime
						}
					}
				} `graphql:"contributions(first: 1)"`
				Repository qlRepository
			} `graphql:"commitContributionsByRepository(maxRepositories: 100)"`
		}
	} `graphql:"user(login:$username)"`
}

type recentIssuesQuery struct {
	User struct {
		Login                   githubv4.String
		ContributionsCollection struct {
			IssueContributionsByRepository []struct {
				Contributions struct {
					Edges []struct {
						Cursor githubv4.String
						Node   struct {
							OccurredAt githubv4.DateTime
							Issue      struct {
								Title githubv4.String
							}
						}
					}
				} `graphql:"contributions(first: 1)"`
				Repository qlRepository
			} `graphql:"issueContributionsByRepository(maxRepositories: 100)"`
		}
	} `graphql:"user(login:$username)"`
}

type sponsorsQuery struct {
	User struct {
		Login                    githubv4.String
		SponsorshipsAsMaintainer struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   struct {
					CreatedAt     githubv4.DateTime
					SponsorEntity struct {
						Typename     githubv4.String `graphql:"__typename"`
						User         qlUser          `graphql:"... on User"`
						Organization qlUser          `graphql:"... on Organization"`
					}
				}
			}
		} `graphql:"sponsorshipsAsMaintainer(first: $count, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

// RecentRepos returns recent repositories for the given user.
func (a *Adapter) RecentRepos(ctx context.Context, username string, count int, isFork bool) ([]domain.Repo, error) {
	var q recentReposQuery
	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"count":    githubv4.Int(count),
		"isFork":   githubv4.Boolean(isFork),
	}
	if err := a.client.Query(ctx, &q, variables); err != nil {
		return nil, err
	}

	var repos []domain.Repo
	for _, edge := range q.User.Repositories.Edges {
		repos = append(repos, repoFromQL(edge.Node))
	}
	return repos, nil
}

// Repo returns a repository by owner/name.
func (a *Adapter) Repo(ctx context.Context, owner, name string) (domain.Repo, error) {
	var q repoQuery
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}
	if err := a.client.Query(ctx, &q, variables); err != nil {
		return domain.Repo{}, err
	}
	return repoFromQL(q.Repository), nil
}

// ViewerLogin returns the login of the authenticated viewer
func (a *Adapter) ViewerLogin(ctx context.Context) (string, error) {
	var q viewerQuery
	if err := a.client.Query(ctx, &q, nil); err != nil {
		return "", err
	}
	return string(q.Viewer.Login), nil
}

// Followers returns the followers for a user
func (a *Adapter) Followers(ctx context.Context, username string, count int) ([]domain.User, error) {
	var q followersQuery
	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"count":    githubv4.Int(count),
	}
	if err := a.client.Query(ctx, &q, variables); err != nil {
		return nil, err
	}
	var users []domain.User
	for _, edge := range q.User.Followers.Edges {
		users = append(users, userFromQL(edge.Node))
	}
	return users, nil
}

// RecentPullRequests returns recent pull requests created by the user.
func (a *Adapter) RecentPullRequests(ctx context.Context, username string, count int) ([]domain.PullRequest, error) {
	var q recentPullRequestsQuery
	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"count":    githubv4.Int(count + 1), // +1 to allow skipping meta-repo later
	}
	if err := a.client.Query(ctx, &q, variables); err != nil {
		return nil, err
	}
	var prs []domain.PullRequest
	for _, edge := range q.User.PullRequests.Edges {
		pr := edge.Node
		dpr := domain.PullRequest{
			Title:     string(pr.Title),
			URL:       string(pr.URL),
			State:     string(pr.State),
			CreatedAt: pr.CreatedAt.Time,
			Repo:      repoFromQL(pr.Repository),
		}
		prs = append(prs, dpr)
		if len(prs) >= count {
			break
		}
	}
	return prs, nil
}

// RecentReleases returns repositories with their latest non-draft, non-prerelease release.
func (a *Adapter) RecentReleases(ctx context.Context, username string, count int) ([]domain.Repo, error) {
	var after *githubv4.String
	var out []domain.Repo
	for {
		var q recentReleasesQuery
		variables := map[string]interface{}{
			"username": githubv4.String(username),
			"after":    after,
		}
		if err := a.client.Query(ctx, &q, variables); err != nil {
			return nil, err
		}
		if len(q.User.RepositoriesContributedTo.Edges) == 0 {
			break
		}
		for _, edge := range q.User.RepositoriesContributedTo.Edges {
			r := repoFromQL(edge.Node.qlRepository)
			// find first valid release (non-draft, non-prerelease) in descending order
			for _, rel := range edge.Node.Releases.Nodes {
				if bool(rel.IsDraft) || bool(rel.IsPrerelease) {
					continue
				}
				if rel.TagName == "" || rel.PublishedAt.IsZero() {
					continue
				}
				r.LastRelease = domain.Release{
					Name:        string(rel.Name),
					TagName:     string(rel.TagName),
					PublishedAt: rel.PublishedAt.Time,
					URL:         string(rel.URL),
				}
				break
			}
			if !r.LastRelease.PublishedAt.IsZero() {
				out = append(out, r)
				if len(out) >= count {
					continue
					// don't early return; still need to set cursor for next loop termination
				}
			}
			after = githubv4.NewString(edge.Cursor)
		}
		// If we've collected enough and there's no more pages or we prefer to stop, we can break.
		if len(out) >= count {
			break
		}
	}
	return out, nil
}

// RecentContributions returns commit contributions grouped by repository.
func (a *Adapter) RecentContributions(ctx context.Context, username string, count int) ([]domain.Contribution, error) {
	var q recentContributionsQuery
	variables := map[string]interface{}{
		"username": githubv4.String(username),
	}
	if err := a.client.Query(ctx, &q, variables); err != nil {
		return nil, err
	}
	var out []domain.Contribution
	for _, v := range q.User.ContributionsCollection.CommitContributionsByRepository {
		c := domain.Contribution{
			Repo:       repoFromQL(v.Repository),
			OccurredAt: v.Contributions.Edges[0].Node.OccurredAt.Time,
		}
		out = append(out, c)
	}
	// Limit at adapter level conservatively if oversize
	if len(out) > count {
		out = out[:count]
	}
	return out, nil
}

// RecentIssues returns recent issue contributions grouped by repository.
func (a *Adapter) RecentIssues(ctx context.Context, username string, count int) ([]domain.Issue, error) {
	var q recentIssuesQuery
	variables := map[string]interface{}{
		"username": githubv4.String(username),
	}
	if err := a.client.Query(ctx, &q, variables); err != nil {
		return nil, err
	}
	var out []domain.Issue
	for _, v := range q.User.ContributionsCollection.IssueContributionsByRepository {
		out = append(out, domain.Issue{
			Repo:       repoFromQL(v.Repository),
			OccurredAt: v.Contributions.Edges[0].Node.OccurredAt.Time,
			Title:      string(v.Contributions.Edges[0].Node.Issue.Title),
		})
		if len(out) >= count {
			break
		}
	}
	return out, nil
}

// Sponsors returns recent sponsors (users and organizations) for the maintainer.
func (a *Adapter) Sponsors(ctx context.Context, username string, count int) ([]domain.Sponsor, error) {
	var q sponsorsQuery
	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"count":    githubv4.Int(count),
	}
	if err := a.client.Query(ctx, &q, variables); err != nil {
		return nil, err
	}
	var out []domain.Sponsor
	for _, edge := range q.User.SponsorshipsAsMaintainer.Edges {
		se := edge.Node.SponsorEntity
		var u domain.User
		switch string(se.Typename) {
		case "User":
			u = userFromQL(se.User)
		case "Organization":
			u = userFromQL(se.Organization)
		default:
			continue
		}
		out = append(out, domain.Sponsor{User: u, CreatedAt: edge.Node.CreatedAt.Time})
		if len(out) >= count {
			break
		}
	}
	return out, nil
}

// Gists returns user's gists ordered by creation date desc limited by count.
func (a *Adapter) Gists(ctx context.Context, username string, count int) ([]domain.Gist, error) {
	var q gistsQuery
	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"count":    githubv4.Int(count),
	}
	if err := a.client.Query(ctx, &q, variables); err != nil {
		return nil, err
	}
	var out []domain.Gist
	for _, edge := range q.User.Gists.Edges {
		g := edge.Node
		out = append(out, domain.Gist{
			Name:        string(g.Name),
			Description: string(g.Description),
			URL:         string(g.URL),
			CreatedAt:   g.CreatedAt.Time,
		})
	}
	return out, nil
}

// RecentStars returns recently starred public repositories by the user.
func (a *Adapter) RecentStars(ctx context.Context, username string, count int) ([]domain.Star, error) {
	var q recentStarsQuery
	var after *githubv4.String
	var out []domain.Star
outer:
	for {
		variables := map[string]interface{}{
			"username": githubv4.String(username),
			"count":    githubv4.Int(count),
			"after":    after,
		}
		if err := a.client.Query(ctx, &q, variables); err != nil {
			return nil, err
		}
		if len(q.User.Stars.Edges) == 0 {
			break outer
		}
		for _, edge := range q.User.Stars.Edges {
			repo := edge.Node
			if bool(repo.IsPrivate) {
				continue
			}
			out = append(out, domain.Star{
				StarredAt: edge.StarredAt.Time,
				Repo:      repoFromQL(repo),
			})
			if len(out) == count {
				break outer
			}
			after = githubv4.NewString(edge.Cursor)
		}
	}
	return out, nil
}

// local helpers to map GraphQL to domain
func repoFromQL(repo qlRepository) domain.Repo {
	var lastRelease domain.Release
	if len(repo.Releases.Nodes) > 0 {
		r := repo.Releases.Nodes[len(repo.Releases.Nodes)-1]
		lastRelease = domain.Release{
			Name:        string(r.Name),
			TagName:     string(r.TagName),
			PublishedAt: r.PublishedAt.Time,
			URL:         string(r.URL),
		}
	}

	return domain.Repo{
		Name:        string(repo.NameWithOwner),
		URL:         string(repo.URL),
		Description: string(repo.Description),
		Stargazers:  int(repo.Stargazers.TotalCount),
		IsPrivate:   bool(repo.IsPrivate),
		LastRelease: lastRelease,
	}
}

func userFromQL(user qlUser) domain.User {
	return domain.User{
		Login:     string(user.Login),
		Name:      string(user.Name),
		AvatarURL: string(user.AvatarURL),
		URL:       string(user.URL),
	}
}
