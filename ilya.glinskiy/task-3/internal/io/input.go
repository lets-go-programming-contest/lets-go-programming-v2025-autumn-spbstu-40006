package io

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/charmap"
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

func CharsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch charset {
	case "windows-1251":
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	default:
		return nil, errors.New("unknown charsets")
	}
}

func ReadInput(path string, input interface{}) error {
	inputFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Couldn't open input file")
	}
	defer func() {
		err = inputFile.Close()
		if err != nil {
			panic("Couldn't close input file")
		}
	}()

	decoder := xml.NewDecoder(inputFile)
	decoder.CharsetReader = CharsetReader

	err = decoder.Decode(input)
	if err != nil {
		return fmt.Errorf("Couldn't read input file")
	}

	return nil
}
