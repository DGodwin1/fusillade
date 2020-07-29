package main
// This file deals with all things config and parser.
import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct{
	// Represents a parsed configuration file.
	// Regardless of how a parser is implemented
	urls []string
	count int
}

type Parser interface{
	// All Parsers should have one method method, Translate() that
	// takes a filepath and creates a Config. There can and ought to be
	// multiple Parser implementations that can work with many
	// different file types.
	Translate(filepath string) Config
}

type JSONParser struct{}

func (j JSONParser) Translate(file string) Config{
	//
	//parsed := &Config{}
	//
	//err := json.Unmarshal(bytesData, parsed)
	//
	//if err != nil {
	//	fmt.Println("error", err)
	//}


	return Config{[]string{"JSON PARSED FILE"}, 2}
}

type XMLParser struct{}

func (x XMLParser) Translate(file string) Config{
	fmt.Println("XML parsing baby")
	return Config{[]string{"XML parsing"}, 5}
}



func GetParser(extension string) (Parser, error){
	Parser := map[string]Parser{
		".xml":XMLParser{},
		".json":JSONParser{},
		//Add additional parsers here. They need to implement Parser interface.
	}

	// Before giving a parser back, let's make sure we actually have one to offer
	p, ok := Parser[extension]
	if !ok{
		return nil, errors.New("file format is not parsable")
	}
	return p, nil
}

func ParseConfigFile(path string) Config{
	return GetParser(filepath.Ext(path)).Translate(path)
}


