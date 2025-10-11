package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type Valute struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

type Record struct {
	ID    int
	Name  string
	Value float64
}

func ParseXML(path string) ([]Record, error) {
	xmlData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read xml file %s: %w", path, err)
	}

	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	records := make([]Record, 0, len(valCurs.Valutes))
	for _, v := range valCurs.Valutes {
		value, err := strconv.ParseFloat(strings.ReplaceAll(v.Value, ",", "."), 64)
		if err != nil {
			return nil, fmt.Errorf("parse float %s: %w", v.Value, err)
		}

		records = append(records, Record{
			ID:    v.NumCode,
			Name:  v.CharCode,
			Value: value,
		})
	}

	return records, nil
}
