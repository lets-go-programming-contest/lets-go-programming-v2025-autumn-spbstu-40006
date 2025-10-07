package utils

import (
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

type ExchangeRate struct {
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	NumCode  string     `xml:"NumCode"`
	CharCode string     `xml:"CharCode"`
	Value    CommaFloat `xml:"Value"`
}

func ParseXML(filepath string) (*ExchangeRate, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel // for windows-1251 encoding of the xml

	var exchRates ExchangeRate
	if err := decoder.Decode(&exchRates); err != nil {
		return nil, fmt.Errorf("error decoding XML: %w", err)
	}

	return &exchRates, nil
}

func DescendingComparatorCurrency(a, b Currency) int {
	floatA, floatB := float64(a.Value), float64(b.Value)
	switch {
	case floatB < floatA:
		return -1
	case floatB > floatA:
		return 1
	default:
		return 0
	}
}
