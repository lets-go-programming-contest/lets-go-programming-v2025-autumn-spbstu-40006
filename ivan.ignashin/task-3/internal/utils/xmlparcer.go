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

type Record struct {
	ID    int     `xml:"id"`
	Name  string  `xml:"name"`
	Value float64 `xml:"value"`
}

type Records struct {
	Items []struct {
		ID    int    `xml:"NumCode"`
		Name  string `xml:"CharCode"`
		Value string `xml:"Value"`
	} `xml:"Valute"`
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
		value, err := strconv.ParseFloat(strings.ReplaceAll(item.Value, ",", "."), 64)
		if err != nil {
			return nil, fmt.Errorf("parse float %s: %w", item.Value, err)
		}

		records = append(records, Record{
			ID:    item.ID,
			Name:  item.Name,
			Value: value,
		})
	}

	return records, nil
}
