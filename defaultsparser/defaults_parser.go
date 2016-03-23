package defaultsparser

import (
	"github.com/dugancathal/stuffs/config"
	"io/ioutil"
"fmt"
	"bufio"
	"bytes"
	"path"
)


type DefaultsParser interface {
	RouteMapping() config.RouteMapping
}

type DirectoryDefaultsParser struct {
	directoryPath string
}

func NewDirectoryDefaultParser(path string) *DirectoryDefaultsParser {
	return &DirectoryDefaultsParser{directoryPath: path}
}

func (self *DirectoryDefaultsParser) RouteMapping() config.RouteMapping {
	routeMapping := make(config.RouteMapping)
	if self.directoryPath == "" {
		return routeMapping
	}

	files, err := ioutil.ReadDir(self.directoryPath)
	if err != nil {
		fmt.Printf("*** Received an error trying to parse route defaults: %s\n", err.Error())
		fmt.Printf("*** Falling back to no defaults", err.Error())
		return routeMapping
	}

	fmt.Printf("Parsing files in config directory: %s\n", self.directoryPath)
	for _, file := range files {
		if !file.IsDir() {
			fmt.Printf("-- Parsing file: %s\n", file.Name())
			fileContents, _ := ioutil.ReadFile(path.Join(self.directoryPath, file.Name()))
			scanner := bufio.NewReader(bytes.NewBuffer(fileContents))
			endpoint, _, _ := scanner.ReadLine()
			body := make([]byte, scanner.Buffered())
			scanner.Read(body)
			routeMapping[string(endpoint)] = body
		}
	}
	return routeMapping
}
