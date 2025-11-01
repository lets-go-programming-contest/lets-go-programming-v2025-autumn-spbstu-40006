package rates

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type Currency struct {
	NumCode  int     `xml:"NumCode" json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    string  `xml:"Value" json:"-"`
	ValueNum float64 `json:"value"`
}

type valCurs struct {
	XMLName xml.Name   `xml:"ValCurs"`
	Valutes []Currency `xml:"Valute"`
}

func Load(path string) ([]Currency, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open input xml %q: %w", path, err)
	}
	defer f.Close()

	dec := xml.NewDecoder(f)
	dec.CharsetReader = func(label string, r io.Reader) (io.Reader, error) {
		return charset.NewReaderLabel(label, r)
	}

	var vc valCurs
	if err := dec.Decode(&vc); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	for i := range vc.Valutes {
		raw := strings.ReplaceAll(vc.Valutes[i].Value, ",", ".")
		val, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil {
			return nil, fmt.Errorf("parse value %q: %w", vc.Valutes[i].Value, err)
		}
		vc.Valutes[i].ValueNum = val
	}
	return vc.Valutes, nil
}
