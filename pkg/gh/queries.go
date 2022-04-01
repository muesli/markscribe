package gh

import (
	"github.com/muesli/markscribe/pkg/types"

	"github.com/shurcooL/githubv4"
)

var gistsQuery struct {
	User struct {
		Login githubv4.String
		Gists struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   types.QLGist
			}
		} `graphql:"gists(first: $count, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

var viewerQuery struct {
	Viewer struct {
		Login githubv4.String
	}
}

var recentFollowersQuery struct {
	User struct {
		Login     githubv4.String
		Followers struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   types.QLUser
			}
		} `graphql:"followers(first: $count)"`
	} `graphql:"user(login:$username)"`
}

var recentContributionsQuery struct {
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
				Repository types.QLRepository
			} `graphql:"commitContributionsByRepository(maxRepositories: 100)"`
		}
	} `graphql:"user(login:$username)"`
}

var recentPullRequestsQuery struct {
	User struct {
		Login        githubv4.String
		PullRequests struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   types.QLPullRequest
			}
		} `graphql:"pullRequests(first: $count, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

var recentReposQuery struct {
	User struct {
		Login        githubv4.String
		Repositories struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   types.QLRepository
			}
		} `graphql:"repositories(first: $count, privacy: PUBLIC, isFork: false, ownerAffiliations: OWNER, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

var recentReleasesQuery struct {
	User struct {
		Login                     githubv4.String
		RepositoriesContributedTo struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   struct {
					types.QLRepository
					Releases types.QLRelease `graphql:"releases(first: 10, orderBy: {field: CREATED_AT, direction: DESC})"`
				}
			}
		} `graphql:"repositoriesContributedTo(first: 100, after:$after includeUserRepositories: true, contributionTypes: COMMIT, privacy: PUBLIC)"`
	} `graphql:"user(login:$username)"`
}

var sponsorsQuery struct {
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
						User         types.QLUser    `graphql:"... on User"`
						Organization types.QLUser    `graphql:"... on Organization"`
					}
				}
			}
		} `graphql:"sponsorshipsAsMaintainer(first: $count, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

var recentStarsQuery struct {
	User struct {
		Login githubv4.String
		Stars struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor    githubv4.String
				StarredAt githubv4.DateTime
				Node      types.QLRepository
			}
		} `graphql:"starredRepositories(first: $count, orderBy: {field: STARRED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}
