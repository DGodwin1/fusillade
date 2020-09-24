package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func MaxUserJourneyResponseLatency(r []UserJourneyResult) int {
	// GetMaxUserJourneyLatency is a custom max function
	// used to get the biggest latency value in a set of UserJourneys.

	// I think it's nicer to get actual data from r rather
	// than just initialise this to some random number.
	var max int64 = r[0].JourneyResponseTimeMS

	for _, v := range r[1:] {
		if v.JourneyResponseTimeMS > max {
			max = v.JourneyResponseTimeMS
		}
	}
	return int(max)
}

func MinUserJourneyResponseLatency(r []UserJourneyResult) int {
	min := r[0].JourneyResponseTimeMS

	for _, v := range r[1:] {
		if v.JourneyResponseTimeMS < min {
			min = v.JourneyResponseTimeMS
		}
	}
	return int(min)
}

//func FindPercentile(results []UserJourneyResult, p int) int{
//	// FindPercentile looks for the according percentile
//	//Get the latencies into a structure I can work from.
//	var latencies []int
//
//	for _, journey := range results{
//		latencies = append(latencies, int(journey.JourneyResponseTimeMS))
//		}
//
//
//	np := math.Ceil(float64(p) * float64(len(latencies))/100)
//	fmt.Println(np)
//
//	// Sort the latencies.
//	if !sort.IntsAreSorted(latencies){
//		sort.Ints(latencies)
//	}
//
//	fmt.Println(latencies)
//	if latencies[int(np)-1] != nil {
//
//	}
//	return latencies[int(np)-1]

//}

func CountResponseCodes(r []UserJourneyResult) map[int]int {
	// CountResponseCodes is here to go through a load of responses
	// and to generate a single, flat map that shows, of all the
	// URLs visited, what response codes we received. It can be used
	// later to help with data analysis, or just to get a quick view of
	// what responses came up most often in a test.

	codes := map[int]int{}
	for _, v := range r {
		for code, count := range v.Codes {
			codes[code] += count
		}
	}
	return codes
}

type SummativeData struct {
	ResponseCodeCount  map[int]int
	MinJourneyResponse int
	MaxJourneyResponse int
}

func ConstructSummativeStats(r []UserJourneyResult) SummativeData {
	// As the name suggests, this constructs a summative
	// report that is returned to the user.

	return SummativeData{
		ResponseCodeCount:  CountResponseCodes(r),
		MinJourneyResponse: MinUserJourneyResponseLatency(r),
		MaxJourneyResponse: MaxUserJourneyResponseLatency(r),
	}
}

func ConstructJSONReport(r []UserJourneyResult){
	var jsonData []byte
	jsonData, err := json.MarshalIndent(r, "", "     ")
	if err != nil {
		log.Println(err)
	}
	err = ioutil.WriteFile("output.json", jsonData, 0644)
}
