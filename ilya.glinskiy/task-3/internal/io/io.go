package io

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v3"
)

type Input struct {
	XMLName xml.Name    `xml:"ValCurs"`
	Valutes []XMLValute `xml:"Valute"`
}

type XMLValute struct {
	ID        string `xml:"ID,attr"`
	NumCode   int    `xml:"NumCode"`
	CharCode  string `xml:"CharCode"`
	Name      string `xml:"Name"`
	Nominal   string `xml:"Nominal"`
	Value     string `xml:"Value"`
	VunitRate string `xml:"VunitRate"`
}

type JSONValute struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func CharsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch charset {
	case "windows-1251":
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	default:
		return nil, fmt.Errorf("unknown charset: %s", charset)
	}
}

func ReadConfig(path string, config interface{}) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Couldn't open config file")
	}

	err = yaml.Unmarshal(content, config)
	if err != nil {
		return fmt.Errorf("Couldn't read config file")
	}

	return nil
}

func ReadInput(path string, input interface{}) error {
	inputFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Couldn't open input file")
	}
	defer inputFile.Close()

	decoder := xml.NewDecoder(inputFile)
	decoder.CharsetReader = CharsetReader

	err = decoder.Decode(input)
	if err != nil {
		return fmt.Errorf("Couldn't read input file")
	}

	return nil
}

func WriteOutput(path string, output []JSONValute) error {
	var pathParts []string = strings.Split(path, "/")
	err := os.MkdirAll(strings.Join(pathParts[0:len(pathParts)-1], "/"), 1)
	if err != nil {
		panic("Couldn't create directory for OutputFile")
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		panic("Couldn't marshal valutes into json")
	}

	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		panic("Couldn't write an output file")
	}

	return nil
}
