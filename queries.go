package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Hit struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Points    int    `json:"points"`
	StoryText string `json:"story_text"`
}

type HnResponse struct {
	Hits []Hit `json:"hits"`
}

type api struct{}

func getApi() api {
	return api{}
}

func (a api) request(url string) (HnResponse, error) {
	fmt.Println("Requesting:", url)
	// Create client
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		return HnResponse{}, err
	}

	// Send request
	res, getErr := client.Do(req)
	if getErr != nil {
		fmt.Println(getErr)
		return HnResponse{}, getErr
	}

	// Close response body at end of function
	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}()

	// Read response body
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println(readErr)
		return HnResponse{}, readErr
	}

	var response HnResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		return HnResponse{}, err
	}
	return response, nil
}
