package currency

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type Currency struct {
	XMLName  xml.Name `json:"-"         xml:"Valute"`
	ID       string   `json:"-"         xml:"ID,attr"`
	NumCode  int      `json:"num_code"  xml:"NumCode"`
	CharCode string   `json:"char_code" xml:"CharCode"`
	Nominal  int      `json:"-"         xml:"Nominal"`
	Name     string   `json:"-"         xml:"Name"`
	Value    Float64  `json:"value"     xml:"Value"`
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
