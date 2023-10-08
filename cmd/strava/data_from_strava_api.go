package main

import(
	"fmt"
	"net/http"
	"time"
)

var url string = "https://www.strava.com/"

func main() {

	// Send an HTTP GET request to the URL
	startTime := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Calculate the time it took to receive a response
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)

	// Print the UTC time and response time
	fmt.Printf("Current time (UTC): %s\n", time.Now().UTC().Format(time.RFC3339))
	fmt.Printf("Strava Response Time: %s\n", elapsedTime)
}