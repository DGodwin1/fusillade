package main

import (
	"fmt"
	"time"
)

func main(){
	// We'll store our various responses in this slice.
	var r []Response

	// Make the requests in sequence.
	start := time.Now().Hour()
	for i := 0; i<100; i++ {
		r = append(r, MakeRequest("https://www.google.com"))
	}
	end := time.Now().Hour()
	delta := end - start
	fmt.Println("Test length: %d. Test started: %d Test ended: %d.", delta, start, end)

	for _, v := range r{
		fmt.Printf("Code: %d. Time: %d.", v.StatusCode, v.ResponseTime)
		fmt.Println()
	}






	//r := MakeRequest("https://www.google.com")
	//fmt.Printf("Request start:%v.\nResponse time: %v.\nCode: %v\n", r.RequestStart, r.ResponseTime, r.StatusCode)
}