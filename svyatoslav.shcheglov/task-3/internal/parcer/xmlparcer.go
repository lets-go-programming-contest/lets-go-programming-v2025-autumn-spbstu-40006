package parser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type DecimalNumber float64

func (dec *DecimalNumber) UnmarshalXML(decoder *xml.Decoder, openingTag xml.StartElement) error {
	var textContent string

	err := decoder.DecodeElement(&textContent, &openingTag)
	if err != nil {
		return fmt.Errorf("decode element content: %w", err)
	}

	normalizedText := strings.ReplaceAll(textContent, ",", ".")

	parsedNumber, err := strconv.ParseFloat(normalizedText, 64)
	if err != nil {
		return fmt.Errorf("convert text to number: %w", err)
	}

	*dec = DecimalNumber(parsedNumber)

	return nil
}

type CurrencyData struct {
	NumericCode int           `json:"num_code"  xml:"NumCode"`
	Symbol      string        `json:"char_code" xml:"CharCode"`
	Rate        DecimalNumber `json:"value"     xml:"Value"`
}

type CurrencyList struct {
	Currencies []CurrencyData `xml:"Valute"`
}

func LoadCurrencyData(filePath string) ([]CurrencyData, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read data file %s: %w", filePath, err)
	}

	var currencyContainer CurrencyList

	dataDecoder := xml.NewDecoder(bytes.NewReader(fileContent))
	dataDecoder.CharsetReader = charset.NewReaderLabel

	if err := dataDecoder.Decode(&currencyContainer); err != nil {
		return nil, fmt.Errorf("parse xml data: %w", err)
	}

	return currencyContainer.Currencies, nil
}
