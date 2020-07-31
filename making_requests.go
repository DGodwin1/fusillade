package main

import (
	"net/http"
	"time"
)

type Response struct {
	StatusCode int
	//Store timings for a particular request.
	RequestStart    time.Time
	RequestFinished time.Time
	ResponseTime    int64
}

func MakeRequest(url string, start func() time.Time, end func() time.Time) Response {
	// MakeRequest takes a URL and returns a Response.
	s := start()
	request, err := http.Get(url)
	e := end()

	if err != nil {
		return Response{} //for now, just return an empty response.
	}

	rt := CalculateMSDelta(s, e)

	return Response{StatusCode: request.StatusCode,
		RequestStart:    s,
		RequestFinished: e,
		ResponseTime:    rt,
	}
}

type UserJourneyResult struct {
	Responses             map[int]Response
	Codes                 map[int]int
	JourneyStart          time.Time
	JourneyEnd            time.Time
	JourneyResponseTimeMS int64
	Finished              bool
}

func WalkJourney(urls []string) UserJourneyResult {
	// WalkJourney goes through a user journey (a slice of URLs)
	// and reports back how it went.

	var Codes = map[int]int{}
	var Finished bool

	// Responses stores an ID for a request (when it was sent)
	// and the response that request generated.
	var Responses = map[int]Response{}

	// Loop through the URLs, add the responses
	// and update the status code count.
	for i, u := range urls {
		r := MakeRequest(u, time.Now, time.Now)
		Responses[i] = r
		Codes[r.StatusCode] += 1

		// Should we request the next URL?
		if !StatusOkay(r.StatusCode) {
			break
		}
	}

	// You've completed the walk, now record how long that took (start, finish, delta)
	StartTime := Responses[0].RequestStart
	EndTime := Responses[len(urls)-1].RequestFinished
	MilliSecondDelta := CalculateMSDelta(StartTime, EndTime)

	// Did we walk the full journey?
	if len(Responses) == len(urls) {
		Finished = true
	} else {
		Finished = false
	}

	return UserJourneyResult{Responses, Codes, StartTime, EndTime, MilliSecondDelta, Finished}

}

func DoConcurrentTask(task func(), count int, ticker time.Ticker) {
	// DoConcurrentTask takes in a function (a task) and runs it, concurrently, a set number of times.
	TasksComplete := 0
	for range ticker.C {
		if TasksComplete == count {
			break
		}
		go func() {
			TasksComplete++
			task()
		}()
	}
}

func CalculateMSDelta(start time.Time, end time.Time) (ResponseTime int64) {
	// CalculateMSDelta, as it suggests, takes two timestamps and
	// calculates the delta between them by subtracting the start
	// from the end. It represents the final result in milliseconds.
	return end.Sub(start).Milliseconds()
}

func StatusOkay(status int) bool {
	// StatusOkay takes a status code and says whether
	// or not this is ok with respect to a user journey.

	switch {
	// 'Continue' https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/100
	case status >= 100 && status <= 199:
		return true
	// Response accepted
	case status >= 200 && status <= 299:
		return true
	// Redirects allowed
	case status >= 300 && status <= 399:
		return true
	// Client error
	case status >= 400 && status <= 499:
		return false
	// Server error
	case status >= 500:
		return false
	default:
		// Working on the assumption that if you didn't get
		// a response code then the request went wrong somewhere.
		return false
	}
}
