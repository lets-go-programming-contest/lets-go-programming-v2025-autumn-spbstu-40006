package utils

import (
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

type ExchangeRate struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Date       string     `xml:"Date,attr"`
	Name       string     `xml:"name,attr"`
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	ID        string     `xml:"ID"`
	NumCode   string     `xml:"NumCode"`
	CharCode  string     `xml:"CharCode"`
	Nominal   uint       `xml:"Nominal"`
	Name      string     `xml:"Name"`
	Value     CommaFloat `xml:"Value"`
	VunitRate CommaFloat `xml:"VunitRate"`
}

func ParseXML(filepath string) (*ExchangeRate, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel
	var exchRates ExchangeRate
	if err := decoder.Decode(&exchRates); err != nil {
		return nil, fmt.Errorf("error decoding XML: %w", err)
	}

	return &exchRates, nil
}
