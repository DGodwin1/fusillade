package main

import (
	"fmt"
	"time"
)

func main() {
	//Parse the file.
	config, _ := ParseConfigFile("test_config.json")

	//Validate that everything in the file is okay.

	//Prepare the test.
	urls := config.Urls
	count := config.Count
	ticker := time.NewTicker(100 * time.Millisecond) //TODO: take from config.

	resultChannel := make(chan UserJourneyResult)

	reader := EndUserReader{}


	// Hit the URLS
	DoConcurrentTask(func() {
		resultChannel <- WalkJourney(urls, reader)
	}, count, *ticker)


	// You've done the speedy stuff, now unload from the channel.
	var responses []UserJourneyResult
	for i := 0; i < count; i++ {
		result := <-resultChannel
		responses = append(responses, result)
	}

	for _, v := range responses{
		fmt.Println(v.JourneyResponseTimeMS)
	}


	// Now prepare a report



}
