package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

// Merge takes a list of raw feed strings, parses and then merges them
func Merge(rawFeeds []string) feeds.Feed {
	fp := gofeed.NewParser()

	var items itemList
	var feedList []string
	for _, feed := range rawFeeds {
		feed, err := fp.ParseString(feed)
		if err == nil {
			feedList = append(feedList, feed.Title)

			for _, item := range feed.Items {
				published := item.PublishedParsed
				updated := item.UpdatedParsed

				if published == nil && updated == nil {
					continue
				}

				time := published
				if time == nil {
					time = updated
				}

				items = append(items, convertItem(item, *time, feed))
			}
		}
	}

	sort.Sort(sort.Reverse(items))

	feed := feeds.Feed{
		Title:       fmt.Sprintf("RSSMerge Feed (%v feeds)", len(feedList)),
		Link:        &feeds.Link{Href: "PLACEHOLDER"},
		Description: strings.Join(feedList, ", "),
		Author:      &feeds.Author{Name: "RSSMerge", Email: "rssmerge@charlieegan3.com"},
		Created:     time.Now(),
		Items:       items,
	}

	return feed
}

func convertItem(item *gofeed.Item, created time.Time, feed *gofeed.Feed) *feeds.Item {
	if item.Link == "" {
		item.Link = feed.Link
	}
	return &feeds.Item{
		Title:       item.Title,
		Link:        &feeds.Link{Href: item.Link},
		Created:     created,
		Description: feed.Title,
	}
}

type itemList []*feeds.Item

func (slice itemList) Len() int {
	return len(slice)
}

func (slice itemList) Less(i, j int) bool {
	return slice[i].Created.Before(slice[j].Created)
}

func (slice itemList) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
