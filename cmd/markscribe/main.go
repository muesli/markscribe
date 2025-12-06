package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	texttmpl "text/template"
	"time"

	kbgoodreads "github.com/KyleBanks/goodreads"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	githubsvc "hufschlaeger.net/markscribe/internal/service/github"
	goodreadssvc "hufschlaeger.net/markscribe/internal/service/goodreads"
	literalsvc "hufschlaeger.net/markscribe/internal/service/literal"
	rsssvc "hufschlaeger.net/markscribe/internal/service/rss"
	templatesvc "hufschlaeger.net/markscribe/internal/service/template"

	githubadapter "hufschlaeger.net/markscribe/internal/adapters/github"
	goodreadsadapter "hufschlaeger.net/markscribe/internal/adapters/goodreads"
	literaladapter "hufschlaeger.net/markscribe/internal/adapters/literal"
	rssadapter "hufschlaeger.net/markscribe/internal/adapters/rss"
	"hufschlaeger.net/markscribe/internal/usecase/ports"
)

var (
	gitHubClient    *githubv4.Client
	goodReadsClient *kbgoodreads.Client
	goodReadsID     string
	username        string
	gh              ports.GithubPort
	gr              ports.GoodReadsPort
	lit             ports.LiteralPort
	rss             ports.RssFeedPort

	write = flag.String("write", "", "write output to")
)

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: markscribe [template]")
		os.Exit(1)
	}

	tplIn, err := os.ReadFile(flag.Args()[0])
	if err != nil {
		fmt.Println("Can't read file:", err)
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
	goodReadsClient = kbgoodreads.NewClient(goodReadsToken)

	// Wire the GitHub adapter (available even without token; client can be nil for limited usage)
	gh = githubadapter.New(gitHubClient)

	if len(gitHubToken) > 0 {
		username, err = gh.ViewerLogin(context.Background())
		if err != nil {
			fmt.Println("Can't retrieve GitHub profile:", err)
			os.Exit(1)
		}
	}

	// Wire the GoodReads adapter (requires token + user ID)
	gr = goodreadsadapter.New(goodReadsClient, goodReadsID)

	// Wire the Literal.club adapter
	lit = literaladapter.New()
	rss = rssadapter.New()
	// Build per-port services and the template composition service
	ghSvc := githubsvc.New(gh, username)
	grSvc := goodreadssvc.New(gr)
	litSvc := literalsvc.New(lit)
	rssSvc := rsssvc.New(rss)
	tplSvc := templatesvc.New(ghSvc, grSvc, litSvc, rssSvc)

	// Create template with template service functions to declutter main
	tpl, err := texttmpl.New("tpl").Funcs(texttmpl.FuncMap{
		/* GitHub */
		"recentContributions": tplSvc.RecentContributions,
		"recentPullRequests":  tplSvc.RecentPullRequests,
		"recentRepos":         tplSvc.RecentRepos,
		"recentForks":         tplSvc.RecentForks,
		"recentReleases":      tplSvc.RecentReleases,
		"followers":           tplSvc.Followers,
		"recentStars":         tplSvc.RecentStars,
		"gists":               tplSvc.Gists,
		"recentIssues":        tplSvc.RecentIssues,
		"sponsors":            tplSvc.Sponsors,
		"repo":                tplSvc.Repo,
		/* RSS */
		"rss": tplSvc.LatestRssFeeds,
		/* GoodReads */
		"goodReadsReviews":          tplSvc.GoodReadsReviews,
		"goodReadsCurrentlyReading": tplSvc.GoodReadsCurrentlyReading,
		/* Literal.club */
		"literalClubCurrentlyReading": tplSvc.LiteralCurrentlyReading,
		/* Utils */
		"humanize": tplSvc.Humanize,
		"reverse":  tplSvc.Reverse,
		"now":      time.Now,
		"contains": strings.Contains,
		"toLower":  strings.ToLower,
	}).Parse(string(tplIn))
	if err != nil {
		fmt.Println("Can't parse template:", err)
		os.Exit(1)
	}

	w := os.Stdout
	if len(*write) > 0 {
		f, err := os.Create(*write)
		if err != nil {
			fmt.Println("Can't create:", err)
			os.Exit(1)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				fmt.Println("Can't close file:", err)
			}
		}(f)
		w = f
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Can't render template:", err)
		os.Exit(1)
	}
}
