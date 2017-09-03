package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestHandler(t *testing.T) {
	feedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/feed_ars" {
			pageContent, _ := ioutil.ReadFile("fixtures/feed_ars")
			fmt.Fprintln(w, string(pageContent))
		} else if r.URL.Path == "/feed_guardian" {
			pageContent, _ := ioutil.ReadFile("fixtures/feed_guardian")
			fmt.Fprintln(w, string(pageContent))
		}
	}))

	gist := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/gist_path" {
			pageContent := feedServer.URL + "/feed_ars\n" +
				feedServer.URL + "/feed_guardian\n"

			fmt.Fprintln(w, pageContent)
		}
	}))

	req, err := http.NewRequest("GET", "/?source="+gist.URL+"/gist_path&day=true", nil)
	if err != nil {
		t.Error(err.Error())
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RSSMergeHandler)

	handler.ServeHTTP(rr, req)

	body := rr.Body.String()

	fp := gofeed.NewParser()
	feed, _ := fp.ParseString(body)

	if fmt.Sprintf("%v", feed.Items[0].Title) != "Canon Bill Hall obituary" {
		t.Errorf("got %v want %v", feed.Items[0].Title, "Canon Bill Hall obituary")
	}
	if fmt.Sprintf("%v", feed.Items[0].Link) != "https://www.theguardian.com/world/2017/aug/27/canon-bill-hall-obituary" {
		t.Errorf("got %v want %v", feed.Items[0].Link, "https://www.theguardian.com/world/2017/aug/27/canon-bill-hall-obituary")
	}
	if fmt.Sprintf("%v", feed.Items[0].PublishedParsed.String()) != "2017-08-27 17:15:33 +0000 UTC" {
		t.Errorf("got %v want %v", feed.Items[0].PublishedParsed.String(), "2017-08-27 17:15:33 +0000 UTC")
	}

	if fmt.Sprintf("%v", feed.Items[len(feed.Items)-1].Title) != "Secret Service conducts live test of ShotSpotter system at White House" {
		t.Errorf("got %v want %v", feed.Items[len(feed.Items)-1].Title, "Secret Service conducts live test of ShotSpotter system at White House")
	}
	if fmt.Sprintf("%v", feed.Items[len(feed.Items)-1].PublishedParsed) != "2017-08-26 17:41:28 +0000 UTC" {
		t.Errorf("got %v want %v", feed.Items[len(feed.Items)-1].PublishedParsed.String(), "2017-08-26 17:41:28 +0000 UTC")
	}
}
