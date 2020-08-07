package main

import (
	"testing"
)

func GetValidator() Validator {
	return ConfigValidator{}
}

func AssertError(t *testing.T, err error, file *Config){
	t.Helper()
	if err == nil{
		t.Errorf("wanted error but didn't get one with %v", file)
	}
}

func TestValidator(t *testing.T){
	t.Run("No URLs in JSON throws error", func(t *testing.T) {
		file, _ := ParseConfigFile("empty_array_fail.json")
		v := GetValidator()

		_, err := v.Validate(file)

		AssertError(t, err, file)

	})

	t.Run("An empty string in the URL array throws error", func(t *testing.T) {
		file, _ := ParseConfigFile("empty_string_url_0_hits.json")
		v := GetValidator()

		_, err := v.Validate(file)

		AssertError(t, err, file)
	})

	t.Run("Count must be greater than 0", func(t *testing.T) {
		file, _ := ParseConfigFile("count0.json")
		v := GetValidator()

		_, err := v.Validate(file)

		AssertError(t, err, file)
	})

	t.Run("URLs must be full", func(t *testing.T) {
		

	})
}
