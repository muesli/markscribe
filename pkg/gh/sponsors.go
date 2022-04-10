package gh

import (
	"context"

	"github.com/muesli/markscribe/pkg/types"

	"github.com/shurcooL/githubv4"
)

// GetSponsors retrieves 'n' sponsors.
func (c *Client) GetSponsors(n int) []types.Sponsor {
	sponsors := []types.Sponsor{}
	variables := Map{
		"count":    githubv4.Int(n),
		"username": c.username,
	}

	err := c.client.Query(context.Background(), &sponsorsQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range sponsorsQuery.User.SponsorshipsAsMaintainer.Edges {
		switch v.Node.SponsorEntity.Typename {
		case "User":
			sponsors = append(sponsors, types.Sponsor{
				User:      types.UserFromQL(v.Node.SponsorEntity.User),
				CreatedAt: v.Node.CreatedAt.Time,
			})
		case "Organization":
			sponsors = append(sponsors, types.Sponsor{
				User:      types.UserFromQL(v.Node.SponsorEntity.Organization),
				CreatedAt: v.Node.CreatedAt.Time,
			})
		}
	}

	return sponsors
}
