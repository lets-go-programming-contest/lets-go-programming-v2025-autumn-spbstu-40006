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
	NumCodeStr string `xml:"NumCode"`
	CharCode   string `xml:"CharCode"`
	Value      string `xml:"Value"`
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
	for _, valute := range valCurs.Valutes {
		num, err := strconv.Atoi(strings.TrimSpace(valute.NumCodeStr))
		if err != nil {
			return nil, fmt.Errorf("parse NumCode %q: %w", valute.NumCodeStr, err)
		}
		value, err := strconv.ParseFloat(strings.ReplaceAll(valute.Value, ",", "."), 64)
		if err != nil {
			return nil, fmt.Errorf("parse Value %q: %w", valute.Value, err)
		}

		records = append(records, Record{
			ID:    num,
			Name:  valute.CharCode,
			Value: value,
		})
	}

	return records, nil
}
