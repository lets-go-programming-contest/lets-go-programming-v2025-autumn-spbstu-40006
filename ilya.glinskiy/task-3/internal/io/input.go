package io

import (
	"encoding/xml"
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
		return input, nil
	}
}

func ReadInput(path string, input *Input) (err error) {
	inputFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("couldn't open input file: %w", err)
	}

	defer func() {
		closeErr := inputFile.Close()
		if closeErr != nil {
			err = fmt.Errorf("couldn't close input file: %w", closeErr)
		}
	}()

	decoder := xml.NewDecoder(inputFile)
	decoder.CharsetReader = CharsetReader

	err = decoder.Decode(input)
	if err != nil {
		return fmt.Errorf("couldn't decode input file: %w", err)
	}

	return err
}
