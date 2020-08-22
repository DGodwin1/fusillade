package main

import (
	"github.com/go-echarts/go-echarts/charts"
	"log"
	"os"
)

func MakeBarGraph(results []UserJourneyResult){
	// MakeBarGraph takes in a load of UserJourney results
	// and then makes a nice bar graph of all of the latency data.

	var latencies []int64
	var xAxis []int

	// Setup the graph
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{
		Title:         "Latency data",
		Subtitle:      "Summative latency data for each user journey",
	})


	// Get hold of the latency values so that you can put them into the graph.
	for _, v := range results{
		latencies = append(latencies, v.JourneyResponseTimeMS)
	}

	// Number each user journey value
	for i := 1; i <= len(latencies); i++{
		xAxis = append(xAxis, i)
	}


	// Stitch together the xAxis and latency values.
	bar.AddXAxis(xAxis).AddYAxis("Latency values", latencies)

	// Create the file.
	f, err := os.Create("latencies.html")
	if err != nil {
		log.Println(err)
	}

	bar.Render(f)


}

func MakePieChart(){
	// TODO: take in []UserJourney]
	// MakePieChart takes in a load of UserJourney results
	// and creates a nice pie chart of all the response data.
}
