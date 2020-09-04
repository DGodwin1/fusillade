package main

import (
	"flag"
	"fmt"
)

func main() {
	//Line up the flags.
	configFile := flag.String("config", "", "provide config location")
	wantsGraphs := flag.Bool("graphOutput", false, "does user want graphs and stuff")
	wantsJSONOutput := flag.Bool("outJSON?", false, "does the user want a JSON representation of the report")
	flag.Parse()
	fmt.Println("configFile", *configFile)
	fmt.Println("wants chart:", *wantsGraphs)
	fmt.Println("wants JSON:", *wantsJSONOutput)

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
		fmt.Printf("Error with validating config: %q", err)
		return
	}

	fmt.Println(config.Urls, config.Rate, config.PauseLength, config.Count)
}
	//
	//// Prepare the test.
	//urls := config.Urls
	//count := config.Count
	//rate := time.Duration(config.Rate)
	//ticker := time.NewTicker(rate * time.Millisecond)
	//
	//resultChannel := make(chan UserJourneyResult)
	//
	//user := FakeUser{DelayTime: config.PauseLength}
	//
	//// Hit the URLS
	//DoConcurrentTask(func() {
	//	resultChannel <- WalkJourney(urls, user)
	//}, count, *ticker)
	//
	//// You've done the speedy stuff, now unload from the channel.
	//var responses []UserJourneyResult
	//for i := 0; i < count; i++ {
	//	result := <-resultChannel
	//	responses = append(responses, result)
	//}
	//
	//// Sort the responses by when the UserJourney actually started.
	//sort.Slice(responses, func(i, j int) bool {
	//	return responses[i].JourneyStart.Before(responses[j].JourneyStart)
	//})
	//
	////Now prepare a latency report
	//var latencies []int
	//var xAxis []int
	//// Get hold of the latency values so that you can put them into the graph.
	//for _, v := range responses {
	//	latencies = append(latencies, int(v.JourneyResponseTimeMS))
	//}
	//
	//// Number each user journey value
	//for i := 1; i <= len(latencies); i++ {
	//	xAxis = append(xAxis, i)
	//}
	//
	//MakeBarGraph("Latencies", "latencies.html", "Latency", xAxis, latencies)
	//
	////Prepare a count for the ResponseCodeCount
	//var codeXAxis []int
	//var countYAxis []int
	//
	//for code, count := range CountResponseCodes(responses) {
	//	codeXAxis = append(codeXAxis, code)
	//	countYAxis = append(countYAxis, count)
	//}
	//
	//MakeBarGraph("ResponseCodeCount", "responseCount.html", "Count", codeXAxis, countYAxis)

	//var jsonData []byte
	//jsonData, err = json.MarshalIndent(responses, "", "     ")
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Println(string(jsonData))
	//
	//err = ioutil.WriteFile("output.json", jsonData, 0644)
	//
//}
