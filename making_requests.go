package main

import (
	"net/http"
	"time"
)

type Response struct{
	StatusCode int
	//TODO: ID int. Might be user when it comes to concurrency to order the requests based on their ID.

	//Store timings for a particular request.
	RequestStart time.Time
	RequestFinished time.Time
	ResponseTime int64
}

func MakeRequest(url string) Response{
	// MakeRequest takes a URL and returns a Response containing
	// all of the necessary information. CalculateMSDelta is moved
	// out to make the testing of time calculation independent of
	// the function actually making a request.
	start := time.Now()
	request, _ := http.Get(url)
	end := time.Now()
	rt := CalculateMSDelta(start, end)

	return Response{StatusCode: request.StatusCode,
		RequestStart: start,
		RequestFinished: end,
		ResponseTime: rt,
	}
}

func MakeConcurrentRequests(url string) []Response{
	// MakeConcurrentRequests is, for all intents and purposes,
	// a wrapper function that facilitates concurrent execution
	// of MakeRequest. Rather than force MakeRequest to handle concurrency,
	// as well as reporting and concurrency, I'm pulling this out into
	// something that can manage it. Single responsibility and all that.

	var responses []Response

	for i := 0; i < 100; i++{
		responses = append(responses, MakeRequest(url))}

	return responses
}

func CalculateMSDelta(start time.Time, end time.Time) (ResponseTime int64){
	// CalculateMSDelta does as it suggests, it takes two timestamps and
	// calculates the delta between them by subtracting the start
	// from the end. It represents the final result in milliseconds.
	return end.Sub(start).Milliseconds()
}
