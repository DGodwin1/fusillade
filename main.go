package main

import (
	"time"
)


	func main(){
		//Parse the file.


		//Now get prepare the test.
		ticker := time.NewTicker(100*time.Millisecond) //TODO: take from config.
		count := 100 //TODO: take from the config.
		resultChannel := make(chan UserJourneyResult)
		//var urls []string //TODO: take from config.

		DoConcurrentTask(func() {
			resultChannel <- WalkJourney([]string{"there should be some urls here"})
		}, count, *ticker)

		// You've done the speedy stuff, now unload from the channel.
		var responses []UserJourneyResult
		for i := 0; i < count; i++ {
			result := <-resultChannel
			responses = append(responses, result)
		}

		//Now prepare the data.


}
