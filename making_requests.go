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

type Connection struct {}


func MakeRequest(url string) Response {
	// MakeRequest takes a URL and returns a Response containing
	// all of the necessary information. CalculateMSDelta is moved
	// out to make the testing of time calculation independent of
	// the function actually making a request.
	start := time.Now()
	request, err := http.Get(url)
	fmt.Println("DONE")
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
	// MakeConcurrentRequests makes requests in a constant
	// fashion. It uses a ticker that, at present, is hard coded
	// to send a request every 100 milliseconds.

	var responses []Response
	resultChannel := make(chan Response)

	// Setup a new ticker that ticks every 100 milliseconds.
	ticker := time.NewTicker(100*time.Millisecond)
	requestsSent := 0

	// Send a request every 100 milliseconds.
	go func(){
		for {
			select {
			case _ = <-ticker.C:
				if requestsSent == count{
					return
				}
				go func() {
					requestsSent++
					fmt.Println("Going to send request now...")
					// MakeRequest might instead look at a WalkJourney() function that takes in a slice of URLs
					// that are then visited by it.
					resultChannel <- MakeRequest(url) //TODO: could just pass in _any_ function that is then called whose result is shoved into channel.
				}()
			}
		}
	}()


	// You've done the speedy stuff. Now unpack it and return.
	for i := 0; i < count; i++ {
		result := <-resultChannel
		responses = append(responses, result)
	}

	return responses
}

type UserJourney struct{
	ResponseCodes map[int]int
}

func WalkJourney(urls []string) UserJourney{
	//make a map and add to the structure.
	var Responses = map[int]int{200:2}

	return UserJourney{
		Responses,
	}
}

func CalculateMSDelta(start time.Time, end time.Time) (ResponseTime int64) {
	// CalculateMSDelta, as it suggests, takes two timestamps and
	// calculates the delta between them by subtracting the start
	// from the end. It represents the final result in milliseconds.
	return end.Sub(start).Milliseconds()
}
