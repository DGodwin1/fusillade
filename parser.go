package main

// This file deals with all things config and parser.
import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	// Represents a parsed configuration file.
	// Regardless of how a parser is implemented
	// they should organise data into this format.
	Urls        []string `json:"urls"`
	Count       int      `json:"user_journey_amount"`
	Rate        int      `json:"rate"`
	PauseLength int      `json:"pause_length"`
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
		return &Config{}, err
	}
	defer file.Close()

	parser, err := GetParser(path)

	if err != nil {
		return &Config{}, err
	}

	FileInBytes, err := ioutil.ReadAll(file)

	if err != nil {
		return &Config{}, err
	}

	parsed := &Config{}
	// Do the parsing now
	parsed = parser.Translate(FileInBytes)

	return parsed, nil
}

func MakeFakeJSONFile(){
	fakeConfig, err := json.MarshalIndent(struct {
		Urls        []string `json:"urls"`
		Count       int      `json:"user_journey_amount"`
		Rate        int      `json:"rate"`
		PauseLength int      `json:"pause_length"`
	}{
		Urls:        []string{"https://www.example.com/", "https://www.example.com/"},
		Count:       100,
		Rate:        100,
		PauseLength: 100,
	}, " ", "	")

	if err != nil {
		log.Println(err)
		return
	}

	err = ioutil.WriteFile("config.json", fakeConfig, 0644)
}
