// hnews is a command line tool that fetches the top stories from Hacker News
package main

import "fmt"

func main() {
	a := getApi()
	urls := getHnUrl()
	urls.topStories()

	data, err := a.request(urls.string())

	if err != nil {
		fmt.Println(err)
		return
	}

	// Loop through the hits and print the title and url
	for _, hit := range data.Hits {
		fmt.Println(hit.Title)
		fmt.Println(hit.URL)
		fmt.Println(hit.Points)
	}

}
