package internal

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var (
	client   *githubv4.Client
	username string

	write = flag.String("write", "", "write output to")
)

func New(tplIn []byte) error {
	tpl, err := template.New("tpl").Funcs(template.FuncMap{
		/* GitHub */
		"recentContributions": recentContributions,
		"recentRepos":         recentRepos,
		"recentReleases":      recentReleases,
		"followers":           recentFollowers,
		"gists":               gists,
		"sponsors":            sponsors,
		/* RSS */
		"rss": rssFeed,
		/* Utils */
		"humanize": humanized,
	}).Parse(string(tplIn))

	if err != nil {
		return fmt.Errorf("Can't parse template: %s", err)
	}

	var httpClient *http.Client
	token := os.Getenv("GITHUB_TOKEN")
	if len(token) > 0 {
		httpClient = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		))
	}

	client = githubv4.NewClient(httpClient)

	if len(token) > 0 {
		username, err = getUsername()
		if err != nil {
			return fmt.Errorf("Can't retrieve GitHub profile: %s", err)
		}
	}

	w := os.Stdout
	if len(*write) > 0 {
		f, err := os.Create(*write)
		if err != nil {
			return fmt.Errorf("Can't create: %s", err)
		}
		defer f.Close()
		w = f
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		return fmt.Errorf("Can't render template: %s", err)
	}

	return nil
}
