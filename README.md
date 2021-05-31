# markscribe

[![Latest Release](https://img.shields.io/github/release/muesli/markscribe.svg)](https://github.com/muesli/markscribe/releases)
[![Build Status](https://github.com/muesli/markscribe/workflows/build/badge.svg)](https://github.com/muesli/markscribe/actions)
[![Go ReportCard](https://goreportcard.com/badge/muesli/markscribe)](https://goreportcard.com/report/muesli/markscribe)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/muesli/markscribe)

Your personal markdown scribe with template-engine and Git(Hub) & RSS powers 📜

You can run markscribe as a GitHub Action: [readme-scribe](https://github.com/muesli/readme-scribe/)

## Usage

Render a template to stdout:

    markscribe template.tpl

Render to a file:

    markscribe -write /tmp/output.md template.tpl

## Installation

### Packages & Binaries

If you use Brew, you can simply install the package:

    brew install muesli/tap/markscribe

Or download a binary from the [releases](https://github.com/muesli/markscribe/releases)
page. Linux (including ARM) binaries are available, as well as Debian and RPM
packages.

### Build From Source

Alternatively you can also build `markscribe` from source. Make sure you have a
working Go environment (Go 1.11 or higher is required). See the
[install instructions](https://golang.org/doc/install.html).

To install markscribe, simply run:

    go get github.com/muesli/markscribe

## Templates

You can find an example template to generate a GitHub profile README under
[`templates/github-profile.tpl`](templates/github-profile.tpl). Make sure to fill in (or remove) placeholders,
like the RSS-feed or social media URLs.

Rendered it looks a little like my own profile page: https://github.com/muesli

## Functions

### RSS feed

```
{{range rss "https://domain.tld/feed.xml" 5}}
Title: {{.Title}}
URL: {{.URL}}
Published: {{humanize .PublishedAt}}
{{end}}
```

### Your recent contributions

```
{{range recentContributions 10}}
Name: {{.Repo.Name}}
Description: {{.Repo.Description}}
URL: {{.Repo.URL}})
Occurred: {{humanize .OccurredAt}}
{{end}}
```

This function requires GitHub authentication!

### Your recent pull requests

```
{{range recentPullRequests 10}}
Title: {{.Title}}
URL: {{.URL}}
State: {{.State}}
CreatedAt: {{humanize .CreatedAt}}
Repository name: {{.Repo.Name}}
Repository description: {{.Repo.Description}}
Repository URL: {{.Repo.URL}}
{{end}}
```

### Repositories you recently starred

```
{{range recentStars 10}}
Name: {{.Name}}
Description: {{.Description}}
URL: {{.URL}})
Stars: {{.Stargazers}}
{{end}}
```

This function requires GitHub authentication!

### Repositories you recently created

```
{{range recentRepos 10}}
Name: {{.Name}}
Description: {{.Description}}
URL: {{.URL}})
Stars: {{.Stargazers}}
{{end}}
```

This function requires GitHub authentication!

### Recent releases you contributed to

```
{{range recentReleases 10}}
Name: {{.Name}}
Git Tag: {{.LastRelease.TagName}}
URL: {{.LastRelease.URL}}
Published: {{humanize .LastRelease.PublishedAt}}
{{end}}
```

This function requires GitHub authentication!

### Your published gists

```
{{range gists 10}}
Name: {{.Name}}
Description: {{.Description}}
URL: {{.URL}}
Created: {{humanize .CreatedAt}}
{{end}}
```

This function requires GitHub authentication!

### Your latest followers

```
{{range followers 5}}
Username: {{.Login}}
Name: {{.Name}}
Avatar: {{.AvatarURL}}
URL: {{.URL}}
{{end}}
```

This function requires GitHub authentication!

### Your sponsors

```
{{range sponsors 5}}
Username: {{.User.Login}}
Name: {{.User.Name}}
Avatar: {{.User.AvatarURL}}
URL: {{.User.URL}}
Created: {{humanize .CreatedAt}}
{{end}}
```

This function requires GitHub authentication!

### Your GoodReads reviews

```
{{range goodReadsReviews 5}}
- {{.Book.Title}} - {{.Book.Link}} - {{.Rating}} - {{humanize .DateUpdated}}
{{- end}}
```

This function requires GoodReads API key!

### Your GoodReads currently reading books

```
{{range goodReadsCurrentlyReading 5}}
- {{.Book.Title}} - {{.Book.Link}} - {{humanize .DateUpdated}}
{{- end}}
```

This function requires GoodReads API key!

## Template Engine

markscribe uses Go's powerful template engine. You can find its documentation
here: https://golang.org/pkg/text/template/

## Template Helpers

markscribe comes with a few handy template helpers:

To format timestamps, call `humanize`:

```
{{humanize .Timestamp}}
```

To reverse the order of a slice, call `reverse`:

```
{{reverse (rss "https://domain.tld/feed.xml" 5)}}
```

## GitHub Authentication

In order to access some of GitHub's API, markscribe requires you to provide a
valid GitHub token in an environment variable called `GITHUB_TOKEN`. You can
create a new token by going to your profile settings:

`Developer settings` > `Personal access tokens` > `Generate new token`

## GoodReads API key

In order to access some of GoodReads' API, markscribe requires you to provide a
valid GoodReads key in an environment variable called `GOODREADS_TOKEN`. You can
create a new token by going [here](https://www.goodreads.com/api/keys).
Then you need to go to your repository and add it, `Settings -> Secrets -> New secret`.
You also need to set your GoodReads user ID in your secrets as `GOODREADS_USER_ID`.

## FAQ

Q: That's awesome, but can you expose more APIs and data?  
A: Of course, just open a new issue and let me know what you'd like to do with markscribe!

Q: That's awesome, but I don't have my own server to run this on. Can you help?  
A: Check out [readme-scribe](https://github.com/muesli/readme-scribe/), a GitHub Action that runs markscribe for you!
