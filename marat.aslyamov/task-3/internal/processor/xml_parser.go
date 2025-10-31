package processor

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
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	ValueStr string `xml:"Value"`
}

type RawRecords struct {
	Items []Record `xml:"Valute"`
}

func (cp *CurrencyProcessor) ParseXMLFile(path string) ([]Currency, error) {
	xmlData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read xml file: %w", err)
	}

	var raw RawRecords

	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	currencies := make([]Currency, len(raw.Items))

	for index, item := range raw.Items {
		valueStr := strings.ReplaceAll(item.ValueStr, ",", ".")
		value, err := strconv.ParseFloat(valueStr, 64)

		currencies[index] = Currency{
			NumCode:  item.NumCode,
			CharCode: item.CharCode,
			Value:    value,
		}

		if err != nil {
			return nil, fmt.Errorf("parse float: %w", err)
		}
	}

	return currencies, nil
}
