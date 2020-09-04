package main

import (
	"github.com/go-echarts/go-echarts/charts"
	"log"
	"os"
)

func MakeBarGraph(GraphTitle, fileName, yAxisName string, xAxis, yAxis []int) {
	// Setup the graph
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{
		Title: GraphTitle,
	})

	// Sort out the axises and their values.
	bar.AddXAxis(xAxis)
	bar.AddYAxis(yAxisName, yAxis)

	// Create the file where the graph will be stored.
	f, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
		return
	}

	bar.Render(f)

}
