package main

import (
	"errors"
)

type Validator interface {
	Validate(*Config) (bool, error)
}

type ConfigValidator struct{}

func (ConfigValidator) Validate(c *Config) (bool, error){
	// Validate takes a configuration file
	// and checks that the values held on the different
	// struct tags are appropriate for a given load test.

	//No URLs
	if len(c.Urls)<1{
		return false, errors.New("you need to request at least 1 URL")
	}

	// Any empty strings anywhere?
	for _, v := range c.Urls{
		if v == ""{
			return false, errors.New("you can't have an empty url")
		}
	}

	if c.Count < 1{
		return false, errors.New("count can't be lower than 1")
	}

	return true, nil
}

