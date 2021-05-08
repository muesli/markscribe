package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/KyleBanks/goodreads"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var (
	gitHubClient    *githubv4.Client
	goodReadsClient *goodreads.Client
	goodReadsID     string
	username        string

	write = flag.String("write", "", "write output to")
)

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: markscribe [template]")
		os.Exit(1)
	}

	tplIn, err := ioutil.ReadFile(flag.Args()[0])
	if err != nil {
		fmt.Println("Can't read file:", err)
		os.Exit(1)
	}

	tpl, err := template.New("tpl").Funcs(template.FuncMap{
		/* GitHub */
		"recentContributions": recentContributions,
		"recentPullRequests":  recentPullRequests,
		"recentRepos":         recentRepos,
		"recentReleases":      recentReleases,
		"followers":           recentFollowers,
		"gists":               gists,
		"sponsors":            sponsors,
		/* RSS */
		"rss": rssFeed,
		/* GoodReads */
		"goodReadsReviews":          goodReadsReviews,
		"goodReadsCurrentlyReading": goodReadsCurrentlyReading,
		/* Utils */
		"humanize": humanized,
		"reverse":  reverse,
		"now":      time.Now,
	}).Parse(string(tplIn))
	if err != nil {
		fmt.Println("Can't parse template:", err)
		os.Exit(1)
	}

	var httpClient *http.Client
	gitHubToken := os.Getenv("GITHUB_TOKEN")
	goodReadsToken := os.Getenv("GOODREADS_TOKEN")
	goodReadsID = os.Getenv("GOODREADS_USER_ID")
	if len(gitHubToken) > 0 {
		httpClient = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: gitHubToken},
		))
	}

	gitHubClient = githubv4.NewClient(httpClient)
	goodReadsClient = goodreads.NewClient(goodReadsToken)

	if len(gitHubToken) > 0 {
		username, err = getUsername()
		if err != nil {
			fmt.Println("Can't retrieve GitHub profile:", err)
			os.Exit(1)
		}
	}

	w := os.Stdout
	if len(*write) > 0 {
		f, err := os.Create(*write)
		if err != nil {
			fmt.Println("Can't create:", err)
			os.Exit(1)
		}
		defer f.Close()
		w = f
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Can't render template:", err)
		os.Exit(1)
	}
}
