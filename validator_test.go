package main

import (
	"errors"
	"testing"
)

func GetValidator() Validator {
	return ConfigValidator{}
}

func AssertError(t *testing.T, err error, c *Config) {
	t.Helper()
	if err == nil {
		t.Errorf("wanted error but didn't get one with %v", c)
	}
}

func TestValidator(t *testing.T) {
	t.Run("No URLs in JSON throws error", func(t *testing.T) {
		c := &Config{
			Urls:  nil,
			Count: 100,
			Rate:  1,
		}
		v := GetValidator()

		_, err := v.Validate(c)

		AssertError(t, err, c)

	})

	t.Run("Count must be greater than 0", func(t *testing.T) {
		c := &Config{
			Urls:  []string{"https://"},
			Count: 0,
			Rate:  2,
		}

		v := GetValidator()

		_, err := v.Validate(c)

		AssertError(t, err, c)
	})

	t.Run("URLs must start with http'", func(t *testing.T) {
		c := &Config{
			Urls:  []string{"httz://www.google.com"},
			Count: 1,
			Rate:  1,
		}

		v := GetValidator()

		_, err := v.Validate(c)

		AssertError(t, err, c)
	})

	t.Run("Rate must be greater than 0", func(t *testing.T) {
		c := &Config{
			Urls:  []string{"https://www.googlg.com"},
			Count: 100,
			Rate:  0,
		}

		v := GetValidator()

		_, err := v.Validate(c)

		AssertError(t, err, c)
	})

	t.Run("Session pause must be greater than 0", func(t *testing.T) {
		c := &Config{
			Urls:        []string{"https://www.googlg.com"},
			Count:       100,
			Rate:        1,
			PauseLength: 0,
		}
		v := GetValidator()

		_, err := v.Validate(c)

		AssertError(t, err, c)
	})
}

func TestStartsWithHTTP(t *testing.T) {
	m := errors.New("doesn't start with http")
	var ProtocolTests = []struct {
		s string
		b bool
		e error
	}{
		{"https://", true, nil},
		{"http", true, nil},
		{"h", false, m},
		{"ftp", false, m},
		{"telnet", false, m},
		{"192.168.0", false, m},
	}
	for _, tt := range ProtocolTests {
		got, _ := StartsWithHTTP(tt.s)
		want := tt.b
		if got != want {
			t.Errorf("Got %t, wanted %t with %q", got, want, tt.s)
		}
	}
}
