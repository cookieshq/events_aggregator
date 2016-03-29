package main

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Events []Event `json:"results"`
}

type Event struct {
	Name        string
	Description string
	Directions  string `json:"how_to_find_us"`
	TimeRaw     int64  `json:"time"`
	DurationRaw int64  `json:"duration"`
	URL         string `json:"event_url"`
	StartTime   time.Time
	EndTime     time.Time
	Duration    time.Duration

	Venue struct {
		Name string
	}
	Group struct {
		Name string
	}
}

func DecodeJSON(body io.ReadCloser) (events []Event, err error) {
	defer body.Close()

	var data Response

	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return []Event{}, err
	}

	for i, _ := range data.Events {
		event := &data.Events[i]

		event.StartTime = time.Unix(event.TimeRaw/1000, 0)

		if event.DurationRaw == 0 {
			// Assume 3 hours, as per Meetup docs
			event.Duration = time.Duration(3) * time.Hour
		} else {
			event.Duration = time.Duration(event.DurationRaw) * time.Millisecond
		}

		event.EndTime = event.StartTime.Add(event.Duration)
	}

	return data.Events, nil
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

	data, err := DecodeJSON(response.Body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	spew.Dump(data)
}
