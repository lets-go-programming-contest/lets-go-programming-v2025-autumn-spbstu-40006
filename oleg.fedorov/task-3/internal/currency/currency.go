package currency

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type XMLCurrency struct {
	XMLName  xml.Name `xml:"Valute"`
	ID       string   `xml:"ID,attr"`
	NumCode  int      `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Nominal  int      `xml:"Nominal"`
	Name     string   `xml:"Name"`
	Value    string   `xml:"Value"`
}

type JSONCurrency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

type XMLCurrencies struct {
	XMLName    xml.Name      `xml:"ValCurs"`
	Currencies []XMLCurrency `xml:"Valute"`
}

func (c *XMLCurrency) ToJSONCurrency() (JSONCurrency, error) {
	value := strings.Replace(c.Value, ",", ".", -1)

	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return JSONCurrency{}, err
	}

	return JSONCurrency{
		NumCode:  c.NumCode,
		CharCode: c.CharCode,
		Value:    parsedValue,
	}, nil
}
