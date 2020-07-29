package main
// This file deals with all things config and parser.
import (
	"fmt"
	"path/filepath"
)


// Represents a parsed configuration file.
// Regardless of how a parser is implemented, they
type Config struct{
	urls []string
	count int
}

// A tiny interface. All Parsers should have one method method, Translate() that
// parses a file and create a Config.

// There can, and ought, to be multiple implementations of Parser
// that can, if they implement the Translate said to satisfy the Parser interface — JSON, XML, CSV.
type Parser interface{
	Translate(filepath string) Config
}

type JSONParser struct{}

func (j JSONParser) Translate(file string) Config{
	fmt.Println("JSON parsing")
	return Config{[]string{"JSON PARSED FILE"}, 2}
}

type XMLParser struct{}

func (x XMLParser) Translate(file string) Config{
	fmt.Println("XML parsing baby")
	return Config{[]string{"XML parsing"}, 5}
}

func GetParser(extension string) Parser{
	Parser := map[string]Parser{
		".xml":XMLParser{},
		".json":JSONParser{},
		//Add additional parsers here. They need to implement Parser interface.
	}
	return Parser[extension]
}


func ParseConfigFile(path string) Config{
	return GetParser(filepath.Ext(path)).Translate(path)
}


