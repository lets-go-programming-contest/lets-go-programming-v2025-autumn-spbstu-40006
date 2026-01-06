package io

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type Input struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type ValuteValue float64

type Valute struct {
	NumCode  int         `xml:"NumCode" json:"num_code"`
	CharCode string      `xml:"CharCode" json:"char_code"`
	Value    ValuteValue `xml:"Value" json:"value"`
}

func CharsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch charset {
	case "windows-1251":
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	default:
		return input, nil
	}
}

func (value *ValuteValue) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var str string

	err := decoder.DecodeElement(&str, &start)
	if err != nil {
		return fmt.Errorf("couldn't decode element: %w", err)
	}

	val, err := strconv.ParseFloat(strings.Replace(str, ",", ".", 1), 64)
	if err != nil {
		return fmt.Errorf("couldn't parse float: %w", err)
	}

	*value = ValuteValue(val)

	return nil
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
