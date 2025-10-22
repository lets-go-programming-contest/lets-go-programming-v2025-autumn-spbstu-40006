package processing

import (
	"encoding/xml"
	"os"
	"strconv"
	"strings"
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
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var raw RawCurrencies
	if err := xml.Unmarshal(data, &raw); err != nil {
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
