package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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

	Venue Venue

	Group struct {
		Name string
	}
}

type Venue struct {
	Name    string
	Address string `json:"address_1"`
	City    string
}

func (v Venue) IsEmpty() bool {
	return v == (Venue{})
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

func ordinal(day int) string {
	suffix := "th"

	switch day {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}

	return strconv.Itoa(day) + suffix
}

func formatDate(t time.Time) string {
	return fmt.Sprintf(t.Format("Monday %s January 2006 - 3:04pm"), ordinal(t.Day()))
}

func GetMeetupEvents(params map[string]string) (response *http.Response, err error) {
	paramsMap := url.Values{}

	for key, value := range params {
		paramsMap.Set(key, value)
	}

	meetupURL := url.URL{
		Scheme:   "https",
		Host:     "api.meetup.com",
		Path:     "/2/open_events",
		RawQuery: paramsMap.Encode(),
	}

	url := meetupURL.String()

	fmt.Println(url)

	response, err = http.Get(url)

	if response.StatusCode > 399 {
		errorJSON, err := ioutil.ReadAll(response.Body)

		response.Body.Close()

		if err != nil {
			return response, err
		}

		return response, errors.New(string(errorJSON))
	}

	return
}

func main() {
	now := time.Now()
	month := now.Month() + 1

	startTime := time.Date(now.Year(), month, 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := time.Date(now.Year(), month+1, 1, 0, 0, 0, 0, time.UTC).Add(-24 * time.Hour).Day()
	endTime := time.Date(now.Year(), month, lastDayOfMonth, 0, 0, 0, 0, time.UTC)

	fmt.Println(startTime, endTime)

	params := map[string]string{
		"key":         os.Getenv("MEETUP_API_KEY"),
		"category":    "34", // "tech"
		"city":        "Bristol",
		"country":     "GB",
		"text_format": "plain",
		"page":        "50",
		"time":        strconv.FormatInt(startTime.Unix()*1000, 10) + "," + strconv.FormatInt(endTime.Unix()*1000, 10),
	}

	response, err := GetMeetupEvents(params)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data, err := DecodeJSON(response.Body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	const tmplSrc = `## [{{.Name}}]({{.URL}})

**{{formatDate .StartTime}}**{{if not .Venue.IsEmpty}} - *{{.Venue.Name}}, {{.Venue.Address}}, {{trim .Venue.City}}*{{end}}

<div class="small">
{{.Description}}
</div>
`

	helpers := template.FuncMap{
		"formatDate": formatDate,
		"trim":       strings.TrimSpace,
	}

	tmpl, err := template.New("blog").Funcs(helpers).Parse(tmplSrc)

	if err != nil {
		panic(err)
	}

	for _, event := range data {
		tmpl.Execute(os.Stdout, event)
		fmt.Print("\n\n")
	}
}
