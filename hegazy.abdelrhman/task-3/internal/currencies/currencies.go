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
	var valueString string
	if err := d.DecodeElement(&valueString, &start); err != nil {
		return fmt.Errorf("currencies: %w", err)
	}

	valueString = strings.ReplaceAll(valueString, ",", ".")

	val, err := strconv.ParseFloat(valueString, 32)
	if err != nil {
		return fmt.Errorf("currencies: %w", err)
	}

	*f = Float(val)

	return nil
}

type Currency struct {
	Name  string  `json:"name"      xml:"Name"`
	Code  int     `json:"curr_code" xml:"CurrencyCode"`
	Value float64 `json:"value"     xml:"Value"`
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

	currencyData := &Currencies{
		Currencies: []Currency{},
	}
	if err = decoder.Decode(currencyData); err != nil {
		return nil, fmt.Errorf("currencies: %w", err)
	}

	return currencyData, nil
}

func (c *Currencies) SaveToOutputFile(path string) error {
	const directoryPermissions = 0o755
	if err := os.MkdirAll(filepath.Dir(path), directoryPermissions); err != nil {
		return fmt.Errorf("currencies: %w", err)
	}

	outputFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("currencies: %w", err)
	}

	defer func() {
		if closeErr := outputFile.Close(); closeErr != nil {
			// Log the error but don't return it since we're in a defer
			fmt.Printf("warning: failed to close file: %v\n", closeErr)
		}
	}()

	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ")

	if err = encoder.Encode(c.Currencies); err != nil {
		return fmt.Errorf("currencies: %w", err)
	}

	return nil
}
