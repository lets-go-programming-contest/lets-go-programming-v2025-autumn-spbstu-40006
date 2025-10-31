package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
)

const (
	dirPerm  = 0o755
	filePerm = 0o600
)

type Data struct {
	ValCurs ValCurs `xml:"ValCurs"`
}

type ValCurs struct {
	Date   string   `xml:"Date,attr"`
	Name   string   `xml:"name,attr"`
	Valute []Valute `xml:"Valute"`
}

type Valute struct {
	ID       string `xml:"ID,attr"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

var (
	ErrBadConfig   = errors.New("config: input-file and output-file must be set")
	ErrEmptyValue  = errors.New("empty Value for currency")
	ErrZeroNominal = errors.New("nominal is zero for currency")
)

func main() {
	inputFile := "input.xml"
	outputFile := "output.xml"

	if inputFile == "" || outputFile == "" {
		panic(ErrBadConfig)
	}

	file, err := os.Open(inputFile)
	if err != nil {
		panic(fmt.Errorf("failed to open file: %w", err))
	}
	defer func() { _ = file.Close() }()

	data, err := parseXML(file)
	if err != nil {
		panic(fmt.Errorf("failed to parse XML: %w", err))
	}

	if err = writeXML(outputFile, data); err != nil {
		panic(fmt.Errorf("failed to write XML: %w", err))
	}

	return data, nil
}

func parseXML(file *os.File) (Data, error) {
	var data Data

	decoder := xml.NewDecoder(file)

	if err := decoder.Decode(&data); err != nil {
		return data, fmt.Errorf("decode XML: %w", err)
	}

	for _, val := range data.ValCurs.Valute {

		if val.Value == "" {
			return data, ErrEmptyValue
		}

		if val.Nominal == 0 {
			return data, ErrZeroNominal
		}
	}

	return data, fmt.Errorf("decode xml: %w", err)
}

func writeXML(filename string, data Data) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filePerm)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "    ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("encode XML: %w", err)
	}

	return nil
}
