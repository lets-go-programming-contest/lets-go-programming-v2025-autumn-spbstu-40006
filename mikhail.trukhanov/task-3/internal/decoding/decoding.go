package decoding

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/Mishaa105/task-3/internal/config"
	"golang.org/x/net/html/charset"
)

type Valute struct {
	NumCode  int        `json:"num_code"  xml:"NumCode"`
	CharCode string     `json:"char_code" xml:"CharCode"`
	Value    FloatValue `json:"value"     xml:"Value"`
}

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

type FloatValue float64

func (f *FloatValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var valueStr string
	if err := d.DecodeElement(&valueStr, &start); err != nil {
		return fmt.Errorf("decode element failed: %w", err)
	}

	valueStr = strings.Replace(valueStr, ",", ".", 1)

	val, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fmt.Errorf("cannot parse Value: %w", err)
	}

	*f = FloatValue(val)

	return nil
}

func Decoding(configPath string) (ValCurs, error) {
	cfg, err := config.CheckInput(configPath)
	if err != nil {
		return ValCurs{}, fmt.Errorf("checking config: %w", err)
	}

	xmlFile, err := os.Open(cfg.InputFile)
	if err != nil {
		return ValCurs{}, fmt.Errorf("opening input file: %w", err)
	}

	defer func() {
		if err := xmlFile.Close(); err != nil {
			fmt.Printf("failed to close file: %v\n", err)
		}
	}()

	var valCurs ValCurs

	decoder := xml.NewDecoder(xmlFile)
	decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(&valCurs)
	if err != nil {
		return ValCurs{}, fmt.Errorf("decoding XML failed: %w", err)
	}

	sort.Slice(valCurs.Valutes, func(i, j int) bool {
		return valCurs.Valutes[i].Value > valCurs.Valutes[j].Value
	})

	return valCurs, nil
}
