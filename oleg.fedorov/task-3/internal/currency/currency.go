package currency

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type Currency struct {
	XMLName  xml.Name `xml:"Valute" json:"-"`
	ID       string   `xml:"ID,attr" json:"-"`
	NumCode  int      `xml:"NumCode" json:"num_code"`
	CharCode string   `xml:"CharCode" json:"char_code"`
	Nominal  int      `xml:"Nominal" json:"-"`
	Name     string   `xml:"Name" json:"-"`
	Value    Float64  `xml:"Value" json:"value"`
}

type Float64 float64

type Currencies struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Currencies []Currency `xml:"Valute"`
}

func (f *Float64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string

	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	s = strings.Replace(s, ",", ".", -1)

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("failed to parse currency value '%s': %w", s, err)
	}

	*f = Float64(val)

	return nil
}
