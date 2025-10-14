package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

type Record struct {
	ID    int     `xml:"numCode"`
	Name  string  `xml:"CharCode"`
	Value float64 `xml:"Value"`
}

type Records struct {
	Items []Record `xml:"Valute"`
}

func ParseXML(path string) ([]Record, error) {
	xmlData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read xml file %s: %w", path, err)
	}

	var rawRecords Records

	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&rawRecords); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	records := make([]Record, 0, len(rawRecords.Items))

	for _, item := range rawRecords.Items {
		records = append(records, Record{
			ID:    item.ID,
			Name:  item.Name,
			Value: item.Value,
		})
	}

	return records, nil
}
