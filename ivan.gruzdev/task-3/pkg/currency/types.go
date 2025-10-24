package currency

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type CurrencyValue float64

type ValCurs struct {
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	NumCode  int           `xml:"NumCode"  json:"num_code"`
	CharCode string        `xml:"CharCode" json:"char_code"`
	Value    CurrencyValue `xml:"Value"    json:"value"`
}

func (currencyValue *CurrencyValue) UnmarshalXML(decoder *xml.Decoder, startElement xml.StartElement) error {
	var str string

	err := decoder.DecodeElement(&str, &startElement)
	if err != nil {
		return fmt.Errorf("decode element: %w", err)
	}

	str = strings.Replace(str, ",", ".", 1)

	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("parse float: %w", err)
	}

	*currencyValue = CurrencyValue(value)

	return nil
}
