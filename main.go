package main

import (
	"fmt"
	"time"
)

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

	func main(){
		Get all the bits together so you can close over 'em.
		resultChannel := make(chan Response)
		count := 100


		DoConcurrentTask(func() {
			resultChannel <- MakeRequest("https://www.google.com")
		}, count, *ticker)


		// You've done the speedy stuff, now pull stuff out.
		var responses []Response

		for i := 0; i < count; i++ {
			result := <-resultChannel
			responses = append(responses, result)
		}

		for _, v := range responses{
			fmt.Println(v.ResponseTime)
		}


		ticker := time.NewTicker(100 * time.Millisecond)
		var results []UserJourneyResult
		resultChannel := make(chan UserJourneyResult)
		count := 10

		var urls = []string{"https://www.google.com", "https://wwww.tatler.com"}

		DoConcurrentTask(func(){
					resultChannel <- WalkJourney(urls)},
					count, *ticker)


		for i := 0; i < count; i++ {
			result := <-resultChannel
			results = append(results, result)
		}

		fmt.Println("Done unpacking")

		// Now let's go through all of the codes you received in the process.
		for _, v := range results{
			fmt.Println(v.Codes)
		}
}
