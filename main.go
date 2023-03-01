package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/KyleBanks/goodreads"
	"github.com/shkh/lastfm-go/lastfm"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var (
	gitHubClient    *githubv4.Client
	goodReadsClient *goodreads.Client
	lastfmClient    *lastfm.Api

	goodReadsID  string
	username     string
	lastfmUser   string
	lastfmAPIKey string
	lastfmSecret string

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
		"recentForks":         recentForks,
		"recentReleases":      recentReleases,
		"followers":           recentFollowers,
		"recentStars":         recentStars,
		"gists":               gists,
		"sponsors":            sponsors,
		"repo":                repo,
		/* RSS */
		"rss": rssFeed,
		/* GoodReads */
		"goodReadsReviews":          goodReadsReviews,
		"goodReadsCurrentlyReading": goodReadsCurrentlyReading,
		/* Literal.club */
		"literalClubCurrentlyReading": literalClubCurrentlyReading,
		/* last.fm */
		"lastfmFavouriteAlbums":  lastfmFavouriteAlbums,
		"lastfmFavouriteTracks":  lastfmFavouriteTracks,
		"lastfmFavouriteArtists": lastfmFavouriteArtists,
		"lastfmRecentTracks":     lastfmRecentTracks,
		/* Utils */
		"humanize": humanized,
		"reverse":  reverse,
		"now":      time.Now,
		"contains": strings.Contains,
		"toLower":  strings.ToLower,
	}).Parse(string(tplIn))
	if err != nil {
		fmt.Println("Can't parse template:", err)
		os.Exit(1)
	}

	gitHubToken := os.Getenv("GITHUB_TOKEN")
	goodReadsToken := os.Getenv("GOODREADS_TOKEN")
	goodReadsID = os.Getenv("GOODREADS_USER_ID")
	lastfmUser = os.Getenv("LASTFM_USER")
	lastfmAPIKey = os.Getenv("LASTFM_API_KEY")
	lastfmSecret = os.Getenv("LASTFM_API_SECRET")

	var httpClient *http.Client
	if len(gitHubToken) > 0 {
		httpClient = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: gitHubToken},
		))
	}

	gitHubClient = githubv4.NewClient(httpClient)
	goodReadsClient = goodreads.NewClient(goodReadsToken)
	lastfmClient = lastfm.New(lastfmAPIKey, lastfmSecret)

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
