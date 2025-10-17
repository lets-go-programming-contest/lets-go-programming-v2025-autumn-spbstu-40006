package utils

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

const dirPerm = 0o755

type ValCurs struct {
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int     `xml:"NumCode"`
	CharCode string  `xml:"CharCode"`
	ValueStr string  `xml:"Value"`
	Value    float64 `xml:"-"`
}

type Result struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
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
	err = decoder.Decode(&valCurs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	for index := range valCurs.Valutes {
		str := strings.ReplaceAll(valCurs.Valutes[index].ValueStr, ",", ".")
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert value %q to float: %w", str, err)
		}

		valCurs.Valutes[index].Value = val
	}

	return valCurs.Valutes, nil
}

func SaveToJSON(path string, valutes []Valute) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", dir, err)
	}

	results := make([]Result, len(valutes))
	for idx, val := range valutes {
		results[idx] = Result{
			NumCode:  val.NumCode,
			CharCode: val.CharCode,
			Value:    val.Value,
		}
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

	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
