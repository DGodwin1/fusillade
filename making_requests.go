package fusillade

import "net/http"

func MakeRequest(url string) int{

	resp, _ := http.Get(url)

	return resp.StatusCode
}
