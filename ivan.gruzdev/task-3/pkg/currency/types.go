package currency

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type ValCurs struct {
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	NumCode  int           `xml:"NumCode" json:"num_code"`
	CharCode string        `xml:"CharCode" json:"char_code"`
	Value    CurrencyValue `xml:"Value" json:"value"`
}

type CurrencyValue float64

func (currencyValue *CurrencyValue) UnmarshalXML(decoder *xml.Decoder, startElement xml.StartElement) error {
	var str string

	err := decoder.DecodeElement(&str, &startElement)
	if err != nil {
		return err
	}

	str = strings.Replace(str, ",", ".", 1)
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}

	*currencyValue = CurrencyValue(value)

	return nil
}
