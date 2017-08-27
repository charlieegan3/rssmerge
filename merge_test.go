package main

import (
	"io/ioutil"
	"testing"
)

func TestMerge(t *testing.T) {
	fixtures := []string{
		"feed_ars",
		"feed_github",
		"feed_guardian",
		"feed_hn",
	}

	var feeds []string
	for _, v := range fixtures {
		feedContent, _ := ioutil.ReadFile("fixtures/" + v)
		feeds = append(feeds, string(feedContent))
	}

	mergedFeed := Merge(feeds)

	if len(mergedFeed.Items) != 119 {
		t.Errorf("Did not expect %v feed items", len(mergedFeed.Items))
	}

	if mergedFeed.Items[0].Title != "Git-daemon-dummy – deprecate git:// for https://" {
		t.Errorf("Did not expect title: %v", mergedFeed.Items[0].Title)
	}

	if mergedFeed.Items[0].Link.Href != "https://git.zx2c4.com/git-daemon-dummy/tree/README.md" {
		t.Errorf("Did not expect link: %v", mergedFeed.Items[0].Link)
	}

	if mergedFeed.Items[0].Description != "Hacker News: Newest" {
		t.Errorf("Did not expect description: %v", mergedFeed.Items[0].Description)
	}
}
