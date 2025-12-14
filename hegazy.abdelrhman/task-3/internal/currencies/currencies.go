package currencies

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

type Float float32

func (f *Float) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return fmt.Errorf("currencies: %w", err)
	}

	s = strings.ReplaceAll(s, ",", ".")

	val, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return fmt.Errorf("currencies: %w", err)
	}

	*f = Float(val)

	return nil
}

type Currency struct {
	NumCode  int    `json:"num_code" xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    Float  `json:"value" xml:"Value"`
}

type Currencies struct {
	Currencies []Currency `xml:"Valute"`
}

func New(path string) (*Currencies, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("currencies: %w", err)
	}

	decoder := xml.NewDecoder(strings.NewReader(string(data)))
	decoder.CharsetReader = charset.NewReaderLabel

	c := &Currencies{
		Currencies: []Currency{},
	}
	if err = decoder.Decode(c); err != nil {
		return nil, fmt.Errorf("currencies: %w", err)
	}

	return c, nil
}

func (c *Currencies) SaveToOutputFile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("currencies: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("currencies: %w", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	if err = encoder.Encode(c.Currencies); err != nil {
		return fmt.Errorf("currencies: %w", err)
	}

	return nil
}
