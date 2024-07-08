package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/KyleBanks/goodreads"
	"github.com/go-sprout/sprout"
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

	funcMap := sprout.FuncMap(sprout.WithAlias("lower", "toLower"))
	/* Github */
	funcMap["recentContributions"] = recentContributions
	funcMap["recentPullRequests"] = recentPullRequests
	funcMap["recentRepos"] = recentRepos
	funcMap["recentForks"] = recentForks
	funcMap["recentReleases"] = recentReleases
	funcMap["followers"] = recentFollowers
	funcMap["recentStars"] = recentStars
	funcMap["gists"] = gists
	funcMap["sponsors"] = sponsors
	funcMap["repo"] = repo
	/* RSS */
	funcMap["rss"] = rssFeed
	/* GoodReads */
	funcMap["goodReadsReviews"] = goodReadsReviews
	funcMap["goodReadsCurrentlyReading"] = goodReadsCurrentlyReading
	/* Literal.club */
	funcMap["literalClubCurrentlyReading"] = literalClubCurrentlyReading
	/* Utils */
	funcMap["humanize"] = humanized

	tpl, err := template.New("tpl").Funcs(funcMap).Parse(string(tplIn))
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
		defer f.Close() //nolint: errcheck
		w = f
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Can't render template:", err)
		os.Exit(1)
	}
}
