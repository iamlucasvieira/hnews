package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestApiRequest tests the api.request method
func TestApiRequest(t *testing.T) {
	// First we create a mock response
	mockResponse := HnResponse{
		Hits: []Hit{
			{
				Title:  "Test title",
				URL:    "http://example.com",
				Points: 100,
			},
		},
	}

	// Then we create a mock serer
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(mockResponse)
		_, err := w.Write(resp)
		if err != nil {
			http.Error(w, "Failed to write response of mock server", http.StatusInternalServerError)
		}
	}))

	// Add a defer to close the mock server
	defer mockServer.Close()

	// Test the request function
	a := getApi()
	response, err := a.request(mockServer.URL)
	if err != nil {
		t.Fatalf("Received unexpected error: %s", err)
	}

	// Test if the response contains 1 hit
	if len(response.Hits) != 1 {
		t.Errorf("Expected 1 hit, got %d", len(response.Hits))
	}

	// Test if the response contains the correct title
	if response.Hits[0].Title != "Test title" {
		t.Errorf("Expected 'Test title', got %s", response.Hits[0].Title)
	}

	// Test if the response contains the correct url
	if response.Hits[0].URL != "http://example.com" {
		t.Errorf("Expected 'http://example.com', got %s", response.Hits[0].URL)
	}

	// Test if the response contains the correct points
	if response.Hits[0].Points != 100 {
		t.Errorf("Expected 100, got %d", response.Hits[0].Points)
	}
}
