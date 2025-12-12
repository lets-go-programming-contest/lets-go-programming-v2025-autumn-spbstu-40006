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

func (float *Float) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	decodeElementString := ""
	if err := d.DecodeElement(&decodeElementString, &start); err != nil {
		return fmt.Errorf("currencies: %w", err)
	}

	decodeElementString = strings.ReplaceAll(decodeElementString, ",", ".")

	f, err := strconv.ParseFloat(decodeElementString, 32)
	if err != nil {
		return fmt.Errorf("currencies: %w", err)
	}

	*float = Float(f)

	return nil
}

type Currency struct {
	NumCode  int    `json:"num_code"  xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    Float  `json:"value"     xml:"Value"`
}

type Currencies struct {
	Currencies []Currency `xml:"Valute"`
}

func New(path string) (*Currencies, error) {
	currenciesContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("currencies: %w", err)
	}

	decoder := xml.NewDecoder(strings.NewReader(string(currenciesContent)))
	decoder.CharsetReader = charset.NewReaderLabel

	currencies := &Currencies{
		Currencies: []Currency{},
	}
	if err = decoder.Decode(currencies); err != nil {
		return nil, fmt.Errorf("currencies: %w", err)
	}

	return currencies, nil
}

func (currencies *Currencies) SaveToOutputFile(path string) error {
	dirs := filepath.Dir(path)

	err := os.MkdirAll(dirs, 0o755) //nolint:mnd
	if err != nil {
		return fmt.Errorf("currencies: %w", err)  // Fixed: return error instead of printing
	}

	output, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o600) //nolint:mnd
	if err != nil {
		return fmt.Errorf("currencies: %w", err)  // Fixed: return error instead of printing
	}
	defer output.Close()  // Added: important to close the file

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")

	if err = encoder.Encode(currencies.Currencies); err != nil {
		return fmt.Errorf("currencies: %w", err)  // Fixed: return error instead of printing
	}

	return nil
}