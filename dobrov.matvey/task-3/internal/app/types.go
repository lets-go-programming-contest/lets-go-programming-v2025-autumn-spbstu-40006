package app

type ValCurs struct {
	Valute []Valute `xml:"Valute"`
}
type Valute struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	ValueRaw string `xml:"Value"`
}

type Rate struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}
