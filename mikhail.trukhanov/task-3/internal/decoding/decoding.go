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
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value"     xml:"Value"`
}

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

func (v *Valute) UnmarshalXML(decode *xml.Decoder, start xml.StartElement) error {
	var temp struct {
		NumCode  int    `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	if err := decode.DecodeElement(&temp, &start); err != nil {
		return fmt.Errorf("decode element failed: %w", err)
	}

	v.NumCode = temp.NumCode
	v.CharCode = temp.CharCode

	valueStr := strings.Replace(temp.Value, ",", ".", 1)

	val, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fmt.Errorf("cannot parse Value: %w", err)
	}

	v.Value = val

	return nil
}

func Decoding(configPath string) ValCurs {
	cfg, err := config.CheckInput(configPath)
	if err != nil {
		panic(err)
	}

	xmlFile, err := os.Open(cfg.InputFile)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	sort.Slice(valCurs.Valutes, func(i, j int) bool {
		return valCurs.Valutes[i].Value > valCurs.Valutes[j].Value
	})

	return valCurs
}
