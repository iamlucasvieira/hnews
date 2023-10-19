package main

import (
	"fmt"
	"net/url"
)

// urlHandler define the UrlBuilder interface
type urlHandler interface {
	base() string
	reset()
	topStories()
	newStories()
	string() string
}

// hnUrl is the struct that implements the UrlBuilder interface and uses the url package
type hnUrl struct {
	*url.URL
}

// string returns the url as a string
func (u hnUrl) string() string {
	return u.URL.String()
}

// reset resets the query string
func (u hnUrl) reset() {
	u.RawQuery = ""
}

// base returns the base url
func (u hnUrl) base() string {
	// Remove all query strings
	urlCopy := *u.URL
	urlCopy.RawQuery = ""
	return urlCopy.String()
}

// topStories sets the query string to fetch the top stories
func (u hnUrl) topStories() {
	query := u.Query()
	query.Set("tags", "front_page")
	u.RawQuery = query.Encode()
}

// newStories sets the query string to fetch the new stories
func (u hnUrl) newStories() {
	query := u.Query()
	query.Set("tags", "story")
	u.RawQuery = query.Encode()
}

// getHnUrl returns a hnUrl struct
func getHnUrl() *hnUrl {
	parsedUrl, err := url.Parse("http://hn.algolia.com/api/v1/search")
	if err != nil {
		fmt.Println(err)
	}
	return &hnUrl{parsedUrl}
}
