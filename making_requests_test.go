package fusillade

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

	t.Run("Response time for request is 10ms", func(t *testing.T){
		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			time.Sleep(10 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))

		got := MakeRequest(FakeServer.URL)
		var want int64 = 10

		if got.ResponseTimeMS != want{
			t.Errorf("Got response time %d. Expected response time of %d", got.ResponseTimeMS, want)
		}

	})

	t.Run("Response time for request is 20ms", func(t *testing.T){
		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			time.Sleep(20 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))

		got := MakeRequest(FakeServer.URL)
		var want int64 = 20

		if got.ResponseTimeMS != want{
			t.Errorf("Got response time %d. Expected response time of %d", got.ResponseTimeMS, want)
		}

	})
}