package main

import (
	//"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	//"time"
)

//func AssertResponseCode(t *testing.T, got, want int) {
//	if got != want {
//		t.Errorf("Got %d, want %d", got, want)
//	}
//}

//func TestGetter(t *testing.T) {
//	t.Run("MakeRequest returns 200 on 'ok' URL", func(t *testing.T) {
//
//		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			w.WriteHeader(http.StatusOK)
//		}))
//
//		got := MakeRequest(FakeServer.URL)
//		want := http.StatusOK
//
//		AssertResponseCode(t, got.StatusCode, want)
//
//	})
//
//	t.Run("MakeRequest returns 500 on bad URL", func(t *testing.T) {
//		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			w.WriteHeader(http.StatusInternalServerError)
//		}))
//
//		got := MakeRequest(FakeServer.URL)
//		want := http.StatusInternalServerError
//
//		AssertResponseCode(t, got.StatusCode, want)
//	})
//}
//
//func TestLatency(t *testing.T) {
//	t.Run("Test latency checker is 10", func(t *testing.T) {
//		start := time.Date(2019, 1, 1, 1, 1, 1, 0, time.UTC)
//		finish := time.Date(2019, 1, 1, 1, 1, 1, 10000000, time.UTC)
//		got := CalculateMSDelta(start, finish)
//		var want int64 = 10
//		if got != want {
//			t.Errorf("got %d, wanted %d", got, want)
//		}
//	})
//
//	t.Run("Test latency checker is 30", func(t *testing.T) {
//		start := time.Date(2019, 1, 1, 1, 1, 1, 10000000, time.UTC)
//		finish := time.Date(2019, 1, 1, 1, 1, 1, 40000000, time.UTC)
//		got := CalculateMSDelta(start, finish)
//		var want int64 = 30
//		if got != want {
//			t.Errorf("got %d, wanted %d", got, want)
//		}
//	})
//}
//
//func TestConcurrency(t *testing.T) {
//	t.Run("Test that 100 requests leads to a collection of 100 Responses", func(t *testing.T) {
//		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			w.WriteHeader(http.StatusOK)
//		}))
//
//		got := MakeConcurrentRequests(FakeServer.URL, 100)
//		want := 100
//
//		if len(got) != want {
//			t.Errorf("got %d, want %d", len(got), want)
//		}
//
//	})
//
//	t.Run("Test that 200 requests leads to a collection of 200 Responses", func(t *testing.T) {
//		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			w.WriteHeader(http.StatusOK)
//		}))
//
//		got := MakeConcurrentRequests(FakeServer.URL, 200)
//		want := 200
//
//		if len(got) != want {
//			t.Errorf("got %d, want %d", len(got), want)
//		}
//
//	})
//
//	t.Run("Test that server responds ", func(t *testing.T){
//		SlowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			// The point of this is that we should expect to see a channel with 10 responses in it, even
//			// though the requests have been sent to the server faster than the server is able to respond to all of them.
//			time.Sleep(5*time.Second)
//			w.WriteHeader(http.StatusOK)
//		}))
//
//		got := MakeConcurrentRequests(SlowServer.URL, 10)
//		want := 10
//
//		if len(got) != want {
//			t.Errorf("got %d, want %d", len(got), want)
//		}
//	})
//}

func TestWalker(t *testing.T){
	t.Run("Test that Walker struct has 2 200 response codes when given two good urls", func(t *testing.T) {
		//setup server.
		Server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		Server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		//setup list of URLS to hit.
		var URLS = []string{Server1.URL, Server2.URL}

		//work will return a UserJourney struct.
		work := WalkJourney(URLS)

		//we should be able to look at the ResponseCodes component of the struct, which is itself a map of response codes to maps.
		got := work.ResponseCodes[200]
		want := 2

		if got != want{
			t.Errorf("got %d, want %d", got, want)
		}


		//give URLs to function
		//say that want should be 2
		//check that it passes

	})
}
