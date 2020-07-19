# markscribe

[![Build Status](https://github.com/muesli/markscribe/workflows/build/badge.svg)](https://github.com/muesli/markscribe/actions)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/muesli/markscribe)
[![Go ReportCard](http://goreportcard.com/badge/muesli/markscribe)](http://goreportcard.com/report/muesli/markscribe)

Your personal markdown scribe with template-engine and Git(Hub) & RSS powers 📜

In order to access GitHub's API, markscribe expects you to provide a valid
GitHub token as an environment variable called `GITHUB_TOKEN`.

## Usage

Render a template to stdout:

    markscribe file.tpl

Render to a file:

    markscribe -write /tmp/output.md file.tpl

## Templates

You can find an example template to generate a GitHub profile README under
`templates/github-profile.tpl`. Make sure to fill in (or remove) placeholders,
like the RSS-feed or social media URLs.

Rendered it looks a little like my own profile page: https://github.com/muesli

## Functions

### Fetch repositories you recently created

```
{{range recentRepos 10}}
Name: {{.Name}}
Description: {{.Description}}
URL: {{.URL}})
Stars: {{.Stargazers}}
{{end}}
```

### Fetch recent releases you contributed to

```
{{range recentReleases 10}}
Name: {{.Name}}
Git Tag: {{.LastRelease.TagName}}
URL: {{.LastRelease.URL}}
Published: {{humanize .LastRelease.PublishedAt}}
{{end}}
```

### Fetch your latest followers

```
{{range followers 5}}
Username: {{.Login}}
URL: {{.URL}}
{{end}}
```

### Retrieve RSS feed

```
{{range rss "https://.../feed.xml" 5}}
Title: {{.Title}}
URL: {{.URL}}
Published: {{humanize .PublishedAt}}
{{end}}
```

## Template Engine

markscribe uses Go's powerful template engine. You can find its documentation
here: https://golang.org/pkg/text/template/

## FAQ

Q: "That's awesome, but can you expose more APIs and data?"  
A: Of course, just let me know what you'd like to do with markscribe and open a new issue!
