package parser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int           `json:"num_code"  xml:"NumCode"`
	CharCode string        `json:"char_code" xml:"CharCode"`
	Value    CurrencyValue `json:"value"     xml:"Value"`
}

type CurrencyValue float64

func (curr *CurrencyValue) UnmarshalXML(decod *xml.Decoder, start xml.StartElement) error {
	var str string
	if err := decod.DecodeElement(&str, &start); err != nil {
		return fmt.Errorf("failed decode element: %w", err)
	}

	strValue := strings.ReplaceAll(str, ",", ".")
	res, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return fmt.Errorf("failed to parse float: %w", err)
	}

	*curr = CurrencyValue(res)

	return nil
}

func ParseXML(path string) ([]Valute, error) {
	xmlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read xml file: %w", err)
	}

	var valCurs ValCurs

	decoder := xml.NewDecoder(bytes.NewReader(xmlFile))
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("failed to decode XML file: %w", err)
	}

	return valCurs.Valutes, nil
}
