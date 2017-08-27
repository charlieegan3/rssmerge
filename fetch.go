package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Fetch will visit a plaintext URL and extract the list of links
// Suggested sources are gists or pastebin
func Fetch(source *url.URL) ([]url.URL, error) {
	resp, err := http.Get(source.String())
	if err != nil {
		return []url.URL{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []url.URL{}, err
	}

	rawFeeds := strings.Split(strings.TrimSpace(string(body)), "\n")
	var parsedFeeds []url.URL

	for _, v := range rawFeeds {
		feedURL, err := url.Parse(v)
		if err == nil && feedURL.Scheme != "" {
			parsedFeeds = append(parsedFeeds, *feedURL)
		}
	}

	return parsedFeeds, nil
}
