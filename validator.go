package main

import (
	"errors"
	"net/url"
)

type Validator interface {
	Validate(*Config) (bool, error)
}

type ConfigValidator struct{}

func (ConfigValidator) Validate(c *Config) (bool, error) {
	// Validate takes a configuration file
	// and checks that the values held on the different
	// struct tags are appropriate for a given load test.

	// No URLs?
	if len(c.Urls) < 1 {
		return false, errors.New("you need to request at least 1 URL")
	}

	// You're trying to test at least one user journey, right?
	if c.Count < 1 {
		return false, errors.New("Count needs to be bigger than larger than 0")
	}

	// Properly formatted according to Go?
	for _, v := range c.Urls {
		_, err := url.ParseRequestURI(v)
		if err != nil {
			return false, err
		}
	}

	// Starts properly.
	for _, v := range c.Urls {
		_, err := StartsWithHTTP(v)
		if err != nil {
			return false, err
		}
	}

	// Rate bigger than 0
	if c.Rate < 1 {
		return false, errors.New("Rate must be greater than 0")
	}

	// PauseLength bigger than 0
	if c.PauseLength < 1 {
		return false, errors.New("PauseLength must be greater than 0")
	}

	return true, nil
}

func StartsWithHTTP(s string) (bool, error) {
	if len(s) < 4 {
		//Not enough letters to be http to begin with
		return false, errors.New("doesn't start with http")
	}

	if !(s[0:4] == "http") {
		return false, errors.New("doesn't start with http")
	}

	return true, nil
}
