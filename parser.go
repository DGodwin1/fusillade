package main

// This file deals with all things config and parser.
import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	// Represents a parsed configuration file.
	// Regardless of how a parser is implemented
	Urls  []string `json:"urls"`
	Count int      `json:"hits"`
}

type Parser interface {
	// All Parsers should have one method method, Translate() that
	// takes a filepath and creates a Config. There can and ought to be
	// multiple Parser implementations that can work with many
	// different file types.
	Translate(bytesData []byte) *Config
}

type JSONParser struct{}

func (JSONParser) Translate(bytesData []byte) *Config {

	parsed := &Config{}

	err := json.Unmarshal(bytesData, parsed)

	if err != nil {
		fmt.Println("error", err)
	}

	return parsed
}

func GetParser(path string) (Parser, error) {
	Parser := map[string]Parser{
		".json": JSONParser{},
		//Add additional parsers here. They need to implement Parser interface.
	}

	// Before giving a parser back, let's make sure we actually have one to offer
	p, ok := Parser[filepath.Ext(path)]
	if !ok {
		return nil, errors.New("there isn't a parser available for this type of file")
	}
	return p, nil
}

func ParseConfigFile(path string) (*Config, error) {
	// Get the file
	file, err := os.Open(path)

	if err != nil {
		fmt.Println(err)

	}
	defer file.Close()

	parser, _ := GetParser(path)

	FileInBytes, _ := ioutil.ReadAll(file)

	parsed := &Config{}
	// Do the parsing now
	parsed = parser.Translate(FileInBytes)

	return parsed, nil
}
