package currency

type CurrencyValue float64

type ValCurs struct {
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	NumCode  int           `json:"num_code"  xml:"NumCode"`
	CharCode string        `json:"char_code" xml:"CharCode"`
	Value    CurrencyValue `json:"value"     xml:"Value"`
}
