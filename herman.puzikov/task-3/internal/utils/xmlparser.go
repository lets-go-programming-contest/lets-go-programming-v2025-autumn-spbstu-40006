package utils

import (
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

func ParseXML(filepath string) (*ExchangeRate, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

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
