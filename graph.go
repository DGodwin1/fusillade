package main

import (
	"github.com/go-echarts/go-echarts/charts"
	"log"
	"os"
)

func MakeBarGraph(GraphTitle, fileName, yAxisName string, xValues, yValues []int){
	// Setup the graph
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{
		Title:    GraphTitle,
	})

	// Stitch together the xAxis and latency values.
	bar.AddXAxis(xValues)

	bar.AddYAxis(yAxisName, yValues)

	// Create the file.
	f, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
	}

	bar.Render(f)

}