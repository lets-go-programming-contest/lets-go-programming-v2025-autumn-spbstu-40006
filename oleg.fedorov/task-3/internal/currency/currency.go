package currency

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type Currency struct {
	XMLName  xml.Name `xml:"Valute"  json:"-"`
	ID       string   `xml:"ID,attr" json:"-"`
	NumCode  int      `xml:"NumCode" json:"num_code"`
	CharCode string   `xml:"CharCode" json:"char_code"`
	Nominal  int      `xml:"Nominal" json:"-"`
	Name     string   `xml:"Name"     json:"-"`
	Value    Float64  `xml:"Value"    json:"value"`
}

type Float64 float64

func (f *Float64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var valueString string
	if err := d.DecodeElement(&valueString, &start); err != nil {
		return fmt.Errorf("failed to decode element: %w", err)
	}

	valueString = strings.ReplaceAll(valueString, ",", ".")

	val, err := strconv.ParseFloat(valueString, 64)
	if err != nil {
		return fmt.Errorf("failed to parse currency value '%s': %w", valueString, err)
	}

	*f = Float64(val)
	return nil
}

type Currencies struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Currencies []Currency `xml:"Valute"`
}
