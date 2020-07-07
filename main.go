package main

import (
	"fmt"
	"time"
)

func main() {
	var r = make(map[int]Response)
	// I should be able to make a request concurrently.
	// Let's start by just printing out the request when it completes.
	for i := 0; i < 100; i++ {
		go func(i int){
			r[i] = MakeRequest("https://www.google.com")
		}(i)
	}

	time.Sleep(3 * time.Second)

	for _, v := range r{
		fmt.Println(v.StatusCode, v.ResponseTime)
	}
}
