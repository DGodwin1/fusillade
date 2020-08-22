package main

import (
	"fmt"
	"time"
)

func main() {
	//TODO: get the config from _somewhere
	// Get the config.


	// Parse the file.
	config, err := ParseConfigFile("test_config.json")

	if err != nil{
		fmt.Println("There's been an issue parsing the config file.")
		fmt.Printf("Here's the error %q", err)
		return
	}

	//Validate that everything in the file is okay.
	v := ConfigValidator{}
	_, err = v.Validate(config)

	if err != nil{
		fmt.Println("There's been an issue validating the configuration.")
		fmt.Printf("Here's the error %q", err)
		return
	}

	//Prepare the test.
	urls := config.Urls
	count := config.Count
	ticker := time.NewTicker(100 * time.Millisecond) //TODO: take from config.

	resultChannel := make(chan UserJourneyResult)

	reader := EndUserReader{}
	MillisecondReadingTime := 10000

	// Hit the URLS
	DoConcurrentTask(func() {
		resultChannel <- WalkJourney(urls, reader, MillisecondReadingTime)
	}, count, *ticker)


	// You've done the speedy stuff, now unload from the channel.
	var responses []UserJourneyResult
	for i := 0; i < count; i++ {
		result := <-resultChannel
		responses = append(responses, result)
	}

	for _, v := range responses {
		fmt.Println(v.JourneyResponseTimeMS)
	}
	//Now prepare a report
	MakeBarGraph(responses)

}
