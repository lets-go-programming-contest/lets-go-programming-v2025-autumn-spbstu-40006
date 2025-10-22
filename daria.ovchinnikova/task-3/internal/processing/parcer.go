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
	NumCode  int     `xml:"NumCode" json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    float64 `xml:"Value" json:"value"`
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
	defer file.Close()

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

	var currencies []Currency
	for _, r := range raw.Items {
		f := parseValue(r.Value)
		c := Currency{
			NumCode:  r.NumCode,
			CharCode: r.CharCode,
			Value:    f,
		}
		currencies = append(currencies, c)
	}

	return currencies
}

func parseValue(s string) float64 {
	s = strings.Replace(s, ",", ".", 1)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}
