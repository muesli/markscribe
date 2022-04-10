package cmd

import (
	"errors"
	"flag"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/muesli/markscribe/pkg/gh"
	"github.com/muesli/markscribe/pkg/gr"
	"github.com/muesli/markscribe/pkg/rss"
	"github.com/muesli/markscribe/pkg/utils"
)

var (
	ErrUsage = errors.New("usage error")
)

func Execute() error {
	write := flag.String("write", "", "write output to")

	flag.Parse()

	if len(flag.Args()) == 0 {
		return ErrUsage
	}

	tplIn, err := ioutil.ReadFile(flag.Args()[0])
	if err != nil {
		return err
	}

	// Create GitHub and GoodReads clients
	ghClient := gh.New()
	grClient := gr.New()

	// Set GitHub username
	err = ghClient.GetUsername()
	if err != nil {
		return err
	}

	tpl, err := template.New("tpl").Funcs(template.FuncMap{
		/* GitHub */
		"recentContributions": ghClient.GetRecentContributions,
		"recentPullRequests":  ghClient.GetRecentPullRequests,
		"recentReleases":      ghClient.GetRecentReleases,
		"recentRepos":         ghClient.GetRecentRepos,
		"recentStars":         ghClient.GetRecentStars,
		"followers":           ghClient.GetRecentFollowers,
		"sponsors":            ghClient.GetSponsors,
		"gists":               ghClient.GetRecentGists,
		/* RSS */
		"rss": rss.Feed,
		/* GoodReads */
		"goodReadsCurrentlyReading": grClient.GetCurrentlyReading,
		"goodReadsReviews":          grClient.GetReviews,
		/* Utils */
		"humanize": utils.Humanized,
		"toLower":  strings.ToLower,
		"toUpper":  strings.ToUpper,
		"reverse":  utils.Reverse,
		"now":      time.Now,
	}).Parse(string(tplIn))
	if err != nil {
		return err
	}

	w := os.Stdout
	if len(*write) > 0 {
		f, err := os.Create(*write)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	return tpl.Execute(w, nil)
}
