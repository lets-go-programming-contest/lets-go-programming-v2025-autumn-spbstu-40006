package currency

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type CurrencyValue float64

func (cv *CurrencyValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string

	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\u00A0", "")
	s = strings.ReplaceAll(s, ",", ".")

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*cv = CurrencyValue(f)
	return nil
}

type ValCurs struct {
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	NumCode  int           `json:"num_code"  xml:"NumCode"`
	CharCode string        `json:"char_code" xml:"CharCode"`
	Value    CurrencyValue `json:"value"     xml:"Value"`
}
