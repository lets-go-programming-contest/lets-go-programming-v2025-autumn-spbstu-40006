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
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value"     xml:"Value"`
}

func (v *Valute) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	type Alias Valute

	aux := struct {
		Value string `xml:"Value"`
		*Alias
	}{
		Value: "",
		Alias: (*Alias)(v),
	}

	if err := decoder.DecodeElement(&aux, &start); err != nil {
		return fmt.Errorf("failed to decode XML element: %w", err)
	}

	str := strings.ReplaceAll(aux.Value, ",", ".")
	val, err := strconv.ParseFloat(str, 64)

	if err != nil {
		return fmt.Errorf("invalid value %q: %w", aux.Value, err)
	}

	v.Value = val

	return nil
}

func ParseXMLFile(path string) ([]Valute, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open XML file: %w", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			panic("failed to close file")
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
			panic("failed to close file")
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(valutes); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
