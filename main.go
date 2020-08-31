package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"time"
)

func main() {
	//TODO: get the config from _somewhere
	// Get the config.
	//
	//Parse the file.
	config, err := ParseConfigFile("test_config.json")

	if err != nil{
		fmt.Printf("error parsing config: %q", err)
		return
	}

	// Validate that everything in the file is okay.
	v := ConfigValidator{}
	_, err = v.Validate(config)

	if err != nil{
		fmt.Printf("Error with validating config: %q", err)
		return
	}

	// Prepare the test.
	urls := config.Urls
	count := config.Count
	ticker := time.NewTicker(100 * time.Millisecond) //TODO: take from config.

	resultChannel := make(chan UserJourneyResult)

	user := FakeUser{DelayTime: 100}

	// Hit the URLS
	DoConcurrentTask(func() {
		resultChannel <- WalkJourney(urls, user)
	}, count, *ticker)


	// You've done the speedy stuff, now unload from the channel.
	var responses []UserJourneyResult
	for i := 0; i < count; i++ {
		result := <-resultChannel
		responses = append(responses, result)
	}

	// Sort the responses by when the UserJourney actually started.
	sort.Slice(responses, func(i, j int) bool{
		return responses[i].JourneyStart.Before(responses[j].JourneyStart)
	})

	//Now prepare a report
	//MakeBarGraph(responses)
	//MakePieChart(responses)

	var jsonData []byte
	jsonData, err = json.MarshalIndent(responses, "", "     ")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(jsonData))

	err = ioutil.WriteFile("output.json", jsonData, 0644)

	fmt.Println("thank god for that")
}

