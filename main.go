package main

import (
	"fmt"
	"path/filepath"
)

// Represents a parsed configuration file.
// Regardless of how a parser is implemented, they
// all ought to be able to create a Config.
type Config struct{
	urls []string
	count int

}

// A tiny interface with one method (ParseFile()) that is used to
// to extract data from a file and create a Config object.
// There can, and ought, to be multiple implementations of Parser
// to suit different use cases — JSON, XML, CSV.
type Parser interface{
	Translate(file string) Config
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


// Consider creating a function that is called GetParser().
// GetParser() takes in the file name and returns the appropriate
// Parser object. When .ParseFile is run on the Parser object, then
// the decoding happens.

func Parse(path string) Config{
	// Translate takes in a filepath and returns
	// a parsed object for that file.
	var Parser = map[string]Parser{
		".xml":XMLParser{},
		".json":JSONParser{},
		//Add additional parsers here. They need to implement Parser interface.
	}
	return Parser[filepath.Ext(path)].Translate("path")

}



func main() {
	Parse("hello.xml")
	Parse("hello.json")
}
