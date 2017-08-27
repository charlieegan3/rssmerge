package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestFetch(t *testing.T) {
	rawFeeds := []string{
		"http://golangweekly.com/rss",
		"http://n-gate.com/hackernews/index.atom",
		"http://rubyweekly.com/rss",
		"https://blog.mode7games.com/feed",
		"https://charlieegan3.com/feed",
		"https://githubengineering.com/atom.xml",
		"https://increment.com/feed.xml",
		"https://unboxed.co/blog/feed.xml"}

	var parsedFeeds []url.URL
	for _, v := range rawFeeds {
		feedURL, _ := url.Parse(v)
		parsedFeeds = append(parsedFeeds, *feedURL)
	}

	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/source_gist" {
			pageContent := strings.Join(rawFeeds, "\n")

			fmt.Fprintln(w, pageContent)
		}
	}))

	sourceURL, _ := url.Parse(localServer.URL + "/source_gist")

	resultingFeeds, err := FetchList(sourceURL)

	if err != nil {
		t.Error(err)
	}

	if len(resultingFeeds) != len(parsedFeeds) {
		t.Errorf("Did not expect %v feeds", len(resultingFeeds))
		return
	}

	for i, v := range parsedFeeds {
		if v != resultingFeeds[i] {
			t.Errorf("Expected %v to be %v", v, resultingFeeds[i])
			return
		}
	}
}

func TestFetchInvalid(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/source_gist" {
			fmt.Fprintln(w, "blurgh bleep bloop\nnot a feed\n\n\t\t")
		}
	}))

	sourceURL, _ := url.Parse(localServer.URL + "/source_gist")

	resultingFeeds, err := FetchList(sourceURL)

	if err != nil {
		t.Error(err)
	}

	if len(resultingFeeds) != 0 {
		t.Errorf("Expected an empty list of feeds, got %v", resultingFeeds)
		return
	}
}
