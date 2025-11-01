package rates

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type Currency struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    string  `json:"-"         xml:"Value"`
	ValueNum float64 `json:"value"     xml:"-"`
}

type valCurs struct {
	XMLName xml.Name   `xml:"ValCurs"`
	Valutes []Currency `xml:"Valute"`
}

func Load(path string) ([]Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open input xml %q: %w", path, err)
	}

	defer func() {
		_ = file.Close()
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var curs valCurs
	if err := decoder.Decode(&curs); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	for idx := range curs.Valutes {
		raw := strings.ReplaceAll(curs.Valutes[idx].Value, ",", ".")
		raw = strings.TrimSpace(raw)

		val, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil {
			return nil, fmt.Errorf("parse value %q: %w", curs.Valutes[idx].Value, err)
		}

		curs.Valutes[idx].ValueNum = val
	}

	return curs.Valutes, nil
}
