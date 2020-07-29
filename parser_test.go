package main

import (
	"testing"
)

func TestGetParser(t *testing.T) {
	t.Run("JSON file returns JSON parser", func(t *testing.T) {
		got, _ := GetParser(".json")
		want := JSONParser{}
		if got != want{
			t.Errorf("Got %v, want %v", got, want)
		}
	})
	t.Run("XML file returns XML parser", func(t *testing.T) {
		got, _ := GetParser(".xml")
		want := XMLParser{}
		if got != want{
			t.Errorf("Got %v, want %v", got, want)
		}
	})

	t.Run("Test that unsupported file format errors", func(t *testing.T) {
		_, err := GetParser(".jpeg")
		if err == nil{
			t.Errorf("wanted error but didn't get one")
		}
	})
}

func TestTranslate(t *testing.T){}


