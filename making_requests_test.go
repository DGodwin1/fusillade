package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func AssertResponseCode(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("Got %d, want %d", got, want)
	}
}

func TestGetter(t *testing.T) {
	t.Run("MakeRequest returns 200 on 'ok' URL", func(t *testing.T) {

		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		got := MakeRequest(FakeServer.URL)
		want := http.StatusOK

		AssertResponseCode(t, got.StatusCode, want)

	})

	t.Run("MakeRequest returns 500 on bad URL", func(t *testing.T) {
		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		got := MakeRequest(FakeServer.URL)
		want := http.StatusInternalServerError

		AssertResponseCode(t, got.StatusCode, want)
	})
}

func TestLatency(t *testing.T) {
	t.Run("Test latency checker is 10", func(t *testing.T) {
		start := time.Date(2019, 1, 1, 1, 1, 1, 0, time.UTC)
		finish := time.Date(2019, 1, 1, 1, 1, 1, 10000000, time.UTC)
		got := CalculateMSDelta(start, finish)
		var want int64 = 10
		if got != want {
			t.Errorf("got %d, wanted %d", got, want)
		}
	})

	t.Run("Test latency checker is 30", func(t *testing.T) {
		start := time.Date(2019, 1, 1, 1, 1, 1, 10000000, time.UTC)
		finish := time.Date(2019, 1, 1, 1, 1, 1, 40000000, time.UTC)
		got := CalculateMSDelta(start, finish)
		var want int64 = 30
		if got != want {
			t.Errorf("got %d, wanted %d", got, want)
		}
	})
}

func TestConcurrency(t *testing.T) {
	t.Run("Test that 20 requests leads to a collection of 20 Responses", func(t *testing.T) {
		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		got := MakeConcurrentRequests(FakeServer.URL, 20)
		want := 20

		if len(got) != want {
			t.Errorf("got %d, want %d", len(got), want)
		}

	})

	t.Run("Test that slow server responds to requests even though they are sent before it can respond", func(t *testing.T) {
		SlowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// The point of this is that we should expect to see a channel with 10 responses in it, even
			// though the requests have been sent to the server faster than the server is able to respond to all of them.
			time.Sleep(1 * time.Second)
			w.WriteHeader(http.StatusOK)
		}))

		got := MakeConcurrentRequests(SlowServer.URL, 10)
		want := 10

		if len(got) != want {
			t.Errorf("got %d, want %d", len(got), want)
		}
	})
}

func TestWalker(t *testing.T) {
	t.Run("Test that Walker struct has 2 200 response codes when given two good urls", func(t *testing.T) {
		// Setup the servers
		Server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		Server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		// Setup the URLS to hit.
		var URLS = []string{Server1.URL, Server2.URL}

		// 'work' will be a UserJourney struct that we can pull data out of.
		work := WalkJourney(URLS)

		// Reach into the ResponseCodes map that is stored in 'work' and see what's held there for 200 codes.
		got := work.Codes[200]
		want := 2

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
	t.Run("Test that Walker struct has 1 200 and 1 404 when it gets one good and one bad URL", func(t *testing.T) {
		//setup server.
		GoodServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		BadServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		//Setup the URLS to hit.
		var URLS = []string{GoodServer.URL, BadServer.URL}

		work := WalkJourney(URLS)
		got200 := work.Codes[200]
		got404 := work.Codes[404]
		want := 1

		if got200 != want || got404 != want {
			t.Errorf("got 200 of %d, 404 of %d. Both should be %d", got200, got404, want)
		}

	})

	t.Run("Test that Walker struct does not request further URLs after hitting a 404", func(t *testing.T) {
		//setup servers
		GoodServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		BadServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		var URLS = []string{GoodServer.URL, BadServer.URL, GoodServer.URL}

		// We should only get two pieces of response data in the struct.
		// One for the first GoodServer.URL and another for the BadServer.URL.
		work := WalkJourney(URLS)
		got := len(work.Responses)
		want := 2

		if got != want{
			t.Errorf("requested %d URLs. Should have requested %d", got, want)
		}

	})
}
