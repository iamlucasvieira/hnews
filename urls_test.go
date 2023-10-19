package main

import "testing"

// TestGetHnUrl tests the getHnUrl function
func TestGetHnUrl(t *testing.T) {
	u := getHnUrl()
	if u == nil {
		t.Error("Expected hnUrl, got nil")
	}
	if u.string() != "http://hn.algolia.com/api/v1/search" {
		t.Errorf("Expected base URL, got %s", u.string())
	}
}

// TestHnUrlBase tests the base method
func TestHnUrlBase(t *testing.T) {
	u := getHnUrl()

	baseExpected := "http://hn.algolia.com/api/v1/search"

	base := u.base()
	if base != baseExpected {
		t.Errorf("Expected base URL without query, got %s", base)
	}

	// Modify the query string
	u.topStories()

	// Base should still be the same
	baseAfter := u.base()

	if baseAfter != baseExpected {
		t.Errorf("Expected base URL without query, got %s", baseAfter)
	}
}

// TestHnUrlReset tests the reset method
func TestHnUrlReset(t *testing.T) {
	u := getHnUrl()
	u.topStories() // This sets a query string
	u.reset()
	if u.URL.RawQuery != "" {
		t.Errorf("Expected empty query string, got %s", u.URL.RawQuery)
	}
}

// TestHnUrlTopStories tests the topStories method
func TestHnUrlTopStories(t *testing.T) {
	u := getHnUrl()
	u.topStories()
	if u.URL.RawQuery != "tags=front_page" {
		t.Errorf("Expected top stories query, got %s", u.URL.RawQuery)
	}
}

// TestHnUrlNewStories tests the newStories method
func TestHnUrlNewStories(t *testing.T) {
	u := getHnUrl()
	u.newStories()
	if u.URL.RawQuery != "tags=story" {
		t.Errorf("Expected new stories query, got %s", u.URL.RawQuery)
	}
}
