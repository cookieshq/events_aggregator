package main

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"os"
)

type Response struct {
	Events []Event `json:"results"`
}

type Event struct {
	Name        string
	Description string
	Directions  string `json:"how_to_find_us"`
	Time        int
	Duration    int
	URL         string `json:"event_url"`

	Venue struct {
		Name string
	}
	Group struct {
		Name string
	}
}

func main() {
	apiKey := os.Getenv("MEETUP_API_KEY")
	techCategoryId := 34

	url := fmt.Sprintf("https://api.meetup.com/2/open_events?key=%s&page=20&category=%d&city=%s&country=%s&text_format=plain", apiKey, techCategoryId, "Bristol", "GB")

	fmt.Println(url)

	response, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer response.Body.Close()

	var data Response

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	spew.Dump(data)
}
