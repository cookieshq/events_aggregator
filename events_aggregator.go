package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

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

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(data))
}
