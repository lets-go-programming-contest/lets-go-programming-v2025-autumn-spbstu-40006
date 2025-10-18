package app

import (
	"sort"
	"strconv"
	"strings"
)

func FillNSortRates(curs ValCurs) []Rate {
	rates := make([]Rate, 0, len(curs.Valute))

	for _, valute := range curs.Valute {
		s := strings.ReplaceAll(valute.ValueRaw, ",", ".")
		//nolint:wsl // конфликт двух линтеров: gofugmpt (требует без пустой строки) и wsl (требует пустую строку)
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			continue
		}

		if valute.Nominal > 0 {
			val /= float64(valute.Nominal)
		}

		rate := Rate{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    val,
		}

		rates = append(rates, rate)
	}

	sort.Slice(rates, func(i, j int) bool {
		return rates[i].Value > rates[j].Value
	})

	return rates
}
