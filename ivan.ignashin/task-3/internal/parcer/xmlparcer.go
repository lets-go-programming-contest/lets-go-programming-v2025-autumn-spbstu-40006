package parcer

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type FloatValue float64

func (f *FloatValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string

	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	s = strings.ReplaceAll(s, ",", ".")

	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*f = FloatValue(value)
	return nil
}

type Record struct {
	ID    int        `json:"num_code"  xml:"NumCode"`
	Name  string     `json:"char_code" xml:"CharCode"`
	Value FloatValue `json:"value"     xml:"Value"`
}

type RawRecords struct {
	Items []Record `xml:"Valute"`
}

func ParseXML(path string) ([]Record, error) {
	xmlData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read xml file %s: %w", path, err)
	}

	var raw RawRecords
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	return raw.Items, nil
}
