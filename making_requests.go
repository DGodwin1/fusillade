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

func MakeRequest(url string) Response {
	// MakeRequest takes a URL and returns a Response containing
	// all of the necessary information. CalculateMSDelta is moved
	// out to make the testing of time calculation independent of
	// the function actually making a request.
	start := time.Now()
	request, err := http.Get(url)
	end := time.Now()

	if err != nil {
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
	// MakeConcurrentRequests makes requests in a constant
	// fashion. It uses a ticker that, at present, is hard coded
	// to send a request every 100 milliseconds.

	var responses []Response
	resultChannel := make(chan Response)

	// Setup a new ticker that ticks every 100 milliseconds.
	ticker := time.NewTicker(100 * time.Millisecond)
	requestsSent := 0
	// Send a request every 100 milliseconds.
	for range ticker.C {
		if requestsSent == count {
			break
		}
		go func() {
			requestsSent++
			// MakeRequest might instead look at a WalkJourney() function that takes in a slice of URLs
			// that are then visited by it.
			resultChannel <- MakeRequest(url) //TODO: could just pass in _any_ function that is then called whose result is shoved into channel.
		}()
	}

	// You've done the speedy stuff. Now unpack it and return.
	for i := 0; i < count; i++ {
		result := <-resultChannel
		responses = append(responses, result)
	}

	return responses
}

type UserJourneyResult struct {
	Responses             map[int]Response
	Codes                 map[int]int
	JourneyStart          time.Time
	JourneyEnd            time.Time
	JourneyResponseTimeMS int64
}

func WalkJourney(urls []string) UserJourneyResult {
	// WalkJourney goes through a user journey (essentially a list of URLs)
	// and reports back how that user journey went.

	var Codes = map[int]int{}
	// Responses can store the numerical ID of the request that has been sent
	// and store the Response data if it wants it.
	var Responses = map[int]Response{}

	// Loop through each URL and add the response code to
	// the UserJourney struct.
	for i, u := range urls {
		// Make the request
		r := MakeRequest(u)

		// Add the status code from the request
		// we have just sent to the results struct.
		Codes[r.StatusCode] += 1
		Responses[i] = r
	}

	// You've got a load of different data points, now add them to the results
	StartTime := Responses[0].RequestStart
	EndTime := Responses[len(urls)-1].RequestFinished
	MilliSecondDelta := CalculateMSDelta(StartTime, EndTime)

	return UserJourneyResult{Responses, Codes, StartTime, EndTime, MilliSecondDelta}

}

func CalculateMSDelta(start time.Time, end time.Time) (ResponseTime int64) {
	// CalculateMSDelta, as it suggests, takes two timestamps and
	// calculates the delta between them by subtracting the start
	// from the end. It represents the final result in milliseconds.
	return end.Sub(start).Milliseconds()
}
