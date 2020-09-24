package main

import (
	"reflect"
	"testing"
)

func TestMaxLatencyValue(t *testing.T) {
	t.Run("Test max with a single 'max' number", func(t *testing.T) {
		var ujs []UserJourneyResult

		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 50})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 60})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 70})

		got := MaxUserJourneyResponseLatency(ujs)
		var want = 70

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("Test max with several 'max' numbers", func(t *testing.T) {
		var ujs []UserJourneyResult

		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 70})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 60})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 70})

		got := MaxUserJourneyResponseLatency(ujs)
		var want = 70

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("Test max with minus values", func(t *testing.T) {
		// This is unlikely to play out in production (you can't have minus time)
		// but it's good to know the system won't fall over if it gets a
		// non-natural number.
		var ujs []UserJourneyResult

		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: -20})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: -30})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: -40})

		got := MaxUserJourneyResponseLatency(ujs)
		var want = -20

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("Test max with numbers across the divide", func(t *testing.T) {

		var ujs []UserJourneyResult

		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: -1})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 0})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 1})

		got := MaxUserJourneyResponseLatency(ujs)
		var want = 1

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

}

func TestFindMinJourneyResponseTime(t *testing.T) {
	t.Run("Find min with natural numbers", func(t *testing.T) {
		var ujs []UserJourneyResult

		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 55})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 20})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 64})

		got := MinUserJourneyResponseLatency(ujs)
		var want = 20

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("Find min with two min numbers", func(t *testing.T) {
		var ujs []UserJourneyResult

		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 0})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 0})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 1})

		got := MinUserJourneyResponseLatency(ujs)
		var want = 0

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("Find min with negative numbers", func(t *testing.T) {
		var ujs []UserJourneyResult

		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: -1})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: -2})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: -3})

		got := MinUserJourneyResponseLatency(ujs)
		var want = -3

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("Find min with numbers across the divide", func(t *testing.T) {
		var ujs []UserJourneyResult

		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 1})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: 0})
		ujs = append(ujs, UserJourneyResult{JourneyResponseTimeMS: -1})

		got := MinUserJourneyResponseLatency(ujs)
		var want = -1

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}

func TestCountResponseCodes(t *testing.T) {
	t.Run("Test correct counting", func(t *testing.T) {
		var ujs []UserJourneyResult
		codes := map[int]int{200: 2, 404: 1, 500: 5}
		ujs = append(ujs, UserJourneyResult{Codes: codes})
		ujs = append(ujs, UserJourneyResult{Codes: codes})
		ujs = append(ujs, UserJourneyResult{Codes: codes})

		got := CountResponseCodes(ujs)
		want := map[int]int{200: 6, 404: 3, 500: 15}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("test only 200 responses", func(t *testing.T) {
		var twoHundreds []UserJourneyResult
		twos := map[int]int{200: 2}

		twoHundreds = append(twoHundreds, UserJourneyResult{Codes: twos})
		twoHundreds = append(twoHundreds, UserJourneyResult{Codes: twos})

		got := CountResponseCodes(twoHundreds)
		want := map[int]int{200: 4}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("test no responses", func(t *testing.T) {
		var twoHundreds []UserJourneyResult

		got := CountResponseCodes(twoHundreds)
		want := map[int]int{}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
