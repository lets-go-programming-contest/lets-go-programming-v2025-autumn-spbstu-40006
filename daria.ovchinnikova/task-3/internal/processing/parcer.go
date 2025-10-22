package processing

import (
	"encoding/xml"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type Currency struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value"     xml:"Value"`
}

type RawCurrencies struct {
	Items []struct {
		NumCode  int    `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	} `xml:"Valute"`
}

func LoadXML(path string) []Currency {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch strings.ToLower(charset) {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return input, nil
		}
	}

	var raw RawCurrencies
	if err := decoder.Decode(&raw); err != nil {
		panic(err)
	}

	currencies := make([]Currency, 0, len(raw.Items))

	for _, item := range raw.Items {
		valueFloat := parseValue(item.Value)
		currencies = append(currencies, Currency{
			NumCode:  item.NumCode,
			CharCode: item.CharCode,
			Value:    valueFloat,
		})
	}

	return currencies
}

func parseValue(s string) float64 {
	valueFloat, err := strconv.ParseFloat(strings.Replace(s, ",", ".", 1), 64)
	if err != nil {
		panic(err)
	}

	return valueFloat
}
