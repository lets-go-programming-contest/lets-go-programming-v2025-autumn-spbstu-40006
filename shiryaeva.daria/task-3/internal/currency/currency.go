package currency

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type CurrencyValue float64

func (cv *CurrencyValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	s = strings.Replace(s, ",", ".", 1)
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("parse currency value '%s': %w", s, err)
	}

	*cv = CurrencyValue(value)
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
