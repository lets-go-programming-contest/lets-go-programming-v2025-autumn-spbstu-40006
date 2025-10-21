package rates

type ExchangeRate struct {
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"  yaml:"num_code"`
	CharCode string  `json:"char_code" xml:"CharCode" yaml:"char_code"`
	Value    Decimal `json:"value"     xml:"Value"    yaml:"value"`
}

func CompareByValueDesc(a, b Currency) int {
	floatA, floatB := float64(a.Value), float64(b.Value)

	switch {
	case floatB < floatA:
		return -1
	case floatB > floatA:
		return 1
	default:
		return 0
	}
}
