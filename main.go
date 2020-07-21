package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
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
		"recentRepos":    recentRepos,
		"recentReleases": recentReleases,
		"followers":      recentFollowers,
		"rss":            rssFeed,
		"humanize":       humanized,
	}).Parse(string(tplIn))
	if err != nil {
		fmt.Println("Can't parse template:", err)
		os.Exit(1)
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
