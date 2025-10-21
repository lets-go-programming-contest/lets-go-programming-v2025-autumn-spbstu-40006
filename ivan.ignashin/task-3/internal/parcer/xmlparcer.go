package parcer

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
	ID    int     `json:"num_code"  xml:"NumCode"`
	Name  string  `json:"char_code" xml:"CharCode"`
	Value float64 `json:"value"     xml:"Value"`
}

type RawRecords struct {
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

	var raw RawRecords

	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	records := make([]Record, 0, len(raw.Items))

	for _, item := range raw.Items {
		valueStr := strings.ReplaceAll(item.Value, ",", ".")

		valueFloat, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("parse float %q: %w", item.Value, err)
		}

		records = append(records, Record{
			ID:    item.ID,
			Name:  item.Name,
			Value: valueFloat,
		})
	}

	return records, nil
}
