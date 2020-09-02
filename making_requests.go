package main

import (
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	StatusCode      int
	RequestStart    time.Time
	RequestFinished time.Time
	ResponseTime    int64
}

func MakeRequest(url string, start func() time.Time, end func() time.Time) Response {
	// MakeRequest takes a URL and returns a Response.
	fmt.Printf("Making request to %q\n", url)
	s := start()
	request, err := http.Get(url)
	e := end()

	if err != nil {
		return Response{}
	}

	rt := calculateMSDelta(s, e)

	return Response{StatusCode: request.StatusCode,
		RequestStart:    s,
		RequestFinished: e,
		ResponseTime:    rt,
	}
}

type UserJourneyResult struct {
	// Collects responses and summative data
	// about a user journey
	Responses             map[int]Response
	Codes                 map[int]int
	JourneyStart          time.Time
	JourneyEnd            time.Time
	JourneyResponseTimeMS int64
	Finished              bool
}

func WalkJourney(urls []string, s SessionPauser) UserJourneyResult {
	// WalkJourney goes through URLs and reports back how it went.
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

		// Pause the session and pretend
		// to read the content for the amount
		// of time stored on the SessionPauser.
		s.PauseSession()

		// Should we request the next URL?
		if !StatusOkay(r.StatusCode) {
			break
		}
	}

	// You've completed the walk. Now collect the summative info.
	FirstURLStartTime := Responses[0].RequestStart
	LastURLEndTime := Responses[len(urls)-1].RequestFinished
	JourneyTime := CalculateJourneyDuration(Responses)

	// Did we walk the full journey?
	Finished = len(Responses) == len(urls)

	return UserJourneyResult{
		Responses,
		Codes,
		FirstURLStartTime,
		LastURLEndTime,
		JourneyTime,
		Finished}

}

func DoConcurrentTask(task func(), count int, ticker time.Ticker) {
	// DoConcurrentTask takes in a a task, fires off the code a set
	// It's used in this code for running many user journeys.
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

func CalculateJourneyDuration(Responses map[int]Response) int64 {
	var JourneyTime int64 = 0
	for _, v := range Responses {
		// We have to loop through through each response
		// and sum the response times. This is so we make sure that
		// when we talk of the user journey time, we don't
		// let any reading time figures affect the data.
		JourneyTime += v.ResponseTime
	}
	return JourneyTime
}

func calculateMSDelta(start time.Time, end time.Time) int64 {
	// CalculateMSDelta, as it suggests, takes two timestamps and
	// calculates the delta between them by subtracting the start
	// from the end. It's here to ensure that the calculation of deltas
	// is only managed in one place.
	return end.Sub(start).Milliseconds()
}

type SessionPauser interface {
	PauseSession()
}

type FakeUser struct {
	DelayTime int
}

func (f FakeUser) PauseSession() {
	time.Sleep(time.Duration(f.DelayTime) * time.Millisecond)
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
