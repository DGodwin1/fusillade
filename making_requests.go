package fusillade

import (
	"net/http"
	"time"
)

type Response struct{
	StatusCode int
	//TODO: ID int
	//TODO: RequestStart
	//TODO: RequestFinished
	//TODO: RequestDelta
	ResponseTime int64
}

func MakeRequest(url string) Response{
	start := time.Now()
	request, _ := http.Get(url)
	end := CalculateMSDelta(start, time.Now())

	return Response{StatusCode: request.StatusCode, ResponseTime: end}
}

func CalculateMSDelta(start time.Time, end time.Time) (ResponseTime int64){
	// CalculateMSDelta does as it suggests, it takes two timestamps and
	// calculates the delta between them by subtracting the start
	// from the end. It represents the final result in milliseconds.
	return end.Sub(start).Milliseconds()

}
