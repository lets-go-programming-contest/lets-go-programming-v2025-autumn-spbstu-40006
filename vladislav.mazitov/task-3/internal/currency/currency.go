package currency

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type FloatType float64

type Currency struct {
	XMLName  xml.Name  `json:"-"         xml:"Valute"`
	ID       string    `json:"-"         xml:"ID,attr"`
	NumCode  int       `json:"num_code"  xml:"NumCode"`
	CharCode string    `json:"char_code" xml:"CharCode"`
	Nominal  int       `json:"-"         xml:"Nominal"`
	Name     string    `json:"-"         xml:"Name"`
	Value    FloatType `json:"value"     xml:"Value"`
}

type ValCurs struct {
	XMLName xml.Name   `xml:"ValCurs"`
	Valutes []Currency `xml:"Valute"`
}

func (float *FloatType) UnmarshalXML(decoder *xml.Decoder, s xml.StartElement) error {
	var value string

	err := decoder.DecodeElement(&value, &s)
	if err != nil {
		return fmt.Errorf("decode currency: %w", err)
	}

	value = strings.ReplaceAll(value, ",", ".")

	strToFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("parse float: %w", err)
	}

	*float = FloatType(strToFloat)

	return nil
}
