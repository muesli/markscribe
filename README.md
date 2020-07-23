# markscribe

[![Build Status](https://github.com/muesli/markscribe/workflows/build/badge.svg)](https://github.com/muesli/markscribe/actions)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/muesli/markscribe)
[![Go ReportCard](http://goreportcard.com/badge/muesli/markscribe)](http://goreportcard.com/report/muesli/markscribe)

Your personal markdown scribe with template-engine and Git(Hub) & RSS powers 📜

You can run markscribe as a GitHub Action: [readme-scribe](https://github.com/muesli/readme-scribe/)

## Usage

Render a template to stdout:

    markscribe template.tpl

Render to a file:

    markscribe -write /tmp/output.md template.tpl

## Templates

You can find an example template to generate a GitHub profile README under
`templates/github-profile.tpl`. Make sure to fill in (or remove) placeholders,
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

## Template Engine

markscribe uses Go's powerful template engine. You can find its documentation
here: https://golang.org/pkg/text/template/

## GitHub Authentication

In order to access some of GitHub's API, markscribe requires you to provide a
valid GitHub token in an environment variable called `GITHUB_TOKEN`. You can
create a new token by going to your profile settings:

`Developer settings` > `Personal access tokens` > `Generate new token`

## FAQ

Q: That's awesome, but can you expose more APIs and data?
A: Of course, just open a new issue and let me know what you'd like to do with markscribe!

Q: That's awesome, but I don't have my own server to run this on. Can you help?
A: Check out [readme-scribe](https://github.com/muesli/readme-scribe/), a GitHub Action that runs markscribe for you!
