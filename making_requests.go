package fusillade

import (
	"net/http"
	"time"
)

type Response struct{
	StatusCode int
	//TODO: ID int
	//TODO: RequestStart
	//TODO: RequestFinished
	//TODO: RequestDelta
	ResponseTimeMS int64
}

func MakeRequest(url string) Response{
	start := time.Now()
	request, _ := http.Get(url)
	end := time.Since(start)

	return Response{StatusCode: request.StatusCode, ResponseTimeMS: end.Milliseconds()}
}
