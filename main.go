package main

import (
	"flag"
	"fmt"
	"sort"
	"time"
)

func main() {
	//Line up the flags.
	configFile := flag.String("c", "", "provide config location.")
	wantsGraphs := flag.Bool("g", false, "set flag to true if you want graphs for user journey and response code counts.")
	wantsJSON := flag.Bool("j", false, "set flag to true if you want user journey data to be saved a JSON file.")
	flag.Parse()

	if *configFile == ""{
		//generate fake JSON
		fmt.Println(`This tool needs a configuration file to work from. To make things easier, I have made one called 'config.json' which you can edit and use. Once you're happy, simply re-run the program with the -config flag set to 'config.json'.`)

		//create the file as promised:
		MakeFakeJSONFile()
		return
	}

	//Parse the file.
	config, err := ParseConfigFile(*configFile)

	if err != nil {
		fmt.Printf("error parsing config: %q", err)
		return
	}

	// Validate that everything in the file is okay.
	v := ConfigValidator{}
	_, err = v.Validate(config)

	if err != nil {
		fmt.Printf("error validating config: %q", err)
		return
	}

	// Prepare the test based on the values stored in the config.
	urls := config.Urls
	count := config.Count
	rate := time.Duration(config.Rate)
	ticker := time.NewTicker(rate * time.Millisecond)
	resultChannel := make(chan UserJourneyResult)
	user := FakeUser{DelayTime: config.PauseLength}

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

	// Now sort the responses based on when the user journey was launched
	// otherwise it's based on when journeys completed.
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].JourneyStart.Before(responses[j].JourneyStart)
	})

	if *wantsGraphs{
		// Make with the data vis.
		//Now prepare a latency report
		var latencies []int
		var xAxis []int

		// Get hold of the latency values so that you can put them into the graph.
		for _, v := range responses {
			latencies = append(latencies, int(v.JourneyResponseTimeMS))
		}

		// Number each user journey value
		for i := 1; i <= len(latencies); i++ {
			xAxis = append(xAxis, i)
		}

		MakeBarGraph("Latencies", "latencies.html", "Latency in MS", xAxis, latencies)

		//Prepare a count for the ResponseCodeCount
		var codeXAxis []int
		var countYAxis []int

		for code, count := range CountResponseCodes(responses) {
			codeXAxis = append(codeXAxis, code)
			countYAxis = append(countYAxis, count)
		}

		MakeBarGraph("ResponseCodeCount", "responseCount.html", "Response code count", codeXAxis, countYAxis)
	}

	if *wantsJSON{
		ConstructJSONReport(responses)
	}

	d := ConstructSummativeStats(responses)
	fmt.Println("Min user journey response time:", d.MinJourneyResponse)
	fmt.Println("Max user journey response time:", d.MaxJourneyResponse)
	fmt.Println("Response codes encountered time:", d.ResponseCodeCount)
	fmt.Println("Test complete")
}
