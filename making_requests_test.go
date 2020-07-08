package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func AssertResponseCode(t *testing.T, got, want int){
	if got != want{
		t.Errorf("Got %d, want %d", got, want)
	}
}

func TestGetter(t *testing.T){
	t.Run("MakeRequest returns 200 on 'ok' URL", func(t *testing.T){

		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			w.WriteHeader(http.StatusOK)
		}))

		got := MakeRequest(FakeServer.URL)
		want := http.StatusOK

		AssertResponseCode(t, got.StatusCode, want)

	})

	t.Run("MakeRequest returns 500 on bad URL", func(t *testing.T){
		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			w.WriteHeader(http.StatusInternalServerError)
		}))

		got := MakeRequest(FakeServer.URL)
		want := http.StatusInternalServerError

		AssertResponseCode(t, got.StatusCode, want)
	})
}

func TestLatency(t *testing.T){
	t.Run("Test latency checker is 10", func(t *testing.T) {
		start := time.Date(2019, 1, 1, 1, 1, 1, 0, time.UTC)
		finish := time.Date(2019, 1, 1, 1, 1, 1, 10000000, time.UTC)
		got := CalculateMSDelta(start, finish)
		var want int64 = 10
		if got != want{
			t.Errorf("got %d, wanted %d", got, want)
		}
	})

	t.Run("Test latency checker is 30", func(t *testing.T) {
		start := time.Date(2019, 1, 1, 1, 1, 1, 10000000, time.UTC)
		finish := time.Date(2019, 1, 1, 1, 1, 1, 40000000, time.UTC)
		got := CalculateMSDelta(start, finish)
		var want int64 = 30
		if got != want{
			t.Errorf("got %d, wanted %d", got, want)
		}
	})
}

func TestConcurrency(t *testing.T){
	t.Run("Test that 100 requests leads to a collection of 100 Responses", func(t *testing.T){
		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			w.WriteHeader(http.StatusOK)
		}))

		got := MakeConcurrentRequests(FakeServer.URL)
		want := 100

		if len(got) != want{
			t.Errorf("got %d, want %d", len(got),want)
		}

	})

}