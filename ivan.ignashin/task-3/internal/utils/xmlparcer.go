package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html/charset"
)

type Record struct {
	ID    int    `xml:"NumCode"`
	Name  string `xml:"CharCode"`
	Value string `xml:"Value"`
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
		value := strings.Replace(item.Value, ",", ".", -1)

		records = append(records, Record{
			ID:    item.ID,
			Name:  item.Name,
			Value: value,
		})
	}

	return records, nil
}
