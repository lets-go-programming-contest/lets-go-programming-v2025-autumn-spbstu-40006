package utils

type ExchangeRate struct {
	Currencies []Currency `json:"Valute" xml:"Valute"`
}

type Currency struct {
	NumCode  string     `json:"num_code"  xml:"NumCode"`
	CharCode string     `json:"char_code" xml:"CharCode"`
	Value    CommaFloat `json:"value"     xml:"Value"`
}
