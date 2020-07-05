package fusillade

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

		AssertResponseCode(t, got, want)

	})

	t.Run("MakeRequest returns 500 on bad URL", func(t *testing.T){
		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			w.WriteHeader(http.StatusInternalServerError)
		}))

		got := MakeRequest(FakeServer.URL)
		want := http.StatusInternalServerError

		AssertResponseCode(t, got, want)
	})
}