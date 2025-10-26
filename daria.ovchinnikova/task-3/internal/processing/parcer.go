package processing

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type ValueFloat float64

func (value *ValueFloat) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var str string
	if err := decoder.DecodeElement(&str, &start); err != nil {
		return fmt.Errorf("decode value: %w", err)
	}

	valueFloat, err := strconv.ParseFloat(strings.Replace(str, ",", ".", 1), 64)
	if err != nil {
		return fmt.Errorf("parse float from: %w", err)
	}

	*value = ValueFloat(valueFloat)

	return nil
}

type Currency struct {
	NumCode  int        `json:"num_code"  xml:"NumCode"`
	CharCode string     `json:"char_code" xml:"CharCode"`
	Value    ValueFloat `json:"value"     xml:"Value"`
}

type Currencies struct {
	Items []Currency `xml:"Valute"`
}

func LoadXML(path string) ([]Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open XML file: %w", err)
	}

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch strings.ToLower(charset) {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return input, nil
		}
	}

	var data Currencies
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("decode XML: %w", err)
	}

	return data.Items, nil
}
