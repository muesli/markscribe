package main

import (
	"io/ioutil"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestRSSExtensions(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	content, err := ioutil.ReadFile("testdata/rss-with-extensions.xml")
	if err != nil {
		t.Fatal(err)
	}
	httpmock.RegisterResponder("GET", "https://example.com/rss-with-extensions",
		httpmock.NewStringResponder(200, string(content)))

	returnCount := 5
	feed := rssFeed("https://example.com/rss-with-extensions", returnCount)

	// Make sure we got exactly 5 results back
	require.Equal(t, returnCount, len(feed))

	// First entry from the test file
	first := feed[0]

	// Make sure the first entry has an extension is the right movie
	require.Equal(t, "Dashcam", first.Extensions["letterboxd"]["filmTitle"][0].Value)
}

func TestRSSPanic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://example.com/panic-rss",
		httpmock.NewStringResponder(200, `This is not xml`))

	require.Panics(t, func() { rssFeed("https://example.com/panic-rss", 5) })
}
