package parser

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

const dirPerm = os.ModePerm

type ValCurs struct {
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int          `json:"num_code"  xml:"NumCode"`
	CharCode string       `json:"char_code" xml:"CharCode"`
	Value    Float64Comma `json:"value"     xml:"Value"`
}

type Float64Comma float64

func (f *Float64Comma) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var valueStr string
	if err := decoder.DecodeElement(&valueStr, &start); err != nil {
		return fmt.Errorf("failed to decode float value: %w", err)
	}

	value, err := strconv.ParseFloat(strings.ReplaceAll(valueStr, ",", "."), 64)
	if err != nil {
		return fmt.Errorf("invalid float format %q: %w", valueStr, err)
	}

	*f = Float64Comma(value)

	return nil
}

func (f *Float64Comma) Float64() float64 {
	return float64(*f)
}

func ParseXMLFile(path string) ([]Valute, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open XML file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "failed to close file %q: %v\n", path, cerr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return valCurs.Valutes, nil
}

func SaveToJSON(path string, valutes []Valute) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", dir, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", path, err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "failed to close file %q: %v\n", path, cerr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(valutes); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
