package main

import (
	"time"
)


	func main(){
		//Get all the bits together so you can close over 'em.
		ticker := time.NewTicker(100*time.Millisecond)
		resultChannel := make(chan UserJourneyResult)
		count := 100

		DoConcurrentTask(func() {
			resultChannel <- WalkJourney([]string{"there should be some urls here"})
		}, count, *ticker)

		// You've done the speedy stuff, now unload from the channel.
		var responses []UserJourneyResult
		for i := 0; i < count; i++ {
			result := <-resultChannel
			responses = append(responses, result)
		}


}
