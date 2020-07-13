package main

import (
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	StatusCode int
	//TODO: ID int. Might be user when it comes to concurrency to order the requests based on their ID.

	//Store timings for a particular request.
	RequestStart    time.Time
	RequestFinished time.Time
	ResponseTime    int64
}

func MakeRequest(url string) Response {
	// MakeRequest takes a URL and returns a Response containing
	// all of the necessary information. CalculateMSDelta is moved
	// out to make the testing of time calculation independent of
	// the function actually making a request.
	start := time.Now()
	request, err := http.Get(url)
	end := time.Now()

	if err != nil{
		fmt.Println(err)
		return Response{} //for now, just return an empty response.
	}

	rt := CalculateMSDelta(start, end)

	return Response{StatusCode: request.StatusCode,
		RequestStart:    start,
		RequestFinished: end,
		ResponseTime:    rt,
	}
}

func MakeConcurrentRequests(url string, count int) []Response {
	// MakeConcurrentRequests is, for all intents and purposes,
	// a wrapper function that facilitates concurrent execution
	// of MakeRequest. Rather than force MakeRequest to handle requests,
	// reporting and concurrency, I'm pulling the concurrency out into
	// something separate. Single responsibility and all that.

	var responses []Response
	resultChannel := make(chan Response)

	// Do the concurrent stuff.
	for i := 0; i < count; i++ {
		//request a lock: good to go.
		//if good, acquire.
		go func(i int) {
			resultChannel <- MakeRequest(url)
			//release the lock: will allow another in.
		}(i)
	}

	// You've done the speedy stuff. Now unpack the channel
	// and place into a slice of Responses.
	for i := 0; i < count; i++ {
		result := <-resultChannel
		responses = append(responses, result)
	}

	return responses
}

func CalculateMSDelta(start time.Time, end time.Time) (ResponseTime int64) {
	// CalculateMSDelta, as it suggests, takes two timestamps and
	// calculates the delta between them by subtracting the start
	// from the end. It represents the final result in milliseconds.
	return end.Sub(start).Milliseconds()
}
