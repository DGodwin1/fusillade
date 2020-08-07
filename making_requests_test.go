package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func AssertResponseCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Got %d, want %d", got, want)
	}
}

func GetUserReader() EndUserReader {
	return EndUserReader{}
}

type Faker interface {
	Start() time.Time
	End() time.Time
}

type FakeTime struct{}

func (FakeTime) Start() time.Time {
	return time.Date(2019, 1, 1, 1, 1, 1, 0, time.UTC)
}

func (FakeTime) End() time.Time {
	return time.Date(2019, 1, 1, 1, 1, 1, 10000000, time.UTC)
}

func TestGetter(t *testing.T) {
	t.Run("MakeRequest latency response is 10MS", func(t *testing.T) {
		FakeServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
		}))

		faker := FakeTime{}

		got := MakeRequest(FakeServer.URL, faker.Start, faker.End)
		var want int64 = 10

		if got.ResponseTime != want {
			t.Errorf("Got %d, want %d", got.ResponseTime, want)
		}

	})

	t.Run("MakeRequest returns 200 on 'ok' URL", func(t *testing.T) {
		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		got := MakeRequest(FakeServer.URL, time.Now, time.Now)
		want := http.StatusOK

		AssertResponseCode(t, got.StatusCode, want)

	})

	t.Run("MakeRequest returns 500 on bad URL", func(t *testing.T) {
		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		got := MakeRequest(FakeServer.URL, time.Now, time.Now)
		want := http.StatusInternalServerError

		AssertResponseCode(t, got.StatusCode, want)
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
		URLS := []string{Server1.URL, Server2.URL}

		// 'work' will be a UserJourney struct that we can pull data out of.
		work := WalkJourney(URLS, GetUserReader())

		// Reach into the ResponseCodes map that is stored in 'work' and see what's held there for 200 codes.
		got := work.Codes[200]
		want := 2

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
	t.Run("Test that Walker struct has 1 200 and 1 404 when it gets one good and one bad URL", func(t *testing.T) {
		GoodServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		BadServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		URLS := []string{GoodServer.URL, BadServer.URL}

		work := WalkJourney(URLS, GetUserReader())
		got200 := work.Codes[200]
		got404 := work.Codes[404]
		want := 1

		if got200 != want || got404 != want {
			t.Errorf("got 200 of %d, 404 of %d. Both should be %d", got200, got404, want)
		}

	})

	t.Run("Test that Walker struct does not request further URLs after hitting a 404", func(t *testing.T) {
		GoodServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		BadServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		URLS := []string{GoodServer.URL, BadServer.URL, GoodServer.URL}

		// We should only get two responses in the struct:
		// one for the first GoodServer.URL and another for the BadServer.URL.
		// We should include the BadServer.URL response it was _part_ of the user's
		// journey. They would have experienced that page/that blocker so it should show
		// in the overall metrics. The final GoodServer.URL, on the other hand, should not.
		work := WalkJourney(URLS, GetUserReader())
		got := len(work.Responses)
		want := 2

		if got != want {
			t.Errorf("requested %d URLs. Should have requested %d", got, want)
		}

		// Do we have the correct status codes?
		got200 := work.Codes[200]
		got404 := work.Codes[404]
		CodeAmount := 1

		if got200 != CodeAmount || got404 != CodeAmount {
			t.Errorf("got 200 of %d, 404 of %d. Both should be %d", got200, got404, want)
		}

	})

	t.Run("User journey finished status is listed properly", func(t *testing.T) {
		GoodServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		BadServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		var FinishedTests = []struct {
			journey []string
			want    bool
		}{
			{[]string{GoodServer.URL}, true},
			{[]string{GoodServer.URL, GoodServer.URL}, true},
			{[]string{BadServer.URL}, true},
			{[]string{GoodServer.URL, BadServer.URL}, true},
			{[]string{BadServer.URL, BadServer.URL}, false},
			{[]string{BadServer.URL, GoodServer.URL}, false},
			{[]string{GoodServer.URL, BadServer.URL, GoodServer.URL}, false},
		}

		for _, tt := range FinishedTests {
			got := WalkJourney(tt.journey, GetUserReader())
			if got.Finished != tt.want {
				t.Errorf("Got %v, wanted %v", got.Finished, tt.want)
			}
		}

	})
}

func TestConcurrency(t *testing.T) {
	t.Run("20 requests leads to 20 Responses", func(t *testing.T) {
		FakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		ticker := time.NewTicker(1 * time.Millisecond)
		resultChannel := make(chan Response)
		count := 20

		DoConcurrentTask(func() {
			resultChannel <- MakeRequest(FakeServer.URL, time.Now, time.Now)
		}, count, *ticker)

		var responses []Response
		for i := 0; i < count; i++ {
			result := <-resultChannel
			responses = append(responses, result)
		}

		// There should be 20 responses from the channel
		got := len(responses)
		want := 20

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

	})

	t.Run("Test that slow server responds to requests even though they are sent before it can respond", func(t *testing.T) {
		SlowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// The point of this is that we should expect to see a channel with 10 responses in it, even
			// though the requests have been sent to the server faster than the server is able to respond to all of them.
			time.Sleep(1 * time.Second)
			w.WriteHeader(http.StatusOK)
		}))

		ticker := time.NewTicker(1 * time.Millisecond)
		resultChannel := make(chan Response)
		count := 10

		DoConcurrentTask(func() {
			resultChannel <- MakeRequest(SlowServer.URL, time.Now, time.Now)
		}, count, *ticker)

		var got []Response
		for i := 0; i < count; i++ {
			result := <-resultChannel
			got = append(got, result)
		}

		want := 10

		if len(got) != want {
			t.Errorf("got %d, want %d", len(got), want)
		}
	})
}

func TestStatusOkay(t *testing.T) {
	var StatusTests = []struct {
		c    int  //input
		want bool //ok?
	}{
		{0, false},
		{100, true},
		{199, true},
		{200, true},
		{201, true},
		{299, true},
		{300, true},
		{399, true},
		{400, false},
		{499, false},
		{500, false},
		{501, false},
		{599, false},
		{599, false},
	}
	for _, tt := range StatusTests {
		got := StatusOkay(tt.c)
		if got != tt.want {
			t.Errorf("got %v, want %v", got, tt.want)
		}
	}
}
