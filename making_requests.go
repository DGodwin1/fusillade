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

//TODO: pass in the dependency of time. Create a function that returns
// time.Now(). Under test, the function can return set timestamps where
// the latency score will be correct by virtue of it being 'tricked'. 
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
			resultChannel <- MakeRequest(url)
		}()
	}

	// You've done the speedy stuff. Now unpack it and return.
	for i := 0; i < count; i++ {
		result := <-resultChannel
		responses = append(responses, result)
	}

	return responses
}

func DoConcurrentUserJourney(url []string, count int) []UserJourneyResult {

	//So this should be passed in and changed.
	var results []UserJourneyResult
	resultChannel := make(chan UserJourneyResult)

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
			resultChannel <- WalkJourney(url)
		}()
	}

	// You've done the speedy stuff. Now unpack it and return.
	for i := 0; i < count; i++ {
		result := <-resultChannel
		results = append(results, result)
	}

	return results
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
		r := MakeRequest(u)
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
	// DoConcurrentTask takes in a function and run it concurrently for a set number of ticks
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
