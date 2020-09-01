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
		return false, errors.New("you need to make request at least 1 URL")
	}

	// Empty URLs?
	for _, v := range c.Urls {
		if v == "" {
			return false, errors.New("you can't test a url that doesn't exist")
		}
	}

	if c.Count < 1 {
		return false, errors.New("you need to test at least one user journey")
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

	return true, nil
}

func StartsWithHTTP(s string) (bool, error) {
	if len(s) < 4 {
		//Not enough letters to be http to begin with
		return false, errors.New("doesn't start with http")
	}
	return s[0:4] == "http", nil
}
