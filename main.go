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
		//Get all the bits together so you can close over 'em.
		ticker := time.NewTicker(100*time.Millisecond)
		resultChannel := make(chan UserJourneyResult)
		count := 100

		DoConcurrentTask(func() {
			resultChannel <- WalkJourney([]string{"there should be some urls here"})
		}, count, *ticker)

		// You've done the speedy stuff, now pull stuff out.
		var responses []UserJourneyResult

		for i := 0; i < count; i++ {
			result := <-resultChannel
			responses = append(responses, result)
		}

		for _, v := range responses{
			fmt.Println(v.Codes)
		}

}
