package currency

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type CurrencyValue float64

func (cv *CurrencyValue) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var stringValue string
	if decodeErr := decoder.DecodeElement(&stringValue, &start); decodeErr != nil {
		return fmt.Errorf("failed to decode XML element: %w", decodeErr)
	}

	stringValue = strings.Replace(stringValue, ",", ".", 1)
	parsedValue, parseErr := strconv.ParseFloat(stringValue, 64)

	if parseErr != nil {
		return fmt.Errorf("parse currency value '%s': %w", stringValue, parseErr)
	}

	*cv = CurrencyValue(parsedValue)

	return nil
}

type Currency struct {
	XMLName  xml.Name      `xml:"Valute"`
	NumCode  int           `xml:"NumCode"`
	CharCode string        `xml:"CharCode"`
	Value    CurrencyValue `xml:"Value"`
}

type ValCurs struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Currencies []Currency `xml:"Valute"`
}

type JSONCurrency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func (c Currency) ToJSON() JSONCurrency {
	return JSONCurrency{
		NumCode:  c.NumCode,
		CharCode: c.CharCode,
		Value:    float64(c.Value),
	}
}
