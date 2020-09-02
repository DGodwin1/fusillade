package main

import (
	"reflect"
	"testing"
)

func TestGetParser(t *testing.T) {
	t.Run("JSON file returns JSON parser", func(t *testing.T) {
		got, _ := GetParser(".json")
		want := JSONParser{}
		if got != want {
			t.Errorf("Got %v, want %v", got, want)
		}
	})

	t.Run("Test that unsupported file format errors", func(t *testing.T) {
		_, err := GetParser(".jpeg")
		if err == nil {
			t.Errorf("wanted error but didn't get one")
		}
	})
}

func TestTranslate(t *testing.T) {
	t.Run("JSON file returns correct config", func(t *testing.T) {
		got, _ := ParseConfigFile("test_config.json")

		want := &Config{
			Urls:        []string{"https://www.google.com", "https://www.voguebusiness.com"},
			Count:       10,
			Rate:        100,
			PauseLength: 100,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got %v, want %v", got, want)
		}
	})
}
