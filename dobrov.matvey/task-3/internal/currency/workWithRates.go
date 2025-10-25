package currency

import (
	"sort"
)

func FillNSortRates(currencies *ValCurs) []Currency {
	rates := make([]Currency, 0, len(currencies.Currencies))

	for _, currency := range currencies.Currencies {
		object := Currency{
			NumCode:  currency.NumCode,
			CharCode: currency.CharCode,
			Value:    currency.Value,
		}
		rates = append(rates, object)
	}

	sort.Slice(rates, func(i, j int) bool {
		return rates[i].Value > rates[j].Value
	})

	return rates
}
